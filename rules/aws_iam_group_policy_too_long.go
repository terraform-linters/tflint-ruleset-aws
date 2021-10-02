package rules

import (
	"fmt"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"regexp"
)

// AwsIAMGroupPolicyTooLongRule checks that the policy length is less than 5,120 characters
type AwsIAMGroupPolicyTooLongRule struct {
	resourceType  string
	attributeName string
	whitespaceRegex *regexp.Regexp
}

// NewAwsIAMGroupPolicyTooLongRule returns new rule with default attributes
func NewAwsIAMGroupPolicyTooLongRule() *AwsIAMGroupPolicyTooLongRule {
	return &AwsIAMGroupPolicyTooLongRule{
		resourceType:  "aws_iam_group_policy",
		attributeName: "policy",
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

// Name returns the rule name
func (r *AwsIAMGroupPolicyTooLongRule) Name() string {
	return "aws_iam_group_policy_too_long"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMGroupPolicyTooLongRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsIAMGroupPolicyTooLongRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsIAMGroupPolicyTooLongRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks the length of the policy
func (r *AwsIAMGroupPolicyTooLongRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var policy string
		err := runner.EvaluateExpr(attribute.Expr, &policy, nil)
		return runner.EnsureNoError(err, func() error {
			policy = r.whitespaceRegex.ReplaceAllString(policy, "")
			if len(policy) > 5120 {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("The policy length is %d characters and is limited to 5120 characters.", len(policy)),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
