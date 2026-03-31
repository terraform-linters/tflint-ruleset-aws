//go:build generators

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/genutils"
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
	RuleName        string
	RuleNameCC      string
	ResourceType    string
	AttributeName   string
	Sensitive       bool
	ItemsMax        int
	KeyMax          int
	KeyMin          int
	KeyPattern      string
	KeyPrefixDeny   []string
	ValueMax        int
	ValueMin        int
	ValuePattern    string
	ValuePrefixDeny []string
}

// sentinelMax is the Smithy unbounded length sentinel (math.MaxInt32).
// Constraints using this value are effectively "no limit" and should be omitted.
const sentinelMax = 2147483647

var (
	unicodeEscapeRegex     = regexp.MustCompile(`\\u([0-9A-F]{4})`)
	negativeLookaheadRegex = regexp.MustCompile(`\(\?![^)]+\)`)
	literalPrefixRegex     = regexp.MustCompile(`^[a-zA-Z0-9_:/.+\-]+$`)
	matchAllRegex          = regexp.MustCompile(`^\^(\(\?[simU]+\))?\.\*\$$`)
)

func buildRuleMeta(resource, attribute string, model map[string]interface{}, schema *tfjson.SchemaAttribute) *ruleMeta {
	ruleName := makeRuleName(resource, attribute)
	meta := &ruleMeta{
		RuleName:      ruleName,
		RuleNameCC:    genutils.ToCamel(ruleName),
		ResourceType:  resource,
		AttributeName: attribute,
		Sensitive:     schema != nil && schema.Sensitive,
		Max:           clampSentinel(fetchNumber(model, "max")),
		Min:           fetchNumber(model, "min"),
		Pattern:       stripMatchAll(replacePattern(fetchString(model, "pattern"))),
		Enum:          fetchStrings(model, "enum"),
	}

	if meta.Pattern != "" {
		if _, err := regexp.Compile(meta.Pattern); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s.%s has malformed pattern %q: %v\n", resource, attribute, meta.Pattern, err)
			os.Exit(1)
		}
	}

	return meta
}

func generateRuleFile(resource, attribute string, model map[string]interface{}, schema *tfjson.SchemaAttribute) {
	meta := buildRuleMeta(resource, attribute, model, schema)
	genutils.GenerateFile(fmt.Sprintf("%s.go", meta.RuleName), "pattern_rule.go.tmpl", meta)
}

func generateRuleTestFile(resource, attribute string, model map[string]interface{}, test test) {
	meta := buildRuleMeta(resource, attribute, model, nil)
	meta.TestOK = formatTest(test.OK)
	meta.TestNG = formatTest(test.NG)
	genutils.GenerateFile(fmt.Sprintf("%s_test.go", meta.RuleName), "pattern_rule_test.go.tmpl", meta)
}

func generateMapRuleFile(resource, attribute string, listModel, keyModel, valueModel map[string]interface{}, schema *tfjson.SchemaAttribute) bool {
	ruleName := makeRuleName(resource, attribute)

	var itemsMax int
	if listModel != nil {
		itemsMax = fetchNumber(listModel, "max")
	}
	keyPrefixDeny, rawKeyPattern := extractPrefixDenies(fetchString(keyModel, "pattern"))
	valuePrefixDeny, rawValuePattern := extractPrefixDenies(fetchString(valueModel, "pattern"))
	keyPattern := stripMatchAll(replacePattern(rawKeyPattern))
	valuePattern := stripMatchAll(replacePattern(rawValuePattern))
	keyMax := clampSentinel(fetchNumber(keyModel, "max"))
	keyMin := fetchNumber(keyModel, "min")
	valueMax := clampSentinel(fetchNumber(valueModel, "max"))
	valueMin := fetchNumber(valueModel, "min")

	// Skip if no constraints exist
	hasItemsConstraints := itemsMax != 0
	hasKeyConstraints := keyPattern != "" || keyMax != 0 || keyMin != 0 || len(keyPrefixDeny) > 0
	hasValueConstraints := valuePattern != "" || valueMax != 0 || valueMin != 0 || len(valuePrefixDeny) > 0
	if !hasItemsConstraints && !hasKeyConstraints && !hasValueConstraints {
		return false
	}

	// Smithy models may use regex features incompatible with Go's regexp (e.g.,
	// surrogate pair ranges). Warn and skip rather than failing fatally.
	if keyPattern != "" {
		if _, err := regexp.Compile(keyPattern); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping %s.%s: incompatible key pattern %q: %v\n", resource, attribute, keyPattern, err)
			keyPattern = ""
		}
	}
	if valuePattern != "" {
		if _, err := regexp.Compile(valuePattern); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping %s.%s: incompatible value pattern %q: %v\n", resource, attribute, valuePattern, err)
			valuePattern = ""
		}
	}

	meta := &mapRuleMeta{
		RuleName:        ruleName,
		RuleNameCC:      genutils.ToCamel(ruleName),
		ResourceType:    resource,
		AttributeName:   attribute,
		Sensitive:       schema != nil && schema.Sensitive,
		ItemsMax:        itemsMax,
		KeyMax:          keyMax,
		KeyMin:          keyMin,
		KeyPattern:      keyPattern,
		KeyPrefixDeny:   keyPrefixDeny,
		ValueMax:        valueMax,
		ValueMin:        valueMin,
		ValuePattern:    valuePattern,
		ValuePrefixDeny: valuePrefixDeny,
	}

	genutils.GenerateFile(fmt.Sprintf("%s.go", ruleName), "map_rule.go.tmpl", meta)
	return true
}

// generateMapRuleFromShapes generates a map validation rule from a resolved map result.
// The result contains three nested shape objects: items, key, and value.
func generateMapRuleFromShapes(resource, attribute string, result cty.Value, schema *tfjson.SchemaAttribute) bool {
	return generateMapRuleFile(resource, attribute,
		shapeToModel(result.GetAttr("items")),
		shapeToModel(result.GetAttr("key")),
		shapeToModel(result.GetAttr("value")),
		schema,
	)
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

// extractPrefixDenies finds negative lookahead patterns containing literal strings
// (e.g., (?!aws:)) and returns them as prefix deny values. Non-literal lookaheads
// (containing regex metacharacters) are left in the pattern for replacePattern to strip.
func extractPrefixDenies(pattern string) (prefixes []string, cleaned string) {
	cleaned = negativeLookaheadRegex.ReplaceAllStringFunc(pattern, func(match string) string {
		// Extract the content between (?! and )
		inner := match[3 : len(match)-1]
		if literalPrefixRegex.MatchString(inner) {
			prefixes = append(prefixes, inner)
			return ""
		}
		return match
	})
	return prefixes, cleaned
}

func replacePattern(pattern string) string {
	if pattern == "" {
		return pattern
	}

	// Convert Unicode escapes from \uXXXX to \x{XXXX} format
	replaced := unicodeEscapeRegex.ReplaceAllString(pattern, `\x{$1}`)

	// Remove remaining negative lookahead patterns that Go's regexp doesn't support.
	// Literal prefix lookaheads (e.g., (?!aws:)) should be extracted by
	// extractPrefixDenies before calling replacePattern.
	replaced = negativeLookaheadRegex.ReplaceAllString(replaced, "")

	// Handle patterns missing anchors
	// The Ruby SDK generator ensured all patterns had both ^ and $ anchors.
	// Smithy models often have incomplete patterns (e.g., "^(?s)" or "[a-z]*").
	// We add missing anchors to maintain backward compatibility.
	hasPrefix := strings.HasPrefix(replaced, "^") || strings.HasPrefix(replaced, "(^")
	hasSuffix := strings.HasSuffix(replaced, "$") || strings.HasSuffix(replaced, "$)")

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

	// Strip redundant nested anchors: ^(^...$)$ → ^...$, ^(^...$|^...$)$ → ^...$|^...$
	replaced = stripRedundantAnchors(replaced)

	// Apply compatibility transforms to maintain backward compatibility with Ruby SDK
	if transformed := applyCompatibilityTransforms(replaced); transformed != "" {
		return transformed
	}

	return replaced
}

// stripRedundantAnchors removes outer ^(...)$ wrapping when the opening paren
// matches the closing paren (i.e., the entire pattern is one group) and the
// inner content is already anchored.
func stripRedundantAnchors(pattern string) string {
	if !strings.HasPrefix(pattern, "^(") || !strings.HasSuffix(pattern, ")$") {
		return pattern
	}
	// Check if the ( at position 1 matches the ) at the end
	depth := 0
	for i := 1; i < len(pattern)-1; i++ {
		switch pattern[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i < len(pattern)-2 {
				// The opening ( closed before the final ), so the outer
				// parens are not a single wrapping group (e.g., (a)|(b))
				return pattern
			}
		}
	}
	inner := pattern[2 : len(pattern)-2]
	if strings.HasPrefix(inner, "^") {
		return inner
	}
	return pattern
}

// compatibilityTransforms maps patterns that need special handling for backward
// compatibility with the Ruby SDK or Terraform provider expectations.
var compatibilityTransforms = map[string]string{
	// IAM role ARN patterns: Smithy models end at the role prefix (e.g., "role/")
	// but Terraform configurations include role paths (e.g., "role/my-app/my-role").
	"^arn:aws:iam::[0-9]*:role/":                 "^arn:aws:iam::[0-9]*:role/.*$",
	"^arn:aws:iam::\\d{12}:role/":                "^arn:aws:iam::\\d{12}:role/.*$",
	"^arn:aws(-[\\w]+)*:iam::[0-9]{12}:role/":    "^arn:aws(-[\\w]+)*:iam::[0-9]{12}:role/.*$",
	"^arn:aws(-[a-z]{1,3}){0,2}:iam::\\d+:role/": "^arn:aws(-[a-z]{1,3}){0,2}:iam::\\d+:role/.*$",
	// Cognito SMS authentication messages: The {####} placeholder must appear
	// within a message, not be the entire message.
	"^\\{####\\}$": "^.*\\{####\\}.*$",
}

func applyCompatibilityTransforms(pattern string) string {
	if transformed, ok := compatibilityTransforms[pattern]; ok {
		return transformed
	}
	return ""
}

// stripMatchAll returns "" for patterns that match all strings (e.g., ^.*$, ^(?s).*$).
func stripMatchAll(pattern string) string {
	if matchAllRegex.MatchString(pattern) {
		return ""
	}
	return pattern
}

// clampSentinel returns 0 for the Smithy unbounded sentinel value.
func clampSentinel(n int) int {
	if n == sentinelMax {
		return 0
	}
	return n
}

func formatTest(body string) string {
	if strings.Contains(body, "\n") {
		return fmt.Sprintf("<<TEXT\n%sTEXT", body)
	}
	return fmt.Sprintf(`"%s"`, body)
}
