# aws_cloudwatch_log_group_retention

Ensures that `aws_cloudwatch_log_group` resources have `retention_in_days` set to at least 365 days.

## Configuration

This rule is disabled by default. Enable it in your `.tflint.hcl`:

```hcl
rule "aws_cloudwatch_log_group_retention" {
  enabled = true
}
```

## Example

```hcl
resource "aws_cloudwatch_log_group" "bad" {
  name              = "/aws/lambda/my-function"
  retention_in_days = 30
}
```

```
$ tflint
1 issue(s) found:

Warning: The CloudWatch log group retention is 30 days. Set to at least 365 days. (aws_cloudwatch_log_group_retention)

  on template.tf line 3:
   3:   retention_in_days = 30
```

## Why

Log groups with short retention periods may lose data needed for auditing, compliance, or debugging. Many compliance frameworks (SOC 2, PCI-DSS, HIPAA) require at least one year of log retention. This rule also flags log groups missing `retention_in_days` entirely, which retain logs indefinitely and may incur unnecessary costs.

## How To Fix

Set `retention_in_days` to 365 or greater:

```hcl
resource "aws_cloudwatch_log_group" "good" {
  name              = "/aws/lambda/my-function"
  retention_in_days = 365
}
```
