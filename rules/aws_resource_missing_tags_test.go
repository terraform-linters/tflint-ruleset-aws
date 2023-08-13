package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsResourceMissingTags(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
		RaiseErr error
	}{
		{
			Name: "Wanted tags: Bar,Foo, found: bar,foo",
			Content: `
resource "aws_instance" "ec2_instance" {
  instance_type = "t2.micro"
  tags = {
    foo = "bar"
    bar = "baz"
  }
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceMissingTagsRule(),
					Message: "The resource is missing the following tags: \"Bar\", \"Foo\".",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 4, Column: 10},
						End:      hcl.Pos{Line: 7, Column: 4},
					},
				},
			},
		},
		{
			Name: "No tags",
			Content: `
resource "aws_instance" "ec2_instance" {
  instance_type = "t2.micro"
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceMissingTagsRule(),
					Message: "The resource is missing the following tags: \"Bar\", \"Foo\".",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 39},
					},
				},
			},
		},
		{
			Name: "Tags are correct",
			Content: `
resource "aws_instance" "ec2_instance" {
  instance_type = "t2.micro"
  tags = {
    Foo = "bar"
    Bar = "baz"
  }
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "AutoScaling Group with tag blocks and correct tags",
			Content: `
resource "aws_autoscaling_group" "asg" {
  tag {
    key = "Foo"
    value = "bar"
    propagate_at_launch = true
  }
  tag {
    key = "Bar"
    value = "baz"
    propagate_at_launch = true
  }
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "AutoScaling Group with tag blocks and incorrect tags",
			Content: `
resource "aws_autoscaling_group" "asg" {
  tag {
    key = "Foo"
    value = "bar"
    propagate_at_launch = true
  }
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceMissingTagsRule(),
					Message: "The resource is missing the following tags: \"Bar\".",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 39},
					},
				},
			},
		},
		{
			Name: "AutoScaling Group with tags attribute and correct tags",
			Content: `
resource "aws_autoscaling_group" "asg" {
  tags = [
    {
      key = "Foo"
      value = "bar"
      propagate_at_launch = true
    },
    {
      key = "Bar"
      value = "baz"
      propagate_at_launch = true
    }
  ]
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "AutoScaling Group with tags attribute and incorrect tags",
			Content: `
resource "aws_autoscaling_group" "asg" {
  tags = [
    {
      key = "Foo"
      value = "bar"
      propagate_at_launch = true
    }
  ]
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceMissingTagsRule(),
					Message: "The resource is missing the following tags: \"Bar\".",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 3, Column: 10},
						End:      hcl.Pos{Line: 9, Column: 4},
					},
				},
			},
		},
		{
			Name: "AutoScaling Group excluded from missing tags rule",
			Content: `
resource "aws_autoscaling_group" "asg" {
  tags = [
    {
      key = "Foo"
      value = "bar"
      propagate_at_launch = true
    }
  ]
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
  exclude = ["aws_autoscaling_group"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "AutoScaling Group with both tag block and tags attribute",
			Content: `
resource "aws_autoscaling_group" "asg" {
  tag {
    key = "Foo"
    value = "bar"
    propagate_at_launch = true
  }
  tags = [
    {
      key = "Foo"
      value = "bar"
      propagate_at_launch = true
    }
  ]
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceMissingTagsRule(),
					Message: "Only tag block or tags attribute may be present, but found both",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 39},
					},
				},
			},
		},
		{
			Name: "AutoScaling Group with no tags",
			Content: `
resource "aws_autoscaling_group" "asg" {
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceMissingTagsRule(),
					Message: "The resource is missing the following tags: \"Bar\", \"Foo\".",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 39},
					},
				},
			},
		},
		{
			Name: "Default tags multiple providers",
			Content: `
provider "aws" {
  default_tags {
    tags = {
      "Fooz": "Barz"
      "Bazz": "Quxz"
    }
  }
}

provider "aws" {
  alias = "foo"
  default_tags {
    tags = {
      "Bazz": "Quxz"
      "Fooz": "Barz"
    }
  }
}

resource "aws_instance" "ec2_instance" {
  instance_type = "t2.micro"
}

resource "aws_instance" "ec2_instance_alias" {
  provider = aws.foo
  instance_type = "t2.micro"
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Bazz", "Fooz"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Default Tags Are to Be overriden by resource specific tags",
			Content: `
provider "aws" {
  default_tags {
    tags = {
      "Foo": "Bar"
    }
  }
}

resource "aws_instance" "ec2_instance" {
  instance_type = "t2.micro"
  tags = {
    "Foo": "Bazz"
  }
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Resource specific tags are not needed if default tags are placed",
			Content: `
provider "aws" {
  default_tags {
    tags = {
      "Foo": "Bar"
    }
  }
}

resource "aws_instance" "ec2_instance" {
  instance_type = "t2.micro"
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Resource tags in combination with provider level tags",
			Content: `
provider "aws" {
  default_tags {
    tags = {
      "Foo": "Bar"
    }
  }
}

resource "aws_instance" "ec2_instance_fail" {
  instance_type = "t2.micro"
}


resource "aws_instance" "ec2_instance" {
  instance_type = "t2.micro"
  tags = {
    "Bazz": "Quazz"
  }
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bazz"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceMissingTagsRule(),
					Message: "The resource is missing the following tags: \"Bazz\".",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 10, Column: 1},
						End:      hcl.Pos{Line: 10, Column: 44},
					},
				},
			},
		},
		{
			Name: "Provider reference existent without tags definition",
			Content: `provider "aws" {
  alias = "west"
  region = "us-west-2"
}

resource "aws_ssm_parameter" "param" {
  provider = aws.west
  name = "test"
  type = "String"
  value = "test"
  tags = {
    Foo = "Bar"
  }
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Unknown value maps should be silently ignored",
			Content: `variable "aws_region" {
  default = "us-east-1"
  type = string
}

provider "aws" {
  region = "us-east-1"
  alias = "foo"
  default_tags {
    tags = {
      Owner = "Owner"
    }
  }
}

resource "aws_s3_bucket" "a" {
  provider = aws.foo
  name = "a"
}

resource "aws_s3_bucket" "b" {
  name = "b"
  tags = var.default_tags
}

variable "default_tags" {
  type = map(string)
}

provider "aws" {
  region = "us-east-1"
  alias = "bar"
  default_tags {
    tags = var.default_tags
  }
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Owner"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Not wholly known tags as value maps should work as long as each key is statically defined",
			Content: `variable "owner" {
  type = string
}

variable "project" {
  type = string
}

provider "aws" {
  region = "us-east-1"
  alias = "foo"
  default_tags {
    tags = {
      Owner = var.owner
      Project = var.project
    }
  }
}

resource "aws_s3_bucket" "a" {
  provider = aws.foo
  name = "a"
}

resource "aws_s3_bucket" "b" {
  name = "b"
  tags = {
    Owner = var.owner
    Project = var.project
  }
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Owner", "Project"]
}`,
			Expected: helper.Issues{},
		},
		{
			// In child modules, the provider declarations are passed implicitly or explicitly from the root module.
			// In this case, it is not possible to refer to the provider's default_tags.
			// Here, the strategy is to ignore default_tags rather than skip inspection of the tags.
			Name: "provider aliases within child modules",
			Content: `
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      configuration_aliases = [ aws.foo ]
    }
  }
}

resource "aws_instance" "ec2_instance" {
  provider = aws.foo
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo", "Bar"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceMissingTagsRule(),
					Message: "The resource is missing the following tags: \"Bar\", \"Foo\".",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 11, Column: 1},
						End:      hcl.Pos{Line: 11, Column: 39},
					},
				},
			},
		},
	}

	rule := NewAwsResourceMissingTagsRule()

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
