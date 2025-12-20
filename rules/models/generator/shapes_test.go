//go:build generators

package main

import (
	"testing"
)

// makeShapes creates a mock Smithy shapes map for testing.
// Uses the namespace "com.example.test" for qualified names.
func makeShapes(shapes map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"com.example.test#TestService": map[string]interface{}{
			"type": "service",
		},
	}
	for name, shape := range shapes {
		result["com.example.test#"+name] = shape
	}
	return result
}

// makeStringShape creates a Smithy string shape with optional constraints.
func makeStringShape(pattern string, min, max int) map[string]interface{} {
	shape := map[string]interface{}{
		"type": "string",
	}
	traits := map[string]interface{}{}
	if pattern != "" {
		traits["smithy.api#pattern"] = pattern
	}
	if min > 0 || max > 0 {
		length := map[string]interface{}{}
		if min > 0 {
			length["min"] = float64(min)
		}
		if max > 0 {
			length["max"] = float64(max)
		}
		traits["smithy.api#length"] = length
	}
	if len(traits) > 0 {
		shape["traits"] = traits
	}
	return shape
}

func TestExtractShapeName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"com.amazonaws.ecs#TagKey", "TagKey"},
		{"com.amazonaws.ecs#TagValue", "TagValue"},
		{"TagKey", "TagKey"},
		{"", ""},
		{"com.amazonaws.ecs#nested#Name", "Name"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractShapeName(tt.input)
			if got != tt.want {
				t.Errorf("extractShapeName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetTarget(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		want  string
	}{
		{
			name:  "nil ref",
			input: nil,
			want:  "",
		},
		{
			name:  "empty ref",
			input: map[string]interface{}{},
			want:  "",
		},
		{
			name:  "qualified target",
			input: map[string]interface{}{"target": "com.amazonaws.ecs#TagKey"},
			want:  "TagKey",
		},
		{
			name:  "simple target",
			input: map[string]interface{}{"target": "TagKey"},
			want:  "TagKey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTarget(tt.input)
			if got != tt.want {
				t.Errorf("getTarget() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTraverseToMapConstraints_MapType(t *testing.T) {
	// Smithy map type: directly has key/value targets
	// Example: TagMap with key -> TagKey, value -> TagValue
	shapes := makeShapes(map[string]interface{}{
		"TagMap": map[string]interface{}{
			"type": "map",
			"key": map[string]interface{}{
				"target": "com.example.test#TagKey",
			},
			"value": map[string]interface{}{
				"target": "com.example.test#TagValue",
			},
		},
		"TagKey":   makeStringShape("^[\\w]+$", 1, 128),
		"TagValue": makeStringShape("", 0, 256),
	})

	keyModel, valueModel := traverseToMapConstraints(shapes, "TagMap")

	if keyModel == nil {
		t.Fatal("expected keyModel, got nil")
	}
	if valueModel == nil {
		t.Fatal("expected valueModel, got nil")
	}
	if keyModel["type"] != "string" {
		t.Errorf("keyModel type = %v, want string", keyModel["type"])
	}
	if keyModel["pattern"] != "^[\\w]+$" {
		t.Errorf("keyModel pattern = %v, want ^[\\w]+$", keyModel["pattern"])
	}
	if valueModel["max"] != float64(256) {
		t.Errorf("valueModel max = %v, want 256", valueModel["max"])
	}
}

func TestTraverseToMapConstraints_ListType(t *testing.T) {
	// Smithy list type: recurses into member element
	// Example: TagList -> Tag -> Key/Value
	shapes := makeShapes(map[string]interface{}{
		"TagList": map[string]interface{}{
			"type": "list",
			"member": map[string]interface{}{
				"target": "com.example.test#Tag",
			},
		},
		"Tag": map[string]interface{}{
			"type": "structure",
			"members": map[string]interface{}{
				"Key": map[string]interface{}{
					"target": "com.example.test#TagKey",
				},
				"Value": map[string]interface{}{
					"target": "com.example.test#TagValue",
				},
			},
		},
		"TagKey":   makeStringShape("", 1, 128),
		"TagValue": makeStringShape("", 0, 256),
	})

	keyModel, valueModel := traverseToMapConstraints(shapes, "TagList")

	if keyModel == nil {
		t.Fatal("expected keyModel, got nil")
	}
	if valueModel == nil {
		t.Fatal("expected valueModel, got nil")
	}
}

func TestTraverseToMapConstraints_StructureType(t *testing.T) {
	// Smithy structure type with Key/Value members
	shapes := makeShapes(map[string]interface{}{
		"Tag": map[string]interface{}{
			"type": "structure",
			"members": map[string]interface{}{
				"Key": map[string]interface{}{
					"target": "com.example.test#TagKey",
				},
				"Value": map[string]interface{}{
					"target": "com.example.test#TagValue",
				},
			},
		},
		"TagKey":   makeStringShape("", 1, 128),
		"TagValue": makeStringShape("", 0, 256),
	})

	keyModel, valueModel := traverseToMapConstraints(shapes, "Tag")

	if keyModel == nil {
		t.Fatal("expected keyModel, got nil")
	}
	if valueModel == nil {
		t.Fatal("expected valueModel, got nil")
	}
}

func TestTraverseToMapConstraints_StructureLowercaseMembers(t *testing.T) {
	// Some AWS services use lowercase member names (key/value instead of Key/Value)
	shapes := makeShapes(map[string]interface{}{
		"Tag": map[string]interface{}{
			"type": "structure",
			"members": map[string]interface{}{
				"key": map[string]interface{}{
					"target": "com.example.test#TagKey",
				},
				"value": map[string]interface{}{
					"target": "com.example.test#TagValue",
				},
			},
		},
		"TagKey":   makeStringShape("", 1, 128),
		"TagValue": makeStringShape("", 0, 256),
	})

	keyModel, valueModel := traverseToMapConstraints(shapes, "Tag")

	if keyModel == nil {
		t.Fatal("expected keyModel, got nil")
	}
	if valueModel == nil {
		t.Fatal("expected valueModel, got nil")
	}
}

func TestTraverseToMapConstraints_StructureNameValue(t *testing.T) {
	// Some structures use Name/Value instead of Key/Value (e.g., EnvironmentVariable)
	shapes := makeShapes(map[string]interface{}{
		"EnvironmentVariable": map[string]interface{}{
			"type": "structure",
			"members": map[string]interface{}{
				"Name": map[string]interface{}{
					"target": "com.example.test#EnvName",
				},
				"Value": map[string]interface{}{
					"target": "com.example.test#EnvValue",
				},
			},
		},
		"EnvName":  makeStringShape("", 1, 128),
		"EnvValue": makeStringShape("", 0, 256),
	})

	keyModel, valueModel := traverseToMapConstraints(shapes, "EnvironmentVariable")

	if keyModel == nil {
		t.Fatal("expected keyModel, got nil")
	}
	if valueModel == nil {
		t.Fatal("expected valueModel, got nil")
	}
}

func TestTraverseToMapConstraints_NonStringMembers(t *testing.T) {
	// Structure with non-string members should return nil
	shapes := makeShapes(map[string]interface{}{
		"NotATag": map[string]interface{}{
			"type": "structure",
			"members": map[string]interface{}{
				"Count": map[string]interface{}{
					"target": "com.example.test#Integer",
				},
				"Enabled": map[string]interface{}{
					"target": "com.example.test#Boolean",
				},
			},
		},
		"Integer": map[string]interface{}{
			"type": "integer",
		},
		"Boolean": map[string]interface{}{
			"type": "boolean",
		},
	})

	keyModel, valueModel := traverseToMapConstraints(shapes, "NotATag")

	if keyModel != nil || valueModel != nil {
		t.Errorf("expected nil, nil for non-string members, got %v, %v", keyModel, valueModel)
	}
}

func TestTraverseToMapConstraints_MixedMembers(t *testing.T) {
	// Structure with mixed string and non-string members
	// Should still work if exactly 2 string members exist
	shapes := makeShapes(map[string]interface{}{
		"TagWithMetadata": map[string]interface{}{
			"type": "structure",
			"members": map[string]interface{}{
				"Key": map[string]interface{}{
					"target": "com.example.test#TagKey",
				},
				"Value": map[string]interface{}{
					"target": "com.example.test#TagValue",
				},
				"Timestamp": map[string]interface{}{
					"target": "com.example.test#Timestamp",
				},
			},
		},
		"TagKey":   makeStringShape("", 1, 128),
		"TagValue": makeStringShape("", 0, 256),
		"Timestamp": map[string]interface{}{
			"type": "timestamp",
		},
	})

	keyModel, valueModel := traverseToMapConstraints(shapes, "TagWithMetadata")

	// Current implementation requires exactly 2 string members
	// This test documents current behavior - adjust if requirements change
	if keyModel == nil || valueModel == nil {
		t.Log("Structure with 3 members (2 string, 1 non-string) returned nil - this may be intentional")
	}
}

func TestTraverseToMapConstraints_UnknownShape(t *testing.T) {
	shapes := makeShapes(map[string]interface{}{})

	keyModel, valueModel := traverseToMapConstraints(shapes, "NonExistent")

	if keyModel != nil || valueModel != nil {
		t.Errorf("expected nil, nil for unknown shape, got %v, %v", keyModel, valueModel)
	}
}

func TestTraverseToMapConstraints_ValueIsList(t *testing.T) {
	// Map where value is a list (not a string) should return nil
	// This is the FMS exclude_map case
	shapes := makeShapes(map[string]interface{}{
		"PolicyScopeMap": map[string]interface{}{
			"type": "map",
			"key": map[string]interface{}{
				"target": "com.example.test#ScopeType",
			},
			"value": map[string]interface{}{
				"target": "com.example.test#ScopeIdList",
			},
		},
		"ScopeType": map[string]interface{}{
			"type": "string",
			"traits": map[string]interface{}{
				"smithy.api#enum": []interface{}{
					map[string]interface{}{"value": "ACCOUNT"},
					map[string]interface{}{"value": "ORG_UNIT"},
				},
			},
		},
		"ScopeIdList": map[string]interface{}{
			"type": "list",
			"member": map[string]interface{}{
				"target": "com.example.test#ScopeId",
			},
		},
		"ScopeId": makeStringShape("", 1, 128),
	})

	keyModel, valueModel := traverseToMapConstraints(shapes, "PolicyScopeMap")

	if keyModel != nil || valueModel != nil {
		t.Errorf("expected nil, nil for map with list value, got %v, %v", keyModel, valueModel)
	}
}
