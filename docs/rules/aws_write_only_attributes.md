# aws_write_only_arguments

Recommends using available [write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral/write-only) instead of the original sensitive attribute. This is only valid for Terraform v1.11+.

## Example

This example uses `aws_secretsmanager_secret_version`, but the rule applies to all resources with write-only arguments:

```hcl
resource "aws_secretsmanager_secret_version" "test" {
  secret_string = var.secret
}
```

```
$ tflint
1 issue(s) found:

Warning: [Fixable] "secret_string" is a non-ephemeral attribute, which means this secret is stored in state. Please use "secret_string_wo". (aws_write_only_arguments)

  on test.tf line 3:
   3:   secret_string = var.secret

```

## Why

By default, sensitive attributes are still stored in state, just hidden from view in plan output. Other resources are able to refer to these attributes. Current versions of Terraform also include support for write-only arguments, which are not persisted to state. Other resources cannot refer to their values.

Using write-only arguments mitigates the risk of a malicious actor obtaining privileged credentials by accessing Terraform state files directly. Prefer using them over the original sensitive attribute unless you need to refer to it in other blocks, such as a [root `output`](https://developer.hashicorp.com/terraform/language/values/outputs#ephemeral-avoid-storing-values-in-state-or-plan-files), that cannot be ephemeral.

## How To Fix

Replace the attribute with its write-only argument equivalent. Reference an ephemeral resource or ephemeral variable to ensure that the sensitive value is not persisted to state.

```hcl
ephemeral "random_password" "test" {
  length           = 32
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_secretsmanager_secret_version" "test" {
  secret_string_wo         = ephemeral.random_password.test.value
  secret_string_wo_version = 1
}
```

```hcl
variable "test" {
  type        = string
  ephemeral   = true # Optional, non-ephemeral values can also be used for write-only arguments
  description = "Input variable for a secret"
}

resource "aws_secretsmanager_secret_version" "test" {
  secret_string_wo         = var.test
  secret_string_wo_version = 1
}
```