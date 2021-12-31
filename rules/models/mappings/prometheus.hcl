import = "aws-sdk-go/models/apis/amp/2020-08-01/api-2.json"

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
