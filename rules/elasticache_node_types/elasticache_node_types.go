//go:generate go run -tags generators ./generator

// Package elasticache_node_types provides the node types that Amazon
// ElastiCache offers, such as cache.t4g.micro, generated from the AWS Price
// List and the ElastiCache API model.
package elasticache_node_types

import (
	_ "embed"
	"encoding/json"

	"github.com/terraform-linters/tflint-ruleset-aws/rules/stringset"
)

//go:embed node_types.json
var nodeTypesJSON []byte

var nodeTypes struct {
	All                stringset.Set `json:"node_types"`
	PreviousGeneration stringset.Set `json:"previous_generation_node_types"`
}

func init() {
	if err := json.Unmarshal(nodeTypesJSON, &nodeTypes); err != nil {
		panic(err)
	}
}

// Valid returns whether ElastiCache offers the given node type.
func Valid(nodeType string) bool {
	return nodeTypes.All[nodeType]
}

// PreviousGeneration returns whether the given node type belongs to a previous
// generation, such as cache.m1.small.
func PreviousGeneration(nodeType string) bool {
	return nodeTypes.PreviousGeneration[nodeType]
}
