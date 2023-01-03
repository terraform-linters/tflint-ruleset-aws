package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSecurityGroupRuleInvalidCidrBlock(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "Find issues when 0.0.0.0/0 allows ingress access through 22",
			Content: `
resource "aws_security_group_rule" "rule" {
  from_port         = 10
  to_port           = 300
  protocol          = "tcp"
  type              = "ingress"
  cidr_blocks       = ["0.0.0.0/0", "10.0.0.0/16"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupRuleInvalidCidrBlockRule(),
					Message: "cidr_blocks can not contain '0.0.0.0/0' when allowing 'ingress' access to ports [22 3389]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 23},
						End:      hcl.Pos{Line: 7, Column: 51},
					},
				},
			},
		},
		{
			Name: "Find issues when ::/0 allows ingress access through 22",
			Content: `
resource "aws_security_group_rule" "rule" {
  from_port         = 10
  to_port           = 300
  protocol          = "tcp"
  type              = "ingress"
  ipv6_cidr_blocks  = ["::/0"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupRuleInvalidCidrBlockRule(),
					Message: "ipv6_cidr_blocks can not contain '::/0' when allowing 'ingress' access to ports [22 3389]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 23},
						End:      hcl.Pos{Line: 7, Column: 31},
					},
				},
			},
		},
		{
			Name: "Find no issues 0.0.0.0/0 allows ingress access through 80",
			Content: `
resource "aws_security_group_rule" "rule" {
  from_port         = 80
  to_port           = 80
  protocol          = "tcp"
  type              = "ingress"
  cidr_blocks       = ["0.0.0.0/0"]
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSecurityGroupRuleInvalidCidrBlockRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
