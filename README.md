# TFLint Ruleset for terraform-provider-aws
[![Build Status](https://github.com/terraform-linters/tflint-ruleset-aws/workflows/build/badge.svg?branch=master)](https://github.com/terraform-linters/tflint-ruleset-aws/actions)
[![GitHub release](https://img.shields.io/github/release/terraform-linters/tflint-ruleset-aws.svg)](https://github.com/terraform-linters/tflint-ruleset-aws/releases/latest)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

TFLint ruleset plugin for Terraform AWS Provider

This ruleset focus on possible errors and best practices about AWS resources. Many rules are enabled by default and warn against code that might fail when running `terraform apply`, or clearly unrecommened.

## Requirements

- TFLint v0.23+
- Go v1.15

## Installation

Download the plugin and place it in `~/.tflint.d/plugins/tflint-ruleset-aws` (or `./.tflint.d/plugins/tflint-ruleset-aws`). When using the plugin, configure as follows in `.tflint.hcl`:

```hcl
plugin "aws" {
    enabled = true
}
```

For more configuration about the plugin, see [Plugin Configuration](docs/configuration.md).

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
