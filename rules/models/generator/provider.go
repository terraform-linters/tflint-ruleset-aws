//go:build generators

package main

import utils "github.com/terraform-linters/tflint-ruleset-aws/rules/generator-utils"

type providerMeta struct {
	RuleNameCCList []string
}

func generateProviderFile(ruleNames []string) {
	meta := &providerMeta{}

	for _, ruleName := range ruleNames {
		meta.RuleNameCCList = append(meta.RuleNameCCList, utils.ToCamel(ruleName))
	}

	utils.GenerateFile("provider.go", "provider.go.tmpl", meta)
}
