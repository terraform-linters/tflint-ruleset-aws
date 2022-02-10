package rules

import (
	"fmt"
        "strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsInstancePreviousTypeRule checks whether the resource uses previous generation instance type
type AwsInstancePreviousTypeRule struct {
	resourceType          string
	attributeName         string
	previousInstanceTypes map[string]bool
}

// NewAwsInstancePreviousTypeRule returns new rule with default attributes
func NewAwsInstancePreviousTypeRule() *AwsInstancePreviousTypeRule {
	return &AwsInstancePreviousTypeRule{
		resourceType:  "aws_instance",
		attributeName: "instance_type",
		previousInstanceTypes: map[string]bool{
			"c1":  true,
			"c2":  true,
			"c3":  true,
			"cc2": true,
			"cg1": true,
			"cr1": true,
			"g2":  true,
			"hi1": true,
			"hs1": true,
			"i2":  true,
			"m1":  true,
			"m2":  true,
			"m3":  true,
			"r3":  true,
			"t1":  true,			
		},
	}
}

// Name returns the rule name
func (r *AwsInstancePreviousTypeRule) Name() string {
	return "aws_instance_previous_type"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsInstancePreviousTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsInstancePreviousTypeRule) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsInstancePreviousTypeRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the resource's `instance_type` is included in the list of previous generation instance type
func (r *AwsInstancePreviousTypeRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var instanceType string
		err := runner.EvaluateExpr(attribute.Expr, &instanceType, nil)

		return runner.EnsureNoError(err, func() error {
			if r.previousInstanceTypes[strings.Split(instanceType, ".")[0]] {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("\"%s\" is previous generation instance type.", instanceType),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
