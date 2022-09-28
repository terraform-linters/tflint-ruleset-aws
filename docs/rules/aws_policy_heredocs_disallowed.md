# aws_policy_heredocs_disallowed

This rule ensures that heredocs are not in use across the terraform configuration

## Example

We prefer a jsonencoded configuration of a policy
```hcl
  resource "aws_iam_policy" "test" {
    policy = jsonencode({
        Statement = [
            {
                Action = [
                    "iam:PassRole"
                ],
            }
        ]
    })
  }
```
over one with a heredoc:
```hcl
resource "aws_iam_policy" "test" {
  policy = <<EOF
{
  "Statement": [
    {
      "Action":  [
        "iam:PassRole"
      ],
    }
  ]
}
EOF
}
```


```
$ tflint
1 issue(s) found:

Notice: Avoid Use of HEREDOC syntax for policy documents (aws_policy_heredocs_disallowed)

  on test.tf line 3:
 3:   policy = <<EOF
 4: {

```


## Why

Preference is made to policy documents written with jsonencode() functionality vs the Heredoc functionality.

## How To Fix

convert your heredoc code to json
```hcl
  resource "aws_iam_policy" "test" {
    policy = jsonencode({
        Statement = [
            {
                Action = [
                    "iam:PassRole"
                ],
            }
        ]
    })
  }
```
