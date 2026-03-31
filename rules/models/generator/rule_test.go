//go:build generators

package main

import (
	"slices"
	"testing"
)

func TestExtractPrefixDenies(t *testing.T) {
	for _, tc := range []struct {
		name         string
		input        string
		wantPrefixes []string
		wantCleaned  string
	}{
		{
			name:         "no lookahead",
			input:        `^[\p{L}]+$`,
			wantPrefixes: nil,
			wantCleaned:  `^[\p{L}]+$`,
		},
		{
			name:         "literal prefix aws:",
			input:        `^(?!aws:)[\p{L}]+$`,
			wantPrefixes: []string{"aws:"},
			wantCleaned:  `^[\p{L}]+$`,
		},
		{
			name:         "case insensitive aws variant",
			input:        `^(?!aws:).{1,128}$`,
			wantPrefixes: []string{"aws:"},
			wantCleaned:  `^.{1,128}$`,
		},
		{
			name:         "non-literal lookahead preserved",
			input:        `^(?!\s+$)[\p{L}]*$`,
			wantPrefixes: nil,
			wantCleaned:  `^(?!\s+$)[\p{L}]*$`,
		},
		{
			name:         "mixed literal and non-literal",
			input:        `^(?!aws:)(?!\s+$)[\p{L}]+$`,
			wantPrefixes: []string{"aws:"},
			wantCleaned:  `^(?!\s+$)[\p{L}]+$`,
		},
		{
			name:         "empty string",
			input:        "",
			wantPrefixes: nil,
			wantCleaned:  "",
		},
		{
			name:         "multiple literal prefixes",
			input:        `^(?!aws:)(?!pcs_)[a-zA-Z]+$`,
			wantPrefixes: []string{"aws:", "pcs_"},
			wantCleaned:  `^[a-zA-Z]+$`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gotPrefixes, gotCleaned := extractPrefixDenies(tc.input)
			if !slices.Equal(gotPrefixes, tc.wantPrefixes) {
				t.Errorf("prefixes = %v, want %v", gotPrefixes, tc.wantPrefixes)
			}
			if gotCleaned != tc.wantCleaned {
				t.Errorf("cleaned = %q, want %q", gotCleaned, tc.wantCleaned)
			}
		})
	}
}

func TestReplacePattern(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "unicode escape",
			input: `[\u000A]`,
			want:  `^[\x{000A}]$`,
		},
		{
			name:  "negative lookahead removal",
			input: `^(?!aws:)[\p{L}]+$`,
			want:  `^[\p{L}]+$`,
		},
		{
			name:  "no anchors",
			input: `[a-z]+`,
			want:  `^[a-z]+$`,
		},
		{
			name:  `special \S case`,
			input: `\S`,
			want:  `^.*\S.*$`,
		},
		{
			name:  "has $ missing ^",
			input: `[a-z]+$`,
			want:  `^[a-z]+$`,
		},
		{
			name:  "bare modifier",
			input: `^(?s)`,
			want:  `^(?s).*$`,
		},
		{
			name:  "prefix ending in colon",
			input: `^arn:`,
			want:  `^arn:.*$`,
		},
		{
			name:  "prefix ending in slash",
			input: `^s3://`,
			want:  `^s3://.*$`,
		},
		{
			name:  "quantifier ending *",
			input: `^[a-z]*`,
			want:  `^[a-z]*$`,
		},
		{
			name:  "quantifier ending +",
			input: `^[a-z]+`,
			want:  `^[a-z]+$`,
		},
		{
			name:  "quantifier ending }",
			input: `^[a-z]{1,128}`,
			want:  `^[a-z]{1,128}$`,
		},
		{
			name:  "ends with char class",
			input: `^[a-z]`,
			want:  `^[a-z].*$`,
		},
		{
			name:  "ends with letter",
			input: `^arn`,
			want:  `^arn.*$`,
		},
		{
			name:  "fully anchored no change",
			input: `^[a-z]+$`,
			want:  `^[a-z]+$`,
		},
		{
			name:  "IAM role ARN compatibility",
			input: `^arn:aws:iam::[0-9]*:role/`,
			want:  `^arn:aws:iam::[0-9]*:role/.*$`,
		},
		{
			name:  "Cognito message compatibility",
			input: `\{####\}`,
			want:  `^.*\{####\}.*$`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := replacePattern(tc.input); got != tc.want {
				t.Errorf("replacePattern(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
