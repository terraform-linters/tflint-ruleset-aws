import = "aws-sdk-go/models/apis/apigateway/2015-07-09/api-2.json"

mapping "aws_api_gateway_documentation_part" {
  location = DocumentationPartLocation
}

mapping "aws_api_gateway_domain_name" {
  endpoint_configuration           = EndpointConfiguration
  mutual_tls_authentication        = MutualTlsAuthenticationInput
  security_policy                  = SecurityPolicy
}

mapping "aws_api_gateway_gateway_response" {
  response_type = GatewayResponseType
  status_code = StatusCode
}

mapping "aws_api_gateway_integration_response" {
  status_code = StatusCode
  content_handling = ContentHandlingStrategy
}

mapping "aws_api_gateway_method_response" {
  status_code = StatusCode
}

mapping "aws_api_gateway_authorizer" {
  authorizer_result_ttl_in_seconds = NullableInteger
  provider_arns                    = ListOfARNs
  type                             = AuthorizerType
}

mapping "aws_api_gateway_integration" {
  type = IntegrationType
  connection_type  = ConnectionType
  content_handling = ContentHandlingStrategy
}

mapping "aws_api_gateway_rest_api" {
  api_key_source = ApiKeySourceType
}

mapping "aws_api_gateway_stage" {
  cache_cluster_size = CacheClusterSize
}

test "aws_api_gateway_gateway_response" "status_code" {
  valid   = ["200"]
  invalid = ["004"]
}

test "aws_api_gateway_authorizer" "type" {
  valid   = ["TOKEN"]
  invalid = ["RESPONSE"]
}

test "aws_api_gateway_gateway_response" "response_type" {
  valid   = ["UNAUTHORIZED"]
  invalid = ["4XX"]
}

test "aws_api_gateway_integration" "type" {
  valid   = ["HTTP"]
  invalid = ["AWS_HTTP"]
}

test "aws_api_gateway_integration" "connection_type" {
  valid   = ["INTERNET"]
  invalid = ["INTRANET"]
}

test "aws_api_gateway_integration" "content_handling" {
  valid   = ["CONVERT_TO_BINARY"]
  invalid = ["CONVERT_TO_FILE"]
}

test "aws_api_gateway_rest_api" "api_key_source" {
  valid   = ["AUTHORIZER"]
  invalid = ["BODY"]
}

test "aws_api_gateway_stage" "cache_cluster_size" {
  valid   = ["6.1"]
  invalid = ["6.2"]
}
