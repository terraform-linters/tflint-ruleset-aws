package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketName(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
		Error    bool
	}{
		{
			Name: "regex",
			Content: `
resource "aws_s3_bucket" "foo" {
  bucket = "foo-bucket"
}

resource "aws_s3_bucket" "bar" {
	bucket = "bar-bucket"
}`,
			Config: `
rule "aws_s3_bucket_name" {
	enabled = true
	regex = "^foo"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket name "bar-bucket" does not match regex "^foo"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 11},
						End:      hcl.Pos{Line: 7, Column: 23},
					},
				},
			},
		},
		{
			Name:    "invalid regex",
			Content: ``,
			Config: `
rule "aws_s3_bucket_name" {
	enabled = true
	regex = "\\"
}`,
			Expected: helper.Issues{},
			Error:    true,
		},
		{
			Name: "prefix",
			Content: `
resource "aws_s3_bucket" "foo" {
  bucket = "foo-bar"
}

resource "aws_s3_bucket" "bar" {
	bucket = "bar-baz"
}`,
			Config: `
rule "aws_s3_bucket_name" {
	enabled = true
	prefix = "foo-"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket name "bar-baz" does not have prefix "foo-"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 11},
						End:      hcl.Pos{Line: 7, Column: 20},
					},
				},
			},
		},
		{
			Name: "regex and prefix",
			Content: `
resource "aws_s3_bucket" "foo" {
  bucket = "foo-bar"
}

resource "aws_s3_bucket" "bar" {
	bucket = "bar-baz"
}`,
			Config: `
rule "aws_s3_bucket_name" {
	enabled = true
	prefix = "foo-"
	regex = "^f"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket name "bar-baz" does not have prefix "foo-"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 11},
						End:      hcl.Pos{Line: 7, Column: 20},
					},
				},
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket name "bar-baz" does not match regex "^f"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 11},
						End:      hcl.Pos{Line: 7, Column: 20},
					},
				},
			},
		},
		{
			Name: "length",
			Content: `
resource "aws_s3_bucket" "max_length" {
  bucket = "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test"
}

resource "aws_s3_bucket" "min_length" {
  bucket = "xy"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must be between 3 (min) and 63 (max) characters long. (name: "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test", regex: "^.{0,2}$|.{64,}$")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 81},
					},
				},
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must be between 3 (min) and 63 (max) characters long. (name: "xy", regex: "^.{0,2}$|.{64,}$")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 12},
						End:      hcl.Pos{Line: 7, Column: 16},
					},
				},
			},
		},
		{
			Name: "invalid_characters",
			Content: `
resource "aws_s3_bucket" "invalid_characters" {
  bucket = "who_am_I?.com"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names can consist only of lowercase letters, numbers, dots (.), and hyphens (-). (name: "who_am_I?.com", regex: "[^a-z0-9\\.\\-]")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 27},
					},
				},
			},
		},
		{
			Name: "invalid_begin_end",
			Content: `
resource "aws_s3_bucket" "invalid_begin" {
  bucket = ".domain.com"
}

resource "aws_s3_bucket" "invalid_end" {
  bucket = "domain.com."
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must begin and end with a lowercase letter or number. (name: ".domain.com", regex: "^[^a-z0-9]|[^a-z0-9]$")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must begin and end with a lowercase letter or number. (name: "domain.com.", regex: "^[^a-z0-9]|[^a-z0-9]$")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 12},
						End:      hcl.Pos{Line: 7, Column: 25},
					},
				},
			},
		},
		{
			Name: "adjacent_periods",
			Content: `
resource "aws_s3_bucket" "adjacent_periods" {
  bucket = "domain..com"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must not contain two adjacent periods. (name: "domain..com", regex: "\\.\\.")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
			},
		},
		{
			Name: "ipv4",
			Content: `
resource "aws_s3_bucket" "ipv4" {
  bucket = "192.168.0.254"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must not be formatted as an IP address. (name: "192.168.0.254")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 27},
					},
				},
			},
		},
		{
			Name: "invalid_prefix_xn",
			Content: `
resource "aws_s3_bucket" "invalid_prefix_xn" {
  bucket = "xn--domain.com"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must not start with the prefix 'xn--'. (name: "xn--domain.com", regex: "^xn--")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 28},
					},
				},
			},
		},
		{
			Name: "invalid_prefix_sthree",
			Content: `
resource "aws_s3_bucket" "invalid_prefix_sthree" {
  bucket = "sthree-domain.com"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must not start with the prefix 'sthree-' and the prefix 'sthree-configurator'. (name: "sthree-domain.com", regex: "^(sthree-|sthree-configurator)")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 31},
					},
				},
			},
		},
		{
			Name: "invalid_suffix_s3alias",
			Content: `
resource "aws_s3_bucket" "invalid_suffix_s3alias" {
  bucket = "domain-s3alias"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must not end with the suffix '-s3alias'. (name: "domain-s3alias", regex: "-s3alias$")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 28},
					},
				},
			},
		},
		{
			Name: "invalid_suffix_ols3",
			Content: `
resource "aws_s3_bucket" "invalid_suffix_ols3" {
  bucket = "domain--ol-s3"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must not end with the suffix '--ol-s3'. (name: "domain--ol-s3", regex: "--ol-s3$")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 27},
					},
				},
			},
		},
	}

	rule := NewAwsS3BucketNameRule()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content, ".tflint.hcl": tc.Config})

			err := rule.Check(runner)
			if err != nil && !tc.Error {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			if err == nil && tc.Error {
				t.Fatal("Expected error but got none")
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
