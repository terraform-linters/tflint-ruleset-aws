import = "aws-sdk-ruby/apis/cloudfront/2020-05-31/api-2.json"

mapping "aws_cloudfront_cache_policy" {
  parameters_in_cache_key_and_forwarded_to_origin = ParametersInCacheKeyAndForwardedToOrigin
}

mapping "aws_cloudfront_distribution" {
  aliases = Aliases
  comment = CommentType
  custom_error_response = CustomErrorResponses
  default_cache_behavior = DefaultCacheBehavior
  http_version = HttpVersion
  logging_config = LoggingConfig
  ordered_cache_behavior = CacheBehaviors
  origin = Origins
  origin_group = OriginGroups
  price_class = PriceClass
  restrictions = Restrictions
  tags = Tags
  viewer_certificate = ViewerCertificate
}

mapping "aws_cloudfront_function" {
  name = FunctionName
  code = FunctionBlob
  runtime = FunctionRuntime
}

mapping "aws_cloudfront_key_group" {
  items = PublicKeyIdList
}

mapping "aws_cloudfront_realtime_log_config" {
  endpoint = EndPointList
  fields = FieldList
}

mapping "aws_cloudfront_response_headers_policy" {
  cors_config = ResponseHeadersPolicyCorsConfig
  custom_headers_config = ResponseHeadersPolicyCustomHeadersConfig
  security_headers_config = ResponseHeadersPolicySecurityHeadersConfig
}

test "aws_cloudfront_distribution" "http_version" {
  ok = "http2"
  ng = "http1.2"
}

test "aws_cloudfront_distribution" "price_class" {
  ok = "PriceClass_All"
  ng = "PriceClass_300"
}
