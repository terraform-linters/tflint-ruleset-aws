# aws_security_group_rule_deprecated

The `aws_security_group_rule` resource should be replaced with `aws_vpc_security_group_egress_rule` or `aws_vpc_security_group_ingress_rule`. It lacks support of unique IDs, tags, and descriptions, and has difficulties managing multiple CIDR blocks.

## Example

```hcl
resource "aws_security_group_rule" "foo" {
  security_group_id = "sg-12345678"
  type              = "ingress"
}
```

```sh
❯ tflint
1 issue(s) found:

Warning: Consider using aws_vpc_security_group_egress_rule or aws_vpc_security_group_ingress_rule instead. (aws_security_group_rule_deprecated)

  on bastion.tf line 4:
   4:   resource "aws_security_group_rule" "foo" {
```

## Why

Avoid using the [`aws_security_group_rule`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group_rule) resource because it has difficulties managing multiple CIDR blocks and historically lacks unique IDs, tags, and descriptions. To prevent these issues, follow the current best practice of using the `aws_vpc_security_group_egress_rule` and `aws_vpc_security_group_ingress_rule` resources.

## How To Fix

Depending on `type`, you can fix the issue by using either `aws_vpc_security_group_egress_rule` or `aws_vpc_security_group_ingress_rule`.
