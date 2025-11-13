import = "api-models-aws/models/resource-groups/service/2017-11-27/resource-groups-2017-11-27.json"

mapping "aws_resourcegroups_group" {
  name           = GroupName
  description    = any //GroupDescription
  resource_query = ResourceQuery
}
