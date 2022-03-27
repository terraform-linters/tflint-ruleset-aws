package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

const statementBlockName = "statement"

// AwsIAMPolicyDocumentGovFriendlyArnsRule checks for non-GovCloud arns
type AwsIAMPolicyDocumentGovFriendlyArnsRule struct {
	tflint.DefaultRule

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
func (r *AwsIAMPolicyDocumentGovFriendlyArnsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsIAMPolicyDocumentGovFriendlyArnsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks the pattern is valid
func (r *AwsIAMPolicyDocumentGovFriendlyArnsRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: statementBlockName,
				Body: &hclext.BodySchema{Attributes: []hclext.AttributeSchema{{Name: r.attributeName}}},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, statement := range resource.Body.Blocks {
			attribute, exists := statement.Body.Attributes[r.attributeName]
			if !exists {
				continue
			}

			var val []string
			err := runner.EvaluateExpr(attribute.Expr, &val, nil)

			err = runner.EnsureNoError(err, func() error {
				for _, arn := range val {
					if r.pattern.MatchString(arn) {
						runner.EmitIssue(
							r,
							fmt.Sprintf(`ARN detected in IAM policy document that could potentially fail in AWS GovCloud due to resource pattern: %s`, r.pattern),
							attribute.Expr.Range(),
						)
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
