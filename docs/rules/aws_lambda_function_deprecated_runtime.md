# aws_lambda_function_deprecated_runtime

Checks to see if a lambda function has been set with a runtime that has reached end of support or is approaching a deprecation milestone (function creation blocked, function updates blocked).

Runtime lifecycle data is generated from the [AWS Lambda runtimes documentation](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html). Dates that were in the future at generation time are speculative and noted as such in the warning message.

## Example

```hcl
resource "aws_lambda_function" "function" {
	function_name = "test_function"
	role = "test_role"
	runtime = "python2.7"
}
```

```
$ tflint
1 issue(s) found:

Warning: The "python2.7" runtime has reached end of support and function updates are blocked (aws_lambda_function_deprecated_runtime)

  on template.tf line 4:
   4:   runtime = "python2.7"

```

## Why

AWS deprecates Lambda runtimes in phases:

1. **End of support** — the runtime no longer receives security patches
2. **Block function create** — new functions cannot be created with this runtime
3. **Block function update** — existing functions cannot be updated

The warning message reflects which phase the runtime is in:

**Confirmed** (date had already passed when data was last generated):

- `The "python2.7" runtime has reached end of support`
- `The "python2.7" runtime has reached end of support and new function creation is blocked`
- `The "python2.7" runtime has reached end of support and function updates are blocked`

**Speculative** (date was still in the future when data was last generated):

- `The "nodejs22.x" runtime was scheduled to reach end of support on Apr 30, 2027 (as of Mar 29, 2026)`
- `The "python3.9" runtime has reached end of support and new function creation was scheduled to be blocked on Aug 31, 2026 (as of Mar 29, 2026)`
- `The "python3.9" runtime has reached end of support and function updates were scheduled to be blocked on Sep 30, 2026 (as of Mar 29, 2026)`

Speculative messages include the scheduled date and when the data was last updated so you can judge whether to act on them or add a `tflint-ignore` comment.

## How To Fix

Update to a newer runtime. Supported runtimes can be found [here](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html).
