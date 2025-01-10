# aws_security_group_inline_rules

Disallow `ingress` and `egress` arguments of the `aws_security_group` resource.


## Example

```hcl
resource "aws_security_group" "foo" {
  name = "test"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
```

```
$ tflint
2 issue(s) found:

Notice: Replace this egress block with aws_vpc_security_group_egress_rule. (aws_security_group_inline_rules)

  on test.tf line 4:
   4:   egress {
   
Notice: Replace this ingress block with aws_vpc_security_group_ingress_rule. (aws_security_group_inline_rules)

  on test.tf line 11:
   11:   ingress {
   
```

## Why

In-line rules are difficult to manage and maintain, especially when multiple CIDR blocks are used. They lack unique IDs, tags, and descriptions, which makes it hard to identify and manage them.

See [best practices](https://registry.terraform.io/providers/hashicorp/aws/5.82.2/docs/resources/security_group).

## How To Fix

Replace an `egress` block by

```hcl
resource "aws_vpc_security_group_egress_rule" "example" {
  security_group_id = aws_security_group.example.id

  cidr_ipv4   = "0.0.0.0/0"
  from_port   = 443
  ip_protocol = "tcp"
  to_port     = 443
}
```

using the attributes according to your code. For `ingress` blocks use `aws_vpc_security_group_ingress_rule` in the same way.
