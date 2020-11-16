package aws

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// RuleSet is the custom ruleset for the AWS provider plugin.
type RuleSet struct {
	tflint.BuiltinRuleSet
	APIRules []tflint.Rule
	config   *Config
}

// ApplyConfig reflects the plugin configuration to the ruleset.
func (r *RuleSet) ApplyConfig(config *tflint.Config) error {
	r.ApplyCommonConfig(config)

	cfg := Config{}
	diags := gohcl.DecodeBody(config.Body, nil, &cfg)
	if diags.HasErrors() {
		return diags
	}
	r.config = &cfg
	return nil
}

// Check runs inspections for each rule with the custom AWS runner.
func (r *RuleSet) Check(rr tflint.Runner) error {
	runner, err := NewRunner(rr, r.config)
	if err != nil {
		return err
	}

	for _, rule := range r.Rules {
		if err := rule.Check(runner); err != nil {
			return fmt.Errorf("Failed to check `%s` rule: %s", rule.Name(), err)
		}
	}
	if !r.config.DeepCheck {
		return nil
	}

	for _, rule := range r.APIRules {
		if err := rule.Check(runner); err != nil {
			return fmt.Errorf("Failed to check `%s` rule: %s", rule.Name(), err)
		}
	}
	return nil
}
