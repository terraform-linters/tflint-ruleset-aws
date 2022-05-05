package integration

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIntegration(t *testing.T) {
	cases := []struct {
		Name    string
		Command *exec.Cmd
		Dir     string
	}{
		{
			Name:    "basic",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "basic",
		},
		{
			// Regression: https://github.com/terraform-linters/tflint-ruleset-aws/issues/41
			Name:    "rule config",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "rule-config",
		},
		{
			// Regression: https://github.com/terraform-linters/tflint/issues/1028
			Name:    "deep checking rule config",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "deep-checking-rule-config",
		},
		{
			// Regression: https://github.com/terraform-linters/tflint/issues/1029
			Name:    "heredoc",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "heredoc",
		},
		{
			// Regression: https://github.com/terraform-linters/tflint/issues/1054
			Name:    "disabled-rules",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "disabled-rules",
		},
		{
			// Regression: https://github.com/terraform-linters/tflint-ruleset-aws/issues/48
			Name:    "cty-based-eval",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "cty-based-eval",
		},
		{
			// Regression: https://github.com/terraform-linters/tflint/issues/1102
			Name:    "map-attribute",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "map-attribute",
		},
		{
			// Regression: https://github.com/terraform-linters/tflint/issues/1103
			Name:    "rule config with --enable-rule",
			Command: exec.Command("tflint", "--enable-rule", "aws_s3_bucket_name", "--format", "json", "--force"),
			Dir:     "rule-config",
		},
		{
			Name:    "rule config with --only",
			Command: exec.Command("tflint", "--only", "aws_s3_bucket_name", "--format", "json", "--force"),
			Dir:     "rule-config",
		},
	}

	dir, _ := os.Getwd()
	defer os.Chdir(dir)

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			testDir := filepath.Join(dir, tc.Dir)

			t.Cleanup(func() {
				if err := os.Chdir(dir); err != nil {
					t.Fatal(err)
				}
			})

			if err := os.Chdir(testDir); err != nil {
				t.Fatal(err)
			}

			var stdout, stderr bytes.Buffer
			tc.Command.Stdout = &stdout
			tc.Command.Stderr = &stderr
			if err := tc.Command.Run(); err != nil {
				t.Fatalf("%s, stdout=%s stderr=%s", err, stdout.String(), stderr.String())
			}

			var b []byte
			var err error
			if runtime.GOOS == "windows" && IsWindowsResultExist() {
				b, err = os.ReadFile(filepath.Join(testDir, "result_windows.json"))
			} else {
				b, err = os.ReadFile(filepath.Join(testDir, "result.json"))
			}
			if err != nil {
				t.Fatal(err)
			}

			var expected interface{}
			if err := json.Unmarshal(b, &expected); err != nil {
				t.Fatal(err)
			}

			var got interface{}
			if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, expected); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func IsWindowsResultExist() bool {
	_, err := os.Stat("result_windows.json")
	return !os.IsNotExist(err)
}
