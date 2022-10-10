# aws_iam_policy_attachment_has_alternatives

Consider alternative resources to `aws_iam_policy_attachment`.

## Configuration

```hcl
rule "aws_iam_policy_attachment_has_alternatives" {
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
Warning: Consider aws_iam_role_policy_attachment, aws_iam_user_policy_attachment, or aws_iam_group_policy_attachment instead. (aws_iam_policy_attachment_has_alternatives)
  on template.tf line 2:
   2:   name "test_attachment" // Consider alternatives!

```

## Why

The `aws_iam_policy_attachment` resource creates exclusive attachments of IAM policies. Across the entire AWS account, all of the users/roles/groups to which a single policy is attached must be declared by a single `aws_iam_policy_attachment` resource. This means that even any users/roles/groups that have the attached policy via any other mechanism (including other Terraform resources) will have that attached policy revoked by this resource. https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy_attachment

## How To Fix

Consider `aws_iam_role_policy_attachment`, `aws_iam_user_policy_attachment`, or `aws_iam_group_policy_attachment` instead. These resources do not enforce exclusive attachment of an IAM policy.
