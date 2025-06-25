import = "api-models-aws/models/s3/service/2006-03-01/s3-2006-03-01.json"

mapping "aws_s3_account_public_access_block" {
  account_id              = AccountId
  block_public_acls       = Setting
  block_public_policy     = Setting
  ignore_public_acls      = Setting
  restrict_public_buckets = Setting
}

mapping "aws_s3_bucket" {
  bucket                               = BucketName
  bucket_prefix                        = any
  acl                                  = any // TODO: BucketCannedACL
  policy                               = Policy
  tags                                 = TagSet
  force_destroy                        = any
  website                              = WebsiteConfiguration
  cors_rule                            = CORSRules
  versioning                           = VersioningConfiguration
  logging                              = LoggingEnabled
  lifecycle_rule                       = LifecycleRules
  acceleration_status                  = BucketAccelerateStatus
  region                               = any // BucketLocationConstraint
  request_payer                        = Payer
  replication_configuration            = ReplicationConfiguration
  server_side_encryption_configuration = ServerSideEncryptionConfiguration
  object_lock_configuration            = ObjectLockConfiguration
}

mapping "aws_s3_bucket_analytics_configuration" {
  bucket = BucketName
  filter = AnalyticsFilter
  storage_class_analysis = StorageClassAnalysis
}

mapping "aws_s3_bucket_intelligent_tiering_configuration" {
  bucket = BucketName
  status = IntelligentTieringStatus
  filter = IntelligentTieringFilter
  tiering = TieringList
}

mapping "aws_s3_bucket_inventory" {
  bucket                   = BucketName
  name                     = InventoryId
  included_object_versions = InventoryIncludedObjectVersions
  schedule                 = InventorySchedule
  destination              = InventoryDestination
  enabled                  = IsEnabled
  filter                   = InventoryFilter
  optional_fields          = InventoryOptionalFields
}

mapping "aws_s3_bucket_metric" {
  bucket = BucketName
  name   = MetricsId
  filter = MetricsFilter
}

mapping "aws_s3_bucket_notification" {
  bucket          = BucketName
  topic           = TopicConfiguration
  queue           = QueueConfiguration
  lambda_function = LambdaFunctionConfiguration
}

mapping "aws_s3_bucket_object" {
  bucket                 = BucketName
  key                    = ObjectKey
  source                 = any
  content                = any
  content_base64         = any
  acl                    = ObjectCannedACL
  cache_control          = CacheControl
  content_disposition    = ContentDisposition
  content_encoding       = ContentEncoding
  content_language       = ContentLanguage
  content_type           = ContentType
  website_redirect       = WebsiteRedirectLocation
  storage_class          = StorageClass
  etag                   = ETag
  server_side_encryption = ServerSideEncryption
  kms_key_id             = SSEKMSKeyId
  tags                   = TagSet
}

mapping "aws_s3_bucket_ownership_controls" {
  bucket = BucketName
}

mapping "aws_s3_bucket_policy" {
  bucket = BucketName
  policy = Policy
}

mapping "aws_s3_bucket_public_access_block" {
  bucket                  = BucketName
  block_public_acls       = Setting
  block_public_policy     = Setting
  ignore_public_acls      = Setting
  restrict_public_buckets = Setting
}

mapping "aws_s3_bucket_replication_configuration" {
  bucket = BucketName
  role = Role
  rule = ReplicationRules
}

mapping "aws_s3_object_copy" {
  bucket = BucketName
  key = ObjectKey
  source = CopySource
  acl = ObjectCannedACL
  cache_control = CacheControl
  content_disposition = ContentDisposition
  content_encoding = ContentEncoding
  content_language = ContentLanguage
  content_type = ContentType
  copy_if_match = CopySourceIfMatch
  copy_if_modified_since = CopySourceIfModifiedSince
  copy_if_none_match = CopySourceIfNoneMatch
  copy_if_unmodified_since = CopySourceIfUnmodifiedSince
  customer_algorithm = SSECustomerAlgorithm
  customer_key = SSECustomerKey
  customer_key_md5 = SSECustomerKeyMD5
  expected_bucket_owner = AccountId
  expected_source_bucket_owner = AccountId
  expires = Expires
  kms_encryption_context = SSEKMSEncryptionContext
  kms_key_id = SSEKMSKeyId
  metadata = Metadata
  metadata_directive = MetadataDirective
  object_lock_legal_hold_status = ObjectLockLegalHoldStatus
  object_lock_mode = ObjectLockMode
  object_lock_retain_until_date = ObjectLockRetainUntilDate
  request_payer = RequestPayer
  server_side_encryption = ServerSideEncryption
  source_customer_algorithm = CopySourceSSECustomerAlgorithm
  source_customer_key = CopySourceSSECustomerKey
  source_customer_key_md5 = CopySourceSSECustomerKeyMD5
  storage_class = StorageClass
  tagging_directive = TaggingDirective
  website_redirect = WebsiteRedirectLocation
}
