package db_instance_types

import (
	"testing"

	"github.com/terraform-linters/tflint-ruleset-aws/rules/stringset"
)

func TestValid(t *testing.T) {
	for _, tc := range []struct {
		name          string
		instanceClass string
		expected      bool
	}{
		{
			name:          "current generation class",
			instanceClass: "db.m6g.large",
			expected:      true,
		},
		{
			name:          "previous generation class",
			instanceClass: "db.m1.small",
			expected:      true,
		},
		{
			name:          "class with a tenancy and memory suffix",
			instanceClass: "db.r5.2xlarge.tpc2.mem8x",
			expected:      true,
		},
		{
			name:          "EC2 instance type without the db prefix",
			instanceClass: "m4.2xlarge",
			expected:      false,
		},
		{
			name:          "unknown size",
			instanceClass: "db.m6g.nano",
			expected:      false,
		},
		{
			name:          "empty",
			instanceClass: "",
			expected:      false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if result := Valid(tc.instanceClass); result != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, result)
			}
		})
	}
}

func TestPreviousGeneration(t *testing.T) {
	for _, tc := range []struct {
		name          string
		instanceClass string
		expected      bool
	}{
		{
			name:          "previous generation class",
			instanceClass: "db.m1.small",
			expected:      true,
		},
		{
			name:          "burstable previous generation class",
			instanceClass: "db.t2.micro",
			expected:      true,
		},
		{
			name:          "current generation class",
			instanceClass: "db.t4g.micro",
			expected:      false,
		},
		{
			name:          "unknown size in a previous generation family",
			instanceClass: "db.m1.unknownsize",
			expected:      false,
		},
		{
			name:          "previous generation family under another prefix",
			instanceClass: "foo.m1.large",
			expected:      false,
		},
		{
			name:          "family rather than an instance class",
			instanceClass: "m1",
			expected:      false,
		},
		{
			name:          "empty",
			instanceClass: "",
			expected:      false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if result := PreviousGeneration(tc.instanceClass); result != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, result)
			}
		})
	}
}

// TestEmbeddedInstanceClasses guards the read path. The generator's floors only
// run when someone regenerates, so a truncated or partially written JSON file
// would otherwise ship a plugin that rejects every instance class.
func TestEmbeddedInstanceClasses(t *testing.T) {
	for _, tc := range []struct {
		name    string
		classes stringset.Set
		minimum int
	}{
		{
			name:    "instance classes",
			classes: instanceClasses.All,
			minimum: 400,
		},
		{
			name:    "previous generation classes",
			classes: instanceClasses.PreviousGeneration,
			minimum: 30,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.classes) < tc.minimum {
				t.Errorf("expected at least %d, got %d", tc.minimum, len(tc.classes))
			}
		})
	}
}
