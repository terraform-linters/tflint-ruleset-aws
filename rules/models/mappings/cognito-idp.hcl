import = "aws-sdk-go/models/apis/cognito-idp/2016-04-18/api-2.json"

mapping "aws_cognito_identity_provider" {
  user_pool_id  = UserPoolIdType
  provider_name = ProviderNameType
  provider_type = IdentityProviderTypeType
}

mapping "aws_cognito_resource_server" {
  identifier = ResourceServerIdentifierType
  name       = ResourceServerNameType
}

mapping "aws_cognito_user_group" {
  name         = GroupNameType
  user_pool_id = UserPoolIdType
  description  = DescriptionType
  precedence   = PrecedenceType
  role_arn     = ArnType
}

mapping "aws_cognito_user_pool" {
  alias_attributes           = AliasAttributesListType
  auto_verified_attributes   = VerifiedAttributesListType
  name                       = UserPoolNameType
  email_verification_subject = EmailVerificationSubjectType
  email_verification_message = EmailVerificationMessageType
  mfa_configuration          = UserPoolMfaType
  sms_authentication_message = SmsVerificationMessageType
  sms_verification_message   = SmsVerificationMessageType
}

mapping "aws_cognito_user_pool_client" {
  default_redirect_uri   = RedirectUrlType
  name                   = ClientNameType
  refresh_token_validity = RefreshTokenValidityType
  user_pool_id           = UserPoolIdType
}

mapping "aws_cognito_user_pool_domain" {
  domain          = any // DomainType is not appropriate for a fully-customized domain. See also https://github.com/terraform-linters/tflint-ruleset-aws/issues/156
  user_pool_id    = UserPoolIdType
  certificate_arn = ArnType
}
  
mapping "aws_cognito_user_pool_ui_customization" {
  client_id = ClientIdType
  css = CSSType
  image_file = ImageFileType
  user_pool_id = UserPoolIdType
}

test "aws_cognito_identity_provider" "user_pool_id" {
  valid   = ["foo_bar"]
  invalid = ["foobar"]
}

test "aws_cognito_identity_provider" "provider_name" {
  valid   = ["Google"]
  invalid = ["\t"]
}

test "aws_cognito_identity_provider" "provider_type" {
  valid   = ["LoginWithAmazon"]
  invalid = ["Apple"]
}

test "aws_cognito_resource_server" "identifier" {
  valid   = ["https://example.com"]
  invalid = ["\t"]
}

test "aws_cognito_resource_server" "name" {
  valid   = ["example"]
  invalid = ["example/server"]
}

test "aws_cognito_user_group" "name" {
  valid   = ["user-group"]
  invalid = ["user\tgroup"]
}

test "aws_cognito_user_group" "role_arn" {
  valid   = ["arn:aws:iam::123456789012:role/s3access"]
  invalid = ["aws:iam::123456789012:instance-profile/s3access-profile"]
}

test "aws_cognito_user_pool" "name" {
  valid   = ["mypool"]
  invalid = ["my/pool"]
}

test "aws_cognito_user_pool" "email_verification_message" {
  valid   = ["Verification code is {####}"]
  invalid = ["Verification code"]
}

test "aws_cognito_user_pool" "mfa_configuration" {
  valid   = ["ON"]
  invalid = ["IN"]
}

test "aws_cognito_user_pool" "sms_authentication_message" {
  valid   = ["Authentication code is {####}"]
  invalid = ["Authentication code"]
}

test "aws_cognito_user_pool" "sms_verification_message" {
  valid   = ["Verification code is {####}"]
  invalid = ["Verification code"]
}

test "aws_cognito_user_pool_client" "default_redirect_uri" {
  valid   = ["https://example.com/callback"]
  invalid = ["https://example com"]
}

test "aws_cognito_user_pool_client" "name" {
  valid   = ["client"]
  invalid = ["client/example"]
}
