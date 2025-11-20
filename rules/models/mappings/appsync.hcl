import = "api-models-aws/models/appsync/service/2017-07-25/appsync-2017-07-25.json"

mapping "aws_appsync_datasource" {
  name = ResourceName
  type = DataSourceType
}

mapping "aws_appsync_graphql_api" {
  authentication_type = AuthenticationType
}

mapping "aws_appsync_resolver" {
  type              = ResourceName
  field             = ResourceName
  data_source       = ResourceName
  request_template  = any //MappingTemplate
  response_template = any //MappingTemplate
}

mapping "aws_appsync_function" {
  api_id                    = String
  data_source               = ResourceName
  name                      = ResourceName
  request_mapping_template  = any //MappingTemplate
  response_mapping_template = any //MappingTemplate
  description               = String
  function_version          = String
}

test "aws_appsync_datasource" "name" {
  ok = "tf_appsync_example"
  ng = "01_tf_example"
}

test "aws_appsync_datasource" "type" {
  ok = "AWS_LAMBDA"
  ng = "AMAZON_SIMPLEDB"
}

test "aws_appsync_graphql_api" "authentication_type" {
  ok = "API_KEY"
  ng = "AWS_KEY"
}
