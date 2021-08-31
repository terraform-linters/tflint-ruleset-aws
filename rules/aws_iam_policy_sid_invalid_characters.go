package rules

import (
	"encoding/json"
	"fmt"
	"regexp"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

type AwsIAMPolicySidInvalidCharactersStatementStruct struct {
	Sid string `json:"Sid"`
}
type AwsIAMPolicySidInvalidCharactersPolicyStruct struct {
	Statement []AwsIAMPolicySidInvalidCharactersStatementStruct `json:"Statement"`
}

// AwsIAMPolicySidInvalidCharactersRule checks for invalid characters in SID
type AwsIAMPolicySidInvalidCharactersRule struct {
	resourceType    string
	attributeName   string
	validCharacters *regexp.Regexp
}

// NewAwsIAMPolicySidInvalidCharactersRule returns new rule with default attributes
func NewAwsIAMPolicySidInvalidCharactersRule() *AwsIAMPolicySidInvalidCharactersRule {
	return &AwsIAMPolicySidInvalidCharactersRule{
		resourceType:    "aws_iam_policy",
		attributeName:   "policy",
		validCharacters: regexp.MustCompile(`^[a-zA-Z0-9]+$`),
	}
}

// Name returns the rule name
func (r *AwsIAMPolicySidInvalidCharactersRule) Name() string {
	return "aws_iam_policy_sid_invalid_characters"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIAMPolicySidInvalidCharactersRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsIAMPolicySidInvalidCharactersRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsIAMPolicySidInvalidCharactersRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks the unmarshaled policy and loops through statements checking for invalid statement ids
func (r *AwsIAMPolicySidInvalidCharactersRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var val string
		err := runner.EvaluateExpr(attribute.Expr, &val, nil)
		unMarshaledPolicy := AwsIAMPolicySidInvalidCharactersPolicyStruct{}
		if jsonErr := json.Unmarshal([]byte(val), &unMarshaledPolicy); jsonErr != nil {
			return jsonErr
		}
		statements := unMarshaledPolicy.Statement

		return runner.EnsureNoError(err, func() error {
			for _, statement := range statements {
				if statement.Sid == "" {
					continue
				}

				if !r.validCharacters.MatchString(statement.Sid) {
					runner.EmitIssueOnExpr(
						r,
						fmt.Sprintf("The policy's sid (\"%s\") does not match \"%s\".", statement.Sid, r.validCharacters.String()),
						attribute.Expr,
					)
				}
			}
			return nil
		})
	})
}
