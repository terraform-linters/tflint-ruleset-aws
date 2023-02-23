import = "aws-sdk-go/models/apis/ecr/2015-09-21/api-2.json"

mapping "aws_ecr_lifecycle_policy" {
  repository = RepositoryName
  policy     = LifecyclePolicyText
}

mapping "aws_ecr_pull_through_cache_rule" {
  ecr_repository_prefix = PullThroughCacheRuleRepositoryPrefix
  upstream_registry_url = Url
}

mapping "aws_ecr_registry_policy" {
  policy = RegistryPolicyText
}

mapping "aws_ecr_registry_scanning_configuration" {
  scan_type = ScanType
  rule = RegistryScanningRuleList
}

mapping "aws_ecr_replication_configuration" {
  replication_configuration = ReplicationConfiguration
}

mapping "aws_ecr_repository" {
  name = RepositoryName
  tags = TagList
}

mapping "aws_ecr_repository_policy" {
  repository = RepositoryName
  policy     = RepositoryPolicyText
}

test "aws_ecr_lifecycle_policy" "repository" {
  valid   = ["example"]
  invalid = ["example@com"]
}
