# aws_iam_policy_attachment_exclusive_attachment

This rule checks whether the `aws_iam_policy_attachment` resource is used.

The `aws_iam_policy_attachment` resource creates exclusive attachments for IAM policies. Within the entire AWS account, all users, roles, and groups that a single policy is attached to must be specified by a single aws_iam_policy_attachment resource. 

## Configuration

```hcl
rule "aws_iam_policy_attachment_exclusive_attachment" {
  enabled = true
}
```

## Example

```hcl
resource "aws_iam_policy_attachment" "attachment" {
	name = "test_attachment"
}
```

```shell
$ tflint
1 issue(s) found:
Warning: Within the entire AWS account, all users, roles, and groups that a single policy is attached to must be specified by a single aws_iam_policy_attachment resource. Consider aws_iam_role_policy_attachment, aws_iam_user_policy_attachment, or aws_iam_group_policy_attachment instead. (aws_iam_policy_attachment_has_alternatives)
  on template.tf line 2:
   2:  resource "aws_iam_policy_attachment" "attachment" {
```

## Why

The [`aws_iam_policy_attachment`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy_attachment) resource creates exclusive attachments of IAM policies. Across the entire AWS account, all the users/roles/groups to which a single policy is attached must be declared by a single `aws_iam_policy_attachment` resource. This means that even any users/roles/groups that have the attached policy via any other mechanism (including other Terraform resources) will have that attached policy revoked by this resource.

## How To Fix

Consider using `aws_iam_role_policy_attachment`, `aws_iam_user_policy_attachment`, or `aws_iam_group_policy_attachment` instead. These resources do not enforce exclusive attachment of an IAM policy.
