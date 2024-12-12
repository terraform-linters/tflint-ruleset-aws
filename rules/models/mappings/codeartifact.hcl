import = "aws-sdk-ruby/apis/codeartifact/2018-09-22/api-2.json"

mapping "aws_codeartifact_domain" {
  domain = DomainName
  encryption_key = Arn
  tags = TagList
}

mapping "aws_codeartifact_domain_permissions_policy" {
  domain = DomainName
  policy_document = PolicyDocument
  domain_owner = AccountId
  policy_revision = PolicyRevision
}

mapping "aws_codeartifact_repository" {
  domain = DomainName
  repository = RepositoryName
  domain_owner = AccountId
  description = Description
  upstream = UpstreamRepositoryList
  tags = TagList
}

mapping "aws_codeartifact_repository_permissions_policy" {
  repository = RepositoryName
  domain = DomainName
  policy_document = PolicyDocument
  domain_owner = AccountId
  policy_revision = PolicyRevision
}
