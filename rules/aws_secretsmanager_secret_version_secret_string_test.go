package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSecretsmanagerSecretVersionSecretString(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
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
					Rule:    NewAwsSecretsmanagerSecretVersionSecretStringRule(),
					Message: `"secret_string" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only attribute "secret_string_wo".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 19},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
			},
		},
		{
			Name: "everything is fine",
			Content: `
resource "aws_secretsmanager_secret_version" "test" {
  secret_string_wo = "test"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSecretsmanagerSecretVersionSecretStringRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
