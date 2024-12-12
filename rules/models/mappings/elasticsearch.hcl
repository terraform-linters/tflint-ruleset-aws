import = "aws-sdk-ruby/apis/es/2015-01-01/api-2.json"

mapping "aws_elasticsearch_domain" {
  domain_name             = DomainName
  access_policies         = PolicyDocument
  advanced_options        = AdvancedOptions
  ebs_options             = EBSOptions
  encrypt_at_rest         = EncryptionAtRestOptions
  node_to_node_encryption = NodeToNodeEncryptionOptions
  cluster_config          = ElasticsearchClusterConfig
  snapshot_options        = SnapshotOptions
  vpc_options             = VPCOptions
  log_publishing_options  = LogPublishingOptions
  cognito_options         = CognitoOptions
  elasticsearch_version   = ElasticsearchVersionString
  tags                    = TagList
}

mapping "aws_elasticsearch_domain_policy" {
  domain_name     = DomainName
  access_policies = PolicyDocument
}

mapping "aws_elasticsearch_domain_saml_options" {
  domain_name = DomainName
}
