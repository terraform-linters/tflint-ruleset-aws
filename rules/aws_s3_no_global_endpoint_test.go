package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3NoGlobalEndpoint(t *testing.T) {
	filename := "resource.tf"

	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "unrelated expressions",
			Content: `
output "test1" {
  value = var.whatever
}

output "test2" {
  value = "testing"
}

output "test3" {
  value = aws_iam_role.test.arn
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "regional endpoint used",
			Content: `
output "test" {
  value = aws_s3_bucket.test.bucket_regional_domain_name
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "global endpoint used",
			Content: `
output "test" {
  value = aws_s3_bucket.test.bucket_domain_name
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 48},
					},
				},
			},
		},
		{
			Name: "global endpoint used from aws_s3_bucket data source",
			Content: `
output "test" {
  value = data.aws_s3_bucket.test.bucket_domain_name
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 53},
					},
				},
			},
		},
		{
			Name: "global endpoint interpolation",
			Content: `
output "test" {
  value = "interpolation test: ${aws_s3_bucket.test.bucket_domain_name}"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 3, Column: 34},
						End:      hcl.Pos{Line: 3, Column: 71},
					},
				},
			},
		},
		{
			Name: "global endpoint interpolation multiple expressions",
			Content: `
output "test" {
  value = "interpolation test: ${aws_s3_bucket.test.bucket_domain_name} ${var.whatever}"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 3, Column: 34},
						End:      hcl.Pos{Line: 3, Column: 71},
					},
				},
			},
		},
		{
			Name: "global endpoint with count",
			Content: `
output "test1" {
  value = aws_s3_bucket.test[0].bucket_domain_name
}

output "test2" {
  value = aws_s3_bucket.test[length(var.bucket_names) - 1].bucket_domain_name
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 51},
					},
				},
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 7, Column: 11},
						End:      hcl.Pos{Line: 7, Column: 78},
					},
				},
			},
		},
		{
			Name: "global endpoint in a for expression",
			Content: `
output "test" {
  value = {
    for bucket_name, bucket in aws_s3_bucket.test: bucket_name => bucket.bucket_domain_name
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 4, Column: 67},
						End:      hcl.Pos{Line: 4, Column: 92},
					},
				},
			},
		},
		{
			Name: "global endpoint with foreach",
			Content: `
resource "aws_s3_bucket" "test" {
  for_each = toset(var.bucket_names)

  bucket = each.key
}

resource "aws_ssm_parameter" "test" {
  for_each = aws_s3_bucket.test

  name  = each.value.id
  type  = "String"
  value = each.value.bucket_domain_name
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 13, Column: 11},
						End:      hcl.Pos{Line: 13, Column: 40},
					},
				},
			},
		},
		{
			Name: "global endpoint in a dynamic block",
			Content: `
resource "aws_s3_bucket" "test" {
  for_each = toset(var.bucket_names)

  bucket = each.key
}

resource "aws_cloudfront_distribution" "test" {
  dynamic "origin" {
    for_each = aws_s3_bucket.test

    content {
      origin_id   = origin.value.id
      domain_name = origin.value.bucket_domain_name
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 14, Column: 11},
						End:      hcl.Pos{Line: 14, Column: 40},
					},
				},
			},
		},
		{
			Name: "global endpoint literal",
			Content: `
output "test" {
  value = "test.s3.amazonaws.com"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 33},
					},
				},
			},
		},
		{
			Name: "global endpoint literal with interpolation",
			Content: `
output "test" {
  value = "${var.bucket_name}.s3.amazonaws.com"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3NoGlobalEndpointRule(),
					Message: "`bucket_domain_name` returns the legacy s3 global endpoint, use `bucket_regional_domain_name` instead",
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 47},
					},
				},
			},
		},
	}

	rule := NewAwsS3NoGlobalEndpointRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
