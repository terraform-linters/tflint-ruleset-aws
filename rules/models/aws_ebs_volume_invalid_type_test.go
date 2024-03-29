// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"testing"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEbsVolumeInvalidTypeRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_ebs_volume" "foo" {
	type = "gp1"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEbsVolumeInvalidTypeRule(),
					Message: fmt.Sprintf(`"%s" is an invalid value as %s`, truncateLongMessage("gp1"), "type"),
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_ebs_volume" "foo" {
	type = "gp2"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsEbsVolumeInvalidTypeRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
