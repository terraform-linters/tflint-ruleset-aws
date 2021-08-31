package rules

import (
	"testing"
	"time"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaFunctionEndOfSupport(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Now      time.Time
		Expected helper.Issues
	}{
		{
			Name: "EOS",
			Content: `
resource "aws_lambda_function" "function" {
	function_name = "test_function"
	role = "test_role"
	runtime = "nodejs10.x"
}
`,
			Now: time.Date(2021, time.August, 10, 0, 0, 0, 0, time.UTC),
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionDeprecatedRuntimeRule(),
					Message: "The \"nodejs10.x\" runtime has reached the end of support",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 12},
						End:      hcl.Pos{Line: 5, Column: 24},
					},
				},
			},
		},
		{
			Name: "EOF",
			Content: `
resource "aws_lambda_function" "function" {
	function_name = "test_function"
	role = "test_role"
	runtime = "nodejs10.x"
}
`,
			Now: time.Date(2021, time.September, 1, 0, 0, 0, 0, time.UTC),
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionDeprecatedRuntimeRule(),
					Message: "The \"nodejs10.x\" runtime has reached the end of life",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 12},
						End:      hcl.Pos{Line: 5, Column: 24},
					},
				},
			},
		},
		{
			Name: "Live",
			Content: `
resource "aws_lambda_function" "function" {
	function_name = "test_function"
	role = "test_role"
	runtime = "nodejs10.x"
}
`,
			Now:      time.Date(2021, time.June, 25, 0, 0, 0, 0, time.UTC),
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLambdaFunctionDeprecatedRuntimeRule()

	for _, tc := range cases {
		rule.Now = tc.Now

		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
