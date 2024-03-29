// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"testing"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDatasyncAgentInvalidActivationKeyRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_datasync_agent" "foo" {
	activation_key = "F0EFT7FPPRGG7MC3I9R327DOH"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDatasyncAgentInvalidActivationKeyRule(),
					Message: fmt.Sprintf(`"%s" does not match valid pattern %s`, truncateLongMessage("F0EFT7FPPRGG7MC3I9R327DOH"), `^[A-Z0-9]{5}(-[A-Z0-9]{5}){4}$`),
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_datasync_agent" "foo" {
	activation_key = "F0EFT-7FPPR-GG7MC-3I9R3-27DOH"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDatasyncAgentInvalidActivationKeyRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
