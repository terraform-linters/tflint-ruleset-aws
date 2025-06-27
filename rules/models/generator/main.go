//go:build generators

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

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
		shapes := api["shapes"].(map[string]interface{})

		for _, mapping := range mappingFile.Mappings {
			for attribute, value := range mapping.Attrs {
				fmt.Printf("Checking `%s.%s`\n", mapping.Resource, attribute)
				shapeName := value.Expr.Variables()[0].RootName()
				if shapeName == "any" {
					continue
				}
				model := shapes[shapeName].(map[string]interface{})
				schema, err := fetchSchema(mapping.Resource, attribute, model, awsProvider)
				if err != nil {
					fmt.Printf("Skipping `%s.%s`: %v\n", mapping.Resource, attribute, err)
					continue
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
