# aws_elastic_beanstalk_environment_name_invalid_format

Ensure Elastic Beanstalk environment name matches allowed format.

## Example

```hcl
resource "aws_elastic_beanstalk_environment" "tfenvtest" {
	name                = "tf_test_name"
	application         = "tf-test-name"
	solution_stack_name = "64bit Amazon Linux 2015.03 v2.0.3 running Go 1.4"
}
```

```
$ tflint
1 issue(s) found:

Error: tf_test_name does not match constraint: must contain only letters, digits, and the dash character and may not start or end with a dash (^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$) (aws_elastic_beanstalk_environment_name_invalid_format)

  on example.tf line 2:
   2: 	name                = "tf_test_name"

```

## Why

When attempting to create the resource, Terraform will return the error:
```
Error: InvalidParameterValue: Value tf_test_name at 'EnvironmentName' failed to satisfy constraint: Member must contain only letters, digits, and the dash character and may not start or end with a dash
status code: 400
```

## How To Fix

Ensure your environment name consists only of letters, digits, and the dash character, and does not start or end with a dash.
The regex used is `^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$`
