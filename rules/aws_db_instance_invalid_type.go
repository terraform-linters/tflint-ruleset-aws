package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/db_instance_types"
)

// AwsDBInstanceInvalidTypeRule checks whether "aws_db_instance" has invalid instance type.
type AwsDBInstanceInvalidTypeRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsDBInstanceInvalidTypeRule returns new rule with default attributes
func NewAwsDBInstanceInvalidTypeRule() *AwsDBInstanceInvalidTypeRule {
	return &AwsDBInstanceInvalidTypeRule{
		resourceType:  "aws_db_instance",
		attributeName: "instance_class",
	}
}

// Name returns the rule name
func (r *AwsDBInstanceInvalidTypeRule) Name() string {
	return "aws_db_instance_invalid_type"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsDBInstanceInvalidTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsDBInstanceInvalidTypeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsDBInstanceInvalidTypeRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether "aws_db_instance" has invalid instance type.
func (r *AwsDBInstanceInvalidTypeRule) Check(runner tflint.Runner) error {
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
			if !db_instance_types.Valid(instanceType) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is invalid instance type.", instanceType),
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
