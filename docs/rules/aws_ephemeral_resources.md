# aws_ephemeral_resources

Recommends using available [ephemeral resources](https://developer.hashicorp.com/terraform/language/resources/ephemeral/reference) instead of the original data source. This is only valid for Terraform v1.10+.

## Example

This example uses `aws_secretsmanager_random_password`, but the rule applies to all data sources with an ephemeral equivalent:

```hcl
data "aws_secretsmanager_random_password" "test" {
  password_length = 50
  exclude_numbers = true
}
```

```
$ tflint
1 issue(s) found:

Warning: "aws_secretsmanager_random_password" is a non-ephemeral data source, which means that all (sensitive) attributes are stored in state. Please use ephemeral resource "aws_secretsmanager_random_password" instead. (aws_ephemeral_resources)

  on test.tf line 2:
   2: data "aws_secretsmanager_random_password" "test"

```

## Why

By default, sensitive attributes are still stored in state, just hidden from view in plan output. Other resources are able to refer to these attributes. Current versions of Terraform also include support for ephemeral resources, which are not persisted to state. Other resources can refer to their values, but executing of the lookup is defered until the apply stage.

Using ephemeral resources mitigates the risk of a malicious actor obtaining privileged credentials by accessing Terraform state files directly. Prefer using them over the original data sources for sensitive data.

## How To Fix

Replace the data source with its ephemeral resource equivalent. Use resources with write-only arguments or in provider configuration to ensure that the sensitive value is not persisted to state.

In case of the previously shown `aws_secretsmanager_random_password` data source, replace `data` by `ephemeral`:

```hcl
ephemeral "aws_secretsmanager_random_password" "test" {
  password_length = 50
  exclude_numbers = true
}
```
