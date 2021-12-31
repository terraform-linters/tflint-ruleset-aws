import = "aws-sdk-go/models/apis/schemas/2019-12-02/api-2.json"

mapping "aws_schemas_discoverer" {
  source_arn = __stringMin20Max1600
  description = __stringMin0Max256
  tags = Tags
}

mapping "aws_schemas_registry" { 
  description = __stringMin0Max256
  tags = Tags
}

mapping "aws_schemas_schema" {
  content = __stringMin1Max100000
  type = Type
  description = __stringMin0Max256
  tags = Tags
}
