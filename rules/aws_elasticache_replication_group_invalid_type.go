package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsElastiCacheReplicationGroupInvalidTypeRule checks whether "aws_elasticache_replication_group" has invalid node type.
type AwsElastiCacheReplicationGroupInvalidTypeRule struct {
	resourceType  string
	attributeName string
}

// NewAwsElastiCacheReplicationGroupInvalidTypeRule returns new rule with default attributes
func NewAwsElastiCacheReplicationGroupInvalidTypeRule() *AwsElastiCacheReplicationGroupInvalidTypeRule {
	return &AwsElastiCacheReplicationGroupInvalidTypeRule{
		resourceType:  "aws_elasticache_replication_group",
		attributeName: "node_type",
	}
}

// Name returns the rule name
func (r *AwsElastiCacheReplicationGroupInvalidTypeRule) Name() string {
	return "aws_elasticache_replication_group_invalid_type"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsElastiCacheReplicationGroupInvalidTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsElastiCacheReplicationGroupInvalidTypeRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsElastiCacheReplicationGroupInvalidTypeRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether "aws_elasticache_replication_group" has invalid node type.
func (r *AwsElastiCacheReplicationGroupInvalidTypeRule) Check(runner tflint.Runner) error {
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
