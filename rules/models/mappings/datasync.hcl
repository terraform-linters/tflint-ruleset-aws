import = "aws-sdk-go/models/apis/datasync/2018-11-09/api-2.json"

mapping "aws_datasync_agent" {
  name           = TagValue
  activation_key = ActivationKey
  tags           = InputTagList
}

mapping "aws_datasync_location_efs" {
  ec2_config          = Ec2Config
  efs_file_system_arn = EfsFilesystemArn
  subdirectory        = EfsSubdirectory
  tags                = InputTagList
}

mapping "aws_datasync_location_fsx_windows_file_system" {
  fsx_filesystem_arn = FsxFilesystemArn
  password = SmbPassword
  user = SmbUser
  domain = SmbDomain
  security_group_arns = Ec2SecurityGroupArnList
  subdirectory = FsxWindowsSubdirectory
  tags = InputTagList
}

mapping "aws_datasync_location_nfs" {
  on_prem_config  = OnPremConfig
  server_hostname = ServerHostname
  subdirectory    = EfsSubdirectory
  tags            = InputTagList
}

mapping "aws_datasync_location_s3" {
  s3_bucket_arn = S3BucketArn
  s3_config     = S3Config
  subdirectory  = EfsSubdirectory
  tags          = InputTagList
}

mapping "aws_datasync_location_smb" {
  agent_arns = AgentArnList
  domain = SmbDomain
  mount_options = SmbMountOptions
  password = SmbPassword
  server_hostname = ServerHostname
  subdirectory = SmbSubdirectory
  tags = InputTagList
  user = SmbUser
}

mapping "aws_datasync_task" {
  destination_location_arn = LocationArn
  source_location_arn      = LocationArn
  cloudwatch_log_group_arn = LogGroupArn
  name                     = TagValue
  options                  = Options
  tags                     = InputTagList
}

test "aws_datasync_agent" "name" {
  valid   = ["example"]
  invalid = ["example^example"]
}

test "aws_datasync_agent" "activation_key" {
  valid   = ["F0EFT-7FPPR-GG7MC-3I9R3-27DOH"]
  invalid = ["F0EFT7FPPRGG7MC3I9R327DOH"]
}

test "aws_datasync_location_efs" "efs_file_system_arn" {
  valid   = ["arn:aws:elasticfilesystem:us-east-1:123456789012:file-system/fs-12345678"]
  invalid = ["arn:aws:eks:us-east-1:123456789012:cluster/my-cluster"]
}

test "aws_datasync_location_efs" "subdirectory" {
  valid   = ["foo"]
  invalid = ["bar\t"]
}

test "aws_datasync_location_nfs" "server_hostname" {
  valid   = ["nfs.example.com"]
  invalid = ["nfs^example^com"]
}

test "aws_datasync_location_nfs" "subdirectory" {
  valid   = ["/exported/path"]
  invalid = ["/exported^path"]
}

test "aws_datasync_location_s3" "s3_bucket_arn" {
  valid   = ["arn:aws:s3:::my_corporate_bucket"]
  invalid = ["arn:aws:eks:us-east-1:123456789012:cluster/my-cluster"]
}

test "aws_datasync_task" "cloudwatch_log_group_arn" {
  valid   = ["arn:aws:logs:us-east-1:123456789012:log-group:my-log-group"]
  invalid = ["arn:aws:s3:::my_corporate_bucket"]
}

test "aws_datasync_task" "source_location_arn" {
  valid   = ["arn:aws:datasync:us-east-2:111222333444:location/loc-07db7abfc326c50fb"]
  invalid = ["arn:aws:datasync:us-east-2:111222333444:task/task-08de6e6697796f026"]
}
