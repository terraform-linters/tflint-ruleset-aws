package rules

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsAcmCertificateLifecycleRule checks whether `create_before_destroy = true` is set in a lifecycle block of
// `aws_acm_certificate` resource
type AwsAcmCertificateLifecycleRule struct {
	tflint.DefaultRule

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
func (r *AwsAcmCertificateLifecycleRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsAcmCertificateLifecycleRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the aws_acm_certificate resource contains create_before_destroy = true in lifecycle block
func (r *AwsAcmCertificateLifecycleRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("aws_acm_certificate", &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "lifecycle",
				Body: &hclext.BodySchema{Attributes: []hclext.AttributeSchema{{Name: "create_before_destroy"}}},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		if len(resource.Body.Blocks) == 0 {
			if err := runner.EmitIssue(r, "resource `aws_acm_certificate` needs to contain `create_before_destroy = true` in `lifecycle` block", resource.DefRange); err != nil {
				return err
			}
			continue
		}

		for _, lifecycle := range resource.Body.Blocks {
			attr, exists := lifecycle.Body.Attributes["create_before_destroy"]
			if !exists {
				if err := runner.EmitIssue(r, "resource `aws_acm_certificate` needs to contain `create_before_destroy = true` in `lifecycle` block", resource.DefRange); err != nil {
					return err
				}
				continue
			}

			var createBeforeDestroy bool
			// MEMO: lifecycle attributes do not support interpolations
			//       https://github.com/hashicorp/terraform/issues/3116
			diags := gohcl.DecodeExpression(attr.Expr, nil, &createBeforeDestroy)
			if diags.HasErrors() {
				return diags
			}

			if !createBeforeDestroy {
				if err := runner.EmitIssue(r, "resource `aws_acm_certificate` needs to contain `create_before_destroy = true` in `lifecycle` block", resource.DefRange); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
