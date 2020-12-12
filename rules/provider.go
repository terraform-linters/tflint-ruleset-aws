package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/models"
)

// Rules is a list of all rules
var Rules = append([]tflint.Rule{
	NewAwsInstanceExampleTypeRule(),
	NewAwsResourceMissingTagsRule(),
}, models.Rules...)
