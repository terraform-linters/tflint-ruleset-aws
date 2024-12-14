package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSecurityGroupRuleDeprecated(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "resource is used",
			Content: `
resource "aws_security_group_rule" "test" {
	security_group_id = "sg-12345678"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupRuleDeprecatedRule(),
					Message: "Consider using aws_vpc_security_group_egress_rule or aws_vpc_security_group_ingress_rule instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 22},
						End:      hcl.Pos{Line: 3, Column: 35},
					},
				},
			},
		},
		{
			Name: "everything is fine",
			Content: `
resource "aws_security_group" "test" {
	name = "test"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSecurityGroupRuleDeprecatedRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
