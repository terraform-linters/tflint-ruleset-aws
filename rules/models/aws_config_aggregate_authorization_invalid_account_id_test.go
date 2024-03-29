// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"testing"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsConfigAggregateAuthorizationInvalidAccountIDRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_config_aggregate_authorization" "foo" {
	account_id = "01234567891"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsConfigAggregateAuthorizationInvalidAccountIDRule(),
					Message: fmt.Sprintf(`"%s" does not match valid pattern %s`, truncateLongMessage("01234567891"), `^\d{12}$`),
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_config_aggregate_authorization" "foo" {
	account_id = "012345678910"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsConfigAggregateAuthorizationInvalidAccountIDRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
