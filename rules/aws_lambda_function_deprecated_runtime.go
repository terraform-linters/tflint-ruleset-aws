package rules

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
	eosRuntimes   map[string]time.Time
	eolRuntimes   map[string]time.Time

	Now time.Time
}

// NewAwsLambdaFunctionDeprecatedRuntimeRule returns new rule with default attributes
func NewAwsLambdaFunctionDeprecatedRuntimeRule() *AwsLambdaFunctionDeprecatedRuntimeRule {
	return &AwsLambdaFunctionDeprecatedRuntimeRule{
		resourceType:  "aws_lambda_function",
		attributeName: "runtime",
		eosRuntimes: map[string]time.Time{
			"nodejs10.x":    time.Date(2021, time.July, 30, 0, 0, 0, 0, time.UTC),
			"ruby2.5":       time.Date(2021, time.July, 30, 0, 0, 0, 0, time.UTC),
			"python2.7":     time.Date(2021, time.July, 15, 0, 0, 0, 0, time.UTC),
			"dotnetcore2.1": time.Date(2021, time.September, 20, 0, 0, 0, 0, time.UTC),
		},
		eolRuntimes: map[string]time.Time{
			"dotnetcore1.0":  time.Date(2019, time.July, 30, 0, 0, 0, 0, time.UTC),
			"dotnetcore2.0":  time.Date(2019, time.May, 30, 0, 0, 0, 0, time.UTC),
			"nodejs":         time.Date(2016, time.October, 31, 0, 0, 0, 0, time.UTC),
			"nodejs4.3":      time.Date(2020, time.March, 06, 0, 0, 0, 0, time.UTC),
			"nodejs4.3-edge": time.Date(2019, time.April, 30, 0, 0, 0, 0, time.UTC),
			"nodejs6.10":     time.Date(2019, time.August, 12, 0, 0, 0, 0, time.UTC),
			"nodejs8.10":     time.Date(2020, time.March, 06, 0, 0, 0, 0, time.UTC),
			"nodejs10.x":     time.Date(2021, time.August, 30, 0, 0, 0, 0, time.UTC),
			"ruby2.5":        time.Date(2021, time.August, 30, 0, 0, 0, 0, time.UTC),
			"python2.7":      time.Date(2021, time.September, 30, 0, 0, 0, 0, time.UTC),
			"dotnetcore2.1":  time.Date(2021, time.October, 30, 0, 0, 0, 0, time.UTC),
		},
		Now: time.Now().UTC(),
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

		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)

		err = runner.EnsureNoError(err, func() error {
			if _, ok := r.eolRuntimes[val]; ok && r.Now.After(r.eolRuntimes[val]) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("The \"%s\" runtime has reached the end of life", val),
					attribute.Expr.Range(),
				)
			} else if _, ok := r.eosRuntimes[val]; ok && r.Now.After(r.eosRuntimes[val]) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("The \"%s\" runtime has reached the end of support", val),
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
