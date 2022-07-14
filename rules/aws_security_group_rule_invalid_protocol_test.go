package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSecurityGroupRuleInvalidProtocol(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "valid protocol",
			Content: `
resource "aws_security_group_rule" "this" {
    protocol = "tcp"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid number as protocol",
			Content: `
resource "aws_security_group_rule" "this" {
    protocol = "-1"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid protocol",
			Content: `
resource "aws_security_group_rule" "this" {
    protocol = "http"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupRuleInvalidProtocolRule(),
					Message: "\"http\" is an invalid protocol.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 16},
						End:      hcl.Pos{Line: 3, Column: 22},
					},
				},
			},
		},
	}

	rule := NewAwsSecurityGroupRuleInvalidProtocolRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
