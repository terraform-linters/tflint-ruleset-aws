import = "aws-sdk-ruby/apis/appconfig/2019-10-09/api-2.json"

mapping "aws_appconfig_application" {
  name = Name
  description = Description
  tags = TagMap
}

mapping "aws_appconfig_configuration_profile" {
  application_id = Id
  location_uri = Uri
  name = Name
  description = Description
  retrieval_role_arn = RoleArn
  tags = TagMap
  validator = ValidatorList
}

mapping "aws_appconfig_deployment" {
  application_id = Id
  configuration_profile_id = Id
  configuration_version = Version
  deployment_strategy_id = DeploymentStrategyId
  description = Description
  environment_id = Id
  tags = TagMap
}

mapping "aws_appconfig_deployment_strategy" {
  deployment_duration_in_minutes = MinutesBetween0And24Hours
  growth_factor = GrowthFactor
  name = Name
  replicate_to = ReplicateTo
  description = Description
  final_bake_time_in_minutes = MinutesBetween0And24Hours
  growth_type = GrowthType
  tags = TagMap
}

mapping "aws_appconfig_environment" {
  application_id = Id
  name = Name
  description = Description
  monitor = MonitorList
  tags = TagMap
}

mapping "aws_appconfig_hosted_configuration_version" {
  application_id = Id
  configuration_profile_id = Id
  content = Blob
  content_type = StringWithLengthBetween1And255
  description = Description
}
