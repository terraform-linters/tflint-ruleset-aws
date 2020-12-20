package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/models"
)

// Rules is a list of all rules
var Rules = append([]tflint.Rule{
	NewAwsDBInstanceDefaultParameterGroupRule(),
	NewAwsDBInstanceInvalidTypeRule(),
	NewAwsDBInstancePreviousTypeRule(),
	NewAwsDynamoDBTableInvalidStreamViewTypeRule(),
	NewAwsElastiCacheClusterDefaultParameterGroupRule(),
	NewAwsElastiCacheClusterInvalidTypeRule(),
	NewAwsElastiCacheClusterPreviousTypeRule(),
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
	NewAwsInstanceExampleTypeRule(),
}, models.Rules...)
