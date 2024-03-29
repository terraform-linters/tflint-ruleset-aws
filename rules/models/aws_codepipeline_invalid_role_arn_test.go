// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"testing"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsCodepipelineInvalidRoleArnRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_codepipeline" "foo" {
	role_arn = "arn:aws:iam::123456789012:instance-profile/s3access-profile"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCodepipelineInvalidRoleArnRule(),
					Message: fmt.Sprintf(`"%s" does not match valid pattern %s`, truncateLongMessage("arn:aws:iam::123456789012:instance-profile/s3access-profile"), `^arn:aws(-[\w]+)*:iam::[0-9]{12}:role/.*$`),
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_codepipeline" "foo" {
	role_arn = "arn:aws:iam::123456789012:role/s3access"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsCodepipelineInvalidRoleArnRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
