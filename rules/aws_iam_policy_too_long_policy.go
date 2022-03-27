package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsIAMPolicyTooLongPolicyRule checks that the policy length is less than 6,144 characters
type AwsIAMPolicyTooLongPolicyRule struct {
	tflint.DefaultRule

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
func (r *AwsIAMPolicyTooLongPolicyRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsIAMPolicyTooLongPolicyRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks the length of the policy
func (r *AwsIAMPolicyTooLongPolicyRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.attributeName}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		var policy string
		err := runner.EvaluateExpr(attribute.Expr, &policy, nil)
		whitespaceRegex := regexp.MustCompile(`\s+`)
		policy = whitespaceRegex.ReplaceAllString(policy, "")

		err = runner.EnsureNoError(err, func() error {
			if len(policy) > 6144 {
				runner.EmitIssue(
					r,
					fmt.Sprintf("The policy length is %d characters and is limited to 6144 characters.", len(policy)),
					attribute.Expr.Range(),
				)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
