# aws_iam_policy_gov_friendly_arns

Ensure `iam_policy` resources do not contain `arn:aws:` ARN's.

## Configuration

```hcl
rule "aws_iam_policy_gov_friendly_arns" {
  enabled = true
}
```

## Examples

```hcl
resource "aws_iam_policy" "policy" {
  name        = "test_policy"
  path        = "/"
  description = "My test policy"
  policy = <<EOF
  {
    "Version": "2012-10-17",
    "Statement": [
	  {
	    "Action": [
		  "ec2:Describe*"
	    ],
	    "Effect": "Allow",
	    "Resource": "arn:aws:s3:::<bucketname>/*""
	  }
    ]
  }
EOF
}
```

```sh
❯ tflint policy.tf
1 issue(s) found:

Warning: ARN detected in IAM policy that could potentially fail in AWS GovCloud due to resource pattern: arn:aws:.* (aws_iam_policy_gov_friendly_arns)

  on policy.tf line 5:
   5:   policy      = <<EOF
   6:   {
   7:     "Version": "2012-10-17",
   8:     "Statement": [
   9:     {
  10:       "Action": [
  11:             "ec2:Describe*"
  12:       ],
  13:       "Effect": "Allow",
  14:       "Resource": "arn:aws:s3:::<bucketname>/*""
  15:     }
  16:     ]
  17:   }
  18: EOF
```

## Why

1. Some firms have strict requirements for what AWS resources have access to government accounts. Usually only resources within the `arn:aws-us-gov:` scope are allowed.
2. When developing reusable terraform modules for many AWS accounts, arn separators are usually converted into variables when creating resources in gov and non-gov accounts like so :

```hcl
locals {
  arn_sep = var.is_govcloud ? "aws-us-gov" : "aws"
}

resource "aws_iam_policy" "policy" {
  name        = "test_policy"
  path        = "/"
  description = "My test policy"
  policy      = <<EOF
  {
    "Version": "2012-10-17",
    "Statement": [
	  {
	    "Action": [
		  "ec2:Describe*"
	    ],
	    "Effect": "Allow",
	    "Resource": "arn:${local.arn_sep}:s3:::my_bucket/*""
	  }
    ]
  }
EOF
}
```

## How To Fix

Ensure there are no `arn:aws:` scoped ARN's in your policy documents.
