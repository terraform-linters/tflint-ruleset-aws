# aws_security_group_egress_and_ingress_blocks_deprecated

Avoid using the `ingress` and `egress` arguments of the `aws_security_group` resource to configure in-line rules, as they have difficulties managing multiple CIDR blocks and lack unique IDs, tags, and descriptions. To prevent these issues, follow the current best practice of using the `aws_vpc_security_group_egress_rule` and `aws_vpc_security_group_ingress_rule` resources, with one CIDR block per rule.

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

// TODO: Write the output when inspects the above code

```

## Why

Refrain from using the `ingress` and `egress` arguments of the `aws_security_group` resource for in-line rules, as they have difficulties managing multiple CIDR blocks and historically lack unique IDs, tags, and descriptions. To prevent these issues, follow the best practice of using the `aws_vpc_security_group_egress_rule` and `aws_vpc_security_group_ingress_rule` resources, with one CIDR block per rule.

Avoid using the `aws_security_group` resource with in-line rules (using the ingress and egress arguments) alongside the `aws_vpc_security_group_egress_rule`, `aws_vpc_security_group_ingress_rule`, or `aws_security_group_rule` resources. This practice can lead to rule conflicts, perpetual differences, and rules being overwritten.

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

using the attributes according to your code.

`ingress` blocks are replaced by `aws_vpc_security_group_egress_rule` in the same way.
