import = "api-models-aws/models/mediastore/service/2017-09-01/mediastore-2017-09-01.json"

mapping "aws_media_store_container" {
  name = ContainerName
}

mapping "aws_media_store_container_policy" {
  container_name = ContainerName
  policy         = any
}
