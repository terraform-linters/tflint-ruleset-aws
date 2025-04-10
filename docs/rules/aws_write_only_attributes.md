# aws_write_only_attributes

Disallows using the normal attribute containing sensitive information, in favor of an available write-only attribute. This is only valid for Terraform version 1.11.0 and later. If there alternative options available to prevent this information ending up in state, these are suggested as well.

## Example

In this example `aws_secretsmanager_secret_version` is used, but you can replace that with any resource with write-only attributes:

```hcl
resource "aws_secretsmanager_secret_version" "test" {
  secret_string = var.secret
}
```

```
$ tflint
1 issue(s) found:

Warning: [Fixable] "secret_string" is a non-ephemeral attribute, which means this secret is stored in state. Please use "secret_string_wo". (aws_write_only_attributes)

  on test.tf line 3:
   3:   secret_string = var.secret

```

## Why

Saving secrets to state or plan files is a bad practice. It can cause serious security issues. Keeping secrets from these files is possible in most of the cases by using write-only attributes.

## How To Fix

Replace the `secret_string` with `secret_string_wo`, preferable via using an ephemeral resource as source or using an ephemeral variable:

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
  ephemeral   = true # Optional, non-ephemeral values can also be used for write-only attributes
  description = "Input variable for a secret"
}

resource "aws_secretsmanager_secret_version" "test" {
  secret_string_wo         = var.test
  secret_string_wo_version = 1
}
```