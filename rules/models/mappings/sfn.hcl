import = "api-models-aws/models/sfn/service/2016-11-23/sfn-2016-11-23.json"

mapping "aws_sfn_activity" {
  name = Name
  tags = listmap(TagList, TagKey, TagValue)
}

mapping "aws_sfn_state_machine" {
  name       = Name
  definition = Definition
  role_arn   = Arn
  tags       = listmap(TagList, TagKey, TagValue)
}
