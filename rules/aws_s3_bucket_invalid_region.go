package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketInvalidRegionRule checks the pattern is valid
type AwsS3BucketInvalidRegionRule struct {
	resourceType  string
	attributeName string
	enum          []string
}

// NewAwsS3BucketInvalidRegionRule returns new rule with default attributes
func NewAwsS3BucketInvalidRegionRule() *AwsS3BucketInvalidRegionRule {
	return &AwsS3BucketInvalidRegionRule{
		resourceType:  "aws_s3_bucket",
		attributeName: "region",
		enum: []string{
			"EU",
			"us-east-1",
			"us-east-2",
			"eu-west-1",
			"eu-west-2",
			"eu-west-3",
			"eu-north-1",
			"eu-south-1",
			"us-west-1",
			"us-west-2",
			"ap-east-1",
			"ap-south-1",
			"ap-southeast-1",
			"ap-southeast-2",
			"ap-northeast-1",
			"ap-northeast-2",
			"ap-northeast-3",
			"ca-central-1",
			"sa-east-1",
			"cn-north-1",
			"cn-northwest-1",
			"eu-central-1",
			"me-south-1",
			"af-south-1",
			"us-gov-east-1",
			"us-gov-west-1",
		},
	}
}

// Name returns the rule name
func (r *AwsS3BucketInvalidRegionRule) Name() string {
	return "aws_s3_bucket_invalid_region"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsS3BucketInvalidRegionRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsS3BucketInvalidRegionRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsS3BucketInvalidRegionRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsS3BucketInvalidRegionRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		return runner.EnsureNoError(err, func() error {
			found := false
			for _, item := range r.enum {
				if item == val {
					found = true
				}
			}
			if !found {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf(`"%s" is an invalid value as region`, val),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
