package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayModelInvalidNameRule checks the name is alphanumeric
type AwsAPIGatewayModelInvalidNameRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	pattern       *regexp.Regexp
}

// NewAwsAPIGatewayModelInvalidNameRule returns new rule with default attributes
func NewAwsAPIGatewayModelInvalidNameRule() *AwsAPIGatewayModelInvalidNameRule {
	return &AwsAPIGatewayModelInvalidNameRule{
		resourceType:  "aws_api_gateway_model",
		attributeName: "name",
		pattern:       regexp.MustCompile("^[a-zA-Z0-9]+$"),
	}
}

// Name returns the rule name
func (r *AwsAPIGatewayModelInvalidNameRule) Name() string {
	return "aws_api_gateway_model_invalid_name"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAPIGatewayModelInvalidNameRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAPIGatewayModelInvalidNameRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsAPIGatewayModelInvalidNameRule) Link() string {
	return ""
}

// Check checks the name attributes is matched with ^[a-zA-Z0-9]+$ regexp
func (r *AwsAPIGatewayModelInvalidNameRule) Check(runner tflint.Runner) error {
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

		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		err = runner.EnsureNoError(err, func() error {
			if !r.pattern.MatchString(val) {
				runner.EmitIssue(
					r,
					fmt.Sprintf(`%s does not match valid pattern %s`, val, `^[a-zA-Z0-9]+$`),
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
