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
	if !strings.HasPrefix(replaced, "^") && !strings.HasSuffix(replaced, "$") {
		// The Smithy models contain pattern "\S" (single non-whitespace character)
		// for many string fields (ResourceName, ResourceArn, etc). This is clearly
		// incorrect for fields that can be 1-128 characters. The Ruby SDK generator
		// produced "^.*\S.*$" (contains non-whitespace) for these same fields.
		// We maintain this transformation for backward compatibility.
		if replaced == "\\S" {
			return "^.*\\S.*$"
		}
		replaced = fmt.Sprintf("^%s$", replaced)
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

	// WorkSpaces directory IDs: Smithy pattern has redundant end anchors ($)
	// inside each alternation branch AND at the end of the pattern.
	// Pattern "(foo$)|(bar$)$" is equivalent to "(foo)|(bar)" when the outer
	// anchors are present. Remove inner anchors for cleaner regex.
	workspacesDirectoryPatterns := map[string]string{
		"^(d-[0-9a-f]{8,63}$)|(wsd-[0-9a-z]{8,63}$)$": "(d-[0-9a-f]{8,63}$)|(wsd-[0-9a-z]{8,63}$)",
	}

	for _, transforms := range []map[string]string{arnRolePatterns, cognitoMessagePatterns, workspacesDirectoryPatterns} {
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
