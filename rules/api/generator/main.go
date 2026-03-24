//go:build generators

package main

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/genutils"
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

var awsProvider = genutils.LoadProviderSchema("../../tools/provider-schema/schema.json")

func main() {
	files, err := filepath.Glob("./definitions/*.hcl")
	if err != nil {
		panic(err)
	}

	providerMeta := &providerMeta{}
	var generatedFiles []string
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
				RuleNameCC:    genutils.ToCamel(rule.Name),
				ResourceType:  rule.Resource,
				AttributeName: rule.Attribute,
				DataType:      dataType(rule.Resource, rule.Attribute),
				ActionName:    rule.SourceAction,
				Template:      rule.Template,
			}

			genutils.GenerateFile(
				fmt.Sprintf("%s.go", rule.Name),
				"rule.go.tmpl",
				meta,
			)
			generatedFiles = append(generatedFiles, fmt.Sprintf("%s.go", rule.Name))

			providerMeta.RuleNameCCList = append(providerMeta.RuleNameCCList, meta.RuleNameCC)
		}
	}

	sort.Strings(providerMeta.RuleNameCCList)
	genutils.GenerateFile(
		"provider.go",
		"provider.go.tmpl",
		providerMeta,
	)
	generatedFiles = append(generatedFiles, "provider.go")
	genutils.CleanDir(".", generatedFiles)
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

	if ty := attrSchema.AttributeType.FriendlyName(); ty == "string" {
		return "string"
	} else if attrSchema.AttributeType.IsSetType() || attrSchema.AttributeType.IsListType() {
		return "list"
	}
	panic(fmt.Errorf("Unexpected data type: %#v", attrSchema.AttributeType))
}
