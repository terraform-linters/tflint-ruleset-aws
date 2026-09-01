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
	"golang.org/x/net/html"
)

const (
	instanceTypeColumn = "Instance Type"

	// nodeTypePrefix separates node types from the other products the offer
	// prices, such as serverless capacity and backup storage.
	nodeTypePrefix = "cache."

	// Every current generation node type is priced in at least one region, so a
	// result below this floor indicates the offer file layout changed. The floor
	// sits just under the count at the time of writing.
	minNodeTypes = 90

	// apiModel holds the ElastiCache API model. The api-models-aws submodule
	// sits under rules/models, which generates the pattern rules from it.
	apiModel = "../models/api-models-aws/models/elasticache/service/2015-02-02/elasticache-2015-02-02.json"

	// nodeTypeShape documents the node types in prose, the only place AWS
	// publishes the retired ones. Several shapes carry the same list. This is
	// the one that creates a cluster.
	nodeTypeShape  = "com.amazonaws.elasticache#CreateCacheClusterMessage"
	nodeTypeMember = "CacheNodeType"

	// previousGenerationHeading opens the list item holding a category's retired
	// node types. The model marks them no other way.
	previousGenerationHeading = "Previous generation"

	// AWS retired six families and has added none since, so anything short of
	// the count at the time of writing means the prose was restructured.
	minPreviousGenerationNodeTypes = 18
)

type nodeTypes struct {
	NodeTypes                   []string `json:"node_types"`
	PreviousGenerationNodeTypes []string `json:"previous_generation_node_types"`
}

func main() {
	current, err := currentNodeTypes()
	if err != nil {
		panic(fmt.Sprintf("reading the price list: %s", err))
	}
	previous, err := previousGenerationNodeTypes()
	if err != nil {
		panic(fmt.Sprintf("reading the API model: %s", err))
	}

	// Each source answers only what it knows. The price list carries every node
	// type AWS currently sells but drops retired ones outright rather than
	// labeling them, and the API model names the retired ones but trails the
	// price list by whole families.
	maps.Copy(current, previous)

	types := nodeTypes{
		NodeTypes:                   slices.Sorted(maps.Keys(current)),
		PreviousGenerationNodeTypes: slices.Sorted(maps.Keys(previous)),
	}
	atLeast(len(types.NodeTypes), minNodeTypes, "node types")
	atLeast(len(types.PreviousGenerationNodeTypes), minPreviousGenerationNodeTypes, "previous generation node types")

	data, err := json.MarshalIndent(types, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshaling JSON: %s", err))
	}
	if err := os.WriteFile("node_types.json", append(data, '\n'), 0644); err != nil {
		panic(fmt.Sprintf("writing JSON: %s", err))
	}

	fmt.Printf(
		"Wrote %d node types, %d of them previous generation, to node_types.json\n",
		len(types.NodeTypes), len(types.PreviousGenerationNodeTypes),
	)
}

func atLeast(count, minimum int, name string) {
	if count < minimum {
		panic(fmt.Sprintf("found only %d %s, expected at least %d", count, name, minimum))
	}
}

func currentNodeTypes() (map[string]bool, error) {
	products, err := pricelist.Distinct(pricelist.ElastiCache, instanceTypeColumn)
	if err != nil {
		return nil, err
	}

	types := map[string]bool{}
	for _, product := range products {
		if nodeType := product.Get(instanceTypeColumn); strings.HasPrefix(nodeType, nodeTypePrefix) {
			types[nodeType] = true
		}
	}

	return types, nil
}

func previousGenerationNodeTypes() (map[string]bool, error) {
	documentation, err := nodeTypeDocumentation()
	if err != nil {
		return nil, err
	}

	document, err := html.Parse(strings.NewReader(documentation))
	if err != nil {
		return nil, fmt.Errorf("parsing %s documentation: %w", nodeTypeMember, err)
	}

	types := map[string]bool{}
	for item := range document.Descendants() {
		if item.Type != html.ElementNode || item.Data != "li" || !previousGeneration(item) {
			continue
		}

		for node := range item.Descendants() {
			if node.Type != html.ElementNode || node.Data != "code" {
				continue
			}

			if nodeType := strings.TrimSpace(text(node)); strings.HasPrefix(nodeType, nodeTypePrefix) {
				types[nodeType] = true
			}
		}
	}

	return types, nil
}

// previousGeneration reports whether a list item introduces retired node types.
// Each category, such as "Memory optimized", nests a list item per generation
// under a paragraph naming it. Only that paragraph opens with the heading, so
// the family paragraphs beside it cannot match.
func previousGeneration(item *html.Node) bool {
	for node := range item.ChildNodes() {
		if node.Type == html.ElementNode && node.Data == "p" &&
			strings.HasPrefix(strings.TrimSpace(text(node)), previousGenerationHeading) {
			return true
		}
	}

	return false
}

func text(node *html.Node) string {
	var builder strings.Builder
	for descendant := range node.Descendants() {
		if descendant.Type == html.TextNode {
			builder.WriteString(descendant.Data)
		}
	}

	return builder.String()
}

func nodeTypeDocumentation() (string, error) {
	raw, err := os.ReadFile(apiModel)
	if err != nil {
		return "", err
	}

	var model struct {
		Shapes map[string]struct {
			Members map[string]struct {
				Traits struct {
					Documentation string `json:"smithy.api#documentation"`
				} `json:"traits"`
			} `json:"members"`
		} `json:"shapes"`
	}
	if err := json.Unmarshal(raw, &model); err != nil {
		return "", err
	}

	member, ok := model.Shapes[nodeTypeShape].Members[nodeTypeMember]
	if !ok {
		return "", fmt.Errorf("no %s member of %s in %s", nodeTypeMember, nodeTypeShape, apiModel)
	}
	if member.Traits.Documentation == "" {
		return "", fmt.Errorf("no documentation for %s of %s", nodeTypeMember, nodeTypeShape)
	}

	return member.Traits.Documentation, nil
}
