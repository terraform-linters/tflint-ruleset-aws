package rules

import (
	"fmt"
	"log"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsIAMPolicyGovFriendlyArnsRule checks for non-GovCloud arns
type AwsIAMPolicyGovFriendlyArnsRule struct {
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
func (r *AwsIAMPolicyGovFriendlyArnsRule) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsIAMPolicyGovFriendlyArnsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks the pattern is valid
func (r *AwsIAMPolicyGovFriendlyArnsRule) Check(runner tflint.Runner) error {
	log.Printf("[TRACE] Check `%s` rule", r.Name())
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		return runner.EnsureNoError(err, func() error {
			if r.pattern.MatchString(val) {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf(`ARN detected in IAM policy that could potentially fail in AWS GovCloud due to resource pattern: %s`, r.pattern),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
