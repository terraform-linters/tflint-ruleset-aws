package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsIAMRoleDeprecatedPolicyAttributes(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline_policies",
			Content: `
resource "aws_iam_role" "example" {
  name               = "yak_role"
  assume_role_policy = data.aws_iam_policy_document.instance_assume_role_policy.json # (not shown)

  inline_policy {
    name = "my_inline_policy"

    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action   = ["ec2:Describe*"]
          Effect   = "Allow"
          Resource = "*"
        },
      ]
    })
  }

  inline_policy {
    name   = "policy-8675309"
    policy = data.aws_iam_policy_document.inline_policy.json
  }
}

data "aws_iam_policy_document" "inline_policy" {
  statement {
    actions   = ["ec2:DescribeAccountAttributes"]
    resources = ["*"]
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMRoleDeprecatedPolicyAttributesRule(),
					Message: "The inline_policy argument is deprecated. Use the aws_iam_role_policy resource instead. If Terraform should exclusively manage all inline policy associations (the current behavior of this argument), use the aws_iam_role_policies_exclusive resource as well.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 16},
					},
				},
				{
					Rule:    NewAwsIAMRoleDeprecatedPolicyAttributesRule(),
					Message: "The inline_policy argument is deprecated. Use the aws_iam_role_policy resource instead. If Terraform should exclusively manage all inline policy associations (the current behavior of this argument), use the aws_iam_role_policies_exclusive resource as well.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 21, Column: 3},
						End:      hcl.Pos{Line: 21, Column: 16},
					},
				},
			},
		},
		{
			Name: "managed_policy_arns",
			Content: `
resource "aws_iam_role" "example" {
  name                = "yak_role"
  assume_role_policy  = data.aws_iam_policy_document.instance_assume_role_policy.json # (not shown)
  managed_policy_arns = []
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMRoleDeprecatedPolicyAttributesRule(),
					Message: "The managed_policy_arns argument is deprecated. Use the aws_iam_role_policy_attachment resource instead. If Terraform should exclusively manage all managed policy attachments (the current behavior of this argument), use the aws_iam_role_policy_attachments_exclusive resource as well.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 27},
					},
				},
			},
		},
	}

	rule := NewAwsIAMRoleDeprecatedPolicyAttributesRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
