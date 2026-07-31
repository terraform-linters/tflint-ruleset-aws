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

func TestParseProducts(t *testing.T) {
	for _, tc := range []struct {
		name     string
		offer    string
		columns  []string
		expected []string
		wantErr  bool
	}{
		{
			name: "rows repeated per price collapse to a set",
			offer: offerMetadata + offerHeader +
				`"ABC","OnDemand","Database Instance","db.m6g.large","Yes"
"DEF","Reserved","Database Instance","db.m6g.large","Yes"
"GHI","OnDemand","Database Instance","db.t4g.micro","Yes"
`,
			columns:  []string{"Instance Type"},
			expected: []string{"db.m6g.large", "db.t4g.micro"},
		},
		{
			name: "combinations of several columns",
			offer: offerMetadata + offerHeader +
				`"ABC","OnDemand","Database Instance","db.m1.small","No"
"DEF","OnDemand","Database Instance","db.m6g.large","Yes"
`,
			columns:  []string{"Instance Type", "Current Generation"},
			expected: []string{"db.m1.small\x00No", "db.m6g.large\x00Yes"},
		},
		{
			name: "products with no values at all are skipped",
			offer: offerMetadata + offerHeader +
				`"ABC","OnDemand","Database Storage","",""
"DEF","OnDemand","Database Instance","db.m6g.large","Yes"
`,
			columns:  []string{"Instance Type", "Current Generation"},
			expected: []string{"db.m6g.large\x00Yes"},
		},
		{
			name: "short rows pad missing columns",
			offer: offerMetadata + offerHeader +
				`"ABC","OnDemand","Database Instance"
"DEF","OnDemand","Database Instance","db.m6g.large","Yes"
`,
			columns:  []string{"Instance Type", "Current Generation"},
			expected: []string{"db.m6g.large\x00Yes"},
		},
		{
			name:     "no products below the header",
			offer:    offerMetadata + offerHeader,
			columns:  []string{"Instance Type"},
			expected: []string{},
		},
		{
			name:    "no header",
			offer:   offerMetadata,
			columns: []string{"Instance Type"},
			wantErr: true,
		},
		{
			name:    "column missing from the header",
			offer:   offerMetadata + offerHeader,
			columns: []string{"Instance Type", "Cache Engine"},
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			products, err := parseProducts(strings.NewReader(tc.offer), tc.columns)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected an error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if got := slices.Sorted(maps.Keys(products)); !slices.Equal(got, tc.expected) {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestDistinct(t *testing.T) {
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
		columns  []string
		expected [][]string
		wantErr  bool
	}{
		{
			name:    "products merged across regions",
			regions: []string{"us-east-1", "eu-west-1"},
			columns: []string{"Instance Type", "Current Generation"},
			expected: [][]string{
				{"db.m1.small", "No"},
				{"db.m6g.large", "Yes"},
				{"db.t4g.micro", "Yes"},
			},
		},
		{
			name:    "unreachable region offer",
			regions: []string{"us-east-1", "ap-northeast-1"},
			columns: []string{"Instance Type"},
			wantErr: true,
		},
		{
			name:    "no columns requested",
			regions: []string{"us-east-1"},
			columns: nil,
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

			products, err := Distinct(RDS, tc.columns...)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected an error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got := make([][]string, 0, len(products))
			for _, product := range products {
				got = append(got, []string{product.Get("Instance Type"), product.Get("Current Generation")})
			}
			if !slices.EqualFunc(got, tc.expected, slices.Equal) {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestProductGet(t *testing.T) {
	product := Product{columns: columnIndex([]string{"Instance Type"}), values: []string{"db.m6g.large"}}

	for _, tc := range []struct {
		name     string
		column   string
		expected string
	}{
		{
			name:     "requested column",
			column:   "Instance Type",
			expected: "db.m6g.large",
		},
		{
			name:     "column the product was not read with",
			column:   "Current Generation",
			expected: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if result := product.Get(tc.column); result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
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
