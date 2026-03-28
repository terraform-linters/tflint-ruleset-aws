package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsIAMPolicyUsePolicyReference(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "policy uses raw JSON",
			Content: `
resource "aws_iam_policy" "policy" {
	name = "test_policy"
	policy = <<EOF
	{
		"Version": "2012-10-17",
		"Statement": [
		{
			"Effect": "Allow",
			"Resource": "arn:aws:s3:::<bucketname>/*""
		}
		]
	}
EOF
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMPolicyUsePolicyReferenceRule(),
					Message: "The policy does not reference an object",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 2},
						End:      hcl.Pos{Line: 14, Column: 4},
					},
				},
			},
		},
		{
			Name: "policy uses jsonencode",
			Content: `
resource "aws_iam_policy" "policy" {
	name   = "test_policy"
	policy = jsonencode()
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMPolicyUsePolicyReferenceRule(),
					Message: "The policy does not reference an object",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 2},
						End:      hcl.Pos{Line: 4, Column: 23},
					},
				},
			},
		},
		{
			Name: "policy uses data.aws_iam_policy_document",
			Content: `
resource "aws_iam_policy" "policy" {
	name   = "test_policy"
	policy = data.aws_iam_policy_document.document.json # (not shown)
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "policy uses file call",
			Content: `
resource "aws_iam_policy" "policy" {
	name   = "test_policy"
	policy = file("path/to/policy.json")
}
`,
			Expected: helper.Issues{},
		},

		{
			Name: "policy uses templatefile call",
			Content: `
resource "aws_iam_policy" "policy" {
	name   = "test_policy"
	policy = templatefile("path/to/policy.json.tmpl")
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsIAMPolicyUsePolicyReferenceRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
