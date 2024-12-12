import = "aws-sdk-ruby/apis/timestream-write/2018-11-01/api-2.json"

mapping "aws_timestreamwrite_database" {
  database_name = ResourceCreateAPIName
  kms_key_id = StringValue2048
  tags = TagList
}

mapping "aws_timestreamwrite_table" {
  database_name = ResourceCreateAPIName
  retention_properties = RetentionProperties
  table_name = ResourceCreateAPIName
  tags = TagList
}
