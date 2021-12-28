import = "aws-sdk-go/models/apis/events/2015-10-07/api-2.json"

mapping "aws_cloudwatch_event_api_destination" {
  name = ApiDestinationName
  description = ApiDestinationDescription
  invocation_endpoint = HttpsEndpoint
  http_method = ApiDestinationHttpMethod
  invocation_rate_limit_per_second = ApiDestinationInvocationRateLimitPerSecond
  connection_arn = ConnectionArn
}

mapping "aws_cloudwatch_event_archive" {
  name = ArchiveName
  event_source_arn = Arn
  description = ArchiveDescription
  event_pattern = EventPattern
  retention_days = RetentionDays
}

mapping "aws_cloudwatch_event_bus" {
  name = EventBusName
  event_source_name = EventSourceName
  tags = TagList
}

mapping "aws_cloudwatch_event_bus_policy" {
  event_bus_name = EventBusName
}

mapping "aws_cloudwatch_event_connection" {
  name = ConnectionName
  description = ConnectionDescription
  authorization_type = ConnectionAuthorizationType
  auth_parameters = CreateConnectionAuthRequestParameters
}

mapping "aws_cloudwatch_event_permission" {
  principal    = Principal
  statement_id = StatementId
  action       = Action
}

mapping "aws_cloudwatch_event_rule" {
  name                = RuleName
  schedule_expression = ScheduleExpression
  description         = RuleDescription
  role_arn            = RoleArn
}

mapping "aws_cloudwatch_event_target" {
  rule       = RuleName
  target_id  = TargetId
  arn        = TargetArn
  input      = TargetInput
  input_path = TargetInputPath
  role_arn   = RoleArn
}

test "aws_cloudwatch_event_permission" "principal" {
  ok = "*"
  ng = "-"
}

test "aws_cloudwatch_event_permission" "statement_id" {
  ok = "OrganizationAccess"
  ng = "Organization Access"
}

test "aws_cloudwatch_event_permission" "action" {
  ok = "events:PutEvents"
  ng = "cloudwatchevents:PutEvents"
}

test "aws_cloudwatch_event_rule" "name" {
  ok = "capture-aws-sign-in"
  ng = "capture aws sign in"
}

test "aws_cloudwatch_event_target" "target_id" {
  ok = "run-scheduled-task-every-hour"
  ng = "run scheduled task every hour"
}
