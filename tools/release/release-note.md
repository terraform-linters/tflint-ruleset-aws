## What's Changed

### Breaking Changes
* Update AWS provider/module and generated content by @wata727 in https://github.com/terraform-linters/tflint-ruleset-aws/pull/837
  * Remove Amazon Chime rules
    * `aws_chime_voice_connector_group_invalid_name`
    * `aws_chime_voice_connector_invalid_aws_region`
    * `aws_chime_voice_connector_invalid_name`
    * `aws_chime_voice_connector_logging_invalid_voice_connector_id`
    * `aws_chime_voice_connector_origination_invalid_voice_connector_id`
    * `aws_chime_voice_connector_streaming_invalid_voice_connector_id`
    * `aws_chime_voice_connector_termination_credentials_invalid_voice_connector_id`
    * `aws_chime_voice_connector_termination_invalid_default_phone_number`
    * `aws_chime_voice_connector_termination_invalid_voice_connector_id`

### Enhancements
* feat: add aws_security_group_inline_rules rule by @kayman-mk in https://github.com/terraform-linters/tflint-ruleset-aws/pull/793

### Chores
* Bump github.com/aws/aws-sdk-go-v2/service/ecs from 1.53.1 to 1.53.2 in the aws-sdk group by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/806
* Bump github.com/zclconf/go-cty from 1.15.1 to 1.16.0 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/807
* Bump the aws-sdk group with 7 updates by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/809
* Bump golang.org/x/net from 0.33.0 to 0.34.0 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/811
* Bump github.com/terraform-linters/tflint-plugin-sdk from 0.21.0 to 0.22.0 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/812
* Bump github.com/hashicorp/aws-sdk-go-base/v2 from 2.0.0-beta.59 to 2.0.0-beta.60 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/810
* Bump github.com/hashicorp/aws-sdk-go-base/v2 from 2.0.0-beta.60 to 2.0.0-beta.61 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/816
* Bump github.com/zclconf/go-cty from 1.16.0 to 1.16.1 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/817
* Bump the aws-sdk group with 7 updates by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/815
* Bump github.com/zclconf/go-cty from 1.16.1 to 1.16.2 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/819
* Bump the aws-sdk group with 7 updates by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/818
* Bump the aws-sdk group with 7 updates by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/823
* Bump the aws-sdk group with 7 updates by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/824
* Bump github.com/hashicorp/aws-sdk-go-base/v2 from 2.0.0-beta.61 to 2.0.0-beta.62 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/825
* Bump the aws-sdk group with 2 updates by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/826
* Bump golang.org/x/net from 0.34.0 to 0.35.0 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/827
* Bump github.com/aws/smithy-go from 1.22.2 to 1.22.3 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/830
* Bump github.com/google/go-cmp from 0.6.0 to 0.7.0 by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/831
* Bump the aws-sdk group with 7 updates by @dependabot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/829
* deps: Go 1.24 by @wata727 in https://github.com/terraform-linters/tflint-ruleset-aws/pull/832
* rule template: fix typo by @bendrucker in https://github.com/terraform-linters/tflint-ruleset-aws/pull/834
* Remove hard-coded versions from integration tests by @wata727 in https://github.com/terraform-linters/tflint-ruleset-aws/pull/835
* Add make release for release automation by @wata727 in https://github.com/terraform-linters/tflint-ruleset-aws/pull/836


**Full Changelog**: https://github.com/terraform-linters/tflint-ruleset-aws/compare/v0.37.0...v0.38.0
