# aws_s3_bucket_name_length

Ensures all S3 bucket names contains a valid number of characters.

## Example

```hcl
resource "aws_s3_bucket" "too_long" {
  bucket = "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test"
}
```

```sh
$ tflint
1 issue(s) found:

Error: Bucket name "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test" length must be within 3 - 63 character range (aws_s3_bucket_name_length)

  on main.tf line 2:
  2:   bucket = "a-really-ultra-hiper-super-long-foo-bar-baz-bucket-name.domain.test"
```

## Why

Amazon S3 bucket names have [restrictive naming rules](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html).

* This rule is specifically for how many characters long a bucket name must have.

## How To Fix

Ensure the bucket name matches the minimum and maximum characters length count.
