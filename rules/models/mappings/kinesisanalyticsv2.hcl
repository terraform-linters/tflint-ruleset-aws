import = "aws-sdk-go/models/apis/kinesisanalyticsv2/2018-05-23/api-2.json"

mapping "aws_kinesisanalyticsv2_application" {
  name = ApplicationName
  runtime_environment = RuntimeEnvironment
  service_execution_role = RoleARN
  application_configuration = ApplicationConfiguration
  cloudwatch_logging_options = CloudWatchLoggingOptions
  description = ApplicationDescription
  tags = Tags
}

mapping "aws_kinesisanalyticsv2_application_snapshot" {
  application_name = ApplicationName
  snapshot_name = SnapshotName
}
