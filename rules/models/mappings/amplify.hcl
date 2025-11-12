import = "api-models-aws/models/amplify/service/2017-07-25/amplify-2017-07-25.json"

mapping "aws_amplify_app" {
  name = Name
  access_token = AccessToken
  auto_branch_creation_config = AutoBranchCreationConfig
  auto_branch_creation_patterns = AutoBranchCreationPatterns
  basic_auth_credentials = BasicAuthCredentials
  build_spec = BuildSpec
  custom_rule = CustomRules
  description = Description
  enable_auto_branch_creation = EnableAutoBranchCreation
  enable_basic_auth = EnableBasicAuth
  enable_branch_auto_build = EnableBranchAutoBuild
  enable_branch_auto_deletion = EnableBranchAutoDeletion
  environment_variables = EnvironmentVariables
  iam_service_role_arn = ServiceRoleArn
  oauth_token = OauthToken
  platform = Platform
  repository = Repository
  tags = TagMap
}

mapping "aws_amplify_backend_environment" {
  app_id = AppId
  environment_name = EnvironmentName
  deployment_artifacts = DeploymentArtifacts
  stack_name = StackName
}

mapping "aws_amplify_branch" {
  app_id = AppId
  branch_name = BranchName
  backend_environment_arn = BackendEnvironmentArn
  basic_auth_credentials = BasicAuthCredentials
  description = Description
  display_name = DisplayName
  enable_auto_build = EnableAutoBuild
  enable_basic_auth = EnableBasicAuth
  enable_notification = EnableNotification
  enable_performance_mode = EnablePerformanceMode
  enable_pull_request_preview = EnablePullRequestPreview
  environment_variables = EnvironmentVariables
  framework = Framework
  pull_request_environment_name = PullRequestEnvironmentName
  stage = Stage
  tags = TagMap
  ttl = TTL
}

mapping "aws_amplify_domain_association" {
  app_id = AppId
  domain_name = any //DomainName
}

mapping "aws_amplify_webhook" {
  app_id = AppId
  branch_name = BranchName
  description = Description
}
