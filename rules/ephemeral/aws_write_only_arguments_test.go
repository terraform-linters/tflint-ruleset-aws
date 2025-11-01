package ephemeral

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsWriteOnlyAttribute(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "basic",
			Content: `
resource "aws_secretsmanager_secret_version" "test" {
  secret_string = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsWriteOnlyArgumentsRule(),
					Message: `"secret_string" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only argument "secret_string_wo".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 19},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
			},
			Fixed: `
resource "aws_secretsmanager_secret_version" "test" {
  secret_string_wo         = "test"
  secret_string_wo_version = 1
}
`,
		},
		{
			Name: "everything is fine",
			Content: `
resource "aws_secretsmanager_secret_version" "test" {
  secret_string_wo         = "test"
  secret_string_wo_version = 1
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "without version attribute",
			Content: `
resource "aws_transfer_host_key" "test" {
  host_key_body = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsWriteOnlyArgumentsRule(),
					Message: `"host_key_body" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only argument "host_key_body_wo".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 19},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
			},
			Fixed: `
resource "aws_transfer_host_key" "test" {
  host_key_body_wo = "test"
}
`,
		},
	}

	rule := NewAwsWriteOnlyArgumentsRule()

	for _, tc := range cases {
		filename := "resource.tf"
		runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}
		helper.AssertIssues(t, tc.Expected, runner.Issues)

		want := map[string]string{}
		if tc.Fixed != "" {
			want[filename] = tc.Fixed
		}
		helper.AssertChanges(t, want, runner.Changes())
	}
}
