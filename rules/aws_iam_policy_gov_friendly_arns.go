package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsIAMPolicyGovFriendlyArnsRule checks for non-GovCloud arns
type AwsIAMPolicyGovFriendlyArnsRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	pattern       *regexp.Regexp
}

// NewAwsIAMPolicyInvalidPolicyRule returns new rule with default attributes
func NewAwsIAMPolicyGovFriendlyArnsRule() *AwsIAMPolicyGovFriendlyArnsRule {
	return &AwsIAMPolicyGovFriendlyArnsRule{
		resourceType:  "aws_iam_policy",
		attributeName: "policy",
		// AWS GovCloud arn separator is arn:aws-us-gov
		// https://docs.aws.amazon.com/govcloud-us/latest/UserGuide/using-govcloud-arns.html
		pattern: regexp.MustCompile(`arn:aws:.*`),
	}
}

// Name returns the rule name
func (r *AwsIAMPolicyGovFriendlyArnsRule) Name() string {
	return "aws_iam_policy_gov_friendly_arns"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMPolicyGovFriendlyArnsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsIAMPolicyGovFriendlyArnsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsIAMPolicyGovFriendlyArnsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks the pattern is valid
func (r *AwsIAMPolicyGovFriendlyArnsRule) Check(runner tflint.Runner) error {
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

		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		err = runner.EnsureNoError(err, func() error {
			if r.pattern.MatchString(val) {
				runner.EmitIssue(
					r,
					fmt.Sprintf(`ARN detected in IAM policy that could potentially fail in AWS GovCloud due to resource pattern: %s`, r.pattern),
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
