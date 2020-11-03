package aws

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Runner is a wrapper of RPC client for inserting custom actions for AWS provider.
type Runner struct {
	tflint.Runner
	PluginConfig *Config
}

// CustomCall is ...
func (r *Runner) CustomCall() string {
	return fmt.Sprintf("config=%s", r.PluginConfig)
}
