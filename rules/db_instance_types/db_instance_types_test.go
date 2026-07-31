package db_instance_types

import "testing"

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
		name     string
		family   string
		expected bool
	}{
		{
			name:     "previous generation family",
			family:   "m1",
			expected: true,
		},
		{
			name:     "current generation family",
			family:   "m6g",
			expected: false,
		},
		{
			name:     "instance class rather than a family",
			family:   "db.m1.small",
			expected: false,
		},
		{
			name:     "empty",
			family:   "",
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if result := PreviousGeneration(tc.family); result != tc.expected {
				t.Errorf("expected %t, got %t", tc.expected, result)
			}
		})
	}
}
