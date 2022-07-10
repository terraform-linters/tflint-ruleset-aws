package rules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsSecurityGroupInvalidProtocolRule checks whether "protocol" in "ingress" or "egress" has invalid value
type AwsSecurityGroupInvalidProtocolRule struct {
	tflint.DefaultRule

	resourceType string
	protocols    map[string]struct{}
}

// NewAwsSecurityGroupInvalidProtocolRule returns new rule with default attributes
func NewAwsSecurityGroupInvalidProtocolRule() *AwsSecurityGroupInvalidProtocolRule {
	return &AwsSecurityGroupInvalidProtocolRule{
		resourceType: "aws_security_group",
		protocols: map[string]struct{}{
			"tcp":    {},
			"udp":    {},
			"icmp":   {},
			"icmpv6": {},
		},
	}
}

// Name returns the rule name
func (r *AwsSecurityGroupInvalidProtocolRule) Name() string {
	return "aws_security_group_invalid_protocol"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupInvalidProtocolRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSecurityGroupInvalidProtocolRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSecurityGroupInvalidProtocolRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether "protocol" in "ingress" or "egress" has invalid value
func (r *AwsSecurityGroupInvalidProtocolRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "ingress",
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "protocol"},
					},
				},
			},
			{
				Type: "egress",
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "protocol"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, rule := range resource.Body.Blocks {
			attribute, exists := rule.Body.Attributes["protocol"]
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
	}

	return nil
}
