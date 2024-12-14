# aws_security_group_rule_deprecated

// TODO: Write the rule's description here

## Example

```hcl
resource "aws_security_group_rule" "foo" {
  security_group_id = "sg-12345678"
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

Avoid using the `aws_security_group_rule` resource, as it struggles with managing multiple CIDR blocks, and, due to the historical lack of unique IDs, tags and descriptions.

For further information, see the [Terraform documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group_rule).

## How To Fix

Depending on `foo.type`, you can fix the issue by using either `aws_vpc_security_group_egress_rule` or `aws_vpc_security_group_ingress_rule`.
