// Package pricelist reads the AWS Price List bulk API, which publishes an offer
// file per service and region listing every product AWS sells. Generators use
// it as a source of truth for values that AWS does not publish in its API
// models, such as RDS DB instance classes.
package pricelist

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
)

// Offer identifies a service's price list.
type Offer string

const (
	RDS         Offer = "AmazonRDS"
	ElastiCache Offer = "AmazonElastiCache"
)

// Each region offer file is tens of megabytes, so they are fetched concurrently
// and parsed as a stream rather than buffered.
const fetchConcurrency = 8

// priceListHeader is the first column of the header row that separates the
// metadata rows at the top of an offer file from the products below it.
const priceListHeader = "SKU"

// valueSeparator joins a product's values into a key that identifies it. It
// cannot appear in an offer file, which holds only printable text.
const valueSeparator = "\x00"

var baseURL = "https://pricing.us-east-1.amazonaws.com"

// client bounds the wait for a response to fail a stalled connection rather
// than hang the weekly generation, while leaving the body itself untimed
// because the larger offers stream for a while.
var client = &http.Client{
	Transport: &http.Transport{
		// A custom transport starts with no proxy at all, unlike the default one.
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: 30 * time.Second,
		MaxIdleConnsPerHost:   fetchConcurrency,
	},
}

// Product is a distinct combination of column values, such as an instance type
// paired with the generation it belongs to.
type Product struct {
	columns map[string]int
	values  []string
}

// Get returns the value of the named column, or an empty string when the
// product was not read with that column.
func (p Product) Get(column string) string {
	index, ok := p.columns[column]
	if !ok {
		return ""
	}

	return p.values[index]
}

// Distinct returns the distinct combinations of the named columns across every
// region of the aws partition, sorted. Products whose values are all empty are
// skipped.
func Distinct(offer Offer, columns ...string) ([]Product, error) {
	if len(columns) == 0 {
		return nil, errors.New("no columns requested")
	}

	urls, err := offerURLs(offer)
	if err != nil {
		return nil, err
	}

	var (
		mutex  sync.Mutex
		all    = map[string][]string{}
		failed []error
	)

	queue := make(chan string)
	var workers sync.WaitGroup
	for range fetchConcurrency {
		workers.Go(func() {
			for url := range queue {
				products, err := offerProducts(url, columns)

				mutex.Lock()
				if err != nil {
					failed = append(failed, fmt.Errorf("reading %s: %w", url, err))
				} else {
					maps.Copy(all, products)
				}
				mutex.Unlock()
			}
		})
	}

	for _, url := range urls {
		queue <- url
	}
	close(queue)
	workers.Wait()

	if len(failed) > 0 {
		return nil, errors.Join(failed...)
	}

	index := columnIndex(columns)
	products := make([]Product, 0, len(all))
	for _, key := range slices.Sorted(maps.Keys(all)) {
		products = append(products, Product{columns: index, values: all[key]})
	}

	return products, nil
}

func offerURLs(offer Offer) ([]string, error) {
	body, err := get(fmt.Sprintf("%s/offers/v1.0/aws/%s/current/region_index.json", baseURL, offer))
	if err != nil {
		return nil, fmt.Errorf("fetching %s region index: %w", offer, err)
	}
	defer body.Close()

	var index struct {
		Regions map[string]struct {
			CurrentVersionURL string `json:"currentVersionUrl"`
		} `json:"regions"`
	}
	if err := json.NewDecoder(body).Decode(&index); err != nil {
		return nil, fmt.Errorf("decoding %s region index: %w", offer, err)
	}
	if len(index.Regions) == 0 {
		return nil, fmt.Errorf("no regions in %s region index", offer)
	}

	urls := make([]string, 0, len(index.Regions))
	for _, region := range index.Regions {
		urls = append(urls, baseURL+strings.TrimSuffix(region.CurrentVersionURL, "index.json")+"index.csv")
	}

	return urls, nil
}

func offerProducts(url string, columns []string) (map[string][]string, error) {
	body, err := get(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseProducts(body, columns)
}

// parseProducts reads an offer file, which opens with metadata rows before the
// header row of the price list itself. Products are keyed by their joined
// values so that a region's rows, which repeat per price, collapse to a set.
func parseProducts(offer io.Reader, columns []string) (map[string][]string, error) {
	reader := csv.NewReader(offer)
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true

	var indexes []int
	products := map[string][]string{}
	values := make([]string, len(columns))

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		if indexes == nil {
			if len(record) == 0 || record[0] != priceListHeader {
				continue
			}

			if indexes, err = headerIndexes(record, columns); err != nil {
				return nil, err
			}
			continue
		}

		empty := true
		for position, index := range indexes {
			if index < len(record) {
				values[position] = record[index]
			} else {
				values[position] = ""
			}
			empty = empty && values[position] == ""
		}

		if !empty {
			// values is reused across rows, so a newly seen product keeps a copy.
			if key := strings.Join(values, valueSeparator); products[key] == nil {
				products[key] = slices.Clone(values)
			}
		}
	}

	if indexes == nil {
		return nil, errors.New("no price list header found")
	}

	return products, nil
}

func headerIndexes(header []string, columns []string) ([]int, error) {
	indexes := make([]int, len(columns))
	for position, column := range columns {
		index := slices.Index(header, column)
		if index < 0 {
			return nil, fmt.Errorf("no %q column in price list", column)
		}
		indexes[position] = index
	}

	return indexes, nil
}

func columnIndex(columns []string) map[string]int {
	index := make(map[string]int, len(columns))
	for position, column := range columns {
		index[column] = position
	}

	return index
}

func get(url string) (io.ReadCloser, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return resp.Body, nil
}
