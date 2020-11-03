package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/aws"
)

// AwsCustomRunnerRule checks whether ...
type AwsCustomRunnerRule struct{}

// NewAwsCustomRunnerRule returns a new rule
func NewAwsCustomRunnerRule() *AwsCustomRunnerRule {
	return &AwsCustomRunnerRule{}
}

// Name returns the rule name
func (r *AwsCustomRunnerRule) Name() string {
	return "aws_custom_runner"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsCustomRunnerRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsCustomRunnerRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsCustomRunnerRule) Link() string {
	return ""
}

// Check checks whether ...
func (r *AwsCustomRunnerRule) Check(rr tflint.Runner) error {
	runner := rr.(*aws.Runner)

	return runner.EmitIssue(
		r,
		fmt.Sprintf("config is %s", runner.CustomCall()),
		hcl.Range{},
	)
}
