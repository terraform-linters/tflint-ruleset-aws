package lambda_deprecated_runtime

import (
	"testing"
	"time"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// python3.9 dates:
//   EOS:          2025-12-15
//   Block create: 2026-08-31
//   Block update: 2026-09-30
func Test_AwsLambdaFunctionDeprecatedRuntime(t *testing.T) {
	tf := `
resource "aws_lambda_function" "function" {
	function_name = "test_function"
	role = "test_role"
	runtime = "python3.9"
}
`
	issueRange := hcl.Range{
		Filename: "resource.tf",
		Start:    hcl.Pos{Line: 5, Column: 12},
		End:      hcl.Pos{Line: 5, Column: 23},
	}

	for _, tc := range []struct {
		name     string
		now      time.Time
		expected helper.Issues
	}{
		{
			name:     "before end of support",
			now:      time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{},
		},
		{
			name: "after end of support before block create",
			now:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{
				{
					Rule:    NewRule(),
					Message: `The "python3.9" runtime has reached end of support`,
					Range:   issueRange,
				},
			},
		},
		{
			name: "after block create before block update",
			now:  time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{
				{
					Rule:    NewRule(),
					Message: `The "python3.9" runtime has reached end of support and new function creation is blocked`,
					Range:   issueRange,
				},
			},
		},
		{
			name: "after block update",
			now:  time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
			expected: helper.Issues{
				{
					Rule:    NewRule(),
					Message: `The "python3.9" runtime has reached end of support and function updates are blocked`,
					Range:   issueRange,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rule := NewRule()
			rule.Now = tc.now

			runner := helper.TestRunner(t, map[string]string{"resource.tf": tf})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.expected, runner.Issues)
		})
	}
}
