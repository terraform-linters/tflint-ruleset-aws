# aws_lambda_function_deprecated_runtime

Checks to see if a lambda function has been set with a runtime that is deprecated. This can show up as either "end of support" or "end of life" depending on the phase of deprecation it is currently in.

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

Error: The "python2.7" runtime has reached the end of support. (aws_lambda_function_deprecated_runtime)

  on template.tf line 4:
   4:   runtime  = "python2.7" // end of support reached!

```

## Why

AWS no longer supports these runtimes.

## How To Fix

Update to a newer runtime. Supported runtimes can be found [here](https://docs.aws.amazon.com/lambda/latest/dg/runtime-support-policy.html)
