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

func TestParseListMapExpr(t *testing.T) {
	tests := []struct {
		name           string
		expr           string
		wantListShape  string // list shape (first arg)
		wantKeyShape   string // key constraint shape (second arg)
		wantValueShape string // value constraint shape (third arg)
		wantValid      bool
	}{
		{
			name:      "simple variable reference",
			expr:      "Tags",
			wantValid: false,
		},
		{
			name:           "listmap with all three arguments",
			expr:           `listmap(Tags, TagKey, TagValue)`,
			wantListShape:  "Tags",
			wantKeyShape:   "TagKey",
			wantValueShape: "TagValue",
			wantValid:      true,
		},
		{
			name:      "listmap with wrong function name",
			expr:      `othermap(Tags, TagKey, TagValue)`,
			wantValid: false,
		},
		{
			name:      "listmap with too few arguments",
			expr:      `listmap(Tags, TagKey)`,
			wantValid: false,
		},
		{
			name:      "listmap with too many arguments",
			expr:      `listmap(Tags, TagKey, TagValue, Extra)`,
			wantValid: false,
		},
		{
			name:      "listmap with no arguments",
			expr:      `listmap()`,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listShape, keyShape, valueShape := parseListMapExpr(tt.expr)

			gotValid := listShape != ""

			if gotValid != tt.wantValid {
				t.Errorf("valid: got %v, want %v", gotValid, tt.wantValid)
			}
			if listShape != tt.wantListShape {
				t.Errorf("listShape: got %q, want %q", listShape, tt.wantListShape)
			}
			if keyShape != tt.wantKeyShape {
				t.Errorf("keyShape: got %q, want %q", keyShape, tt.wantKeyShape)
			}
			if valueShape != tt.wantValueShape {
				t.Errorf("valueShape: got %q, want %q", valueShape, tt.wantValueShape)
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
			name: "list type recursing to map",
			shapes: makeShapes(map[string]interface{}{
				"TagMapList": map[string]interface{}{
					"type":   "list",
					"member": map[string]interface{}{"target": "com.example.test#TagMap"},
				},
				"TagMap": map[string]interface{}{
					"type":  "map",
					"key":   map[string]interface{}{"target": "com.example.test#TagKey"},
					"value": map[string]interface{}{"target": "com.example.test#TagValue"},
				},
				"TagKey":   makeStringShape("", 1, 128),
				"TagValue": makeStringShape("", 0, 256),
			}),
			shapeName: "TagMapList",
			wantKey:   true,
			wantValue: true,
		},
		{
			name: "list type recursing to structure (not supported)",
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
			wantKey:   false,
			wantValue: false,
		},
		{
			name: "structure (not supported without configuration)",
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
			wantKey:   false,
			wantValue: false,
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
