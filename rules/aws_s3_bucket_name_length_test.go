package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketNameLength(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "basic",
			Content: `
resource "aws_s3_bucket" "too_long" {
  bucket = "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketNameLengthRule(),
					Message: `Bucket name "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test" must be between 3 and 63 characters`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 81},
					},
				},
			},
		},
	}

	rule := NewAwsS3BucketNameLengthRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
