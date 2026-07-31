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
)

// Offer identifies a service's price list.
type Offer string

const (
	ElastiCache Offer = "AmazonElastiCache"
	RDS         Offer = "AmazonRDS"
)

// Each region offer file is tens of megabytes, so they are fetched concurrently
// and parsed as a stream rather than buffered.
const fetchConcurrency = 8

// priceListHeader is the first column of the header row that separates the
// metadata rows at the top of an offer file from the products below it.
const priceListHeader = "SKU"

var baseURL = "https://pricing.us-east-1.amazonaws.com"

// Product is a single product of an offer file, such as an RDS DB instance
// running a particular engine in a particular region.
type Product struct {
	columns map[string]int
	record  []string
}

// Get returns the value of the named column, or an empty string when the offer
// file has no such column.
func (p Product) Get(column string) string {
	index, ok := p.columns[column]
	if !ok || index >= len(p.record) {
		return ""
	}

	return p.record[index]
}

// Values returns the distinct values of the given column, sorted, across every
// region of the aws partition. Products whose value is empty are skipped, as
// are those rejected by keep. A nil keep accepts every product.
//
// keep is called from multiple goroutines, so it must be safe for concurrent
// use.
func Values(offer Offer, column string, keep func(Product) bool) ([]string, error) {
	urls, err := offerURLs(offer)
	if err != nil {
		return nil, err
	}

	var (
		mutex  sync.Mutex
		all    = map[string]bool{}
		failed []error
	)

	queue := make(chan string)
	var workers sync.WaitGroup
	for range fetchConcurrency {
		workers.Add(1)
		go func() {
			defer workers.Done()

			for url := range queue {
				values, err := offerValues(url, column, keep)

				mutex.Lock()
				if err != nil {
					failed = append(failed, fmt.Errorf("reading %s: %w", url, err))
				} else {
					maps.Copy(all, values)
				}
				mutex.Unlock()
			}
		}()
	}

	for _, url := range urls {
		queue <- url
	}
	close(queue)
	workers.Wait()

	if len(failed) > 0 {
		return nil, errors.Join(failed...)
	}

	return slices.Sorted(maps.Keys(all)), nil
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

func offerValues(url, column string, keep func(Product) bool) (map[string]bool, error) {
	body, err := get(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseValues(body, column, keep)
}

// parseValues reads an offer file, which opens with metadata rows before the
// header row of the price list itself.
func parseValues(offer io.Reader, column string, keep func(Product) bool) (map[string]bool, error) {
	reader := csv.NewReader(offer)
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true

	var columns map[string]int
	values := map[string]bool{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if columns == nil {
			if len(record) == 0 || record[0] != priceListHeader {
				continue
			}

			columns = make(map[string]int, len(record))
			for index, name := range record {
				columns[name] = index
			}
			if _, ok := columns[column]; !ok {
				return nil, fmt.Errorf("no %q column in price list", column)
			}
			continue
		}

		product := Product{columns: columns, record: record}
		if value := product.Get(column); value != "" && (keep == nil || keep(product)) {
			values[value] = true
		}
	}

	if columns == nil {
		return nil, errors.New("no price list header found")
	}

	return values, nil
}

func get(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return resp.Body, nil
}
