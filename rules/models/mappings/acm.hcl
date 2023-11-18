import = "aws-sdk-go/models/apis/acm/2015-12-08/api-2.json"

mapping "aws_acm_certificate" {
  // domain_name            = DomainNameString
  subject_alternative_names = DomainList
  // validation_method      = ValidationMethod
  private_key               = PrivateKeyBlob
  certificate_body          = CertificateBody
  certificate_chain         = CertificateChain
  certificate_authority_arn = PcaArn
  tags                      = TagList
}

mapping "aws_acm_certificate_validation" {
  certificate_arn = Arn
}

test "aws_acm_certificate" "certificate_authority_arn" {
  ok = "arn:aws:acm-pca:us-east-1:0000000000:certificate-authority/xxxxxx-xxx-xxx-xxxx-xxxxxxxxx"
  ng = "arn:aws:unknown-service:us-east-1:0000000000:certificate-authority/xxxxxx-xxx-xxx-xxxx-xxxxxxxxx"
}
