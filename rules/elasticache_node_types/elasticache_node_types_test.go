package elasticache_node_types

import (
	"testing"

	"github.com/terraform-linters/tflint-ruleset-aws/rules/stringset"
)

func TestValid(t *testing.T) {
	for _, tc := range []struct {
		name     string
		nodeType string
		expected bool
	}{
		{
			name:     "current generation node type",
			nodeType: "cache.t4g.micro",
			expected: true,
		},
		{
			name:     "latest generation node type",
			nodeType: "cache.r8g.16xlarge",
			expected: true,
		},
		{
			name:     "previous generation node type",
			nodeType: "cache.t1.micro",
			expected: true,
		},
		{
			name:     "data tiering node type",
			nodeType: "cache.r6gd.xlarge",
			expected: true,
		},
		{
			name:     "size a data tiering family does not offer",
			nodeType: "cache.r6gd.large",
			expected: false,
		},
		{
			name:     "EC2 instance type without the cache prefix",
			nodeType: "t4g.micro",
			expected: false,
		},
		{
			name:     "unknown size",
			nodeType: "cache.t3.mini",
			expected: false,
		},
		{
			name:     "empty",
			nodeType: "",
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if result := Valid(tc.nodeType); result != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, result)
			}
		})
	}
}

func TestPreviousGeneration(t *testing.T) {
	for _, tc := range []struct {
		name     string
		nodeType string
		expected bool
	}{
		{
			name:     "previous generation node type",
			nodeType: "cache.m1.small",
			expected: true,
		},
		{
			name:     "current generation node type",
			nodeType: "cache.t2.micro",
			expected: false,
		},
		{
			name:     "unknown size in a previous generation family",
			nodeType: "cache.m1.unknownsize",
			expected: false,
		},
		{
			name:     "previous generation family under another prefix",
			nodeType: "foo.m1.large",
			expected: false,
		},
		{
			name:     "family rather than a node type",
			nodeType: "m1",
			expected: false,
		},
		{
			name:     "empty",
			nodeType: "",
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if result := PreviousGeneration(tc.nodeType); result != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, result)
			}
		})
	}
}

// TestEmbeddedNodeTypes guards the read path. The generator's floors only run
// when someone regenerates, so a truncated or partially written JSON file would
// otherwise ship a plugin that rejects every node type.
func TestEmbeddedNodeTypes(t *testing.T) {
	for _, tc := range []struct {
		name    string
		types   stringset.Set
		minimum int
	}{
		{
			name:    "node types",
			types:   nodeTypes.All,
			minimum: 100,
		},
		{
			name:    "previous generation node types",
			types:   nodeTypes.PreviousGeneration,
			minimum: 18,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.types) < tc.minimum {
				t.Errorf("expected at least %d, got %d", tc.minimum, len(tc.types))
			}
		})
	}
}
