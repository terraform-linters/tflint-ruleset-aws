package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsElastiCacheReplicationGroupInvalidTypeRule checks whether "aws_elasticache_replication_group" has invalid node type.
type AwsElastiCacheReplicationGroupInvalidTypeRule struct {
	tflint.DefaultRule

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
func (r *AwsElastiCacheReplicationGroupInvalidTypeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsElastiCacheReplicationGroupInvalidTypeRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether "aws_elasticache_replication_group" has invalid node type.
func (r *AwsElastiCacheReplicationGroupInvalidTypeRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.attributeName}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		var nodeType string
		err := runner.EvaluateExpr(attribute.Expr, &nodeType, nil)

		err = runner.EnsureNoError(err, func() error {
			if !validElastiCacheNodeTypes[nodeType] {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is invalid node type.", nodeType),
					attribute.Expr.Range(),
				)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
