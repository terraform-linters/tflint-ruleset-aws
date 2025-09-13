package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

type meta struct {
	Version string
}

type result struct {
	Issues []issue `json:"issues"`
	Errors []any   `json:"errors"`
}

type issue struct {
	Rule    any `json:"rule"`
	Message any `json:"message"`
	Range   any `json:"range"`
	Callers any `json:"callers"`

	// The following fields are ignored because they are added in TFLint v0.59.1.
	// We can uncomment this once the minimum supported version is v0.59.1+.
	// Fixed   any `json:"fixed"`
	// Fixable any `json:"fixable"`
}

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

			rawWant, err := readResultFile(testDir)
			if err != nil {
				t.Fatal(err)
			}

			var want result
			if err := json.Unmarshal(rawWant, &want); err != nil {
				t.Fatal(err)
			}

			var got result
			if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, want); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func readResultFile(dir string) ([]byte, error) {
	resultFile := "result.json"
	if runtime.GOOS == "windows" {
		if _, err := os.Stat(filepath.Join(dir, "result_windows.json")); !os.IsNotExist(err) {
			resultFile = "result_windows.json"
		}
	}
	if _, err := os.Stat(fmt.Sprintf("%s.tmpl", resultFile)); !os.IsNotExist(err) {
		resultFile = fmt.Sprintf("%s.tmpl", resultFile)
	}

	if !strings.HasSuffix(resultFile, ".tmpl") {
		return os.ReadFile(filepath.Join(dir, resultFile))
	}

	want := new(bytes.Buffer)
	tmpl := template.Must(template.ParseFiles(filepath.Join(dir, resultFile)))
	if err := tmpl.Execute(want, meta{Version: project.Version}); err != nil {
		return nil, err
	}
	return want.Bytes(), nil
}
