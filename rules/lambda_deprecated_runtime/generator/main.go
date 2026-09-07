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
	UpdatedAt time.Time               `json:"updated_at"`
	Runtimes  map[string]runtimeEntry `json:"runtimes"`
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

	// A runtime AWS has not scheduled for deprecation has no lifecycle to
	// record. Emitting it with a zero end of support date would report every
	// function using the runtime as having reached end of support.
	endOfSupport, scheduled := parseDeprecationDate(cells[3])
	if !scheduled {
		return "", runtimeEntry{}, false
	}

	return identifier, runtimeEntry{
		EndOfSupportDate: endOfSupport,
		BlockCreateDate:  parseBlockDate(cells[4]),
		BlockUpdateDate:  parseBlockDate(cells[5]),
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

// parseDateCell parses a date cell from the runtimes table. The values AWS
// uses to mean "no date announced" report scheduled=false. An unrecognized
// value means the table format changed and returns an error, leaving each
// column to decide what a parse failure there costs.
func parseDateCell(s string) (time.Time, bool, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "N/A" || s == "–" || strings.Contains(s, "Not scheduled") {
		return time.Time{}, false, nil
	}
	t, err := time.Parse("Jan 2, 2006", s)
	if err != nil {
		return time.Time{}, false, err
	}
	return t, true, nil
}

// parseDeprecationDate parses the column the whole schedule is keyed on, so an
// unrecognized value panics rather than silently dropping the runtime.
func parseDeprecationDate(s string) (time.Time, bool) {
	date, scheduled, err := parseDateCell(s)
	if err != nil {
		panic(fmt.Sprintf("parsing deprecation date %q: %s", s, err))
	}
	return date, scheduled
}

// parseBlockDate parses one of the optional block columns. An unrecognized
// value is absent rather than fatal, costing the reported message some
// specificity where the same failure in the deprecation column would drop the
// runtime entirely.
func parseBlockDate(s string) *time.Time {
	date, scheduled, _ := parseDateCell(s)
	if !scheduled {
		return nil
	}
	return &date
}
