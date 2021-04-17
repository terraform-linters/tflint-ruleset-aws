package rules

import (
	"fmt"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayModelInvalidNameRule checks the name is alphanumeric
type AwsAPIGatewayModelInvalidNameRule struct {
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
func (r *AwsAPIGatewayModelInvalidNameRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsAPIGatewayModelInvalidNameRule) Link() string {
	return ""
}

// Check checks the name attributes is matched with ^[a-zA-Z0-9]+$ regexp
func (r *AwsAPIGatewayModelInvalidNameRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		return runner.EnsureNoError(err, func() error {
			if !r.pattern.MatchString(val) {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf(`%s does not match valid pattern %s`, val, `^[a-zA-Z0-9]+$`),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
