package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsSecurityGroupNoRuleBlock checks whether a route definition has a routing target
type AwsSecurityGroupNoRuleBlock struct {
	resourceType string
	blockTypes   []string
}

// NewAwsSecurityGroupNoRuleBlock returns new rule with default attributes
func NewAwsSecurityGroupNoRuleBlockRule() *AwsSecurityGroupNoRuleBlock {
	return &AwsSecurityGroupNoRuleBlock{
		resourceType: "aws_security_group",
		blockTypes:   []string{"ingress", "egress"},
	}
}

// Name returns the rule name
func (r *AwsSecurityGroupNoRuleBlock) Name() string {
	return "aws_security_group_no_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupNoRuleBlock) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSecurityGroupNoRuleBlock) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsSecurityGroupNoRuleBlock) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether a target is defined in a resource
func (r *AwsSecurityGroupNoRuleBlock) Check(runner tflint.Runner) error {
	for _, blockType := range r.blockTypes {
		err := runner.WalkResourceBlocks(r.resourceType, blockType, func(block *hcl.Block) error {
			msg := fmt.Sprintf("Inline rule block `%s` found! It is advisable to create a separate resource (`aws_security_group_rule`) instead", blockType)
			return runner.EmitIssue(
				r,
				msg,
				block.TypeRange,
			)
		})

		if err != nil {
			return err
		}
	}

	return nil
}
