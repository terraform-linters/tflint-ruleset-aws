import = "api-models-aws/models/fsx/service/2018-03-01/fsx-2018-03-01.json"

mapping "aws_fsx_backup" {
  file_system_id = FileSystemId
  tags = listmap(Tags, TagKey, TagValue)
  volume_id = VolumeId
}

mapping "aws_fsx_lustre_file_system" {
  storage_capacity              = StorageCapacity
  subnet_ids                    = SubnetIds
  export_path                   = any // ArchivePath
  import_path                   = any // ArchivePath
  imported_file_chunk_size      = Megabytes
  security_group_ids            = SecurityGroupIds
  tags                          = listmap(Tags, TagKey, TagValue)
  weekly_maintenance_start_time = WeeklyTime
}

mapping "aws_fsx_ontap_file_system" {
  preferred_subnet_id = SubnetId
  weekly_maintenance_start_time = WeeklyTime
  deployment_type = OntapDeploymentType
  automatic_backup_retention_days = AutomaticBackupRetentionDays
  daily_automatic_backup_start_time = DailyTime
  disk_iops_configuration = DiskIopsConfiguration
  endpoint_ip_address_range = IpAddressRange
  fsx_admin_password = AdminPassword
  route_table_ids = RouteTableIds
  throughput_capacity = MegabytesPerSecond
}

mapping "aws_fsx_ontap_storage_virtual_machine" {
  active_directory_configuration = CreateSvmActiveDirectoryConfiguration
  file_system_id = FileSystemId
  name = StorageVirtualMachineName
  root_volume_security_style = StorageVirtualMachineRootVolumeSecurityStyle
  tags = listmap(Tags, TagKey, TagValue)
}

mapping "aws_fsx_ontap_volume" {
  name = VolumeName
  junction_path = JunctionPath
  security_style = SecurityStyle
  size_in_megabytes = VolumeCapacity
  storage_virtual_machine_id = StorageVirtualMachineId
  tags = listmap(Tags, TagKey, TagValue)
}

mapping "aws_fsx_openzfs_file_system" {
  deployment_type = OpenZFSDeploymentType
  storage_capacity = StorageCapacity
  subnet_ids = SubnetIds
  throughput_capacity = MegabytesPerSecond
  automatic_backup_retention_days = AutomaticBackupRetentionDays
  backup_id = BackupId
  copy_tags_to_backups = Flag
  copy_tags_to_volumes = Flag
  daily_automatic_backup_start_time = DailyTime
  disk_iops_configuration = DiskIopsConfiguration
  kms_key_id = any // KmsKeyId
  root_volume_configuration = OpenZFSCreateRootVolumeConfiguration
  security_group_ids = SecurityGroupIds
  storage_type = StorageType
  weekly_maintenance_start_time = WeeklyTime
}

mapping "aws_fsx_openzfs_snapshot" {
  name = SnapshotName
  tags = listmap(Tags, TagKey, TagValue)
  volume_id = VolumeId
}

mapping "aws_fsx_openzfs_volume" {
  parent_volume_id = VolumeId
  origin_snapshot = CreateOpenZFSOriginSnapshotConfiguration
  copy_tags_to_snapshots = Flag
  data_compression_type = OpenZFSDataCompressionType
  nfs_exports = OpenZFSNfsExports
  read_only = ReadOnly
  storage_capacity_quota_gib = IntegerNoMax
  storage_capacity_reservation_gib = IntegerNoMax
  user_and_group_quotas = OpenZFSUserAndGroupQuotas
}

mapping "aws_fsx_windows_file_system" {
  storage_capacity                  = StorageCapacity
  subnet_ids                        = SubnetIds
  throughput_capacity               = MegabytesPerSecond
  active_directory_id               = DirectoryId
  automatic_backup_retention_days   = AutomaticBackupRetentionDays
  copy_tags_to_backups              = Flag
  daily_automatic_backup_start_time = DailyTime
  kms_key_id                        = any // KmsKeyId
  security_group_ids                = SecurityGroupIds
  self_managed_active_directory     = SelfManagedActiveDirectoryConfiguration
  skip_final_backup                 = Flag
  tags                              = listmap(Tags, TagKey, TagValue)
  weekly_maintenance_start_time     = WeeklyTime
}
