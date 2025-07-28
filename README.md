# TFLint Ruleset for terraform-provider-aws
[![Build Status](https://github.com/terraform-linters/tflint-ruleset-aws/actions/workflows/build.yml/badge.svg?branch=master)](https://github.com/terraform-linters/tflint-ruleset-aws/actions)
[![GitHub release](https://img.shields.io/github/release/terraform-linters/tflint-ruleset-aws.svg)](https://github.com/terraform-linters/tflint-ruleset-aws/releases/latest)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

TFLint ruleset plugin for Terraform AWS Provider

This ruleset focus on possible errors and best practices about AWS resources. In other words, the rules warn against code that might fail when running `terraform apply` or practices that should be avoided.

## Requirements

- TFLint v0.42+
- Go v1.24

## Installation

After configuring [TFLint](https://github.com/terraform-linters/tflint/blob/master/docs/user-guide/config.md), you can install the plugin by adding a config to `.tflint.hcl` and running `tflint --init`:

```hcl
plugin "aws" {
    enabled = true
    version = "0.40.0"
    source  = "github.com/terraform-linters/tflint-ruleset-aws"
}
```

For more information about configuration options for this ruleset plugin, see [Plugin Configuration](docs/configuration.md).

## Getting Started

Terraform is a great tool for Infrastructure as Code. However, many of these tools don't validate provider-specific issues. For example, see the following configuration file:

```hcl
resource "aws_instance" "foo" {
  ami           = "ami-0ff8a91507f77f867"
  instance_type = "t1.2xlarge" # invalid type!
}
```

Since `t1.2xlarge` is an invalid instance type, an error will occur when you run `terraform apply`. But `terraform validate` and `terraform plan` cannot find this possible error in advance. That's because it's an AWS provider-specific issue and it's valid as the Terraform Language.

The goal of this ruleset is to find such errors:

![demo](docs/assets/demo.gif)

By running TFLint with this ruleset in advance, you can fix the problem before the error occurs in production CI/CD pipelines.

## Rules

There are 700+ rules available, see the [Rules documentation page](docs/rules/README.md) for a complete list. Note that not all of them are enabled by default and need to be configured manually (especially rules involving [Best Practices](docs/rules/README.md#best-practicesnaming-conventions)).

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

Note that if you install the plugin with `make install`, you must omit the `version` and `source` attributes in` .tflint.hcl`:

```hcl
plugin "aws" {
    enabled = true
}
```

## Add a new rule

If you are interested in adding a new rule to this ruleset, you can use the generator. Run the following command:

```
$ go run ./rules/generator
```

Follow the instructions to edit the generated files and open a new pull request.
