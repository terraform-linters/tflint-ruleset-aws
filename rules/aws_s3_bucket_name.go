package rules

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsS3BucketNameRule checks that an S3 bucket name matches naming rules
type AwsS3BucketNameRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

type awsS3BucketNameConfig struct {
	Regex  string `hclext:"regex,optional"`
	Prefix string `hclext:"prefix,optional"`
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
	return true
}

// Severity returns the rule severity
func (r *AwsS3BucketNameRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsS3BucketNameRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check if the name of the s3 bucket is valid
func (r *AwsS3BucketNameRule) Check(runner tflint.Runner) error {
	config := awsS3BucketNameConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	regex, err := regexp.Compile(config.Regex)
	if err != nil {
		return fmt.Errorf("invalid regex (%s): %w", config.Regex, err)
	}

	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.attributeName}},
	}, nil)
	if err != nil {
		return err
	}

	bucketNameMinLength := 3
	bucketNameMaxLength := 63

	// https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
	regexRules := []struct {
		Regexp      regexp.Regexp
		Description string
	}{
		{
			Regexp: *regexp.MustCompile(fmt.Sprintf(
				"^.{0,%d}$|.{%d,}$",
				bucketNameMinLength-1,
				bucketNameMaxLength+1,
			)),
			Description: fmt.Sprintf(
				"Bucket names must be between %d (min) and %d (max) characters long.",
				bucketNameMinLength,
				bucketNameMaxLength,
			),
		},
		{
			Regexp:      *regexp.MustCompile("[^a-z0-9\\.\\-]"),
			Description: "Bucket names can consist only of lowercase letters, numbers, dots (.), and hyphens (-).",
		},
		{
			Regexp:      *regexp.MustCompile("^[^a-z0-9]|[^a-z0-9]$"),
			Description: "Bucket names must begin and end with a lowercase letter or number.",
		},
		{
			Regexp:      *regexp.MustCompile("\\.\\."),
			Description: "Bucket names must not contain two adjacent periods.",
		},
		{
			Regexp:      *regexp.MustCompile("^xn--"),
			Description: "Bucket names must not start with the prefix 'xn--'.",
		},
		{
			Regexp:      *regexp.MustCompile("^(sthree-|sthree-configurator|amzn-s3-demo-)"),
			Description: "Bucket names must not start with the prefix 'sthree-', 'sthree-configurator', or 'amzn-s3-demo-'.",
		},
		{
			Regexp:      *regexp.MustCompile("-s3alias$"),
			Description: "Bucket names must not end with the suffix '-s3alias'.",
		},
		{
			Regexp:      *regexp.MustCompile("--ol-s3$"),
			Description: "Bucket names must not end with the suffix '--ol-s3'.",
		},
		{
			Regexp:      *regexp.MustCompile("\\.mrap$"),
			Description: "Bucket names must not end with the suffix '.mrap'.",
		},
		{
			Regexp:      *regexp.MustCompile("--x-s3$"),
			Description: "Bucket names must not end with the suffix '--x-s3'.",
		},
		{
			Regexp:      *regexp.MustCompile("--table-s3$"),
			Description: "Bucket names must not end with the suffix '--table-s3'.",
		},
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		err := runner.EvaluateExpr(attribute.Expr, func(name string) error {
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

      nameAsIP := net.ParseIP(name)
      if nameAsIP != nil && net.IP.To4(nameAsIP) != nil {
        runner.EmitIssue(
          r,
          fmt.Sprintf(`Bucket names must not be formatted as an IP address. (name: %q)`, name),
          attribute.Expr.Range(),
        )
      }

			for _, rule := range regexRules {
				if rule.Regexp.MatchString(name) {
					runner.EmitIssue(
						r,
						fmt.Sprintf(`%s (name: %q, regex: %q)`, rule.Description, name, rule.Regexp.String()),
						attribute.Expr.Range(),
					)
				}
			}
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
