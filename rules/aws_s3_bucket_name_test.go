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
resource "aws_s3_bucket" "too_long" {
  bucket = "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameRule(),
					Message: `Bucket names must be between 3 (min) and 63 (max) characters long. (name: "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test", regex: "^(.{0,2})?(.{64,})?$")`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 81},
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
