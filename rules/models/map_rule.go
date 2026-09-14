package models

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// mapRule validates the keys and values of a map attribute, such as tags,
// against constraints extracted from AWS API models. Every generated map rule
// is an instance of this type; the instances live in map_rules.go.
type mapRule struct {
	tflint.DefaultRule

	name          string
	resourceType  string
	attributeName string
	// sensitive omits the offending key or value from pattern messages.
	sensitive bool

	itemsMax int

	keyMax        int
	keyMin        int
	keyPattern    *regexp.Regexp
	keyPrefixDeny []string

	valueMax        int
	valueMin        int
	valuePattern    *regexp.Regexp
	valuePrefixDeny []string
}

// Name returns the rule name
func (r *mapRule) Name() string {
	return r.name
}

// Enabled returns whether the rule is enabled by default
func (r *mapRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *mapRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *mapRule) Link() string {
	return ""
}

// Check validates map keys and values
func (r *mapRule) Check(runner tflint.Runner) error {
	logger.Trace("Check `%s` rule", r.Name())

	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		err := runner.EvaluateExpr(attribute.Expr, func(val map[string]string) error {
			for _, message := range r.violations(val) {
				if err := runner.EmitIssue(r, message, attribute.Expr.Range()); err != nil {
					return err
				}
			}
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// violations returns one message per constraint the map breaks, ordered by key.
func (r *mapRule) violations(val map[string]string) []string {
	var messages []string
	if r.itemsMax != 0 && len(val) > r.itemsMax {
		messages = append(messages, fmt.Sprintf("too many %s: %d exceeds the maximum of %d", r.attributeName, len(val), r.itemsMax))
	}
	for _, k := range slices.Sorted(maps.Keys(val)) {
		messages = append(messages, r.keyViolations(k)...)
		messages = append(messages, r.valueViolations(k, val[k])...)
	}
	return messages
}

func (r *mapRule) keyViolations(k string) []string {
	var messages []string
	length := utf8.RuneCountInString(k)
	if r.keyMax != 0 && length > r.keyMax {
		messages = append(messages, fmt.Sprintf("%s key %q must be %d characters or less", r.attributeName, truncateLongMessage(k), r.keyMax))
	}
	if r.keyMin != 0 && length < r.keyMin {
		messages = append(messages, fmt.Sprintf("%s key %q must be at least %d %s", r.attributeName, truncateLongMessage(k), r.keyMin, characters(r.keyMin)))
	}
	for _, prefix := range r.keyPrefixDeny {
		if strings.HasPrefix(k, prefix) {
			messages = append(messages, fmt.Sprintf("%s key %q must not start with %q", r.attributeName, truncateLongMessage(k), prefix))
		}
	}
	if r.keyPattern != nil && !r.keyPattern.MatchString(k) {
		if r.sensitive {
			messages = append(messages, fmt.Sprintf("%s key does not match valid pattern %s", r.attributeName, r.keyPattern))
		} else {
			messages = append(messages, fmt.Sprintf("%s key %q does not match valid pattern %s", r.attributeName, truncateLongMessage(k), r.keyPattern))
		}
	}
	return messages
}

func (r *mapRule) valueViolations(k, v string) []string {
	var messages []string
	length := utf8.RuneCountInString(v)
	if r.valueMax != 0 && length > r.valueMax {
		messages = append(messages, fmt.Sprintf("%s value for key %q must be %d characters or less", r.attributeName, truncateLongMessage(k), r.valueMax))
	}
	if r.valueMin != 0 && length < r.valueMin {
		messages = append(messages, fmt.Sprintf("%s value for key %q must be at least %d %s", r.attributeName, truncateLongMessage(k), r.valueMin, characters(r.valueMin)))
	}
	for _, prefix := range r.valuePrefixDeny {
		if strings.HasPrefix(v, prefix) {
			messages = append(messages, fmt.Sprintf("%s value for key %q must not start with %q", r.attributeName, truncateLongMessage(k), prefix))
		}
	}
	if r.valuePattern != nil && !r.valuePattern.MatchString(v) {
		if r.sensitive {
			messages = append(messages, fmt.Sprintf("%s value does not match valid pattern %s", r.attributeName, r.valuePattern))
		} else {
			messages = append(messages, fmt.Sprintf("%s value %q for key %q does not match valid pattern %s", r.attributeName, truncateLongMessage(v), truncateLongMessage(k), r.valuePattern))
		}
	}
	return messages
}

func characters(n int) string {
	if n == 1 {
		return "character"
	}
	return "characters"
}
