import = "aws-sdk-go/models/apis/securityhub/2018-10-26/api-2.json"

mapping "aws_securityhub_account" {}

mapping "aws_securityhub_action_target" {
  name = NonEmptyString
  identifier = NonEmptyString
  description = NonEmptyString
}

mapping "aws_securityhub_finding_aggregator" {
  linking_mode = NonEmptyString
  specified_regions = StringList
}

mapping "aws_securityhub_insight" {
  filters = AwsSecurityFindingFilters
  group_by_attribute = NonEmptyString
  name = NonEmptyString
}

mapping "aws_securityhub_invite_accepter" {
  master_id = NonEmptyString
}

mapping "aws_securityhub_member" {
  account_id = AccountId
  email = NonEmptyString
}

mapping "aws_securityhub_organization_admin_account" {
  admin_account_id = NonEmptyString
}

mapping "aws_securityhub_product_subscription" {
  product_arn = NonEmptyString
}

mapping "aws_securityhub_standards_control" {
  standards_control_arn = NonEmptyString
  control_status = ControlStatus
  disabled_reason = NonEmptyString
}

mapping "aws_securityhub_standards_subscription" {
  standards_arn = NonEmptyString
}
