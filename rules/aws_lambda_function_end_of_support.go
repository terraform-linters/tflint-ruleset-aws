package rules

import (
	"fmt"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"time"
)

// AwsLambdaFunctionEndOfSupportRule checks to see if the lambda runtime has reached End Of Support
type AwsLambdaFunctionEndOfSupportRule struct {
	resourceType  string
	attributeName string
	eosRuntimes   map[string]time.Time
}

// NewAwsLambdaFunctionEndOfSupportRule returns new rule with default attributes
func NewAwsLambdaFunctionEndOfSupportRule() *AwsLambdaFunctionEndOfSupportRule {
	return &AwsLambdaFunctionEndOfSupportRule{
		resourceType:  "aws_lambda_function",
		attributeName: "runtime",
		eosRuntimes:  map[string]time.Time{
			"nodejs10.x":     time.Date(2021, time.July, 30, 0, 0, 0, 0, time.UTC),
			"ruby2.5":        time.Date(2021, time.July, 30, 0, 0, 0, 0, time.UTC),
			"python2.7":      time.Date(2021, time.July, 15, 0, 0, 0, 0, time.UTC),
			"dotnetcore2.1":  time.Date(2021, time.September, 20, 0, 0, 0, 0, time.UTC),
		},
	}
}

// Name returns the rule name
func (r *AwsLambdaFunctionEndOfSupportRule) Name() string {
	return "aws_lambda_function_end_of_support"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsLambdaFunctionEndOfSupportRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsLambdaFunctionEndOfSupportRule) Severity() string {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsLambdaFunctionEndOfSupportRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks if the chosen runtime has reached EOS. Date check allows future values to be created as well.
func (r *AwsLambdaFunctionEndOfSupportRule) Check(runner tflint.Runner) error {

	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)
		now := time.Now().UTC()

		return runner.EnsureNoError(err, func() error {
			if _, ok := r.eosRuntimes[val]; ok && now.After(r.eosRuntimes[val]) {
				runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("The \"%s\" runtime has reached the end of support", val),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}
