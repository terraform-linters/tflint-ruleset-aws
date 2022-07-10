package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSecurityGroupInvalidProtocol(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "valid protocols",
			Content: `
resource "aws_security_group" "this" {
    ingress {
        protocol = "tcp"
    }
	egress {
        protocol = "-1"
    }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid protocol with upper case (for Terraform v0.11 and older)",
			Content: `
resource "aws_security_group" "this" {
    ingress {
        protocol = "TCP"
    }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid ingress.protocol",
			Content: `
resource "aws_security_group" "this" {
    ingress {
        protocol = "http"
    }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupInvalidProtocolRule(),
					Message: "\"http\" is an invalid protocol.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 20},
						End:      hcl.Pos{Line: 4, Column: 26},
					},
				},
			},
		},
		{
			Name: "invalid egress.protocol",
			Content: `
resource "aws_security_group" "this" {
    egress {
        protocol = "all"
    }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupInvalidProtocolRule(),
					Message: "\"all\" is an invalid protocol.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 20},
						End:      hcl.Pos{Line: 4, Column: 25},
					},
				},
			},
		},
		{
			Name: "invalid protocol with upper case (for Terraform v0.11 and older)",
			Content: `
resource "aws_security_group" "this" {
    ingress {
        protocol = "HTTP"
    }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupInvalidProtocolRule(),
					Message: "\"HTTP\" is an invalid protocol.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 20},
						End:      hcl.Pos{Line: 4, Column: 26},
					},
				},
			},
		},
	}

	rule := NewAwsSecurityGroupInvalidProtocolRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
