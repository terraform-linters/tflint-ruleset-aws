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
func GetCredentialsFromProvider(runner tflint.Runner) (Credentials, error) {
	creds := Credentials{}

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
		return creds, err
	}

	for _, provider := range providers.Blocks {
		if provider.Labels[0] != "aws" {
			continue
		}

		opts := &tflint.EvaluateExprOption{ModuleCtx: tflint.RootModuleCtxType}

		if attr, exists := provider.Body.Attributes["access_key"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, &creds.AccessKey, opts); err != nil {
				return creds, err
			}
		}

		if attr, exists := provider.Body.Attributes["secret_key"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, &creds.SecretKey, opts); err != nil {
				return creds, err
			}
		}

		if attr, exists := provider.Body.Attributes["profile"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, &creds.Profile, opts); err != nil {
				return creds, err
			}
		}

		if attr, exists := provider.Body.Attributes["shared_credentials_file"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, &creds.CredsFile, opts); err != nil {
				return creds, err
			}
		}

		if attr, exists := provider.Body.Attributes["region"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, &creds.Region, opts); err != nil {
				return creds, err
			}
		}

		for _, assumeRole := range provider.Body.Blocks {
			if attr, exists := assumeRole.Body.Attributes["role_arn"]; exists {
				if err := runner.EvaluateExpr(attr.Expr, &creds.AssumeRoleARN, opts); err != nil {
					return creds, err
				}
			}

			if attr, exists := assumeRole.Body.Attributes["session_name"]; exists {
				if err := runner.EvaluateExpr(attr.Expr, &creds.AssumeRoleSessionName, opts); err != nil {
					return creds, err
				}
			}

			if attr, exists := assumeRole.Body.Attributes["external_id"]; exists {
				if err := runner.EvaluateExpr(attr.Expr, &creds.AssumeRoleExternalID, opts); err != nil {
					return creds, err
				}
			}

			if attr, exists := assumeRole.Body.Attributes["policy"]; exists {
				if err := runner.EvaluateExpr(attr.Expr, &creds.AssumeRolePolicy, opts); err != nil {
					return creds, err
				}
			}
		}
	}

	return creds, nil
}
