import = "aws-sdk-ruby/apis/lakeformation/2017-03-31/api-2.json"

mapping "aws_lakeformation_data_lake_settings" {
  admins = DataLakePrincipalList
  # catalog_id = CatalogIdString
  create_database_default_permissions = PrincipalPermissionsList
  create_table_default_permissions = PrincipalPermissionsList
  trusted_resource_owners = TrustedResourceOwners
}

mapping "aws_lakeformation_permissions" {
  permissions = PermissionList
  principal = DataLakePrincipal
  catalog_resource = CatalogResource
  data_location = DataLocationResource
  database = DatabaseResource
  table = TableResource
  table_with_columns = TableWithColumnsResource
  # catalog_id = CatalogIdString
  permissions_with_grant_option = PermissionList
}

mapping "aws_lakeformation_resource" {
  arn = ResourceArnString
  role_arn = IAMRoleArn
}
