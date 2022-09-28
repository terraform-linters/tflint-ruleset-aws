package rules

import (
    "strings"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsPolicyHeredocsDisallowedRule checks that policies are not defined with Heredoc format
type AwsPolicyHeredocsDisallowedRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsPolicyHeredocsDisallowedRule returns new rule with default attributes
func NewAwsPolicyHeredocsDisallowedRule() *AwsPolicyHeredocsDisallowedRule {
	return &AwsPolicyHeredocsDisallowedRule{
		resourceType:  "aws_iam_policy",
		attributeName: "policy",
	}
}

// Name returns the rule name
func (r *AwsPolicyHeredocsDisallowedRule) Name() string {
	return "aws_policy_heredocs_disallowed"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsPolicyHeredocsDisallowedRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsPolicyHeredocsDisallowedRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsPolicyHeredocsDisallowedRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks if a heredoc block is in use.
func (r *AwsPolicyHeredocsDisallowedRule) Check(runner tflint.Runner) error {

	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for name, file := range files {
		if err := r.checkHeredocs(runner, name, file); err != nil {
			return err
		}
	}

	return nil
}

func (r *AwsPolicyHeredocsDisallowedRule) checkHeredocs(runner tflint.Runner, filename string, file *hcl.File) error {
	if strings.HasSuffix(filename, ".json") {
		return nil
	}

	tokens, diags := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return diags
	}

	for _, token := range tokens {
		if token.Type != hclsyntax.TokenOHeredoc && token.Type != hclsyntax.TokenCHeredoc {
			continue
		}

		if err := runner.EmitIssue(
			r,
			"Avoid Use of HEREDOC syntax for policy documents",
			token.Range,
		); err != nil {
			return err
		}
	}

	return nil
}
