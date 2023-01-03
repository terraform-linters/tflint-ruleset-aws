package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSecurityGroupRuleInvalidCidrBlockRule checks whether ...
type AwsSecurityGroupRuleInvalidCidrBlockRule struct {
	tflint.DefaultRule

	resourceType      string
	cidrBlocks        []string
	remoteAccessPorts []int
}

// NewAwsSecurityGroupRuleInvalidCidrBlockRule returns a new rule
func NewAwsSecurityGroupRuleInvalidCidrBlockRule() *AwsSecurityGroupRuleInvalidCidrBlockRule {
	return &AwsSecurityGroupRuleInvalidCidrBlockRule{
		resourceType:      "aws_security_group_rule",
		remoteAccessPorts: []int{22, 3389},
	}
}

// Name returns the rule name
func (r *AwsSecurityGroupRuleInvalidCidrBlockRule) Name() string {
	return "aws_security_group_rule_invalid_cidr_block"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupRuleInvalidCidrBlockRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSecurityGroupRuleInvalidCidrBlockRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSecurityGroupRuleInvalidCidrBlockRule) Link() string {
	return "TODO"
}

// Check checks whether ...
func (r *AwsSecurityGroupRuleInvalidCidrBlockRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			// Required attributes
			{Name: "from_port"},
			{Name: "to_port"},
			{Name: "type"},
			// Optional attributes
			{Name: "cidr_blocks"},
			{Name: "ipv6_cidr_blocks"},
		},
	}, nil)
	if err != nil {
		return err
	}

	// Put a log that can be output with `TFLINT_LOG=debug`
	logger.Debug(fmt.Sprintf("Get %d security groups", len(resources.Blocks)))

	for _, resource := range resources.Blocks {
		typeAttribute, exists := resource.Body.Attributes["type"]

		var sgType string
		err := runner.EvaluateExpr(typeAttribute.Expr, &sgType, nil)

		logger.Debug(fmt.Sprintf("Found %s type of Security Group Rule", sgType))

		if sgType != "ingress" {
			continue
		}

		//  A Port value of ALL
		fromPortAttribute, exists := resource.Body.Attributes["from_port"]
		if !exists {
			continue
		}

		var fromPort int
		err = runner.EvaluateExpr(fromPortAttribute.Expr, &fromPort, nil)

		toPortAttribute, exists := resource.Body.Attributes["to_port"]
		if !exists {
			continue
		}

		var toPort int
		err = runner.EvaluateExpr(toPortAttribute.Expr, &toPort, nil)

		logger.Debug(fmt.Sprintf("Security Group Rule port range is %d to %d", fromPort, toPort))

		// No need to check the CIDR blocks if the ports range does not contain 22 or 3389.
		if !doesPortRangeContainsPorts(fromPort, toPort, r.remoteAccessPorts) {
			continue
		}

		cidrBlocksAttribute, exists := resource.Body.Attributes["cidr_blocks"]
		if exists {
			var cidrBlocks []string
			err = runner.EvaluateExpr(cidrBlocksAttribute.Expr, &cidrBlocks, nil)
			if containsIpv4CidrBlocksAllowAll(cidrBlocks) {
				return runner.EmitIssue(
					r,
					fmt.Sprintf("cidr_blocks can not contain '0.0.0.0/0' when allowing 'ingress' access to ports %v", r.remoteAccessPorts),
					cidrBlocksAttribute.Expr.Range(),
				)
			}
		}

		ipv6CidrBlocksAttribute, exists := resource.Body.Attributes["ipv6_cidr_blocks"]
		if !exists {
			continue
		}

		var ipv6CidrBlocks []string
		err = runner.EvaluateExpr(ipv6CidrBlocksAttribute.Expr, &ipv6CidrBlocks, nil)
		if containsIpv6CidrBlocksAllowAll(ipv6CidrBlocks) {
			return runner.EmitIssue(
				r,
				fmt.Sprintf("ipv6_cidr_blocks can not contain '::/0' when allowing 'ingress' access to ports %v", r.remoteAccessPorts),
				ipv6CidrBlocksAttribute.Expr.Range(),
			)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func doesPortRangeContainsPorts(fromPort int, toPort int, ports []int) bool {
	for _, port := range ports {
		isIncludedInRange := fromPort <= port && port <= toPort
		logger.Debug(fmt.Sprintf("%v for %d isIncludedInRange", isIncludedInRange, port))
		if isIncludedInRange {
			return true
		}
	}

	return false
}

func containsIpv4CidrBlocksAllowAll(cidrBlocks []string) bool {
	for _, cidrBlock := range cidrBlocks {
		if cidrBlock == "0.0.0.0/0" {
			return true
		}
	}
	return false
}

func containsIpv6CidrBlocksAllowAll(ipv6CidrBlocks []string) bool {
	for _, cidrBlock := range ipv6CidrBlocks {
		if cidrBlock == "::/0" {
			return true
		}
	}
	return false
}
