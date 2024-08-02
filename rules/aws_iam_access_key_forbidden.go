package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsIAMAccessKeyForbiddenRule checks that aws_iam_access_key resource is not used
type AwsIAMAccessKeyForbiddenRule struct {
	tflint.DefaultRule

	resourceType string
}

// NewAwsIAMAccessKeyForbiddenRule returns new rule with default attributes
func NewAwsIAMAccessKeyForbiddenRule() *AwsIAMAccessKeyForbiddenRule {
	return &AwsIAMAccessKeyForbiddenRule{
		resourceType: "aws_iam_access_key",
	}
}

// Name returns the rule name
func (r *AwsIAMAccessKeyForbiddenRule) Name() string {
	return "aws_iam_access_key_forbidden"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMAccessKeyForbiddenRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsIAMAccessKeyForbiddenRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsIAMAccessKeyForbiddenRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks that no aws_iam_access_key resources are being used
func (r *AwsIAMAccessKeyForbiddenRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		runner.EmitIssue(
			r,
			"\"aws_iam_access_key\" resources are forbidden.",
			resource.DefRange,
		)
	}

	return nil
}
