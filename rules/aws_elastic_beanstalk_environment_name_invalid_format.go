package rules

import (
	"fmt"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsElasticBeanstalkEnvironmentNameInvalidFormatRule checks EB environment name matches a pattern
type AwsElasticBeanstalkEnvironmentNameInvalidFormatRule struct {
	resourceType  string
	attributeName string
	pattern       *regexp.Regexp
}

// NewAwsElasticBeanstalkEnvironmentNameInvalidFormatRule returns new rule with default attributes
func NewAwsElasticBeanstalkEnvironmentNameInvalidFormatRule() *AwsElasticBeanstalkEnvironmentNameInvalidFormatRule {
	return &AwsElasticBeanstalkEnvironmentNameInvalidFormatRule{
		resourceType:  "aws_elastic_beanstalk_environment",
		attributeName: "name",
		pattern:       regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$"),
	}
}

// Name returns the rule name
func (r *AwsElasticBeanstalkEnvironmentNameInvalidFormatRule) Name() string {
	return "aws_elastic_beanstalk_environment_name_invalid_format"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsElasticBeanstalkEnvironmentNameInvalidFormatRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsElasticBeanstalkEnvironmentNameInvalidFormatRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsElasticBeanstalkEnvironmentNameInvalidFormatRule) Link() string {
	return ""
}

// Check checks the environment name matches the pattern provided
func (r *AwsElasticBeanstalkEnvironmentNameInvalidFormatRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		return runner.EnsureNoError(err, func() error {
			if !r.pattern.MatchString(val) {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf(`%s does not match constraint: must contain only letters, digits, and the dash ` +
					`character and may not start or end with a dash (^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$)`, val),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
