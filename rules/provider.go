package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/api"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/models"
)

var manualRules = []tflint.Rule{
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
	NewAwsElasticBeanstalkEnvironmentInvalidNameFormatRule(),
	NewAwsSecurityGroupInvalidProtocolRule(),
	NewAwsSecurityGroupRuleInvalidProtocolRule(),
	NewAwsPolicyHeredocsDisallowedRule(),
}

// Rules is a list of all rules
var Rules []tflint.Rule

func init() {
	Rules = append(Rules, manualRules...)
	Rules = append(Rules, models.Rules...)
	Rules = append(Rules, api.Rules...)
}
