import = "aws-sdk-go/models/apis/cloudhsmv2/2017-04-28/api-2.json"

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
  valid   = ["backup-rtq2dwi2gq6"]
  invalid = ["rtq2dwi2gq6"]
}

test "aws_cloudhsm_v2_cluster" "hsm_type" {
  valid   = ["hsm1.medium"]
  invalid = ["hsm1.micro"]
}

test "aws_cloudhsm_v2_hsm" "cluster_id" {
  valid   = ["cluster-jxhlf7644ne"]
  invalid = ["jxhlf7644ne"]
}

test "aws_cloudhsm_v2_hsm" "subnet_id" {
  valid   = ["subnet-0e358c43"]
  invalid = ["0e358c43"]
}

test "aws_cloudhsm_v2_hsm" "availability_zone" {
  valid   = ["us-east-1a"]
  invalid = ["us-east-1"]
}

test "aws_cloudhsm_v2_hsm" "ip_address" {
  valid   = ["8.8.8.8"]
  invalid = ["2001:4860:4860::8888"]
}
