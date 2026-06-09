package lambda_deprecated_runtime

import (
	"testing"
	"time"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaFunctionDeprecatedRuntime(t *testing.T) {
	// Synthetic data so the test does not depend on the generated
	// deprecated_runtimes.json, whose dates shift whenever AWS updates
	// its documentation.
	updatedAt := time.Date(2026, time.March, 29, 0, 0, 0, 0, time.UTC)
	blockCreate := time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC)
	blockUpdate := time.Date(2026, time.September, 30, 0, 0, 0, 0, time.UTC)

	testRuntimes := runtimesData{
		UpdatedAt: updatedAt,
		Runtimes: map[string]runtimeLifecycle{
			"python3.9": {
				EndOfSupportDate: time.Date(2025, time.December, 15, 0, 0, 0, 0, time.UTC),
				BlockCreateDate:  &blockCreate,
				BlockUpdateDate:  &blockUpdate,
			},
			"nodejs22.x": {
				EndOfSupportDate: time.Date(2027, time.April, 30, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	stale := updatedAt.Format("Jan 2, 2006")

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
			name:    "speculative block create",
			runtime: "python3.9",
			now:     time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{
				{
					Rule:    NewRule(),
					Message: `The "python3.9" runtime has reached end of support and new function creation was scheduled to be blocked on Aug 31, 2026 (as of ` + stale + `)`,
				},
			},
		},
		{
			name:    "speculative block update",
			runtime: "python3.9",
			now:     time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{
				{
					Rule:    NewRule(),
					Message: `The "python3.9" runtime has reached end of support and function updates were scheduled to be blocked on Sep 30, 2026 (as of ` + stale + `)`,
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
					Message: `The "nodejs22.x" runtime was scheduled to reach end of support on Apr 30, 2027 (as of ` + stale + `)`,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rule := NewRule()
			rule.Now = tc.now
			rule.runtimes = testRuntimes

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

// Exercises the rule against the embedded generated data. Asserts only that
// a long-deprecated runtime is flagged, since the generated dates shift
// whenever AWS updates its documentation.
func Test_AwsLambdaFunctionDeprecatedRuntime_EmbeddedData(t *testing.T) {
	tf := `
resource "aws_lambda_function" "function" {
	function_name = "test_function"
	role = "test_role"
	runtime = "python2.7"
}
`
	runner := helper.TestRunner(t, map[string]string{"resource.tf": tf})
	rule := NewRule()
	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}
	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}
}

func Test_RuntimeMessage(t *testing.T) {
	updatedAt := runtimes.UpdatedAt
	stale := updatedAt.Format("Jan 2, 2006")
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
			expected: `The "test" runtime has reached end of support and new function creation was scheduled to be blocked on Aug 31, 2026 (as of ` + stale + `)`,
		},
		{
			name:     "speculative block update",
			now:      time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
			expected: `The "test" runtime has reached end of support and function updates were scheduled to be blocked on Sep 30, 2026 (as of ` + stale + `)`,
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
