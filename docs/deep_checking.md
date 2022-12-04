# Deep Checking

Deep Checking uses your provider's credentials to perform a more strict inspection.

For example, if the IAM profile references something that doesn't exist, terraform apply will fail, which can't be found by general validation. Deep Checking solves this problem.

```console
$ tflint
1 issue(s) found:

Error: "invalid_profile" is invalid IAM profile name. (aws_instance_invalid_iam_profile)

  on template.tf line 4:
   4:   iam_instance_profile = "invalid_profile"

```

You can enable Deep Checking by changing the plugin configuration.

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

If you have an access key and a secret key, you can pass these keys like the following:

```hcl
plugin "aws" {
  enabled = true

  deep_check = true
  access_key = "AWS_ACCESS_KEY_ID"
  secret_key = "AWS_SECRET_ACCESS_KEY"
  region     = "us-east-1"
}
```

Although there is not recommended, if an access key is hard-coded in a provider configuration, they will also be taken into account. The priority is higher than the environment variable and lower than the above way.

```hcl
provider "aws" {
  region     = "us-west-2"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

### Shared credentials

If you have [shared credentials](https://aws.amazon.com/jp/blogs/security/a-new-and-standardized-way-to-manage-credentials-in-the-aws-sdks/), you can pass a profile name and credentials file path. If omitted, these will be `default` and `~/.aws/credentials`.

```hcl
plugin "aws" {
  enabled = true

  deep_check              = true
  region                  = "us-east-1"
  shared_credentials_file = "~/.aws/myapp"
  profile                 = "AWS_PROFILE"
}
```

If these configurations are defined in the provider block, they will also be taken into account. But the priority is lower than the above way.

```hcl
provider "aws" {
  region                  = "us-west-2"
  shared_credentials_file = "/Users/tf_user/.aws/creds"
  profile                 = "customprofile"
}
```

### Environment variables

This plugin looks up `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`, `AWS_REGION` environment variables. This is useful when you don't want to explicitly pass credentials.

```
$ export AWS_ACCESS_KEY_ID=AWS_ACCESS_KEY
$ export AWS_SECRET_ACCESS_KEY=AWS_SECRET_KEY
```

### Role-based authentication

This plugin fetches credentials in the same way as Terraform. See [this documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#codebuild-ecs-and-eks-roles) for the role-based authentication.

### Assume role

This plugin can assume a role in the same way as Terraform. See [this documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#assume-role).

You can also declare the assume role config in the plugin config:

```hcl
plugin "aws" {
  enabled = true

  deep_check = true

  assume_role {
    role_arn = "arn:aws:iam::123456789012:role/ROLE_NAME"
  }
}
```

## Required permissions

The following policy document provides the minimal set permissions necessary for the deep checking:

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
