import = "aws-sdk-ruby/apis/imagebuilder/2019-12-02/api-2.json"

mapping "aws_imagebuilder_component" {
  name = ResourceName
  platform = Platform
  version = VersionNumber
  change_description = NonEmptyString
  data = InlineComponentData
  description = NonEmptyString
  kms_key_id = NonEmptyString
  supported_os_versions = OsVersionList
  tags = TagMap
  uri = Uri
}

mapping "aws_imagebuilder_distribution_configuration" {
  name = ResourceName
  distribution = DistributionList
  description = NonEmptyString
  # kms_key_id = NonEmptyString
  tags = TagMap
}

mapping "aws_imagebuilder_image" {
  image_recipe_arn = ImageRecipeArn
  infrastructure_configuration_arn = InfrastructureConfigurationArn
  distribution_configuration_arn = DistributionConfigurationArn
  image_tests_configuration = ImageTestsConfiguration
  tags = TagMap
}

mapping "aws_imagebuilder_image_pipeline" {
  image_recipe_arn = ImageRecipeArn
  infrastructure_configuration_arn = InfrastructureConfigurationArn
  name = ResourceName
  description = NonEmptyString
  distribution_configuration_arn = DistributionConfigurationArn
  image_tests_configuration = ImageTestsConfiguration
  schedule = Schedule
  status = PipelineStatus
  tags = TagMap
}

mapping "aws_imagebuilder_image_recipe" {
  component = ComponentConfigurationList
  name = ResourceName
  parent_image = NonEmptyString
  version = VersionNumber
  block_device_mapping = InstanceBlockDeviceMappings
  description = NonEmptyString
  tags = TagMap
  working_directory = NonEmptyString
}

mapping "aws_imagebuilder_infrastructure_configuration" {
  instance_profile_name = InstanceProfileNameType
  name = ResourceName
  description = NonEmptyString
  instance_types = InstanceTypeList
  key_pair = NonEmptyString
  logging = Logging
  resource_tags = ResourceTagMap
  security_group_ids = SecurityGroupIds
  sns_topic_arn = SnsTopicArn
  subnet_id = NonEmptyString
  tags = TagMap
}
