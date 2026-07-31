//go:build generators

package main

import (
	"maps"
	"slices"
	"strings"
	"testing"
)

func TestParseInstanceClasses(t *testing.T) {
	const preamble = `"FormatVersion","v1.0"
"Disclaimer","This pricing list is for informational purposes only."
"Publication Date","2026-07-29T23:42:48Z"
"Version","20260729234248"
"OfferCode","AmazonRDS"
`
	const header = `"SKU","TermType","Product Family","Location","Instance Type","Current Generation"
`

	for _, tc := range []struct {
		name     string
		offer    string
		expected []string
		wantErr  bool
	}{
		{
			name: "instance classes across product families",
			offer: preamble + header +
				`"ABC","OnDemand","Database Instance","US East (N. Virginia)","db.m6g.large","Yes"
"DEF","OnDemand","Database Instance","US East (N. Virginia)","db.m6g.large","Yes"
"GHI","OnDemand","Database Storage","US East (N. Virginia)","","No"
"JKL","OnDemand","Database Instance","US East (N. Virginia)","db.t4g.micro","Yes"
`,
			expected: []string{"db.m6g.large", "db.t4g.micro"},
		},
		{
			name:     "no rows below the header",
			offer:    preamble + header,
			expected: []string{},
		},
		{
			name:    "no header",
			offer:   preamble,
			wantErr: true,
		},
		{
			name:    "header without an instance type column",
			offer:   preamble + `"SKU","TermType","Product Family"` + "\n",
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			classes, err := parseInstanceClasses(strings.NewReader(tc.offer))
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected an error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got := slices.Sorted(maps.Keys(classes))
			if !slices.Equal(got, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
