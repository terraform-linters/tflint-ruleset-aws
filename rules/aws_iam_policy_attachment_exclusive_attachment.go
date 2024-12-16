package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsIAMPolicyAttachmentExclusiveAttachmentRule warns that the resource has alternatives recommended
type AwsIAMPolicyAttachmentExclusiveAttachmentRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// AwsIAMPolicyAttachmentExclusiveAttachmentRule returns new rule with default attributes
func NewAwsIAMPolicyAttachmentExclusiveAttachmentRule() *AwsIAMPolicyAttachmentExclusiveAttachmentRule {
	return &AwsIAMPolicyAttachmentExclusiveAttachmentRule{
		resourceType:  "aws_iam_policy_attachment",
		attributeName: "name",
	}
}

// Name returns the rule name
func (r *AwsIAMPolicyAttachmentExclusiveAttachmentRule) Name() string {
	return "aws_iam_policy_attachment_exclusive_attachment"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMPolicyAttachmentExclusiveAttachmentRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsIAMPolicyAttachmentExclusiveAttachmentRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsIAMPolicyAttachmentExclusiveAttachmentRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check that the resource is not used
func (r *AwsIAMPolicyAttachmentExclusiveAttachmentRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		runner.EmitIssue(
			r,
			"Within the entire AWS account, all users, roles, and groups that a single policy is attached to must be specified by a single aws_iam_policy_attachment resource. Consider aws_iam_role_policy_attachment, aws_iam_user_policy_attachment, or aws_iam_group_policy_attachment instead.",
			resource.DefRange,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
