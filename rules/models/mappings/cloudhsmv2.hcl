import = "api-models-aws/models/cloudhsm-v2/service/2017-04-28/cloudhsm-v2-2017-04-28.json"

mapping "aws_cloudhsm_v2_cluster" {
  source_backup_identifier = BackupId
  hsm_type                 = HsmType
}

mapping "aws_cloudhsm_v2_hsm" {
  cluster_id        = ClusterId
  subnet_id         = SubnetId
  availability_zone = ExternalAz
  ip_address        = IpAddress
}

test "aws_cloudhsm_v2_cluster" "source_backup_identifier" {
  ok = "backup-rtq2dwi2gq6"
  ng = "rtq2dwi2gq6"
}

test "aws_cloudhsm_v2_cluster" "hsm_type" {
  ok = "hsm1.medium"
  ng = "t2.medium"
}

test "aws_cloudhsm_v2_hsm" "cluster_id" {
  ok = "cluster-jxhlf7644ne"
  ng = "jxhlf7644ne"
}

test "aws_cloudhsm_v2_hsm" "subnet_id" {
  ok = "subnet-0e358c43"
  ng = "0e358c43"
}

test "aws_cloudhsm_v2_hsm" "availability_zone" {
  ok = "us-east-1a"
  ng = "us-east-1"
}

test "aws_cloudhsm_v2_hsm" "ip_address" {
  ok = "8.8.8.8"
  ng = "2001:4860:4860::8888"
}
