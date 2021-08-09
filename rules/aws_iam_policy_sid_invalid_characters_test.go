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
					Message: "The policy's sid contains invalid characters.",
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
					Message: "The policy's sid contains invalid characters.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 11},
						End:      hcl.Pos{Line: 27, Column: 4},
					},
				},
			},
		},
	}

	rule := NewAwsIAMPolicySidInvalidCharactersRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
