package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsIamPolicyGovFriendlyArns(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "Wanted arn: arn:aws-us-gov, found arn: arn:aws",
			Content: `
resource "aws_iam_policy" "policy" {
  name        = "test_policy"
  path        = "/"
  description = "My test policy"
  policy = <<EOF
  {
    "Version": "2012-10-17",
    "Statement": [
	  {
	    "Action": [
		  "ec2:Describe*"
	    ],
	    "Effect": "Allow",
	    "Resource": "arn:aws:s3:::<bucketname>/*""
	  }
    ]
  }
EOF
}`,
			Config: `
rule "aws_iam_policy_gov_friendly_arns" {
  enabled = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMPolicyGovFriendlyArnsRule(),
					Message: "ARN detected in IAM policy that could potentially fail in AWS GovCloud due to resource pattern: arn:aws:.*",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 6, Column: 12},
						End:      hcl.Pos{Line: 19, Column: 4},
					},
				},
			},
		},
		{
			Name: "Wanted arn: arn:aws-us-gov, found arn: arn:aws-us-gov",
			Content: `
resource "aws_iam_policy" "policy" {
  name        = "test_policy"
  path        = "/"
  description = "My test policy"
  policy = <<EOF
  {
    "Version": "2012-10-17",
    "Statement": [
	  {
	    "Action": [
		  "ec2:Describe*"
	    ],
	    "Effect": "Allow",
	    "Resource": "arn:aws-us-gov:s3:::<bucketname>/*""
	  }
    ]
  }
EOF
}`,
			Config: `
rule "aws_iam_policy_gov_friendly_arns" {
  enabled = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Wanted arn *, found:  arn *",
			Content: `
resource "aws_iam_policy" "policy" {
  name        = "test_policy"
  path        = "/"
  description = "My test policy"
  policy = <<EOF
  {
    "Version": "2012-10-17",
    "Statement": [
	  {
	    "Action": [
		  "ec2:Describe*"
	    ],
	    "Effect": "Allow",
	    "Resource": "*"
	  }
    ]
  }
EOF
}`,
			Config: `
rule "aws_iam_policy_gov_friendly_arns" {
  enabled = true
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsIAMPolicyGovFriendlyArnsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"module.tf": tc.Content, ".tflint.hcl": tc.Config})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
