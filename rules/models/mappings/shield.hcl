import = "api-models-aws/models/shield/service/2016-06-02/shield-2016-06-02.json"

mapping "aws_shield_protection" {
  name         = ProtectionName
  resource_arn = ResourceArn
}

mapping "aws_shield_protection_group" {
  aggregation = ProtectionGroupAggregation
  members = ProtectionGroupMembers
  pattern = ProtectionGroupPattern
  protection_group_id = ProtectionGroupId
  resource_type = ProtectedResourceType
  tags = TagList
}
