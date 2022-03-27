package aws

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// RuleSet is the custom ruleset for the AWS provider plugin.
type RuleSet struct {
	tflint.BuiltinRuleSet
	config *Config
}

func (r *RuleSet) ConfigSchema() *hclext.BodySchema {
	r.config = &Config{}
	return hclext.ImpliedBodySchema(r.config)
}

// ApplyConfig reflects the plugin configuration to the ruleset.
func (r *RuleSet) ApplyConfig(body *hclext.BodyContent) error {
	diags := hclext.DecodeBody(body, nil, r.config)
	if diags.HasErrors() {
		return diags
	}

	if r.config.DeepCheck {
		return nil
	}

	// Disable deep checking rules
	enabledRules := []tflint.Rule{}
	for _, rule := range r.EnabledRules {
		meta := rule.Metadata()
		// Deep checking rules must have metadata like `map[string]bool{"deep": true}``
		if meta == nil {
			enabledRules = append(enabledRules, rule)
		}
	}
	r.EnabledRules = enabledRules

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
