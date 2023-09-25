# aws_s3_bucket_name

Ensures all S3 bucket names match the specified naming rules.

## Configuration

```hcl
rule "aws_s3_bucket_name" {
  enabled = true
  regex = "^[a-z\\-]+$"
  prefix = "my-org"
}
```

* `regex`: A Go regex that bucket names must match (string)
* `prefix`: A prefix that should be used for bucket names (string)

## Examples

```hcl
resource "aws_s3_bucket" "foo" {
  bucket = "foo"
}

resource "aws_s3_bucket" "too_long" {
  bucket = "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test"
}
```

```sh
$ tflint
2 issue(s) found:

Error: Bucket name "foo" does not have prefix "my-org" (aws_s3_bucket_name)

  on main.tf line 2:
  2:   bucket = "foo"

Error: Bucket name "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test" length must be within 3 - 63 character range (aws_s3_bucket_name)

  on main.tf line 2:
  2:   bucket = "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test"
```

## Why

Amazon S3 bucket names must be globally unique and have [restrictive naming rules](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html).

* Prefixing bucket names with an organization name can help avoid naming conflicts
* You may wish to enforce other naming conventions (e.g., disallowing dots)

## How To Fix

Ensure the bucket name matches the specified rules.
