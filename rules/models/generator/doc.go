//go:build generators

package main

import "github.com/terraform-linters/tflint-ruleset-aws/rules/genutils"

type docMeta struct {
	RuleNames []string
}

func generateDocFile(ruleNames []string) {
	meta := &docMeta{RuleNames: ruleNames}

	genutils.GenerateFile("../../docs/rules/README.md", "../../docs/rules/README.md.tmpl", meta)
}
