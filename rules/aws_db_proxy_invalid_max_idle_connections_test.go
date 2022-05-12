package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBProxyInvalidMaxIdleConnections(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "equal",
			Content: `
resource "aws_db_proxy_default_target_group" "example" {
  connection_pool_config {
    max_connections_percent      = 50
    max_idle_connections_percent = 50
  }
}
                        `,
			Expected: helper.Issues{},
		},
		{
			Name: "max_idle_connections_percent grater than max_connections_percent",
			Content: `
resource "aws_db_proxy_default_target_group" "example" {
  connection_pool_config {
    max_connections_percent      = 10
    max_idle_connections_percent = 50
  }
}
                        `,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBProxyInvalidMaxIdleConnectionsRule(),
					Message: "max_idle_connections_percent 50 must be less than max_connections_percent 10",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 5},
						End:      hcl.Pos{Line: 5, Column: 38},
					},
				},
			},
		},
	}

	rule := NewAwsDBProxyInvalidMaxIdleConnectionsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
