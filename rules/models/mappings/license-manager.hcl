import = "api-models-aws/models/license-manager/service/2018-08-01/license-manager-2018-08-01.json"

mapping "aws_licensemanager_association" {
  license_configuration_arn = String
  resource_arn              = String
}

mapping "aws_licensemanager_license_configuration" {
  name                     = String
  description              = String
  license_count            = BoxLong
  license_count_hard_limit = BoxBoolean
  license_counting_type    = LicenseCountingType
  license_rules            = StringList
  tags                     = TagList
}
