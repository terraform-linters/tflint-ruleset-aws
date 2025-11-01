package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstanceInvalidType(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "m4.2xlarge is invalid",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "m4.2xlarge"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceInvalidTypeRule(),
					Message: "\"m4.2xlarge\" is invalid instance type.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 22},
						End:      hcl.Pos{Line: 3, Column: 34},
					},
				},
			},
		},
		{
			Name: "db.m4.2xlarge is valid",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "db.m4.2xlarge"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "db.c5.large is valid",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "db.c5.large"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "db.c6g.16xlarge is valid",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "db.c6g.16xlarge"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "db.c7g.medium is valid",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "db.c7g.medium"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "db.i3.xlarge is valid",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "db.i3.xlarge"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "db.i4i.32xlarge is valid",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "db.i4i.32xlarge"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "db.c6gd.metal is valid",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "db.c6gd.metal"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "db.x2idn.16xlarge is valid (typo fixed)",
			Content: `
resource "aws_db_instance" "mysql" {
    instance_class = "db.x2idn.16xlarge"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDBInstanceInvalidTypeRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
