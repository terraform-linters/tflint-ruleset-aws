package rules

import (
	"slices"

	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsIAMPolicyUsePolicyReference checks if aws_iam_policy doesn't assign to
// policy from a reference
type AwsIAMPolicyUsePolicyReference struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsIAMPolicyUsePolicyReferenceRule returns new rule with default attributes
func NewAwsIAMPolicyUsePolicyReferenceRule() *AwsIAMPolicyUsePolicyReference {
	return &AwsIAMPolicyUsePolicyReference{
		resourceType:  "aws_iam_policy",
		attributeName: "policy",
	}
}

// Name returns the rule name
func (r *AwsIAMPolicyUsePolicyReference) Name() string {
	return "aws_iam_policy_use_policy_reference"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMPolicyUsePolicyReference) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsIAMPolicyUsePolicyReference) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsIAMPolicyUsePolicyReference) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks if policy is not referencing an object reference
func (r *AwsIAMPolicyUsePolicyReference) Check(runner tflint.Runner) error {
	issueMessage := "The policy does not reference an object"

	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
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

		// Check if the value is a reference and if so permit it.
		if _, ok := attribute.Expr.(*hclsyntax.ScopeTraversalExpr); ok {
			continue
		}

		// Check if the value is a call to file or templatefile and if so permit it.
		permittedFunctions := []string{"file", "templatefile"}
		if functionCall, ok := attribute.Expr.(*hclsyntax.FunctionCallExpr); ok {
			if slices.Contains(permittedFunctions, functionCall.Name) {
				continue
			}
		}

		if err := runner.EmitIssue(r, issueMessage, attribute.Range); err != nil {
			return err
		}
	}

	return nil
}
