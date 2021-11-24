package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsAcmCertificateLifecycle(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "no lifecycle block",
			Content: `
resource "aws_acm_certificate" "test" {
  domain_name       = var.domain_name
  validation_method = "DNS"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAcmCertificateLifecycleRule(),
					Message: "resource `aws_acm_certificate` needs to contain `create_before_destroy = true` in `lifecycle` block",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 38},
					},
				},
			},
		},
		{
			Name: "no create_before_destroy attribute in lifecycle block",
			Content: `
resource "aws_acm_certificate" "test" {
  domain_name       = var.domain_name
  validation_method = "DNS"
  lifecycle {}
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAcmCertificateLifecycleRule(),
					Message: "resource `aws_acm_certificate` needs to contain `create_before_destroy = true` in `lifecycle` block",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 38},
					},
				},
			},
		},
		{
			Name: "create_before_destroy = false",
			Content: `
resource "aws_acm_certificate" "test" {
  domain_name       = var.domain_name
  validation_method = "DNS"
  lifecycle {
    create_before_destroy = true
  }
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsAcmCertificateLifecycleRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
