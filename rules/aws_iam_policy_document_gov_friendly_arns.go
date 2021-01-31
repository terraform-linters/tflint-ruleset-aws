package rules

import (
	"fmt"
	"log"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const statementBlockName = "statement"

// AwsIAMPolicyDocumentGovFriendlyArnsRule checks for non-GovCloud arns
type AwsIAMPolicyDocumentGovFriendlyArnsRule struct {
	resourceType  string
	attributeName string
	pattern       *regexp.Regexp
}

// AwsIAMPolicyDocumentGovFriendlyArnsRule returns new rule with default attributes
func NewAwsIAMPolicyDocumentGovFriendlyArnsRule() *AwsIAMPolicyDocumentGovFriendlyArnsRule {
	return &AwsIAMPolicyDocumentGovFriendlyArnsRule{
		resourceType:  "aws_iam_policy_document",
		attributeName: "resources",
		// AWS GovCloud arn separator is arn:aws-us-gov
		// https://docs.aws.amazon.com/govcloud-us/latest/UserGuide/using-govcloud-arns.html
		pattern: regexp.MustCompile(`arn:aws:.*`),
	}
}

// Name returns the rule name
func (r *AwsIAMPolicyDocumentGovFriendlyArnsRule) Name() string {
	return "aws_iam_policy_document_gov_friendly_arns"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMPolicyDocumentGovFriendlyArnsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsIAMPolicyDocumentGovFriendlyArnsRule) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsIAMPolicyDocumentGovFriendlyArnsRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsIAMPolicyDocumentGovFriendlyArnsRule) Check(runner tflint.Runner) error {
	log.Printf("[TRACE] Check `%s` rule", r.Name())
	return runner.WalkResourceBlocks(r.resourceType, statementBlockName, func(statementBlock *hcl.Block) error {
		attributes, diags := statementBlock.Body.JustAttributes()
		if diags.HasErrors() {
			return diags
		}
		var val []string
		err := runner.EvaluateExpr(attributes[r.attributeName].Expr, &val, nil)
		return runner.EnsureNoError(err, func() error {
			for _, arn := range val {
				if r.pattern.MatchString(arn) {
					runner.EmitIssueOnExpr(
						r,
						fmt.Sprintf(`ARN detected in IAM policy document that could potentially fail in AWS GovCloud due to resource pattern: %s`, r.pattern),
						attributes[r.attributeName].Expr,
					)
				}
			}
			return nil
		})
	})
}
