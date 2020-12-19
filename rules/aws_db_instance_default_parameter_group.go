package rules

import (
	"fmt"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsDBInstanceDefaultParameterGroupRule checks whether the db instance use default parameter group
type AwsDBInstanceDefaultParameterGroupRule struct {
	resourceType  string
	attributeName string
}

// NewAwsDBInstanceDefaultParameterGroupRule returns new rule with default attributes
func NewAwsDBInstanceDefaultParameterGroupRule() *AwsDBInstanceDefaultParameterGroupRule {
	return &AwsDBInstanceDefaultParameterGroupRule{
		resourceType:  "aws_db_instance",
		attributeName: "parameter_group_name",
	}
}

// Name returns the rule name
func (r *AwsDBInstanceDefaultParameterGroupRule) Name() string {
	return "aws_db_instance_default_parameter_group"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsDBInstanceDefaultParameterGroupRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsDBInstanceDefaultParameterGroupRule) Severity() string {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsDBInstanceDefaultParameterGroupRule) Link() string {
	return project.ReferenceLink(r.Name())
}

var defaultDBParameterGroupRegexp = regexp.MustCompile("^default")

// Check checks the parameter group name starts with `default`
func (r *AwsDBInstanceDefaultParameterGroupRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var name string
		err := runner.EvaluateExpr(attribute.Expr, &name, nil)

		return runner.EnsureNoError(err, func() error {
			if defaultDBParameterGroupRegexp.Match([]byte(name)) {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("\"%s\" is default parameter group. You cannot edit it.", name),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
