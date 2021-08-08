package rules

import (
	"fmt"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsIAMPolicyInvalidLengthRule checks that the policy length is less than 2048 characters
type AwsIAMPolicyInvalidLengthRule struct {
	resourceType  string
	attributeName string
}

// NewAwsIAMPolicyInvalidLengthRule returns new rule with default attributes
func NewAwsIAMPolicyInvalidLengthRule() *AwsIAMPolicyInvalidLengthRule {
	return &AwsIAMPolicyInvalidLengthRule{
		resourceType:  "aws_iam_policy",
		attributeName: "policy",
	}
}

// Name returns the rule name
func (r *AwsIAMPolicyInvalidLengthRule) Name() string {
	return "aws_iam_policy_invalid_length"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMPolicyInvalidLengthRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsIAMPolicyInvalidLengthRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsIAMPolicyInvalidLengthRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks the length of the policy
func (r *AwsIAMPolicyInvalidLengthRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var policy string
		err := runner.EvaluateExpr(attribute.Expr, &policy, nil)

		return runner.EnsureNoError(err, func() error {
			if len(policy) > 2048  {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("The policy length is %d characters and is limited to 2048 characters.", len(policy)),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
