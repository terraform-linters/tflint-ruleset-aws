// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsCognitoUserPoolDomainInvalidDomainRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_cognito_user_pool_domain" "foo" {
	domain = "auth example"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCognitoUserPoolDomainInvalidDomainRule(),
					Message: `"auth example" does not match valid pattern ^[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?$`,
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_cognito_user_pool_domain" "foo" {
	domain = "auth"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsCognitoUserPoolDomainInvalidDomainRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
