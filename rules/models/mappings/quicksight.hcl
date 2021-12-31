import = "aws-sdk-go/models/apis/quicksight/2018-04-01/api-2.json"

mapping "aws_quicksight_data_source" {
  data_source_id = ResourceId
  name = ResourceName
  parameters = DataSourceParameters
  type = DataSourceType
  aws_account_id = AwsAccountId
  credentials = DataSourceCredentials
  permission = ResourcePermissionList
  ssl_properties = SslProperties
  tags = TagList
  vpc_connection_properties = VpcConnectionProperties
}

mapping "aws_quicksight_group" {
  group_name     = GroupName
  aws_account_id = AwsAccountId
  description    = GroupDescription
  namespace      = Namespace
}

mapping "aws_quicksight_group_membership" {
  group_name = GroupName
  member_name = GroupMemberName
  aws_account_id = AwsAccountId
  namespace = Namespace
}

mapping "aws_quicksight_user" {
  identity_type = IdentityType
  user_role = UserRole
  user_name = UserName
  aws_account_id = AwsAccountId
  namespace = Namespace
  session_name = RoleSessionName
}
