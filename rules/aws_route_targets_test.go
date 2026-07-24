package rules

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/terraform-linters/tflint-ruleset-aws/rules/genutils"
)

// knownNonTargetAttributes are the aws_route arguments that are not routing
// targets. Together with routeTargetAttributes they account for every settable
// attribute of the resource, so an argument added by the provider fails
// TestRouteTargetAttributes until it is classified as one or the other. Both
// aws_route target rules once silently missed core_network_arn and
// odb_network_arn for exactly that reason.
var knownNonTargetAttributes = []string{
	"destination_cidr_block",
	"destination_ipv6_cidr_block",
	"destination_prefix_list_id",
	"id",
	"region",
	"route_table_id",
}

func TestRouteTargetAttributes(t *testing.T) {
	provider := genutils.LoadProviderSchema(filepath.Join("..", "tools", "provider-schema", "schema.json"))
	route, exists := provider.ResourceSchemas["aws_route"]
	if !exists {
		t.Fatal("aws_route is missing from the provider schema")
	}
	names := routeTargetNames()

	t.Run("every routing target exists in the provider schema", func(t *testing.T) {
		for _, name := range names {
			if route.Block.Attributes[name] == nil {
				t.Errorf("aws_route no longer has a %s attribute", name)
			}
		}
	})

	t.Run("every settable attribute is classified", func(t *testing.T) {
		for name, attribute := range route.Block.Attributes {
			if !attribute.Optional && !attribute.Required {
				continue
			}
			if slices.Contains(names, name) || slices.Contains(knownNonTargetAttributes, name) {
				continue
			}
			t.Errorf("aws_route attribute %s is neither a routing target nor a known non-target: add it to routeTargetAttributes or knownNonTargetAttributes", name)
		}
	})

	t.Run("routing targets are sorted", func(t *testing.T) {
		if !slices.IsSorted(names) {
			t.Errorf("routeTargetAttributes must be sorted because routeTargetNotSpecifiedMessage lists them in order, got %v", names)
		}
	})
}
