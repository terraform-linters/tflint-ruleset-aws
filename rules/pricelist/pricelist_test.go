package pricelist

import (
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

const offerMetadata = `"FormatVersion","v1.0"
"Disclaimer","This pricing list is for informational purposes only."
"Publication Date","2026-07-29T23:42:48Z"
"Version","20260729234248"
"OfferCode","AmazonRDS"
`

const offerHeader = `"SKU","TermType","Product Family","Instance Type","Current Generation"
`

func TestParseValues(t *testing.T) {
	for _, tc := range []struct {
		name     string
		offer    string
		column   string
		keep     func(Product) bool
		expected []string
		wantErr  bool
	}{
		{
			name: "distinct values below the metadata rows",
			offer: offerMetadata + offerHeader +
				`"ABC","OnDemand","Database Instance","db.m6g.large","Yes"
"DEF","Reserved","Database Instance","db.m6g.large","Yes"
"GHI","OnDemand","Database Instance","db.t4g.micro","Yes"
`,
			column:   "Instance Type",
			expected: []string{"db.m6g.large", "db.t4g.micro"},
		},
		{
			name: "products with an empty value are skipped",
			offer: offerMetadata + offerHeader +
				`"ABC","OnDemand","Database Storage","","No"
"DEF","OnDemand","Database Instance","db.m6g.large","Yes"
`,
			column:   "Instance Type",
			expected: []string{"db.m6g.large"},
		},
		{
			name: "keep filters on another column",
			offer: offerMetadata + offerHeader +
				`"ABC","OnDemand","Database Instance","db.m1.small","No"
"DEF","OnDemand","Database Instance","db.m6g.large","Yes"
`,
			column:   "Instance Type",
			keep:     func(p Product) bool { return p.Get("Current Generation") == "No" },
			expected: []string{"db.m1.small"},
		},
		{
			name: "keep reading an absent column",
			offer: offerMetadata + offerHeader +
				`"ABC","OnDemand","Database Instance","db.m6g.large","Yes"
`,
			column:   "Instance Type",
			keep:     func(p Product) bool { return p.Get("Database Engine") == "" },
			expected: []string{"db.m6g.large"},
		},
		{
			name:     "no products below the header",
			offer:    offerMetadata + offerHeader,
			column:   "Instance Type",
			expected: []string{},
		},
		{
			name:    "no header",
			offer:   offerMetadata,
			column:  "Instance Type",
			wantErr: true,
		},
		{
			name:    "column missing from the header",
			offer:   offerMetadata + offerHeader,
			column:  "Cache Engine",
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			values, err := parseValues(strings.NewReader(tc.offer), tc.column, tc.keep)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected an error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if got := slices.Sorted(maps.Keys(values)); !slices.Equal(got, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestValues(t *testing.T) {
	offers := map[string]string{
		"/offers/v1.0/aws/20260729/us-east-1/index.csv": offerMetadata + offerHeader +
			`"ABC","OnDemand","Database Instance","db.m6g.large","Yes"
"DEF","OnDemand","Database Instance","db.m1.small","No"
`,
		"/offers/v1.0/aws/20260729/eu-west-1/index.csv": offerMetadata + offerHeader +
			`"GHI","OnDemand","Database Instance","db.m6g.large","Yes"
"JKL","OnDemand","Database Instance","db.t4g.micro","Yes"
`,
	}

	for _, tc := range []struct {
		name     string
		regions  []string
		keep     func(Product) bool
		expected []string
		wantErr  bool
	}{
		{
			name:     "values merged across regions",
			regions:  []string{"us-east-1", "eu-west-1"},
			expected: []string{"db.m1.small", "db.m6g.large", "db.t4g.micro"},
		},
		{
			name:     "keep applied in every region",
			regions:  []string{"us-east-1", "eu-west-1"},
			keep:     func(p Product) bool { return p.Get("Current Generation") == "Yes" },
			expected: []string{"db.m6g.large", "db.t4g.micro"},
		},
		{
			name:    "unreachable region offer",
			regions: []string{"us-east-1", "ap-northeast-1"},
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/offers/v1.0/aws/AmazonRDS/current/region_index.json" {
					fmt.Fprint(w, regionIndex(tc.regions))
					return
				}

				offer, ok := offers[r.URL.Path]
				if !ok {
					http.NotFound(w, r)
					return
				}
				fmt.Fprint(w, offer)
			}))
			defer server.Close()

			previous := baseURL
			t.Cleanup(func() { baseURL = previous })
			baseURL = server.URL

			values, err := Values(RDS, "Instance Type", tc.keep)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected an error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !slices.Equal(values, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, values)
			}
		})
	}
}

func regionIndex(regions []string) string {
	var entries []string
	for _, region := range regions {
		entries = append(entries, fmt.Sprintf(
			`"%s": {"regionCode": "%s", "currentVersionUrl": "/offers/v1.0/aws/20260729/%s/index.json"}`,
			region, region, region,
		))
	}

	return fmt.Sprintf(`{"regions": {%s}}`, strings.Join(entries, ","))
}
