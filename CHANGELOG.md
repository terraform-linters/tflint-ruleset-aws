See https://github.com/terraform-linters/tflint-ruleset-aws/releases for later releases.

## 0.37.0 (2024-12-31)

### Breaking Changes

- [#789](https://github.com/terraform-linters/tflint-ruleset-aws/pull/789) [#798](https://github.com/terraform-linters/tflint-ruleset-aws/pull/798): Update AWS provider/module and generated content
  - Removed the following rules
    - `aws_service_discovery_http_namespace_invalid_name`
    - `aws_service_discovery_private_dns_namespace_invalid_name`
    - `aws_service_discovery_public_dns_namespace_invalid_name`

### Enhancements

- [#787](https://github.com/terraform-linters/tflint-ruleset-aws/pull/787) [#800](https://github.com/terraform-linters/tflint-ruleset-aws/pull/800): Add new db.m8g, r8g, r7i, r6i-preconfigured, m7i, c6gd instances ([@gnetsman](https://github.com/gnetsman), [@wata727](https://github.com/wata727))
- [#786](https://github.com/terraform-linters/tflint-ruleset-aws/pull/786): feat: add `aws_iam_policy_attachment_exclusive_attachment` rule ([@kayman-mk](https://github.com/kayman-mk))
- [#790](https://github.com/terraform-linters/tflint-ruleset-aws/pull/790): feat: add `aws_security_group_rule_deprecated` rule ([@kayman-mk](https://github.com/kayman-mk))
- [#801](https://github.com/terraform-linters/tflint-ruleset-aws/pull/801): rules: Add missing DB engines ([@wata727](https://github.com/wata727))
- [#802](https://github.com/terraform-linters/tflint-ruleset-aws/pull/802): rules: Add cache.c7gn high-bandwidth node type ([@wata727](https://github.com/wata727))
- [#803](https://github.com/terraform-linters/tflint-ruleset-aws/pull/803): rules: Update Lambda runtime deprecations ([@wata727](https://github.com/wata727))

### Chores

- [#791](https://github.com/terraform-linters/tflint-ruleset-aws/pull/791) [#795](https://github.com/terraform-linters/tflint-ruleset-aws/pull/795) [#797](https://github.com/terraform-linters/tflint-ruleset-aws/pull/797): Bump aws-sdk-go-v2
- [#792](https://github.com/terraform-linters/tflint-ruleset-aws/pull/792): Bump github.com/hashicorp/terraform-json from 0.23.0 to 0.24.0
- [#796](https://github.com/terraform-linters/tflint-ruleset-aws/pull/796): Bump golang.org/x/net from 0.32.0 to 0.33.0

## 0.36.0 (2024-12-08)

### Breaking Changes

- [#759](https://github.com/terraform-linters/tflint-ruleset-aws/pull/759): glue: disable connection type and properties checks ([@bendrucker](https://github.com/bendrucker))
  - Removed `aws_glue_connection_invalid_connection_type` rule

### Enhancements

- [#758](https://github.com/terraform-linters/tflint-ruleset-aws/pull/758) [#768](https://github.com/terraform-linters/tflint-ruleset-aws/pull/768) [#779](https://github.com/terraform-linters/tflint-ruleset-aws/pull/779) [#780](https://github.com/terraform-linters/tflint-ruleset-aws/pull/780): Update AWS provider/module and generated content

### Chores

- [#755](https://github.com/terraform-linters/tflint-ruleset-aws/pull/755): Bump github.com/hashicorp/hcl/v2 from 2.22.0 to 2.23.0
- [#765](https://github.com/terraform-linters/tflint-ruleset-aws/pull/765): Bump github.com/stretchr/testify from 1.9.0 to 1.10.0
- [#769](https://github.com/terraform-linters/tflint-ruleset-aws/pull/769): Bump github.com/aws/aws-sdk-go-v2/service/ec2 from 1.187.1 to 1.194.0
- [#777](https://github.com/terraform-linters/tflint-ruleset-aws/pull/777): Bump the aws-sdk group with 7 updates
- [#778](https://github.com/terraform-linters/tflint-ruleset-aws/pull/778): Bump github.com/zclconf/go-cty from 1.15.0 to 1.15.1
- [#781](https://github.com/terraform-linters/tflint-ruleset-aws/pull/781): Enable Dependabot auto-merge ([@wata727](https://github.com/wata727))
- [#782](https://github.com/terraform-linters/tflint-ruleset-aws/pull/782): Bump golang.org/x/net from 0.31.0 to 0.32.0
- [#783](https://github.com/terraform-linters/tflint-ruleset-aws/pull/783): Squash auto-merge Dependabot PRs ([@wata727](https://github.com/wata727))
- [#784](https://github.com/terraform-linters/tflint-ruleset-aws/pull/784): Bump actions/attest-build-provenance from 1 to 2

## 0.35.0 (2024-11-17)

### Enhancements

- [#725](https://github.com/terraform-linters/tflint-ruleset-aws/pull/725): Update AWS provider/module and generated content

### Chores

- [#726](https://github.com/terraform-linters/tflint-ruleset-aws/pull/726) [#731](https://github.com/terraform-linters/tflint-ruleset-aws/pull/731) [#739](https://github.com/terraform-linters/tflint-ruleset-aws/pull/739) [#746](https://github.com/terraform-linters/tflint-ruleset-aws/pull/746): Bump github.com/aws/aws-sdk-go-v2/service/ecs from 1.47.3 to 1.49.2
- [#727](https://github.com/terraform-linters/tflint-ruleset-aws/pull/727) [#732](https://github.com/terraform-linters/tflint-ruleset-aws/pull/732) [#747](https://github.com/terraform-linters/tflint-ruleset-aws/pull/747): Bump github.com/aws/aws-sdk-go-v2/service/ec2 from 1.182.0 to 1.187.1
- [#728](https://github.com/terraform-linters/tflint-ruleset-aws/pull/728) [#730](https://github.com/terraform-linters/tflint-ruleset-aws/pull/730) [#736](https://github.com/terraform-linters/tflint-ruleset-aws/pull/736) [#748](https://github.com/terraform-linters/tflint-ruleset-aws/pull/748): Bump github.com/aws/aws-sdk-go-v2/service/rds from 1.87.2 to 1.89.1
- [#733](https://github.com/terraform-linters/tflint-ruleset-aws/pull/733) [#751](https://github.com/terraform-linters/tflint-ruleset-aws/pull/751): Bump github.com/aws/aws-sdk-go-v2/service/iam from 1.37.2 to 1.37.3
- [#735](https://github.com/terraform-linters/tflint-ruleset-aws/pull/735): Bump github.com/hashicorp/terraform-json from 0.22.1 to 0.23.0
- [#738](https://github.com/terraform-linters/tflint-ruleset-aws/pull/738) [#743](https://github.com/terraform-linters/tflint-ruleset-aws/pull/743): Bump github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing from 1.28.2 to 1.28.4
- [#740](https://github.com/terraform-linters/tflint-ruleset-aws/pull/740): Bump github.com/hashicorp/aws-sdk-go-base/v2 from 2.0.0-beta.58 to 2.0.0-beta.59
- [#741](https://github.com/terraform-linters/tflint-ruleset-aws/pull/741): Bump github.com/aws/aws-sdk-go-v2/service/elasticache from 1.43.0 to 1.43.2
- [#742](https://github.com/terraform-linters/tflint-ruleset-aws/pull/742): Bump github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 from 1.40.0 to 1.41.1
- [#744](https://github.com/terraform-linters/tflint-ruleset-aws/pull/744): Bump golang.org/x/net from 0.30.0 to 0.31.0

## 0.34.0 (2024-10-15)

### Breaking Changes

- [#719](https://github.com/terraform-linters/tflint-ruleset-aws/pull/719): Remove worklink rules ([@wata727](https://github.com/wata727))
  - The following rules are removed
    - `aws_worklink_fleet_invalid_audit_stream_arn`
    - `aws_worklink_fleet_invalid_device_ca_certificate`
    - `aws_worklink_fleet_invalid_display_name`
    - `aws_worklink_fleet_invalid_name`
    - `aws_worklink_website_certificate_authority_association_invalid_certificate`
    - `aws_worklink_website_certificate_authority_association_invalid_display_name`
    - `aws_worklink_website_certificate_authority_association_invalid_fleet_arn`

### Enhancements

- [#720](https://github.com/terraform-linters/tflint-ruleset-aws/pull/720): Update AWS provider/module and generated content

### Bug Fixes

- [#721](https://github.com/terraform-linters/tflint-ruleset-aws/pull/721): Fix incorrect mapping about aws_servicecatalog_portfolio_share type ([@wata727](https://github.com/wata727))

### Chores

- [#692](https://github.com/terraform-linters/tflint-ruleset-aws/pull/692) [#712](https://github.com/terraform-linters/tflint-ruleset-aws/pull/712): Bump github.com/aws/aws-sdk-go-v2/service/ecs from 1.44.3 to 1.47.3
- [#693](https://github.com/terraform-linters/tflint-ruleset-aws/pull/693): Bump github.com/hashicorp/aws-sdk-go-base/v2 from 2.0.0-beta.54 to 2.0.0-beta.58
- [#694](https://github.com/terraform-linters/tflint-ruleset-aws/pull/694) [#717](https://github.com/terraform-linters/tflint-ruleset-aws/pull/717): Bump github.com/aws/aws-sdk-go-v2/service/iam from 1.34.3 to 1.37.2
- [#695](https://github.com/terraform-linters/tflint-ruleset-aws/pull/695) [#714](https://github.com/terraform-linters/tflint-ruleset-aws/pull/714): Bump github.com/aws/aws-sdk-go-v2/service/rds from 1.82.0 to 1.87.2
- [#696](https://github.com/terraform-linters/tflint-ruleset-aws/pull/696) [#715](https://github.com/terraform-linters/tflint-ruleset-aws/pull/715): Bump github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing from 1.26.3 to 1.28.2
- [#697](https://github.com/terraform-linters/tflint-ruleset-aws/pull/697) [#713](https://github.com/terraform-linters/tflint-ruleset-aws/pull/713): Bump github.com/aws/aws-sdk-go-v2/service/ec2 from 1.173.0 to 1.182.0
- [#698](https://github.com/terraform-linters/tflint-ruleset-aws/pull/698) [#718](https://github.com/terraform-linters/tflint-ruleset-aws/pull/718): Bump github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 from 1.34.0 to 1.40.0
- [#699](https://github.com/terraform-linters/tflint-ruleset-aws/pull/699) [#716](https://github.com/terraform-linters/tflint-ruleset-aws/pull/716): Bump github.com/aws/aws-sdk-go-v2/service/elasticache from 1.40.5 to 1.43.0
- [#702](https://github.com/terraform-linters/tflint-ruleset-aws/pull/702): Bump golang.org/x/net from 0.29.0 to 0.30.0
- [#722](https://github.com/terraform-linters/tflint-ruleset-aws/pull/722): Bump GoReleaser to v2 ([@wata727](https://github.com/wata727))
- [#723](https://github.com/terraform-linters/tflint-ruleset-aws/pull/723): Remove patch version pin from go.mod ([@wata727](https://github.com/wata727))

## 0.33.0 (2024-09-10)

### Enhancements

- [#650](https://github.com/terraform-linters/tflint-ruleset-aws/pull/650) [#651](https://github.com/terraform-linters/tflint-ruleset-aws/pull/651) [#653](https://github.com/terraform-linters/tflint-ruleset-aws/pull/653) [#655](https://github.com/terraform-linters/tflint-ruleset-aws/pull/655) [#656](https://github.com/terraform-linters/tflint-ruleset-aws/pull/656) [#658](https://github.com/terraform-linters/tflint-ruleset-aws/pull/658) [#659](https://github.com/terraform-linters/tflint-ruleset-aws/pull/659) [#660](https://github.com/terraform-linters/tflint-ruleset-aws/pull/660): Update AWS provider/module and generated content

### Chores

- [#652](https://github.com/terraform-linters/tflint-ruleset-aws/pull/652): Fix HSMv2 type tests ([@wata727](https://github.com/wata727))
- [#654](https://github.com/terraform-linters/tflint-ruleset-aws/pull/654) [#668](https://github.com/terraform-linters/tflint-ruleset-aws/pull/668): Bump github.com/hashicorp/hcl/v2 from 2.20.1 to 2.22.0
- [#657](https://github.com/terraform-linters/tflint-ruleset-aws/pull/657) [#670](https://github.com/terraform-linters/tflint-ruleset-aws/pull/670): Bump golang.org/x/net from 0.26.0 to 0.29.0
- [#665](https://github.com/terraform-linters/tflint-ruleset-aws/pull/665): mappings: Migrate import paths from aws-sdk-go to aws-sdk-ruby ([@wata727](https://github.com/wata727))
- [#663](https://github.com/terraform-linters/tflint-ruleset-aws/pull/663): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.20.0 to 0.21.0
- [#667](https://github.com/terraform-linters/tflint-ruleset-aws/pull/667): Migrate aws-sdk-go to v2 ([@wata727](https://github.com/wata727))
- [#671](https://github.com/terraform-linters/tflint-ruleset-aws/pull/671): Bump peter-evans/create-pull-request from 6 to 7
- [#672](https://github.com/terraform-linters/tflint-ruleset-aws/pull/672): deps: Go 1.23.1

## 0.32.0 (2024-06-08)

### Enhancements

- [#633](https://github.com/terraform-linters/tflint-ruleset-aws/pull/633): New rule: aws_provider_missing_tags ([@bootswithdefer](https://github.com/bootswithdefer))
- [#637](https://github.com/terraform-linters/tflint-ruleset-aws/pull/637): Update AWS provider/module and generated content

### Chores

- [#638](https://github.com/terraform-linters/tflint-ruleset-aws/pull/638): release: Introduce Artifact Attestations ([@wata727](https://github.com/wata727))
- [#642](https://github.com/terraform-linters/tflint-ruleset-aws/pull/642) [#645](https://github.com/terraform-linters/tflint-ruleset-aws/pull/645): Bump golang.org/x/net from 0.24.0 to 0.26.0
- [#643](https://github.com/terraform-linters/tflint-ruleset-aws/pull/643): Bump github.com/hashicorp/terraform-json from 0.21.0 to 0.22.1
- [#646](https://github.com/terraform-linters/tflint-ruleset-aws/pull/646): Bump goreleaser/goreleaser-action from 5 to 6
- [#647](https://github.com/terraform-linters/tflint-ruleset-aws/pull/647): deps: Go 1.22.4 ([@wata727](https://github.com/wata727))

## 0.31.0 (2024-05-03)

### Enhancements

- [#613](https://github.com/terraform-linters/tflint-ruleset-aws/pull/613) [#618](https://github.com/terraform-linters/tflint-ruleset-aws/pull/618) [#619](https://github.com/terraform-linters/tflint-ruleset-aws/pull/619) [#623](https://github.com/terraform-linters/tflint-ruleset-aws/pull/623) [#627](https://github.com/terraform-linters/tflint-ruleset-aws/pull/627) [#629](https://github.com/terraform-linters/tflint-ruleset-aws/pull/629): Update AWS provider/module and generated content
- [#617](https://github.com/terraform-linters/tflint-ruleset-aws/pull/617) [#624](https://github.com/terraform-linters/tflint-ruleset-aws/pull/624): Bump github.com/hashicorp/hcl/v2 from 2.19.1 to 2.20.1
  - This is required for provider-defined functions support

### Chores

- [#609](https://github.com/terraform-linters/tflint-ruleset-aws/pull/609) [#620](https://github.com/terraform-linters/tflint-ruleset-aws/pull/620) [#628](https://github.com/terraform-linters/tflint-ruleset-aws/pull/628): Bump golang.org/x/net from 0.20.0 to 0.24.0
- [#614](https://github.com/terraform-linters/tflint-ruleset-aws/pull/614) [#634](https://github.com/terraform-linters/tflint-ruleset-aws/pull/634): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.18.0 to 0.20.0
- [#615](https://github.com/terraform-linters/tflint-ruleset-aws/pull/615): Bump github.com/stretchr/testify from 1.8.4 to 1.9.0
- [#616](https://github.com/terraform-linters/tflint-ruleset-aws/pull/616) [#622](https://github.com/terraform-linters/tflint-ruleset-aws/pull/622): Bump github.com/zclconf/go-cty from 1.14.2 to 1.14.4
- [#621](https://github.com/terraform-linters/tflint-ruleset-aws/pull/621): Bump google.golang.org/protobuf from 1.32.0 to 1.33.0
- [#635](https://github.com/terraform-linters/tflint-ruleset-aws/pull/635): deps: Go 1.22.2

## 0.30.0 (2024-02-11)

### Enhancements

- [#589](https://github.com/terraform-linters/tflint-ruleset-aws/pull/598) [#601](https://github.com/terraform-linters/tflint-ruleset-aws/pull/601) [#602](https://github.com/terraform-linters/tflint-ruleset-aws/pull/602) [#606](https://github.com/terraform-linters/tflint-ruleset-aws/pull/606): Update AWS provider/module and generated content

### Bug Fixes

- [#600](https://github.com/terraform-linters/tflint-ruleset-aws/pull/600): fix: prevent runtime error when instanceType/nodeType have no dot ([@guillaumelecerf](https://github.com/guillaumelecerf))

### Chores

- [#599](https://github.com/terraform-linters/tflint-ruleset-aws/pull/599): Bump golang.org/x/net from 0.19.0 to 0.20.0
- [#603](https://github.com/terraform-linters/tflint-ruleset-aws/pull/603): Bump github.com/hashicorp/terraform-json from 0.20.0 to 0.21.0
- [#604](https://github.com/terraform-linters/tflint-ruleset-aws/pull/604): Bump github.com/zclconf/go-cty from 1.14.1 to 1.14.2
- [#605](https://github.com/terraform-linters/tflint-ruleset-aws/pull/605): Bump peter-evans/create-pull-request from 5 to 6
- [#607](https://github.com/terraform-linters/tflint-ruleset-aws/pull/607): deps: Go 1.22 ([@wata727](https://github.com/wata727))

## 0.29.0 (2024-01-08)

### Enhancements

- [#580](https://github.com/terraform-linters/tflint-ruleset-aws/pull/580) [#581](https://github.com/terraform-linters/tflint-ruleset-aws/pull/581) [#582](https://github.com/terraform-linters/tflint-ruleset-aws/pull/582) [#583](https://github.com/terraform-linters/tflint-ruleset-aws/pull/583) [#585](https://github.com/terraform-linters/tflint-ruleset-aws/pull/585) [#586](https://github.com/terraform-linters/tflint-ruleset-aws/pull/586) [#587](https://github.com/terraform-linters/tflint-ruleset-aws/pull/587): Update AWS provider/module and generated content
- [#591](https://github.com/terraform-linters/tflint-ruleset-aws/pull/591): rds: add db.m6idn instance family ([@GeorgeSimpsonAtFanatics](https://github.com/GeorgeSimpsonAtFanatics))
- [#592](https://github.com/terraform-linters/tflint-ruleset-aws/pull/592): rds: Add m6in, r6idn, r6in instance classes ([@wata727](https://github.com/wata727))
- [#593](https://github.com/terraform-linters/tflint-ruleset-aws/pull/593): cache: Add cache.c7gn node types ([@wata727](https://github.com/wata727))
- [#594](https://github.com/terraform-linters/tflint-ruleset-aws/pull/594): Update Lambda deprecated runtimes ([@wata727](https://github.com/wata727))

### Chore

- [#584](https://github.com/terraform-linters/tflint-ruleset-aws/pull/584): Bump golang.org/x/net from 0.18.0 to 0.19.0
- [#588](https://github.com/terraform-linters/tflint-ruleset-aws/pull/588): Bump actions/setup-go from 4 to 5
- [#589](https://github.com/terraform-linters/tflint-ruleset-aws/pull/589): generator: use terraform-json types for provider schemas ([@HadrienPatte](https://github.com/HadrienPatte))

## 0.28.0 (2023-11-18)

### Breaking Changes

- [#576](https://github.com/terraform-linters/tflint-ruleset-aws/pull/576): Fix aws_acm_certificate.private_key mapping and remove the rule ([@wata727](https://github.com/wata727))
  - Removed `aws_acm_certificate_invalid_private_key` rule

### Enhancements

- [#554](https://github.com/terraform-linters/tflint-ruleset-aws/pull/554): s3_bucket_name: add length validation ([@davimmt](https://github.com/davimmt))
- [#566](https://github.com/terraform-linters/tflint-ruleset-aws/pull/566) [#567](https://github.com/terraform-linters/tflint-ruleset-aws/pull/567) [#568](https://github.com/terraform-linters/tflint-ruleset-aws/pull/568) [#575](https://github.com/terraform-linters/tflint-ruleset-aws/pull/575): Update AWS provider/module and generated content
- [#571](https://github.com/terraform-linters/tflint-ruleset-aws/pull/571): s3_bucket_name: add all documented naming rules ([@davimmt](https://github.com/davimmt))

### Chores

- [#550](https://github.com/terraform-linters/tflint-ruleset-aws/pull/550): Bump goreleaser/goreleaser-action from 4 to 5
- [#555](https://github.com/terraform-linters/tflint-ruleset-aws/pull/555) [#558](https://github.com/terraform-linters/tflint-ruleset-aws/pull/558) [#572](https://github.com/terraform-linters/tflint-ruleset-aws/pull/572): Bump golang.org/x/net from 0.15.0 to 0.18.0
- [#556](https://github.com/terraform-linters/tflint-ruleset-aws/pull/556): Bump github.com/zclconf/go-cty from 1.14.0 to 1.14.1
- [#559](https://github.com/terraform-linters/tflint-ruleset-aws/pull/559): Bump github.com/dave/dst from 0.27.2 to 0.27.3
- [#560](https://github.com/terraform-linters/tflint-ruleset-aws/pull/560): Bump github.com/google/go-cmp from 0.5.9 to 0.6.0
- [#561](https://github.com/terraform-linters/tflint-ruleset-aws/pull/561): Bump github.com/hashicorp/hcl/v2 from 2.18.0 to 2.19.1
- [#565](https://github.com/terraform-linters/tflint-ruleset-aws/pull/565): Bump google.golang.org/grpc from 1.57.0 to 1.57.1

## 0.27.0 (2023-09-17)

### Enhancements

- [#534](https://github.com/terraform-linters/tflint-ruleset-aws/pull/534) [#544](https://github.com/terraform-linters/tflint-ruleset-aws/pull/544): Update AWS provider/module and generated content
- [#545](https://github.com/terraform-linters/tflint-ruleset-aws/pull/545): Add support for cache instance families m7g and r7g ([@dvxam](https://github.com/dvxam))
- [#546](https://github.com/terraform-linters/tflint-ruleset-aws/pull/546): Add missing DB instance classes ([@wata727](https://github.com/wata727))
- [#547](https://github.com/terraform-linters/tflint-ruleset-aws/pull/547): Update Lambda runtime deprecations ([@wata727](https://github.com/wata727))
- [#548](https://github.com/terraform-linters/tflint-ruleset-aws/pull/548): Update DB engine list ([@wata727](https://github.com/wata727))

### Chores

- [#535](https://github.com/terraform-linters/tflint-ruleset-aws/pull/535) [#536](https://github.com/terraform-linters/tflint-ruleset-aws/pull/536): Bump github.com/zclconf/go-cty from 1.13.2 to 1.14.0
- [#537](https://github.com/terraform-linters/tflint-ruleset-aws/pull/537): Bump github.com/hashicorp/hcl/v2 from 2.17.0 to 2.18.0
- [#540](https://github.com/terraform-linters/tflint-ruleset-aws/pull/540): deps: Go 1.21 ([@wata727](https://github.com/wata727))
- [#541](https://github.com/terraform-linters/tflint-ruleset-aws/pull/541): Bump actions/checkout from 3 to 4
- [#542](https://github.com/terraform-linters/tflint-ruleset-aws/pull/542): Bump golang.org/x/net from 0.14.0 to 0.15.0

## 0.26.0 (2023-08-14)

### Breaking Changes

- [#531](https://github.com/terraform-linters/tflint-ruleset-aws/pull/531): Update AWS provider/module and generated content
  - Removed `aws_budgets_budget_invalid_name` rule

### BugFixes

- [#530](https://github.com/terraform-linters/tflint-ruleset-aws/pull/530): resource_missing_tags: Fix panic for unknown and null values ([@wata727](https://github.com/wata727))

### Chores

- [#529](https://github.com/terraform-linters/tflint-ruleset-aws/pull/529): Bump golang.org/x/net from 0.12.0 to 0.14.0

## 0.25.0 (2023-07-29)

### Enhancements

- [#514](https://github.com/terraform-linters/tflint-ruleset-aws/pull/514): Update AWS provider/module and generated content

### BugFixes

- [#526](https://github.com/terraform-linters/tflint-ruleset-aws/pull/526): Bump tflint-plugin-sdk to v0.18.0 ([@wata727](https://github.com/wata727))

### Chores

- [#523](https://github.com/terraform-linters/tflint-ruleset-aws/pull/523): release: Add raw binary entries to checksums.txt ([@wata727](https://github.com/wata727))

## 0.24.3 (2023-07-22)

- [#517](https://github.com/terraform-linters/tflint-ruleset-aws/pull/517): resource_missing_tags: rewrite to only consider keys ([@bendrucker](https://github.com/bendrucker))
- [#519](https://github.com/terraform-linters/tflint-ruleset-aws/pull/519): resource_missing_tags: Ignore non-existent provider aliases ([@wata727](https://github.com/wata727))

## 0.24.2 (2023-07-19)

### Enhancements

- [#505](https://github.com/terraform-linters/tflint-ruleset-aws/pull/505) [#513](https://github.com/terraform-linters/tflint-ruleset-aws/pull/513): Update AWS provider/module and generated content

### BugFixes

- [#511](https://github.com/terraform-linters/tflint-ruleset-aws/pull/511): fix: Fixed regression in tags with dynamic values ([@JorgeReus](https://github.com/JorgeReus))

### Chores

- [#512](https://github.com/terraform-linters/tflint-ruleset-aws/pull/512): Bump golang.org/x/net from 0.11.0 to 0.12.0

## 0.24.1 (2023-07-01)

- [#508](https://github.com/terraform-linters/tflint-ruleset-aws/pull/508): fix: Fixed the error 'Provider doesn't exist' and the provider exist ([@JorgeReus](https://github.com/JorgeReus))

## 0.24.0 (2023-06-29)

### Breaking Changes

- [#501](https://github.com/terraform-linters/tflint-ruleset-aws/pull/501): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.16.1 to 0.17.0
  - TFLint v0.40/v0.41 is no longer supported
- [#504](https://github.com/terraform-linters/tflint-ruleset-aws/pull/504): Update AWS provider/module and generated content  ([@wata727](https://github.com/wata727))
  - The following rules have been removed.
    - `aws_macie_member_account_association_invalid_member_account_id`
    - `aws_macie_s3_bucket_association_invalid_bucket_name`
    - `aws_macie_s3_bucket_association_invalid_member_account_id`
    - `aws_macie_s3_bucket_association_invalid_prefix`
    - `aws_redshift_security_group_invalid_description`
    - `aws_redshift_security_group_invalid_name`
    - `aws_secretsmanager_secret_invalid_rotation_lambda_arn`
  - The `aws_resource_missing_tags` rule no longer checks `aws_db_security_group`.

### Enhancements

- [#489](https://github.com/terraform-linters/tflint-ruleset-aws/pull/489): feat: Added default tags functionality ([@JorgeReus](https://github.com/JorgeReus))

### Chores

- [#497](https://github.com/terraform-linters/tflint-ruleset-aws/pull/497): Bump github.com/zclconf/go-cty from 1.13.1 to 1.13.2
- [#498](https://github.com/terraform-linters/tflint-ruleset-aws/pull/498): Bump github.com/hashicorp/hcl/v2 from 2.16.2 to 2.17.0
- [#500](https://github.com/terraform-linters/tflint-ruleset-aws/pull/500): Bump golang.org/x/net from 0.10.0 to 0.11.0
- [#502](https://github.com/terraform-linters/tflint-ruleset-aws/pull/502): Bump github.com/stretchr/testify from 1.7.2 to 1.8.4

## 0.23.1 (2023-05-22)

### Enhancements

- [#484](https://github.com/terraform-linters/tflint-ruleset-aws/pull/484): `aws_route_not_specified_target`: Add core_network_arn as target ([@ttretau](https://github.com/ttretau))
- [#485](https://github.com/terraform-linters/tflint-ruleset-aws/pull/485) [#487](https://github.com/terraform-linters/tflint-ruleset-aws/pull/487) [#490](https://github.com/terraform-linters/tflint-ruleset-aws/pull/490) [#495](https://github.com/terraform-linters/tflint-ruleset-aws/pull/495): Update AWS provider/module and generated content

### Chores

- [#493](https://github.com/terraform-linters/tflint-ruleset-aws/pull/493): Bump golang.org/x/net from 0.9.0 to 0.10.0

## 0.23.0 (2023-04-22)

### Enhancements

- [#471](https://github.com/terraform-linters/tflint-ruleset-aws/pull/471) [#480](https://github.com/terraform-linters/tflint-ruleset-aws/pull/480): Update AWS provider/module and generated content

### Chores

- [#436](https://github.com/terraform-linters/tflint-ruleset-aws/pull/436): Use NewRunner hook ([@wata727](https://github.com/wata727))
- [#468](https://github.com/terraform-linters/tflint-ruleset-aws/pull/468): Bump actions/setup-go from 3 to 4
- [#469](https://github.com/terraform-linters/tflint-ruleset-aws/pull/469): Bump github.com/zclconf/go-cty from 1.13.0 to 1.13.1
- [#473](https://github.com/terraform-linters/tflint-ruleset-aws/pull/473): Bump peter-evans/create-pull-request from 4 to 5
- [#475](https://github.com/terraform-linters/tflint-ruleset-aws/pull/475): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.15.0 to 0.16.1
- [#477](https://github.com/terraform-linters/tflint-ruleset-aws/pull/477): docs: copy edits for deep check ([@bendrucker](https://github.com/bendrucker))
- [#481](https://github.com/terraform-linters/tflint-ruleset-aws/pull/481): Follow up of the EnsureNoError deprecation ([@wata727](https://github.com/wata727))

## 0.22.1 (2023-03-18)

- [#465](https://github.com/terraform-linters/tflint-ruleset-aws/pull/465): Fix Cosign v2 signing ([@wata727](https://github.com/wata727))

## 0.22.0 (2023-03-18)

### Breaking Changes

- [#462](https://github.com/terraform-linters/tflint-ruleset-aws/pull/462): appsync: Remove invalid regexp rules ([@wata727](https://github.com/wata727))

### Enhancements

- [#444](https://github.com/terraform-linters/tflint-ruleset-aws/pull/444) [#451](https://github.com/terraform-linters/tflint-ruleset-aws/pull/451) [#454](https://github.com/terraform-linters/tflint-ruleset-aws/pull/454): Update AWS provider/module and generated content

### BugFixes

- [#455](https://github.com/terraform-linters/tflint-ruleset-aws/pull/455): aws_acm_certificate: fix false positive for private CA ([@bendrucker](https://github.com/bendrucker) [@wata727](https://github.com/wata727))

### Chores

- [#445](https://github.com/terraform-linters/tflint-ruleset-aws/pull/445) [#452](https://github.com/terraform-linters/tflint-ruleset-aws/pull/452) [#460](https://github.com/terraform-linters/tflint-ruleset-aws/pull/460): Bump github.com/hashicorp/hcl/v2 from 2.15.0 to 2.16.2
- [#447](https://github.com/terraform-linters/tflint-ruleset-aws/pull/447) [#449](https://github.com/terraform-linters/tflint-ruleset-aws/pull/449) [#459](https://github.com/terraform-linters/tflint-ruleset-aws/pull/459): Bump golang.org/x/net from 0.5.0 to 0.8.0
- [#450](https://github.com/terraform-linters/tflint-ruleset-aws/pull/450): Fix submodule checkout issue ([@wata727](https://github.com/wata727))
- [#457](https://github.com/terraform-linters/tflint-ruleset-aws/pull/457): Bump github.com/zclconf/go-cty from 1.12.1 to 1.13.0
- [#458](https://github.com/terraform-linters/tflint-ruleset-aws/pull/458): Bump sigstore/cosign-installer from 2 to 3
- [#463](https://github.com/terraform-linters/tflint-ruleset-aws/pull/463): Fix generated_code_checks workflow step ([@wata727](https://github.com/wata727))
- [#464](https://github.com/terraform-linters/tflint-ruleset-aws/pull/464): go 1.20 ([@wata727](https://github.com/wata727))

## 0.21.2 (2023-02-03)

### Enhancements

- [#431](https://github.com/terraform-linters/tflint-ruleset-aws/pull/431) [#442](https://github.com/terraform-linters/tflint-ruleset-aws/pull/442): Update AWS provider/module and generated content

### Chores

- [#433](https://github.com/terraform-linters/tflint-ruleset-aws/pull/433) [#441](https://github.com/terraform-linters/tflint-ruleset-aws/pull/441): Bump golang.org/x/net from 0.2.0 to 0.5.0
- [#434](https://github.com/terraform-linters/tflint-ruleset-aws/pull/434): Bump goreleaser/goreleaser-action from 3 to 4
- [#435](https://github.com/terraform-linters/tflint-ruleset-aws/pull/435): Pass `GITHUB_TOKEN` to e2e test workflow ([@wata727](https://github.com/wata727))
- [#437](https://github.com/terraform-linters/tflint-ruleset-aws/pull/437): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.14.0 to 0.15.0

## 0.21.1 (2022-12-12)

### BugFixes

- [#430](https://github.com/terraform-linters/tflint-ruleset-aws/pull/430): `elasticache_cluster_previous_type`: fix panic on empty string ([@bendrucker](https://github.com/bendrucker))

### Chores

- [#407](https://github.com/terraform-linters/tflint-ruleset-aws/pull/407): autogenerated maintenance

## 0.21.0 (2022-12-05)

### Enhancements

- [#403](https://github.com/terraform-linters/tflint-ruleset-aws/pull/403): autogenerated maintenance
- [#405](https://github.com/terraform-linters/tflint-ruleset-aws/pull/405) [#406](https://github.com/terraform-linters/tflint-ruleset-aws/pull/406): Add assume role configuration to plugin config ([@kaito3desuyo](https://github.com/kaito3desuyo))

## 0.20.0 (2022-11-27)

### Enhancements

- [#400](https://github.com/terraform-linters/tflint-ruleset-aws/pull/400): autogenerated maintenance

### Chores

- [#399](https://github.com/terraform-linters/tflint-ruleset-aws/pull/399): Bump up GoReleaser version in release.yml ([@wata727](https://github.com/wata727))
- [#401](https://github.com/terraform-linters/tflint-ruleset-aws/pull/401): Bump golang.org/x/net from 0.1.0 to 0.2.0

## 0.19.0 (2022-11-14)

### Enhancements

- [#390](https://github.com/terraform-linters/tflint-ruleset-aws/pull/390): autogenerated maintenance

### BugFixes

- [#397](https://github.com/terraform-linters/tflint-ruleset-aws/pull/397): Prefer credentials in "plugin" blocks over "provider" blocks ([@wata727](https://github.com/wata727))

### Chores

- [#394](https://github.com/terraform-linters/tflint-ruleset-aws/pull/394): Add signatures for keyless signing ([@wata727](https://github.com/wata727))
- [#395](https://github.com/terraform-linters/tflint-ruleset-aws/pull/395): Bump github.com/hashicorp/hcl/v2 from 2.14.1 to 2.15.0
- [#398](https://github.com/terraform-linters/tflint-ruleset-aws/pull/398): Bump up GoReleaser version ([@wata727](https://github.com/wata727))

## 0.18.0 (2022-10-24)

### Breaking Changes

- [#367](https://github.com/terraform-linters/tflint-ruleset-aws/pull/367): remove hardcoded S3 region rule ([@PatMyron](https://github.com/PatMyron))

### Enhancements

- [#382](https://github.com/terraform-linters/tflint-ruleset-aws/pull/382): autogenerated maintenance
- [#388](https://github.com/terraform-linters/tflint-ruleset-aws/pull/388): Bump tflint-plugin-sdk to v0.14.0 ([@wata727](https://github.com/wata727))

### Chores

- [#387](https://github.com/terraform-linters/tflint-ruleset-aws/pull/387): Bump github.com/dave/dst from 0.27.0 to 0.27.2

## 0.17.1 (2022-09-29)

### Enhancements

- [#373](https://github.com/terraform-linters/tflint-ruleset-aws/pull/373): autogenerated maintenance
- [#380](https://github.com/terraform-linters/tflint-ruleset-aws/pull/380): Update db instance type list with m6i and r6i ([@milestruecar](https://github.com/milestruecar))

### Chores

- [#374](https://github.com/terraform-linters/tflint-ruleset-aws/pull/374): Bump github.com/google/go-cmp from 0.5.8 to 0.5.9
- [#377](https://github.com/terraform-linters/tflint-ruleset-aws/pull/377): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.12.0 to 0.13.0
- [#378](https://github.com/terraform-linters/tflint-ruleset-aws/pull/378): Bump github.com/hashicorp/hcl/v2 from 2.14.0 to 2.14.1

## 0.17.0 (2022-09-08)

The minimum supported version of TFLint has changed in this version. TFLint v0.40.0+ is required for this plugin to work.

### Breaking Changes

- [#369](https://github.com/terraform-linters/tflint-ruleset-aws/pull/369): Bump tflint-plugin-sdk to v0.12.0 ([@wata727](https://github.com/wata727))

### Enhancements

- [#366](https://github.com/terraform-linters/tflint-ruleset-aws/pull/366): autogenerated maintenance

### Chores

- [#365](https://github.com/terraform-linters/tflint-ruleset-aws/pull/365): Bump github.com/zclconf/go-cty from 1.10.0 to 1.11.0
- [#368](https://github.com/terraform-linters/tflint-ruleset-aws/pull/368): Bump github.com/hashicorp/hcl/v2 from 2.13.0 to 2.14.0
- [#371](https://github.com/terraform-linters/tflint-ruleset-aws/pull/371): build: Improve Go workflows ([@wata727](https://github.com/wata727))

## 0.16.1 (2022-08-27)

### Enhancements

- [363](https://github.com/terraform-linters/tflint-ruleset-aws/pull/363): autogenerated maintenance

## 0.16.0 (2022-08-14)

### Enhancements

- [#358](https://github.com/terraform-linters/tflint-ruleset-aws/pull/358): autogenerated maintenance
  - Removed `aws_cloudwatch_metric_alarm_invalid_extended_statistic` rule
- [#362](https://github.com/terraform-linters/tflint-ruleset-aws/pull/362): Lambda runtime deprecation updates ([@PatMyron](https://github.com/PatMyron))

### Chores

- [#359](https://github.com/terraform-linters/tflint-ruleset-aws/pull/359): go 1.19 ([@PatMyron](https://github.com/PatMyron))

## 0.15.0 (2022-07-15)

### Enhancements

- [#352](https://github.com/terraform-linters/tflint-ruleset-aws/pull/352): autogenerated maintenance
- [#355](https://github.com/terraform-linters/tflint-ruleset-aws/pull/355): Add `aws_security_group_rule_invalid_protocol` rule ([@x-color](https://github.com/x-color))
- [#356](https://github.com/terraform-linters/tflint-ruleset-aws/pull/356): Add `aws_security_group_invalid_protocol` rule ([@x-color](https://github.com/x-color))

### Chores

- [#354](https://github.com/terraform-linters/tflint-ruleset-aws/pull/354): Bump github.com/hashicorp/hcl/v2 from 2.12.0 to 2.13.0

## 0.14.0 (2022-05-31)

### Enhancements

- [#342](https://github.com/terraform-linters/tflint-ruleset-aws/pull/342): feat: support provider aliases in deep checking ([@suzuki-shunsuke](https://github.com/suzuki-shunsuke))
- [#343](https://github.com/terraform-linters/tflint-ruleset-aws/pull/343): autogenerated maintenance

### Chores

- [#344](https://github.com/terraform-linters/tflint-ruleset-aws/pull/344): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.10.1 to 0.11.0
- [#347](https://github.com/terraform-linters/tflint-ruleset-aws/pull/347): Bump goreleaser/goreleaser-action from 2 to 3
- [#351](https://github.com/terraform-linters/tflint-ruleset-aws/pull/351): Bump github.com/dave/dst from 0.26.2 to 0.27.0

## 0.13.4 (2022-05-05)

### Enhancements

- [#336](https://github.com/terraform-linters/tflint-ruleset-aws/pull/336): autogenerated maintenance

### Chores

- [#338](https://github.com/terraform-linters/tflint-ruleset-aws/pull/338): Bump github.com/hashicorp/hcl/v2 from 2.11.1 to 2.12.0
- [#339](https://github.com/terraform-linters/tflint-ruleset-aws/pull/339): Bump github.com/google/go-cmp from 0.5.7 to 0.5.8
- [#340](https://github.com/terraform-linters/tflint-ruleset-aws/pull/340): Replace logger from the standard logger ([@wata727](https://github.com/wata727))
- [#341](https://github.com/terraform-linters/tflint-ruleset-aws/pull/341): Add E2E tests ([@wata727](https://github.com/wata727))

## 0.13.3 (2022-04-17)

### Enhancements

- [#324](https://github.com/terraform-linters/tflint-ruleset-aws/pull/324): autogenerated maintenance
- [#335](https://github.com/terraform-linters/tflint-ruleset-aws/pull/335): Lambda runtime deprecation updates (python3.6) ([@PatMyron](https://github.com/PatMyron))

### Chores

- [#328](https://github.com/terraform-linters/tflint-ruleset-aws/pull/328): chores: Remove snaker ([@wata727](https://github.com/wata727))
- [#329](https://github.com/terraform-linters/tflint-ruleset-aws/pull/329): Fix rule template for rule generator ([@wata727](https://github.com/wata727))
- [#330](https://github.com/terraform-linters/tflint-ruleset-aws/pull/330): Bump github.com/terraform-linters/tflint-plugin-sdk from 0.10.0 to 0.10.1
- [#333](https://github.com/terraform-linters/tflint-ruleset-aws/pull/333): style: format rules/api/rule.go.tmpl and run `go generate ./...` ([@suzuki-shunsuke](https://github.com/suzuki-shunsuke))
- [#334](https://github.com/terraform-linters/tflint-ruleset-aws/pull/334): Bump actions/setup-go from 2 to 3

## 0.13.2 (2022-03-29)

### BugFixes

- [#327](https://github.com/terraform-linters/tflint-ruleset-aws/pull/327): Suppress unevaluable/unknown/null errors on provider block eval ([@wata727](https://github.com/wata727))

### Chores

- [#321](https://github.com/terraform-linters/tflint-ruleset-aws/pull/321): Bump peter-evans/create-pull-request from 3 to 4
- [#322](https://github.com/terraform-linters/tflint-ruleset-aws/pull/322): Bump actions/cache from 2 to 3

## 0.13.1 (2022-03-28)

### BugFixes

- [#320](https://github.com/terraform-linters/tflint-ruleset-aws/pull/320): Call EnsureNoError after evaluating `aws_route` expression so unevaluable expressions are skipped ([@jandersen-plaid](https://github.com/jandersen-plaid))

## 0.13.0 (2022-03-27)

The minimum supported version of TFLint has changed in this version. TFLint v0.35.0+ is required for this plugin to work.

### Breaking Changes

- [#274](https://github.com/terraform-linters/tflint-ruleset-aws/pull/274): Bump tflint-plugin-sdk for gRPC-based new plugin system ([@wata727](https://github.com/wata727))
- [#310](https://github.com/terraform-linters/tflint-ruleset-aws/pull/310): aws_spot_instance_request.instance_interruption_behaviour renamed ([@PatMyron](https://github.com/PatMyron))
- [#317](https://github.com/terraform-linters/tflint-ruleset-aws/pull/317) [#318](https://github.com/terraform-linters/tflint-ruleset-aws/pull/318): Update aws-sdk-go and AWS provider rules ([@wata727](https://github.com/wata727))
  - Removed `aws_amplify_domain_association_invalid_domain_name` rule.

### Enhancements

- [#309](https://github.com/terraform-linters/tflint-ruleset-aws/pull/309): refactor previous generation instance type rules ([@PatMyron](https://github.com/PatMyron))
- [#315](https://github.com/terraform-linters/tflint-ruleset-aws/pull/315): rules: Add new  `aws_elastic_beanstalk_environment_invalid_name_format` rule ([@samhpickering](https://github.com/samhpickering))

### Chores

- [#312](https://github.com/terraform-linters/tflint-ruleset-aws/pull/312): Bump actions/checkout from 2 to 3
- [#313](https://github.com/terraform-linters/tflint-ruleset-aws/pull/313): Bump github.com/hashicorp/aws-sdk-go-base from 1.0.0 to 1.1.0
- [#314](https://github.com/terraform-linters/tflint-ruleset-aws/pull/314): go 1.18 ([@PatMyron](https://github.com/PatMyron))
- [#319](https://github.com/terraform-linters/tflint-ruleset-aws/pull/319): Bump GoReleaser version ([@wata727](https://github.com/wata727))

## 0.12.0 (2022-01-28)

### Enhancements

- [#291](https://github.com/terraform-linters/tflint-ruleset-aws/pull/291): mapping aws_devicefarm ([@PatMyron](https://github.com/PatMyron))
- [#292](https://github.com/terraform-linters/tflint-ruleset-aws/pull/292): mapping aws_fsx ([@PatMyron](https://github.com/PatMyron))
- [#293](https://github.com/terraform-linters/tflint-ruleset-aws/pull/293) [#294](https://github.com/terraform-linters/tflint-ruleset-aws/pull/294) [#300](https://github.com/terraform-linters/tflint-ruleset-aws/pull/300) [#301](https://github.com/terraform-linters/tflint-ruleset-aws/pull/301): Update terraform-provider-aws and aws-sdk submodule
  - terraform-provider-aws: v3.70.0 -> v3.74.0
  - aws-sdk: v1.42.25 -> v1.42.43
- [#295](https://github.com/terraform-linters/tflint-ruleset-aws/pull/295): mapping aws_memorydb ([@PatMyron](https://github.com/PatMyron))

### Chores

- [#286](https://github.com/terraform-linters/tflint-ruleset-aws/pull/286): Updated missing documentation ([@Rihoj](https://github.com/Rihoj))
- [#290](https://github.com/terraform-linters/tflint-ruleset-aws/pull/290): automating maintenance with Github actions ([@PatMyron](https://github.com/PatMyron))
- [#298](https://github.com/terraform-linters/tflint-ruleset-aws/pull/298): Bump github.com/google/go-cmp from 0.5.6 to 0.5.7
- [#302](https://github.com/terraform-linters/tflint-ruleset-aws/pull/302): git submodule update in automated maintenance ([@PatMyron](https://github.com/PatMyron))

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

