# WORKING IN PROGRESS AND MAINTAINER WANTED!

This project was started to make TFLint pluggable. Everything is working in progress and not available as a plugin.

AWS provider rules are currently integrated into the TFLint core. It will be cut out to this repository in the future, but it doesn't mean that limit the proposal of rules about AWS provider into the core repository.

The migration process will take a long time. Please open an issue to the core repository for the latest proposal.

# TFLint Ruleset for terraform-provider-aws

TFLint ruleset plugin for Terraform AWS Provider

## Requirements

- TFLint v0.14+
- Go v1.13

## Installation

Download the plugin and place it in `~/.tflint.d/plugins/tflint-ruleset-aws` (or `./.tflint.d/plugins/tflint-ruleset-aws`). When using the plugin, configure as follows in `.tflint.hcl`:

```hcl
plugin "aws" {
    enabled = true
}
```

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |
|aws_instance_example_type|Show instance type|ERROR|✔||

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```
