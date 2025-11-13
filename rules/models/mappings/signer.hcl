import = "api-models-aws/models/signer/service/2017-08-25/signer-2017-08-25.json"

mapping "aws_signer_signing_job" {
  profile_name = ProfileName
  source = Source
  destination = Destination
}

mapping "aws_signer_signing_profile" {
  platform_id = PlatformId
  name = ProfileName
  signature_validity_period = SignatureValidityPeriod
  tags = TagMap
}

mapping "aws_signer_signing_profile_permission" {
  profile_name = ProfileName
  profile_version = ProfileVersion
}
