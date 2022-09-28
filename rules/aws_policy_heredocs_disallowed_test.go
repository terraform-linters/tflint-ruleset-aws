package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsPolicyHeredocsDisallowed(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "No Heredoc",
			Content: `
resource "aws_iam_policy" "test" {
	policy = jsonencode({
		Statement = [
			{
				Action = [
					"iam:PassRole"
				],
			}
		]
	})
}
`,
			Expected: helper.Issues{},
		},
		{
			Name:    "EOF for a policy",
			Content: `
resource "aws_iam_policy" "test" {
	policy = <<EOF
{
	"Statement": [
		{
			"Action":  [
				"iam:PassRole"
			],
		}
	]
}
EOF
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsPolicyHeredocsDisallowedRule(),
					Message: "Avoid Use of HEREDOC syntax for policy documents",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 4, Column: 1},
					},
				},
				{
					Rule:    NewAwsPolicyHeredocsDisallowedRule(),
					Message: "Avoid Use of HEREDOC syntax for policy documents",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 13, Column: 1},
						End:      hcl.Pos{Line: 13, Column: 4},
					},
				},
			},
		},
	}

	rule := NewAwsPolicyHeredocsDisallowedRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
