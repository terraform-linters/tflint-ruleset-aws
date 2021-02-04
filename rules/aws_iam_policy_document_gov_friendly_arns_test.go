package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsIamDocumentPolicyGovFriendlyArns(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "Wanted arn: arn:aws-us-gov, found arn: arn:aws",
			Content: `
resource "aws_iam_policy_document" "example" {
  statement {
	sid = "1"
	actions = [
	  "s3:ListAllMyBuckets",
	  "s3:GetBucketLocation",
	]
	resources = [
	  "arn:aws:s3:::*",
	]
  }
}`,
			Config: `
rule "aws_iam_policy_document_gov_friendly_arns" {
  enabled = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMPolicyDocumentGovFriendlyArnsRule(),
					Message: "ARN detected in IAM policy document that could potentially fail in AWS GovCloud due to resource pattern: arn:aws:.*",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 9, Column: 14},
						End:      hcl.Pos{Line: 11, Column: 3},
					},
				},
			},
		},
		{
			Name: "Wanted arn: arn:aws-us-gov, found arn: arn:aws",
			Content: `
resource "aws_iam_policy_document" "example" {
  statement {
	sid = "1"
	actions = [
	  "s3:ListAllMyBuckets",
	  "s3:GetBucketLocation",
	]
	resources = [
	  "arn:aws-us-gov:s3:::*",
	]
  }
}`,
			Config: `
rule "aws_iam_policy_document_gov_friendly_arns" {
  enabled = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Wanted arn: arn:aws-us-gov, found arn: arn:aws",
			Content: `
resource "aws_iam_policy_document" "example" {
  statement {
	sid = "1"
	actions = [
	  "s3:ListAllMyBuckets",
	  "s3:GetBucketLocation",
	]
	resources = [
	  "*",
	]
	principals {
	  type        = "Service"
	  identifiers = ["firehose.amazonaws.com"]
	}
  }
}`,
			Config: `
rule "aws_iam_policy_document_gov_friendly_arns" {
  enabled = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Wanted arn: arn:aws-us-gov, found arn: arn:aws",
			Content: `
resource "aws_iam_policy_document" "example" {
  statement {
	sid = "1"
	actions = [
	  "s3:ListAllMyBuckets",
	  "s3:GetBucketLocation",
	]
	resources = [
	  "arn:aws:s3:::*",
	]
  }
  statement {
	sid = "1"
	actions = [
	  "s3:ListAllMyBuckets",
	  "s3:GetBucketLocation",
	]
	resources = [
	  "arn:aws:s3:::*",
	]
  }
}`,
			Config: `
rule "aws_iam_policy_document_gov_friendly_arns" {
  enabled = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMPolicyDocumentGovFriendlyArnsRule(),
					Message: "ARN detected in IAM policy document that could potentially fail in AWS GovCloud due to resource pattern: arn:aws:.*",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 9, Column: 14},
						End:      hcl.Pos{Line: 11, Column: 3},
					},
				},
				{
					Rule:    NewAwsIAMPolicyDocumentGovFriendlyArnsRule(),
					Message: "ARN detected in IAM policy document that could potentially fail in AWS GovCloud due to resource pattern: arn:aws:.*",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 19, Column: 14},
						End:      hcl.Pos{Line: 21, Column: 3},
					},
				},
			},
		},
		{
			Name: "Wanted arn: arn:aws-us-gov, found arn: arn:aws",
			Content: `
resource "aws_iam_policy_document" "example" {
  statement {
	sid = "1"
	actions = [
	  "s3:ListAllMyBuckets",
	  "s3:GetBucketLocation",
	]
	resources = [
	  "arn:aws-us-gov:s3:::*",
	]
  }
  statement {
	sid = "1"
	actions = [
	  "s3:ListAllMyBuckets",
	  "s3:GetBucketLocation",
	]
	resources = [
	  "arn:aws-us-gov:s3:::*",
	]
  }
}`,
			Config: `
rule "aws_iam_policy_document_gov_friendly_arns" {
  enabled = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Wanted arn: arn:aws-us-gov, found arn: arn:aws",
			Content: `
resource "aws_iam_policy_document" "example" {
  statement {
	sid = "1"
	actions = [
	  "s3:ListAllMyBuckets",
	  "s3:GetBucketLocation",
	]
	resources = [
	  "arn:aws-us-gov:s3:::*",
	]
  }
  statement {
	sid = "1"
	actions = [
	  "s3:ListAllMyBuckets",
	  "s3:GetBucketLocation",
	]
	resources = [
	  "arn:aws-us-gov:s3:::*",
	  "arn:aws:s3:::*",
	]
  }
}`,
			Config: `
rule "aws_iam_policy_document_gov_friendly_arns" {
  enabled = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsIAMPolicyDocumentGovFriendlyArnsRule(),
					Message: "ARN detected in IAM policy document that could potentially fail in AWS GovCloud due to resource pattern: arn:aws:.*",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 19, Column: 14},
						End:      hcl.Pos{Line: 22, Column: 3},
					},
				},
			},
		},
	}

	rule := NewAwsIAMPolicyDocumentGovFriendlyArnsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"module.tf": tc.Content, ".tflint.hcl": tc.Config})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
