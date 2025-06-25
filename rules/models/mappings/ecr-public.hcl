import = "api-models-aws/models/ecr-public/service/2020-10-30/ecr-public-2020-10-30.json"

mapping "aws_ecrpublic_repository" {
  repository_name = RepositoryName
  catalog_data = RepositoryCatalogDataInput
}

mapping "aws_ecrpublic_repository_policy" {
  repository_name = RepositoryName
  policy = RepositoryPolicyText
}
