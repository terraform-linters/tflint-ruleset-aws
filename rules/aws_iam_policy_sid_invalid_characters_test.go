package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsIAMPolicySidInvalidCharacters(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "single statement with invalid",
			Content: `
resource "aws_iam_policy" "policy" {
	name = "test_policy"
	role = "test_role"
	policy = <<-EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "This contains invalid-characters.",
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::<bucketname>/*"
    }
  ]
}
EOF
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMPolicySidInvalidCharactersRule(),
					Message: `The policy's sid ("This contains invalid-characters.") does not match "^[a-zA-Z0-9]+$".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 11},
						End:      hcl.Pos{Line: 19, Column: 4},
					},
				},
			},
		},
		{
			Name: "two statements second invalid",
			Content: `
resource "aws_iam_policy" "policy2" {
	name = "test_policy"
	role = "test_role"
	policy = <<-EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ThisIsAValidSid",
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::<bucketname>/*"
    },
    {
      "Sid": "This contains invalid-characters.",
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::<bucketname>/*"
    }
  ]
}
EOF
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMPolicySidInvalidCharactersRule(),
					Message: `The policy's sid ("This contains invalid-characters.") does not match "^[a-zA-Z0-9]+$".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 11},
						End:      hcl.Pos{Line: 27, Column: 4},
					},
				},
			},
		},
		{
			Name: "No Sid",
			Content: `
resource "aws_iam_policy" "policy" {
	name = "test_policy"
	role = "test_role"
	policy = <<-EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::<bucketname>/*"
    }
  ]
}
EOF
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "Single Statement",
			Content: `
resource "aws_iam_policy" "policy" {
	name = "test_policy"
	role = "test_role"
	policy = <<-EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Sid": "ThisIsAValidSid",
    "Action": [
      "ec2:Describe*"
    ],
    "Effect": "Allow",
    "Resource": "arn:aws:s3:::<bucketname>/*"
  }
}
EOF
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsIAMPolicySidInvalidCharactersRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
