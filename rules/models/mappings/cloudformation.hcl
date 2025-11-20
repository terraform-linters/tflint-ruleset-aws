import = "api-models-aws/models/cloudformation/service/2010-05-15/cloudformation-2010-05-15.json"

mapping "aws_cloudformation_stack" {
  template_url = TemplateURL
  capabilities = Capabilities
  on_failure   = OnFailure
  policy_body  = StackPolicyBody
  policy_url   = StackPolicyURL
  iam_role_arn = RoleARN
}

mapping "aws_cloudformation_stack_set" {
  capabilities            = Capabilities
  description             = Description
  execution_role_name     = ExecutionRoleName
  template_url            = TemplateURL
  administration_role_arn = RoleARN
}

mapping "aws_cloudformation_stack_set_instance" {
  account_id = Account
}

test "aws_cloudformation_stack" "on_failure" {
  ok = "DO_NOTHING"
  ng = "DO_ANYTHING"
}

test "aws_cloudformation_stack_set" "execution_role_name" {
  ok = "AWSCloudFormationStackSetExecutionRole"
  ng = "AWSCloudFormation/StackSet/ExecutionRole"
}

test "aws_cloudformation_stack_set_instance" "account_id" {
  ok = "123456789012"
  ng = "1234567890123"
}
