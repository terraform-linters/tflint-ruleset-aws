import = "aws-sdk-go/models/apis/logs/2014-03-28/api-2.json"

mapping "aws_cloudwatch_log_destination" {
  name = DestinationName
}

mapping "aws_cloudwatch_log_destination_policy" {
  destination_name = DestinationName
}

mapping "aws_cloudwatch_log_group" {
  name       = LogGroupName
  kms_key_id = KmsKeyId
}

mapping "aws_cloudwatch_log_metric_filter" {
  name           = FilterName
  pattern        = FilterPattern
  log_group_name = LogGroupName
}

mapping "aws_cloudwatch_log_resource_policy" {
  policy_document = PolicyDocument
}

mapping "aws_cloudwatch_log_stream" {
  name           = LogStreamName
  log_group_name = LogGroupName
}

mapping "aws_cloudwatch_log_subscription_filter" {
  name           = FilterName
  filter_pattern = FilterPattern
  log_group_name = LogGroupName
  distribution   = Distribution
}

test "aws_cloudwatch_log_destination" "name" {
  valid   = ["test_destination"]
  invalid = ["test:destination"]
}

test "aws_cloudwatch_log_group" "name" {
  valid   = ["Yada"]
  invalid = ["Yoda:prod"]
}

test "aws_cloudwatch_log_metric_filter" "name" {
  valid   = ["MyAppAccessCount"]
  invalid = ["MyAppAccessCount:prod"]
}

test "aws_cloudwatch_log_stream" "name" {
  valid   = ["Yada"]
  invalid = ["Yoda:prod"]
}

test "aws_cloudwatch_log_subscription_filter" "name" {
  valid   = ["test_lambdafunction_logfilter"]
  invalid = ["test_lambdafunction_logfilter:test"]
}

test "aws_cloudwatch_log_subscription_filter" "distribution" {
  valid   = ["Random"]
  invalid = ["LogStream"]
}
