import = "aws-sdk-go/models/apis/glue/2017-03-31/api-2.json"

mapping "aws_glue_catalog_database" {
  name         = any // NameString
  catalog_id   = any // CatalogIdString
  description  = any // DescriptionString
  location_uri = any // URI
  parameters   = ParametersMap
}

mapping "aws_glue_catalog_table" {
  name               = any // NameString
  database_name      = any // NameString
  catalog_id         = any // CatalogIdString
  description        = any // DescriptionString
  owner              = any // NameString
  retention          = NonNegativeInteger
  storage_descriptor = StorageDescriptor
  partition_keys     = ColumnList
  view_original_text = ViewTextString
  view_expanded_text = ViewTextString
  table_type         = TableTypeString
  parameters         = ParametersMap
}

mapping "aws_glue_classifier" {
  grok_classifier = CreateGrokClassifierRequest
  json_classifier = CreateJsonClassifierRequest
  name            = any // NameString
  xml_classifier  = CreateXMLClassifierRequest
}

mapping "aws_glue_connection" {
  catalog_id                       = any // CatalogIdString
  connection_properties            = ConnectionProperties
  connection_type                  = ConnectionType
  description                      = any // DescriptionString
  match_criteria                   = MatchCriteria
  name                             = any // NameString
  physical_connection_requirements = PhysicalConnectionRequirements
}

mapping "aws_glue_crawler" {
  database_name          = DatabaseName
  name                   = any // NameString
  role                   = Role
  classifiers            = ClassifierNameList
  configuration          = CrawlerConfiguration
  description            = any // DescriptionString
  dynamodb_target        = DynamoDBTargetList
  jdbc_target            = JdbcTargetList
  s3_target              = S3TargetList
  schedule               = CronExpression
  schema_change_policy   = SchemaChangePolicy
  table_prefix           = TablePrefix
  security_configuration = CrawlerSecurityConfiguration
}

mapping "aws_glue_data_catalog_encryption_settings" {
  catalog_id = CatalogIdString
}

mapping "aws_glue_dev_endpoint" {
  public_keys = PublicKeysList
  role_arn = RoleArn
  tags = TagsMap
  worker_type = WorkerType
}

mapping "aws_glue_job" {
  command                = JobCommand
  connections            = ConnectionsList
  default_arguments      = GenericMap
  description            = any // DescriptionString
  execution_property     = ExecutionProperty
  max_capacity           = NullableDouble
  max_retries            = MaxRetries
  name                   = any // NameString
  role_arn               = RoleString
  timeout                = Timeout
  security_configuration = any // NameString
}

mapping "aws_glue_ml_transform" {
  name = NameString
  input_record_tables = GlueTables
  parameters = TransformParameters
  role_arn = RoleString
  description = DescriptionString
  glue_version = GlueVersionString
  tags = TagsMap
  timeout = Timeout
  worker_type = WorkerType
}

mapping "aws_glue_partition" {
  database_name = NameString
  partition_values = ValueStringList
  catalog_id = CatalogIdString
  storage_descriptor = StorageDescriptor
  parameters = ParametersMap
}

mapping "aws_glue_partition_index" {
  table_name = NameString
  database_name = NameString
  partition_index = PartitionIndex
  catalog_id = CatalogIdString
}

mapping "aws_glue_registry" {
  registry_name = SchemaRegistryNameString
  description = DescriptionString
  tags = TagsMap
}

mapping "aws_glue_resource_policy" {
  enable_hybrid = EnableHybridValues
}

mapping "aws_glue_schema" {
  schema_name = SchemaRegistryNameString
  data_format = DataFormat
  compatibility = Compatibility
  schema_definition = SchemaDefinitionString
  description = DescriptionString
  tags = TagsMap
}

mapping "aws_glue_security_configuration" {
  encryption_configuration = EncryptionConfiguration
  name                     = any // NameString
}

mapping "aws_glue_trigger" {
  actions     = ActionList
  description = any // DescriptionString
  enabled     = Boolean
  name        = any // NameString
  predicate   = Predicate
  schedule    = GenericString
  type        = TriggerType
}

mapping "aws_glue_user_defined_function" {
  name = NameString
  catalog_id = CatalogIdString
  database_name = NameString
  class_name = NameString
  owner_name = NameString
  owner_type = PrincipalType
  resource_uris = ResourceUriList
}

mapping "aws_glue_workflow" {
  name = NameString
  default_run_properties = WorkflowRunProperties
  tags = TagsMap
}
