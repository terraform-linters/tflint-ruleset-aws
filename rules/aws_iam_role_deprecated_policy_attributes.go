package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsIAMRoleDeprecatedPolicyAttributesRule checks that neither `inline_policy` or `managed_policy_arns` are used in an `aws_iam_role`
type AwsIAMRoleDeprecatedPolicyAttributesRule struct {
	tflint.DefaultRule

	resourceType string
}

// NewAwsIAMRoleDeprecatedPolicyAttributesRule returns new rule with default attributes
func NewAwsIAMRoleDeprecatedPolicyAttributesRule() *AwsIAMRoleDeprecatedPolicyAttributesRule {
	return &AwsIAMRoleDeprecatedPolicyAttributesRule{
		resourceType: "aws_iam_role",
	}
}

// Name returns the rule name
func (r *AwsIAMRoleDeprecatedPolicyAttributesRule) Name() string {
	return "aws_iam_role_deprecated_policy_attributes"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMRoleDeprecatedPolicyAttributesRule) Enabled() bool {
	// `inline_policy` and `managed_policy_arns` were deprecated in 5.68.0 and 5.72.0 respectively
	// At time of writing, this is quite recent, so don't automatically enable the rule.
	return false
}

// Severity returns the rule severity
func (r *AwsIAMRoleDeprecatedPolicyAttributesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsIAMRoleDeprecatedPolicyAttributesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks that aws_iam_role resources don't have any `inline_policy` blocks or the `managed_policy_arns` attribute
func (r *AwsIAMRoleDeprecatedPolicyAttributesRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "inline_policy",
			},
		},
		Attributes: []hclext.AttributeSchema{
			{
				Name: "managed_policy_arns",
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes["managed_policy_arns"]

		if exists {
			if err := runner.EmitIssue(r, "The managed_policy_arns argument is deprecated. Use the aws_iam_role_policy_attachment resource instead. If Terraform should exclusively manage all managed policy attachments (the current behavior of this argument), use the aws_iam_role_policy_attachments_exclusive resource as well.", attribute.Range); err != nil {
				return err
			}
		}

		for _, rule := range resource.Body.Blocks {
			if rule.Type != "inline_policy" {
				continue
			}

			if err := runner.EmitIssue(r, "The inline_policy argument is deprecated. Use the aws_iam_role_policy resource instead. If Terraform should exclusively manage all inline policy associations (the current behavior of this argument), use the aws_iam_role_policies_exclusive resource as well.", rule.DefRange); err != nil {
				return err
			}
		}
	}

	return nil
}
