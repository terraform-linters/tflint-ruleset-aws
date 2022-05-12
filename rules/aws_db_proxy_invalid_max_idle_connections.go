package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsDBProxyInvalidMaxIdleConnectionsRule checks whether the resource has invalid max_idle_connections_percent.
type AwsDBProxyInvalidMaxIdleConnectionsRule struct {
	tflint.DefaultRule

	resourceType     string
	blockType        string
	attrMaxConns     string
	attrMaxIdleConns string
}

// NewAwsDBProxyInvalidMaxIdleConnectionsRule returns new rule with default attributes
func NewAwsDBProxyInvalidMaxIdleConnectionsRule() *AwsDBProxyInvalidMaxIdleConnectionsRule {
	return &AwsDBProxyInvalidMaxIdleConnectionsRule{
		resourceType:     "aws_db_proxy_default_target_group",
		blockType:        "connection_pool_config",
		attrMaxConns:     "max_connections_percent",
		attrMaxIdleConns: "max_idle_connections_percent",
	}
}

// Name returns the rule name
func (r *AwsDBProxyInvalidMaxIdleConnectionsRule) Name() string {
	return "aws_db_proxy_invalid_max_idle_connections"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsDBProxyInvalidMaxIdleConnectionsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsDBProxyInvalidMaxIdleConnectionsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsDBProxyInvalidMaxIdleConnectionsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the resource has invalid max_idle_connections_percent.
func (r *AwsDBProxyInvalidMaxIdleConnectionsRule) Check(runner tflint.Runner) error {
	bodySchema := &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.blockType,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name: r.attrMaxConns,
						},
						{
							Name: r.attrMaxIdleConns,
						},
					},
				},
			},
		},
	}
	resources, err := runner.GetResourceContent(r.resourceType, bodySchema, nil)
	if err != nil {
		return err
	}
	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			if block.Type != r.blockType {
				continue
			}
			attrMaxConns, existsMaxConns := block.Body.Attributes[r.attrMaxConns]
			attrMaxIdleConns, existsMaxIdleConns := block.Body.Attributes[r.attrMaxIdleConns]
			if !(existsMaxConns && existsMaxIdleConns) {
				continue
			}
			var (
				maxConns     int
				maxIdleConns int
			)
			if err := runner.EvaluateExpr(attrMaxConns.Expr, &maxConns, nil); err != nil {
				return err
			}
			if err := runner.EvaluateExpr(attrMaxIdleConns.Expr, &maxIdleConns, nil); err != nil {
				return err
			}
			if maxIdleConns > maxConns {
				msg := fmt.Sprintf("%s %d must be less than %s %d", r.attrMaxIdleConns, maxIdleConns, r.attrMaxConns, maxConns)
				runner.EmitIssue(r, msg, attrMaxIdleConns.Range)
			}
		}
	}
	return nil
}
