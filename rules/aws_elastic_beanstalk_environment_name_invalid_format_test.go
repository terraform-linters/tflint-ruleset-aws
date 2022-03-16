package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsElasticBeanstalkEnvironmentNameInvalidFormat(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "tf-test-name dash valid",
			Content: `
resource "aws_elastic_beanstalk_environment" "tfenvtest" {
	name                = "tf-test-name"
	application         = "tf-test-name"
	solution_stack_name = "64bit Amazon Linux 2015.03 v2.0.3 running Go 1.4"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "underscores invalid",
			Content: `
resource "aws_elastic_beanstalk_environment" "tfenvtest" {
	name                = "tf_test_name"
	application         = "tf-test-name"
	solution_stack_name = "64bit Amazon Linux 2015.03 v2.0.3 running Go 1.4"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsElasticBeanstalkEnvironmentNameInvalidFormatRule(),
					Message: "tf_test_name does not match constraint: must contain only letters, digits, and " +
					"the dash character and may not start or end with a dash (^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$)",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 24},
						End:      hcl.Pos{Line: 3, Column: 38},
					},
				},
			},
		},
		{
			Name: "end with dash invalid",
			Content: `
resource "aws_elastic_beanstalk_environment" "tfenvtest" {
	name                = "tf-test-name-"
	application         = "tf-test-name"
	solution_stack_name = "64bit Amazon Linux 2015.03 v2.0.3 running Go 1.4"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsElasticBeanstalkEnvironmentNameInvalidFormatRule(),
					Message: "tf-test-name- does not match constraint: must contain only letters, digits, and " +
					"the dash character and may not start or end with a dash (^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$)",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 24},
						End:      hcl.Pos{Line: 3, Column: 39},
					},
				},
			},
		},
	}

	rule := NewAwsElasticBeanstalkEnvironmentNameInvalidFormatRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
