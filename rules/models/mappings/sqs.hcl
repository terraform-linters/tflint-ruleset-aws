import = "aws-sdk-ruby/apis/sqs/2012-11-05/api-2.json"

mapping "aws_sqs_queue" {
  name                              = String
  name_prefix                       = any
  visibility_timeout_seconds        = NullableInteger
  message_retention_seconds         = NullableInteger
  max_message_size                  = NullableInteger
  delay_seconds                     = NullableInteger
  receive_wait_time_seconds         = NullableInteger
  policy                            = String
  redrive_policy                    = String
  fifo_queue                        = Boolean
  content_based_deduplication       = Boolean
  kms_master_key_id                 = String
  kms_data_key_reuse_period_seconds = NullableInteger
  tags                              = TagMap
}

mapping "aws_sqs_queue_policy" {
  queue_url = String
  policy    = String
}
