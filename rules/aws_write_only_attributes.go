package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/zclconf/go-cty/cty"
)

// AwsWriteOnlyAttributesRule checks if a write-only attribute is available for sensitive input attributes
// It emits a warning if the normal attribute is used, as that is a non-ephemeral attribute and the secret value is stored in state
type AwsWriteOnlyAttributesRule struct {
	tflint.DefaultRule

	writeOnlyAttributes map[string]writeOnlyAttribute
}

type writeOnlyAttribute struct {
	original             string
	writeOnlyAlternative string
	otherAlternative     string
}

// NewAwsWriteOnlyAttributesRule returns new rule with default attributes
func NewAwsWriteOnlyAttributesRule() *AwsWriteOnlyAttributesRule {
	writeOnlyAttributes := map[string]writeOnlyAttribute{
		"aws_secretsmanager_secret_version": {
			original:             "secret_string",
			writeOnlyAlternative: "secret_string_wo",
		},
		"aws_rds_cluster": {
			original:             "master_password",
			writeOnlyAlternative: "master_password_wo",
			otherAlternative:     "manage_master_user_password",
		},
		"aws_redshift_cluster": {
			original:             "master_password",
			writeOnlyAlternative: "master_password_wo",
			otherAlternative:     "manage_master_password",
		},
		"aws_docdb_cluster": {
			original:             "master_password",
			writeOnlyAlternative: "master_password_wo",
		},
		"aws_redshiftserverless_namespace": {
			original:             "admin_password",
			writeOnlyAlternative: "admin_password_wo",
			otherAlternative:     "manage_admin_password",
		},
		"aws_ssm_parameter": {
			original:             "value",
			writeOnlyAlternative: "value_wo",
		},
	}
	return &AwsWriteOnlyAttributesRule{
		writeOnlyAttributes: writeOnlyAttributes,
	}
}

// Name returns the rule name
func (r *AwsWriteOnlyAttributesRule) Name() string {
	return "aws_write_only_attributes"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsWriteOnlyAttributesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsWriteOnlyAttributesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsWriteOnlyAttributesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the secret string attribute exists
func (r *AwsWriteOnlyAttributesRule) Check(runner tflint.Runner) error {
	for resourceType, attributes := range r.writeOnlyAttributes {
		resources, err := runner.GetResourceContent(resourceType, &hclext.BodySchema{
			Attributes: []hclext.AttributeSchema{
				{Name: attributes.original},
			},
		}, nil)
		if err != nil {
			return err
		}

		for _, resource := range resources.Blocks {
			attribute, exists := resource.Body.Attributes[attributes.original]
			if !exists {
				continue
			}

			err := runner.EvaluateExpr(attribute.Expr, func(val cty.Value) error {
				mitigation := fmt.Sprintf("\"%s\" is a non-ephemeral attribute, which means this secret is stored in state. Please use write-only attribute \"%s\".", attributes.original, attributes.writeOnlyAlternative)
				if attributes.otherAlternative != "" {
					mitigation += fmt.Sprintf(" Alternatively, you can use \"%s\" to manage the secret in an different way.", attributes.otherAlternative)
				}

				if !val.IsNull() {
					if err := runner.EmitIssueWithFix(
						r,
						mitigation,
						attribute.Expr.Range(),
						func(f tflint.Fixer) error {
							return f.ReplaceText(attribute.NameRange, attributes.writeOnlyAlternative)
						},
					); err != nil {
						return fmt.Errorf("failed to call EmitIssueWithFix(): %w", err)
					}
				}
				return nil
			}, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
