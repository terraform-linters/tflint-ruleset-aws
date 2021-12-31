package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsDBInstanceDefaultParameterGroupRule checks whether the db instance use default parameter group
type AwsDBInstanceDefaultParameterGroupRule struct {
	tflint.DefaultRule

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
func (r *AwsDBInstanceDefaultParameterGroupRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsDBInstanceDefaultParameterGroupRule) Link() string {
	return project.ReferenceLink(r.Name())
}

var defaultDBParameterGroupRegexp = regexp.MustCompile("^default")

// Check checks the parameter group name starts with `default`
func (r *AwsDBInstanceDefaultParameterGroupRule) Check(runner tflint.Runner) error {
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

		var name string
		err := runner.EvaluateExpr(attribute.Expr, &name, nil)

		err = runner.EnsureNoError(err, func() error {
			if defaultDBParameterGroupRegexp.Match([]byte(name)) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is default parameter group. You cannot edit it.", name),
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
