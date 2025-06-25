import = "api-models-aws/models/acm-pca/service/2017-08-22/acm-pca-2017-08-22.json"

mapping "aws_acmpca_certificate" {
  certificate_authority_arn   = Arn
  certificate_signing_request = CsrBlob
  signing_algorithm           = SigningAlgorithm
}

mapping "aws_acmpca_certificate_authority" {
  type = CertificateAuthorityType
}

mapping "aws_acmpca_certificate_authority_certificate" {
  certificate_authority_arn = Arn
  certificate               = CertificateBodyBlob
  certificate_chain         = CertificateChainBlob
}

test "aws_acmpca_certificate_authority" "type" {
  ok = "SUBORDINATE"
  ng = "ORDINATE"
}
