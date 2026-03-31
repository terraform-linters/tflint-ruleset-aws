package models

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketInvalidTagsRule(t *testing.T) {
	rule := NewAwsS3BucketInvalidTagsRule()

	for _, tc := range []struct {
		name     string
		content  string
		expected helper.Issues
	}{
		{
			name: "valid tags",
			content: `
resource "aws_s3_bucket" "foo" {
	tags = {
		Name = "example"
	}
}`,
			expected: helper.Issues{},
		},
		{
			name: "empty key",
			content: `
resource "aws_s3_bucket" "foo" {
	tags = {
		"" = "value"
	}
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: `tags key "" must be at least 1 character`,
				},
			},
		},
		{
			name: "valid empty value",
			content: `
resource "aws_s3_bucket" "foo" {
	tags = {
		Name = ""
	}
}`,
			expected: helper.Issues{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.expected, runner.Issues)
		})
	}
}
