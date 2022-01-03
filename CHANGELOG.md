## 0.11.0 (2022-01-04)

Many thanks to @PatMyron, a new maintainer! This release adds 589 SDK-based validation rules, significantly increasing resource coverage.

### Enhancements

- [#216](https://github.com/terraform-linters/tflint-ruleset-aws/pull/216): Bump github.com/aws/aws-sdk-go from 1.42.19 to 1.42.23
- [#217](https://github.com/terraform-linters/tflint-ruleset-aws/pull/217): rules: Add ap-southeast-3 (Jakarta) as valid region ([@PatMyron](https://github.com/PatMyron))
- [#218](https://github.com/terraform-linters/tflint-ruleset-aws/pull/218): rules: Add mapping `aws_acmpca_certificate` and `aws_acmpca_certificate_authority_certificate` ([@PatMyron](https://github.com/PatMyron))
- [#219](https://github.com/terraform-linters/tflint-ruleset-aws/pull/219): rules: Add mapping `aws_acm_certificate` and `aws_acm_certificate_validation` ([@PatMyron](https://github.com/PatMyron))
- [#220](https://github.com/terraform-linters/tflint-ruleset-aws/pull/220): rules: Add mapping `aws_api_gateway_authorizer`, `aws_api_gateway_documentation_part`, `aws_api_gateway_domain_name` ([@PatMyron](https://github.com/PatMyron))
- [#222](https://github.com/terraform-linters/tflint-ruleset-aws/pull/222): rules: Add mapping `aws_apigatewayv2` ([@PatMyron](https://github.com/PatMyron))
- [#223](https://github.com/terraform-linters/tflint-ruleset-aws/pull/223): rules: Add mapping `aws_accessanalyzer_analyzer` ([@PatMyron](https://github.com/PatMyron))
- [#224](https://github.com/terraform-linters/tflint-ruleset-aws/pull/224): rules: Add mapping `aws_account_alternate_contact` ([@PatMyron](https://github.com/PatMyron))
- [#225](https://github.com/terraform-linters/tflint-ruleset-aws/pull/225): rules: Add mapping `aws_amplify` ([@PatMyron](https://github.com/PatMyron))
- [#226](https://github.com/terraform-linters/tflint-ruleset-aws/pull/226): rules: Add mapping `aws_appconfig` ([@PatMyron](https://github.com/PatMyron))
- [#227](https://github.com/terraform-linters/tflint-ruleset-aws/pull/227): rules: Add mapping `aws_appmesh` ([@PatMyron](https://github.com/PatMyron))
- [#228](https://github.com/terraform-linters/tflint-ruleset-aws/pull/228): rules: Add mapping `aws_apprunner` ([@PatMyron](https://github.com/PatMyron))
- [#232](https://github.com/terraform-linters/tflint-ruleset-aws/pull/232): rules: Add mapping `aws_appstream` ([@PatMyron](https://github.com/PatMyron))
- [#233](https://github.com/terraform-linters/tflint-ruleset-aws/pull/233): rules: Add mapping `aws_backup` ([@PatMyron](https://github.com/PatMyron))
- [#234](https://github.com/terraform-linters/tflint-ruleset-aws/pull/234): rules: Add mapping `aws_chime` ([@PatMyron](https://github.com/PatMyron))
- [#235](https://github.com/terraform-linters/tflint-ruleset-aws/pull/235): rules: Add mapping `aws_cloudfront` ([@PatMyron](https://github.com/PatMyron))
- [#236](https://github.com/terraform-linters/tflint-ruleset-aws/pull/236): rules: Add mapping `aws_codeartifact` ([@PatMyron](https://github.com/PatMyron))
- [#237](https://github.com/terraform-linters/tflint-ruleset-aws/pull/237): rules: Add mapping `aws_config` ([@PatMyron](https://github.com/PatMyron))
- [#238](https://github.com/terraform-linters/tflint-ruleset-aws/pull/238): rules: Add mapping `aws_connect` ([@PatMyron](https://github.com/PatMyron))
- [#239](https://github.com/terraform-linters/tflint-ruleset-aws/pull/239): rules: Add mapping `aws_dx` (Direct Connect) ([@PatMyron](https://github.com/PatMyron))
- [#240](https://github.com/terraform-linters/tflint-ruleset-aws/pull/240): rules: Add mapping `aws_sagemaker` ([@PatMyron](https://github.com/PatMyron))
- [#241](https://github.com/terraform-linters/tflint-ruleset-aws/pull/241): rules: Add mapping `servicecatalog` ([@PatMyron](https://github.com/PatMyron))
- [#242](https://github.com/terraform-linters/tflint-ruleset-aws/pull/242): rules: Add mapping `aws_glue` ([@PatMyron](https://github.com/PatMyron))
- [#243](https://github.com/terraform-linters/tflint-ruleset-aws/pull/243): rules: Add mapping `aws_securityhub` ([@PatMyron](https://github.com/PatMyron))
- [#244](https://github.com/terraform-linters/tflint-ruleset-aws/pull/244): rules: Add mapping `macie2` ([@PatMyron](https://github.com/PatMyron))
- [#245](https://github.com/terraform-linters/tflint-ruleset-aws/pull/245): rules: Add mapping `s3control` ([@PatMyron](https://github.com/PatMyron))
- [#247](https://github.com/terraform-linters/tflint-ruleset-aws/pull/247): rules: Add mapping `imagebuilder` ([@PatMyron](https://github.com/PatMyron))
- [#248](https://github.com/terraform-linters/tflint-ruleset-aws/pull/248): rules: Add mapping `wafv2` ([@PatMyron](https://github.com/PatMyron))
- [#249](https://github.com/terraform-linters/tflint-ruleset-aws/pull/249): rules: Add mapping `aws_networkfirewall` ([@PatMyron](https://github.com/PatMyron))
- [#250](https://github.com/terraform-linters/tflint-ruleset-aws/pull/250): rules: Add mapping `aws_vpc` ([@PatMyron](https://github.com/PatMyron))
- [#251](https://github.com/terraform-linters/tflint-ruleset-aws/pull/251): rules: Add mapping `aws_ec2` ([@PatMyron](https://github.com/PatMyron))
- [#252](https://github.com/terraform-linters/tflint-ruleset-aws/pull/252): rules: Add mapping `aws_route53_resolver` ([@PatMyron](https://github.com/PatMyron))
- [#253](https://github.com/terraform-linters/tflint-ruleset-aws/pull/253): rules: Add mapping `aws_cloudwatch_event` ([@PatMyron](https://github.com/PatMyron))
- [#254](https://github.com/terraform-linters/tflint-ruleset-aws/pull/254): rules: Add mapping `aws_s3` ([@PatMyron](https://github.com/PatMyron))
- [#255](https://github.com/terraform-linters/tflint-ruleset-aws/pull/255): rules: Add mapping `aws_db_proxy` ([@PatMyron](https://github.com/PatMyron))
- [#256](https://github.com/terraform-linters/tflint-ruleset-aws/pull/256): rules: Add mapping `aws_ecr` ([@PatMyron](https://github.com/PatMyron))
- [#257](https://github.com/terraform-linters/tflint-ruleset-aws/pull/257): rules: Add mapping `aws_ecs` ([@PatMyron](https://github.com/PatMyron))
- [#258](https://github.com/terraform-linters/tflint-ruleset-aws/pull/258): rules: Add mapping `aws_eks` ([@PatMyron](https://github.com/PatMyron))
- [#259](https://github.com/terraform-linters/tflint-ruleset-aws/pull/259) [#289](https://github.com/terraform-linters/tflint-ruleset-aws/pull/289): tools: Bump Terraform and provider schema ([@PatMyron](https://github.com/PatMyron), [@wata727](https://github.com/wata727))
- [#260](https://github.com/terraform-linters/tflint-ruleset-aws/pull/260): rules: Add mapping `aws_fsx` ([@PatMyron](https://github.com/PatMyron))
- [#261](https://github.com/terraform-linters/tflint-ruleset-aws/pull/261): rules: Add mapping `aws_guardduty` ([@PatMyron](https://github.com/PatMyron))
- [#262](https://github.com/terraform-linters/tflint-ruleset-aws/pull/262): rules: Add mapping `aws_lambda` ([@PatMyron](https://github.com/PatMyron))
- [#263](https://github.com/terraform-linters/tflint-ruleset-aws/pull/263): rules: Add mapping `aws_ssoadmin` ([@PatMyron](https://github.com/PatMyron))
- [#264](https://github.com/terraform-linters/tflint-ruleset-aws/pull/264): rules: Add mapping `aws_route53recoverycontrolconfig` ([@PatMyron](https://github.com/PatMyron))
- [#265](https://github.com/terraform-linters/tflint-ruleset-aws/pull/265): rules: Add mapping `aws_route53recoveryreadiness` ([@PatMyron](https://github.com/PatMyron))
- [#266](https://github.com/terraform-linters/tflint-ruleset-aws/pull/266): rules: Add mapping `aws_efs` ([@PatMyron](https://github.com/PatMyron))
- [#267](https://github.com/terraform-linters/tflint-ruleset-aws/pull/267): rules: Add mapping `aws_elasticache` ([@PatMyron](https://github.com/PatMyron))
- [#268](https://github.com/terraform-linters/tflint-ruleset-aws/pull/268): rules: Add mapping `aws_lakeformation` ([@PatMyron](https://github.com/PatMyron))
- [#269](https://github.com/terraform-linters/tflint-ruleset-aws/pull/269): rules: Add mapping `aws_prometheus` ([@PatMyron](https://github.com/PatMyron))
- [#270](https://github.com/terraform-linters/tflint-ruleset-aws/pull/270): rules: Add mapping `aws_quicksight` ([@PatMyron](https://github.com/PatMyron))
- [#271](https://github.com/terraform-linters/tflint-ruleset-aws/pull/271): rules: Add mapping `aws_schemas` ([@PatMyron](https://github.com/PatMyron))
- [#272](https://github.com/terraform-linters/tflint-ruleset-aws/pull/272): rules: Add mapping `aws_signer` ([@PatMyron](https://github.com/PatMyron))
- [#273](https://github.com/terraform-linters/tflint-ruleset-aws/pull/273): rules: Add mapping `aws_storagegateway` ([@PatMyron](https://github.com/PatMyron))
- [#275](https://github.com/terraform-linters/tflint-ruleset-aws/pull/275): rules: Add mapping `aws_codecommit` ([@PatMyron](https://github.com/PatMyron))
- [#276](https://github.com/terraform-linters/tflint-ruleset-aws/pull/276): rules: Add mapping `aws_datasync` ([@PatMyron](https://github.com/PatMyron))
- [#277](https://github.com/terraform-linters/tflint-ruleset-aws/pull/277): rules: Add mapping `aws_dynamodb` ([@PatMyron](https://github.com/PatMyron))
- [#278](https://github.com/terraform-linters/tflint-ruleset-aws/pull/278): rules: Add mapping `aws_kms` ([@PatMyron](https://github.com/PatMyron))
- [#279](https://github.com/terraform-linters/tflint-ruleset-aws/pull/279): rules: Add mapping `aws_secretsmanager` ([@PatMyron](https://github.com/PatMyron))
- [#280](https://github.com/terraform-linters/tflint-ruleset-aws/pull/280): rules: Add mapping `aws_xray` ([@PatMyron](https://github.com/PatMyron))
- [#281](https://github.com/terraform-linters/tflint-ruleset-aws/pull/281): rules: Add mapping `aws_codestarconnections` ([@PatMyron](https://github.com/PatMyron))
- [#282](https://github.com/terraform-linters/tflint-ruleset-aws/pull/282): rules: Add mapping `aws_ecrpublic` ([@PatMyron](https://github.com/PatMyron))
- [#283](https://github.com/terraform-linters/tflint-ruleset-aws/pull/283): rules: Add mapping `aws_timestreamwrite` ([@PatMyron](https://github.com/PatMyron))
- [#284](https://github.com/terraform-linters/tflint-ruleset-aws/pull/284): rules: Add mapping `aws_kinesisanalyticsv2` ([@PatMyron](https://github.com/PatMyron))
- [#285](https://github.com/terraform-linters/tflint-ruleset-aws/pull/285): rules: Add mapping `aws_workspaces` ([@PatMyron](https://github.com/PatMyron))
- [#287](https://github.com/terraform-linters/tflint-ruleset-aws/pull/287): rules: Add mapping `aws_emr` ([@PatMyron](https://github.com/PatMyron))
- [#288](https://github.com/terraform-linters/tflint-ruleset-aws/pull/288): rules: Add mapping other services ([@PatMyron](https://github.com/PatMyron))

### Chores

- [#230](https://github.com/terraform-linters/tflint-ruleset-aws/pull/230): docs: remove irrelevant issue from deep check example ([@PatMyron](https://github.com/PatMyron))
- [#231](https://github.com/terraform-linters/tflint-ruleset-aws/pull/231): docs: documenting undocumented rules ([@PatMyron](https://github.com/PatMyron))
- [#246](https://github.com/terraform-linters/tflint-ruleset-aws/pull/246): Bump github.com/aws/aws-sdk-go from 1.42.23 to 1.42.25

## 0.10.1 (2021-12-10)

### Enhancements

- [#211](https://github.com/terraform-linters/tflint-ruleset-aws/pull/211): rules: Add missing t4g ElastiCache node types ([@acastro2](https://github.com/acastro2))
- [#213](https://github.com/terraform-linters/tflint-ruleset-aws/pull/213): rules: Add data tiering node types for ElastiCache ([@wata727](https://github.com/wata727))

### Chores

- [#210](https://github.com/terraform-linters/tflint-ruleset-aws/pull/210): Bump github.com/hashicorp/hcl/v2 from 2.10.1 to 2.11.1
- [#214](https://github.com/terraform-linters/tflint-ruleset-aws/pull/214): Extract ElastiCache node types to utils ([@wata727](https://github.com/wata727))

## 0.10.0 (2021-12-04)

### Enhancements

- [#202](https://github.com/terraform-linters/tflint-ruleset-aws/pull/202): rules: Add acm certificate lifecycle create before destroy rule ([@AleksaC](https://github.com/AleksaC))
- [#208](https://github.com/terraform-linters/tflint-ruleset-aws/pull/208): Bump aws-sdk-go submodule and Terraform provider schema ([@wata727](https://github.com/wata727))

### Chores

- [#199](https://github.com/terraform-linters/tflint-ruleset-aws/pull/199): Bump github.com/zclconf/go-cty from 1.9.1 to 1.10.0
- [#204](https://github.com/terraform-linters/tflint-ruleset-aws/pull/204): Bump github.com/hashicorp/aws-sdk-go-base from 0.7.1 to 1.0.0
- [#209](https://github.com/terraform-linters/tflint-ruleset-aws/pull/209): Bump github.com/aws/aws-sdk-go from 1.41.19 to 1.42.19

## 0.9.0 (2021-11-06)

### Breaking Changes

- [#198](https://github.com/terraform-linters/tflint-ruleset-aws/pull/198): build: Remove unsupported build targets ([@wata727](https://github.com/wata727))

### Enhancements

- [#197](https://github.com/terraform-linters/tflint-ruleset-aws/pull/197): Bump aws-sdk-go submodule and Terraform provider schema ([@wata727](https://github.com/wata727))

### BugFixes

- [#195](https://github.com/terraform-linters/tflint-ruleset-aws/pull/195): rules: Use EmitIssueOnExpr when emitting an issue for an expression ([@wata727](https://github.com/wata727))

### Chores

- [#196](https://github.com/terraform-linters/tflint-ruleset-aws/pull/196): Bump github.com/aws/aws-sdk-go from 1.40.54 to 1.41.19

## 0.8.0 (2021-10-11)

### Enhancements

- [#159](https://github.com/terraform-linters/tflint-ruleset-aws/pull/159): rules: Add `aws_iam_group_policy_too_long` rule ([@Rihoj](https://github.com/Rihoj))
- [#187](https://github.com/terraform-linters/tflint-ruleset-aws/pull/187): rules: Add Aurora Graviton2-based T4g and X2g instances ([@wata727](https://github.com/wata727))
- [#188](https://github.com/terraform-linters/tflint-ruleset-aws/pull/188): Bump aws-sdk-go submodule and Terraform provider schema ([@wata727](https://github.com/wata727))

### Chores

- [#185](https://github.com/terraform-linters/tflint-ruleset-aws/pull/185): Bump github.com/aws/aws-sdk-go from 1.40.32 to 1.40.54

## 0.7.2 (2021-09-25)

### BugFixes

- [#183](https://github.com/terraform-linters/tflint-ruleset-aws/pull/183): rules: Allow unbracketed IAM policy statements ([@wata727](https://github.com/wata727))

### Chores

- [#180](https://github.com/terraform-linters/tflint-ruleset-aws/pull/180): build: Go 1.17 ([@wata727](https://github.com/wata727))

## 0.7.1 (2021-09-03)

### BugFixes

- [#176](https://github.com/terraform-linters/tflint-ruleset-aws/pull/176): rules: Fix false positive for IAM policy document without Sid ([@wata727](https://github.com/wata727))
- [#178](https://github.com/terraform-linters/tflint-ruleset-aws/pull/178): rules: Fix an error when policy is not evaluable

### Chores

- [#175](https://github.com/terraform-linters/tflint-ruleset-aws/pull/175): rules: Fix broken `aws_lambda_function_deprecated_runtime` test ([@wata727](https://github.com/wata727))

## 0.7.0 (2021-08-28)

### Breaking Changes

- [#158](https://github.com/terraform-linters/tflint-ruleset-aws/pull/158): rules: Remove `aws_cognito_user_pool_domain_invalid_domain` rule ([@wata727](https://github.com/wata727))

### Enhancements

- [#152](https://github.com/terraform-linters/tflint-ruleset-aws/pull/152): rules: Add deep checking rules for the `aws_elasticache_replication_group` resource ([@Rihoj](https://github.com/Rihoj))
  - `aws_elasticache_replication_group_invalid_parameter_group`
  - `aws_elasticache_replication_group_invalid_security_group`
  - `aws_elasticache_replication_group_invalid_subnet_group`
- [#153](https://github.com/terraform-linters/tflint-ruleset-aws/pull/153): rules: Add `aws_iam_policy_too_long_policy` rule ([@Rihoj](https://github.com/Rihoj))
- [#154](https://github.com/terraform-linters/tflint-ruleset-aws/pull/154): rules: Add `aws_lambda_function_deprecated_runtime` rule ([@Rihoj](https://github.com/Rihoj))
- [#155](https://github.com/terraform-linters/tflint-ruleset-aws/pull/155): rules: Add `aws_iam_policy_sid_invalid_characters` rule ([@Rihoj](https://github.com/Rihoj))
- [#166](https://github.com/terraform-linters/tflint-ruleset-aws/pull/166): rules: Update valid DB instance types ([@wata727](https://github.com/wata727))
- [#167](https://github.com/terraform-linters/tflint-ruleset-aws/pull/167): rules: Add support for Oracle multitenant container database engines ([@wata727](https://github.com/wata727))
- [#168](https://github.com/terraform-linters/tflint-ruleset-aws/pull/168): rules: Add `RabbitMQ` to `aws_mq_configuration_invalid_engine_type` rule ([@wata727](https://github.com/wata727))
- [#169](https://github.com/terraform-linters/tflint-ruleset-aws/pull/169): rules: Add route target types for `aws_route` rules ([@wata727](https://github.com/wata727))
- [#170](https://github.com/terraform-linters/tflint-ruleset-aws/pull/170): rules: Update valid regions for the `aws_s3_bucket_invalid_region` rule ([@wata727](https://github.com/wata727))
- [#171](https://github.com/terraform-linters/tflint-ruleset-aws/pull/171): rules: Bump aws-sdk-go submodule and Terraform provider schema

### Chores

- [#157](https://github.com/terraform-linters/tflint-ruleset-aws/pull/157): docs: Add required IAM policy document for deep checking ([@wata727](https://github.com/wata727))
- [#161](https://github.com/terraform-linters/tflint-ruleset-aws/pull/161): docs: Fix typo in `aws_elasticache_cluster_default_parameter_group.md` ([@w0rmr1d3r](https://github.com/w0rmr1d3r))
- [#163](https://github.com/terraform-linters/tflint-ruleset-aws/pull/163): Bump github.com/zclconf/go-cty from 1.9.0 to 1.9.1
- [#172](https://github.com/terraform-linters/tflint-ruleset-aws/pull/172): Bump github.com/aws/aws-sdk-go from 1.40.17 to 1.40.32

## 0.6.0 (2021-08-08)

### Enhancements

- [#143](https://github.com/terraform-linters/tflint-ruleset-aws/pull/143): Add rules for aws_elasticache_replication_group resource
- [#151](https://github.com/terraform-linters/tflint-ruleset-aws/pull/151): Bump aws-sdk-go submodule and Terraform provider schema

### Chores

- [#138](https://github.com/terraform-linters/tflint-ruleset-aws/pull/138): Bump github.com/zclconf/go-cty from 1.8.4 to 1.9.0
- [#142](https://github.com/terraform-linters/tflint-ruleset-aws/pull/142): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.9.0 to 0.9.1
- [#144](https://github.com/terraform-linters/tflint-ruleset-aws/pull/144): Remove hashicorp/terraform-provider-aws build dependency
- [#145](https://github.com/terraform-linters/tflint-ruleset-aws/pull/145): Bump github.com/hashicorp/hcl/v2 from 2.10.0 to 2.10.1
- [#150](https://github.com/terraform-linters/tflint-ruleset-aws/pull/150): Bump github.com/aws/aws-sdk-go from 1.39.0 to 1.40.17

## 0.5.0 (2021-07-03)

The minimum supported version of TFLint has changed in this version. TFLint v0.30.0+ is required for this plugin to work.

### Breaking Changes

- [#137](https://github.com/terraform-linters/tflint-ruleset-aws/pull/137): Bump tflint-plugin-sdk to v0.9.0

### Enhancements

- [#136](https://github.com/terraform-linters/tflint-ruleset-aws/pull/136): Bump aws-sdk-go submodule and Terraform provider

### Chores

- [#129](https://github.com/terraform-linters/tflint-ruleset-aws/pull/129): Bump github.com/golang/mock from 1.5.0 to 1.6.0
- [#132](https://github.com/terraform-linters/tflint-ruleset-aws/pull/132): Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.6.1 to 2.7.0
- [#134](https://github.com/terraform-linters/tflint-ruleset-aws/pull/134): Bump github.com/zclconf/go-cty from 1.8.3 to 1.8.4
- [#135](https://github.com/terraform-linters/tflint-ruleset-aws/pull/135): Bump github.com/aws/aws-sdk-go from 1.38.55 to 1.39.0

## 0.4.3 (2021-06-12)

### Chores

- [#127](https://github.com/terraform-linters/tflint-ruleset-aws/pull/127): Fix GoReleaser action inputs

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

