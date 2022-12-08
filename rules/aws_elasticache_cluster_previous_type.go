package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsElastiCacheClusterPreviousTypeRule checks whether the resource uses previous generation node type
type AwsElastiCacheClusterPreviousTypeRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsElastiCacheClusterPreviousTypeRule returns new rule with default attributes
func NewAwsElastiCacheClusterPreviousTypeRule() *AwsElastiCacheClusterPreviousTypeRule {
	return &AwsElastiCacheClusterPreviousTypeRule{
		resourceType:  "aws_elasticache_cluster",
		attributeName: "node_type",
	}
}

// Name returns the rule name
func (r *AwsElastiCacheClusterPreviousTypeRule) Name() string {
	return "aws_elasticache_cluster_previous_type"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsElastiCacheClusterPreviousTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsElastiCacheClusterPreviousTypeRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsElastiCacheClusterPreviousTypeRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the resource's `node_type` is included in the list of previous generation node type
func (r *AwsElastiCacheClusterPreviousTypeRule) Check(runner tflint.Runner) error {
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
			parts := strings.Split(nodeType, ".")
			if len(parts) != 3 {
				return nil
			}

			if previousElastiCacheNodeTypes[parts[1]] {
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
