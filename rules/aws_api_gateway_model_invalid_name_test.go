package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsAPIGatewayModelInvalidName(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "underscore",
			Content: `
resource "aws_api_gateway_model" "demo" {
  name = "demo_model"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAPIGatewayModelInvalidNameRule(),
					Message: `demo_model does not match valid pattern ^[a-zA-Z0-9]+$`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 10},
						End:      hcl.Pos{Line: 3, Column: 22},
					},
				},
			},
		},
		{
			Name: "alphanumeric",
			Content: `
resource "aws_api_gateway_model" "demo" {
  name = "user"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsAPIGatewayModelInvalidNameRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
