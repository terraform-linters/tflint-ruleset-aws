package aws

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Runner is a wrapper of RPC client for inserting custom actions for AWS provider.
type Runner struct {
	tflint.Runner
	PluginConfig *Config
	AwsClient    *Client
}

// NewRunner returns a custom AWS runner.
func NewRunner(runner tflint.Runner, config *Config) (*Runner, error) {
	var client *Client
	if config.DeepCheck {
		var err error
		client, err = NewClient(config.toCredentials())
		if err != nil {
			return nil, err
		}
	}

	return &Runner{
		Runner:       runner,
		PluginConfig: config,
		AwsClient:    client,
	}, nil
}

// EachStringSliceExprs iterates an evaluated value and the corresponding expression
// If the given expression is a static list, get an expression for each value
// If not, the given expression is used as it is
func (r *Runner) EachStringSliceExprs(expr hcl.Expression, proc func(val string, expr hcl.Expression)) error {
	var vals []string
	err := r.EvaluateExpr(expr, &vals)

	exprs, diags := hcl.ExprList(expr)
	if diags.HasErrors() {
		log.Printf("[DEBUG] Expr is not static list: %s", diags)
		for range vals {
			exprs = append(exprs, expr)
		}
	}

	return r.EnsureNoError(err, func() error {
		for idx, val := range vals {
			proc(val, exprs[idx])
		}
		return nil
	})
}

// CustomCall is ...
func (r *Runner) CustomCall() string {
	return fmt.Sprintf("config=%#v", r.PluginConfig)
}
