# aws_iam_policy_document_gov_friendly

Ensure `iam_policy_document` data sources do not contain `arn:aws:` ARN's.

## Configuration

```hcl
rule "aws_iam_policy_document_gov_friendly" {
  enabled = true
}
```

## Examples

```hcl
resource "aws_iam_policy_document" "example" {
  statement {
    sid = "1"
    actions = [
      "s3:ListAllMyBuckets",
      "s3:GetBucketLocation",
    ]
    resources = [
      "arn:aws:s3:::*",
    ]
  }

  statement {
    sid = "1"
    actions = [
      "s3:ListAllMyBuckets",
      "s3:GetBucketLocation",
    ]
    resources = [
      "arn:aws-us-gov:s3:::*",
    ]
  }
}
```

```sh
❯ tflint policy.tf
1 issue(s) found:

Warning: ARN detected in IAM policy document that could potentially fail in AWS GovCloud due to resource pattern: arn:aws:.* (aws_iam_policy_document_gov_friendly_arns)

  on policy.tf line 8:
   8:     resources = [
   9:       "arn:aws:s3:::*",
  10:     ]
```

## Why

1. Some firms have strict requirements for what AWS resources have access to government accounts. Usually only resources within the `arn:aws-us-gov:` scope are allowed.
2. When developing reusable terraform modules for many AWS accounts, arn separators are usually converted into variables when creating resources in gov and non-gov accounts like so :

```hcl
locals {
  arn_sep = var.is_govcloud ? "aws-us-gov" : "aws"
}

resource "aws_iam_policy_document" "example" {
  statement {
    sid = "1"
    actions = [
      "s3:ListAllMyBuckets",
      "s3:GetBucketLocation",
    ]
    resources = [
      "arn:${local.arn_sep}:s3:::my_bucket",
    ]
  }
}
```

## How To Fix

Ensure there are no `arn:aws:` scoped ARN's in your policy documents.
