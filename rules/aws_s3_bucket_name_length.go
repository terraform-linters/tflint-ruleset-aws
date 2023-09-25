package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsS3BucketNameLengthRule checks that an S3 bucket name contains a valid number of characteres
type AwsS3BucketNameLengthRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsS3BucketNameLengthRule returns new rule with default attributes
func NewAwsS3BucketNameLengthRule() *AwsS3BucketNameLengthRule {
	return &AwsS3BucketNameLengthRule{
		resourceType:  "aws_s3_bucket",
		attributeName: "bucket",
	}
}

// Name returns the rule name
func (r *AwsS3BucketNameLengthRule) Name() string {
	return "aws_s3_bucket_name_length"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsS3BucketNameLengthRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsS3BucketNameLengthRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsS3BucketNameLengthRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check if s3 bucket name is within the characteres length accepted range
func (r *AwsS3BucketNameLengthRule) Check(runner tflint.Runner) error {

	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

  bucketNameMinLength := 3
  bucketNameMaxLength := 63

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		err := runner.EvaluateExpr(attribute.Expr, func (name string) error {
      if len(name) < bucketNameMinLength || len(name) > bucketNameMaxLength {
				runner.EmitIssue(
					r,
          fmt.Sprintf("Bucket name %q must be between %d and %d characters", name, bucketNameMinLength, bucketNameMaxLength),
					attribute.Expr.Range(),
				)
			}
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
