package rules

import (
	"fmt"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"regexp"
)

// AwsIAMPolicyTooLongPolicyRule checks that the policy length is less than 6,144 characters
type AwsIAMPolicyTooLongPolicyRule struct {
	resourceType  string
	attributeName string
}

// NewAwsIAMPolicyTooLongPolicyRule returns new rule with default attributes
func NewAwsIAMPolicyTooLongPolicyRule() *AwsIAMPolicyTooLongPolicyRule {
	return &AwsIAMPolicyTooLongPolicyRule{
		resourceType:  "aws_iam_policy",
		attributeName: "policy",
	}
}

// Name returns the rule name
func (r *AwsIAMPolicyTooLongPolicyRule) Name() string {
	return "aws_iam_policy_too_long_policy"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMPolicyTooLongPolicyRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsIAMPolicyTooLongPolicyRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsIAMPolicyTooLongPolicyRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks the length of the policy
func (r *AwsIAMPolicyTooLongPolicyRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var policy string
		err := runner.EvaluateExpr(attribute.Expr, &policy, nil)
		whitespaceRegex := regexp.MustCompile(`\s+`)
		policy = whitespaceRegex.ReplaceAllString(policy, "")
		return runner.EnsureNoError(err, func() error {
			if len(policy) > 6144 {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("The policy length is %d characters and is limited to 6144 characters.", len(policy)),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
