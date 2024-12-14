package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// TODO: Write the rule's description here
// AwsSecurityGroupRuleDeprecatedRule checks ...
type AwsSecurityGroupRuleDeprecatedRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsSecurityGroupRuleDeprecatedRule returns new rule with default attributes
func NewAwsSecurityGroupRuleDeprecatedRule() *AwsSecurityGroupRuleDeprecatedRule {
	return &AwsSecurityGroupRuleDeprecatedRule{
		// TODO: Write resource type and attribute name here
		resourceType:  "...",
		attributeName: "...",
	}
}

// Name returns the rule name
func (r *AwsSecurityGroupRuleDeprecatedRule) Name() string {
	return "aws_security_group_rule_deprecated"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupRuleDeprecatedRule) Enabled() bool {
	// TODO: Determine whether the rule is enabled by default
	return true
}

// Severity returns the rule severity
func (r *AwsSecurityGroupRuleDeprecatedRule) Severity() tflint.Severity {
	// TODO: Determine the rule's severiry
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSecurityGroupRuleDeprecatedRule) Link() string {
	// TODO: If the rule is so trivial that no documentation is needed, return "" instead.
	return project.ReferenceLink(r.Name())
}

// TODO: Write the details of the inspection
// Check checks ...
func (r *AwsSecurityGroupRuleDeprecatedRule) Check(runner tflint.Runner) error {
	// TODO: Write the implementation here. See this documentation for what tflint.Runner can do.
	//       https://pkg.go.dev/github.com/terraform-linters/tflint-plugin-sdk/tflint#Runner

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
