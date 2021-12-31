package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/api"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/models"
)

var rules = [][]tflint.Rule{
	[]tflint.Rule{
		NewAwsDBInstanceDefaultParameterGroupRule(),
		NewAwsDBInstanceInvalidEngineRule(),
		NewAwsDBInstanceInvalidTypeRule(),
		NewAwsDBInstancePreviousTypeRule(),
		NewAwsDynamoDBTableInvalidStreamViewTypeRule(),
		NewAwsElastiCacheClusterDefaultParameterGroupRule(),
		NewAwsElastiCacheClusterInvalidTypeRule(),
		NewAwsElastiCacheClusterPreviousTypeRule(),
		NewAwsIAMPolicyDocumentGovFriendlyArnsRule(),
		NewAwsIAMPolicyGovFriendlyArnsRule(),
		NewAwsIAMRolePolicyGovFriendlyArnsRule(),
		NewAwsInstancePreviousTypeRule(),
		NewAwsMqBrokerInvalidEngineTypeRule(),
		NewAwsMqConfigurationInvalidEngineTypeRule(),
		NewAwsResourceMissingTagsRule(),
		NewAwsRouteNotSpecifiedTargetRule(),
		NewAwsRouteSpecifiedMultipleTargetsRule(),
		NewAwsS3BucketInvalidACLRule(),
		NewAwsS3BucketInvalidRegionRule(),
		NewAwsS3BucketNameRule(),
		NewAwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule(),
		NewAwsAPIGatewayModelInvalidNameRule(),
		NewAwsElastiCacheReplicationGroupDefaultParameterGroupRule(),
		NewAwsElastiCacheReplicationGroupInvalidTypeRule(),
		NewAwsElastiCacheReplicationGroupPreviousTypeRule(),
		NewAwsIAMPolicySidInvalidCharactersRule(),
		NewAwsIAMPolicyTooLongPolicyRule(),
		NewAwsLambdaFunctionDeprecatedRuntimeRule(),
		NewAwsIAMGroupPolicyTooLongRule(),
		NewAwsAcmCertificateLifecycleRule(),
	},
	models.Rules,
	api.Rules,
}

// Rules is a list of all rules
var Rules []tflint.Rule

func init() {
	for _, r := range rules {
		Rules = append(Rules, r...)
	}
}
