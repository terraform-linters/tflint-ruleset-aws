package rules

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

var PresetRules = map[string][]tflint.Rule{
	"s3": {
		NewAwsS3BucketNameRule(),
		NewAwsS3BucketInvalidACLRule(),
	},
}
