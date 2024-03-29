// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"testing"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEcsServiceInvalidPropagateTagsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_ecs_service" "foo" {
	propagate_tags = "CONTAINER"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEcsServiceInvalidPropagateTagsRule(),
					Message: fmt.Sprintf(`"%s" is an invalid value as %s`, truncateLongMessage("CONTAINER"), "propagate_tags"),
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_ecs_service" "foo" {
	propagate_tags = "SERVICE"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsEcsServiceInvalidPropagateTagsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
