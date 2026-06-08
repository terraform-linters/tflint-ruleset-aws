package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsCloudwatchLogGroupRetentionRule checks that CloudWatch log groups have retention >= 365 days
type AwsCloudwatchLogGroupRetentionRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsCloudwatchLogGroupRetentionRule returns new rule with default attributes
func NewAwsCloudwatchLogGroupRetentionRule() *AwsCloudwatchLogGroupRetentionRule {
	return &AwsCloudwatchLogGroupRetentionRule{
		resourceType:  "aws_cloudwatch_log_group",
		attributeName: "retention_in_days",
	}
}

// Name returns the rule name
func (r *AwsCloudwatchLogGroupRetentionRule) Name() string {
	return "aws_cloudwatch_log_group_retention"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsCloudwatchLogGroupRetentionRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsCloudwatchLogGroupRetentionRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsCloudwatchLogGroupRetentionRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks that retention_in_days is at least 365
func (r *AwsCloudwatchLogGroupRetentionRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.attributeName}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			runner.EmitIssue(
				r,
				"The CloudWatch log group does not have retention_in_days set. It will retain logs forever. Set to at least 365 days.",
				resource.DefRange,
			)
			continue
		}

		err := runner.EvaluateExpr(attribute.Expr, func(days int) error {
			if days > 0 && days < 365 {
				runner.EmitIssue(
					r,
					fmt.Sprintf("The CloudWatch log group retention is %d days. Set to at least 365 days.", days),
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
