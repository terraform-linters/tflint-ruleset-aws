package rules

import (
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsSecurityGroupEmbeddedIngressRule checks to make sure embedded ingress rules are not set, prefering aws_security_group_rule
type AwsSecurityGroupEmbeddedIngressRule struct {
	resourceType  string
	attributeName string
}

// NewAwsSecurityGroupEmbeddedIngressRule returns new rule with default attributes
func NewAwsSecurityGroupEmbeddedIngressRule() *AwsSecurityGroupEmbeddedIngressRule {
	return &AwsSecurityGroupEmbeddedIngressRule{
		resourceType:  "aws_security_group",
		attributeName: "ingress",
	}
}

// Name returns the rule name
func (r *AwsSecurityGroupEmbeddedIngressRule) Name() string {
	return "aws_security_group_embedded_ingress"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupEmbeddedIngressRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSecurityGroupEmbeddedIngressRule) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsSecurityGroupEmbeddedIngressRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Checks to see if ingress block exists 
func (r *AwsSecurityGroupEmbeddedIngressRule) Check(runner tflint.Runner) error {
	// TODO: Write the implementation here. See this documentation for what tflint.Runner can do.
	//       https://pkg.go.dev/github.com/terraform-linters/tflint-plugin-sdk/tflint#Runner

	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		return runner.EnsureNoError(err, func() error {
			if val == "" {
				runner.EmitIssueOnExpr(
					r,
					"TODO",
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
