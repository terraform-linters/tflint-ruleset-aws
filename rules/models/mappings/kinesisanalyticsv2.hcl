import = "api-models-aws/models/kinesis-analytics-v2/service/2018-05-23/kinesis-analytics-v2-2018-05-23.json"

mapping "aws_kinesisanalyticsv2_application" {
  name = ApplicationName
  runtime_environment = RuntimeEnvironment
  service_execution_role = RoleARN
  application_configuration = ApplicationConfiguration
  cloudwatch_logging_options = CloudWatchLoggingOptions
  description = ApplicationDescription
  tags = listmap(Tags, TagKey, TagValue)
}

mapping "aws_kinesisanalyticsv2_application_snapshot" {
  application_name = ApplicationName
  snapshot_name = SnapshotName
}
