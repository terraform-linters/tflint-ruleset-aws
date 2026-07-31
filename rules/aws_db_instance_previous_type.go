package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/db_instance_types"
)

// AwsDBInstancePreviousTypeRule checks whether the resource uses previous generation instance type
type AwsDBInstancePreviousTypeRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsDBInstancePreviousTypeRule returns new rule with default attributes
func NewAwsDBInstancePreviousTypeRule() *AwsDBInstancePreviousTypeRule {
	return &AwsDBInstancePreviousTypeRule{
		resourceType:  "aws_db_instance",
		attributeName: "instance_class",
	}
}

// Name returns the rule name
func (r *AwsDBInstancePreviousTypeRule) Name() string {
	return "aws_db_instance_previous_type"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsDBInstancePreviousTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsDBInstancePreviousTypeRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsDBInstancePreviousTypeRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the resource's `instance_class` is included in the list of previous generation instance type
func (r *AwsDBInstancePreviousTypeRule) Check(runner tflint.Runner) error {
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

		err := runner.EvaluateExpr(attribute.Expr, func(instanceType string) error {
			if db_instance_types.PreviousGeneration(instanceType) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is previous generation instance type.", instanceType),
					attribute.Expr.Range(),
				)
			}
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
