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

	resourceType       string
	attributeName      string
	deprecatedRuntimes map[string]time.Time

	Now time.Time
}

// NewAwsLambdaFunctionDeprecatedRuntimeRule returns new rule with default attributes
func NewAwsLambdaFunctionDeprecatedRuntimeRule() *AwsLambdaFunctionDeprecatedRuntimeRule {
	// @see https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html
	return &AwsLambdaFunctionDeprecatedRuntimeRule{
		resourceType:  "aws_lambda_function",
		attributeName: "runtime",
		deprecatedRuntimes: map[string]time.Time{
			"nodejs22.x":      time.Date(2027, time.April, 30, 0, 0, 0, 0, time.UTC),
			"nodejs20.x":      time.Date(2026, time.April, 30, 0, 0, 0, 0, time.UTC),
			"python3.13":      time.Date(2029, time.June, 30, 0, 0, 0, 0, time.UTC),
			"python3.12":      time.Date(2028, time.October, 31, 0, 0, 0, 0, time.UTC),
			"python3.11":      time.Date(2026, time.June, 30, 0, 0, 0, 0, time.UTC),
			"python3.10":      time.Date(2026, time.June, 30, 0, 0, 0, 0, time.UTC),
			"python3.9":       time.Date(2025, time.December, 15, 0, 0, 0, 0, time.UTC),
			"java21":          time.Date(2029, time.June, 30, 0, 0, 0, 0, time.UTC),
			"java17":          time.Date(2026, time.June, 30, 0, 0, 0, 0, time.UTC),
			"java11":          time.Date(2026, time.June, 30, 0, 0, 0, 0, time.UTC),
			"java8.al2":       time.Date(2026, time.June, 30, 0, 0, 0, 0, time.UTC),
			"dotnet8":         time.Date(2026, time.November, 10, 0, 0, 0, 0, time.UTC),
			"ruby3.3":         time.Date(2027, time.March, 31, 0, 0, 0, 0, time.UTC),
			"ruby3.2":         time.Date(2026, time.March, 31, 0, 0, 0, 0, time.UTC),
			"provided.al2023": time.Date(2029, time.June, 30, 0, 0, 0, 0, time.UTC),
			"provided.al2":    time.Date(2026, time.June, 30, 0, 0, 0, 0, time.UTC),
			// Already reached end of support
			"nodejs18.x":     time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC),
			"nodejs16.x":     time.Date(2024, time.June, 12, 0, 0, 0, 0, time.UTC),
			"python3.8":      time.Date(2024, time.October, 14, 0, 0, 0, 0, time.UTC),
			"java8":          time.Date(2024, time.January, 8, 0, 0, 0, 0, time.UTC),
			"dotnet7":        time.Date(2024, time.May, 14, 0, 0, 0, 0, time.UTC),
			"dotnet6":        time.Date(2024, time.November, 12, 0, 0, 0, 0, time.UTC),
			"go1.x":          time.Date(2024, time.January, 8, 0, 0, 0, 0, time.UTC),
			"provided":       time.Date(2024, time.January, 8, 0, 0, 0, 0, time.UTC),
			"nodejs14.x":     time.Date(2023, time.December, 4, 0, 0, 0, 0, time.UTC),
			"python3.7":      time.Date(2023, time.December, 4, 0, 0, 0, 0, time.UTC),
			"ruby2.7":        time.Date(2023, time.November, 27, 0, 0, 0, 0, time.UTC),
			"dotnetcore3.1":  time.Date(2023, time.April, 3, 0, 0, 0, 0, time.UTC),
			"nodejs12.x":     time.Date(2023, time.March, 31, 0, 0, 0, 0, time.UTC),
			"python3.6":      time.Date(2022, time.July, 18, 0, 0, 0, 0, time.UTC),
			"dotnet5.0":      time.Date(2022, time.May, 10, 0, 0, 0, 0, time.UTC),
			"dotnetcore2.1":  time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC),
			"nodejs10.x":     time.Date(2021, time.July, 30, 0, 0, 0, 0, time.UTC),
			"ruby2.5":        time.Date(2021, time.July, 30, 0, 0, 0, 0, time.UTC),
			"python2.7":      time.Date(2021, time.July, 15, 0, 0, 0, 0, time.UTC),
			"nodejs8.10":     time.Date(2020, time.March, 6, 0, 0, 0, 0, time.UTC),
			"nodejs4.3":      time.Date(2020, time.March, 5, 0, 0, 0, 0, time.UTC),
			"nodejs4.3-edge": time.Date(2020, time.March, 5, 0, 0, 0, 0, time.UTC),
			"nodejs6.10":     time.Date(2019, time.August, 12, 0, 0, 0, 0, time.UTC),
			"dotnetcore1.0":  time.Date(2019, time.July, 27, 0, 0, 0, 0, time.UTC),
			"dotnetcore2.0":  time.Date(2019, time.May, 30, 0, 0, 0, 0, time.UTC),
			"nodejs":         time.Date(2016, time.October, 31, 0, 0, 0, 0, time.UTC),
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

		err := runner.EvaluateExpr(attribute.Expr, func(val string) error {
			if _, ok := r.deprecatedRuntimes[val]; ok && r.Now.After(r.deprecatedRuntimes[val]) {
				runner.EmitIssue(
					r,
					fmt.Sprintf("The \"%s\" runtime has reached the end of support", val),
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
