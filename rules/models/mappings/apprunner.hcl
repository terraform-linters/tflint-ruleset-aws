import = "aws-sdk-go/models/apis/apprunner/2020-05-15/api-2.json"

mapping "aws_apprunner_auto_scaling_configuration_version" {
  auto_scaling_configuration_name = AutoScalingConfigurationName
  max_concurrency = ASConfigMaxConcurrency
  max_size = ASConfigMaxSize
  min_size = ASConfigMinSize
  tags = TagList
}

mapping "aws_apprunner_connection" {
  connection_name = ConnectionName
  provider_type = ProviderType
  tags = TagList
}

mapping "aws_apprunner_custom_domain_association" {
  domain_name = DomainName
  # service_arn = AppRunnerResourceArn # https://github.com/golang/go/issues/7252
}

mapping "aws_apprunner_service" {
  service_name = ServiceName
  source_configuration = SourceConfiguration
  # auto_scaling_configuration_arn = AppRunnerResourceArn # https://github.com/golang/go/issues/7252
  encryption_configuration = EncryptionConfiguration
  health_check_configuration = HealthCheckConfiguration
  instance_configuration = InstanceConfiguration
  tags = TagList
}
