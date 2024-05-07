# aws_resource_invalid_tags

Require tags to be assigned to a specific set of values.

## Example

```hcl
rule "aws_resource_invalid_tags" {
    enabled = true
    tags    = {
        Department  = ["finance", "hr", "payments", "engineering"]
        Environment = ["sandbox", "staging", "production"]
    }
    exclude = ["aws_autoscaling_group"]
}

provider "aws" {
    ...
    default_tags {
        tags = { Environment = "sandbox" }
    }
}

resource "aws_s3_bucket" "bucket" {
    ...
    tags = { Project: "homepage", Department = "science" }
}
```

```
$ tflint
1 issue(s) found:

Notice: aws_s3_bucket.bucket Received 'science' for tag 'Department', expected one of 'finance,hr,payments,engineering'.

  on test.tf line 3:
   3:   tags = { Project: "homepage", Department = "science" }
```

## Why

Enforce standard tag values across all resources.

## How To Fix

Align the provider, resource or autoscaling group tags to the configured expectation.
