//go:build generators

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
)

const (
	pricingHost    = "https://pricing.us-east-1.amazonaws.com"
	regionIndexURL = pricingHost + "/offers/v1.0/aws/AmazonRDS/current/region_index.json"

	instanceTypeColumn  = "Instance Type"
	instanceClassPrefix = "db."

	// Each region offer file is tens of megabytes, so they are fetched
	// concurrently and parsed as a stream rather than buffered.
	fetchConcurrency = 8

	// Every region publishes hundreds of classes, so a result below this floor
	// indicates the offer file layout changed.
	minInstanceClasses = 300
)

// regionIndex is the AWS Price List index of per-region offer files.
type regionIndex struct {
	Regions map[string]struct {
		CurrentVersionURL string `json:"currentVersionUrl"`
	} `json:"regions"`
}

func main() {
	classes := fetchInstanceClasses(fetchRegionIndex())
	if len(classes) < minInstanceClasses {
		panic(fmt.Sprintf("found only %d instance classes, expected at least %d", len(classes), minInstanceClasses))
	}

	data, err := json.MarshalIndent(classes, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshaling JSON: %s", err))
	}

	if err := os.WriteFile("instance_classes.json", append(data, '\n'), 0644); err != nil {
		panic(fmt.Sprintf("writing JSON: %s", err))
	}

	fmt.Printf("Wrote %d instance classes to instance_classes.json\n", len(classes))
}

func fetchRegionIndex() regionIndex {
	body, err := get(regionIndexURL)
	if err != nil {
		panic(fmt.Sprintf("fetching region index: %s", err))
	}
	defer body.Close()

	var index regionIndex
	if err := json.NewDecoder(body).Decode(&index); err != nil {
		panic(fmt.Sprintf("decoding region index: %s", err))
	}
	if len(index.Regions) == 0 {
		panic("no regions found in region index")
	}

	return index
}

func fetchInstanceClasses(index regionIndex) []string {
	all := map[string]bool{}
	var mutex sync.Mutex

	urls := make(chan string)
	var workers sync.WaitGroup
	for range fetchConcurrency {
		workers.Add(1)
		go func() {
			defer workers.Done()

			for url := range urls {
				classes, err := fetchRegionInstanceClasses(url)
				if err != nil {
					panic(fmt.Sprintf("reading %s: %s", url, err))
				}

				mutex.Lock()
				maps.Copy(all, classes)
				mutex.Unlock()
			}
		}()
	}

	for _, region := range index.Regions {
		urls <- pricingHost + strings.TrimSuffix(region.CurrentVersionURL, "index.json") + "index.csv"
	}
	close(urls)
	workers.Wait()

	return slices.Sorted(maps.Keys(all))
}

func fetchRegionInstanceClasses(url string) (map[string]bool, error) {
	body, err := get(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return parseInstanceClasses(body)
}

// parseInstanceClasses collects the instance classes named by an offer file,
// which opens with metadata rows before the header row of the price list.
func parseInstanceClasses(offer io.Reader) (map[string]bool, error) {
	reader := csv.NewReader(offer)
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true

	classes := map[string]bool{}
	column := -1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if column < 0 {
			if len(record) == 0 || record[0] != "SKU" {
				continue
			}
			if column = slices.Index(record, instanceTypeColumn); column < 0 {
				return nil, fmt.Errorf("no %q column in price list", instanceTypeColumn)
			}
			continue
		}

		if column < len(record) && strings.HasPrefix(record[column], instanceClassPrefix) {
			classes[record[column]] = true
		}
	}

	if column < 0 {
		return nil, fmt.Errorf("no price list header found")
	}

	return classes, nil
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
