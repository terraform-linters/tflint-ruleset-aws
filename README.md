# TFLint Ruleset Template

This is a template repository for building a custom ruleset. You can create a plugin repository from "Use this template".

## Requirements

- TFLint v0.14+
- Go v1.13

## Installation

Download the plugin and place it in `~/.tflint.d/plugins/tflint-ruleset-template` (or `./.tflint.d/plugins/tflint-ruleset-template`). When using the plugin, configure as follows in `.tflint.hcl`:

```hcl
plugin "template" {
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
