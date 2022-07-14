package rules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsSecurityGroupRuleInvalidProtocolRule checks whether "protocol" has invalid value
type AwsSecurityGroupRuleInvalidProtocolRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
	protocols     map[string]struct{}
}

// NewAwsSecurityGroupRuleInvalidProtocolRule returns new rule with default attributes
func NewAwsSecurityGroupRuleInvalidProtocolRule() *AwsSecurityGroupRuleInvalidProtocolRule {
	return &AwsSecurityGroupRuleInvalidProtocolRule{
		resourceType:  "aws_security_group_rule",
		attributeName: "protocol",
		protocols: map[string]struct{}{
			"all":    {},
			"tcp":    {},
			"udp":    {},
			"icmp":   {},
			"icmpv6": {},
		},
	}
}

// Name returns the rule name
func (r *AwsSecurityGroupRuleInvalidProtocolRule) Name() string {
	return "aws_security_group_rule_invalid_protocol"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupRuleInvalidProtocolRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSecurityGroupRuleInvalidProtocolRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSecurityGroupRuleInvalidProtocolRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether "protocol" has invalid value
func (r *AwsSecurityGroupRuleInvalidProtocolRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		var protocol string
		err := runner.EvaluateExpr(attribute.Expr, &protocol, nil)

		err = runner.EnsureNoError(err, func() error {
			if _, err := strconv.Atoi(protocol); err == nil {
				return nil
			}
			if _, ok := r.protocols[strings.ToLower(protocol)]; !ok {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is an invalid protocol.", protocol),
					attribute.Expr.Range(),
				)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
