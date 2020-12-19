package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstancePreviousType(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "db.t1.micro is previous type",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "db.t1.micro"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstancePreviousTypeRule(),
					Message: "\"db.t1.micro\" is previous generation instance type.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 22},
						End:      hcl.Pos{Line: 3, Column: 35},
					},
				},
			},
		},
		{
			Name: "db.t2.micro is not previous type",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "db.t2.micro"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDBInstancePreviousTypeRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
