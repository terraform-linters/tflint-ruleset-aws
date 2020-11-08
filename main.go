package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/aws"
	"github.com/terraform-linters/tflint-ruleset-aws/rules"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/api"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &aws.RuleSet{
			BuiltinRuleSet: tflint.BuiltinRuleSet{
				Name:    "aws",
				Version: "0.1.0",
				Rules: []tflint.Rule{
					rules.NewAwsInstanceExampleTypeRule(),
					rules.NewAwsCustomRunnerRule(),
				},
			},
			APIRules: api.Rules,
		},
	})
}
