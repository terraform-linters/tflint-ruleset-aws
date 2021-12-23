import = "aws-sdk-go/models/apis/acm-pca/2017-08-22/api-2.json"

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
