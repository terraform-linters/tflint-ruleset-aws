# aws_iam_policy_use_policy_reference

Check the source of the value for `policy` in an `aws_iam_policy` to catch values not generated from resources or data sources such as `data.aws_iam_policy_document`. This includes function calls, except those that read from files (`file` and `templatefile`).

## Configuration

```hcl
rule "aws_iam_policy_use_policy_reference" {
  enabled = true
}
```

## Example

```hcl
resource "aws_iam_policy" "raw_json" {
  name = "test_policy"

  policy = <<-EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": ["ec2:Describe*"],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::<bucketname>/*"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "jsonencode" {
  name = "test_policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "ec2:Describe*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}
```

```
$ tflint
2 issue(s) found:

Notice: The policy does not reference an object (aws_iam_policy_use_policy_reference)

  on main.tf line 4:
   4:   policy = <<-EOF
   5: {
   6:   "Version": "2012-10-17",
   7:   "Statement": [
   8:     {
   9:       "Action": ["ec2:Describe*"],
  10:       "Effect": "Allow",
  11:       "Resource": "arn:aws:s3:::<bucketname>/*"
  12:     }
  13:   ]
  14: }
  15: EOF

Notice: The policy does not reference an object (aws_iam_policy_use_policy_reference)

  on main.tf line 21:
  21:   policy = jsonencode({
  22:     Version = "2012-10-17"
  23:     Statement = [
  24:       {
  25:         Action = [
  26:           "ec2:Describe*",
  27:         ]
  28:         Effect   = "Allow"
  29:         Resource = "*"
  30:       },
  31:     ]
  32:   })
```

## Why

The `policy` argument takes a string, which must be well-formed. Specifically, it must be a JSON-formatted IAM policy document. By passing HCL to `jsonencode` we can ensure that it is JSON-formatted, but not that it is the correct structure for an IAM policy document. This check aims to make it simpler to assert correctness before applying.

## How To Fix

Define an `aws_iam_policy_document` data source and assign its `json` attribute to `policy`.
