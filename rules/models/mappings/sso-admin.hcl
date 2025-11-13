import = "api-models-aws/models/sso-admin/service/2020-07-20/sso-admin-2020-07-20.json"

mapping "aws_ssoadmin_account_assignment" {
  instance_arn = InstanceArn
  permission_set_arn = PermissionSetArn
  principal_id = PrincipalId
  principal_type = PrincipalType
  target_id = TargetId
  target_type = TargetType
}

mapping "aws_ssoadmin_managed_policy_attachment" {
  instance_arn = InstanceArn
  managed_policy_arn = ManagedPolicyArn
  permission_set_arn = PermissionSetArn
}

mapping "aws_ssoadmin_permission_set" {
  description = PermissionSetDescription
  instance_arn = InstanceArn
  name = PermissionSetName
  relay_state = RelayState
  # session_duration = Duration
  tags = TagList
}

mapping "aws_ssoadmin_permission_set_inline_policy" {
  inline_policy = PermissionSetPolicyDocument
  instance_arn = InstanceArn
  permission_set_arn = PermissionSetArn
}
