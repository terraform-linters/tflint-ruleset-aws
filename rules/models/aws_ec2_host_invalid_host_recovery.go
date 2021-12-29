// This file generated by `generator/`. DO NOT EDIT

package models

import (
	"fmt"
	"log"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsEc2HostInvalidHostRecoveryRule checks the pattern is valid
type AwsEc2HostInvalidHostRecoveryRule struct {
	resourceType  string
	attributeName string
	enum          []string
}

// NewAwsEc2HostInvalidHostRecoveryRule returns new rule with default attributes
func NewAwsEc2HostInvalidHostRecoveryRule() *AwsEc2HostInvalidHostRecoveryRule {
	return &AwsEc2HostInvalidHostRecoveryRule{
		resourceType:  "aws_ec2_host",
		attributeName: "host_recovery",
		enum: []string{
			"on",
			"off",
		},
	}
}

// Name returns the rule name
func (r *AwsEc2HostInvalidHostRecoveryRule) Name() string {
	return "aws_ec2_host_invalid_host_recovery"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsEc2HostInvalidHostRecoveryRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsEc2HostInvalidHostRecoveryRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsEc2HostInvalidHostRecoveryRule) Link() string {
	return ""
}

// Check checks the pattern is valid
func (r *AwsEc2HostInvalidHostRecoveryRule) Check(runner tflint.Runner) error {
	log.Printf("[TRACE] Check `%s` rule", r.Name())

	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		return runner.EnsureNoError(err, func() error {
			found := false
			for _, item := range r.enum {
				if item == val {
					found = true
				}
			}
			if !found {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf(`"%s" is an invalid value as host_recovery`, truncateLongMessage(val)),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}