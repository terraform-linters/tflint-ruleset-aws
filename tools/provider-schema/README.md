# Provider Schema

This directory contains a terraform-provider-aws schema JSON file and its configuration files.

This `schema.json` file is used to get information against the terraform-provider-aws schema. See [`/rules/generator-utils/schema.go`](../../rules/generator-utils/schema.go) for more information.

## Update schema file

```sh
tfenv install
terraform init -upgrade
terraform providers schema -json > schema.json
```
