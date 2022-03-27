package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/zclconf/go-cty/cty"
)

// AwsRouteSpecifiedMultipleTargetsRule checks whether a route definition has multiple routing targets
type AwsRouteSpecifiedMultipleTargetsRule struct {
	tflint.DefaultRule

	resourceType string
}

// NewAwsRouteSpecifiedMultipleTargetsRule returns new rule with default attributes
func NewAwsRouteSpecifiedMultipleTargetsRule() *AwsRouteSpecifiedMultipleTargetsRule {
	return &AwsRouteSpecifiedMultipleTargetsRule{
		resourceType: "aws_route",
	}
}

// Name returns the rule name
func (r *AwsRouteSpecifiedMultipleTargetsRule) Name() string {
	return "aws_route_specified_multiple_targets"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsRouteSpecifiedMultipleTargetsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsRouteSpecifiedMultipleTargetsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsRouteSpecifiedMultipleTargetsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether a resource defines multiple targets
func (r *AwsRouteSpecifiedMultipleTargetsRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "gateway_id"},
			{Name: "egress_only_gateway_id"},
			{Name: "nat_gateway_id"},
			{Name: "instance_id"},
			{Name: "vpc_peering_connection_id"},
			{Name: "network_interface_id"},
			{Name: "transit_gateway_id"},
			{Name: "vpc_endpoint_id"},
			{Name: "carrier_gateway_id"},
			{Name: "local_gateway_id"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		var nullAttributes int
		for _, attribute := range resource.Body.Attributes {
			var val cty.Value
			err := runner.EvaluateExpr(attribute.Expr, &val, nil)
			if err != nil {
				return err
			}

			if val.IsNull() {
				nullAttributes = nullAttributes + 1
			}
		}

		if len(resource.Body.Attributes)-nullAttributes > 1 {
			runner.EmitIssue(
				r,
				"More than one routing target specified. It must be one.",
				resource.DefRange,
			)
		}
	}

	return nil
}
