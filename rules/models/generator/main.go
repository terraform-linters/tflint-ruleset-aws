//go:build generators

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	tfjson "github.com/hashicorp/terraform-json"
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

func main() {
	files, err := filepath.Glob("./mappings/*.hcl")
	if err != nil {
		panic(err)
	}

	mappingFiles := []mappingFile{}
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

	generatedRules := []string{}
	for _, mappingFile := range mappingFiles {
		raw, err := ioutil.ReadFile(mappingFile.Import)
		if err != nil {
			panic(err)
		}

		var api map[string]interface{}
		err = json.Unmarshal(raw, &api)
		if err != nil {
			panic(err)
		}

		// Detect format and extract shapes accordingly
		var shapes map[string]interface{}
		if smithyVersion, ok := api["smithy"]; ok && smithyVersion != nil {
			// Smithy format (api-models-aws)
			shapes = api["shapes"].(map[string]interface{})
		} else {
			// Legacy Ruby SDK format
			shapes = api["shapes"].(map[string]interface{})
		}

		for _, mapping := range mappingFile.Mappings {
			for attribute, value := range mapping.Attrs {
				fmt.Printf("Checking `%s.%s`\n", mapping.Resource, attribute)
				shapeName := value.Expr.Variables()[0].RootName()
				if shapeName == "any" {
					continue
				}
				
				// Handle both naming conventions
				model := findShape(shapes, shapeName, mappingFile.Import)
				if model == nil {
					fmt.Printf("Shape `%s` not found, skipping\n", shapeName)
					continue
				}
				
				schema := fetchSchema(mapping.Resource, attribute, model, awsProvider)
				if validMapping(model) {
					fmt.Printf("Generating rule for `%s.%s`\n", mapping.Resource, attribute)
					generateRuleFile(mapping.Resource, attribute, model, schema)
					for _, test := range mappingFile.Tests {
						if mapping.Resource == test.Resource && attribute == test.Attribute {
							generateRuleTestFile(mapping.Resource, attribute, model, test)
						}
					}
					generatedRules = append(generatedRules, makeRuleName(mapping.Resource, attribute))
				}
			}
		}
	}

	sort.Strings(generatedRules)
	generateProviderFile(generatedRules)
	generateDocFile(generatedRules)
}

func fetchSchema(resource, attribute string, model map[string]interface{}, provider *tfjson.ProviderSchema) *tfjson.SchemaAttribute {
	resourceSchema, ok := provider.ResourceSchemas[resource]
	if !ok {
		panic(fmt.Sprintf("resource `%s` not found in the Terraform schema", resource))
	}
	attrSchema, ok := resourceSchema.Block.Attributes[attribute]
	if !ok {
		if _, ok := resourceSchema.Block.NestedBlocks[attribute]; !ok {
			panic(fmt.Sprintf("`%s.%s` not found in the Terraform schema", resource, attribute))
		}
	}

	switch model["type"].(string) {
	case "string":
		ty := attrSchema.AttributeType.FriendlyName()
		if ty != "string" && ty != "number" {
			panic(fmt.Sprintf("`%s.%s` is expected as string, but not (%s)", resource, attribute, ty))
		}
	default:
		// noop
	}

	return attrSchema
}

// findShape locates a shape in either Ruby SDK or Smithy format
func findShape(shapes map[string]interface{}, shapeName, importPath string) map[string]interface{} {
	// First try direct lookup (Ruby SDK format)
	if shape, ok := shapes[shapeName]; ok {
		return shape.(map[string]interface{})
	}
	
	// If not found and this looks like a Smithy file, try with service namespace
	if isSmithyFile(importPath) {
		serviceNamespace := extractServiceNamespace(shapes)
		if serviceNamespace != "" {
			qualifiedName := fmt.Sprintf("%s#%s", serviceNamespace, shapeName)
			if shape, ok := shapes[qualifiedName]; ok {
				// Convert Smithy shape to Ruby SDK-like format for compatibility
				return convertSmithyShape(shape.(map[string]interface{}))
			}
		}
	}
	
	return nil
}

// isSmithyFile checks if the import path is for a Smithy file
func isSmithyFile(importPath string) bool {
	return strings.Contains(importPath, "api-models-aws") && 
		   !strings.Contains(importPath, "aws-sdk-ruby")
}

// extractServiceNamespace gets the service namespace from the Smithy shapes
func extractServiceNamespace(shapes map[string]interface{}) string {
	// Look for the service definition shape to find the namespace
	for shapeName, shape := range shapes {
		if shapeMap, ok := shape.(map[string]interface{}); ok {
			if shapeType, ok := shapeMap["type"].(string); ok && shapeType == "service" {
				// Extract namespace from the service shape name (e.g., "com.amazonaws.acmpca#ACMPrivateCA")
				if parts := strings.Split(shapeName, "#"); len(parts) == 2 {
					return parts[0]
				}
			}
		}
	}
	return ""
}

// convertSmithyShape converts a Smithy shape format to Ruby SDK-like format
func convertSmithyShape(smithyShape map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Copy type
	if shapeType, ok := smithyShape["type"]; ok {
		result["type"] = shapeType
	}
	
	// Convert traits to Ruby SDK format
	if traits, ok := smithyShape["traits"].(map[string]interface{}); ok {
		// Handle length constraints
		if lengthTrait, ok := traits["smithy.api#length"].(map[string]interface{}); ok {
			if min, ok := lengthTrait["min"]; ok {
				result["min"] = min
			}
			if max, ok := lengthTrait["max"]; ok {
				result["max"] = max
			}
		}
		
		// Handle pattern
		if pattern, ok := traits["smithy.api#pattern"].(string); ok {
			result["pattern"] = pattern
		}
		
		// Handle enum (stored differently in Smithy)
		if enumTrait, ok := traits["smithy.api#enum"]; ok {
			// Handle enum values from the trait
			if enumList, ok := enumTrait.([]interface{}); ok {
				enumValues := make([]string, 0, len(enumList))
				for _, enumItem := range enumList {
					if enumMap, ok := enumItem.(map[string]interface{}); ok {
						if value, ok := enumMap["value"].(string); ok {
							enumValues = append(enumValues, value)
						}
					}
				}
				result["enum"] = enumValues
			}
		}
	}
	
	// Handle enum type shapes (where type is "enum")
	if shapeType, ok := smithyShape["type"].(string); ok && shapeType == "enum" {
		if members, ok := smithyShape["members"].(map[string]interface{}); ok {
			enumValues := make([]string, 0, len(members))
			for memberName, memberData := range members {
				// Extract enumValue from traits, fallback to memberName if not present
				enumValue := memberName
				if memberMap, ok := memberData.(map[string]interface{}); ok {
					if traits, ok := memberMap["traits"].(map[string]interface{}); ok {
						if enumValueTrait, ok := traits["smithy.api#enumValue"].(string); ok {
							enumValue = enumValueTrait
						}
					}
				}
				enumValues = append(enumValues, enumValue)
			}
			result["enum"] = enumValues
			result["type"] = "string" // Convert to string type for compatibility
		}
	}
	
	return result
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
