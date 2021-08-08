# aws_lambda_function_eol

Checks to see if a lambda function has been set with a runtime that has reached End of Life.

## Example

```hcl
resource "aws_lambda_function" "function" {
	function_name = "test_function"
	role = "test_role"
	runtime = "dotnetcore1.0"
}
```

```
$ tflint
1 issue(s) found:

Error: The "dotnetcore1.0" runtime has reached the end of life. (aws_lambda_function_eol)

  on template.tf line 4:
   4:   runtime  = "dotnetcore1.0" // end of life reached!

```

## Why

AWS no longer supports these runtimes.

## How To Fix

Update to a newer runtime. Supported runtimmes can be found [here](https://docs.aws.amazon.com/lambda/latest/dg/runtime-support-policy.html)
