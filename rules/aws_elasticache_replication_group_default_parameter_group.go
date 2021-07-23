package rules

import (
	"fmt"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsElastiCacheReplicationGroupDefaultParameterGroupRule checks whether the cluster use default parameter group
type AwsElastiCacheReplicationGroupDefaultParameterGroupRule struct {
	resourceType  string
	attributeName string
}

// NewAwsElastiCacheReplicationGroupDefaultParameterGroupRule returns new rule with default attributes
func NewAwsElastiCacheReplicationGroupDefaultParameterGroupRule() *AwsElastiCacheReplicationGroupDefaultParameterGroupRule {
	return &AwsElastiCacheReplicationGroupDefaultParameterGroupRule{
		resourceType:  "aws_elasticache_replication_group",
		attributeName: "parameter_group_name",
	}
}

// Name returns the rule name
func (r *AwsElastiCacheReplicationGroupDefaultParameterGroupRule) Name() string {
	return "aws_elasticache_replication_group_default_parameter_group"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsElastiCacheReplicationGroupDefaultParameterGroupRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsElastiCacheReplicationGroupDefaultParameterGroupRule) Severity() string {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsElastiCacheReplicationGroupDefaultParameterGroupRule) Link() string {
	return project.ReferenceLink(r.Name())
}

var defaultElastiCacheReplicationParameterGroupRegexp = regexp.MustCompile("^default")

// Check checks the parameter group name starts with `default`
func (r *AwsElastiCacheReplicationGroupDefaultParameterGroupRule) Check(runner tflint.Runner) error {
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
