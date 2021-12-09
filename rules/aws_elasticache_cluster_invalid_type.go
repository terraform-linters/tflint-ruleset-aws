package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsElastiCacheClusterInvalidTypeRule checks whether "aws_elasticache_cluster" has invalid node type.
type AwsElastiCacheClusterInvalidTypeRule struct {
	resourceType  string
	attributeName string
}

// NewAwsElastiCacheClusterInvalidTypeRule returns new rule with default attributes
func NewAwsElastiCacheClusterInvalidTypeRule() *AwsElastiCacheClusterInvalidTypeRule {
	return &AwsElastiCacheClusterInvalidTypeRule{
		resourceType:  "aws_elasticache_cluster",
		attributeName: "node_type",
	}
}

// Name returns the rule name
func (r *AwsElastiCacheClusterInvalidTypeRule) Name() string {
	return "aws_elasticache_cluster_invalid_type"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsElastiCacheClusterInvalidTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsElastiCacheClusterInvalidTypeRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsElastiCacheClusterInvalidTypeRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether "aws_elasticache_cluster" has invalid node type.
func (r *AwsElastiCacheClusterInvalidTypeRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var nodeType string
		err := runner.EvaluateExpr(attribute.Expr, &nodeType, nil)

		return runner.EnsureNoError(err, func() error {
			if !validElastiCacheNodeTypes[nodeType] {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("\"%s\" is invalid node type.", nodeType),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
