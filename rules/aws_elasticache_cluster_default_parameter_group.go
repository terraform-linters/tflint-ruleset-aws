package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsElastiCacheClusterDefaultParameterGroupRule checks whether the cluster use default parameter group
type AwsElastiCacheClusterDefaultParameterGroupRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsElastiCacheClusterDefaultParameterGroupRule returns new rule with default attributes
func NewAwsElastiCacheClusterDefaultParameterGroupRule() *AwsElastiCacheClusterDefaultParameterGroupRule {
	return &AwsElastiCacheClusterDefaultParameterGroupRule{
		resourceType:  "aws_elasticache_cluster",
		attributeName: "parameter_group_name",
	}
}

// Name returns the rule name
func (r *AwsElastiCacheClusterDefaultParameterGroupRule) Name() string {
	return "aws_elasticache_cluster_default_parameter_group"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsElastiCacheClusterDefaultParameterGroupRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsElastiCacheClusterDefaultParameterGroupRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsElastiCacheClusterDefaultParameterGroupRule) Link() string {
	return project.ReferenceLink(r.Name())
}

var defaultElastiCacheParameterGroupRegexp = regexp.MustCompile("^default")

// Check checks the parameter group name starts with `default`
func (r *AwsElastiCacheClusterDefaultParameterGroupRule) Check(runner tflint.Runner) error {
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

		var parameterGroup string
		err := runner.EvaluateExpr(attribute.Expr, &parameterGroup, nil)

		err = runner.EnsureNoError(err, func() error {
			if defaultElastiCacheParameterGroupRegexp.Match([]byte(parameterGroup)) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is default parameter group. You cannot edit it.", parameterGroup),
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
