//go:build generators

package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

// formatEntry renders an entry for comparison, avoiding time.Time equality
// pitfalls around monotonic clock readings and location pointers.
func formatEntry(e runtimeEntry) string {
	optional := func(t *time.Time) string {
		if t == nil {
			return "none"
		}
		return t.Format("2006-01-02")
	}
	return fmt.Sprintf(
		"eos=%s create=%s update=%s",
		e.EndOfSupportDate.Format("2006-01-02"),
		optional(e.BlockCreateDate),
		optional(e.BlockUpdateDate),
	)
}

func TestParseDeprecationDate(t *testing.T) {
	for _, tc := range []struct {
		name      string
		input     string
		scheduled bool
		expected  string
	}{
		{
			name:      "full date",
			input:     "Jun 30, 2029",
			scheduled: true,
			expected:  "2029-06-30",
		},
		{
			name:      "single digit day",
			input:     "Mar 5, 2020",
			scheduled: true,
			expected:  "2020-03-05",
		},
		{
			name:      "surrounding whitespace",
			input:     "  Nov 10, 2026\n",
			scheduled: true,
			expected:  "2026-11-10",
		},
		{
			name:      "not scheduled",
			input:     "Not scheduled",
			scheduled: false,
		},
		{
			name:      "not applicable",
			input:     "N/A",
			scheduled: false,
		},
		{
			name:      "empty",
			input:     "",
			scheduled: false,
		},
		{
			name:      "en dash",
			input:     "–",
			scheduled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			date, scheduled := parseDeprecationDate(tc.input)
			if scheduled != tc.scheduled {
				t.Fatalf("expected scheduled %t, got %t", tc.scheduled, scheduled)
			}
			if !scheduled {
				return
			}
			if got := date.Format("2006-01-02"); got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

// An unrecognized deprecation date means the table format changed. Panicking
// fails the scheduled maintenance job so the change is noticed.
func TestParseDeprecationDateUnrecognized(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic for an unrecognized deprecation date")
		}
	}()

	parseDeprecationDate("30 June 2029")
}

func TestParseBlockDate(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "full date",
			input:    "Aug 31, 2029",
			expected: "2029-08-31",
		},
		{
			name:     "not scheduled",
			input:    "Not scheduled",
			expected: "none",
		},
		{
			name:     "not applicable",
			input:    "N/A",
			expected: "none",
		},
		{
			// A block date the generator cannot read costs the reported
			// message some specificity. The same value in the deprecation
			// column panics, since it would drop the runtime instead.
			name:     "unrecognized value",
			input:    "Aug 31, 2029 (estimated)",
			expected: "none",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			date := parseBlockDate(tc.input)
			got := "none"
			if date != nil {
				got = date.Format("2006-01-02")
			}
			if got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestParseRow(t *testing.T) {
	for _, tc := range []struct {
		name     string
		cells    []string
		ok       bool
		expected string
	}{
		{
			name:     "full lifecycle",
			cells:    []string{"Python 3.13", "python3.13", "Amazon Linux 2023", "Jun 30, 2029", "Jul 31, 2029", "Aug 31, 2029"},
			ok:       true,
			expected: "eos=2029-06-30 create=2029-07-31 update=2029-08-31",
		},
		{
			// dotnet7 and dotnet5.0 reached end of support without AWS ever
			// scheduling the blocks.
			name:     "end of support without blocks",
			cells:    []string{".NET 7", "dotnet7", "Amazon Linux 2", "May 14, 2024", "N/A", "N/A"},
			ok:       true,
			expected: "eos=2024-05-14 create=none update=none",
		},
		{
			name:     "blocks not yet scheduled",
			cells:    []string{".NET 9", "dotnet9", "Amazon Linux 2023", "Nov 10, 2026", "Not scheduled", "Not scheduled"},
			ok:       true,
			expected: "eos=2026-11-10 create=none update=none",
		},
		{
			// A supported runtime with no announced deprecation date. Emitting
			// it with a zero end of support date reported every function using
			// the runtime as deprecated.
			name:  "no deprecation date scheduled",
			cells: []string{"Python 3.15", "python3.15", "Amazon Linux 2023", "Not scheduled", "Not scheduled", "Not scheduled"},
			ok:    false,
		},
		{
			// An unreadable block column leaves that date absent and keeps
			// the runtime on the schedule.
			name:     "unrecognized block date",
			cells:    []string{"Python 3.13", "python3.13", "Amazon Linux 2023", "Jun 30, 2029", "Jul 31, 2029 (estimated)", "Aug 31, 2029"},
			ok:       true,
			expected: "eos=2029-06-30 create=none update=2029-08-31",
		},
		{
			name:  "header row",
			cells: []string{"Name", "Identifier", "Operating system", "Deprecation date", "Block function create", "Block function update"},
			ok:    false,
		},
		{
			name:  "unrecognized identifier",
			cells: []string{"Rust", "rust1.x", "Amazon Linux 2023", "Jun 30, 2029", "Jul 31, 2029", "Aug 31, 2029"},
			ok:    false,
		},
		{
			name:  "too few cells",
			cells: []string{"Python 3.13", "python3.13", "Amazon Linux 2023"},
			ok:    false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			identifier, entry, ok := parseRow(tc.cells)
			if ok != tc.ok {
				t.Fatalf("expected ok %t, got %t", tc.ok, ok)
			}
			if !ok {
				return
			}
			if identifier != tc.cells[1] {
				t.Errorf("expected identifier %s, got %s", tc.cells[1], identifier)
			}
			if got := formatEntry(entry); got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

const runtimesPage = `
<html><body>
<table>
  <tr>
    <th>Name</th><th>Identifier</th><th>Operating system</th>
    <th>Deprecation date</th><th>Block function create</th><th>Block function update</th>
  </tr>
  <tr>
    <td>Node.js 26</td><td>nodejs26.x</td><td>Amazon Linux 2023</td>
    <td>Not scheduled</td><td>Not scheduled</td><td>Not scheduled</td>
  </tr>
  <tr>
    <td>Python 3.15</td><td>python3.15</td><td>Amazon Linux 2023</td>
    <td>Not scheduled</td><td>Not scheduled</td><td>Not scheduled</td>
  </tr>
  <tr>
    <td>Python 3.13</td><td>python3.13</td><td>Amazon Linux 2023</td>
    <td>Jun 30, 2029</td><td>Jul 31, 2029</td><td>Aug 31, 2029</td>
  </tr>
  <tr>
    <td>.NET 7</td><td>dotnet7</td><td>Amazon Linux 2</td>
    <td>May 14, 2024</td><td>N/A</td><td>N/A</td>
  </tr>
</table>
</body></html>
`

// Runtimes AWS has not scheduled for deprecation are absent from the output.
// They were previously emitted with a zero end of support date, which the rule
// read as a date in the past, reporting every function using a newly released
// runtime as having reached end of support.
func TestParseRuntimes(t *testing.T) {
	entries := parseRuntimes(html.NewTokenizer(strings.NewReader(runtimesPage)))

	expected := map[string]string{
		"python3.13": "eos=2029-06-30 create=2029-07-31 update=2029-08-31",
		"dotnet7":    "eos=2024-05-14 create=none update=none",
	}

	if len(entries) != len(expected) {
		t.Fatalf("expected %d runtimes, got %d: %v", len(expected), len(entries), entries)
	}

	for identifier, want := range expected {
		entry, ok := entries[identifier]
		if !ok {
			t.Errorf("expected runtime %s to be present", identifier)
			continue
		}
		if got := formatEntry(entry); got != want {
			t.Errorf("%s: expected %s, got %s", identifier, want, got)
		}
	}
}
