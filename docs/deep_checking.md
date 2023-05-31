# Deep Checking

_Deep checking_ uses your provider's credentials to apply additional checks that require read access to a target AWS account. TFLint will read AWS configuration from a `provider "aws" {}` block or the same environment variables used by the AWS provider.

For example, the `aws_instance_invalid_iam_profile` rule checks whether a specified IAM profile exists in the target AWS account. This helps detect issues that would result in a failed `terraform plan`.

```console
$ tflint
1 issue(s) found:

Error: "invalid_profile" is invalid IAM profile name. (aws_instance_invalid_iam_profile)

  on template.tf line 4:
   4:   iam_instance_profile = "invalid_profile"

```

You can enable deep checking by enabling `deep_check` in the plugin block:

```hcl
plugin "aws" {
  enabled = true

  deep_check = true
}
```

## Credentials

Credentials can be set in several ways. Each is referenced in the following order:

- Static credentials
- Static credentials (Terraform)
- Environment variables
- Shared credentials
- Shared credentials (Terraform)
- ECS and CodeBuild task roles
- EC2 role

### Static credentials

Access and secret keys can be passed as literals in the plugin or provider configuration:

```hcl
plugin "aws" {
  enabled = true

  deep_check = true
  access_key = "AWS_ACCESS_KEY_ID"
  secret_key = "AWS_SECRET_ACCESS_KEY"
  region     = "us-east-1"
}
```

```hcl
provider "aws" {
  region     = "us-west-2"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

However, committing credentials is not recommended.

### Shared Credentials

If you have [shared credentials](https://docs.aws.amazon.com/sdkref/latest/guide/file-format.html), you can pass a profile name and credentials file path. If omitted, these will be `default` and `~/.aws/credentials`.

```hcl
plugin "aws" {
  enabled = true

  deep_check              = true
  region                  = "us-east-1"
  shared_credentials_file = "~/.aws/myapp"
  profile                 = "AWS_PROFILE"
}
```

```hcl
provider "aws" {
  region                  = "us-west-2"
  shared_credentials_file = "/Users/tf_user/.aws/creds"
  profile                 = "customprofile"
}
```

### Environment Variables

This plugin reads the `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_REGION` environment variables.

```
export AWS_ACCESS_KEY_ID=AWS_ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=AWS_SECRET_KEY
```

### Assume Role

This plugin can assume a role using the provider configuration declared in the target module. See [the provider documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#assume-role) for examples.

You can also specify a role in the plugin configuration:

```hcl
plugin "aws" {
  enabled = true

  deep_check = true

  assume_role {
    role_arn = "arn:aws:iam::123456789012:role/ROLE_NAME"
  }
}
```

## Required Permissions

The following policy document provides the minimal set permissions necessary for deep checking:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeInstances",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeSubnets",
        "ec2:DescribeKeyPairs",
        "ec2:DescribeEgressOnlyInternetGateways",
        "ec2:DescribeInternetGateways",
        "ec2:DescribeNatGateways",
        "ec2:DescribeNetworkInterfaces",
        "ec2:DescribeRouteTables",
        "ec2:DescribeVpcPeeringConnections",
        "rds:DescribeDBSubnetGroups",
        "rds:DescribeOptionGroups",
        "rds:DescribeDBParameterGroups",
        "elasticache:DescribeCacheParameterGroups",
        "elasticache:DescribeCacheSubnetGroups",
        "iam:ListInstanceProfiles"
      ],
      "Resource": "*"
    }
  ]
}
```
