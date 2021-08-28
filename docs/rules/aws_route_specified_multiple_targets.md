# aws_route_specified_multiple_targets

Disallow routes that have multiple targets.

## Example

```hcl
resource "aws_route" "foo" {
  route_table_id         = "rtb-1234abcd"
  destination_cidr_block = "10.0.1.0/22"
  gateway_id             = "igw-1234abcd"
  egress_only_gateway_id = "eigw-1234abcd" # second routing target?
}
```

```
$ tflint
1 issue(s) found:

Error: More than one routing target specified. It must be one. (aws_route_specified_multiple_targets)

  on template.tf line 1:
   1: resource "aws_route" "foo" {
 
```

## Why

It results in an error.

## How To Fix

Remove the other targets so that there is only one target. See https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route#argument-reference
