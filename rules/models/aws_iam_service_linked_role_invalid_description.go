// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsIAMServiceLinkedRoleInvalidDescriptionRule checks the pattern is valid
type AwsIAMServiceLinkedRoleInvalidDescriptionRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	max           int
	pattern       *regexp.Regexp
}

// NewAwsIAMServiceLinkedRoleInvalidDescriptionRule returns new rule with default attributes
func NewAwsIAMServiceLinkedRoleInvalidDescriptionRule() *AwsIAMServiceLinkedRoleInvalidDescriptionRule {
	return &AwsIAMServiceLinkedRoleInvalidDescriptionRule{
		resourceType:  "aws_iam_service_linked_role",
		attributeName: "description",
		max:           1000,
		pattern:       regexp.MustCompile(`^[\x{0009}\x{000A}\x{000D}\x{0020}-\x{007E}\x{00A1}-\x{00FF}]*$`),
	}
}

// Name returns the rule name
func (r *AwsIAMServiceLinkedRoleInvalidDescriptionRule) Name() string {
	return "aws_iam_service_linked_role_invalid_description"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMServiceLinkedRoleInvalidDescriptionRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsIAMServiceLinkedRoleInvalidDescriptionRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsIAMServiceLinkedRoleInvalidDescriptionRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsIAMServiceLinkedRoleInvalidDescriptionRule) Check(runner tflint.Runner) error {
	logger.Trace("Check `%s` rule", r.Name())

	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		err := runner.EvaluateExpr(attribute.Expr, func (val string) error {
			if len(val) > r.max {
				runner.EmitIssue(
					r,
					"description must be 1000 characters or less",
					attribute.Expr.Range(),
				)
			}
			if !r.pattern.MatchString(val) {
				runner.EmitIssue(
					r,
					fmt.Sprintf(`"%s" does not match valid pattern %s`, truncateLongMessage(val), `^[\x{0009}\x{000A}\x{000D}\x{0020}-\x{007E}\x{00A1}-\x{00FF}]*$`),
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
