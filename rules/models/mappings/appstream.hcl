import = "aws-sdk-go/models/apis/appstream/2016-12-01/api-2.json"

mapping "aws_appstream_directory_config" {
  directory_name = DirectoryName
  organizational_unit_distinguished_names = OrganizationalUnitDistinguishedNamesList
  service_account_credentials = ServiceAccountCredentials
}

mapping "aws_appstream_fleet" {
  compute_capacity = ComputeCapacity
  name = Name
  description = Description
  display_name = DisplayName
  domain_join_info = DomainJoinInfo
  fleet_type = FleetType
  iam_role_arn = Arn
  image_arn = Arn
  stream_view = StreamView
  vpc_config = VpcConfig
  tags = Tags
}

mapping "aws_appstream_image_builder" {
  name = Name
  access_endpoint = AccessEndpointList
  appstream_agent_version = AppstreamAgentVersion
  description = Description
  display_name = DisplayName
  domain_join_info = DomainJoinInfo
  iam_role_arn = Arn
  image_arn = Arn
  vpc_config = VpcConfig
  tags = Tags
}

mapping "aws_appstream_stack" {
  application_settings = ApplicationSettings
  description = Description
  display_name = DisplayName
  embed_host_domains = EmbedHostDomains
  feedback_url = FeedbackURL
  redirect_url = RedirectURL
  storage_connectors = StorageConnectorList
  user_settings = UserSettingList
}

mapping "aws_appstream_user" {
  authentication_type = AuthenticationType
  user_name = Username
  first_name = UserAttributeValue
  last_name = UserAttributeValue
}

mapping "aws_appstream_user_stack_association" {
  authentication_type = AuthenticationType
  user_name = Username
}
