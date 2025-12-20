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

		evalCtx := buildEvalContext(shapes)

		for _, mapping := range mappingFile.Mappings {
			for attribute, value := range mapping.Attrs {
				fmt.Printf("Checking `%s.%s`\n", mapping.Resource, attribute)

				// Evaluate the expression to see what type of mapping it is
				result, diags := value.Expr.Value(evalCtx)
				if diags.HasErrors() {
					fmt.Fprintf(os.Stderr, "Error evaluating `%s.%s`: %v\n", mapping.Resource, attribute, diags)
					os.Exit(1)
				}

				// Check if result is a listmap (fully resolved constraints)
				if result.Type().Equals(listmapResultType) {
					schema, err := fetchSchemaForMap(mapping.Resource, attribute, awsProvider)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error processing `%s.%s`: %v\n", mapping.Resource, attribute, err)
						os.Exit(1)
					}

					fmt.Printf("Generating listmap rule for `%s.%s`\n", mapping.Resource, attribute)
					if generateListmapRuleFile(mapping.Resource, attribute, result, schema) {
						generatedRules = append(generatedRules, makeRuleName(mapping.Resource, attribute))
					}
					continue
				}

				// Result should be a shape object for simple mappings
				if !result.Type().Equals(shapeType) {
					fmt.Fprintf(os.Stderr, "Error: `%s.%s` evaluated to unexpected type %s\n",
						mapping.Resource, attribute, result.Type().FriendlyName())
					os.Exit(1)
				}

				shapeName := result.GetAttr("name").AsString()
				if shapeName == "any" {
					continue
				}

				// Convert shape object back to model format for existing code paths
				model := shapeToModel(result)
				if model == nil {
					fmt.Printf("Shape %q has no constraints, skipping\n", shapeName)
					continue
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
					// Try traversing to find map constraints (only works for Smithy map types)
					keyModel, valueModel := traverseToMapConstraints(shapes, shapeName)
					if keyModel != nil && valueModel != nil {
						if validMapping(keyModel) || validMapping(valueModel) {
							fmt.Printf("Generating map rule for `%s.%s`\n", mapping.Resource, attribute)
							if generateMapRuleFile(mapping.Resource, attribute, nil, keyModel, valueModel, schema) {
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

func fetchSchemaForMap(resource, attribute string, provider *tfjson.ProviderSchema) (*tfjson.SchemaAttribute, error) {
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
	return attrSchema, nil
}

func validMapping(model map[string]interface{}) bool {
	modelType, ok := model["type"].(string)
	if !ok {
		return false
	}
	switch modelType {
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
			if shapeMap, ok := shape.(map[string]interface{}); ok {
				return shapeMap
			}
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

// getTarget extracts the target shape name from a Smithy member reference.
func getTarget(ref map[string]interface{}) string {
	if ref == nil {
		return ""
	}
	target, _ := ref["target"].(string)
	return extractShapeName(target)
}

// traverseToMapConstraints traverses Smithy shapes to find key/value constraints.
// Only Smithy map types have explicit key/value - structures require configuration.
func traverseToMapConstraints(shapes map[string]interface{}, shapeName string) (keyModel, valueModel map[string]interface{}) {
	raw := findRawShape(shapes, shapeName)
	if raw == nil {
		return nil, nil
	}

	shapeType, _ := raw["type"].(string)

	switch shapeType {
	case "map":
		// Smithy map types have explicit "key" and "value" fields per the spec
		keyRef, _ := raw["key"].(map[string]interface{})
		valueRef, _ := raw["value"].(map[string]interface{})
		return resolveStringPair(shapes, getTarget(keyRef), getTarget(valueRef))

	case "list":
		// Recurse into list member
		memberRef, _ := raw["member"].(map[string]interface{})
		if target := getTarget(memberRef); target != "" {
			return traverseToMapConstraints(shapes, target)
		}
	}

	// Structures and other types: no automatic key/value inference
	return nil, nil
}

// resolveStringPair resolves two shape names to their models,
// returning nil if either is missing or not a string type.
func resolveStringPair(shapes map[string]interface{}, firstName, secondName string) (map[string]interface{}, map[string]interface{}) {
	if firstName == "" || secondName == "" {
		return nil, nil
	}
	first := findShape(shapes, firstName)
	second := findShape(shapes, secondName)
	if first == nil || second == nil {
		return nil, nil
	}
	if first["type"] != "string" || second["type"] != "string" {
		return nil, nil
	}
	return first, second
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
// The generator uses HCL expression evaluation with shape objects as variables.
// Each Smithy shape becomes a rich cty object containing all constraint info.
//
// Shape Object Fields:
//   - name: string - The shape name
//   - type: string - The shape type (e.g., "string", "list")
//   - pattern: string - Regex pattern constraint
//   - min: number - Minimum length constraint
//   - max: number - Maximum length constraint
//   - enum: list(string) - Allowed enum values
//
// Available Functions:
//   - uppercase(shape) - Transforms enum values to uppercase
//   - replace(shape, old, new) - Replaces substrings in enum values
//   - listmap(list, key, value) - Combines three shapes for map validation
//
// Usage Examples in mappings/*.hcl:
//   compression_type = uppercase(CompressionTypeValue)
//     -> Transforms shape's enum values to uppercase
//
//   encryption_mode = uppercase(replace(EncryptionModeValue, "-", "_"))
//     -> Chains transforms on shape's enum values
//
//   tags = listmap(TagList, TagKey, TagValue)
//     -> Combines constraints from three shapes for tag validation

// shapeType is the cty object type for shape variables.
// Shapes are rich objects containing all their constraint information.
var shapeType = cty.Object(map[string]cty.Type{
	"name":    cty.String,
	"type":    cty.String,
	"pattern": cty.String,
	"min":     cty.Number,
	"max":     cty.Number,
	"enum":    cty.List(cty.String),
})

// listmapResultType is the cty object type returned by the listmap function.
// Contains fully resolved constraints from list, key, and value shapes.
var listmapResultType = cty.Object(map[string]cty.Type{
	"type":         cty.String, // "listmap" marker
	"itemsMin":     cty.Number,
	"itemsMax":     cty.Number,
	"keyPattern":   cty.String,
	"keyMin":       cty.Number,
	"keyMax":       cty.Number,
	"keyEnum":      cty.List(cty.String),
	"valuePattern": cty.String,
	"valueMin":     cty.Number,
	"valueMax":     cty.Number,
	"valueEnum":    cty.List(cty.String),
})

// buildEvalContext creates an HCL evaluation context with transform functions
// and shape variables. Each shape becomes a rich object with all its constraints.
func buildEvalContext(shapes map[string]interface{}) *hcl.EvalContext {
	variables := make(map[string]cty.Value)

	// Special value to skip validation
	variables["any"] = makeShapeValue("any", nil)

	// Populate shapes as rich objects
	for qualifiedName := range shapes {
		name := extractShapeName(qualifiedName)
		if name != "" {
			model := findShape(shapes, name)
			variables[name] = makeShapeValue(name, model)
		}
	}

	return &hcl.EvalContext{
		Functions: map[string]function.Function{
			"uppercase": makeShapeTransformFunction(stdlib.UpperFunc),
			"replace":   makeShapeTransformFunction(stdlib.ReplaceFunc),
			"listmap":   makeListmapFunction(),
		},
		Variables: variables,
	}
}

// makeShapeValue creates a cty shape object from a model
func makeShapeValue(name string, model map[string]interface{}) cty.Value {
	typeStr := ""
	pattern := ""
	min := int64(0)
	max := int64(0)
	var enumList cty.Value

	if model != nil {
		if t, ok := model["type"].(string); ok {
			typeStr = t
		}
		pattern = fetchString(model, "pattern")
		min = int64(fetchNumber(model, "min"))
		max = int64(fetchNumber(model, "max"))

		if enums := fetchStrings(model, "enum"); len(enums) > 0 {
			vals := make([]cty.Value, len(enums))
			for i, e := range enums {
				vals[i] = cty.StringVal(e)
			}
			enumList = cty.ListVal(vals)
		} else {
			enumList = cty.ListValEmpty(cty.String)
		}
	} else {
		enumList = cty.ListValEmpty(cty.String)
	}

	return cty.ObjectVal(map[string]cty.Value{
		"name":    cty.StringVal(name),
		"type":    cty.StringVal(typeStr),
		"pattern": cty.StringVal(pattern),
		"min":     cty.NumberIntVal(min),
		"max":     cty.NumberIntVal(max),
		"enum":    enumList,
	})
}

// shapeToModel converts a cty shape object back to a model map for existing code paths
func shapeToModel(shape cty.Value) map[string]interface{} {
	model := make(map[string]interface{})

	model["type"] = shape.GetAttr("type").AsString()

	if pattern := shape.GetAttr("pattern").AsString(); pattern != "" {
		model["pattern"] = pattern
	}

	minVal := shape.GetAttr("min")
	if bf := minVal.AsBigFloat(); bf.Sign() > 0 {
		f, _ := bf.Float64()
		model["min"] = f
	}

	maxVal := shape.GetAttr("max")
	if bf := maxVal.AsBigFloat(); bf.Sign() > 0 {
		f, _ := bf.Float64()
		model["max"] = f
	}

	enumVal := shape.GetAttr("enum")
	if enumVal.LengthInt() > 0 {
		enums := make([]string, 0, enumVal.LengthInt())
		for it := enumVal.ElementIterator(); it.Next(); {
			_, v := it.Element()
			enums = append(enums, v.AsString())
		}
		model["enum"] = enums
	}

	return model
}

// makeShapeTransformFunction wraps a string transform to work on shape objects.
// It transforms the enum values and returns a new shape with transformed enums.
func makeShapeTransformFunction(strFunc function.Function) function.Function {
	return function.New(&function.Spec{
		VarParam: &function.Parameter{
			Name: "args",
			Type: cty.DynamicPseudoType,
		},
		Type: func(args []cty.Value) (cty.Type, error) {
			if len(args) == 0 {
				return cty.NilType, fmt.Errorf("expected at least one argument")
			}
			// If first arg is a shape, return a shape
			if args[0].Type().Equals(shapeType) {
				return shapeType, nil
			}
			// Otherwise return list of strings (for backwards compatibility)
			return cty.List(cty.String), nil
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			if len(args) == 0 {
				return cty.NilVal, fmt.Errorf("expected at least one argument")
			}

			// If first arg is a shape object, transform its enum values
			if args[0].Type().Equals(shapeType) {
				shape := args[0]
				enumVal := shape.GetAttr("enum")

				// If no enum values, return shape unchanged
				if enumVal.LengthInt() == 0 {
					return shape, nil
				}

				// Transform each enum value
				var results []cty.Value
				for it := enumVal.ElementIterator(); it.Next(); {
					_, val := it.Element()
					transformArgs := make([]cty.Value, len(args))
					transformArgs[0] = val
					copy(transformArgs[1:], args[1:])

					result, err := strFunc.Call(transformArgs)
					if err != nil {
						return cty.NilVal, err
					}
					results = append(results, result)
				}

				// Return new shape with transformed enums
				return cty.ObjectVal(map[string]cty.Value{
					"name":    shape.GetAttr("name"),
					"type":    shape.GetAttr("type"),
					"pattern": shape.GetAttr("pattern"),
					"min":     shape.GetAttr("min"),
					"max":     shape.GetAttr("max"),
					"enum":    cty.ListVal(results),
				}), nil
			}

			// Legacy path: first arg is a list of strings
			if !args[0].Type().IsListType() {
				return cty.NilVal, fmt.Errorf("first argument must be a shape or list")
			}

			if args[0].LengthInt() == 0 {
				return cty.ListValEmpty(cty.String), nil
			}

			var results []cty.Value
			for it := args[0].ElementIterator(); it.Next(); {
				_, val := it.Element()
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

// makeListmapFunction creates the listmap(list, key, value) function.
// It takes three shape objects and combines their constraints into a resolved result.
func makeListmapFunction() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{Name: "list", Type: shapeType},
			{Name: "key", Type: shapeType},
			{Name: "value", Type: shapeType},
		},
		Type: function.StaticReturnType(listmapResultType),
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			listShape := args[0]
			keyShape := args[1]
			valueShape := args[2]

			return cty.ObjectVal(map[string]cty.Value{
				"type":         cty.StringVal("listmap"),
				"itemsMin":     listShape.GetAttr("min"),
				"itemsMax":     listShape.GetAttr("max"),
				"keyPattern":   keyShape.GetAttr("pattern"),
				"keyMin":       keyShape.GetAttr("min"),
				"keyMax":       keyShape.GetAttr("max"),
				"keyEnum":      keyShape.GetAttr("enum"),
				"valuePattern": valueShape.GetAttr("pattern"),
				"valueMin":     valueShape.GetAttr("min"),
				"valueMax":     valueShape.GetAttr("max"),
				"valueEnum":    valueShape.GetAttr("enum"),
			}), nil
		},
	})
}

