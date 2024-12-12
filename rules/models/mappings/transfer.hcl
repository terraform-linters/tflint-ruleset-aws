import = "aws-sdk-ruby/apis/transfer/2018-11-05/api-2.json"

mapping "aws_transfer_access" {
  external_id = ExternalId
  server_id = ServerId
  home_directory = HomeDirectory
  home_directory_mappings = HomeDirectoryMappings
  home_directory_type = HomeDirectoryType
  policy = Policy
  posix_profile = PosixProfile
  role = Role
}

mapping "aws_transfer_server" {
  endpoint_details       = EndpointDetails
  endpoint_type          = EndpointType
  invocation_role        = Role
  url                    = Url
  identity_provider_type = IdentityProviderType
  logging_role           = Role
  force_destroy          = any
  tags                   = Tags
}

mapping "aws_transfer_ssh_key" {
  server_id = ServerId
  user_name = UserName
  body      = SshPublicKeyBody
}

mapping "aws_transfer_user" {
  server_id      = ServerId
  user_name      = UserName
  home_directory = HomeDirectory
  policy         = Policy
  role           = Role
  tags           = Tags
}
