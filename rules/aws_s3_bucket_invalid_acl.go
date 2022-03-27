package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketInvalidACLRule checks the pattern is valid
type AwsS3BucketInvalidACLRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	enum          []string
}

// NewAwsS3BucketInvalidACLRule returns new rule with default attributes
func NewAwsS3BucketInvalidACLRule() *AwsS3BucketInvalidACLRule {
	return &AwsS3BucketInvalidACLRule{
		resourceType:  "aws_s3_bucket",
		attributeName: "acl",
		enum: []string{
			"private",
			"public-read",
			"public-read-write",
			"aws-exec-read",
			"authenticated-read",
			"log-delivery-write",
			"bucket-owner-read",
			"bucket-owner-full-control",
		},
	}
}

// Name returns the rule name
func (r *AwsS3BucketInvalidACLRule) Name() string {
	return "aws_s3_bucket_invalid_acl"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsS3BucketInvalidACLRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsS3BucketInvalidACLRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsS3BucketInvalidACLRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsS3BucketInvalidACLRule) Check(runner tflint.Runner) error {
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

		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		err = runner.EnsureNoError(err, func() error {
			found := false
			for _, item := range r.enum {
				if item == val {
					found = true
				}
			}
			if !found {
				runner.EmitIssue(
					r,
					fmt.Sprintf(`"%s" is an invalid value as acl`, val),
					attribute.Expr.Range(),
				)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
