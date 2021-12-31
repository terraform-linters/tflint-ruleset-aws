package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/zclconf/go-cty/cty"
)

// AwsRouteNotSpecifiedTargetRule checks whether a route definition has a routing target
type AwsRouteNotSpecifiedTargetRule struct {
	tflint.DefaultRule

	resourceType string
}

// NewAwsRouteNotSpecifiedTargetRule returns new rule with default attributes
func NewAwsRouteNotSpecifiedTargetRule() *AwsRouteNotSpecifiedTargetRule {
	return &AwsRouteNotSpecifiedTargetRule{
		resourceType: "aws_route",
	}
}

// Name returns the rule name
func (r *AwsRouteNotSpecifiedTargetRule) Name() string {
	return "aws_route_not_specified_target"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsRouteNotSpecifiedTargetRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsRouteNotSpecifiedTargetRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsRouteNotSpecifiedTargetRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether a target is defined in a resource
func (r *AwsRouteNotSpecifiedTargetRule) Check(runner tflint.Runner) error {
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

		if len(resource.Body.Attributes)-nullAttributes == 0 {
			runner.EmitIssue(
				r,
				"The routing target is not specified, each aws_route must contain either egress_only_gateway_id, gateway_id, instance_id, nat_gateway_id, network_interface_id, transit_gateway_id, vpc_peering_connection_id or vpc_endpoint_id.",
				resource.DefRange,
			)
		}
	}

	return nil
}
