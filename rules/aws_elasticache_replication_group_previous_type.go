package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsElastiCacheReplicationGroupPreviousTypeRule checks whether the resource uses previous generation node type
type AwsElastiCacheReplicationGroupPreviousTypeRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsElastiCacheReplicationGroupPreviousTypeRule returns new rule with default attributes
func NewAwsElastiCacheReplicationGroupPreviousTypeRule() *AwsElastiCacheReplicationGroupPreviousTypeRule {
	return &AwsElastiCacheReplicationGroupPreviousTypeRule{
		resourceType:  "aws_elasticache_replication_group",
		attributeName: "node_type",
	}
}

// Name returns the rule name
func (r *AwsElastiCacheReplicationGroupPreviousTypeRule) Name() string {
	return "aws_elasticache_replication_group_previous_type"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsElastiCacheReplicationGroupPreviousTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsElastiCacheReplicationGroupPreviousTypeRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsElastiCacheReplicationGroupPreviousTypeRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the resource's `node_type` is included in the list of previous generation node type
func (r *AwsElastiCacheReplicationGroupPreviousTypeRule) Check(runner tflint.Runner) error {
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
			if previousElastiCacheNodeTypes[strings.Split(nodeType, ".")[1]] {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is previous generation node type.", nodeType),
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
