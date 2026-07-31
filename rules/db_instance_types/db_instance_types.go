//go:generate go run -tags generators ./generator

// Package db_instance_types provides the DB instance classes that Amazon RDS
// offers, such as db.t4g.micro.
package db_instance_types

import (
	_ "embed"
	"encoding/json"
)

// instanceClassesJSON contains every DB instance class that appears in the AWS
// Price List for Amazon RDS, across all regions of the aws partition.
//
//go:embed instance_classes.json
var instanceClassesJSON []byte

var instanceClasses map[string]bool

func init() {
	var classes []string
	if err := json.Unmarshal(instanceClassesJSON, &classes); err != nil {
		panic(err)
	}

	instanceClasses = make(map[string]bool, len(classes))
	for _, class := range classes {
		instanceClasses[class] = true
	}
}

// Valid returns whether the given DB instance class is offered by RDS.
func Valid(instanceClass string) bool {
	return instanceClasses[instanceClass]
}
