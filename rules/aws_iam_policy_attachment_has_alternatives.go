package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsIAMPolicyAttachmentHasAlternativesRule warns that the resource has alternatives recommended
type AwsIAMPolicyAttachmentHasAlternativesRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// AwsIAMPolicyAttachmentHasAlternativesRule returns new rule with default attributes
func NewAwsIAMPolicyAttachmentHasAlternativesRule() *AwsIAMPolicyAttachmentHasAlternativesRule {
	return &AwsIAMPolicyAttachmentHasAlternativesRule{
		resourceType:  "aws_iam_policy_attachment",
		attributeName: "name",
	}
}

// Name returns the rule name
func (r *AwsIAMPolicyAttachmentHasAlternativesRule) Name() string {
	return "aws_iam_policy_attachment_has_alternatives"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMPolicyAttachmentHasAlternativesRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsIAMPolicyAttachmentHasAlternativesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsIAMPolicyAttachmentHasAlternativesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks the length of the policy
func (r *AwsIAMPolicyAttachmentHasAlternativesRule) Check(runner tflint.Runner) error {
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

		var name string
		err := runner.EvaluateExpr(attribute.Expr, &name, nil)
		runner.EmitIssue(
			r,
			"Consider aws_iam_role_policy_attachment, aws_iam_user_policy_attachment, or aws_iam_group_policy_attachment instead.",
			attribute.Expr.Range(),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
