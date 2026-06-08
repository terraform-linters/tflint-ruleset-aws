package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsCloudwatchLogGroupRetention(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "retention below 365 days",
			Content: `
resource "aws_cloudwatch_log_group" "short" {
  name              = "/aws/lambda/my-function"
  retention_in_days = 30
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCloudwatchLogGroupRetentionRule(),
					Message: "The CloudWatch log group retention is 30 days. Set to at least 365 days.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 23},
						End:      hcl.Pos{Line: 4, Column: 25},
					},
				},
			},
		},
		{
			Name: "retention not set",
			Content: `
resource "aws_cloudwatch_log_group" "forever" {
  name = "/aws/lambda/my-function"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCloudwatchLogGroupRetentionRule(),
					Message: "The CloudWatch log group does not have retention_in_days set. It will retain logs forever. Set to at least 365 days.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 46},
					},
				},
			},
		},
		{
			Name: "retention at 365 days is valid",
			Content: `
resource "aws_cloudwatch_log_group" "ok" {
  name              = "/aws/lambda/my-function"
  retention_in_days = 365
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "retention set to 0 (never expire) is valid",
			Content: `
resource "aws_cloudwatch_log_group" "never" {
  name              = "/aws/lambda/my-function"
  retention_in_days = 0
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsCloudwatchLogGroupRetentionRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
