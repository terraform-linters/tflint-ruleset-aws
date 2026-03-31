package genutils

import (
	"encoding/json"
	"os"

	tfjson "github.com/hashicorp/terraform-json"
)

func LoadProviderSchema(path string) *tfjson.ProviderSchema {
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var schema tfjson.ProviderSchemas
	if err := json.Unmarshal(src, &schema); err != nil {
		panic(err)
	}
	return schema.Schemas["registry.terraform.io/hashicorp/aws"]
}
