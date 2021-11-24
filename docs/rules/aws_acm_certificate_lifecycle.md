# aws_acm_certificate_lifecycle

This rule ensures lifecycle [`create_before_destroy`](https://www.terraform.io/docs/language/meta-arguments/lifecycle.html#create_before_destroy)
argument is set to `true` for acm certificates as per
[`aws_acm_certificate`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/acm_certificate)
official docs:

> It's recommended to specify `create_before_destroy = true` in a
> lifecycle block to replace a certificate which is currently in use


## Example

```hcl
resource "aws_acm_certificate" "test" {
  domain_name       = var.domain_name
  validation_method = "DNS"
}
```

```console
$ tflint
1 issue(s) found:

Warning: resource `aws_acm_certificate` needs to contain `create_before_destroy = true` in `lifecycle` block (aws_acm_certificate_lifecycle)

  on test.tf line 1:
   1: resource "aws_acm_certificate" "test" {
```

## Why

If you don't add this, terraform will timeout when trying to replace a certificate
that's in use at the moment of apply.

## How To Fix

```hcl
resource "aws_acm_certificate" "test" {
  domain_name       = local.assets_domain_name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}
```
