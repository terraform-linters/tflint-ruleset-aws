package rules

import (
	"fmt"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsSecurityGroupEgressAndIngressBlocksDeprecatedRule checks that egress and ingress blocks are not used in aws_security_group
type AwsSecurityGroupEgressAndIngressBlocksDeprecatedRule struct {
	tflint.DefaultRule

	resourceType string
}

// NewAwsSecurityGroupEgressAndIngressBlocksDeprecatedRule returns new rule with default attributes
func NewAwsSecurityGroupEgressAndIngressBlocksDeprecatedRule() *AwsSecurityGroupEgressAndIngressBlocksDeprecatedRule {
	return &AwsSecurityGroupEgressAndIngressBlocksDeprecatedRule{
		resourceType: "aws_security_group",
	}
}

// Name returns the rule name
func (r *AwsSecurityGroupEgressAndIngressBlocksDeprecatedRule) Name() string {
	return "aws_security_group_egress_and_ingress_blocks_deprecated"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupEgressAndIngressBlocksDeprecatedRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsSecurityGroupEgressAndIngressBlocksDeprecatedRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsSecurityGroupEgressAndIngressBlocksDeprecatedRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check that no egress and ingress blocks are used in a aws_security_group
func (r *AwsSecurityGroupEgressAndIngressBlocksDeprecatedRule) Check(runner tflint.Runner) error {
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
			if block.Type == "egress" || block.Type == "ingress" {
				runner.EmitIssue(
					r,
					fmt.Sprintf("Replace this %s block with aws_vpc_security_group_%s_rule. Otherwise, rule conflicts, perpetual differences, and result in rules being overwritten may occur.", block.Type, block.Type),
					block.DefRange,
				)
			} else {
				return fmt.Errorf("unexpected resource type: %s", block.Type)
			}
		}
	}

	return nil
}
