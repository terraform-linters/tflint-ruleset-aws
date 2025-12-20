//go:build generators

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
	utils "github.com/terraform-linters/tflint-ruleset-aws/rules/generator-utils"
)

type mappingFile struct {
	Import   string    `hcl:"import"`
	Mappings []mapping `hcl:"mapping,block"`
	Tests    []test    `hcl:"test,block"`
}

type mapping struct {
	Resource string         `hcl:"resource,label"`
	Attrs    hcl.Attributes `hcl:",remain"`
}

type test struct {
	Resource  string `hcl:"resource,label"`
	Attribute string `hcl:"attribute,label"`
	OK        string `hcl:"ok"`
	NG        string `hcl:"ng"`
}

// smithyShape represents a Smithy model shape definition.
type smithyShape struct {
	Type    string                 `json:"type"`
	Traits  map[string]interface{} `json:"traits,omitempty"`
	Members map[string]smithyMember `json:"members,omitempty"`
}

// smithyMember represents a member of a Smithy shape.
type smithyMember struct {
	Target string                 `json:"target,omitempty"`
	Traits map[string]interface{} `json:"traits,omitempty"`
}

// smithyLengthTrait represents the smithy.api#length trait.
type smithyLengthTrait struct {
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
}

// smithyEnumItem represents an enum value in smithy.api#enum trait.
type smithyEnumItem struct {
	Value string `json:"value"`
	Name  string `json:"name,omitempty"`
}

func main() {
	files, err := filepath.Glob("./mappings/*.hcl")
	if err != nil {
		panic(err)
	}

	var mappingFiles []mappingFile
	for _, file := range files {
		parser := hclparse.NewParser()
		f, diags := parser.ParseHCLFile(file)
		if diags.HasErrors() {
			panic(diags)
		}

		var mf mappingFile
		diags = gohcl.DecodeBody(f.Body, nil, &mf)
		if diags.HasErrors() {
			panic(diags)
		}
		mappingFiles = append(mappingFiles, mf)
	}

	awsProvider := utils.LoadProviderSchema("../../tools/provider-schema/schema.json")

	var generatedRules []string
	for _, mappingFile := range mappingFiles {
		raw, err := os.ReadFile(mappingFile.Import)
		if err != nil {
			panic(err)
		}

		var api map[string]interface{}
		err = json.Unmarshal(raw, &api)
		if err != nil {
			panic(err)
		}
		shapes := api["shapes"].(map[string]interface{})

		evalCtx := buildEvalContext()

		for _, mapping := range mappingFile.Mappings {
			for attribute, value := range mapping.Attrs {
				fmt.Printf("Checking `%s.%s`\n", mapping.Resource, attribute)

				// Extract shape name from the expression
				vars := value.Expr.Variables()
				if len(vars) == 0 {
					continue
				}
				if len(vars) > 1 {
					varNames := make([]string, len(vars))
					for i, v := range vars {
						varNames[i] = v.RootName()
					}
					fmt.Fprintf(os.Stderr, "Error: `%s.%s` expression references %d variables (%v), expected exactly 1\n",
						mapping.Resource, attribute, len(vars), varNames)
					os.Exit(1)
				}
				shapeName := vars[0].RootName()
				if shapeName == "any" {
					continue
				}

				model := findShape(shapes, shapeName)
				if model == nil {
					fmt.Printf("Shape %q not found, skipping\n", shapeName)
					continue
				}

				// Populate the eval context with this shape's enum values
				if enumValues, ok := model["enum"].([]string); ok {
					// Use fresh Variables map for each evaluation to avoid pollution
					evalCtx.Variables = map[string]cty.Value{
						shapeName: stringsToCtyList(enumValues),
					}

					// Evaluate the expression to get transformed enum values
					result, diags := value.Expr.Value(evalCtx)
					if !diags.HasErrors() && result.Type().IsListType() {
						transformedEnums := make([]string, 0, result.LengthInt())
						for it := result.ElementIterator(); it.Next(); {
							_, val := it.Element()
							transformedEnums = append(transformedEnums, val.AsString())
						}

						// Create a new model with transformed enums
						transformedModel := make(map[string]interface{}, len(model))
						for k, v := range model {
							transformedModel[k] = v
						}
						transformedModel["enum"] = transformedEnums
						model = transformedModel
					}
				}

				schema, err := fetchSchema(mapping.Resource, attribute, model, awsProvider)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error processing `%s.%s`: %v\n", mapping.Resource, attribute, err)
					os.Exit(1)
				}
				if validMapping(model) {
					fmt.Printf("Generating rule for `%s.%s`\n", mapping.Resource, attribute)
					generateRuleFile(mapping.Resource, attribute, model, schema)
					for _, test := range mappingFile.Tests {
						if mapping.Resource == test.Resource && attribute == test.Attribute {
							generateRuleTestFile(mapping.Resource, attribute, model, test)
						}
					}
					generatedRules = append(generatedRules, makeRuleName(mapping.Resource, attribute))
				} else {
					// Try traversing to find tag-like constraints
					keyModel, valueModel := traverseToTagConstraints(shapes, shapeName)
					if keyModel != nil && valueModel != nil {
						if validMapping(keyModel) || validMapping(valueModel) {
							fmt.Printf("Generating map rule for `%s.%s`\n", mapping.Resource, attribute)
							if generateMapRuleFile(mapping.Resource, attribute, keyModel, valueModel, schema) {
								generatedRules = append(generatedRules, makeRuleName(mapping.Resource, attribute))
							}
						}
					}
				}
			}
		}
	}

	sort.Strings(generatedRules)
	generateProviderFile(generatedRules)
	generateDocFile(generatedRules)
}

func fetchSchema(resource, attribute string, model map[string]interface{}, provider *tfjson.ProviderSchema) (*tfjson.SchemaAttribute, error) {
	resourceSchema, ok := provider.ResourceSchemas[resource]
	if !ok {
		return nil, fmt.Errorf("resource `%s` not found in the Terraform schema", resource)
	}
	attrSchema, ok := resourceSchema.Block.Attributes[attribute]
	if !ok {
		if _, ok := resourceSchema.Block.NestedBlocks[attribute]; !ok {
			return nil, fmt.Errorf("`%s.%s` not found in the Terraform schema", resource, attribute)
		}
	}

	switch model["type"].(string) {
	case "string":
		if attrSchema != nil {
			ty := attrSchema.AttributeType.FriendlyName()
			if ty != "string" && ty != "number" {
				return nil, fmt.Errorf("`%s.%s` is expected as string, but not (%s)", resource, attribute, ty)
			}
		}
	default:
		// noop
	}

	return attrSchema, nil
}

func validMapping(model map[string]interface{}) bool {
	switch model["type"].(string) {
	case "string":
		if _, ok := model["max"]; ok {
			return true
		}
		if min, ok := model["min"]; ok && int(min.(float64)) > 2 {
			return true
		}
		if _, ok := model["pattern"]; ok {
			return true
		}
		if _, ok := model["enum"]; ok {
			return true
		}
		return false
	default:
		// Unsupported types
		return false
	}
}

// findShape locates a shape in Smithy format with namespace-qualified lookup
func findShape(shapes map[string]interface{}, shapeName string) map[string]interface{} {
	raw := findRawShape(shapes, shapeName)
	if raw == nil {
		return nil
	}
	return convertSmithyShape(raw)
}

// findRawShape locates a shape without converting it (needed for traversal)
func findRawShape(shapes map[string]interface{}, shapeName string) map[string]interface{} {
	// Try with service namespace qualification (Smithy format)
	serviceNamespace := extractServiceNamespace(shapes)
	if serviceNamespace != "" {
		qualifiedName := fmt.Sprintf("%s#%s", serviceNamespace, shapeName)
		if shape, ok := shapes[qualifiedName]; ok {
			return shape.(map[string]interface{})
		}
	}

	// Fallback to direct lookup (legacy format or unqualified shapes)
	if shape, ok := shapes[shapeName]; ok {
		if shapeMap, ok := shape.(map[string]interface{}); ok {
			return shapeMap
		}
	}

	return nil
}

// extractShapeName extracts the simple name from a fully qualified shape reference.
// e.g., "com.amazonaws.ecs#TagKey" -> "TagKey"
func extractShapeName(qualifiedName string) string {
	if idx := strings.LastIndex(qualifiedName, "#"); idx >= 0 {
		return qualifiedName[idx+1:]
	}
	return qualifiedName
}

// tagStructureMembers defines the member names used to identify tag-like structures
// in Smithy models. AWS services use varying capitalization for these members
// (e.g., "key"/"Key", "value"/"Value"), so lookups are case-insensitive.
var tagStructureMembers = struct {
	Key   string
	Value string
}{
	Key:   "key",
	Value: "value",
}

// getMemberCaseInsensitive looks up a member in a Smithy members map using
// case-insensitive matching. Returns nil if not found.
func getMemberCaseInsensitive(members map[string]interface{}, name string) map[string]interface{} {
	for k, v := range members {
		if strings.EqualFold(k, name) {
			if m, ok := v.(map[string]interface{}); ok {
				return m
			}
		}
	}
	return nil
}

// getTarget extracts the target shape name from a Smithy reference.
func getTarget(ref map[string]interface{}) string {
	if ref == nil {
		return ""
	}
	target, _ := ref["target"].(string)
	return extractShapeName(target)
}

// traverseToTagConstraints traverses Smithy shapes to find key/value string constraints.
// It handles three shape types:
//   - map: directly has key/value targets (e.g., TagMap)
//   - list: recurses into the member element (e.g., TagList -> Tag)
//   - structure: looks for key/value members (e.g., Tag with Key/Value fields)
//
// Returns nil, nil if the shape doesn't resolve to a string key/value pair.
func traverseToTagConstraints(shapes map[string]interface{}, shapeName string) (keyModel, valueModel map[string]interface{}) {
	raw := findRawShape(shapes, shapeName)
	if raw == nil {
		return nil, nil
	}

	shapeType, _ := raw["type"].(string)

	switch shapeType {
	case "map":
		keyRef, _ := raw["key"].(map[string]interface{})
		valueRef, _ := raw["value"].(map[string]interface{})
		return resolveKeyValuePair(shapes, getTarget(keyRef), getTarget(valueRef))

	case "list":
		memberRef, _ := raw["member"].(map[string]interface{})
		if target := getTarget(memberRef); target != "" {
			return traverseToTagConstraints(shapes, target)
		}

	case "structure":
		members, _ := raw["members"].(map[string]interface{})
		if members == nil {
			return nil, nil
		}
		keyMember := getMemberCaseInsensitive(members, tagStructureMembers.Key)
		valueMember := getMemberCaseInsensitive(members, tagStructureMembers.Value)
		return resolveKeyValuePair(shapes, getTarget(keyMember), getTarget(valueMember))
	}

	return nil, nil
}

// resolveKeyValuePair resolves key and value shape names to their models,
// returning nil if either is missing or not a string type.
func resolveKeyValuePair(shapes map[string]interface{}, keyName, valueName string) (map[string]interface{}, map[string]interface{}) {
	if keyName == "" || valueName == "" {
		return nil, nil
	}
	keyModel := findShape(shapes, keyName)
	valueModel := findShape(shapes, valueName)
	if keyModel == nil || valueModel == nil {
		return nil, nil
	}
	if keyModel["type"] != "string" || valueModel["type"] != "string" {
		return nil, nil
	}
	return keyModel, valueModel
}

// extractServiceNamespace extracts the namespace from the Smithy service definition
func extractServiceNamespace(shapes map[string]interface{}) string {
	for shapeName, shape := range shapes {
		if shapeMap, ok := shape.(map[string]interface{}); ok {
			if shapeType, ok := shapeMap["type"].(string); ok && shapeType == "service" {
				// Extract namespace from shape name (e.g., "com.amazonaws.acmpca#ACMPrivateCA")
				if parts := strings.Split(shapeName, "#"); len(parts) == 2 {
					return parts[0]
				}
			}
		}
	}
	return ""
}

// convertSmithyShape converts Smithy model format to internal format.
// Smithy uses traits for metadata while our internal format uses direct fields.
func convertSmithyShape(rawShape map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Parse the raw shape into typed struct
	shapeBytes, err := json.Marshal(rawShape)
	if err != nil {
		return result
	}

	var shape smithyShape
	if err := json.Unmarshal(shapeBytes, &shape); err != nil {
		return result
	}

	result["type"] = shape.Type

	// Extract length constraints from traits
	if lengthData, ok := shape.Traits["smithy.api#length"]; ok {
		lengthBytes, _ := json.Marshal(lengthData)
		var length smithyLengthTrait
		if json.Unmarshal(lengthBytes, &length) == nil {
			if length.Min != nil {
				result["min"] = *length.Min
			}
			if length.Max != nil {
				result["max"] = *length.Max
			}
		}
	}

	// Extract pattern constraint
	if pattern, ok := shape.Traits["smithy.api#pattern"].(string); ok {
		result["pattern"] = pattern
	}

	// Extract enum as trait (older Smithy style)
	if enumData, ok := shape.Traits["smithy.api#enum"]; ok {
		enumBytes, _ := json.Marshal(enumData)
		var enumItems []smithyEnumItem
		if json.Unmarshal(enumBytes, &enumItems) == nil {
			enumValues := make([]string, 0, len(enumItems))
			for _, item := range enumItems {
				enumValues = append(enumValues, item.Value)
			}
			sort.Strings(enumValues)
			result["enum"] = enumValues
		}
	}

	// Handle enum as type (newer Smithy style: type="enum" with members)
	if shape.Type == "enum" && len(shape.Members) > 0 {
		memberNames := make([]string, 0, len(shape.Members))
		for memberName := range shape.Members {
			memberNames = append(memberNames, memberName)
		}
		sort.Strings(memberNames)

		enumValues := make([]string, 0, len(memberNames))
		for _, memberName := range memberNames {
			member := shape.Members[memberName]
			enumValue := memberName

			// Check for explicit enumValue in member traits
			if enumValueData, ok := member.Traits["smithy.api#enumValue"]; ok {
				if enumValueStr, ok := enumValueData.(string); ok {
					enumValue = enumValueStr
				}
			}
			enumValues = append(enumValues, enumValue)
		}

		result["enum"] = enumValues
		result["type"] = "string" // Normalize enum type to string
	}

	return result
}

// stringsToCtyList converts a slice of strings to a cty list value
func stringsToCtyList(values []string) cty.Value {
	if len(values) == 0 {
		return cty.ListValEmpty(cty.String)
	}
	ctyValues := make([]cty.Value, 0, len(values))
	for _, v := range values {
		ctyValues = append(ctyValues, cty.StringVal(v))
	}
	return cty.ListVal(ctyValues)
}

// HCL Transform System
// ====================
//
// The generator supports HCL function-based transforms for enum values in mapping files.
// This allows enum values from Smithy models to be transformed to match Terraform provider
// expectations without hardcoding transformation logic in Go.
//
// Available Functions:
//   - uppercase(list) - Converts all strings in list to uppercase
//   - replace(list, old, new) - Replaces substring in all strings
//
// Usage Examples in mappings/*.hcl:
//   compression_type = uppercase(CompressionTypeValue)
//     -> ["gzip", "none"] becomes ["GZIP", "NONE"]
//
//   encryption_mode = uppercase(replace(EncryptionModeValue, "-", "_"))
//     -> ["sse-kms", "sse-s3"] becomes ["SSE_KMS", "SSE_S3"]
//
// Adding New Transform Functions:
//   1. Add function to buildEvalContext() Functions map
//   2. Wrap stdlib function with makeListTransformFunction()
//   3. Example: "lowercase": makeListTransformFunction(stdlib.LowerFunc)
//
// Technical Details:
//   - Transform functions operate on lists of strings (enum values)
//   - Each expression must reference exactly one shape variable
//   - Functions can be composed: uppercase(replace(x, "-", "_"))
//   - HCL's native expression evaluation handles all transforms

// buildEvalContext creates an HCL evaluation context with transform functions
func buildEvalContext() *hcl.EvalContext {
	return &hcl.EvalContext{
		Functions: map[string]function.Function{
			"uppercase": makeListTransformFunction(stdlib.UpperFunc),
			"replace":   makeListTransformFunction(stdlib.ReplaceFunc),
		},
		Variables: make(map[string]cty.Value),
	}
}

// makeListTransformFunction wraps a string transform function to work on lists of strings
func makeListTransformFunction(strFunc function.Function) function.Function {
	return function.New(&function.Spec{
		VarParam: &function.Parameter{
			Name: "args",
			Type: cty.DynamicPseudoType,
		},
		Type: func(args []cty.Value) (cty.Type, error) {
			return cty.List(cty.String), nil
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			if len(args) == 0 {
				return cty.NilVal, fmt.Errorf("expected at least one argument")
			}
			if !args[0].Type().IsListType() {
				return cty.NilVal, fmt.Errorf("first argument must be a list")
			}

			// Handle empty list
			if args[0].LengthInt() == 0 {
				return cty.ListValEmpty(cty.String), nil
			}

			var results []cty.Value
			for it := args[0].ElementIterator(); it.Next(); {
				_, val := it.Element()

				// Build args for the string function: [element, ...other args]
				elementArgs := make([]cty.Value, len(args))
				elementArgs[0] = val
				copy(elementArgs[1:], args[1:])

				result, err := strFunc.Call(elementArgs)
				if err != nil {
					return cty.NilVal, err
				}
				results = append(results, result)
			}
			return cty.ListVal(results), nil
		},
	})
}
