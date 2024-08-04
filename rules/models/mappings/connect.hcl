import = "aws-sdk-ruby/apis/connect/2017-08-08/api-2.json"

mapping "aws_connect_bot_association" {
  instance_id = InstanceId
  lex_bot = LexBot
}

mapping "aws_connect_contact_flow" {
  content = ContactFlowContent
  description = ContactFlowContent
  instance_id = InstanceId
  name = ContactFlowName
  tags = TagMap
  type = ContactFlowType
}

mapping "aws_connect_hours_of_operation" {
  config = HoursOfOperationConfigList
  description = HoursOfOperationDescription
  instance_id = InstanceId
  name = CommonNameLength127
  tags = TagMap
  time_zone = TimeZone
}

mapping "aws_connect_instance" {
  identity_management_type = DirectoryType
  inbound_calls_enabled = InboundCallsEnabled
  # instance_alias = DirectoryAlias
  outbound_calls_enabled = OutboundCallsEnabled
}

mapping "aws_connect_lambda_function_association" {
  function_arn = FunctionArn
  instance_id = InstanceId
}
