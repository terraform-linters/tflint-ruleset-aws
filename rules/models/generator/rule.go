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
		// Handle both []interface{} (legacy format) and []string (converted Smithy format)
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
	
	// Handle placeholder patterns that shouldn't be used as regex
	if pattern == "^See rules in parameter description$" {
		return ""
	}
	
	reg := regexp.MustCompile(`\\u([0-9A-F]{4})`)
	replaced := reg.ReplaceAllString(pattern, `\x{$1}`)

	// Handle incomplete patterns like "^(?s)" which should be "^(?s).*$"
	if replaced == "^(?s)" {
		return "^(?s).*$"
	}

	if !strings.HasPrefix(replaced, "^") && !strings.HasSuffix(replaced, "$") {
		// Handle patterns that should preserve original format for minimal diff
		if replaced == "\\S" {
			// Use AWS Service Model Definition compatible format for minimal diff
			return "^.*\\S.*$"
		}
		return fmt.Sprintf("^%s$", replaced)
	}
	
	// Handle patterns that need to preserve original format for backward compatibility
	// These transformations maintain the same validation behavior while minimizing diffs
	
	// Handle patterns with $ anchors that should be removed for backward compatibility
	if strings.HasPrefix(replaced, "^") && strings.HasSuffix(replaced, "$") {
		// Simple quantified character class patterns that originally didn't have $ anchors
		// e.g., "^[a-zA-Z0-9_]{2,}$" should become "^[a-zA-Z0-9_]{2,}"
		// Only apply to simple patterns that start with ^[ and end with }$ with no other complex parts
		if strings.HasPrefix(replaced, "^[") && strings.Count(replaced, "[") == 1 && strings.Count(replaced, "]") == 1 {
			bracketPos := strings.Index(replaced, "]")
			if bracketPos > 0 && bracketPos < len(replaced)-3 && strings.Contains(replaced[bracketPos:], "]{") {
				return strings.TrimSuffix(replaced, "$")
			}
		}
	}
	
	// Convert common Smithy prefix patterns back to original format for backward compatibility
	if strings.HasPrefix(replaced, "^") && !strings.HasSuffix(replaced, "$") {
		// Special case: patterns like "^[^:]" should stay as-is (no $ anchor)
		if strings.Contains(replaced, "[^") {
			return replaced
		}
		
		// Patterns ending with role/ should stay as-is for backward compatibility
		if strings.HasSuffix(replaced, ":role/") {
			return replaced
		}
		
		// Simple prefix patterns that should get .* appended  
		simplePrefixes := []string{
			"^arn:aws",
			"^arn:",
			"^s3://",
			"^build-",
		}
		
		for _, prefix := range simplePrefixes {
			if replaced == prefix {
				return prefix + ".*"
			}
		}
	}
	
	return replaced
}

func formatTest(body string) string {
	if strings.Contains(body, "\n") {
		return fmt.Sprintf("<<TEXT\n%sTEXT", body)
	}
	return fmt.Sprintf(`"%s"`, body)
}
