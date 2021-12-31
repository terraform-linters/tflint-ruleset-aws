import = "aws-sdk-go/models/apis/xray/2016-04-12/api-2.json"

mapping "aws_xray_encryption_config" {
  type = EncryptionType
  key_id = EncryptionKeyId
}

mapping "aws_xray_group" {
  group_name = GroupName
  filter_expression = FilterExpression
  tags = TagList
}

mapping "aws_xray_sampling_rule" {
  rule_name      = RuleName
  resource_arn   = ResourceARN
  priority       = Priority
  fixed_rate     = FixedRate
  reservoir_size = ReservoirSize
  service_name   = ServiceName
  service_type   = ServiceType
  host           = Host
  http_method    = HTTPMethod
  url_path       = URLPath
  version        = Version
  attributes     = AttributeMap
}
