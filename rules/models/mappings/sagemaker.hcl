import = "api-models-aws/models/sagemaker/service/2017-07-24/sagemaker-2017-07-24.json"

mapping "aws_sagemaker_app" {
  app_name = AppName
  app_type = AppType
  domain_id = DomainId
  user_profile_name = UserProfileName
  resource_spec = ResourceSpec
  tags = TagList
}

mapping "aws_sagemaker_app_image_config" {
  app_image_config_name = CreateAppImageConfigRequest
  kernel_gateway_image_config = KernelGatewayImageConfig
  tags = TagList
}

mapping "aws_sagemaker_code_repository" {
  code_repository_name = EntityName
  git_config = GitConfig
  tags = TagList
}

mapping "aws_sagemaker_device_fleet" {
  device_fleet_name = EntityName
  role_arn = RoleArn
  output_config = EdgeOutputConfig
  description = DeviceFleetDescription
  enable_iot_role_alias = EnableIotRoleAlias
  tags = TagList
}

mapping "aws_sagemaker_domain" {
  domain_name = DomainName
  auth_mode = AuthMode
  vpc_id = VpcId
  subnet_ids = Subnets
  default_user_settings = UserSettings
  retention_policy = RetentionPolicy
  kms_key_id = KmsKeyId
  app_network_access_type = AppNetworkAccessType
  tags = TagList
}

mapping "aws_sagemaker_endpoint" {
  endpoint_config_name = EndpointConfigName
  name                 = EndpointName
  tags                 = TagList
}

mapping "aws_sagemaker_endpoint_configuration" {
  production_variants = ProductionVariantList
  kms_key_arn         = KmsKeyId
  name                = EndpointConfigName
  tags                = TagList
}

mapping "aws_sagemaker_feature_group" {
  feature_group_name = FeatureGroupName
  record_identifier_feature_name = FeatureName
  event_time_feature_name = FeatureName
  description = Description
  role_arn = RoleArn
  feature_definition = FeatureDefinitions
  offline_store_config = OfflineStoreConfig
  online_store_config = OnlineStoreConfig
  tags = TagList
}

mapping "aws_sagemaker_flow_definition" {
  flow_definition_name = FlowDefinitionName
  human_loop_config = HumanLoopConfig
  role_arn = RoleArn
  output_config = FlowDefinitionOutputConfig
  human_loop_activation_config = HumanLoopActivationConfig
  human_loop_request_source = HumanLoopRequestSource
  tags = TagList
}

mapping "aws_sagemaker_human_task_ui" {
  human_task_ui_name = HumanTaskUiName
  tags = TagList
  ui_template = UiTemplate
}

mapping "aws_sagemaker_image" {
  image_name = ImageName
  role_arn = RoleArn
  display_name = ImageDisplayName
  description = ImageDescription
  tags = TagList
}

mapping "aws_sagemaker_image_version" {
  image_name = ImageName
  base_image = ImageBaseImage
}

mapping "aws_sagemaker_model" {
  name                     = ModelName
  primary_container        = ContainerDefinition
  execution_role_arn       = RoleArn
  container                = ContainerDefinitionList
  enable_network_isolation = Boolean
  vpc_config               = VpcConfig
  tags                     = TagList
}

mapping "aws_sagemaker_model_package_group" {
  model_package_group_name = EntityName
  model_package_group_description = EntityDescription
  tags = TagList
}

mapping "aws_sagemaker_model_package_group_policy" {
  model_package_group_name = EntityName
}

mapping "aws_sagemaker_notebook_instance" {
  name                  = NotebookInstanceName
  role_arn              = RoleArn
  instance_type         = InstanceType
  subnet_id             = SubnetId
  security_groups       = SecurityGroupIds 
  kms_key_id            = KmsKeyId
  lifecycle_config_name = NotebookInstanceLifecycleConfigName
  tags                  = TagList
}

mapping "aws_sagemaker_notebook_instance_lifecycle_configuration" {
  name      = NotebookInstanceLifecycleConfigName
  on_create = NotebookInstanceLifecycleConfigList
  on_start  = NotebookInstanceLifecycleConfigList
}

mapping "aws_sagemaker_studio_lifecycle_config" {
  studio_lifecycle_config_name = StudioLifecycleConfigName
  studio_lifecycle_config_app_type = StudioLifecycleConfigAppType
  studio_lifecycle_config_content = StudioLifecycleConfigContent
  tags = TagList
}

mapping "aws_sagemaker_user_profile" {
  user_profile_name = UserProfileName
  domain_id = DomainId
  single_sign_on_user_identifier = SingleSignOnUserIdentifier
  single_sign_on_user_value = String256
  user_settings = UserSettings
  tags = TagList
}

mapping "aws_sagemaker_workforce" {
  workforce_name = WorkforceName
  cognito_config = CognitoConfig
  oidc_config = OidcConfig
  source_ip_config = SourceIpConfig
}

mapping "aws_sagemaker_workteam" {
  description = String200
  workforce_name = WorkforceName
  workteam_name = WorkteamName
  member_definition = MemberDefinitions
  notification_configuration = NotificationConfiguration
  tags = TagList
}
