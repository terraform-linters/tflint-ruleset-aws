package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSpotFleetRequestInvalidExcessCapacityTerminationPolicy(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "invalid value with capital D - Default",
			Content: `
resource "aws_spot_fleet_request" "foo" {
	excess_capacity_termination_policy = "Default"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule(),
					Message: `"Default" is an invalid value as excess_capacity_termination_policy`,
				},
			},
		},
		{
			Name: "invalid value with capital N - NoTermination",
			Content: `
resource "aws_spot_fleet_request" "foo" {
	excess_capacity_termination_policy = "NoTermination"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule(),
					Message: `"NoTermination" is an invalid value as excess_capacity_termination_policy`,
				},
			},
		},
		{
			Name: "invalid value",
			Content: `
resource "aws_spot_fleet_request" "foo" {
	excess_capacity_termination_policy = "invalidValue"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule(),
					Message: `"invalidValue" is an invalid value as excess_capacity_termination_policy`,
				},
			},
		},
		{
			Name: "valid value - default",
			Content: `
resource "aws_spot_fleet_request" "foo" {
	excess_capacity_termination_policy = "default"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "valid value - noTermination",
			Content: `
resource "aws_spot_fleet_request" "foo" {
	excess_capacity_termination_policy = "noTermination"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
