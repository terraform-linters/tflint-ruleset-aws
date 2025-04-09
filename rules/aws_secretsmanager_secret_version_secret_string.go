package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/zclconf/go-cty/cty"
)

// AwsSecretsmanagerSecretVersionSecretStringRule checks if the secret_string attribute is used in secretsmanager_secret_version
// It emits a warning if the attribute is used, as it is a non-ephemeral attribute and the secret is stored in state
// It suggests using secret_string_wo instead
type AwsSecretsmanagerSecretVersionSecretStringRule struct {
	tflint.DefaultRule

	resourceType           string
	attributeName          string
	writeOnlyAttributeName string
}

// NewAwsSecretsmanagerSecretVersionSecretStringRule returns new rule with default attributes
func NewAwsSecretsmanagerSecretVersionSecretStringRule() *AwsSecretsmanagerSecretVersionSecretStringRule {
	return &AwsSecretsmanagerSecretVersionSecretStringRule{
		resourceType:           "aws_secretsmanager_secret_version",
		attributeName:          "secret_string",
		writeOnlyAttributeName: "secret_string_wo",
	}
}

// Name returns the rule name
func (r *AwsSecretsmanagerSecretVersionSecretStringRule) Name() string {
	return "aws_secretsmanager_secret_version_secret_string"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecretsmanagerSecretVersionSecretStringRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSecretsmanagerSecretVersionSecretStringRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsSecretsmanagerSecretVersionSecretStringRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the secret string attribute exists
func (r *AwsSecretsmanagerSecretVersionSecretStringRule) Check(runner tflint.Runner) error {
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

		err := runner.EvaluateExpr(attribute.Expr, func(val cty.Value) error {
			if !val.IsNull() {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only attribute \"%s\".", r.attributeName, r.writeOnlyAttributeName),
					attribute.Expr.Range(),
				)
			}
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
