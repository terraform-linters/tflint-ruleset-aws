package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsAcmCertificateLifecycleRule checks whether `create_before_destroy = true` is set in a lifecycle block of
// `aws_acm_certificate` resource
type AwsAcmCertificateLifecycleRule struct {
	resourceType  string
	attributeName string
}

// NewAwsAcmCertificateLifecycleRule returns new rule with default attributes
func NewAwsAcmCertificateLifecycleRule() *AwsAcmCertificateLifecycleRule {
	return &AwsAcmCertificateLifecycleRule{
		resourceType:  "aws_acm_certificate",
		attributeName: "lifecycle",
	}
}

// Name returns the rule name
func (r *AwsAcmCertificateLifecycleRule) Name() string {
	return "aws_acm_certificate_lifecycle"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsAcmCertificateLifecycleRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsAcmCertificateLifecycleRule) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsAcmCertificateLifecycleRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the aws_acm_certificate resource contains create_before_destroy = true in lifecycle block
func (r *AwsAcmCertificateLifecycleRule) Check(runner tflint.Runner) error {
	return runner.WalkResources("aws_acm_certificate", func(resource *configs.Resource) error {
		if !resource.Managed.CreateBeforeDestroy {
			if err := runner.EmitIssue(r, "resource `aws_acm_certificate` needs to contain `create_before_destroy = true` in `lifecycle` block", resource.DeclRange); err != nil {
				return err
			}
		}
		return nil
	})
}
