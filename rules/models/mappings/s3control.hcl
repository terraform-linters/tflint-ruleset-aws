import = "api-models-aws/models/s3-control/service/2018-08-20/s3-control-2018-08-20.json"

mapping "aws_s3control_access_point_policy" {
  access_point_arn = S3AccessPointArn
}

mapping "aws_s3control_bucket" {
  bucket  = BucketName
  outpost_id = NonEmptyMaxLength64String
}

mapping "aws_s3control_bucket_lifecycle_configuration" {
  bucket = BucketName
}

mapping "aws_s3control_bucket_policy" {
  bucket = BucketName
}

mapping "aws_s3control_multi_region_access_point" {
  account_id = AccountId
}

mapping "aws_s3control_multi_region_access_point_policy" {
  account_id = AccountId
}

mapping "aws_s3control_object_lambda_access_point" {
  account_id = AccountId
  name = ObjectLambdaAccessPointName
}

mapping "aws_s3control_object_lambda_access_point_policy" {
  account_id = AccountId
  name = ObjectLambdaAccessPointName
}
