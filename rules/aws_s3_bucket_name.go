package rules

import (
	"fmt"
	"regexp"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsS3BucketNameRule checks that an S3 bucket name matches naming rules
type AwsS3BucketNameRule struct {
	resourceType  string
	attributeName string
}

type awsS3BucketNameConfig struct {
	Regex  string `hcl:"regex,optional"`
	Prefix string `hcl:"prefix,optional"`

	Remain hcl.Body `hcl:",remain"`
}

// NewAwsS3BucketNameRule returns a new rule with default attributes
func NewAwsS3BucketNameRule() *AwsS3BucketNameRule {
	return &AwsS3BucketNameRule{
		resourceType:  "aws_s3_bucket",
		attributeName: "bucket",
	}
}

// Name returns the rule name
func (r *AwsS3BucketNameRule) Name() string {
	return "aws_s3_bucket_name"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsS3BucketNameRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsS3BucketNameRule) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsS3BucketNameRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check if the name of the s3 bucket matches the regex defined in the rule
func (r *AwsS3BucketNameRule) Check(runner tflint.Runner) error {
	config := awsS3BucketNameConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	regex, err := regexp.Compile(config.Regex)
	if err != nil {
		return fmt.Errorf("invalid regex (%s): %w", config.Regex, err)
	}

	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var name string
		err := runner.EvaluateExpr(attribute.Expr, &name, nil)

		return runner.EnsureNoError(err, func() error {
			if config.Prefix != "" {
				if !strings.HasPrefix(name, config.Prefix) {
					runner.EmitIssue(
						r,
						fmt.Sprintf(`Bucket name "%s" does not have prefix "%s"`, name, config.Prefix),
						attribute.Expr.Range(),
					)
				}
			}

			if config.Regex != "" {
				if !regex.MatchString(name) {
					runner.EmitIssue(
						r,
						fmt.Sprintf(`Bucket name "%s" does not match regex "%s"`, name, config.Regex),
						attribute.Expr.Range(),
					)
				}
			}

			return nil
		})
	})
}
