## What's Changed

Support for Cosign signatures has been removed from this release. The `checksums.txt.keyless.sig` and `checksums.txt.pem` will not be included in the release.
These files are not used in normal use cases, so in most cases this will not affect you, but if you are affected, you can use Artifact Attestations instead.

### Breaking Changes
* Bump github.com/terraform-linters/tflint-plugin-sdk from 0.22.0 to 0.23.1 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/966
  * Requires TFLint v0.46+

### Enhancements
* Update AWS provider/module and generated content by @github-actions[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/954
* Update AWS provider/module and generated content by @github-actions[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/959
* Update Lambda runtime deprecation dates by @Copilot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/969
* Add missing ElastiCache node type: cache.r6gd.large by @Copilot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/971
* Fix typos in AWS RDS DB instance types by @Copilot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/972
* Add missing AWS S3 bucket naming restrictions by @Copilot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/976

### Chores
* Bump github.com/aws/aws-sdk-go-v2/service/ec2 from 1.251.2 to 1.253.0 in the aws-sdk group by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/952
* Bump github.com/hashicorp/terraform-json from 0.26.0 to 0.27.2 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/953
* Bump the aws-sdk group with 7 updates by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/956
* Bump the aws-sdk group with 2 updates by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/957
* Bump github.com/hashicorp/aws-sdk-go-base/v2 from 2.0.0-beta.66 to 2.0.0-beta.67 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/958
* Bump github.com/aws/aws-sdk-go-v2/service/rds from 1.108.0 to 1.108.2 in the aws-sdk group by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/960
* Bump golang.org/x/net from 0.44.0 to 0.46.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/961
* Bump sigstore/cosign-installer from 3.10.0 to 4.0.0 by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/962
* Bump the aws-sdk group with 7 updates by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/963
* Bump the aws-sdk group with 7 updates by @dependabot[bot] in https://github.com/terraform-linters/tflint-ruleset-aws/pull/965
* Drop support for Cosign signatures by @wata727 in https://github.com/terraform-linters/tflint-ruleset-aws/pull/968
* Add documentation to AWS MQ engine type validation rules by @Copilot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/974
* Reorder S3 bucket ACL enum values for consistency by @Copilot in https://github.com/terraform-linters/tflint-ruleset-aws/pull/975

## New Contributors
* @Copilot made their first contribution in https://github.com/terraform-linters/tflint-ruleset-aws/pull/969

**Full Changelog**: https://github.com/terraform-linters/tflint-ruleset-aws/compare/v0.43.0...v0.44.0
