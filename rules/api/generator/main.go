// +build generators

package main

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	utils "github.com/terraform-linters/tflint-ruleset-aws/rules/generator-utils"
)

type definition struct {
	Rules []rule `hcl:"rule,block"`
}

type rule struct {
	Name         string `hcl:"name,label"`
	Resource     string `hcl:"resource"`
	Attribute    string `hcl:"attribute"`
	SourceAction string `hcl:"source_action"`
	Template     string `hcl:"template"`
}

type ruleMeta struct {
	RuleName      string
	RuleNameCC    string
	ResourceType  string
	AttributeName string
	DataType      string
	ActionName    string
	Template      string
}

type providerMeta struct {
	RuleNameCCList []string
}

var awsProvider = utils.LoadProviderSchema("../../tools/provider-schema/schema.json")

func main() {
	files, err := filepath.Glob("./definitions/*.hcl")
	if err != nil {
		panic(err)
	}

	providerMeta := &providerMeta{}
	for _, file := range files {
		parser := hclparse.NewParser()
		f, diags := parser.ParseHCLFile(file)
		if diags.HasErrors() {
			panic(diags)
		}

		var def definition
		diags = gohcl.DecodeBody(f.Body, nil, &def)
		if diags.HasErrors() {
			panic(diags)
		}

		for _, rule := range def.Rules {
			meta := &ruleMeta{
				RuleName:      rule.Name,
				RuleNameCC:    utils.ToCamel(rule.Name),
				ResourceType:  rule.Resource,
				AttributeName: rule.Attribute,
				DataType:      dataType(rule.Resource, rule.Attribute),
				ActionName:    rule.SourceAction,
				Template:      rule.Template,
			}

			utils.GenerateFile(
				fmt.Sprintf("%s.go", rule.Name),
				"rule.go.tmpl",
				meta,
			)

			providerMeta.RuleNameCCList = append(providerMeta.RuleNameCCList, meta.RuleNameCC)
		}
	}

	sort.Strings(providerMeta.RuleNameCCList)
	utils.GenerateFile(
		"provider.go",
		"provider.go.tmpl",
		providerMeta,
	)
}

func dataType(resource, attribute string) string {
	resourceSchema, ok := awsProvider.ResourceSchemas[resource]
	if !ok {
		panic(fmt.Sprintf("resource `%s` not found in the Terraform schema", resource))
	}
	attrSchema, ok := resourceSchema.Block.Attributes[attribute]
	if !ok {
		panic(fmt.Sprintf("`%s.%s` not found in the Terraform schema", resource, attribute))
	}

	switch ty := attrSchema.Type.(type) {
	case string:
		if ty != "string" {
			panic(fmt.Errorf("Unexpected data type: %#v", attrSchema.Type))
		}
		return "string"
	case []interface{}:
		if len(ty) != 2 || !(ty[0] == "set" && ty[1] == "string") {
			panic(fmt.Errorf("Unexpected data type: %#v", attrSchema.Type))
		}
		return "list"
	default:
		panic(fmt.Errorf("Unexpected data type: %#v", attrSchema.Type))
	}
}
