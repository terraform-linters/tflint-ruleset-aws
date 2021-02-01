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

// RuleNames is a list of rule names provided by the plugin.
func (r *RuleSet) RuleNames() []string {
	names := []string{}
	for _, rule := range r.Rules {
		names = append(names, rule.Name())
	}
	for _, rule := range r.APIRules {
		names = append(names, rule.Name())
	}
	return names
}

// ApplyConfig reflects the plugin configuration to the ruleset.
func (r *RuleSet) ApplyConfig(config *tflint.Config) error {
	r.ApplyCommonConfig(config)

	// Apply "plugin" block config
	cfg := Config{}
	diags := gohcl.DecodeBody(config.Body, nil, &cfg)
	if diags.HasErrors() {
		return diags
	}
	r.config = &cfg

	// Apply config for API rules
	for _, rule := range r.APIRules {
		enabled := rule.Enabled()
		if cfg := config.Rules[rule.Name()]; cfg != nil {
			enabled = cfg.Enabled
		} else if config.DisabledByDefault {
			enabled = false
		}

		if cfg.DeepCheck && enabled {
			r.EnabledRules = append(r.EnabledRules, rule)
		}
	}

	return nil
}

// Check runs inspections for each rule with the custom AWS runner.
func (r *RuleSet) Check(rr tflint.Runner) error {
	runner, err := NewRunner(rr, r.config)
	if err != nil {
		return err
	}

	for _, rule := range r.EnabledRules {
		if err := rule.Check(runner); err != nil {
			return fmt.Errorf("Failed to check `%s` rule: %s", rule.Name(), err)
		}
	}
	return nil
}
