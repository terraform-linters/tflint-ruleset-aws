# aws_iam_policy_sid_invalid_characters

AWS IAM policies have statement ids within each statement that must fit the regexp `^[a-zA-Z0-9]+$` otherwise the policy will fail.

## Example

```hcl
resource "aws_iam_policy" "policy" {
	name = "test_policy"
	role = "test_role"
	policy = <<-EOF
{
	"Version": "2012-10-17",
	"Statement": [
	  {
			"Sid": "This contains invalid-characters.", //invalid Sid
	    "Action": [
	      "ec2:Describe*"
			],
			"Effect": "Allow",
			"Resource": "arn:aws:s3:::<bucketname>/*"
	  }
	]
}
EOF
}
```

```
$ tflint
1 issue(s) found:

Error: The policy's sid contains invalid characters. (aws_iam_policy_sid_invalid_characters)

  on template.tf line 3:
   3:   policy = <<-EOF
 {
 	"Version": "2012-10-17",
 	"Statement": [
 	  {
 			"Sid": "This contains invalid-characters.", //invalid Sid
 	    "Action": [
 	      "ec2:Describe*"
 			],
 			"Effect": "Allow",
 			"Resource": "arn:aws:s3:::<bucketname>/*"
 	  }
 	]
 }
 EOF

```

## Why

AWS will fail to create or update the policy if the Sid is not valid

## How To Fix

Make sure you Sid's fit the regexp `^[a-zA-Z0-9]+$`
