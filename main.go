package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/aws"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/terraform-linters/tflint-ruleset-aws/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &aws.RuleSet{
			BuiltinRuleSet: tflint.BuiltinRuleSet{
				Name:    "aws",
				Version: project.Version,
				Rules:   rules.Rules,
			},
		},
	})
}
