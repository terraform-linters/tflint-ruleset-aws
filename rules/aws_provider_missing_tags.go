package rules

import (
	"fmt"
	"sort"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/slices"
)

// AwsProviderMissingTagsRule checks whether providers are tagged correctly
type AwsProviderMissingTagsRule struct {
	tflint.DefaultRule
}

type awsProviderTagsRuleConfig struct {
	Tags []string `hclext:"tags"`
}

const (
	providerDefaultTagsBlockName = "default_tags"
	providerTagsAttributeName    = "tags"
)

// NewAwsProviderMissingTagsRule returns new rules for all providers that support tags
func NewAwsProviderMissingTagsRule() *AwsProviderMissingTagsRule {
	return &AwsProviderMissingTagsRule{}
}

// Name returns the rule name
func (r *AwsProviderMissingTagsRule) Name() string {
	return "aws_provider_missing_tags"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsProviderMissingTagsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsProviderMissingTagsRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsProviderMissingTagsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks providers for missing tags
func (r *AwsProviderMissingTagsRule) Check(runner tflint.Runner) error {
	config := awsProviderTagsRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	providerSchema := &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{
				Name:     "alias",
				Required: false,
			},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: providerDefaultTagsBlockName,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name: providerTagsAttributeName,
						},
					},
				},
			},
		},
	}

	providerBody, err := runner.GetProviderContent("aws", providerSchema, nil)
	if err != nil {
		return nil
	}

	// Get provider default tags
	var providerAlias string
	for _, provider := range providerBody.Blocks.OfType("provider") {
		// Get the alias attribute, in terraform when there is a single aws provider its called "default"
		providerAttr, ok := provider.Body.Attributes["alias"]
		if !ok {
			providerAlias = "default"
		} else {
			err := runner.EvaluateExpr(providerAttr.Expr, func(alias string) error {
				providerAlias = alias
				return nil
			}, nil)
			if err != nil {
				return nil
			}
		}
		logger.Debug("Walk `%s` provider", providerAlias)

		if len(provider.Body.Blocks) == 0 {
			if err := runner.EmitIssue(r, fmt.Sprintf("The provider `%s` is missing the `default_tags` block", providerAlias), provider.DefRange); err != nil {
				return err
			}
			continue
		}

		for _, block := range provider.Body.Blocks {
			var providerTags []string
			attr, ok := block.Body.Attributes[providerTagsAttributeName]
			if !ok {
				r.emitIssue(runner, providerTags, config, provider.DefRange)
				continue
			}

			err := runner.EvaluateExpr(attr.Expr, func(val cty.Value) error {
				keys, known := getKeysForValue(val)

				if !known {
					logger.Warn("The missing aws tags rule can only evaluate provided variables, skipping %s.", provider.Labels[0]+"."+providerAlias+"."+providerDefaultTagsBlockName+"."+providerTagsAttributeName)
					return nil
				}

				logger.Debug("Walk `%s` provider with tags `%v`", providerAlias, keys)
				providerTags = keys
				return nil
			}, nil)

			if err != nil {
				return nil
			}

			// Check tags
			r.emitIssue(runner, providerTags, config, provider.DefRange)
		}
	}

	return nil
}

func (r *AwsProviderMissingTagsRule) emitIssue(runner tflint.Runner, tags []string, config awsProviderTagsRuleConfig, location hcl.Range) {
	var missing []string
	for _, tag := range config.Tags {
		if !slices.Contains(tags, tag) {
			missing = append(missing, fmt.Sprintf("%q", tag))
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		wanted := strings.Join(missing, ", ")
		issue := fmt.Sprintf("The provider is missing the following tags: %s.", wanted)
		runner.EmitIssue(r, issue, location)
	}
}
