package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsElastiCacheReplicationGroupPreviousType(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "cache.t1.micro is previous type",
			Content: `
resource "aws_elasticache_replication_group" "redis" {
    node_type = "cache.t1.micro"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsElastiCacheReplicationGroupPreviousTypeRule(),
					Message: "\"cache.t1.micro\" is previous generation node type.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 17},
						End:      hcl.Pos{Line: 3, Column: 33},
					},
				},
			},
		},
		{
			Name: "cache.t2.micro is not previous type",
			Content: `
resource "aws_elasticache_replication_group" "redis" {
    node_type = "cache.t2.micro"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "test is not previous type",
			Content: `
resource "aws_elasticache_replication_group" "redis" {
    node_type = "test"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsElastiCacheReplicationGroupPreviousTypeRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
