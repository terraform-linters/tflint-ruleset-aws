import = "aws-sdk-ruby/apis/apigatewayv2/2018-11-29/api-2.json"

mapping "aws_apigatewayv2_api" {
  name = StringWithLengthBetween1And128
  protocol_type = ProtocolType
  api_key_selection_expression = SelectionExpression
  cors_configuration = Cors
  credentials_arn = Arn
  description = StringWithLengthBetween0And1024
  route_key = SelectionKey
  route_selection_expression = SelectionExpression
  tags = Tags
  target = UriWithLengthBetween1And2048
  version = StringWithLengthBetween1And64
}

mapping "aws_apigatewayv2_api_mapping" {
  api_id      = Id
  stage       = StringWithLengthBetween1And128
  api_mapping_key = SelectionKey
}

mapping "aws_apigatewayv2_authorizer" {
  authorizer_type = AuthorizerType
  name = StringWithLengthBetween1And128
  authorizer_credentials_arn = Arn
  authorizer_payload_format_version = StringWithLengthBetween1And64
  authorizer_result_ttl_in_seconds = IntegerWithLengthBetween0And3600
  authorizer_uri = UriWithLengthBetween1And2048
  identity_sources = IdentitySourceList
  jwt_configuration = JWTConfiguration
}

mapping "aws_apigatewayv2_deployment" {
  description = StringWithLengthBetween0And1024
}

mapping "aws_apigatewayv2_domain_name" {
  domain_name = StringWithLengthBetween1And512
  mutual_tls_authentication = MutualTlsAuthenticationInput
  tags = Tags
}

mapping "aws_apigatewayv2_integration" {
  integration_type = IntegrationType
  connection_id = StringWithLengthBetween1And1024
  connection_type = ConnectionType
  content_handling_strategy = ContentHandlingStrategy
  credentials_arn = Arn
  description = StringWithLengthBetween0And1024
  integration_method = StringWithLengthBetween1And64
  integration_subtype = StringWithLengthBetween1And128
  integration_uri = UriWithLengthBetween1And2048
  passthrough_behavior = PassthroughBehavior
  payload_format_version = StringWithLengthBetween1And64
  request_parameters = IntegrationParameters
  request_templates = TemplateMap
  response_parameters = ResponseParameters
  template_selection_expression = SelectionExpression
  timeout_milliseconds = IntegerWithLengthBetween50And30000
  tls_config = TlsConfigInput
}

mapping "aws_apigatewayv2_integration_response" {
  integration_response_key = SelectionKey
  content_handling_strategy = ContentHandlingStrategy
  response_templates = TemplateMap
  template_selection_expression = SelectionExpression
}

mapping "aws_apigatewayv2_model" {
  content_type = StringWithLengthBetween1And256
  name = StringWithLengthBetween1And128
  schema = StringWithLengthBetween0And32K
  description = StringWithLengthBetween0And1024
}

mapping "aws_apigatewayv2_route" {
  route_key = SelectionKey
  authorization_scopes = AuthorizationScopes
  authorization_type = AuthorizationType
  authorizer_id = Id
  model_selection_expression = SelectionExpression
  operation_name = StringWithLengthBetween1And64
  request_models = RouteModels
  request_parameter = RouteParameters
  route_response_selection_expression = SelectionExpression
  target = StringWithLengthBetween1And128
}

mapping "aws_apigatewayv2_route_response" {
  route_response_key = SelectionKey
  model_selection_expression = SelectionExpression
  response_models = RouteModels
}

mapping "aws_apigatewayv2_stage" {
  name = StringWithLengthBetween1And128
  access_log_settings = AccessLogSettings
  client_certificate_id = Id
  default_route_settings = RouteSettings
  deployment_id = Id
  description = StringWithLengthBetween0And1024
  route_settings = RouteSettingsMap
  stage_variables = StageVariablesMap
  tags = Tags
}

mapping "aws_apigatewayv2_vpc_link" {
  name = StringWithLengthBetween1And128
  security_group_ids = SecurityGroupIdList
  subnet_ids = SubnetIdList
  tags = Tags
}
