package aws

import (
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-ruleset-aws/rules"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type testRule struct {
	tflint.DefaultRule
	name string
}

func (r *testRule) Name() string              { return r.name }
func (r *testRule) Enabled() bool             { return true }
func (r *testRule) Severity() tflint.Severity { return tflint.ERROR }
func (r *testRule) Check(tflint.Runner) error { return nil }

type awsDbInstanceInvalidTypeRule struct {
	testRule
}

type awsResourceMissingTagRule struct {
	testRule
}

type awsS3BucketNameRule struct {
	testRule
}

func TestApplyConfig(t *testing.T) {
	mustParseExpr := func(input string) hcl.Expression {
		expr, diags := hclsyntax.ParseExpression([]byte(input), "", hcl.InitialPos)
		if diags.HasErrors() {
			panic(diags)
		}
		return expr
	}

	tests := []struct {
		name   string
		global *tflint.Config
		config *hclext.BodyContent
		want   []string
	}{
		{
			name:   "default",
			global: &tflint.Config{},
			config: &hclext.BodyContent{},
			want: []string{
				"aws_db_instance_invalid_type",
				"aws_resource_missing_tag_rule",
				"aws_s3_bucket_name",
			},
		},
		{
			name:   "disabled by default",
			global: &tflint.Config{DisabledByDefault: true},
			config: &hclext.BodyContent{},
			want:   []string{},
		},
		{
			name:   "preset",
			global: &tflint.Config{},
			config: &hclext.BodyContent{
				Attributes: hclext.Attributes{
					"preset": &hclext.Attribute{Name: "preset", Expr: mustParseExpr(`"recommended"`)},
				},
			},
			want: []string{
				"aws_db_instance_invalid_type",
				"aws_resource_missing_tag_rule",
			},
		},
		{
			name: "rule config",
			global: &tflint.Config{
				Rules: map[string]*tflint.RuleConfig{
					"aws_db_instance_invalid_type": {
						Name:    "aws_db_instance_invalid_type",
						Enabled: false,
					},
				},
			},
			config: &hclext.BodyContent{},
			want: []string{
				"aws_resource_missing_tag_rule",
				"aws_s3_bucket_name",
			},
		},
		{
			name:   "only",
			global: &tflint.Config{Only: []string{"aws_s3_bucket_name"}},
			config: &hclext.BodyContent{},
			want: []string{
				"aws_s3_bucket_name",
			},
		},
		{
			name:   "disabled by default + preset",
			global: &tflint.Config{DisabledByDefault: true},
			config: &hclext.BodyContent{
				Attributes: hclext.Attributes{
					"preset": &hclext.Attribute{Name: "preset", Expr: mustParseExpr(`"recommended"`)},
				},
			},
			want: []string{
				"aws_db_instance_invalid_type",
				"aws_resource_missing_tag_rule",
			},
		},
		{
			name: "disabled by default + rule config",
			global: &tflint.Config{
				Rules: map[string]*tflint.RuleConfig{
					"aws_db_instance_invalid_type": {
						Name:    "aws_db_instance_invalid_type",
						Enabled: true,
					},
				},
				DisabledByDefault: true,
			},
			config: &hclext.BodyContent{},
			want: []string{
				"aws_db_instance_invalid_type",
			},
		},
		{
			name: "disabled by default + only",
			global: &tflint.Config{
				DisabledByDefault: true,
				Only:              []string{"aws_db_instance_invalid_type"},
			},
			config: &hclext.BodyContent{},
			want: []string{
				"aws_db_instance_invalid_type",
			},
		},
		{
			name: "preset + rule config",
			global: &tflint.Config{
				Rules: map[string]*tflint.RuleConfig{
					"aws_db_instance_invalid_type": {
						Name:    "aws_db_instance_invalid_type",
						Enabled: false,
					},
				},
			},
			config: &hclext.BodyContent{
				Attributes: hclext.Attributes{
					"preset": &hclext.Attribute{Name: "preset", Expr: mustParseExpr(`"recommended"`)},
				},
			},
			want: []string{
				"aws_resource_missing_tag_rule",
			},
		},
		{
			name: "preset + only",
			global: &tflint.Config{
				Only: []string{"aws_s3_bucket_name"},
			},
			config: &hclext.BodyContent{
				Attributes: hclext.Attributes{
					"preset": &hclext.Attribute{Name: "preset", Expr: mustParseExpr(`"recommended"`)},
				},
			},
			want: []string{
				"aws_s3_bucket_name",
			},
		},
		{
			name: "rule config + only",
			global: &tflint.Config{
				Rules: map[string]*tflint.RuleConfig{
					"aws_db_instance_invalid_type": {
						Name:    "aws_db_instance_invalid_type",
						Enabled: false,
					},
				},
				Only: []string{"aws_db_instance_invalid_type", "aws_s3_bucket_name"},
			},
			config: &hclext.BodyContent{},
			want: []string{
				"aws_db_instance_invalid_type",
				"aws_s3_bucket_name",
			},
		},
		{
			name: "disabled by default + preset + rule config",
			global: &tflint.Config{
				Rules: map[string]*tflint.RuleConfig{
					"aws_db_instance_invalid_type": {
						Name:    "aws_db_instance_invalid_type",
						Enabled: false,
					},
				},
				DisabledByDefault: true,
			},
			config: &hclext.BodyContent{
				Attributes: hclext.Attributes{
					"preset": &hclext.Attribute{Name: "preset", Expr: mustParseExpr(`"recommended"`)},
				},
			},
			want: []string{
				"aws_resource_missing_tag_rule",
			},
		},
		{
			name: "disabled by default + preset + rule config + only",
			global: &tflint.Config{
				Rules: map[string]*tflint.RuleConfig{
					"aws_db_instance_invalid_type": {
						Name:    "aws_db_instance_invalid_type",
						Enabled: false,
					},
				},
				DisabledByDefault: true,
				Only:              []string{"aws_db_instance_invalid_type", "aws_s3_bucket_name"},
			},
			config: &hclext.BodyContent{
				Attributes: hclext.Attributes{
					"preset": &hclext.Attribute{Name: "preset", Expr: mustParseExpr(`"recommended"`)},
				},
			},
			want: []string{
				"aws_db_instance_invalid_type",
				"aws_s3_bucket_name",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		BuiltinRuleSet: tflint.BuiltinRuleSet{
			Name:    "aws",
			Version: project.Version,
			Rules:   rules.Rules,
		},
			ruleset := &RuleSet{
				PresetRules: map[string][]tflint.Rule{
					"all": {
						&awsS3BucketNameRule{testRule: testRule{name: "aws_s3_bucket_name"}},
						&awsDbInstanceInvalidTypeRule{testRule: testRule{name: "aws_db_instance_invalid_type"}},
						&awsResourceMissingTagRule{testRule: testRule{name: "aws_resource_missing_tag_rule"}},
					},
					"cis": {
						&awsDbInstanceInvalidTypeRule{testRule: testRule{name: "aws_db_instance_invalid_type"}},
						&awsResourceMissingTagRule{testRule: testRule{name: "aws_resource_missing_tag_rule"}},
					},
				},
				config: &Config{},
			}

			err := ruleset.ApplyGlobalConfig(test.global)
			if err != nil {
				t.Fatal(err)
			}

			err = ruleset.ApplyConfig(test.config)
			if err != nil {
				t.Fatal(err)
			}

			got := make([]string, len(ruleset.EnabledRules))
			for i, r := range ruleset.EnabledRules {
				got[i] = r.Name()
			}

			if diff := cmp.Diff(got, test.want); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
