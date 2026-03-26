//go:build generators

package main

import "github.com/terraform-linters/tflint-ruleset-aws/rules/genutils"

type providerMeta struct {
	RuleNameCCList []string
}

func generateProviderFile(ruleNames []string) {
	meta := &providerMeta{}

	for _, ruleName := range ruleNames {
		meta.RuleNameCCList = append(meta.RuleNameCCList, genutils.ToCamel(ruleName))
	}

	genutils.GenerateFile("provider.go", "provider.go.tmpl", meta)
}
