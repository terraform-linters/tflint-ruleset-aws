package aws

import (
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsProviderBlockSchema is a schema of `aws` provider block
var AwsProviderBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "access_key"},
		{Name: "secret_key"},
		{Name: "profile"},
		{Name: "shared_credentials_file"},
		{Name: "region"},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "assume_role"},
	},
}

// AwsProviderAssumeRoleBlockShema is a schema of `assume_role` block
var AwsProviderAssumeRoleBlockShema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "role_arn", Required: true},
		{Name: "session_name"},
		{Name: "external_id"},
		{Name: "policy"},
	},
}

// ProviderData represents a provider block with an eval context (runner)
type ProviderData struct {
	provider   *configs.Provider
	runner     tflint.Runner
	attributes hcl.Attributes
	blocks     hcl.Blocks
}

// GetCredentialsFromProvider retrieves credentials from the "provider" block in the Terraform configuration
func GetCredentialsFromProvider(runner tflint.Runner) (Credentials, error) {
	creds := Credentials{}

	providerConfig, err := runner.RootProvider("aws")
	if err != nil || providerConfig == nil {
		return creds, err
	}

	d, err := newProviderData(providerConfig, runner)
	if err != nil {
		return creds, err
	}

	accessKey, exists, err := d.Get("access_key")
	if err != nil {
		return creds, err
	}
	if exists {
		creds.AccessKey = accessKey
	}

	secretKey, exists, err := d.Get("secret_key")
	if err != nil {
		return creds, err
	}
	if exists {
		creds.SecretKey = secretKey
	}

	profile, exists, err := d.Get("profile")
	if err != nil {
		return creds, err
	}
	if exists {
		creds.Profile = profile
	}

	credsFile, exists, err := d.Get("shared_credentials_file")
	if err != nil {
		return creds, err
	}
	if exists {
		creds.CredsFile = credsFile
	}

	region, exists, err := d.Get("region")
	if err != nil {
		return creds, err
	}
	if exists {
		creds.Region = region
	}

	assumeRole, exists, err := d.GetBlock("assume_role", AwsProviderAssumeRoleBlockShema)
	if err != nil {
		return creds, err
	}
	if exists {
		roleARN, exists, err := assumeRole.Get("role_arn")
		if err != nil {
			return creds, err
		}
		if exists {
			creds.AssumeRoleARN = roleARN
		}

		sessionName, exists, err := assumeRole.Get("session_name")
		if err != nil {
			return creds, err
		}
		if exists {
			creds.AssumeRoleSessionName = sessionName
		}

		externalID, exists, err := assumeRole.Get("external_id")
		if err != nil {
			return creds, err
		}
		if exists {
			creds.AssumeRoleExternalID = externalID
		}

		policy, exists, err := assumeRole.Get("policy")
		if err != nil {
			return creds, err
		}
		if exists {
			creds.AssumeRolePolicy = policy
		}
	}

	return creds, nil
}

func newProviderData(provider *configs.Provider, runner tflint.Runner) (*ProviderData, error) {
	providerData := &ProviderData{
		provider:   provider,
		runner:     runner,
		attributes: map[string]*hcl.Attribute{},
		blocks:     []*hcl.Block{},
	}

	if provider != nil {
		content, _, diags := provider.Config.PartialContent(AwsProviderBlockSchema)
		if diags.HasErrors() {
			return nil, diags
		}

		providerData.attributes = content.Attributes
		providerData.blocks = content.Blocks
	}

	return providerData, nil
}

// Get returns a value corresponding to the given key
// It should be noted that the value is evaluated if it is evaluable
// The second return value is a flag that determines whether a value exists
// We assume the provider has only simple attributes, so it just returns string
func (d *ProviderData) Get(key string) (string, bool, error) {
	attribute, exists := d.attributes[key]
	if !exists {
		log.Printf("[INFO] `%s` is not found in the provider block.", key)
		return "", false, nil
	}

	var val string
	err := d.runner.EvaluateExprOnRootCtx(attribute.Expr, &val)

	err = d.runner.EnsureNoError(err, func() error { return nil })
	if err != nil {
		return "", true, err
	}
	return val, true, nil
}

// GetBlock returns a value just like Get.
// The difference is that GetBlock returns ProviderData rather than a string value.
func (d *ProviderData) GetBlock(key string, schema *hcl.BodySchema) (*ProviderData, bool, error) {
	providerData := &ProviderData{
		provider:   d.provider,
		runner:     d.runner,
		attributes: map[string]*hcl.Attribute{},
		blocks:     []*hcl.Block{},
	}

	var ret *hcl.Block
	for _, block := range d.blocks {
		if block.Type == key {
			ret = block
		}
	}
	if ret == nil {
		log.Printf("[INFO] `%s` is not found in the provider block.", key)
		return providerData, false, nil
	}

	content, _, diags := ret.Body.PartialContent(schema)
	if diags.HasErrors() {
		return providerData, true, diags
	}

	providerData.attributes = content.Attributes
	providerData.blocks = content.Blocks

	return providerData, true, nil
}
