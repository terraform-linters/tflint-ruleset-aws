package aws

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsProviderBlockSchema is a schema of `aws` provider block
var AwsProviderBlockSchema = &hclext.BodySchema{
	Attributes: []hclext.AttributeSchema{
		{Name: "access_key"},
		{Name: "secret_key"},
		{Name: "profile"},
		{Name: "shared_credentials_file"},
		{Name: "region"},
		{Name: "alias"},
	},
	Blocks: []hclext.BlockSchema{
		{
			Type: "assume_role",
			Body: AwsProviderAssumeRoleBlockShema,
		},
	},
}

// AwsProviderAssumeRoleBlockShema is a schema of `assume_role` block
var AwsProviderAssumeRoleBlockShema = &hclext.BodySchema{
	Attributes: []hclext.AttributeSchema{
		{Name: "role_arn", Required: true},
		{Name: "session_name"},
		{Name: "external_id"},
		{Name: "policy"},
	},
}

// GetCredentialsFromProvider retrieves credentials from the "provider" block in the Terraform configuration
func GetCredentialsFromProvider(runner tflint.Runner) (map[string]Credentials, error) {
	providers, err := runner.GetModuleContent(
		&hclext.BodySchema{
			Blocks: []hclext.BlockSchema{
				{
					Type:       "provider",
					LabelNames: []string{"name"},
					Body:       AwsProviderBlockSchema,
				},
			},
		},
		&tflint.GetModuleContentOption{ModuleCtx: tflint.RootModuleCtxType},
	)
	if err != nil {
		return nil, err
	}

	m := map[string]Credentials{}

	for _, provider := range providers.Blocks {
		if provider.Labels[0] != "aws" {
			continue
		}

		creds := Credentials{}

		opts := &tflint.EvaluateExprOption{ModuleCtx: tflint.RootModuleCtxType}

		if attr, exists := provider.Body.Attributes["access_key"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, func(accessKey string) error {
				creds.AccessKey = accessKey
				return nil
			}, opts); err != nil {
				return nil, err
			}
		}

		if attr, exists := provider.Body.Attributes["secret_key"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, func(secretKey string) error {
				creds.SecretKey = secretKey
				return nil
			}, opts); err != nil {
				return nil, err
			}
		}

		if attr, exists := provider.Body.Attributes["profile"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, func(profile string) error {
				creds.Profile = profile
				return nil
			}, opts); err != nil {
				return nil, err
			}
		}

		if attr, exists := provider.Body.Attributes["shared_credentials_file"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, func(credsFile string) error {
				creds.CredsFile = credsFile
				return nil
			}, opts); err != nil {
				return nil, err
			}
		}

		if attr, exists := provider.Body.Attributes["region"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, func(region string) error {
				creds.Region = region
				return nil
			}, opts); err != nil {
				return nil, err
			}
		}

		for _, assumeRole := range provider.Body.Blocks {
			if attr, exists := assumeRole.Body.Attributes["role_arn"]; exists {
				if err := runner.EvaluateExpr(attr.Expr, func(roleARN string) error {
					creds.AssumeRoleARN = roleARN
					return nil
				}, opts); err != nil {
					return nil, err
				}
			}

			if attr, exists := assumeRole.Body.Attributes["session_name"]; exists {
				if err := runner.EvaluateExpr(attr.Expr, func(sessionName string) error {
					creds.AssumeRoleSessionName = sessionName
					return nil
				}, opts); err != nil {
					return nil, err
				}
			}

			if attr, exists := assumeRole.Body.Attributes["external_id"]; exists {
				if err := runner.EvaluateExpr(attr.Expr, func(id string) error {
					creds.AssumeRoleExternalID = id
					return nil
				}, opts); err != nil {
					return nil, err
				}
			}

			if attr, exists := assumeRole.Body.Attributes["policy"]; exists {
				if err := runner.EvaluateExpr(attr.Expr, func(policy string) error {
					creds.AssumeRolePolicy = policy
					return nil
				}, opts); err != nil {
					return nil, err
				}
			}
		}
		if attr, exists := provider.Body.Attributes["alias"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, func(alias string) error {
				m[alias] = creds
				return nil
			}, opts); err != nil {
				return nil, err
			}
		} else {
			m["aws"] = creds
		}
	}

	return m, nil
}
