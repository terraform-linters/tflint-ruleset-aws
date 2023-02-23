import = "aws-sdk-go/models/apis/dynamodb/2012-08-10/api-2.json"

mapping "aws_dynamodb_global_table" {
  name    = TableName
  replica = ReplicaList
}

mapping "aws_dynamodb_kinesis_streaming_destination" {
  stream_arn = StreamArn
  table_name = TableName
}

mapping "aws_dynamodb_table" {
  name                   = TableName
  billing_mode           = BillingMode
  hash_key               = KeySchemaAttributeName
  range_key              = KeySchemaAttributeName
  write_capacity         = PositiveLongObject
  read_capacity          = PositiveLongObject
  attribute              = AttributeDefinitions
  local_secondary_index  = LocalSecondaryIndexList
  global_secondary_index = GlobalSecondaryIndexList
  stream_enabled         = StreamEnabled
  stream_view_type       = any # StreamViewType
  server_side_encryption = SSESpecification
  tags                   = TagList
  point_in_time_recovery = PointInTimeRecoverySpecification
}

mapping "aws_dynamodb_table_item" {
  table_name = TableName
  hash_key   = KeySchemaAttributeName
  range_key  = KeySchemaAttributeName
  item       = AttributeMap
}

mapping "aws_dynamodb_tag" {
  resource_arn = ResourceArnString
}

test "aws_dynamodb_global_table" "name" {
  valid   = ["myTable"]
  invalid = ["myTable@development"]
}

test "aws_dynamodb_table" "billing_mode" {
  valid   = ["PROVISIONED"]
  invalid = ["FLEXIBLE"]
}

