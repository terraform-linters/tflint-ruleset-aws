package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsElasticBeanstalkEnvironmentInvalidNameFormatRule checks EB environment name matches a pattern
type AwsElasticBeanstalkEnvironmentInvalidNameFormatRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	pattern       *regexp.Regexp
}

// NewAwsElasticBeanstalkEnvironmentInvalidNameFormatRule returns new rule with default attributes
func NewAwsElasticBeanstalkEnvironmentInvalidNameFormatRule() *AwsElasticBeanstalkEnvironmentInvalidNameFormatRule {
	return &AwsElasticBeanstalkEnvironmentInvalidNameFormatRule{
		resourceType:  "aws_elastic_beanstalk_environment",
		attributeName: "name",
		pattern:       regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$"),
	}
}

// Name returns the rule name
func (r *AwsElasticBeanstalkEnvironmentInvalidNameFormatRule) Name() string {
	return "aws_elastic_beanstalk_environment_invalid_name_format"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsElasticBeanstalkEnvironmentInvalidNameFormatRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsElasticBeanstalkEnvironmentInvalidNameFormatRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsElasticBeanstalkEnvironmentInvalidNameFormatRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks the environment name matches the pattern provided
func (r *AwsElasticBeanstalkEnvironmentInvalidNameFormatRule) Check(runner tflint.Runner) error {
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

		err := runner.EvaluateExpr(attribute.Expr, func(val string) error {
			if !r.pattern.MatchString(val) {
				runner.EmitIssue(
					r,
					fmt.Sprintf(`%s does not match constraint: must contain only letters, digits, and the dash `+
						`character and may not start or end with a dash (^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$)`, val),
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
