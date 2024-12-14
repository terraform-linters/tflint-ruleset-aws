package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsSecurityGroupRuleDeprecatedRule checks that aws_security_group_rule is not used
type AwsSecurityGroupRuleDeprecatedRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsSecurityGroupRuleDeprecatedRule returns new rule with default attributes
func NewAwsSecurityGroupRuleDeprecatedRule() *AwsSecurityGroupRuleDeprecatedRule {
	return &AwsSecurityGroupRuleDeprecatedRule{
		resourceType:  "aws_security_group_rule",
		attributeName: "security_group_id",
	}
}

// Name returns the rule name
func (r *AwsSecurityGroupRuleDeprecatedRule) Name() string {
	return "aws_security_group_rule_deprecated"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupRuleDeprecatedRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsSecurityGroupRuleDeprecatedRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsSecurityGroupRuleDeprecatedRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check that aws_security_group_rule resource is not used
func (r *AwsSecurityGroupRuleDeprecatedRule) Check(runner tflint.Runner) error {
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

		runner.EmitIssue(
			r,
			"Consider using aws_vpc_security_group_egress_rule or aws_vpc_security_group_ingress_rule instead.",
			attribute.Expr.Range(),
		)
	}

	return nil
}
