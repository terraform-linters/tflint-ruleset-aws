package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsElastiCacheReplicationGroupDefaultParameterGroup(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "default.redis3.2 is default parameter group",
			Content: `
resource "aws_elasticache_replication_group" "cache" {
    parameter_group_name = "default.redis3.2"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsElastiCacheReplicationGroupDefaultParameterGroupRule(),
					Message: "\"default.redis3.2\" is default parameter group. You cannot edit it.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 28},
						End:      hcl.Pos{Line: 3, Column: 46},
					},
				},
			},
		},
		{
			Name: "application3.2 is not default parameter group",
			Content: `
resource "aws_elasticache_replication_group" "cache" {
    parameter_group_name = "application3.2"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsElastiCacheReplicationGroupDefaultParameterGroupRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
