//go:build generators

package main

import (
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func TestStringsToCtyList(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  int
	}{
		{
			name:  "empty slice",
			input: []string{},
			want:  0,
		},
		{
			name:  "single value",
			input: []string{"value"},
			want:  1,
		},
		{
			name:  "multiple values",
			input: []string{"a", "b", "c"},
			want:  3,
		},
		{
			name:  "uppercase values",
			input: []string{"GZIP", "NONE"},
			want:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringsToCtyList(tt.input)

			if !result.Type().IsListType() {
				t.Errorf("stringsToCtyList() returned non-list type: %v", result.Type())
			}

			if result.LengthInt() != tt.want {
				t.Errorf("stringsToCtyList() length = %d, want %d", result.LengthInt(), tt.want)
			}

			// Verify all elements are strings and match input
			i := 0
			for it := result.ElementIterator(); it.Next(); {
				_, val := it.Element()
				if val.Type() != cty.String {
					t.Errorf("stringsToCtyList() element %d is not a string", i)
				}
				if val.AsString() != tt.input[i] {
					t.Errorf("stringsToCtyList() element %d = %q, want %q", i, val.AsString(), tt.input[i])
				}
				i++
			}
		})
	}
}

func TestMakeListTransformFunction_Uppercase(t *testing.T) {
	upperFunc := makeListTransformFunction(stdlib.UpperFunc)

	tests := []struct {
		name    string
		input   []string
		want    []string
		wantErr bool
	}{
		{
			name:  "lowercase to uppercase",
			input: []string{"gzip", "none"},
			want:  []string{"GZIP", "NONE"},
		},
		{
			name:  "mixed case",
			input: []string{"SsE-kMs", "sse-S3"},
			want:  []string{"SSE-KMS", "SSE-S3"},
		},
		{
			name:  "already uppercase",
			input: []string{"GZIP"},
			want:  []string{"GZIP"},
		},
		{
			name:  "empty list",
			input: []string{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputList := stringsToCtyList(tt.input)
			args := []cty.Value{inputList}

			result, err := upperFunc.Call(args)
			if (err != nil) != tt.wantErr {
				t.Errorf("upperFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if !result.Type().IsListType() {
					t.Errorf("upperFunc.Call() returned non-list type: %v", result.Type())
				}

				got := make([]string, 0, result.LengthInt())
				for it := result.ElementIterator(); it.Next(); {
					_, val := it.Element()
					got = append(got, val.AsString())
				}

				if len(got) != len(tt.want) {
					t.Errorf("upperFunc.Call() returned %d elements, want %d", len(got), len(tt.want))
				}

				for i, v := range got {
					if v != tt.want[i] {
						t.Errorf("upperFunc.Call() element %d = %q, want %q", i, v, tt.want[i])
					}
				}
			}
		})
	}
}

func TestMakeListTransformFunction_Replace(t *testing.T) {
	replaceFunc := makeListTransformFunction(stdlib.ReplaceFunc)

	tests := []struct {
		name    string
		input   []string
		old     string
		new     string
		want    []string
		wantErr bool
	}{
		{
			name:  "replace hyphens with underscores",
			input: []string{"sse-kms", "sse-s3"},
			old:   "-",
			new:   "_",
			want:  []string{"sse_kms", "sse_s3"},
		},
		{
			name:  "replace multiple occurrences",
			input: []string{"a-b-c"},
			old:   "-",
			new:   "_",
			want:  []string{"a_b_c"},
		},
		{
			name:  "no matches",
			input: []string{"none", "gzip"},
			old:   "-",
			new:   "_",
			want:  []string{"none", "gzip"},
		},
		{
			name:  "empty list",
			input: []string{},
			old:   "-",
			new:   "_",
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputList := stringsToCtyList(tt.input)
			args := []cty.Value{inputList, cty.StringVal(tt.old), cty.StringVal(tt.new)}

			result, err := replaceFunc.Call(args)
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if !result.Type().IsListType() {
					t.Errorf("replaceFunc.Call() returned non-list type: %v", result.Type())
				}

				got := make([]string, 0, result.LengthInt())
				for it := result.ElementIterator(); it.Next(); {
					_, val := it.Element()
					got = append(got, val.AsString())
				}

				if len(got) != len(tt.want) {
					t.Errorf("replaceFunc.Call() returned %d elements, want %d", len(got), len(tt.want))
				}

				for i, v := range got {
					if v != tt.want[i] {
						t.Errorf("replaceFunc.Call() element %d = %q, want %q", i, v, tt.want[i])
					}
				}
			}
		})
	}
}

func TestMakeListTransformFunction_ErrorCases(t *testing.T) {
	upperFunc := makeListTransformFunction(stdlib.UpperFunc)

	tests := []struct {
		name    string
		args    []cty.Value
		wantErr string
	}{
		{
			name:    "no arguments",
			args:    []cty.Value{},
			wantErr: "expected at least one argument",
		},
		{
			name:    "non-list argument",
			args:    []cty.Value{cty.StringVal("not a list")},
			wantErr: "first argument must be a list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := upperFunc.Call(tt.args)
			if err == nil {
				t.Errorf("upperFunc.Call() expected error containing %q, got nil", tt.wantErr)
				return
			}
			if err.Error() != tt.wantErr {
				t.Errorf("upperFunc.Call() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestBuildEvalContext(t *testing.T) {
	shapes := map[string]interface{}{
		"com.example#TagKey":   map[string]interface{}{"type": "string"},
		"com.example#TagValue": map[string]interface{}{"type": "string"},
	}
	ctx := buildEvalContext(shapes)

	if ctx == nil {
		t.Fatal("buildEvalContext() returned nil")
	}

	if ctx.Functions == nil {
		t.Error("buildEvalContext() Functions map is nil")
	}

	if ctx.Variables == nil {
		t.Error("buildEvalContext() Variables map is nil")
	}

	// Check that expected functions are registered
	expectedFuncs := []string{"uppercase", "replace", "listmap"}
	for _, name := range expectedFuncs {
		if _, ok := ctx.Functions[name]; !ok {
			t.Errorf("buildEvalContext() missing function %q", name)
		}
	}

	// Check that shape names are populated as shape objects
	for _, shapeName := range []string{"TagKey", "TagValue"} {
		val, ok := ctx.Variables[shapeName]
		if !ok {
			t.Errorf("buildEvalContext() missing variable %q", shapeName)
			continue
		}
		if !val.Type().Equals(shapeType) {
			t.Errorf("variable %q has type %s, want shape object", shapeName, val.Type().FriendlyName())
			continue
		}
		if val.GetAttr("name").AsString() != shapeName {
			t.Errorf("variable %q name = %q, want %q", shapeName, val.GetAttr("name").AsString(), shapeName)
		}
	}
}

func TestTransformComposition(t *testing.T) {
	// Test that transforms can be composed: uppercase(replace(x, "-", "_"))
	ctx := buildEvalContext(nil)

	input := []string{"sse-kms", "sse-s3"}
	ctx.Variables["EncryptionModeValue"] = stringsToCtyList(input)

	// This simulates the HCL expression: uppercase(replace(EncryptionModeValue, "-", "_"))
	replaceFunc := ctx.Functions["replace"]
	upperFunc := ctx.Functions["uppercase"]

	// First apply replace
	replaceResult, err := replaceFunc.Call([]cty.Value{
		ctx.Variables["EncryptionModeValue"],
		cty.StringVal("-"),
		cty.StringVal("_"),
	})
	if err != nil {
		t.Fatalf("replace() failed: %v", err)
	}

	// Then apply uppercase
	finalResult, err := upperFunc.Call([]cty.Value{replaceResult})
	if err != nil {
		t.Fatalf("uppercase() failed: %v", err)
	}

	want := []string{"SSE_KMS", "SSE_S3"}
	got := make([]string, 0, finalResult.LengthInt())
	for it := finalResult.ElementIterator(); it.Next(); {
		_, val := it.Element()
		got = append(got, val.AsString())
	}

	if len(got) != len(want) {
		t.Errorf("composition returned %d elements, want %d", len(got), len(want))
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("composition element %d = %q, want %q", i, v, want[i])
		}
	}
}

// Example test to demonstrate usage
func Example_stringsToCtyList() {
	values := []string{"gzip", "none"}
	result := stringsToCtyList(values)

	for it := result.ElementIterator(); it.Next(); {
		_, val := it.Element()
		fmt.Println(val.AsString())
	}
	// Output:
	// gzip
	// none
}
