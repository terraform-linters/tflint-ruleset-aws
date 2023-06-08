package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsResourceMissingTags(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
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
			Name: "Provider reference no existent",
			Content: `provider "aws" {
  alias = "zoom"
  default_tags {
    tags = {
      "Foo": "Bar"
    }
  }
}

resource "aws_instance" "ec2_instance" {
  provider = aws.west
  instance_type = "t2.micro"
}`,
			Config: `
rule "aws_resource_missing_tags" {
  enabled = true
  tags = ["Foo"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsResourceMissingTagsRule(),
					Message: "The aws provider with alias \"west\" doesn't exist.",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 11, Column: 17},
						End:      hcl.Pos{Line: 11, Column: 22},
					},
				},
			},
		},
	}

	rule := NewAwsResourceMissingTagsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"module.tf": tc.Content, ".tflint.hcl": tc.Config})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
