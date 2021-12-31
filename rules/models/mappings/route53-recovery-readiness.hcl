import = "aws-sdk-go/models/apis/route53-recovery-readiness/2019-12-02/api-2.json"

mapping "aws_route53recoveryreadiness_cell" {
  tags = Tags
}

mapping "aws_route53recoveryreadiness_readiness_check" {
  tags = Tags
}

mapping "aws_route53recoveryreadiness_recovery_group" {
  tags = Tags
}

mapping "aws_route53recoveryreadiness_resource_set" {
  resource_set_type = __stringPatternAWSAZaZ09AZaZ09
  tags = Tags
}
