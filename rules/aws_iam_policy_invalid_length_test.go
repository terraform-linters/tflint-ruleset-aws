package rules

import (
	"testing"
	"time"
	"math/rand"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func Test_AwsIAMPolicyInvalidLength(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "policy is too long",
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
				` + randSeq(1894) + `
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
					Rule:    NewAwsIAMPolicyInvalidLengthRule(),
					Message: "The policy length is 2050 characters and is limited to 2048 characters.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 11},
						End:      hcl.Pos{Line: 19, Column: 4},
					},
				},
			},
		},
	}

	rule := NewAwsIAMPolicyInvalidLengthRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
