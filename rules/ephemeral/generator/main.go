//go:build generators

package main

import (
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	utils "github.com/terraform-linters/tflint-ruleset-aws/rules/generator-utils"
)

type writeOnlyArgument struct {
	OriginalAttribute         string
	WriteOnlyAlternative      string
	WriteOnlyVersionAttribute string
}

func main() {
	awsProvider := utils.LoadProviderSchema("../../tools/provider-schema/schema.json")

	resourcesWithWriteOnly := map[string][]writeOnlyArgument{}
	// Iterate over all resources in the AWS provider schema
	for resourceName, resource := range awsProvider.ResourceSchemas {
		if arguments := writeOnlyArguments(resource); len(arguments) > 0 {
			// gather sensitive attributes with write only argument alternatives
			resourcesWithWriteOnly[resourceName] = findReplaceableAttribute(arguments, resource)
		}
	}

	// Generate the write-only arguments variable to file
	utils.GenerateFile("../../rules/ephemeral/write_only_arguments_gen.go", "../../rules/ephemeral/write_only_arguments_gen.go.tmpl", resourcesWithWriteOnly)
}

func findReplaceableAttribute(arguments []string, resource *tfjson.Schema) []writeOnlyArgument {
	writeOnlyArguments := []writeOnlyArgument{}

	for _, argument := range arguments {
		// Check if the argument ends with "_wo" and if the original attribute without "_wo" suffix exists in the resource schema
		attribute := strings.TrimSuffix(argument, "_wo")
		versionAttribute := attribute + "_wo_version"
		if strings.HasSuffix(argument, "_wo") && resource.Block.Attributes[attribute] != nil && resource.Block.Attributes[versionAttribute] != nil {
			writeOnlyArguments = append(writeOnlyArguments, writeOnlyArgument{
				OriginalAttribute:         attribute,
				WriteOnlyAlternative:      argument,
				WriteOnlyVersionAttribute: versionAttribute,
			})
		}
	}

	return writeOnlyArguments
}

func writeOnlyArguments(resource *tfjson.Schema) []string {
	if resource == nil || resource.Block == nil {
		return []string{}
	}

	writeOnlyArguments := []string{}

	// Check if the resource has any write-only attributes
	for name, attribute := range resource.Block.Attributes {
		if attribute.WriteOnly {
			writeOnlyArguments = append(writeOnlyArguments, name)
		}
	}

	return writeOnlyArguments
}
