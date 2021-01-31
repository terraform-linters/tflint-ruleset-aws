# TFLint Ruleset for terraform-provider-aws
[![Build Status](https://github.com/terraform-linters/tflint-ruleset-aws/workflows/build/badge.svg?branch=master)](https://github.com/terraform-linters/tflint-ruleset-aws/actions)
[![GitHub release](https://img.shields.io/github/release/terraform-linters/tflint-ruleset-aws.svg)](https://github.com/terraform-linters/tflint-ruleset-aws/releases/latest)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

TFLint ruleset plugin for Terraform AWS Provider

This ruleset focus on possible errors and best practices about AWS resources. Many rules are enabled by default and warn against code that might fail when running `terraform apply`, or clearly unrecommened.

## Requirements

- TFLint v0.24+
- Go v1.15

## Installation

Download the plugin and place it in `~/.tflint.d/plugins/tflint-ruleset-aws` (or `./.tflint.d/plugins/tflint-ruleset-aws`). When using the plugin, configure as follows in `.tflint.hcl`:

```hcl
plugin "aws" {
    enabled = true
}
```

For more configuration about the plugin, see [Plugin Configuration](docs/configuration.md).

**NOTE:** This plugin is bundled with the TFLint binary for backward compatibility, so you can use it without installing it separately. And it is automatically enabled when your Terraform configuration requires AWS provider.

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

700+ rules are available. See [Rules](docs/rules/README.md).

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```
