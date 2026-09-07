//go:generate go run -tags generators ./generator

// Package db_instance_types provides the DB instance classes that Amazon RDS
// offers, such as db.t4g.micro, generated from the AWS Price List.
package db_instance_types

import (
	_ "embed"
	"encoding/json"

	"github.com/terraform-linters/tflint-ruleset-aws/rules/stringset"
)

//go:embed instance_classes.json
var instanceClassesJSON []byte

var instanceClasses struct {
	All                stringset.Set `json:"instance_classes"`
	PreviousGeneration stringset.Set `json:"previous_generation_classes"`
}

func init() {
	if err := json.Unmarshal(instanceClassesJSON, &instanceClasses); err != nil {
		panic(err)
	}
}

// Valid returns whether RDS offers the given DB instance class.
func Valid(instanceClass string) bool {
	return instanceClasses.All[instanceClass]
}

// PreviousGeneration returns whether the given DB instance class belongs to a
// previous generation, such as db.m1.small.
func PreviousGeneration(instanceClass string) bool {
	return instanceClasses.PreviousGeneration[instanceClass]
}
