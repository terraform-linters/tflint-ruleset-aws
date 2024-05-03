# aws_provider_missing_default_tags

Require specific tags for all AWS provider default_tags blocks.

## Configuration

```hcl
rule "aws_provider_missing_default_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}
```

## Examples

The `tags` attribute is a map of `key`=`value` pairs, like resource tags:

```hcl
provider "aws" {
  default_tags {
    tags = {
      foo = "Bar"
      bar = "Baz"
    }
  }
}
```

An issue is reported because the tag keys ["foo", "bar"] don't match the required tags ["Foo", "Bar"]

```
$ tflint
1 issue(s) found:

Notice: The provider is missing the following tags: "Bar", "Foo". (aws_provider_missing_default_tags)

  on test.tf line 1:
   1: provider "aws" {
```

## Why

- You want to set a standardized set of tags for your AWS resources via the provider default tags.
- You want more DRY (don't repeat yourself) Terraform code, by eliminating tag configuration from resources.
- Using default tags results in better tagging coverage. The resource missing tags rule needs support
  to be added for non-standard uses of tags in the provider, for example EC2 root block devices.

Use this rule in conjuction with aws_resource_missing_tags_rule, for example to enforce common tags and
resource specific tags, without duplicating tags.

```hcl
rule "aws_resource_missing_tags" {
  enabled = false
  tags = [
    "kubernetes.io/cluster/eks",
  ]
  include = [
    "aws_subnet",
  ]
}

rule "aws_provider_missing_default_tags" {
  enabled = true
  tags = [
    "CostCenter",
  ]
}
```

## How To Fix

Ensure the provider default tags contains all your required tags.
