package lambda_deprecated_runtime

import (
	"fmt"
	"time"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsLambdaFunctionDeprecatedRuntimeRule checks to see if the lambda runtime has reached End Of Support
type AwsLambdaFunctionDeprecatedRuntimeRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string

	Now time.Time
}

// NewAwsLambdaFunctionDeprecatedRuntimeRule returns new rule with default attributes
func NewAwsLambdaFunctionDeprecatedRuntimeRule() *AwsLambdaFunctionDeprecatedRuntimeRule {
	return &AwsLambdaFunctionDeprecatedRuntimeRule{
		resourceType:  "aws_lambda_function",
		attributeName: "runtime",
		Now:           time.Now().UTC(),
	}
}

// Name returns the rule name
func (r *AwsLambdaFunctionDeprecatedRuntimeRule) Name() string {
	return "aws_lambda_function_deprecated_runtime"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLambdaFunctionDeprecatedRuntimeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsLambdaFunctionDeprecatedRuntimeRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsLambdaFunctionDeprecatedRuntimeRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks if the chosen runtime has reached EOS. Date check allows future values to be created as well.
func (r *AwsLambdaFunctionDeprecatedRuntimeRule) Check(runner tflint.Runner) error {
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

		err := runner.EvaluateExpr(attribute.Expr, func(val string) error {
			rt, ok := deprecatedRuntimes[val]
			if !ok || !r.Now.After(rt.EndOfSupportDate) {
				return nil
			}

			var message string
			switch {
			case rt.BlockUpdateDate != nil && r.Now.After(*rt.BlockUpdateDate):
				message = fmt.Sprintf("The %q runtime has reached end of support and function updates are blocked", val)
			case rt.BlockCreateDate != nil && r.Now.After(*rt.BlockCreateDate):
				message = fmt.Sprintf("The %q runtime has reached end of support and new function creation is blocked", val)
			default:
				message = fmt.Sprintf("The %q runtime has reached end of support", val)
			}

			runner.EmitIssue(r, message, attribute.Expr.Range())
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
