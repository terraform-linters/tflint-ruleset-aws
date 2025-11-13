import = "api-models-aws/models/kinesis/service/2013-12-02/kinesis-2013-12-02.json"

mapping "aws_kinesis_stream" {
  name                      = StreamName
  shard_count               = PositiveIntegerObject
  retention_period          = RetentionPeriodHours
  shard_level_metrics       = MetricsNameList
  enforce_consumer_deletion = BooleanObject
  encryption_type           = EncryptionType
  kms_key_id                = KeyId
  tags                      = TagList
}
