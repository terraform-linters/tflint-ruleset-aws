# aws_iam_role_deprecated_policy_attributes

Disallow `inline_policy` configuration blocks and `managed_policy_arns` argument of the `aws_iam_role` resource.

## Example

```hcl
resource "aws_iam_role" "example" {
  name               = "example"
  assume_role_policy = data.aws_iam_policy_document.instance_assume_role_policy.json # (not shown)

  inline_policy {
    name = "my_inline_policy"

    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action   = ["ec2:Describe*"]
          Effect   = "Allow"
          Resource = "*"
        },
      ]
    })
  }

  inline_policy {
    name   = "policy-8675309"
    policy = data.aws_iam_policy_document.inline_policy.json # (not shown)
  }

  managed_policy_arns = []
}
```

```
$ tflint

3 issue(s) found:

Warning: The inline_policy argument is deprecated. Use the aws_iam_role_policy resource instead. If Terraform should exclusively manage all inline policy associations (the current behavior of this argument), use the aws_iam_role_policies_exclusive resource as well. (aws_iam_role_deprecated_policy_attributes)

  on test.tf line 5:
   5:   inline_policy {

Warning: The inline_policy argument is deprecated. Use the aws_iam_role_policy resource instead. If Terraform should exclusively manage all inline policy associations (the current behavior of this argument), use the aws_iam_role_policies_exclusive resource as well. (aws_iam_role_deprecated_policy_attributes)

  on test.tf line 20:
  20:   inline_policy {

Warning: The managed_policy_arns argument is deprecated. Use the aws_iam_role_policy_attachment resource instead. If Terraform should exclusively manage all managed policy attachments (the current behavior of this argument), use the aws_iam_role_policy_attachments_exclusive resource as well. (aws_iam_role_deprecated_policy_attributes)

  on test.tf line 25:
  25:   managed_policy_arns = []
```

## Why

The `inline_policy` and `managed_policy_arns` arguments are both deprecated since late 2024.

## How To Fix

For managing a role's inline policy, use the [aws_iam_role_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) resource instead.
For attaching managed policies to a role, use the [aws_iam_role_policy_attachment](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) resource.
