package aws

import "github.com/hashicorp/hcl/v2"

// Config is the configuration for the ruleset.
type Config struct {
	DeepCheck             bool   `hcl:"deep_check,optional"`
	AccessKey             string `hcl:"access_key,optional"`
	SecretKey             string `hcl:"secret_key,optional"`
	Region                string `hcl:"region,optional"`
	Profile               string `hcl:"profile,optional"`
	SharedCredentialsFile string `hcl:"shared_credentials_file,optional"`

	Remain hcl.Body `hcl:",remain"`
}

func (c *Config) toCredentials() Credentials {
	return Credentials{
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
		Region:    c.Region,
		Profile:   c.Profile,
		CredsFile: c.SharedCredentialsFile,
	}
}
