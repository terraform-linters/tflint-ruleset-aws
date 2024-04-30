# aws_provider_missing_tags

Require specific tags for all AWS provider default_tags blocks.

## Configuration

```hcl
rule "aws_provider_missing_tags" {
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

Notice: The provider is missing the following tags: "Bar", "Foo". (aws_provider_missing_tags)

  on test.tf line 1:
   1: provider "aws" {
```

## Why

You want to set a standardized set of tags for your AWS resources via the provider default tags.

## How To Fix

Ensure the provider default tags contains all your required tags.
