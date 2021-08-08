package rules

import (
	"testing"
	"time"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaFunctionEol(t *testing.T) {
	type caseStudy struct {
		Name     string
		Content  string
		Expected helper.Issues
	}
	eolRuntimes :=  map[string]time.Time{
		"dotnetcore1.0":  time.Date(2019, time.July, 30, 0, 0, 0, 0, time.UTC),
		"dotnetcore2.0":  time.Date(2019, time.May, 30, 0, 0, 0, 0, time.UTC),
		"nodejs":         time.Date(2016, time.October, 31, 0, 0, 0, 0, time.UTC),
		"nodejs4.3":      time.Date(2020, time.March, 06, 0, 0, 0, 0, time.UTC),
		"nodejs4.3-edge": time.Date(2019, time.April, 30, 0, 0, 0, 0, time.UTC),
		"nodejs6.10":     time.Date(2019, time.August, 12, 0, 0, 0, 0, time.UTC),
		"nodejs8.10":     time.Date(2020, time.March, 06, 0, 0, 0, 0, time.UTC),
		"nodejs10.x":     time.Date(2021, time.August, 30, 0, 0, 0, 0, time.UTC),
		"ruby2.5":        time.Date(2021, time.August, 30, 0, 0, 0, 0, time.UTC),
		"python2.7":      time.Date(2021, time.September, 30, 0, 0, 0, 0, time.UTC),
		"dotnetcore2.1":  time.Date(2021, time.October, 30, 0, 0, 0, 0, time.UTC),
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
					Rule:    NewAwsLambdaFunctionEolRule(),
					Message: "The \"" + runtime + "\" runtime has reached the end of life",
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
	// 	cases := []struct {
	// 	Name     string
	// 	Content  string
	// 	Expected helper.Issues
	// }
	// 	{
	// 		{
	// 			Name: "dotnetcore1.0 end of life",
	// 			Content: `
	// resource "aws_lambda_function" "function" {
	// 	function_name = "test_function"
	// 	role = "test_role"
	// 	runtime = "dotnetcore1.0"
	// }
	// `,
	// 			Expected: helper.Issues{
	// 				{
	// 					Rule:    NewAwsLambdaFunctionEolRule(),
	// 					Message: "The \"dotnetcore1.0\" runtime has reached the end of life",
	// 					Range: hcl.Range{
	// 						Filename: "resource.tf",
	// 						Start:    hcl.Pos{Line: 5, Column: 12},
	// 						End:      hcl.Pos{Line: 5, Column: 27},
	// 					},
	// 				},
	// 			},
	// 		},
	// 		{
	// 			Name: "dotnetcore2.0 end of life",
	// 			Content: `
	// resource "aws_lambda_function" "function" {
	// 	function_name = "test_function"
	// 	role = "test_role"
	// 	runtime = "dotnetcore2.0"
	// }
	// `,
	// 			Expected: helper.Issues{
	// 				{
	// 					Rule:    NewAwsLambdaFunctionEolRule(),
	// 					Message: "The \"dotnetcore1.0\" runtime has reached the end of life",
	// 					Range: hcl.Range{
	// 						Filename: "resource.tf",
	// 						Start:    hcl.Pos{Line: 5, Column: 12},
	// 						End:      hcl.Pos{Line: 5, Column: 27},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	}

	rule := NewAwsLambdaFunctionEolRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
