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

// NewRule returns new rule with default attributes
func NewRule() *AwsLambdaFunctionDeprecatedRuntimeRule {
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

// Check checks if the chosen runtime has reached EOS.
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
			rt, ok := runtimes.Runtimes[val]
			if !ok || !r.Now.After(rt.EndOfSupportDate) {
				return nil
			}

			message := runtimeMessage(val, rt, runtimes.UpdatedAt, r.Now)
			runner.EmitIssue(r, message, attribute.Expr.Range())
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func runtimeMessage(name string, rt runtimeLifecycle, updatedAt, now time.Time) string {
	past := func(t *time.Time) bool {
		return t != nil && now.After(*t)
	}
	confirmed := func(t *time.Time) bool {
		return past(t) && !t.After(updatedAt)
	}

	stale := fmt.Sprintf("as of %s", updatedAt.Format("Jan 2, 2006"))

	switch {
	case confirmed(rt.BlockUpdateDate):
		return fmt.Sprintf("The %q runtime has reached end of support and function updates are blocked", name)
	case past(rt.BlockUpdateDate):
		return fmt.Sprintf("The %q runtime has reached end of support and function updates were scheduled to be blocked on %s (%s)", name, rt.BlockUpdateDate.Format("Jan 2, 2006"), stale)
	case confirmed(rt.BlockCreateDate):
		return fmt.Sprintf("The %q runtime has reached end of support and new function creation is blocked", name)
	case past(rt.BlockCreateDate):
		return fmt.Sprintf("The %q runtime has reached end of support and new function creation was scheduled to be blocked on %s (%s)", name, rt.BlockCreateDate.Format("Jan 2, 2006"), stale)
	case rt.EndOfSupportDate.After(updatedAt):
		return fmt.Sprintf("The %q runtime was scheduled to reach end of support on %s (%s)", name, rt.EndOfSupportDate.Format("Jan 2, 2006"), stale)
	default:
		return fmt.Sprintf("The %q runtime has reached end of support", name)
	}
}
