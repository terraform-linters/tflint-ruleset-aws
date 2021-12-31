import = "aws-sdk-go/models/apis/ecr-public/2020-10-30/api-2.json"

mapping "aws_ecrpublic_repository" {
  repository_name = RepositoryName
  catalog_data = RepositoryCatalogDataInput
}

mapping "aws_ecrpublic_repository_policy" {
  repository_name = RepositoryName
  policy = RepositoryPolicyText
}
