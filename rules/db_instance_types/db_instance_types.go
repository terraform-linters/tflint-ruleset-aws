//go:generate go run -tags generators ./generator

// Package db_instance_types provides the DB instance classes that Amazon RDS
// offers, such as db.t4g.micro, generated from the AWS Price List.
package db_instance_types

import (
	_ "embed"
	"encoding/json"
)

//go:embed instance_classes.json
var instanceClassesJSON []byte

var instanceClasses struct {
	classes            map[string]bool
	previousGeneration map[string]bool
}

func init() {
	var file struct {
		InstanceClasses            []string `json:"instance_classes"`
		PreviousGenerationFamilies []string `json:"previous_generation_families"`
	}
	if err := json.Unmarshal(instanceClassesJSON, &file); err != nil {
		panic(err)
	}

	instanceClasses.classes = set(file.InstanceClasses)
	instanceClasses.previousGeneration = set(file.PreviousGenerationFamilies)
}

// Valid returns whether RDS offers the given DB instance class.
func Valid(instanceClass string) bool {
	return instanceClasses.classes[instanceClass]
}

// PreviousGeneration returns whether the given instance family, such as m1, is
// a previous generation.
func PreviousGeneration(family string) bool {
	return instanceClasses.previousGeneration[family]
}

func set(values []string) map[string]bool {
	members := make(map[string]bool, len(values))
	for _, value := range values {
		members[value] = true
	}

	return members
}
