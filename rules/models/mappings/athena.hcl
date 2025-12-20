import = "api-models-aws/models/athena/service/2017-05-18/athena-2017-05-18.json"

mapping "aws_athena_database" {
  name = DatabaseString
}

mapping "aws_athena_named_query" {
  name        = NameString
  database    = DatabaseString
  query       = QueryString
  description = DescriptionString
}

mapping "aws_athena_workgroup" {
  name          = WorkGroupName
  configuration = WorkGroupConfiguration
  description   = WorkGroupDescriptionString
  state         = WorkGroupState
  tags          = listmap(TagList, TagKey, TagValue)
}
