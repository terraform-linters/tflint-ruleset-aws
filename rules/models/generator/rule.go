//go:build generators

package main

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	utils "github.com/terraform-linters/tflint-ruleset-aws/rules/generator-utils"
	"github.com/zclconf/go-cty/cty"
)

type ruleMeta struct {
	RuleName      string
	RuleNameCC    string
	ResourceType  string
	AttributeName string
	Sensitive     bool
	Max           int
	Min           int
	Pattern       string
	Enum          []string
	TestOK        string
	TestNG        string
}

type mapRuleMeta struct {
	RuleName      string
	RuleNameCC    string
	ResourceType  string
	AttributeName string
	Sensitive     bool
	ItemsMax      int
	ItemsMin      int
	KeyMax        int
	KeyMin        int
	KeyPattern    string
	KeyEnum       []string
	ValueMax      int
	ValueMin      int
	ValuePattern  string
	ValueEnum     []string
}

func generateRuleFile(resource, attribute string, model map[string]interface{}, schema *tfjson.SchemaAttribute) {
	ruleName := makeRuleName(resource, attribute)

	meta := &ruleMeta{
		RuleName:      ruleName,
		RuleNameCC:    utils.ToCamel(ruleName),
		ResourceType:  resource,
		AttributeName: attribute,
		Sensitive:     schema != nil && schema.Sensitive,
		Max:           fetchNumber(model, "max"),
		Min:           fetchNumber(model, "min"),
		Pattern:       replacePattern(fetchString(model, "pattern")),
		Enum:          fetchStrings(model, "enum"),
	}

	// Validate generated regexp
	regexp.MustCompile(meta.Pattern)

	utils.GenerateFile(fmt.Sprintf("%s.go", ruleName), "pattern_rule.go.tmpl", meta)
}

func generateRuleTestFile(resource, attribute string, model map[string]interface{}, test test) {
	ruleName := makeRuleName(resource, attribute)

	meta := &ruleMeta{
		RuleName:      ruleName,
		RuleNameCC:    utils.ToCamel(ruleName),
		ResourceType:  resource,
		AttributeName: attribute,
		Max:           fetchNumber(model, "max"),
		Min:           fetchNumber(model, "min"),
		Pattern:       replacePattern(fetchString(model, "pattern")),
		Enum:          fetchStrings(model, "enum"),
		TestOK:        formatTest(test.OK),
		TestNG:        formatTest(test.NG),
	}

	// Validate generated regexp
	regexp.MustCompile(meta.Pattern)

	utils.GenerateFile(fmt.Sprintf("%s_test.go", ruleName), "pattern_rule_test.go.tmpl", meta)
}

func generateMapRuleFile(resource, attribute string, listModel, keyModel, valueModel map[string]interface{}, schema *tfjson.SchemaAttribute) bool {
	ruleName := makeRuleName(resource, attribute)

	var itemsMax, itemsMin int
	if listModel != nil {
		itemsMax = fetchNumber(listModel, "max")
		itemsMin = fetchNumber(listModel, "min")
	}
	keyPattern := replacePattern(fetchString(keyModel, "pattern"))
	valuePattern := replacePattern(fetchString(valueModel, "pattern"))
	keyMax := fetchNumber(keyModel, "max")
	keyMin := fetchNumber(keyModel, "min")
	valueMax := fetchNumber(valueModel, "max")
	valueMin := fetchNumber(valueModel, "min")

	// Skip if no constraints exist
	hasItemsConstraints := itemsMax != 0 || itemsMin != 0
	hasKeyConstraints := keyPattern != "" || keyMax != 0 || keyMin != 0
	hasValueConstraints := valuePattern != "" || valueMax != 0 || valueMin != 0
	if !hasItemsConstraints && !hasKeyConstraints && !hasValueConstraints {
		return false
	}

	// Test generated regexps - skip if malformed
	if keyPattern != "" {
		if _, err := regexp.Compile(keyPattern); err != nil {
			fmt.Printf("Skipping %s.%s: malformed key pattern %q: %v\n", resource, attribute, keyPattern, err)
			return false
		}
	}
	if valuePattern != "" {
		if _, err := regexp.Compile(valuePattern); err != nil {
			fmt.Printf("Skipping %s.%s: malformed value pattern %q: %v\n", resource, attribute, valuePattern, err)
			return false
		}
	}

	meta := &mapRuleMeta{
		RuleName:      ruleName,
		RuleNameCC:    utils.ToCamel(ruleName),
		ResourceType:  resource,
		AttributeName: attribute,
		Sensitive:     schema != nil && schema.Sensitive,
		ItemsMax:      itemsMax,
		ItemsMin:      itemsMin,
		KeyMax:        keyMax,
		KeyMin:        keyMin,
		KeyPattern:    keyPattern,
		ValueMax:      valueMax,
		ValueMin:      valueMin,
		ValuePattern:  valuePattern,
	}

	utils.GenerateFile(fmt.Sprintf("%s.go", ruleName), "map_rule.go.tmpl", meta)
	return true
}

// generateListmapRuleFile generates a map validation rule from a resolved listmap result.
// The result contains fully resolved constraints from list, key, and value shapes.
func generateListmapRuleFile(resource, attribute string, result cty.Value, schema *tfjson.SchemaAttribute) bool {
	ruleName := makeRuleName(resource, attribute)

	itemsMin := ctyToInt(result.GetAttr("itemsMin"))
	itemsMax := ctyToInt(result.GetAttr("itemsMax"))
	keyPattern := replacePattern(result.GetAttr("keyPattern").AsString())
	keyMin := ctyToInt(result.GetAttr("keyMin"))
	keyMax := ctyToInt(result.GetAttr("keyMax"))
	keyEnum := ctyToStringSlice(result.GetAttr("keyEnum"))
	valuePattern := replacePattern(result.GetAttr("valuePattern").AsString())
	valueMin := ctyToInt(result.GetAttr("valueMin"))
	valueMax := ctyToInt(result.GetAttr("valueMax"))
	valueEnum := ctyToStringSlice(result.GetAttr("valueEnum"))

	// Skip if no constraints exist
	hasItemsConstraints := itemsMax != 0 || itemsMin != 0
	hasKeyConstraints := keyPattern != "" || keyMax != 0 || keyMin != 0 || len(keyEnum) > 0
	hasValueConstraints := valuePattern != "" || valueMax != 0 || valueMin != 0 || len(valueEnum) > 0
	if !hasItemsConstraints && !hasKeyConstraints && !hasValueConstraints {
		return false
	}

	// Test generated regexps - skip if malformed
	if keyPattern != "" {
		if _, err := regexp.Compile(keyPattern); err != nil {
			fmt.Printf("Skipping %s.%s: malformed key pattern %q: %v\n", resource, attribute, keyPattern, err)
			return false
		}
	}
	if valuePattern != "" {
		if _, err := regexp.Compile(valuePattern); err != nil {
			fmt.Printf("Skipping %s.%s: malformed value pattern %q: %v\n", resource, attribute, valuePattern, err)
			return false
		}
	}

	meta := &mapRuleMeta{
		RuleName:      ruleName,
		RuleNameCC:    utils.ToCamel(ruleName),
		ResourceType:  resource,
		AttributeName: attribute,
		Sensitive:     schema != nil && schema.Sensitive,
		ItemsMax:      itemsMax,
		ItemsMin:      itemsMin,
		KeyMax:        keyMax,
		KeyMin:        keyMin,
		KeyPattern:    keyPattern,
		KeyEnum:       keyEnum,
		ValueMax:      valueMax,
		ValueMin:      valueMin,
		ValuePattern:  valuePattern,
		ValueEnum:     valueEnum,
	}

	utils.GenerateFile(fmt.Sprintf("%s.go", ruleName), "map_rule.go.tmpl", meta)
	return true
}

func ctyToStringSlice(val cty.Value) []string {
	if val.IsNull() || val.LengthInt() == 0 {
		return nil
	}
	result := make([]string, 0, val.LengthInt())
	for it := val.ElementIterator(); it.Next(); {
		_, v := it.Element()
		result = append(result, v.AsString())
	}
	return result
}

func ctyToInt(val cty.Value) int {
	bf := val.AsBigFloat()
	i, _ := bf.Int(new(big.Int))
	return int(i.Int64())
}

func makeRuleName(resource, attribute string) string {
	// XXX: Change the naming convention for the backward compatibility.
	if resource == "aws_instance" && attribute == "instance_type" {
		return "aws_instance_invalid_type"
	}
	if resource == "aws_launch_configuration" && attribute == "instance_type" {
		return "aws_launch_configuration_invalid_type"
	}
	return fmt.Sprintf("%s_invalid_%s", resource, attribute)
}

func fetchNumber(model map[string]interface{}, key string) int {
	if v, ok := model[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		case int64:
			return int(n)
		}
	}
	return 0
}

func fetchStrings(model map[string]interface{}, key string) []string {
	if raw, ok := model[key]; ok {
		// Handle both []interface{} and []string for compatibility
		switch v := raw.(type) {
		case []interface{}:
			ret := make([]string, len(v))
			for i, item := range v {
				str, ok := item.(string)
				if !ok {
					return []string{}
				}
				ret[i] = str
			}
			return ret
		case []string:
			return v
		}
	}
	return []string{}
}

func fetchString(model map[string]interface{}, key string) string {
	if v, ok := model[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func replacePattern(pattern string) string {
	if pattern == "" {
		return pattern
	}

	// Convert Unicode escapes from \uXXXX to \x{XXXX} format
	unicodeEscapeRegex := regexp.MustCompile(`\\u([0-9A-F]{4})`)
	replaced := unicodeEscapeRegex.ReplaceAllString(pattern, `\x{$1}`)

	// Remove negative lookahead patterns that Go's regexp doesn't support
	// Common patterns like (?!aws:) are used to prevent aws: prefix in tag keys
	// We strip these and still validate the character set constraint
	negativeLookaheadRegex := regexp.MustCompile(`\(\?![^)]+\)`)
	replaced = negativeLookaheadRegex.ReplaceAllString(replaced, "")

	// Handle patterns missing anchors
	// The Ruby SDK generator ensured all patterns had both ^ and $ anchors.
	// Smithy models often have incomplete patterns (e.g., "^(?s)" or "[a-z]*").
	// We add missing anchors to maintain backward compatibility.
	hasPrefix := strings.HasPrefix(replaced, "^")
	hasSuffix := strings.HasSuffix(replaced, "$")

	if !hasPrefix && !hasSuffix {
		// No anchors at all
		// Special case: \S in Smithy means "contains non-whitespace", not "is single non-whitespace char"
		if replaced == "\\S" {
			return "^.*\\S.*$"
		}
		replaced = fmt.Sprintf("^%s$", replaced)
	} else if !hasPrefix {
		// Has $ but missing ^
		replaced = "^" + replaced
	} else if !hasSuffix {
		// Has ^ but missing $
		// Check if pattern is just a modifier with no actual content
		// e.g., "^(?s)" or "^(?i)" - these need .*$ to be useful
		bareModifiers := []string{"^(?s)", "^(?i)", "^(?m)", "^(?s)(?i)", "^(?i)(?s)"}
		isBareModifier := false
		for _, mod := range bareModifiers {
			if replaced == mod {
				isBareModifier = true
				break
			}
		}

		if isBareModifier {
			// Bare modifier - add .*$ to match any content
			// e.g., "^(?s)" becomes "^(?s).*$"
			replaced = replaced + ".*$"
		} else {
			// Check if pattern looks like a prefix pattern (ends with delimiter or literal)
			// These should match "starts with X" not "equals exactly X"
			isPrefixPattern := false
			lastChar := replaced[len(replaced)-1]

			// Common delimiters in ARNs, URIs, paths that indicate "starts with"
			if lastChar == ':' || lastChar == '/' {
				isPrefixPattern = true
			}
			// Patterns ending with quantifiers already match variable length
			// e.g., "^[a-z]*" matches zero or more, just needs $
			// But patterns without quantifiers need .*$ for variable length
			// e.g., "^arn" should be "^arn.*$" not "^arn$"
			if lastChar == '*' || lastChar == '+' || lastChar == '?' || lastChar == '}' {
				isPrefixPattern = false // Has quantifier, just add $
			} else if lastChar == ']' || lastChar == ')' || lastChar >= 'a' && lastChar <= 'z' || lastChar >= 'A' && lastChar <= 'Z' {
				// Ends with literal character, character class, or group without quantifier
				// Check if this is a pattern that should match exact length or variable length
				// For now, treat patterns without quantifiers as needing .*$ for safety
				// This matches Ruby SDK behavior
				isPrefixPattern = true
			}

			if isPrefixPattern {
				// Prefix pattern - add .*$ to match remainder
				// e.g., "^arn:" becomes "^arn:.*$", "^s3://" becomes "^s3://.*$"
				// e.g., "^[a-z]" becomes "^[a-z].*$"
				replaced = replaced + ".*$"
			} else {
				// Pattern with quantifier - just add $
				// e.g., "^[a-z]*" becomes "^[a-z]*$"
				replaced = replaced + "$"
			}
		}
	}

	// Apply compatibility transforms to maintain backward compatibility with Ruby SDK
	if transformed := applyCompatibilityTransforms(replaced); transformed != "" {
		return transformed
	}

	return replaced
}

func applyCompatibilityTransforms(pattern string) string {
	// IAM role ARN patterns: Smithy models end at the role prefix (e.g., "role/")
	// but Terraform configurations include role paths (e.g., "role/my-app/my-role").
	// Add .* suffix to match the role name/path after the prefix.
	arnRolePatterns := map[string]string{
		"^arn:aws:iam::[0-9]*:role/":                  "^arn:aws:iam::[0-9]*:role/.*$",
		"^arn:aws:iam::\\d{12}:role/":                 "^arn:aws:iam::\\d{12}:role/.*$",
		"^arn:aws(-[\\w]+)*:iam::[0-9]{12}:role/":     "^arn:aws(-[\\w]+)*:iam::[0-9]{12}:role/.*$",
		"^arn:aws(-[a-z]{1,3}){0,2}:iam::\\d+:role/": "^arn:aws(-[a-z]{1,3}){0,2}:iam::\\d+:role/.*$",
	}

	// Cognito SMS authentication messages: The {####} placeholder must appear
	// within a message, not be the entire message. Transform from matching
	// literal "{####}" to matching any message containing "{####}".
	// Example valid message: "Your code is {####}. Do not share it."
	cognitoMessagePatterns := map[string]string{
		"^\\{####\\}$": "^.*\\{####\\}.*$",
	}

	for _, transforms := range []map[string]string{arnRolePatterns, cognitoMessagePatterns} {
		if transformed, exists := transforms[pattern]; exists {
			return transformed
		}
	}

	return ""
}

func formatTest(body string) string {
	if strings.Contains(body, "\n") {
		return fmt.Sprintf("<<TEXT\n%sTEXT", body)
	}
	return fmt.Sprintf(`"%s"`, body)
}
