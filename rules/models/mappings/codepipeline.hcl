import = "aws-sdk-go/models/apis/codepipeline/2015-07-09/api-2.json"

mapping "aws_codepipeline" {
  name     = PipelineName
  role_arn = RoleArn
}

mapping "aws_codepipeline_webhook" {
  name            = WebhookName
  authentication  = WebhookAuthenticationType
  target_action   = ActionName
  target_pipeline = PipelineName
}

test "aws_codepipeline" "name" {
  valid   = ["tf-test-pipeline"]
  invalid = ["test/pipeline"]
}

test "aws_codepipeline" "role_arn" {
  valid   = ["arn:aws:iam::123456789012:role/s3access"]
  invalid = ["arn:aws:iam::123456789012:instance-profile/s3access-profile"]
}

test "aws_codepipeline_webhook" "name" {
  valid   = ["test-webhook-github-bar"]
  invalid = ["webhook-github-bar/testing"]
}

test "aws_codepipeline_webhook" "authentication" {
  valid   = ["GITHUB_HMAC"]
  invalid = ["GITLAB_HMAC"]
}

test "aws_codepipeline_webhook" "target_action" {
  valid   = ["Source"]
  invalid = ["Source/Example"]
}
