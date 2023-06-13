package aws

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Runner is a wrapper of RPC client for inserting custom actions for AWS provider.
type Runner struct {
	tflint.Runner
	PluginConfig *Config
	AwsClients   map[string]*Client
}

// NewRunner returns a custom AWS runner.
func NewRunner(runner tflint.Runner, config *Config) (*Runner, error) {
	clients := map[string]*Client{
		"aws": nil,
	}
	if config.DeepCheck {
		credentials, err := GetCredentialsFromProvider(runner)
		if err != nil {
			return nil, err
		}
		if _, ok := credentials["aws"]; !ok {
			credentials["aws"] = Credentials{}
		}
		for k, cred := range credentials {
			client, err := NewClient(cred.Merge(config.toCredentials()))
			if err != nil {
				return nil, err
			}
			clients[k] = client
		}
	}

	return &Runner{
		Runner:       runner,
		PluginConfig: config,
		AwsClients:   clients,
	}, nil
}

func (r *Runner) AwsClient(attributes hclext.Attributes) (*Client, error) {
	provider := "aws"
	if attr, exists := attributes["provider"]; exists {
		providerConfigRef, diags := DecodeProviderConfigRef(attr.Expr, "provider")
		if diags.HasErrors() {
			logger.Error("parse resource provider attribute: %s", diags)
			return nil, diags
		}
		if providerConfigRef.Alias != "" {
			provider = providerConfigRef.Alias
		} else {
			provider = providerConfigRef.Name
		}
	}

	awsClient, ok := r.AwsClients[provider]
	if !ok {
		return nil, fmt.Errorf("aws provider %s isn't found", provider)
	}

	return awsClient, nil
}

// EachStringSliceExprs iterates an evaluated value and the corresponding expression
// If the given expression is a static list, get an expression for each value
// If not, the given expression is used as it is
func (r *Runner) EachStringSliceExprs(expr hcl.Expression, proc func(val string, expr hcl.Expression)) error {
	var vals []string
	err := r.EvaluateExpr(expr, func(v []string) error {
		vals = v
		return nil
	}, nil)
	if err != nil {
		return err
	}

	exprs, diags := hcl.ExprList(expr)
	if diags.HasErrors() {
		logger.Debug("Expr is not static list: %s", diags)
		for range vals {
			exprs = append(exprs, expr)
		}
	}

	for idx, val := range vals {
		proc(val, exprs[idx])
	}
	return nil
}
