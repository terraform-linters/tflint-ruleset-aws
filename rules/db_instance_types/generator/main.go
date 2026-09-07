//go:build generators

package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/terraform-linters/tflint-ruleset-aws/rules/pricelist"
)

const (
	instanceTypeColumn = "Instance Type"
	generationColumn   = "Current Generation"
	previousGeneration = "No"

	// instanceClassPrefix separates RDS DB instance classes from the other
	// instance types the offer prices.
	instanceClassPrefix = "db."

	// Every region prices hundreds of classes, so a result below this floor
	// indicates the offer file layout changed. RDS has only ever added classes,
	// so the floor sits just under the count at the time of writing.
	minInstanceClasses = 400
)

type instanceClasses struct {
	InstanceClasses           []string `json:"instance_classes"`
	PreviousGenerationClasses []string `json:"previous_generation_classes"`
}

func main() {
	classes, err := readInstanceClasses()
	if err != nil {
		panic(fmt.Sprintf("reading instance classes: %s", err))
	}
	if len(classes.InstanceClasses) < minInstanceClasses {
		panic(fmt.Sprintf(
			"found only %d instance classes, expected at least %d",
			len(classes.InstanceClasses), minInstanceClasses,
		))
	}
	if len(classes.PreviousGenerationClasses) == 0 {
		panic("found no previous generation classes")
	}

	data, err := json.MarshalIndent(classes, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshaling JSON: %s", err))
	}
	if err := os.WriteFile("instance_classes.json", append(data, '\n'), 0644); err != nil {
		panic(fmt.Sprintf("writing JSON: %s", err))
	}

	fmt.Printf(
		"Wrote %d instance classes, %d of them previous generation, to instance_classes.json\n",
		len(classes.InstanceClasses), len(classes.PreviousGenerationClasses),
	)
}

func readInstanceClasses() (instanceClasses, error) {
	products, err := pricelist.Distinct(pricelist.RDS, instanceTypeColumn, generationColumn)
	if err != nil {
		return instanceClasses{}, err
	}

	classes := map[string]bool{}
	previous := map[string]bool{}
	for _, product := range products {
		instanceClass := product.Get(instanceTypeColumn)
		if !strings.HasPrefix(instanceClass, instanceClassPrefix) {
			continue
		}

		classes[instanceClass] = true
		if product.Get(generationColumn) == previousGeneration {
			previous[instanceClass] = true
		}
	}

	return instanceClasses{
		InstanceClasses:           slices.Sorted(maps.Keys(classes)),
		PreviousGenerationClasses: slices.Sorted(maps.Keys(previous)),
	}, nil
}
