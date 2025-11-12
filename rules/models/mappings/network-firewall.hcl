import = "api-models-aws/models/network-firewall/service/2020-11-12/network-firewall-2020-11-12.json"

mapping "aws_networkfirewall_firewall" {
  description = Description
  firewall_policy_arn = ResourceArn
  name = ResourceName
  subnet_mapping = SubnetMappings
  tags = TagList
  vpc_id = VpcId
}

mapping "aws_networkfirewall_firewall_policy" {
  description = Description
  firewall_policy = FirewallPolicy
  name = ResourceName
  tags = TagList
}

mapping "aws_networkfirewall_logging_configuration" {
  firewall_arn = ResourceArn
  logging_configuration = LoggingConfiguration
}

mapping "aws_networkfirewall_resource_policy" {
  policy = PolicyString
  resource_arn = ResourceArn
}

mapping "aws_networkfirewall_rule_group" {
  capacity = RuleCapacity
  description = Description
  name = ResourceName
  rule_group = RuleGroup
  rules = RulesString
  tags = TagList
  type = RuleGroupType
}
