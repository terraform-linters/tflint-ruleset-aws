package ephemeral

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEphemeralResources(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
		Fixed    string
	}{
		{
			Name: "basic aws_eks_cluster_auth",
			Content: `
data "aws_eks_cluster_auth" "test" {
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEphemeralResourcesRule(),
					Message: `"aws_eks_cluster_auth" is a non-ephemeral data source, which means that all (sensitive) attributes are stored in state. Please use ephemeral resource "aws_eks_cluster_auth" instead.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 5},
					},
				},
			},
		},
	}

	rule := NewAwsEphemeralResourcesRule()

	for _, tc := range cases {
		filename := "resource.tf"
		runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}
		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
