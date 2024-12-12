import = "aws-sdk-ruby/apis/swf/2012-01-25/api-2.json"

mapping "aws_swf_domain" {
  name                                        = DomainName
  name_prefix                                 = any
  description                                 = Description
  workflow_execution_retention_period_in_days = DurationInDays
}
