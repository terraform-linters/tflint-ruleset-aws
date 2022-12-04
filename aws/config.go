package aws

type AssumeRole struct {
	RoleARN     string `hclext:"role_arn,optional"`
	ExternalID  string `hclext:"external_id,optional"`
	Policy      string `hclext:"policy,optional"`
	SessionName string `hclext:"session_name,optional"`
}

// Config is the configuration for the ruleset.
type Config struct {
	DeepCheck              bool       `hclext:"deep_check,optional"`
	AccessKey              string     `hclext:"access_key,optional"`
	SecretKey              string     `hclext:"secret_key,optional"`
	Region                 string     `hclext:"region,optional"`
	Profile                string     `hclext:"profile,optional"`
	SharedCredentialsFile  string     `hclext:"shared_credentials_file,optional"`
	AssumeRole            *AssumeRole `hclext:"assume_role,block"`
}

func (c *Config) toCredentials() Credentials {
	credentials := Credentials{
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
		Region:    c.Region,
		Profile:   c.Profile,
		CredsFile: c.SharedCredentialsFile,
	}

	if c.AssumeRole != nil {
		credentials.AssumeRoleARN         = c.AssumeRole.RoleARN
		credentials.AssumeRoleExternalID  = c.AssumeRole.ExternalID
		credentials.AssumeRolePolicy      = c.AssumeRole.Policy
		credentials.AssumeRoleSessionName = c.AssumeRole.SessionName
	}

	return credentials
}
