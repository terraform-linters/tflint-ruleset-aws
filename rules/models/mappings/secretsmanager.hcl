import = "api-models-aws/models/secrets-manager/service/2017-10-17/secrets-manager-2017-10-17.json"

mapping "aws_secretsmanager_secret" {
  name                    = NameType
  name_prefix             = any
  description             = DescriptionType
  kms_key_id              = KmsKeyIdType
  policy                  = NonEmptyResourcePolicyType
  recovery_window_in_days = RecoveryWindowInDaysType
  tags                    = TagListType
}

mapping "aws_secretsmanager_secret_policy" {
  policy = NonEmptyResourcePolicyType
  secret_arn = SecretIdType
}

mapping "aws_secretsmanager_secret_rotation" {
  secret_id = SecretIdType
  rotation_lambda_arn = RotationLambdaARNType
  rotation_rules = RotationRulesType
}

mapping "aws_secretsmanager_secret_version" {
  secret_id      = SecretIdType
  secret_string  = SecretStringType
  secret_binary  = SecretBinaryType
  version_stages = SecretVersionStagesType
}
