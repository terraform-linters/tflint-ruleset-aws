import = "aws-sdk-ruby/apis/application-autoscaling/2016-02-06/api-2.json"

mapping "aws_appautoscaling_policy" {
  policy_type        = PolicyType
  scalable_dimension = ScalableDimension
  service_namespace  = ServiceNamespace
}

mapping "aws_appautoscaling_scheduled_action" {
  scalable_dimension = ScalableDimension
  service_namespace = ServiceNamespace
}

mapping "aws_appautoscaling_target" {
  scalable_dimension = ScalableDimension
  service_namespace  = ServiceNamespace
}

test "aws_appautoscaling_policy" "policy_type" {
  ok = "StepScaling"
  ng = "StopScaling"
}

test "aws_appautoscaling_policy" "scalable_dimension" {
  ok = "ecs:service:DesiredCount"
  ng = "ecs:service:DesireCount"
}

test "aws_appautoscaling_policy" "service_namespace" {
  ok = "ecs"
  ng = "eks"
}
