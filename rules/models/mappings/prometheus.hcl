import = "api-models-aws/models/amp/service/2020-08-01/amp-2020-08-01.json"

mapping "aws_prometheus_alert_manager_definition" {
  workspace_id = WorkspaceId
  definition = AlertManagerDefinitionData
}

mapping "aws_prometheus_rule_group_namespace" {
  name = RuleGroupsNamespaceName
  workspace_id = WorkspaceId
  data = RuleGroupsNamespaceData
}

mapping "aws_prometheus_workspace" {
  alias = WorkspaceAlias
}
