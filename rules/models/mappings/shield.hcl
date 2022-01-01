import = "aws-sdk-go/models/apis/shield/2016-06-02/api-2.json"

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
