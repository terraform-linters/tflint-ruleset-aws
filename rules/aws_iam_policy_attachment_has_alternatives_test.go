package rules

import (
	"math/rand"
	"testing"
	"time"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsIAMPolicyAttachmentHasAlternativesRule(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "resource has alternatives",
			Content: `
resource "aws_iam_policy_attachment" "attachment" {
	name = "test_attachment"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMPolicyAttachmentHasAlternativesRule(),
					Message: "Consider aws_iam_role_policy_attachment, aws_iam_user_policy_attachment, or aws_iam_group_policy_attachment instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 9},
						End:      hcl.Pos{Line: 3, Column: 26},
					},
				},
			},
		},
		{
			Name: "no issues with resource",
			Content: `
resource "aws_iam_role_policy_attachment" "attachment" {
	role       = "test_role"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsIAMPolicyAttachmentHasAlternativesRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
