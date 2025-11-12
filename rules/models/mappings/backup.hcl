import = "api-models-aws/models/backup/service/2018-11-15/backup-2018-11-15.json"

mapping "aws_backup_selection" {
  name = BackupSelectionName
}

mapping "aws_backup_vault" {
  name = BackupVaultName
}

mapping "aws_backup_vault_lock_configuration" {
  backup_vault_name = BackupVaultName
}

mapping "aws_backup_vault_notifications" {
  backup_vault_name = BackupVaultName
  sns_topic_arn = ARN
  backup_vault_events = BackupVaultEvents
}

mapping "aws_backup_vault_policy" {
  backup_vault_name = BackupVaultName
}

test "aws_backup_selection" "name" {
  ok = "tf_example_backup_selection"
  ng = "tf_example_backup_selection_tf_example_backup_selection"
}

test "aws_backup_vault" "name" {
  ok = "example_backup_vault"
  ng = "example_backup_vault_example_backup_vault_example_backup_vault"
}
