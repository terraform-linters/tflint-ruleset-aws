package lambda_deprecated_runtime

import (
	"testing"
	"time"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// python3.9 dates (deprecated, all confirmed as of 2026-03-29):
//
//	EOS:          2025-12-15
//	Block create: 2026-08-31
//	Block update: 2026-09-30
//
// nodejs22.x dates (supported, all speculative as of 2026-03-29):
//
//	EOS:          2027-04-30
//	Block create: 2027-06-01
//	Block update: 2027-07-01
func Test_AwsLambdaFunctionDeprecatedRuntime(t *testing.T) {
	for _, tc := range []struct {
		name     string
		runtime  string
		now      time.Time
		expected helper.Issues
	}{
		{
			name:     "before end of support",
			runtime:  "python3.9",
			now:      time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{},
		},
		{
			name:    "confirmed end of support",
			runtime: "python3.9",
			now:     time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{
				{
					Rule:    NewRule(),
					Message: `The "python3.9" runtime has reached end of support`,
				},
			},
		},
		{
			name:    "confirmed block create",
			runtime: "python3.9",
			now:     time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{
				{
					Rule:    NewRule(),
					Message: `The "python3.9" runtime has reached end of support and new function creation was scheduled to be blocked on Aug 31, 2026 (as of Mar 29, 2026)`,
				},
			},
		},
		{
			name:    "confirmed block update",
			runtime: "python3.9",
			now:     time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{
				{
					Rule:    NewRule(),
					Message: `The "python3.9" runtime has reached end of support and function updates were scheduled to be blocked on Sep 30, 2026 (as of Mar 29, 2026)`,
				},
			},
		},
		{
			name:    "speculative end of support",
			runtime: "nodejs22.x",
			now:     time.Date(2027, time.May, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{
				{
					Rule:    NewRule(),
					Message: `The "nodejs22.x" runtime was scheduled to reach end of support on Apr 30, 2027 (as of Mar 29, 2026)`,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rule := NewRule()
			rule.Now = tc.now

			tf := `
resource "aws_lambda_function" "function" {
	function_name = "test_function"
	role = "test_role"
	runtime = "` + tc.runtime + `"
}
`
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tf})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			// Only check Rule and Message, not Range (varies by runtime string length)
			if len(runner.Issues) != len(tc.expected) {
				t.Fatalf("expected %d issues, got %d", len(tc.expected), len(runner.Issues))
			}
			for i, issue := range runner.Issues {
				if issue.Message != tc.expected[i].Message {
					t.Errorf("expected message %q, got %q", tc.expected[i].Message, issue.Message)
				}
			}
		})
	}
}

func Test_RuntimeMessage(t *testing.T) {
	updatedAt := time.Date(2026, time.March, 29, 0, 0, 0, 0, time.UTC)
	blockCreate := time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC)
	blockUpdate := time.Date(2026, time.September, 30, 0, 0, 0, 0, time.UTC)

	rt := runtimeLifecycle{
		EndOfSupportDate: time.Date(2025, time.December, 15, 0, 0, 0, 0, time.UTC),
		BlockCreateDate:  &blockCreate,
		BlockUpdateDate:  &blockUpdate,
	}

	for _, tc := range []struct {
		name     string
		now      time.Time
		expected string
	}{
		{
			name:     "confirmed EOS only",
			now:      time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: `The "test" runtime has reached end of support`,
		},
		{
			name:     "speculative block create",
			now:      time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
			expected: `The "test" runtime has reached end of support and new function creation was scheduled to be blocked on Aug 31, 2026 (as of Mar 29, 2026)`,
		},
		{
			name:     "speculative block update",
			now:      time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
			expected: `The "test" runtime has reached end of support and function updates were scheduled to be blocked on Sep 30, 2026 (as of Mar 29, 2026)`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := runtimeMessage("test", rt, updatedAt, tc.now); got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}

	// With updatedAt after all dates, everything is confirmed
	futureUpdatedAt := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	now := time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC)
	got := runtimeMessage("test", rt, futureUpdatedAt, now)
	expected := `The "test" runtime has reached end of support and function updates are blocked`
	if got != expected {
		t.Errorf("confirmed block update: expected %q, got %q", expected, got)
	}
}
