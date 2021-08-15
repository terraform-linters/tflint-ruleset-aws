package rules

import (
	"math/rand"
	"testing"
	"time"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func AwsIAMGroupPolicyTooLongRandSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Test_AwsIAMGroupPolicyTooLong(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "policy is too long",
			Content: `
resource "aws_iam_group_policy" "policy" {
	name        = "test_policy"
	group = "test_group"
	policy = <<EOF
	{
		"Version": "2012-10-17",
		"Statement": [
		{
			"Action": [
				` + AwsIAMGroupPolicyTooLongRandSeq(5120) + `
			],
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
					Rule:    NewAwsIAMGroupPolicyTooLongRule(),
					Message: "The policy length is 5231 characters and is limited to 5120 characters.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 11},
						End:      hcl.Pos{Line: 18, Column: 4},
					},
				},
			},
		},
	}

	rule := NewAwsIAMGroupPolicyTooLongRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
