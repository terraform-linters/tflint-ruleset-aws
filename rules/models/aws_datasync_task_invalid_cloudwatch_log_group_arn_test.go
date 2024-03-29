// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"testing"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDatasyncTaskInvalidCloudwatchLogGroupArnRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_datasync_task" "foo" {
	cloudwatch_log_group_arn = "arn:aws:s3:::my_corporate_bucket"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDatasyncTaskInvalidCloudwatchLogGroupArnRule(),
					Message: fmt.Sprintf(`"%s" does not match valid pattern %s`, truncateLongMessage("arn:aws:s3:::my_corporate_bucket"), `^arn:(aws|aws-cn|aws-us-gov|aws-iso|aws-iso-b):logs:[a-z\-0-9]+:[0-9]{12}:log-group:([^:\*]*)(:\*)?$`),
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_datasync_task" "foo" {
	cloudwatch_log_group_arn = "arn:aws:logs:us-east-1:123456789012:log-group:my-log-group"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDatasyncTaskInvalidCloudwatchLogGroupArnRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
