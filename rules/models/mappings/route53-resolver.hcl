import = "aws-sdk-go/models/apis/route53resolver/2018-04-01/api-2.json"

mapping "aws_route53_resolver_dnssec_config" {
  resource_id = ResourceId
}

mapping "aws_route53_resolver_endpoint" {
  direction          = ResolverEndpointDirection
  ip_address         = IpAddressesRequest
  security_group_ids = SecurityGroupIds
  name               = any // Name
  tags               = TagList
}

mapping "aws_route53_resolver_firewall_config" {
  resource_id = ResourceId
  firewall_fail_open = FirewallFailOpenStatus
}

mapping "aws_route53_resolver_firewall_domain_list" {
  name = Name
  tags = TagList
}

mapping "aws_route53_resolver_firewall_rule" {
  name = Name
  action = Action
  block_override_dns_type = BlockOverrideDnsType
  block_override_domain = BlockOverrideDomain
  block_override_ttl = BlockOverrideTtl
  block_response = BlockResponse
  firewall_domain_list_id = ResourceId
  firewall_rule_group_id = ResourceId
  priority = Priority
}

mapping "aws_route53_resolver_firewall_rule_group" {
  name = Name
  tags = TagList
}

mapping "aws_route53_resolver_firewall_rule_group_association" {
  name = Name
  firewall_rule_group_id = ResourceId
  mutation_protection = MutationProtectionStatus
  priority = Priority
  vpc_id = ResourceId
  tags = TagList
}

mapping "aws_route53_resolver_rule" {
  domain_name          = DomainName
  rule_type            = RuleTypeOption
  name                 = any // Name
  resolver_endpoint_id = ResourceId
  target_ip            = TargetList
  tags                 = TagList
}

mapping "aws_route53_resolver_rule_association" {
  resolver_rule_id = ResourceId
  vpc_id           = ResourceId
  name             = any // Name
}

mapping "aws_route53_resolver_query_log_config" {
  destination_arn = DestinationArn
  name = ResolverQueryLogConfigName
  tags = TagList
}

mapping "aws_route53_resolver_query_log_config_association" {
  resolver_query_log_config_id = ResourceId
  resource_id = ResourceId
}
