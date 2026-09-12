package models

import (
	"regexp"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// mapRuleByName returns the generated map rule with the given name.
func mapRuleByName(t *testing.T, name string) tflint.Rule {
	t.Helper()
	for _, rule := range mapRules {
		if rule.Name() == name {
			return rule
		}
	}
	t.Fatalf("no generated map rule named %q", name)
	return nil
}

func Test_mapRule_violations(t *testing.T) {
	rule := &mapRule{
		name:            "test_rule",
		resourceType:    "aws_test",
		attributeName:   "labels",
		itemsMax:        2,
		keyMax:          5,
		keyMin:          2,
		keyPattern:      regexp.MustCompile(`^[a-z]+$`),
		keyPrefixDeny:   []string{"aws:", "k8s:"},
		valueMax:        5,
		valueMin:        2,
		valuePattern:    regexp.MustCompile(`^[a-z]*$`),
		valuePrefixDeny: []string{"sys/"},
	}

	for _, tc := range []struct {
		name     string
		val      map[string]string
		expected []string
	}{
		{
			name: "valid",
			val:  map[string]string{"abc": "def"},
		},
		{
			name: "too many items",
			val:  map[string]string{"aa": "aa", "bb": "bb", "cc": "cc"},
			expected: []string{
				"too many labels: 3 exceeds the maximum of 2",
			},
		},
		{
			name: "key constraints in order",
			val:  map[string]string{"aws:Name1": "ok"},
			expected: []string{
				`labels key "aws:Name1" must be 5 characters or less`,
				`labels key "aws:Name1" must not start with "aws:"`,
				`labels key "aws:Name1" does not match valid pattern ^[a-z]+$`,
			},
		},
		{
			name: "short key emits min length and pattern",
			val:  map[string]string{"": "ok"},
			expected: []string{
				`labels key "" must be at least 2 characters`,
				`labels key "" does not match valid pattern ^[a-z]+$`,
			},
		},
		{
			name: "value constraints in order",
			val:  map[string]string{"key": "sys/ABC"},
			expected: []string{
				`labels value for key "key" must be 5 characters or less`,
				`labels value for key "key" must not start with "sys/"`,
				`labels value "sys/ABC" for key "key" does not match valid pattern ^[a-z]*$`,
			},
		},
		{
			name: "short value",
			val:  map[string]string{"key": "a"},
			expected: []string{
				`labels value for key "key" must be at least 2 characters`,
			},
		},
		{
			name: "keys are reported in sorted order",
			val:  map[string]string{"b!": "ok", "a!": "ok"},
			expected: []string{
				`labels key "a!" does not match valid pattern ^[a-z]+$`,
				`labels key "b!" does not match valid pattern ^[a-z]+$`,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.violations(tc.val)
			if len(got) != len(tc.expected) {
				t.Fatalf("violations = %q, want %q", got, tc.expected)
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("violations[%d] = %q, want %q", i, got[i], tc.expected[i])
				}
			}
		})
	}
}

func Test_mapRule_lengthCountsCharacters(t *testing.T) {
	rule := &mapRule{name: "test_rule", attributeName: "tags", keyMax: 5, valueMax: 5}

	if got := rule.violations(map[string]string{"ééééé": "ééééé"}); len(got) != 0 {
		t.Errorf("violations = %q, want none", got)
	}
	got := rule.violations(map[string]string{"éééééé": "éééééé"})
	expected := []string{
		`tags key "éééééé" must be 5 characters or less`,
		`tags value for key "éééééé" must be 5 characters or less`,
	}
	if len(got) != len(expected) {
		t.Fatalf("violations = %q, want %q", got, expected)
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("violations[%d] = %q, want %q", i, got[i], expected[i])
		}
	}
}

func Test_mapRule_sensitive(t *testing.T) {
	rule := &mapRule{
		name:          "test_rule",
		attributeName: "environment_variables",
		sensitive:     true,
		keyPattern:    regexp.MustCompile(`^[A-Z]+$`),
		valuePattern:  regexp.MustCompile(`^[a-z]+$`),
	}

	got := rule.violations(map[string]string{"secret-key": "Secret Value"})
	expected := []string{
		"environment_variables key does not match valid pattern ^[A-Z]+$",
		"environment_variables value does not match valid pattern ^[a-z]+$",
	}
	if len(got) != len(expected) {
		t.Fatalf("violations = %q, want %q", got, expected)
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("violations[%d] = %q, want %q", i, got[i], expected[i])
		}
	}
}

func Test_mapRule_singleCharacterMessage(t *testing.T) {
	rule := &mapRule{name: "test_rule", attributeName: "tags", keyMin: 1, valueMin: 1}

	got := rule.violations(map[string]string{"": ""})
	expected := []string{
		`tags key "" must be at least 1 character`,
		`tags value for key "" must be at least 1 character`,
	}
	if len(got) != len(expected) {
		t.Fatalf("violations = %q, want %q", got, expected)
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("violations[%d] = %q, want %q", i, got[i], expected[i])
		}
	}
}

func Test_mapRules_registered(t *testing.T) {
	if len(mapRules) == 0 {
		t.Fatal("no generated map rules")
	}
	for _, rule := range mapRules {
		found := false
		for _, registered := range Rules {
			if registered == rule {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("map rule %q is not registered in Rules", rule.Name())
		}
	}
	// Check that the shared type works end to end through the runner.
	runner := helper.TestRunner(t, map[string]string{"resource.tf": `
resource "aws_ecs_service" "foo" {
	tags = {
		Name = "example"
	}
}`})
	if err := mapRuleByName(t, "aws_ecs_service_invalid_tags").Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	helper.AssertIssues(t, helper.Issues{}, runner.Issues)
}
