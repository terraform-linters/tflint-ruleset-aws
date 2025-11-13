//go:build generators

package main

import (
	"fmt"
	"regexp"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	utils "github.com/terraform-linters/tflint-ruleset-aws/rules/generator-utils"
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

func generateRuleFile(resource, attribute string, model map[string]interface{}, schema *tfjson.SchemaAttribute) {
	ruleName := makeRuleName(resource, attribute)

	meta := &ruleMeta{
		RuleName:      ruleName,
		RuleNameCC:    utils.ToCamel(ruleName),
		ResourceType:  resource,
		AttributeName: attribute,
		Sensitive:     schema.Sensitive,
		Max:           fetchNumber(model, "max"),
		Min:           fetchNumber(model, "min"),
		Pattern:       replacePattern(fetchString(model, "pattern")),
		Enum:          fetchStrings(model, "enum"),
	}

	// Testing generated regexp
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

	// Testing generated regexp
	regexp.MustCompile(meta.Pattern)

	utils.GenerateFile(fmt.Sprintf("%s_test.go", ruleName), "pattern_rule_test.go.tmpl", meta)
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
		return int(v.(float64))
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
				ret[i] = item.(string)
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
		return v.(string)
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
