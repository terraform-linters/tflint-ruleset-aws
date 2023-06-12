package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3NoGlobalEndpoint(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "unrelated expression",
			Content: `
output "test" {
  value = "testing"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "multiple unrelated expressions",
			Content: `
output "test1" {
  value = var.whatever
}

output "test2" {
  value = "testing"
}

output "test3" {
  value = aws_iam_role.test.arn
}
  `,
			Expected: helper.Issues{},
		},
		// TODO: nested expressions ?
		{
			Name: "regional endpoint used",
			Content: `
output "test" {
  value = aws_s3_bucket.test.bucket_regional_domain_name
}`,
			Expected: helper.Issues{},
		},
		// TODO: strings of the form "bucket_name.s3.amazonaws.com" ?
		{
			Name: "legacy global endpoint used",
			Content: `
output "test" {
  value = aws_s3_bucket.test.bucket_domain_name
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 48},
					},
				},
			},
		},
	}

	rule := NewAwsS3NoGlobalEndpointRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
