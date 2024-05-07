package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

const testInvalidTagRule = `
rule "aws_resource_invalid_tags" {
  enabled = true
  tags = { A: ["1", "foo"], B: ["2", "bar"] }
} `

func Test_AwsResourceInvalidTags(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
		RaiseErr error
	}{
		{
			Name: "no tags assigned",
			Content: `
			provider "aws" { region = "us-east-1" }
			resource "aws_instance" "ec2_instance" { }
			resource "aws_instance" "ec2_instance" { }`,
			Config:   testInvalidTagRule,
			Expected: helper.Issues{},
		},
		{
			Name: "resource with invalid explicit tags is excluded via rule",
			Content: `
			resource "aws_instance" "ec2_instance_one" { tags = { A: "0", B: "0" } }
			resource "aws_instance" "ec2_instance_two" { tags = { A: "xar", B: "zar" } }
			resource "aws_s3_bucket" "s3_bucket_one" { tags = { A: "xar", B: "zar" } }`,
			Config: `
			rule "aws_resource_invalid_tags" {
				enabled = true
				tags = { A: ["1", "foo"], B: ["2", "bar"] }
				exclude = ["aws_instance"]
			}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceInvalidTagsRule(),
					Message: "Received 'xar' for tag 'A', expected one of '1,foo'. Received 'zar' for tag 'B', expected one of '2,bar'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 4, Column: 54},
						End:      hcl.Pos{Line: 4, Column: 76},
					},
				},
			},
		},
		{
			Name: "valid provider tags assigned and no explicit tags assigned to resources",
			Content: `
			provider "aws" {
				default_tags { tags = { A = "1", B = "2" } }
			}
			resource "aws_instance" "ec2_instance_one" {}
			resource "aws_instance" "ec2_instance_two" {}`,
			Config:   testInvalidTagRule,
			Expected: helper.Issues{},
		},
		{
			Name: "valid provider tags assigned and invalid explicit tags assigned to resources",
			Content: `
			provider "aws" {
				default_tags { tags = { A = "1", B = "2" } }
			}
			resource "aws_instance" "ec2_instance_one" { tags = { A = "0" }}`,
			Config: testInvalidTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceInvalidTagsRule(),
					Message: "Received '0' for tag 'A', expected one of '1,foo'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 56},
						End:      hcl.Pos{Line: 5, Column: 67},
					},
				},
			},
		},
		{
			Name: "valid default provider tags assigned and no tags assigned to autoscaling group",
			Content: `
			provider "aws" {
				default_tags { tags = { A = "1", B = "2" } }
			}
			resource "aws_autoscaling_group" "asg" {}`,
			Config:   testInvalidTagRule,
			Expected: helper.Issues{},
		},
		// NOTE: This test surfaces the unknown relationship between Provider default tags
		//       and AutoScaling Groups tags.
		{
			Name: "invalid default provider tags assigned and no tags assigned to autoscaling group",
			Content: `
			provider "aws" {
				default_tags { tags = { A = "0", B = "0" } }
			}
			resource "aws_autoscaling_group" "asg" {}`,
			Config:   testInvalidTagRule,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid default provider tags assigned and no tags assigned to resource",
			Content: `
			provider "aws" {
				default_tags { tags = { A = "0", B = "foo" } }
			}
			resource "aws_s3_bucket" "bucket" {}`,
			Config: testInvalidTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceInvalidTagsRule(),
					Message: "Received '0' for tag 'A', expected one of '1,foo'. Received 'foo' for tag 'B', expected one of '2,bar'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 5, Column: 4},
						End:      hcl.Pos{Line: 5, Column: 37},
					},
				},
			},
		},
		{
			Name: "multiple providers with explicit tags assigned and no custom tags assigned to resources",
			Content: `
			provider "aws" {
				alias = "one"
				default_tags { tags = { A = "1" } }
			}
			provider "aws" {
				alias = "two"
				default_tags { tags = { B = "2" } }
			}
			resource "aws_instance" "ec2_instance" { provider = "one" }
			resource "aws_instance" "ec2_instance" { provider = "two" }`,
			Config:   testInvalidTagRule,
			Expected: helper.Issues{},
		},
		{
			Name: "explicit resource tag assignment with invalid values",
			Content: `
			resource "aws_instance" "ec2_instance" { tags = { A = "0" } }
			resource "aws_instance" "ec2_instance" { tags = { B = "0" } }`,
			Config: testInvalidTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceInvalidTagsRule(),
					Message: "Received '0' for tag 'A', expected one of '1,foo'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 52},
						End:      hcl.Pos{Line: 2, Column: 63},
					},
				},
				{
					Rule:    NewAwsResourceInvalidTagsRule(),
					Message: "Received '0' for tag 'B', expected one of '2,bar'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 3, Column: 52},
						End:      hcl.Pos{Line: 3, Column: 63},
					},
				},
			},
		},
		{
			Name:     "explicit resource tag assignment with valid values",
			Content:  `resource "aws_instance" "ec2_instance" { tags = { A = "1", B = "bar" } }`,
			Config:   testInvalidTagRule,
			Expected: helper.Issues{},
		},
		{
			Name:     "explicit resource tag assignment with unconfigured tag rule",
			Content:  `resource "aws_instance" "ec2_instance" { tags = { A = "1", B = "bar", C = "3" } }`,
			Config:   testInvalidTagRule,
			Expected: helper.Issues{},
		},
		{
			Name: "explicit autoscaling group resource with invalid tags",
			Content: `
			resource "aws_autoscaling_group" "asg" {
				tag {
					key = "A"
					value = "0"
				}
				tag {
					key = "B"
					value = "foo"
				}
			}`,
			Config: testInvalidTagRule,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceInvalidTagsRule(),
					Message: "Received '0' for tag 'A', expected one of '1,foo'. Received 'foo' for tag 'B', expected one of '2,bar'.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 42},
					},
				},
			},
		},
		{
			Name: "autoscaling group with valid explicit tag assignment",
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
			Config:   testInvalidTagRule,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsResourceInvalidTagsRule()

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
