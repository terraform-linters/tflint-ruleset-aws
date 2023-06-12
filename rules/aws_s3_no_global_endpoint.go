package rules

import (
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsS3NoGlobalEndpointRule checks whether deprecated s3 global endpoint is used instead
// of the regional endoint
type AwsS3NoGlobalEndpointRule struct {
	tflint.DefaultRule

	// TODO: do we need this, if so why
	resourceType  string
	attributeName string
}

// NewAwsS3NoGlobalEndpointRule returns new rule with default attributes
func NewAwsS3NoGlobalEndpointRule() *AwsS3NoGlobalEndpointRule {
	return &AwsS3NoGlobalEndpointRule{
		resourceType:  "aws_s3_bucket",
		attributeName: "",
	}
}

// Name returns the rule name
func (r *AwsS3NoGlobalEndpointRule) Name() string {
	return "aws_acm_certificate_lifecycle"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsS3NoGlobalEndpointRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsS3NoGlobalEndpointRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsS3NoGlobalEndpointRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the aws_acm_certificate resource contains create_before_destroy = true in lifecycle block
func (r *AwsS3NoGlobalEndpointRule) Check(runner tflint.Runner) error {
	var err error
	runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		vars := expr.Variables()

		if len(vars) == 0 {
			return nil
		}

		// is this ever greater than 0
		v := vars[0]

		if v.RootName() == "aws_s3_bucket" && len(v) == 3 && v[2].(hcl.TraverseAttr).Name == "bucket_domain_name" {
			err = runner.EmitIssue(r, "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead", v.SourceRange())
		}

		return nil
	}))

	return err
}
