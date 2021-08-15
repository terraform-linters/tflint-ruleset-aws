# aws_iam_group_policy_too_long

This makes sure that an IAM group policy is not longer than the 5120 AWS character limit.
## Example

```hcl
resource "aws_iam_group_policy" "policy" {
	name        = "test_policy"
	group       = "test_group"
	description = "My test policy"
	policy = <<EOF
	{
		STRING LONGER THAN 5120
	}
EOF
}
```

```
$ tflint
1 issue(s) found:
Error: The policy length is %d characters and is limited to 5120 characters. (aws_iam_policy_too_long_policy)
  on template.tf line 6:
   6:   policy = "{POLICY}" // Policy is too long,!

```

## Why

Terraform does not check against this rule, but it will error if an apply is attempted.

## How To Fix

Update policy to reduce characters. Some methods are splitting into multiple policies, or switching to using a combination of deny and allow statements to optimize the policy.
