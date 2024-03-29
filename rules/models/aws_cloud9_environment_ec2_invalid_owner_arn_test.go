// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"testing"
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsCloud9EnvironmentEc2InvalidOwnerArnRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "It includes invalid characters",
			Content: `
resource "aws_cloud9_environment_ec2" "foo" {
	owner_arn = "arn:aws:elasticbeanstalk:us-east-1:123456789012:environment/My App/MyEnvironment"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCloud9EnvironmentEc2InvalidOwnerArnRule(),
					Message: fmt.Sprintf(`"%s" does not match valid pattern %s`, truncateLongMessage("arn:aws:elasticbeanstalk:us-east-1:123456789012:environment/My App/MyEnvironment"), `^arn:(aws|aws-cn|aws-us-gov|aws-iso|aws-iso-b):(iam|sts)::\d+:(root|(user\/[\w+=/:,.@-]{1,64}|federated-user\/[\w+=/:,.@-]{2,32}|assumed-role\/[\w+=:,.@-]{1,64}\/[\w+=,.@-]{1,64}))$`),
				},
			},
		},
		{
			Name: "It is valid",
			Content: `
resource "aws_cloud9_environment_ec2" "foo" {
	owner_arn = "arn:aws:iam::123456789012:user/David"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsCloud9EnvironmentEc2InvalidOwnerArnRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
	}
}
