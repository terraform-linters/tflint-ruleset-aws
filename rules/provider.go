package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/models"
)

// Rules is a list of all rules
var Rules = append([]tflint.Rule{
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
}, models.Rules...)
