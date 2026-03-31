package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEcsServiceInvalidTagsRule(t *testing.T) {
	rule := NewAwsEcsServiceInvalidTagsRule()

	for _, tc := range []struct {
		name     string
		content  string
		expected helper.Issues
	}{
		{
			name: "valid tags",
			content: `
resource "aws_ecs_service" "foo" {
	tags = {
		Name = "example"
	}
}`,
			expected: helper.Issues{},
		},
		{
			name: "key too long",
			content: fmt.Sprintf(`
resource "aws_ecs_service" "foo" {
	tags = {
		%q = "value"
	}
}`, strings.Repeat("a", 129)),
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("tag key %q must be 128 characters or less", truncateLongMessage(strings.Repeat("a", 129))),
				},
			},
		},
		{
			name: "empty key",
			content: `
resource "aws_ecs_service" "foo" {
	tags = {
		"" = "value"
	}
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: `tag key "" must be at least 1 characters`,
				},
			},
		},
		{
			name: "key with invalid characters",
			content: `
resource "aws_ecs_service" "foo" {
	tags = {
		"invalid!key" = "value"
	}
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf(`tag key %q does not match valid pattern %s`, "invalid!key", `^([\p{L}\p{Z}\p{N}_.:/=+\-@]*)$`),
				},
			},
		},
		{
			name: "value too long",
			content: fmt.Sprintf(`
resource "aws_ecs_service" "foo" {
	tags = {
		Name = %q
	}
}`, strings.Repeat("v", 257)),
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: `tag value for key "Name" must be 256 characters or less`,
				},
			},
		},
		{
			name: "value with invalid characters",
			content: `
resource "aws_ecs_service" "foo" {
	tags = {
		Name = "bad<value>"
	}
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf(`tag value %q for key %q does not match valid pattern %s`, "bad<value>", "Name", `^([\p{L}\p{Z}\p{N}_.:/=+\-@]*)$`),
				},
			},
		},
		{
			name: "too many tags",
			content: func() string {
				var tags []string
				for i := range 51 {
					tags = append(tags, fmt.Sprintf(`    tag%d = "value"`, i))
				}
				return fmt.Sprintf(`
resource "aws_ecs_service" "foo" {
	tags = {
%s
	}
}`, strings.Join(tags, "\n"))
			}(),
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: "too many tags: 51 exceeds the maximum of 50",
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
