import = "api-models-aws/models/glacier/service/2012-06-01/glacier-2012-06-01.json"

mapping "aws_glacier_vault" {
  name          = string
  access_policy = string
  notification  = VaultNotificationConfig
  tags          = TagMap
}

mapping "aws_glacier_vault_lock" {
  complete_lock         = boolean
  policy                = string
  vault_name            = string
  ignore_deletion_error = boolean
}
