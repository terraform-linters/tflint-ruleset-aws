import = "aws-sdk-go/models/apis/route53-recovery-control-config/2020-11-02/api-2.json"

mapping "aws_route53recoverycontrolconfig_cluster" {
  name = __stringMin1Max64PatternS
}

mapping "aws_route53recoverycontrolconfig_control_panel" {
  cluster_arn = __stringMin1Max256PatternAZaZ09
  name = __stringMin1Max64PatternS
}

mapping "aws_route53recoverycontrolconfig_routing_control" {
  cluster_arn = __stringMin1Max256PatternAZaZ09
  name = __stringMin1Max64PatternS
  control_panel_arn = __stringMin1Max256PatternAZaZ09
}

mapping "aws_route53recoverycontrolconfig_safety_rule" {
  control_panel_arn = __stringMin1Max256PatternAZaZ09
  name = __stringMin1Max64PatternS
  rule_config = RuleConfig
  asserted_controls = __listOf__stringMin1Max256PatternAZaZ09
  gating_controls = __listOf__stringMin1Max256PatternAZaZ09
  target_controls = __listOf__stringMin1Max256PatternAZaZ09
}
