package rules

import (
	"testing"
	"time"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaFunctionEndOfSupport(t *testing.T) {
	type caseStudy struct {
		Name     string
		Content  string
		Expected helper.Issues
	}
	eolRuntimes :=  map[string]time.Time{
		"nodejs10.x":     time.Date(2021, time.July, 30, 0, 0, 0, 0, time.UTC),
		"ruby2.5":        time.Date(2021, time.July, 30, 0, 0, 0, 0, time.UTC),
		"python2.7":      time.Date(2021, time.July, 15, 0, 0, 0, 0, time.UTC),
		"dotnetcore2.1":  time.Date(2021, time.September, 20, 0, 0, 0, 0, time.UTC),
	}
	var cases []caseStudy
	for runtime, eolDate := range eolRuntimes {
		now := time.Now().UTC()
		if now.Before(eolDate) {
			continue;
		}
		study := caseStudy{
			Name: runtime + " end of life",
			Content: `
resource "aws_lambda_function" "function" {
	function_name = "test_function"
	role = "test_role"
	runtime = "` + runtime + `"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionEndOfSupportRule(),
					Message: "The \"" + runtime + "\" runtime has reached the end of support",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 12},
						End:      hcl.Pos{Line: 5, Column: len(runtime) + 14},
					},
				},
			},
		}
		cases = append(cases, study)
	}

	rule := NewAwsLambdaFunctionEndOfSupportRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
