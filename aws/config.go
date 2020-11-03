package aws

import "github.com/hashicorp/hcl/v2"

// Config is the configuration for the ruleset.
type Config struct {
	AccessKey             string `hcl:"access_key"`
	SecretKey             string `hcl:"secret_key"`
	Region                string `hcl:"region"`
	Profile               string `hcl:"profile"`
	SharedCredentialsFile string `hcl:"shared_credentials_file"`

	Remain hcl.Body `hcl:",remain"`
}
