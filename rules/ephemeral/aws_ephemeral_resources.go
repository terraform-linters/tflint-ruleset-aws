package ephemeral

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
)

// AwsEphemeralResourcesRule checks for data sources which can be replaced by ephemeral resources
type AwsEphemeralResourcesRule struct {
	tflint.DefaultRule

	replacingEphemeralResources []string
}

// NewAwsEphemeralResourcesRule returns new rule with default attributes
func NewAwsEphemeralResourcesRule() *AwsEphemeralResourcesRule {
	return &AwsEphemeralResourcesRule{
		replacingEphemeralResources: replacingEphemeralResources,
	}
}

// Name returns the rule name
func (r *AwsEphemeralResourcesRule) Name() string {
	return "aws_ephemeral_resources"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsEphemeralResourcesRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsEphemeralResourcesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsEphemeralResourcesRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks if there is an ephemeral resource which can replace an data source
func (r *AwsEphemeralResourcesRule) Check(runner tflint.Runner) error {
	for _, resourceType := range r.replacingEphemeralResources {
		resources, err := GetDataSourceContent(runner, resourceType, &hclext.BodySchema{}, nil)
		if err != nil {
			return err
		}

		for _, resource := range resources.Blocks {
			if err := runner.EmitIssueWithFix(
				r,
				fmt.Sprintf("\"%s\" is a non-ephemeral data source, which means that all (sensitive) attributes are stored in state. Please use ephemeral resource \"%s\" instead.", resourceType, resourceType),
				resource.TypeRange,
				func(f tflint.Fixer) error {
					return f.ReplaceText(resource.TypeRange, "ephemeral")
				},
			); err != nil {
				return fmt.Errorf("failed to call EmitIssueWithFix(): %w", err)
			}
		}
	}

	return nil
}

func GetDataSourceContent(r tflint.Runner, name string, schema *hclext.BodySchema, opts *tflint.GetModuleContentOption) (*hclext.BodyContent, error) {
	body, err := r.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "data", LabelNames: []string{"type", "name"}, Body: schema},
		},
	}, opts)
	if err != nil {
		return nil, err
	}

	content := &hclext.BodyContent{Blocks: []*hclext.Block{}}
	for _, resource := range body.Blocks {
		if resource.Labels[0] != name {
			continue
		}

		content.Blocks = append(content.Blocks, resource)
	}

	return content, nil
}
