import = "api-models-aws/models/ram/service/2018-01-04/ram-2018-01-04.json"

mapping "aws_ram_principal_association" {
  principal          = String
  resource_share_arn = String
}

mapping "aws_ram_resource_association" {
  resource_arn       = String
  resource_share_arn = String
}

mapping "aws_ram_resource_share" {
  name                      = String
  allow_external_principals = Boolean
  tags                      = TagList
}

mapping "aws_ram_resource_share_accepter" {
  share_arn = String
}
