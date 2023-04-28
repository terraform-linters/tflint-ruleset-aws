package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsRouteNotSpecifiedTarget(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "route target is not specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsRouteNotSpecifiedTargetRule(),
					Message: "The routing target is not specified, each aws_route must contain either egress_only_gateway_id, gateway_id, instance_id, nat_gateway_id, network_interface_id, transit_gateway_id, vpc_peering_connection_id, core_network_arn or vpc_endpoint_id.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 27},
					},
				},
			},
		},
		{
			Name: "gateway_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    gateway_id = "igw-1234abcd"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "egress_only_gateway_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    egress_only_gateway_id = "eigw-1234abcd"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "nat_gateway_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    nat_gateway_id = "nat-1234abcd"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "instance_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    instance_id = "i-1234abcd"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "vpc_peering_connection_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    vpc_peering_connection_id = "pcx-1234abcd"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "network_interface_id is specified",
			Content: `
resource "aws_route" "foo" {
    route_table_id = "rtb-1234abcd"
    network_interface_id = "eni-1234abcd"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "transit_gateway_id is specified",
			Content: `
resource "aws_route" "foo" {
	route_table_id = "rtb-1234abcd"
	transit_gateway_id = "tgw-1234abcd"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "transit_gateway_id is specified, but the value is null",
			Content: `
resource "aws_route" "foo" {
	route_table_id = "rtb-1234abcd"
	transit_gateway_id = null
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsRouteNotSpecifiedTargetRule(),
					Message: "The routing target is not specified, each aws_route must contain either egress_only_gateway_id, gateway_id, instance_id, nat_gateway_id, network_interface_id, transit_gateway_id, vpc_peering_connection_id, core_network_arn or vpc_endpoint_id.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 27},
					},
				},
			},
		},
		{
			Name: "vpc_endpoint_id is specified",
			Content: `
resource "aws_route" "foo" {
	route_table_id = "rtb-1234abcd"
	vpc_endpoint_id = "vpce-12345678abcdefgh"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "core_network_arn is specified",
			Content: `
resource "aws_route" "foo" {
	route_table_id = "rtb-1234abcd"
	core_network_arn = "arn:aws:networkmanager::230703758040:core-network/core-network-0ce39423cf52b4c76"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsRouteNotSpecifiedTargetRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
