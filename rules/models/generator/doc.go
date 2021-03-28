// +build generators

package main

import utils "github.com/terraform-linters/tflint-ruleset-aws/rules/generator-utils"

type docMeta struct {
	RuleNames []string
}

func generateDocFile(ruleNames []string) {
	meta := &docMeta{RuleNames: ruleNames}

	utils.GenerateFile("../../docs/rules/README.md", "../../docs/rules/README.md.tmpl", meta)
}
