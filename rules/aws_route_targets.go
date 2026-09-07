package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
)

// routeTargetAttributes are the mutually exclusive routing target arguments of
// an aws_route resource, of which AWS requires exactly one.
// aws_route_not_specified_target and aws_route_specified_multiple_targets check
// that requirement from opposite directions and must agree on the set, since an
// attribute missing from the schema is invisible to both.
//
// The provider expresses the mutual exclusion as an ExactlyOneOf, which does not
// cross the plugin protocol, so the set cannot be derived from the provider
// schema and is maintained here instead. TestRouteTargetAttributes fails when
// the schema gains an aws_route argument this list does not classify.
//
// Sorted alphabetically: routeTargetNotSpecifiedMessage lists them in order.
var routeTargetAttributes = []hclext.AttributeSchema{
	{Name: "carrier_gateway_id"},
	{Name: "core_network_arn"},
	{Name: "egress_only_gateway_id"},
	{Name: "gateway_id"},
	{Name: "instance_id"},
	{Name: "local_gateway_id"},
	{Name: "nat_gateway_id"},
	{Name: "network_interface_id"},
	{Name: "odb_network_arn"},
	{Name: "transit_gateway_id"},
	{Name: "vpc_endpoint_id"},
	{Name: "vpc_peering_connection_id"},
}

var routeTargetNotSpecifiedMessage = fmt.Sprintf(
	"The routing target is not specified, each aws_route must contain either %s.",
	routeTargetAlternatives(),
)

func routeTargetNames() []string {
	names := make([]string, len(routeTargetAttributes))
	for i, attribute := range routeTargetAttributes {
		names[i] = attribute.Name
	}
	return names
}

func routeTargetAlternatives() string {
	names := routeTargetNames()
	return strings.Join(names[:len(names)-1], ", ") + " or " + names[len(names)-1]
}
