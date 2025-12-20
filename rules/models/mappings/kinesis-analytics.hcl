import = "api-models-aws/models/kinesis-analytics/service/2015-08-14/kinesis-analytics-2015-08-14.json"

mapping "aws_kinesis_analytics_application" {
  name                       = ApplicationName
  code                       = ApplicationCode
  description                = ApplicationDescription
  cloudwatch_logging_options = CloudWatchLoggingOptions
  inputs                     = Inputs
  outputs                    = Outputs
  reference_data_sources     = ReferenceDataSource
  tags                       = listmap(Tags, TagKey, TagValue)
}
