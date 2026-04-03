//go:build generators

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type output struct {
	UpdatedAt time.Time                `json:"updated_at"`
	Runtimes  map[string]runtimeEntry  `json:"runtimes"`
}

type runtimeEntry struct {
	EndOfSupportDate time.Time  `json:"end_of_support_date"`
	BlockCreateDate  *time.Time `json:"block_create_date,omitempty"`
	BlockUpdateDate  *time.Time `json:"block_update_date,omitempty"`
}

func main() {
	resp, err := http.Get("https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html")
	if err != nil {
		panic(fmt.Sprintf("fetching Lambda runtimes page: %s", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("unexpected status %d fetching Lambda runtimes page", resp.StatusCode))
	}

	entries := parseRuntimes(html.NewTokenizer(resp.Body))
	if len(entries) == 0 {
		panic("no runtimes found on page")
	}

	data, err := json.MarshalIndent(output{
		UpdatedAt: time.Now().UTC(),
		Runtimes:  entries,
	}, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshaling JSON: %s", err))
	}

	if err := os.WriteFile("deprecated_runtimes.json", append(data, '\n'), 0644); err != nil {
		panic(fmt.Sprintf("writing JSON: %s", err))
	}

	fmt.Printf("Wrote %d runtimes to deprecated_runtimes.json\n", len(entries))
}

func parseRuntimes(z *html.Tokenizer) map[string]runtimeEntry {
	entries := map[string]runtimeEntry{}
	var buf strings.Builder
	var cells []string
	isHeader := false
	const (
		scanning = iota
		inTable
		inCell
	)
	state := scanning

	for {
		switch z.Next() {
		case html.ErrorToken:
			return entries
		case html.StartTagToken:
			tn, _ := z.TagName()
			switch string(tn) {
			case "table":
				state = inTable
			case "tr":
				cells, isHeader = cells[:0], false
			case "th":
				isHeader = true
				fallthrough
			case "td":
				state = inCell
				buf.Reset()
			}
		case html.EndTagToken:
			tn, _ := z.TagName()
			switch string(tn) {
			case "th", "td":
				cells = append(cells, strings.TrimSpace(buf.String()))
				state = inTable
			case "tr":
				if !isHeader {
					if id, entry, ok := parseRow(cells); ok {
						entries[id] = entry
					}
				}
			case "table":
				state = scanning
			}
		case html.TextToken:
			if state == inCell {
				buf.Write(z.Text())
			}
		}
	}
}

func parseRow(cells []string) (string, runtimeEntry, bool) {
	if len(cells) < 6 {
		return "", runtimeEntry{}, false
	}

	identifier := cells[1]
	if !isRuntimeIdentifier(identifier) {
		return "", runtimeEntry{}, false
	}

	return identifier, runtimeEntry{
		EndOfSupportDate: parseDate(cells[3]),
		BlockCreateDate:  parseDatePtr(cells[4]),
		BlockUpdateDate:  parseDatePtr(cells[5]),
	}, true
}

func isRuntimeIdentifier(s string) bool {
	prefixes := []string{"nodejs", "python", "java", "dotnet", "ruby", "go1.", "provided"}
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func parseDate(s string) time.Time {
	t, _ := parseDateOpt(s)
	return t
}

func parseDatePtr(s string) *time.Time {
	t, ok := parseDateOpt(s)
	if !ok {
		return nil
	}
	return &t
}

func parseDateOpt(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" || s == "N/A" || s == "–" || strings.Contains(s, "Not scheduled") {
		return time.Time{}, false
	}
	t, err := time.Parse("Jan 2, 2006", s)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}
