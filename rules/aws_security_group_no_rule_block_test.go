package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSecurityGroupNoRuleBlock(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
		Error    bool
	}{
		{
			Name: "ingress-egress",
			Content: `
resource "aws_security_group" "allow_tls" {
  ingress {
    description      = "TLS from VPC"
    from_port        = 443
    to_port          = 443
    protocol         = "tcp"
    cidr_blocks      = [aws_vpc.main.cidr_block]
    ipv6_cidr_blocks = [aws_vpc.main.ipv6_cidr_block]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupNoRuleBlockRule(),
					Message: "Inline rule block `ingress` found! It is advisable to create a separate resource (`aws_security_group_rule`) instead",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 10},
					},
				},
				{
					Rule:    NewAwsSecurityGroupNoRuleBlockRule(),
					Message: "Inline rule block `egress` found! It is advisable to create a separate resource (`aws_security_group_rule`) instead",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 3},
						End:      hcl.Pos{Line: 12, Column: 9},
					},
				},
			},
		},
		{
			Name: "empty",
			Content: `
resource "aws_security_group" "my_sg" {
  name        = "my_sg"
  description = "My security group"
  vpc_id      = aws_vpc.main.id
}
`,
			Expected: helper.Issues{},
		},
	}
	rule := NewAwsSecurityGroupNoRuleBlockRule()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content, ".tflint.hcl": tc.Config})

			err := rule.Check(runner)
			if err != nil && !tc.Error {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			if err == nil && tc.Error {
				t.Fatal("Expected error but got none")
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
