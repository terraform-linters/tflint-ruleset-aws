import = "aws-sdk-go/models/apis/fsx/2018-03-01/api-2.json"

mapping "aws_fsx_backup" {
  file_system_id = FileSystemId
  tags = Tags
  volume_id = VolumeId
}

mapping "aws_fsx_lustre_file_system" {
  storage_capacity              = StorageCapacity
  subnet_ids                    = SubnetIds
  export_path                   = any // ArchivePath
  import_path                   = any // ArchivePath
  imported_file_chunk_size      = Megabytes
  security_group_ids            = SecurityGroupIds
  tags                          = Tags
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
  tags = Tags
}

mapping "aws_fsx_ontap_volume" {
  name = VolumeName
  junction_path = JunctionPath
  security_style = SecurityStyle
  size_in_megabytes = VolumeCapacity
  storage_virtual_machine_id = StorageVirtualMachineId
  tags = Tags
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
  tags                              = Tags
  weekly_maintenance_start_time     = WeeklyTime
}
