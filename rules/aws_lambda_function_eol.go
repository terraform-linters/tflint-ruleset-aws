package rules

import (
	"fmt"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"time"
)

// AwsLambdaFunctionEolRule checks to see if the lambda runtime has reached End Of Life
type AwsLambdaFunctionEolRule struct {
	resourceType  string
	attributeName string
	eolRuntimes   map[string]time.Time
}

// NewAwsLambdaFunctionEolRule returns new rule with default attributes
func NewAwsLambdaFunctionEolRule() *AwsLambdaFunctionEolRule {
	return &AwsLambdaFunctionEolRule{
		resourceType:  "aws_lambda_function",
		attributeName: "runtime",
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
	}
}

// Name returns the rule name
func (r *AwsLambdaFunctionEolRule) Name() string {
	return "aws_lambda_function_eol"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLambdaFunctionEolRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsLambdaFunctionEolRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsLambdaFunctionEolRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks if the chosen runtime has reached EOL. Date check allows future values to be created as well.
func (r *AwsLambdaFunctionEolRule) Check(runner tflint.Runner) error {

	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)
		now := time.Now().UTC()

		return runner.EnsureNoError(err, func() error {
			if _, ok := r.eolRuntimes[val]; ok && now.After(r.eolRuntimes[val]) {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("The \"%s\" runtime has reached the end of life", val),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
