import = "aws-sdk-go/models/apis/application-autoscaling/2016-02-06/api-2.json"

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
  valid   = ["StepScaling"]
  invalid = ["StopScaling"]
}

test "aws_appautoscaling_policy" "scalable_dimension" {
  valid   = ["ecs:service:DesiredCount"]
  invalid = ["ecs:service:DesireCount"]
}

test "aws_appautoscaling_policy" "service_namespace" {
  valid   = ["ecs"]
  invalid = ["eks"]
}
