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
			if got := extractShapeName(tt.input); got != tt.want {
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
		{"nil ref", nil, ""},
		{"empty ref", map[string]interface{}{}, ""},
		{"qualified target", map[string]interface{}{"target": "com.amazonaws.ecs#TagKey"}, "TagKey"},
		{"simple target", map[string]interface{}{"target": "TagKey"}, "TagKey"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTarget(tt.input); got != tt.want {
				t.Errorf("getTarget() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTraverseToMapConstraints(t *testing.T) {
	tests := []struct {
		name      string
		shapes    map[string]interface{}
		shapeName string
		wantKey   bool
		wantValue bool
	}{
		{
			name: "map type with key/value",
			shapes: makeShapes(map[string]interface{}{
				"TagMap": map[string]interface{}{
					"type": "map",
					"key":  map[string]interface{}{"target": "com.example.test#TagKey"},
					"value": map[string]interface{}{"target": "com.example.test#TagValue"},
				},
				"TagKey":   makeStringShape("^[\\w]+$", 1, 128),
				"TagValue": makeStringShape("", 0, 256),
			}),
			shapeName: "TagMap",
			wantKey:   true,
			wantValue: true,
		},
		{
			name: "list type recursing to structure",
			shapes: makeShapes(map[string]interface{}{
				"TagList": map[string]interface{}{
					"type":   "list",
					"member": map[string]interface{}{"target": "com.example.test#Tag"},
				},
				"Tag": map[string]interface{}{
					"type": "structure",
					"members": map[string]interface{}{
						"Key":   map[string]interface{}{"target": "com.example.test#TagKey"},
						"Value": map[string]interface{}{"target": "com.example.test#TagValue"},
					},
				},
				"TagKey":   makeStringShape("", 1, 128),
				"TagValue": makeStringShape("", 0, 256),
			}),
			shapeName: "TagList",
			wantKey:   true,
			wantValue: true,
		},
		{
			name: "structure with Key/Value members",
			shapes: makeShapes(map[string]interface{}{
				"Tag": map[string]interface{}{
					"type": "structure",
					"members": map[string]interface{}{
						"Key":   map[string]interface{}{"target": "com.example.test#TagKey"},
						"Value": map[string]interface{}{"target": "com.example.test#TagValue"},
					},
				},
				"TagKey":   makeStringShape("", 1, 128),
				"TagValue": makeStringShape("", 0, 256),
			}),
			shapeName: "Tag",
			wantKey:   true,
			wantValue: true,
		},
		{
			name: "structure with lowercase key/value members",
			shapes: makeShapes(map[string]interface{}{
				"Tag": map[string]interface{}{
					"type": "structure",
					"members": map[string]interface{}{
						"key":   map[string]interface{}{"target": "com.example.test#TagKey"},
						"value": map[string]interface{}{"target": "com.example.test#TagValue"},
					},
				},
				"TagKey":   makeStringShape("", 1, 128),
				"TagValue": makeStringShape("", 0, 256),
			}),
			shapeName: "Tag",
			wantKey:   true,
			wantValue: true,
		},
		{
			name: "structure with Name/Value members",
			shapes: makeShapes(map[string]interface{}{
				"EnvironmentVariable": map[string]interface{}{
					"type": "structure",
					"members": map[string]interface{}{
						"Name":  map[string]interface{}{"target": "com.example.test#EnvName"},
						"Value": map[string]interface{}{"target": "com.example.test#EnvValue"},
					},
				},
				"EnvName":  makeStringShape("", 1, 128),
				"EnvValue": makeStringShape("", 0, 256),
			}),
			shapeName: "EnvironmentVariable",
			wantKey:   true,
			wantValue: true,
		},
		{
			name: "structure with non-string members only",
			shapes: makeShapes(map[string]interface{}{
				"NotATag": map[string]interface{}{
					"type": "structure",
					"members": map[string]interface{}{
						"Count":   map[string]interface{}{"target": "com.example.test#Integer"},
						"Enabled": map[string]interface{}{"target": "com.example.test#Boolean"},
					},
				},
				"Integer": map[string]interface{}{"type": "integer"},
				"Boolean": map[string]interface{}{"type": "boolean"},
			}),
			shapeName: "NotATag",
			wantKey:   false,
			wantValue: false,
		},
		{
			name: "structure with mixed members (2 string + 1 non-string)",
			shapes: makeShapes(map[string]interface{}{
				"TagWithMetadata": map[string]interface{}{
					"type": "structure",
					"members": map[string]interface{}{
						"Key":       map[string]interface{}{"target": "com.example.test#TagKey"},
						"Value":     map[string]interface{}{"target": "com.example.test#TagValue"},
						"Timestamp": map[string]interface{}{"target": "com.example.test#Timestamp"},
					},
				},
				"TagKey":    makeStringShape("", 1, 128),
				"TagValue":  makeStringShape("", 0, 256),
				"Timestamp": map[string]interface{}{"type": "timestamp"},
			}),
			shapeName: "TagWithMetadata",
			wantKey:   true,
			wantValue: true,
		},
		{
			name:      "unknown shape",
			shapes:    makeShapes(map[string]interface{}{}),
			shapeName: "NonExistent",
			wantKey:   false,
			wantValue: false,
		},
		{
			name: "map with list value (not string)",
			shapes: makeShapes(map[string]interface{}{
				"PolicyScopeMap": map[string]interface{}{
					"type":  "map",
					"key":   map[string]interface{}{"target": "com.example.test#ScopeType"},
					"value": map[string]interface{}{"target": "com.example.test#ScopeIdList"},
				},
				"ScopeType": map[string]interface{}{
					"type": "string",
					"traits": map[string]interface{}{
						"smithy.api#enum": []interface{}{
							map[string]interface{}{"value": "ACCOUNT"},
						},
					},
				},
				"ScopeIdList": map[string]interface{}{
					"type":   "list",
					"member": map[string]interface{}{"target": "com.example.test#ScopeId"},
				},
				"ScopeId": makeStringShape("", 1, 128),
			}),
			shapeName: "PolicyScopeMap",
			wantKey:   false,
			wantValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyModel, valueModel := traverseToMapConstraints(tt.shapes, tt.shapeName)

			gotKey := keyModel != nil
			gotValue := valueModel != nil

			if gotKey != tt.wantKey {
				t.Errorf("keyModel: got %v, want %v", gotKey, tt.wantKey)
			}
			if gotValue != tt.wantValue {
				t.Errorf("valueModel: got %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}
