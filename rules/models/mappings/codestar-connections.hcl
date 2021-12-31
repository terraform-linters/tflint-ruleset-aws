import = "aws-sdk-go/models/apis/codestar-connections/2019-12-01/api-2.json"

mapping "aws_codestarconnections_connection" {
  name = ConnectionName
  provider_type = ProviderType
  host_arn = HostArn
  tags = TagList
}

mapping "aws_codestarconnections_host" {
  name = HostName
  provider_endpoint = Url
  provider_type = ProviderType
  vpc_configuration = VpcConfiguration
}
