// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"testing"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSpotFleetRequestInvalidInstanceInterruptionBehaviourRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_spot_fleet_request" "foo" {
	instance_interruption_behaviour = "restart"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSpotFleetRequestInvalidInstanceInterruptionBehaviourRule(),
					Message: fmt.Sprintf(`"%s" is an invalid value as %s`, truncateLongMessage("restart"), "instance_interruption_behaviour"),
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_spot_fleet_request" "foo" {
	instance_interruption_behaviour = "hibernate"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSpotFleetRequestInvalidInstanceInterruptionBehaviourRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
