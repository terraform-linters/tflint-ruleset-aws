package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsWriteOnlyAttribute(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "basic aws_secretsmanager_secret_version",
			Content: `
resource "aws_secretsmanager_secret_version" "test" {
  secret_string = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsWriteOnlyArgumentsRule(),
					Message: `"secret_string" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only argument "secret_string_wo".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 19},
						End:      hcl.Pos{Line: 3, Column: 25},
					},
				},
			},
			Fixed: `
resource "aws_secretsmanager_secret_version" "test" {
  secret_string_wo = "test"
}
`,
		},
		{
			Name: "everything is fine aws_secretsmanager_secret_version",
			Content: `
resource "aws_secretsmanager_secret_version" "test" {
  secret_string_wo = "test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "basic aws_rds_cluster",
			Content: `
resource "aws_rds_cluster" "test" {
  master_password = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsWriteOnlyArgumentsRule(),
					Message: `"master_password" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only argument "master_password_wo".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 21},
						End:      hcl.Pos{Line: 3, Column: 27},
					},
				},
			},
			Fixed: `
resource "aws_rds_cluster" "test" {
  master_password_wo = "test"
}
`,
		},
		{
			Name: "everything is fine aws_rds_cluster",
			Content: `
resource "aws_rds_cluster" "test" {
  master_password_wo = "test"
}
`,
			Expected: helper.Issues{},
		},

		{
			Name: "basic aws_db_instance",
			Content: `
resource "aws_db_instance" "test" {
  password = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsWriteOnlyArgumentsRule(),
					Message: `"password" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only argument "password_wo".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 14},
						End:      hcl.Pos{Line: 3, Column: 20},
					},
				},
			},
			Fixed: `
resource "aws_db_instance" "test" {
  password_wo = "test"
}
`,
		},
		{
			Name: "everything is fine aws_db_instance",
			Content: `
resource "aws_db_instance" "test" {
  password_wo = "test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "basic aws_redshift_cluster",
			Content: `
resource "aws_redshift_cluster" "test" {
  master_password = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsWriteOnlyArgumentsRule(),
					Message: `"master_password" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only argument "master_password_wo".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 21},
						End:      hcl.Pos{Line: 3, Column: 27},
					},
				},
			},
			Fixed: `
resource "aws_redshift_cluster" "test" {
  master_password_wo = "test"
}
`,
		},
		{
			Name: "everything is fine aws_redshift_cluster",
			Content: `
resource "aws_redshift_cluster" "test" {
  master_password_wo = "test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "basic aws_docdb_cluster",
			Content: `
resource "aws_docdb_cluster" "test" {
  master_password = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsWriteOnlyArgumentsRule(),
					Message: `"master_password" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only argument "master_password_wo".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 21},
						End:      hcl.Pos{Line: 3, Column: 27},
					},
				},
			},
			Fixed: `
resource "aws_docdb_cluster" "test" {
  master_password_wo = "test"
}
`,
		},
		{
			Name: "everything is fine aws_docdb_cluster",
			Content: `
resource "aws_docdb_cluster" "test" {
  master_password_wo = "test"
}
`,
			Expected: helper.Issues{},
		},

		{
			Name: "basic aws_redshiftserverless_namespace",
			Content: `
resource "aws_redshiftserverless_namespace" "test" {
  admin_password = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsWriteOnlyArgumentsRule(),
					Message: `"admin_password" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only argument "admin_password_wo".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 20},
						End:      hcl.Pos{Line: 3, Column: 26},
					},
				},
			},
			Fixed: `
resource "aws_redshiftserverless_namespace" "test" {
  admin_password_wo = "test"
}
`,
		},
		{
			Name: "everything is fine aws_redshiftserverless_namespace",
			Content: `
resource "aws_redshiftserverless_namespace" "test" {
  admin_password_wo = "test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "basic aws_ssm_parameter",
			Content: `
resource "aws_ssm_parameter" "test" {
  value = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsWriteOnlyArgumentsRule(),
					Message: `"value" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only argument "value_wo".`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 17},
					},
				},
			},

			Fixed: `
resource "aws_ssm_parameter" "test" {
  value_wo = "test"
}
`,
		},
		{
			Name: "everything is fine aws_ssm_parameter",
			Content: `
resource "aws_ssm_parameter" "test" {
  value_wo = "test"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsWriteOnlyArgumentsRule()

	for _, tc := range cases {
		filename := "resource.tf"
		runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}
		helper.AssertIssues(t, tc.Expected, runner.Issues)

		want := map[string]string{}
		if tc.Fixed != "" {
			want[filename] = tc.Fixed
		}
		helper.AssertChanges(t, want, runner.Changes())
	}
}
