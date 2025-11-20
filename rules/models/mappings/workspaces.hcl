import = "api-models-aws/models/workspaces/service/2015-04-08/workspaces-2015-04-08.json"

mapping "aws_workspaces_directory" {
  directory_id = DirectoryId
  subnet_ids = SubnetIds
  ip_group_ids = IpGroupIdList
  tags = TagList
}

mapping "aws_workspaces_ip_group" {
  name = IpGroupName
  description = IpGroupDesc
  rules = IpRuleList
  tags = TagList
}

mapping "aws_workspaces_workspace" {
  directory_id = DirectoryId
  bundle_id = BundleId
  user_name = UserName
  volume_encryption_key = VolumeEncryptionKey
  tags = TagList
  workspace_properties = WorkspaceProperties
}
