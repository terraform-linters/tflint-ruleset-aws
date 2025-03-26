package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsProviderMissingDefaultTags(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "Default tags for provider",
			Content: `
provider "aws" {
  default_tags {
    tags = {
      "Fooz": "Barz"
      "Bazz": "Quxz"
    }
  }
}`,
			Config: `
rule "aws_provider_missing_default_tags" {
  enabled = true
  tags = ["Bazz", "Fooz"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Missing default tag for provider",
			Content: `
provider "aws" {
  default_tags {
    tags = {
      "Fooz": "Barz"
    }
  }
}`,
			Config: `
rule "aws_provider_missing_default_tags" {
  enabled = true
  tags = ["Bazz", "Fooz"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsProviderMissingDefaultTagsRule(),
					Message: "The provider is missing the following tags: \"Bazz\".",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 6, Column: 6},
					},
				},
			},
		},
		{
			Name: "Empty default tags for provider",
			Content: `
provider "aws" {
  default_tags {
    tags = {}
  }
}`,
			Config: `
rule "aws_provider_missing_default_tags" {
  enabled = true
  tags = ["Bazz", "Fooz"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsProviderMissingDefaultTagsRule(),
					Message: "The provider is missing the following tags: \"Bazz\", \"Fooz\".",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 4, Column: 5},
						End:      hcl.Pos{Line: 4, Column: 14},
					},
				},
			},
		},
		{
			Name: "Missing default tags block for provider with no alias",
			Content: `
provider "aws" {
}`,
			Config: `
rule "aws_provider_missing_default_tags" {
  enabled = true
  tags = ["Bazz", "Fooz"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsProviderMissingDefaultTagsRule(),
					Message: "The aws provider is missing the `default_tags` block",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 15},
					},
				},
			},
		},
		{
			Name: "Missing default tags block for provider with alias",
			Content: `
provider "aws" {
  alias = "test"
}`,
			Config: `
rule "aws_provider_missing_default_tags" {
  enabled = true
  tags = ["Bazz", "Fooz"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsProviderMissingDefaultTagsRule(),
					Message: "The aws provider with alias `test` is missing the `default_tags` block",
					Range: hcl.Range{
						Filename: "module.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 15},
					},
				},
			},
		},
		{
			Name: "default tags wholly unknown",
			Content: `
provider "aws" {
  default_tags {
    tags = var.tags
	}
}

variable "tags" {
  type = map(string)
}`,
			Config: `
rule "aws_provider_missing_default_tags" {
  enabled = true
  tags = ["Bazz", "Fooz"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "null key",
			Content: `
provider "aws" {
  default_tags {
      tags = {
        (var.tag): "bar"
      }
    }
	}

  // This variable MUST be set, otherwise the configuration is invalid
  // The rule should error when encountering this in the tags expression
  variable "tag" {
    type = string
  }
  `,
			Config: `
rule "aws_provider_missing_default_tags" {
  enabled = true
  tags = ["Bazz", "Fooz"]
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsProviderMissingDefaultTagsRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"module.tf": tc.Content, ".tflint.hcl": tc.Config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
