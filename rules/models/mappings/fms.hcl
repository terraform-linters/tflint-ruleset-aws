import = "aws-sdk-ruby/apis/fms/2018-01-01/api-2.json"

mapping "aws_fms_admin_account" {
  account_id = AWSAccountId
}

mapping "aws_fms_policy" {
  name = ResourceName
  exclude_map = CustomerPolicyScopeMap
  include_map = CustomerPolicyScopeMap
  resource_tags = ResourceTags
  resource_type = ResourceType
  resource_type_list = ResourceTypeList
  security_service_policy_data = SecurityServicePolicyData
}
