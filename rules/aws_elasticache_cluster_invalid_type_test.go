package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsElastiCacheClusterInvalidType(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "t2.micro is invalid",
			Content: `
resource "aws_elasticache_cluster" "redis" {
    node_type = "t2.micro"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsElastiCacheClusterInvalidTypeRule(),
					Message: "\"t2.micro\" is invalid node type.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 17},
						End:      hcl.Pos{Line: 3, Column: 27},
					},
				},
			},
		},
		{
			Name: "empty is invalid",
			Content: `
resource "aws_elasticache_cluster" "redis" {
    node_type = ""
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsElastiCacheClusterInvalidTypeRule(),
					Message: "\"\" is invalid node type.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 17},
						End:      hcl.Pos{Line: 3, Column: 19},
					},
				},
			},
		},
		{
			Name: "cache.t2.micro is valid",
			Content: `
resource "aws_elasticache_cluster" "redis" {
    node_type = "cache.t2.micro"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsElastiCacheClusterInvalidTypeRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
