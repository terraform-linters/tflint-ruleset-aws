// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"testing"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLbInvalidLoadBalancerTypeRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_lb" "foo" {
	load_balancer_type = "classic"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLbInvalidLoadBalancerTypeRule(),
					Message: fmt.Sprintf(`"%s" is an invalid value as %s`, truncateLongMessage("classic"), "load_balancer_type"),
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_lb" "foo" {
	load_balancer_type = "application"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLbInvalidLoadBalancerTypeRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
