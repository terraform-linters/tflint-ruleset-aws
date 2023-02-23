import = "aws-sdk-go/models/apis/monitoring/2010-08-01/api-2.json"

mapping "aws_cloudwatch_metric_alarm" {
  alarm_name                            = AlarmName
  comparison_operator                   = ComparisonOperator
  metric_name                           = MetricName
  namespace                             = Namespace
  statistic                             = Statistic
  alarm_description                     = AlarmDescription
  unit                                  = StandardUnit
  extended_statistic                    = ExtendedStatistic
  treat_missing_data                    = TreatMissingData
  evaluate_low_sample_count_percentiles = EvaluateLowSampleCountPercentile
}

test "aws_cloudwatch_metric_alarm" "comparison_operator" {
  valid   = ["GreaterThanOrEqualToThreshold"]
  invalid = ["GreaterThanOrEqual"]
}

test "aws_cloudwatch_metric_alarm" "namespace" {
  valid   = ["AWS/EC2"]
  invalid = [":EC2"]
}

test "aws_cloudwatch_metric_alarm" "statistic" {
  valid   = ["Average"]
  invalid = ["Median"]
}

test "aws_cloudwatch_metric_alarm" "unit" {
  valid   = ["Gigabytes"]
  invalid = ["GB"]
}

test "aws_cloudwatch_metric_alarm" "extended_statistic" {
  valid   = ["p100"]
  invalid = ["p101"]
}
