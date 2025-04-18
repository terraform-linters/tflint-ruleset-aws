// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule checks the pattern is valid
type AwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	enum          []string
}

// NewAwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule returns new rule with default attributes
func NewAwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule() *AwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule {
	return &AwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule{
		resourceType:  "aws_kinesisanalyticsv2_application",
		attributeName: "runtime_environment",
		enum: []string{
			"SQL-1_0",
			"FLINK-1_6",
			"FLINK-1_8",
			"ZEPPELIN-FLINK-1_0",
			"FLINK-1_11",
			"FLINK-1_13",
			"ZEPPELIN-FLINK-2_0",
			"FLINK-1_15",
			"ZEPPELIN-FLINK-3_0",
			"FLINK-1_18",
			"FLINK-1_19",
			"FLINK-1_20",
		},
	}
}

// Name returns the rule name
func (r *AwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule) Name() string {
	return "aws_kinesisanalyticsv2_application_invalid_runtime_environment"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsKinesisanalyticsv2ApplicationInvalidRuntimeEnvironmentRule) Check(runner tflint.Runner) error {
	logger.Trace("Check `%s` rule", r.Name())

	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		err := runner.EvaluateExpr(attribute.Expr, func (val string) error {
			found := false
			for _, item := range r.enum {
				if item == val {
					found = true
				}
			}
			if !found {
				runner.EmitIssue(
					r,
					fmt.Sprintf(`"%s" is an invalid value as runtime_environment`, truncateLongMessage(val)),
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
