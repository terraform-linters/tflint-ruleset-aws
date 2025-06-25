import = "api-models-aws/models/swf/service/2012-01-25/swf-2012-01-25.json"

mapping "aws_swf_domain" {
  name                                        = DomainName
  name_prefix                                 = any
  description                                 = Description
  workflow_execution_retention_period_in_days = DurationInDays
}
