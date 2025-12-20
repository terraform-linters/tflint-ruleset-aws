import = "api-models-aws/models/wafv2/service/2019-07-29/wafv2-2019-07-29.json"

mapping "aws_wafv2_ip_set" {
  name = EntityName
  description = EntityDescription
  scope = Scope
  ip_address_version = IPAddressVersion
  addresses = IPAddresses
  tags = listmap(TagList, TagKey, TagValue)
}

mapping "aws_wafv2_regex_pattern_set" {
  name = EntityName
  description = EntityDescription
  scope = Scope
  regular_expression = RegularExpressionList
  tags = listmap(TagList, TagKey, TagValue)
}

mapping "aws_wafv2_rule_group" {
  capacity = CapacityUnit
  custom_response_body = CustomResponseBodies
  description = EntityDescription
  name = EntityName
  rule = Rules
  scope = Scope
  tags = listmap(TagList, TagKey, TagValue)
  visibility_config = VisibilityConfig
}

mapping "aws_wafv2_web_acl" {
  custom_response_body = CustomResponseBodies
  default_action = DefaultAction
  description = EntityDescription
  name = EntityName
  rule = Rules
  scope = Scope
  tags = listmap(TagList, TagKey, TagValue)
  visibility_config = VisibilityConfig
}

mapping "aws_wafv2_web_acl_association" {
  resource_arn = ResourceArn
  web_acl_arn = ResourceArn
}

mapping "aws_wafv2_web_acl_logging_configuration" {
  log_destination_configs = LogDestinationConfigs
  logging_filter = LoggingFilter
  redacted_fields = RedactedFields
  resource_arn = ResourceArn
}
