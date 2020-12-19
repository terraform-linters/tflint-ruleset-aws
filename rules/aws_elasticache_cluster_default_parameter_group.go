package rules

import (
	"fmt"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsElastiCacheClusterDefaultParameterGroupRule checks whether the cluster use default parameter group
type AwsElastiCacheClusterDefaultParameterGroupRule struct {
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
func (r *AwsElastiCacheClusterDefaultParameterGroupRule) Severity() string {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsElastiCacheClusterDefaultParameterGroupRule) Link() string {
	return project.ReferenceLink(r.Name())
}

var defaultElastiCacheParameterGroupRegexp = regexp.MustCompile("^default")

// Check checks the parameter group name starts with `default`
func (r *AwsElastiCacheClusterDefaultParameterGroupRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var parameterGroup string
		err := runner.EvaluateExpr(attribute.Expr, &parameterGroup, nil)

		return runner.EnsureNoError(err, func() error {
			if defaultElastiCacheParameterGroupRegexp.Match([]byte(parameterGroup)) {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("\"%s\" is default parameter group. You cannot edit it.", parameterGroup),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
