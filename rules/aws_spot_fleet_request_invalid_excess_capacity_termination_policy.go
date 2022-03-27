package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule checks the pattern is valid
type AwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	enum          []string
}

// NewAwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule returns new rule with default attributes
func NewAwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule() *AwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule {
	return &AwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule{
		resourceType:  "aws_spot_fleet_request",
		attributeName: "excess_capacity_termination_policy",
		enum: []string{
			"Default",
			"NoTermination",
		},
	}
}

// Name returns the rule name
func (r *AwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule) Name() string {
	return "aws_spot_fleet_request_invalid_excess_capacity_termination_policy"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule) Check(runner tflint.Runner) error {
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
					fmt.Sprintf(`"%s" is an invalid value as excess_capacity_termination_policy`, val),
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
