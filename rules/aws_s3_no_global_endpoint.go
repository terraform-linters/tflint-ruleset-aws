package rules

import (
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	hclsyntax "github.com/hashicorp/hcl/v2/hclsyntax"
	lang "github.com/terraform-linters/tflint-plugin-sdk/terraform/lang"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsS3NoGlobalEndpointRule checks whether deprecated s3 global endpoint is used
type AwsS3NoGlobalEndpointRule struct {
	tflint.DefaultRule
}

// NewAwsS3NoGlobalEndpointRule returns new rule with default attributes
func NewAwsS3NoGlobalEndpointRule() *AwsS3NoGlobalEndpointRule {
	return &AwsS3NoGlobalEndpointRule{}
}

// Name returns the rule name
func (r *AwsS3NoGlobalEndpointRule) Name() string {
	return "aws_s3_no_global_endpoint"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsS3NoGlobalEndpointRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsS3NoGlobalEndpointRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsS3NoGlobalEndpointRule) Link() string {
	return project.ReferenceLink(r.Name())
}

func isS3BucketResourceOrData(str string) bool {
	return strings.HasPrefix(str, "aws_s3_bucket") || strings.HasPrefix(str, "data.aws_s3_bucket")
}

func getRealRange(srcRange hcl.Range) hcl.Range {
	// ref.SourceRange doesn't include the bucket_domain_name attribute which is 19 characters long
	srcRange.End.Column += 19
	return srcRange
}

// Check checks whether the deprecated s3 global endpoint is used
func (r *AwsS3NoGlobalEndpointRule) Check(runner tflint.Runner) error {
	var err error
	runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		_, isTemplateExpr := expr.(*hclsyntax.TemplateExpr)
		if isTemplateExpr {
			// TODO: check sth like .s3.amazonaws.com inside template expression ?

			// "${aws_s3_bucket.test.bucket_domain_name}" would report the same error twice if we didn't return here
			return nil
		}

		refs := lang.ReferencesInExpr(expr)

		if len(refs) != 1 {
			return nil
		}

		ref := refs[0]

		// lang.ReferencesInExpr(expr) is broken when the reference contains subexpression ?
		// e.g. aws_s3_bucket.test[length(var.bucket_names) - 1]
		// TODO: seems like a niche case, should we handle this ?
		if isS3BucketResourceOrData(ref.Subject.String()) && len(ref.Remaining) == 1 && ref.Remaining[0].(hcl.TraverseAttr).Name == "bucket_domain_name" {
			// ref.SourceRange doesn't include the bucket_domain_name attribute which is 19 characters long
			srcRange := getRealRange(ref.SourceRange)
			err = runner.EmitIssue(r, "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead", srcRange)
		}

		// TODO: handle foreach, count, dynamic blocks

		return nil
	}))

	return err
}
