# Configuration

This plugin can take advantage of additional features by configure the `plugin` block. Currently, this configuration is only available for [Deep Checking](deep_checking.md).

Here's an example:

```hcl
plugin "aws" {
    // Plugin common attributes

    deep_check = false
    access_key = "AWS_ACCESS_KEY_ID"
    secret_key = "AWS_SECRET_ACCESS_KEY"
    region     = "us-east-1"
    profile    = "AWS_PROFILE"
    shared_credentials_file = "~/.aws/credentials"

    assume_role {
        role_arn     = "arn:aws:iam::123456789012:role/ROLE_NAME"
        external_id  = "EXTERNAL_ID"
        policy       = "..."
        session_name = "SESSION_NAME"
    }
}
```

## `deep_check`

Default: false

Enable [Deep Checking](deep_checking.md).

## `access_key`

Default: Credentials declared in the `provider` block or `AWS_ACCESS_KEY_ID` environment variables when the deep checking is enabled.

AWS access key used in the deep checking.

## `secret_key`

Default: Credentials declared in the `provider` block or `AWS_SECRET_ACCESS_KEY` environment variables when the deep checking is enabled.

AWS secret key used in the deep checking.

## `region`

Default: Region declared in the `provider` block or `AWS_REGION` environment variables when the deep checking is enabled.

AWS region used in the deep checking.

## `profile`

Default: Profile declared in the `provider` block or `AWS_PROFILE` environment variables when the deep checking is enabled.

AWS shared credentials profile name used in the deep checking.

## `shared_credentials_file`

Default: Profile declared in the `provider` block or `~/.aws/credentials` when the deep checking is enabled.

AWS shared credentials file path used in the deep checking.

## `assume_role`

Default: Assume role config declared in the `provider` block.

AWS assume role config used in the deep checking.
