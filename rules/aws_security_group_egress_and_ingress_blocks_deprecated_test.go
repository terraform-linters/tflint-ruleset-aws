package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSecurityGroupEgressAndIngressBlocksDeprecated(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "finds deprecated ingress block",
			Content: `
resource "aws_security_group" "test" {
  name = "test"

  ingress {
	from_port   = 0
	to_port     = 0
	protocol    = "-1"
	cidr_blocks = ["0.0.0.0/0"]
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupEgressAndIngressBlocksDeprecatedRule(),
					Message: "Replace this ingress block with aws_vpc_security_group_ingress_rule. Otherwise, rule conflicts, perpetual differences, and result in rules being overwritten may occur.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 10},
					},
				},
			},
		},
		{
			Name: "finds deprecated egress block",
			Content: `
resource "aws_security_group" "test" {
  name = "test"

  egress {
	from_port   = 0
	to_port     = 0
	protocol    = "-1"
	cidr_blocks = ["0.0.0.0/0"]
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupEgressAndIngressBlocksDeprecatedRule(),
					Message: "Replace this egress block with aws_vpc_security_group_egress_rule. Otherwise, rule conflicts, perpetual differences, and result in rules being overwritten may occur.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 9},
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

	rule := NewAwsSecurityGroupEgressAndIngressBlocksDeprecatedRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
