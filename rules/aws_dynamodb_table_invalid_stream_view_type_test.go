package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDynamoDBTableInvalidStreamViewTypeRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_dynamodb_table" "foo" {
	stream_view_type = "OLD_AND_NEW_IMAGE"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDynamoDBTableInvalidStreamViewTypeRule(),
					Message: `"OLD_AND_NEW_IMAGE" is an invalid value as stream_view_type`,
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_dynamodb_table" "foo" {
	stream_view_type = "NEW_IMAGE"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "empty string",
			Content: `
resource "aws_dynamodb_table" "foo" {
	stream_view_type = ""
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDynamoDBTableInvalidStreamViewTypeRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
