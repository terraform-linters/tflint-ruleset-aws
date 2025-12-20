import = "api-models-aws/models/data-pipeline/service/2012-10-29/data-pipeline-2012-10-29.json"

mapping "aws_datapipeline_pipeline" {
  name        = any // id
  description = any // string
  tags        = listmap(tagList, tagKey, tagValue)
}
