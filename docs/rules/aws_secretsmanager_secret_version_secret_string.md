# aws_secretsmanager_secret_version_secret_string

Disallows using the `secret_string` attribute, in favor of the `secret_string_wo` attribute.

## Example

```hcl
resource "aws_secretsmanager_secret_version" "test" {
  secret_string = var.secret
}
```

```
$ tflint
1 issue(s) found:

Warning: "secret_string" is a non-ephemeral attribute, which means this secret is stored in state. Please use "secret_string_wo". (aws_secretsmanager_secret_version_secret_string)

  on test.tf line 3:
   3:   secret_string = var.secret

```

## Why

Saving secrets to state or plan files is a bad practice. It can cause serious security issues. Keeping secrets from these files is possible in most of the cases by using write-only attributes.

## How To Fix

Replace the `secret_string` with `secret_string_wo` and `secret_string_wo_version`, either via using an ephemeral resource as source or using an ephemeral variable:

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