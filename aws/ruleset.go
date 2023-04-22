package aws

import (
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

// NewRunner injects a custom AWS runner
func (r *RuleSet) NewRunner(runner tflint.Runner) (tflint.Runner, error) {
	return NewRunner(runner, r.config)
}
