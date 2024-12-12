import = "aws-sdk-ruby/apis/memorydb/2021-01-01/api-2.json"

mapping "aws_memorydb_acl" {
  user_names = UserNameListInput
  tags = TagList
}

mapping "aws_memorydb_cluster" {
  acl_name = ACLName
  security_group_ids = SecurityGroupIdsList
  snapshot_arns = SnapshotArnsList
  tags = TagList
}

mapping "aws_memorydb_parameter_group" {
  tags = TagList
}

mapping "aws_memorydb_subnet_group" {
  subnet_ids = SubnetIdentifierList
  tags = TagList
}

mapping "aws_memorydb_user" {
  access_string = AccessString
  authentication_mode = AuthenticationMode
  user_name = UserName
  tags = TagList
}
