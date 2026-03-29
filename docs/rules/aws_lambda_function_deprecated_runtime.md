# aws_lambda_function_deprecated_runtime

Checks to see if a lambda function uses a runtime that has reached end of support. Runtime lifecycle data is generated from the [AWS Lambda runtimes documentation](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html).

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

The message reflects which deprecation phase the runtime is in. When a date had not yet passed at the time the rule data was last generated, the message notes the scheduled date and when the data was last updated:

```
Warning: The "nodejs22.x" runtime was scheduled to reach end of support on Apr 30, 2027 (as of Mar 29, 2026) (aws_lambda_function_deprecated_runtime)
```

## Why

AWS deprecates Lambda runtimes in phases: end of support (no more security patches), then blocking new function creation, then blocking function updates.

## How To Fix

Update to a newer runtime. Supported runtimes can be found [here](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html).
