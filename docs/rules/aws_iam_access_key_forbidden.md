# aws_iam_access_key_forbidden

Disallow using `aws_iam_access_key` resources.

## Example

```hcl
resource "aws_iam_access_key" "my_key" {
  user = "my-user"
}
```

```
$ tflint
1 issue(s) found:

Error: "aws_iam_access_key" resources are forbidden. (aws_iam_access_key_forbidden)

  on main.tf line 1:
   1: resource "aws_iam_access_key" "my_key" {
```

## Why

Managing IAM access keys with terraform leads to sensitive secrets being stored in plain text in terraform state files.

## How To Fix

Do not manage IAM access keys with terraform, or use the [pgp_key](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_access_key#pgp_key) argument and protect your backend state file judiciously.
