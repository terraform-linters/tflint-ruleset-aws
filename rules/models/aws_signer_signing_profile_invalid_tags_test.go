package models

import (
	"fmt"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSignerSigningProfileInvalidTagsRule(t *testing.T) {
	rule := NewAwsSignerSigningProfileInvalidTagsRule()

	for _, tc := range []struct {
		name     string
		content  string
		expected helper.Issues
	}{
		{
			name: "valid tags",
			content: `
resource "aws_signer_signing_profile" "foo" {
	tags = {
		Name = "example"
	}
}`,
			expected: helper.Issues{},
		},
		{
			name: "aws prefix denied",
			content: `
resource "aws_signer_signing_profile" "foo" {
	tags = {
		"aws:reserved" = "value"
	}
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf(`tags key %q must not start with %q`, "aws:reserved", "aws:"),
				},
			},
		},
		{
			name: "key pattern invalid",
			content: `
resource "aws_signer_signing_profile" "foo" {
	tags = {
		"bad key!" = "value"
	}
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf(`tags key %q does not match valid pattern %s`, "bad key!", `^[a-zA-Z+-=._:/]+$`),
				},
			},
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
