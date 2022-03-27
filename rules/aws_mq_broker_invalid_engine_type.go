package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsMqBrokerInvalidEngineTypeRule checks the pattern is valid
type AwsMqBrokerInvalidEngineTypeRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	enum          []string
}

// NewAwsMqBrokerInvalidEngineTypeRule returns new rule with default attributes
func NewAwsMqBrokerInvalidEngineTypeRule() *AwsMqBrokerInvalidEngineTypeRule {
	return &AwsMqBrokerInvalidEngineTypeRule{
		resourceType:  "aws_mq_broker",
		attributeName: "engine_type",
		enum: []string{
			"ActiveMQ",
			"RabbitMQ",
		},
	}
}

// Name returns the rule name
func (r *AwsMqBrokerInvalidEngineTypeRule) Name() string {
	return "aws_mq_broker_invalid_engine_type"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsMqBrokerInvalidEngineTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsMqBrokerInvalidEngineTypeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsMqBrokerInvalidEngineTypeRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsMqBrokerInvalidEngineTypeRule) Check(runner tflint.Runner) error {
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
					`engine_type is not a valid value`,
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
