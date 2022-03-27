package aws

// Config is the configuration for the ruleset.
type Config struct {
	DeepCheck             bool   `hclext:"deep_check,optional"`
	AccessKey             string `hclext:"access_key,optional"`
	SecretKey             string `hclext:"secret_key,optional"`
	Region                string `hclext:"region,optional"`
	Profile               string `hclext:"profile,optional"`
	SharedCredentialsFile string `hclext:"shared_credentials_file,optional"`
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
