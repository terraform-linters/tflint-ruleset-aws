package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstanceInvalidEngine(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "aurora-postgres is invalid",
			Content: `
resource "aws_db_instance" "aurora_postgresql" {
    engine = "aurora-postgres"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceInvalidEngineRule(),
					Message: "\"aurora-postgres\" is invalid engine.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 14},
						End:      hcl.Pos{Line: 3, Column: 31},
					},
				},
			},
		},
		{
			Name: "aurora-postgresql is valid",
			Content: `
resource "aws_db_instance" "aurora_postgresql" {
    engine = "aurora-postgresql"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "aurora (standalone) is invalid",
			Content: `
resource "aws_db_instance" "aurora" {
    engine = "aurora"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceInvalidEngineRule(),
					Message: "\"aurora\" is invalid engine.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 14},
						End:      hcl.Pos{Line: 3, Column: 22},
					},
				},
			},
		},
		{
			Name: "oracle-se1 is invalid (deprecated)",
			Content: `
resource "aws_db_instance" "oracle" {
    engine = "oracle-se1"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceInvalidEngineRule(),
					Message: "\"oracle-se1\" is invalid engine.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 14},
						End:      hcl.Pos{Line: 3, Column: 26},
					},
				},
			},
		},
		{
			Name: "oracle-se is invalid (deprecated)",
			Content: `
resource "aws_db_instance" "oracle" {
    engine = "oracle-se"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceInvalidEngineRule(),
					Message: "\"oracle-se\" is invalid engine.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 14},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
			},
		},
		{
			Name: "oracle-se2 is valid",
			Content: `
resource "aws_db_instance" "oracle" {
    engine = "oracle-se2"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDBInstanceInvalidEngineRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
