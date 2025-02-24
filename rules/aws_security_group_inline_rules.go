package rules

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsSecurityGroupInlineRulesRule checks that egress and ingress blocks are not used in aws_security_group
type AwsSecurityGroupInlineRulesRule struct {
	tflint.DefaultRule

	resourceType string
}

// NewAwsSecurityGroupInlineRulesRule returns new rule with default attributes
func NewAwsSecurityGroupInlineRulesRule() *AwsSecurityGroupInlineRulesRule {
	return &AwsSecurityGroupInlineRulesRule{
		resourceType: "aws_security_group",
	}
}

// Name returns the rule name
func (r *AwsSecurityGroupInlineRulesRule) Name() string {
	return "aws_security_group_inline_rules"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupInlineRulesRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsSecurityGroupInlineRulesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsSecurityGroupInlineRulesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check that no egress and ingress blocks are used in a aws_security_group
func (r *AwsSecurityGroupInlineRulesRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "egress"},
			{Type: "ingress"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			runner.EmitIssue(
				r,
				fmt.Sprintf("Replace this %s block with aws_vpc_security_group_%s_rule.", block.Type, block.Type),
				block.DefRange,
			)
		}
	}

	return nil
}
