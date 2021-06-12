## 0.4.2 (2021-06-12)

### Chores

- [#125](https://github.com/terraform-linters/tflint-ruleset-aws/pull/125): Update GoReleaser config to build arm64 for M1
- [#126](https://github.com/terraform-linters/tflint-ruleset-aws/pull/126): Bump GoReleaser version

## 0.4.1 (2021-06-05)

### Enhancements

- [#122](https://github.com/terraform-linters/tflint-ruleset-aws/pull/122): Bump aws-sdk-go submodule and Terraform provider

### Chores

- [#109](https://github.com/terraform-linters/tflint-ruleset-aws/pull/109): Bump github.com/hashicorp/hcl/v2 from 2.9.1 to 2.10.0
- [#113](https://github.com/terraform-linters/tflint-ruleset-aws/pull/113): Bump github.com/zclconf/go-cty from 1.8.1 to 1.8.3
- [#118](https://github.com/terraform-linters/tflint-ruleset-aws/pull/118): Bump actions/cache from 2.1.5 to 2.1.6
- [#119](https://github.com/terraform-linters/tflint-ruleset-aws/pull/119): Bump github.com/google/go-cmp from 0.5.5 to 0.5.6
- [#121](https://github.com/terraform-linters/tflint-ruleset-aws/pull/121): Bump github.com/aws/aws-sdk-go from 1.38.25 to 1.38.55
- [#123](https://github.com/terraform-linters/tflint-ruleset-aws/pull/123): Add notes about auto installation

## 0.4.0 (2021-04-25)

### Enhancements

- [#101](https://github.com/terraform-linters/tflint-ruleset-aws/pull/101): rule: Add aws_api_gateway_model_invalid_name rule
- [#107](https://github.com/terraform-linters/tflint-ruleset-aws/pull/107): Bump aws-sdk-go submodule and Terraform provider

### BugFixes

- [#106](https://github.com/terraform-linters/tflint-ruleset-aws/pull/106): rule: Fix gob error when using map attributes in aws_resource_missing_tags rule

### Chores

- [#99](https://github.com/terraform-linters/tflint-ruleset-aws/pull/99): add TFLINT_PLUGIN_DIR option to README
- [#100](https://github.com/terraform-linters/tflint-ruleset-aws/pull/100): doc: Fix gov_friendly_arns rule naming
- [#103](https://github.com/terraform-linters/tflint-ruleset-aws/pull/103): Bump actions/cache from v2.1.4 to v2.1.5
- [#108](https://github.com/terraform-linters/tflint-ruleset-aws/pull/108): Bump github.com/aws/aws-sdk-go from 1.38.19 to 1.38.25

## 0.3.1 (2021-04-04)

### Enhancements

- [#92](https://github.com/terraform-linters/tflint-ruleset-aws/pull/92): rules: Update aws_mq_broker_invalid_engine_type.go
- [#95](https://github.com/terraform-linters/tflint-ruleset-aws/pull/95): Bump terraform-provider-aws and aws-sdk-go submodule

### BugFixes

- [#85](https://github.com/terraform-linters/tflint-ruleset-aws/pull/85) [#96](https://github.com/terraform-linters/tflint-ruleset-aws/pull/96): Bump tflint-plugin-sdk

### Chores

- [#82](https://github.com/terraform-linters/tflint-ruleset-aws/pull/82): Bump github.com/google/go-cmp from 0.5.4 to 0.5.5
- [#83](https://github.com/terraform-linters/tflint-ruleset-aws/pull/83): Bump github.com/hashicorp/hcl/v2 from 2.9.0 to 2.9.1
- [#86](https://github.com/terraform-linters/tflint-ruleset-aws/pull/86): Bump github.com/zclconf/go-cty from 1.8.0 to 1.8.1
- [#88](https://github.com/terraform-linters/tflint-ruleset-aws/pull/88): Update link to full list of rules
- [#89](https://github.com/terraform-linters/tflint-ruleset-aws/pull/89): Update rules documentation index
- [#91](https://github.com/terraform-linters/tflint-ruleset-aws/pull/91): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.4.4 to 2.5.0
- [#93](https://github.com/terraform-linters/tflint-ruleset-aws/pull/93): Correct spelling
- [#94](https://github.com/terraform-linters/tflint-ruleset-aws/pull/94): Bump github.com/aws/aws-sdk-go from 1.37.25 to 1.38.12

## 0.3.0 (2021-03-06)

### Enhancements

- [#55](https://github.com/terraform-linters/tflint-ruleset-aws/pull/55): Add optional linting rules for govcloud IAM policies
- [#61](https://github.com/terraform-linters/tflint-ruleset-aws/pull/61): Add lint for db_instance's engine attribute
- [#67](https://github.com/terraform-linters/tflint-ruleset-aws/pull/67) [#80](https://github.com/terraform-linters/tflint-ruleset-aws/pull/80) [#81](https://github.com/terraform-linters/tflint-ruleset-aws/pull/81): Bump github.com/aws/aws-sdk-go from 1.37.1 to 1.37.25

### Chores

- [#63](https://github.com/terraform-linters/tflint-ruleset-aws/pull/63): Add document for aws_db_instance_invalid_type
- [#64](https://github.com/terraform-linters/tflint-ruleset-aws/pull/64): Add rule generator
- [#65](https://github.com/terraform-linters/tflint-ruleset-aws/pull/65) [#77](https://github.com/terraform-linters/tflint-ruleset-aws/pull/77): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.4.2 to 2.4.4
- [#68](https://github.com/terraform-linters/tflint-ruleset-aws/pull/68): Bump github.com/golang/mock from 1.4.4 to 1.5.0
- [#71](https://github.com/terraform-linters/tflint-ruleset-aws/pull/71): go: enable module and build caching
- [#73](https://github.com/terraform-linters/tflint-ruleset-aws/pull/73): Upgrade to Go 1.16
- [#74](https://github.com/terraform-linters/tflint-ruleset-aws/pull/74): Bump actions/cache from v2.1.3 to v2.1.4
- [#75](https://github.com/terraform-linters/tflint-ruleset-aws/pull/75): Bump github.com/hashicorp/hcl/v2 from 2.8.2 to 2.9.0
- [#76](https://github.com/terraform-linters/tflint-ruleset-aws/pull/76): Bump github.com/zclconf/go-cty from 1.7.1 to 1.8.0

## 0.2.1 (2021-02-02)

### BugFixes

- [#59](https://github.com/terraform-linters/tflint-ruleset-aws/pull/59): Check EnabledRules instead of Rules

## 0.2.0 (2021-01-31)

The minimum supported version of TFLint has changed in this version. TFLint v0.24.0+ is required for this plugin to work.

### Breaking Changes

- [#58](https://github.com/terraform-linters/tflint-ruleset-aws/pull/58): Bump tflint-plugin-sdk to v0.8.0

### Enhancements

- [#56](https://github.com/terraform-linters/tflint-ruleset-aws/pull/56): Bump github.com/aws/aws-sdk-go from 1.36.19 to 1.37.1
- [#57](https://github.com/terraform-linters/tflint-ruleset-aws/pull/57): Bump aws-sdk-go

### Chores

- [#47](https://github.com/terraform-linters/tflint-ruleset-aws/pull/47): Bump github.com/hashicorp/hcl/v2 from 2.8.1 to 2.8.2
- [#54](https://github.com/terraform-linters/tflint-ruleset-aws/pull/54): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.4.0 to 2.4.2

## 0.1.2 (2021-01-10)

### BugFixes

- [#40](https://github.com/terraform-linters/tflint-ruleset-aws/pull/40): Override RuleNames
- [#43](https://github.com/terraform-linters/tflint-ruleset-aws/pull/43): rule: Allow other fields of rule configs
- [#44](https://github.com/terraform-linters/tflint-ruleset-aws/pull/44): docs: Fix regex in aws_s3_bucket_name
- [#45](https://github.com/terraform-linters/tflint-ruleset-aws/pull/45): Bump tflint-plugin-sdk

## 0.1.1 (2021-01-03)

- [#24](https://github.com/terraform-linters/tflint-ruleset-aws/pull/24): Setup GoReleaser

## 0.1.0 (2021-01-03)

Initial release 🎉

