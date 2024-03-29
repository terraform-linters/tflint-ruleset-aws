// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsCognitoUserPoolUICustomizationInvalidCssRule checks the pattern is valid
type AwsCognitoUserPoolUICustomizationInvalidCssRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	max           int
}

// NewAwsCognitoUserPoolUICustomizationInvalidCssRule returns new rule with default attributes
func NewAwsCognitoUserPoolUICustomizationInvalidCssRule() *AwsCognitoUserPoolUICustomizationInvalidCssRule {
	return &AwsCognitoUserPoolUICustomizationInvalidCssRule{
		resourceType:  "aws_cognito_user_pool_ui_customization",
		attributeName: "css",
		max:           131072,
	}
}

// Name returns the rule name
func (r *AwsCognitoUserPoolUICustomizationInvalidCssRule) Name() string {
	return "aws_cognito_user_pool_ui_customization_invalid_css"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsCognitoUserPoolUICustomizationInvalidCssRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsCognitoUserPoolUICustomizationInvalidCssRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsCognitoUserPoolUICustomizationInvalidCssRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsCognitoUserPoolUICustomizationInvalidCssRule) Check(runner tflint.Runner) error {
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
					"css must be 131072 characters or less",
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
