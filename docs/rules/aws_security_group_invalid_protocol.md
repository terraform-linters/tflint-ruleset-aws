# aws_security_group_invalid_protocol

Disallow using invalid protocol.

## Example

```hcl
resource "aws_security_group" "sample" {
  description = "test security group"

  ingress {
    from_port = 443
    to_port   = 443
    protocol  = "https"
  }
}
```

```
$ tflint
1 issue(s) found:

Error: "https" is an invalid protocol. It allows "tcp", "udp" or "-1"(all). (aws_security_group_invalid_protocol)

  on terraform.tf line 7:
   7:     protocol  = "https"
```

## Why

Apply will fail. (Plan will succeed with the invalid value though)

## How To Fix

Select valid protocol according to the [document](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_SecurityGroupRule.html)
