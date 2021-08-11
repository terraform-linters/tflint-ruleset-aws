# aws_iam_policy_too_long_policy

This makes sure that a IAM policy is not longer than the 6144 AWS character limit.

## Example

```hcl
resource "aws_iam_policy" "policy" {
	name        = "test_policy"
	path        = "/"
	description = "My test policy"
	policy = <<EOF
	{
		STRING LONGER THAN 6144
	}
EOF
}
```

```
$ tflint
1 issue(s) found:
Error: The policy length is %d characters and is limited to 6144 characters. (aws_iam_policy_too_long_policy)
  on template.tf line 6:
   6:   policy = "{POLICY}" // Policy is too long,!

```

## Why

Terraform does not check against this rule, but it will error if an apply is attempted.

## How To Fix

Update policy to reduce characters. Some methods are splitting into multiple policies, or switching to using a combination of deny and allow statements to optimize the policy.
