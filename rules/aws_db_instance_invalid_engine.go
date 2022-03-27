package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsDBInstanceInvalidEngineRule checks whether "aws_db_instance" has invalid engine.
type AwsDBInstanceInvalidEngineRule struct {
	tflint.DefaultRule

	resource      string
	attributeName string
	engines       map[string]bool
}

// NewAwsDBInstanceInvalidEngineRule returns new rule with default attributes
func NewAwsDBInstanceInvalidEngineRule() *AwsDBInstanceInvalidEngineRule {
	return &AwsDBInstanceInvalidEngineRule{
		resource:      "aws_db_instance",
		attributeName: "engine",
		engines: map[string]bool{
			// https://docs.aws.amazon.com/cli/latest/reference/rds/describe-db-engine-versions.html#options
			"aurora":            true,
			"aurora-mysql":      true,
			"aurora-postgresql": true,
			"mariadb":           true,
			"mysql":             true,
			"oracle-ee":         true,
			"oracle-ee-cdb":     true,
			"oracle-se2":        true,
			"oracle-se2-cdb":    true,
			"oracle-se1":        true,
			"oracle-se":         true,
			"postgres":          true,
			"sqlserver-ee":      true,
			"sqlserver-se":      true,
			"sqlserver-ex":      true,
			"sqlserver-web":     true,
		},
	}
}

// Name returns the rule name
func (r *AwsDBInstanceInvalidEngineRule) Name() string {
	return "aws_db_instance_invalid_engine"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsDBInstanceInvalidEngineRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsDBInstanceInvalidEngineRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsDBInstanceInvalidEngineRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether "aws_db_instance" has invalid engine.
func (r *AwsDBInstanceInvalidEngineRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resource, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.attributeName}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		var engine string
		err := runner.EvaluateExpr(attribute.Expr, &engine, nil)

		err = runner.EnsureNoError(err, func() error {
			if !r.engines[engine] {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is invalid engine.", engine),
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
