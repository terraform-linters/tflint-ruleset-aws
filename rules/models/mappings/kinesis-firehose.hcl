import = "api-models-aws/models/firehose/service/2015-08-04/firehose-2015-08-04.json"

mapping "aws_kinesis_firehose_delivery_stream" {
  name                         = DeliveryStreamName
  tags                         = listmap(TagDeliveryStreamInputTagList, TagKey, TagValue)
  kinesis_source_configuration = KinesisStreamSourceConfiguration
  destination                  = any
  extended_s3_configuration    = ExtendedS3DestinationConfiguration
  redshift_configuration       = RedshiftDestinationConfiguration
  elasticsearch_configuration  = ElasticsearchDestinationConfiguration
  splunk_configuration         = SplunkDestinationConfiguration
  // cloudwatch_logging_options   = CloudWatchLoggingOptions
  // processing_configuration     = ProcessingConfiguration
  // processors                   = ProcessorList
  // parameters                   = ProcessorParameterList
}
