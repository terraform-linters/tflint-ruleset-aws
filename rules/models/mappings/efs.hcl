import = "aws-sdk-go/models/apis/elasticfilesystem/2015-02-01/api-2.json"

mapping "aws_efs_access_point" {
  file_system_id = FileSystemId
  posix_user = PosixUser
  root_directory = RootDirectory
  tags = Tags
}

mapping "aws_efs_backup_policy" {
  file_system_id = FileSystemId
  backup_policy = BackupPolicy
}

mapping "aws_efs_file_system" {
  creation_token                  = CreationToken
  encrypted                       = Encrypted
  kms_key_id                      = KmsKeyId
  performance_mode                = PerformanceMode
  provisioned_throughput_in_mibps = ProvisionedThroughputInMibps
  tags                            = Tags
  throughput_mode                 = ThroughputMode
}

mapping "aws_efs_file_system_policy" {
  file_system_id = FileSystemId
  bypass_policy_lockout_safety_check = BypassPolicyLockoutSafetyCheck
  policy = Policy
}

mapping "aws_efs_mount_target" {
  file_system_id  = FileSystemId
  subnet_id       = SubnetId
  ip_address      = IpAddress
  security_groups = SecurityGroups
}

test "aws_efs_file_system" "performance_mode" {
  valid   = ["generalPurpose"]
  invalid = ["minIO"]
}

test "aws_efs_file_system" "throughput_mode" {
  valid   = ["bursting"]
  invalid = ["generalPurpose"]
}
