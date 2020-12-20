package rules

import (
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsRouteSpecifiedMultipleTargetsRule checks whether a route definition has multiple routing targets
type AwsRouteSpecifiedMultipleTargetsRule struct {
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
func (r *AwsRouteSpecifiedMultipleTargetsRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsRouteSpecifiedMultipleTargetsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether a resource defines `gateway_id`, `egress_only_gateway_id`, `nat_gateway_id`
// `instance_id`, `vpc_peering_connection_id` or `network_interface_id` at the same time
func (r *AwsRouteSpecifiedMultipleTargetsRule) Check(runner tflint.Runner) error {
	return runner.WalkResources(r.resourceType, func(resource *configs.Resource) error {
		body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: "gateway_id",
				},
				{
					Name: "egress_only_gateway_id",
				},
				{
					Name: "nat_gateway_id",
				},
				{
					Name: "instance_id",
				},
				{
					Name: "vpc_peering_connection_id",
				},
				{
					Name: "network_interface_id",
				},
				{
					Name: "transit_gateway_id",
				},
			},
		})
		if diags.HasErrors() {
			return diags
		}

		var nullAttributes int
		for _, attribute := range body.Attributes {
			isNull, err := runner.IsNullExpr(attribute.Expr)
			if err != nil {
				return err
			}

			if isNull {
				nullAttributes = nullAttributes + 1
			}
		}

		if len(body.Attributes)-nullAttributes > 1 {
			runner.EmitIssue(
				r,
				"More than one routing target specified. It must be one.",
				resource.DeclRange,
			)
		}

		return nil
	})
}
