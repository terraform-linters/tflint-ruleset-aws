package utils

import (
	"encoding/json"
	"io/ioutil"
)

type Schema struct {
	ProviderSchemas ProviderSchemas `json:"provider_schemas"`
}

type ProviderSchemas struct {
	AWS ProviderSchema `json:"registry.terraform.io/hashicorp/aws"`
}

type ProviderSchema struct {
	ResourceSchemas map[string]ResourceSchema `json:"resource_schemas"`
}

type ResourceSchema struct {
	Block BlockSchema `json:"block"`
}

type BlockSchema struct {
	Attributes map[string]AttributeSchema `json:"attributes"`
	BlockTypes map[string]ResourceSchema  `json:"block_types"`
}

type AttributeSchema struct {
	Type      interface{} `json:"type"`
	Sensitive bool        `json:"sensitive"`
}

func LoadProviderSchema(path string) ProviderSchema {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var schema Schema
	if err := json.Unmarshal(src, &schema); err != nil {
		panic(err)
	}
	return schema.ProviderSchemas.AWS
}
