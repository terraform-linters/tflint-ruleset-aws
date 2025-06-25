import = "api-models-aws/models/dlm/service/2018-01-12/dlm-2018-01-12.json"

mapping "aws_dlm_lifecycle_policy" {
  description        = PolicyDescription
  execution_role_arn = ExecutionRoleArn
  policy_details     = PolicyDetails
  state              = SettablePolicyStateValues
}

test "aws_dlm_lifecycle_policy" "state" {
  ok = "ENABLED"
  ng = "ERROR"
}
