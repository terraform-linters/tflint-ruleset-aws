package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

const testTagRule = `
rule "aws_resource_tags" {
	enabled  = true
	required = ["A", "B"]
	values   = { A: ["1", "foo"], B: ["2", "bar"] }
} `

func Test_AwsResourceTags(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
		RaiseErr error
	}{
		// basic assertions
		{
			Name: "no tags assigned and no rules set",
			Content: `
			provider "aws" { }
			resource "aws_instance" "ec2_instance" { }
			resource "aws_instance" "ec2_instance" { }`,
			Config:   `rule "aws_resource_tags" { enabled = true }`,
			Expected: helper.Issues{},
		},
		{
			Name: "no tags assigned",
			Content: `
			provider "aws" { }
			resource "aws_instance" "ec2_instance" { }`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Tag 'A' is required. Tag 'B' is required.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 3, Column: 4},
						End:      hcl.Pos{Line: 3, Column: 42},
					},
				},
			},
		},
		{
			Name: "no tags assigned and rule with only required set",
			Content: `
			provider "aws" { }
			resource "aws_instance" "ec2_instance" { }`,
			Config: `rule "aws_resource_tags" {
				required = ["A", "B"]
				enabled = true
			}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Tag 'A' is required. Tag 'B' is required.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 3, Column: 4},
						End:      hcl.Pos{Line: 3, Column: 42},
					},
				},
			},
		},
		{
			Name: "no tags assigned and rule with only values set",
			Content: `
			provider "aws" { }
			resource "aws_instance" "ec2_instance" { }`,
			Config: `rule "aws_resource_tags" {
				values = { A: ["1"], B: ["2"] }
				enabled = true
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "resource tag assignment with invalid values",
			Content: `
			resource "aws_instance" "ec2_instance" { tags = { A = "0" } }
			resource "aws_instance" "ec2_instance" { tags = { B = "0" } }`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Received '0' for tag 'A', expected one of '1, foo'. Tag 'B' is required.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 52},
						End:      hcl.Pos{Line: 2, Column: 63},
					},
				},
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Received '0' for tag 'B', expected one of '2, bar'. Tag 'A' is required.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 3, Column: 52},
						End:      hcl.Pos{Line: 3, Column: 63},
					},
				},
			},
		},
		{
			Name: "resource tag assignment with valid values",
			Content: `resource "aws_instance" "ec2_instance" {
				tags = { A = "1", B = "bar" }
			}`,
			Config:   testTagRule,
			Expected: helper.Issues{},
		},
		{
			Name: "resource tag assignment with unconfigured tag rule",
			Content: `resource "aws_instance" "ec2_instance" {
				tags = { A = "1", B = "bar", C = "3" }
			}`,
			Config:   testTagRule,
			Expected: helper.Issues{},
		},
		// exclude
		{
			Name: "resource with invalid tags is excluded via rule",
			Content: `
			provider "aws" { }
			resource "aws_instance" "ec2_instance_one" { tags = { A: "0", B: "0" } }
			resource "aws_instance" "ec2_instance_two" { tags = { A: "xar", B: "zar" } }
			resource "aws_s3_bucket" "s3_bucket_one" { tags = { A: "xar", B: "zar" } }`,
			Config: `
			rule "aws_resource_tags" {
				enabled = true
				values = { A: ["1", "foo"], B: ["2", "bar"] }
				exclude = ["aws_instance"]
			}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Received 'xar' for tag 'A', expected one of '1, foo'. Received 'zar' for tag 'B', expected one of '2, bar'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 54},
						End:      hcl.Pos{Line: 5, Column: 76},
					},
				},
			},
		},
		// assignment via variables
		{
			Name: "valid provider tags assigned via variables",
			Content: `
			variable "tags" {
				type = map(string, string)
				default = {  A = "1", B = "2" }
			}
			provider "aws" {
				default_tags { tags = var.tags }
			}
			resource "aws_instance" "ec2_instance_one" {}
			resource "aws_instance" "ec2_instance_two" {}`,
			Config:   testTagRule,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid provider tags assigned via variables",
			Content: `
			variable "tags" {
				type = map(string, string)
				default = {  A = "1", B = "zar" }
			}
			provider "aws" {
				default_tags { tags = var.tags }
			}
			resource "aws_instance" "ec2_instance_one" {}`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Received 'zar' for tag 'B', expected one of '2, bar'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 9, Column: 4},
						End:      hcl.Pos{Line: 9, Column: 46},
					},
				},
			},
		},
		{
			Name: "invalid resource tags assigned with variables",
			Content: `
			variable "tags" {
				type = map(string, string)
				default = {  A = "1", B = "zar" }
			}
			provider "aws" { }
			resource "aws_instance" "ec2_instance_one" { tags = var.tags }`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Received 'zar' for tag 'B', expected one of '2, bar'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 7, Column: 56},
						End:      hcl.Pos{Line: 7, Column: 64},
					},
				},
			},
		},
		// provider.default_tags
		{
			Name: "valid provider tags assigned",
			Content: `
			provider "aws" {
				default_tags { tags = { A = "1", B = "2" } }
			}
			resource "aws_instance" "ec2_instance_one" {}
			resource "aws_instance" "ec2_instance_two" {}`,
			Config:   testTagRule,
			Expected: helper.Issues{},
		},
		{
			Name: "valid provider tags assigned and invalid tags assigned to resources",
			Content: `
			provider "aws" {
				default_tags { tags = { A = "1", B = "2" } }
			}
			resource "aws_instance" "ec2_instance_one" { tags = { A = "0" }}`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Received '0' for tag 'A', expected one of '1, foo'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 56},
						End:      hcl.Pos{Line: 5, Column: 67},
					},
				},
			},
		},
		{
			Name: "invalid default provider tags assigned and no tags assigned to resource",
			Content: `
			provider "aws" {
				default_tags { tags = { A = "0", B = "foo" } }
			}
			resource "aws_s3_bucket" "bucket" {}`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Received '0' for tag 'A', expected one of '1, foo'. Received 'foo' for tag 'B', expected one of '2, bar'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 4},
						End:      hcl.Pos{Line: 5, Column: 37},
					},
				},
			},
		},
		{
			Name: "multiple providers with default tags assigned",
			Content: `
			provider "aws" {
				alias = "one"
				default_tags { tags = { A = "1" } }
			}
			provider "aws" {
				alias = "two"
				default_tags { tags = { B = "2" } }
			}
			resource "aws_instance" "ec2_instance" { provider = aws.one }
			resource "aws_instance" "ec2_instance" { provider = aws.two }`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Tag 'B' is required.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 10, Column: 4},
						End:      hcl.Pos{Line: 10, Column: 42},
					},
				},
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Tag 'A' is required.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 11, Column: 4},
						End:      hcl.Pos{Line: 11, Column: 42},
					},
				},
			},
		},
		// aws_autoscaling_group
		{
			Name:    "autoscaling group with no tags",
			Content: `resource "aws_autoscaling_group" "asg" {}`,
			Config:  testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Tag 'A' is required. Tag 'B' is required.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 39},
					},
				},
			},
		},
		{
			Name: "autoscaling group with invalid tags assigned with block syntax",
			Content: `resource "aws_autoscaling_group" "asg" {
				tag {
					key = "A"
					value = "2"
					propagate_at_launch = true
				}
				tag {
					key = "B"
					value = "2"
					propagate_at_launch = true
				}
			}`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Received '2' for tag 'A', expected one of '1, foo'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 39},
					},
				},
			},
		},
		{
			Name: "autoscaling group with missing tags via block syntax",
			Content: `resource "aws_autoscaling_group" "asg" {
				tag {
					key = "A"
					value = "foo"
					propagate_at_launch = true
				}
			}`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Tag 'B' is required.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 39},
					},
				},
			},
		},
		{
			Name: "autoscaling group with missing tags",
			Content: `resource "aws_autoscaling_group" "asg" {
				tags = [
					{ key = "A", value = "foo", propagate_at_launch = true }
				]
			}`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Tag 'B' is required.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 12},
						End:      hcl.Pos{Line: 4, Column: 6},
					},
				},
			},
		},
		{
			Name: "autoscaling group with valid tags assigned with block syntax",
			Content: `
			resource "aws_autoscaling_group" "asg" {
				tag {
					key = "A"
					value = "foo"
				}

				tag {
					key = "B"
					value = "bar"
				}
			}`,
			Config:   testTagRule,
			Expected: helper.Issues{},
		},
		{
			Name: "autoscaling group with valid tag assignment",
			Content: `
			resource "aws_autoscaling_group" "asg" {
				tags = [
					{ key = "A", value = "foo", propagate_at_launch = true },
					{ key = "B", value = "bar", propagate_at_launch = true }
				]
			}`,
			Config:   testTagRule,
			Expected: helper.Issues{},
		},
		{
			Name: "autoscaling group with mixed tag assignment",
			Content: `
			resource "aws_autoscaling_group" "asg" {
				tag {
					key = "A"
					value = "foo"
				}

				tags = [
					{ key = "A", value = "foo", propagate_at_launch = true },
					{ key = "B", value = "bar", propagate_at_launch = true }
				]
			}`,
			Config: testTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceTagsRule(),
					Message: "Only tag block or tags attribute may be present, but found both",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 42},
					},
				},
			},
		},
	}

	rule := NewAwsResourceTagsRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"module.tf": tc.Content, ".tflint.hcl": tc.Config})

			err := rule.Check(runner)

			if tc.RaiseErr == nil && err != nil {
				t.Fatalf("Unexpected error occurred in test \"%s\": %s", tc.Name, err)
			}

			assert.Equal(t, tc.RaiseErr, err)

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
