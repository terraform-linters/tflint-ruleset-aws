package elasticache_node_types

import "testing"

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

// TestPreviousGenerationValid holds the two rules to one story. A previous
// generation node type that is not also valid would have previous_type warning
// on a node type invalid_type rejects, which is the contradiction that kept
// this list hardcoded.
func TestPreviousGenerationValid(t *testing.T) {
	if len(nodeTypes.PreviousGeneration) == 0 {
		t.Fatal("no previous generation node types")
	}

	for nodeType := range nodeTypes.PreviousGeneration {
		if !nodeTypes.All[nodeType] {
			t.Errorf("previous generation node type %s is not valid", nodeType)
		}
	}
}
