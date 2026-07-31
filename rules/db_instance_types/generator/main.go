//go:build generators

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/terraform-linters/tflint-ruleset-aws/rules/pricelist"
)

const (
	instanceClassPrefix = "db."

	// Every region publishes hundreds of classes, so a result below this floor
	// indicates the offer file layout changed.
	minInstanceClasses = 300
)

func main() {
	classes, err := pricelist.Values(pricelist.RDS, "Instance Type", func(product pricelist.Product) bool {
		return strings.HasPrefix(product.Get("Instance Type"), instanceClassPrefix)
	})
	if err != nil {
		panic(fmt.Sprintf("reading instance classes: %s", err))
	}
	if len(classes) < minInstanceClasses {
		panic(fmt.Sprintf("found only %d instance classes, expected at least %d", len(classes), minInstanceClasses))
	}

	data, err := json.MarshalIndent(classes, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshaling JSON: %s", err))
	}

	if err := os.WriteFile("instance_classes.json", append(data, '\n'), 0644); err != nil {
		panic(fmt.Sprintf("writing JSON: %s", err))
	}

	fmt.Printf("Wrote %d instance classes to instance_classes.json\n", len(classes))
}
