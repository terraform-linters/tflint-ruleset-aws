package rules

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/zclconf/go-cty/cty"
)

// AwsProviderMissingDefaultTagsRule checks whether providers are tagged correctly
type AwsProviderMissingDefaultTagsRule struct {
	tflint.DefaultRule
}

type awsProviderTagsRuleConfig struct {
	Tags []string `hclext:"tags"`
}

const (
	providerDefaultTagsBlockName = "default_tags"
	providerTagsAttributeName    = "tags"
)

// NewAwsProviderMissingDefaultTagsRule returns new rules for all providers that support tags
func NewAwsProviderMissingDefaultTagsRule() *AwsProviderMissingDefaultTagsRule {
	return &AwsProviderMissingDefaultTagsRule{}
}

// Name returns the rule name
func (r *AwsProviderMissingDefaultTagsRule) Name() string {
	return "aws_provider_missing_default_tags"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsProviderMissingDefaultTagsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsProviderMissingDefaultTagsRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsProviderMissingDefaultTagsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks providers for missing tags
func (r *AwsProviderMissingDefaultTagsRule) Check(runner tflint.Runner) error {
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
				return err
			}
		}
		logger.Debug("Walk `%s` provider", providerAlias)

		if len(provider.Body.Blocks) == 0 {
			var issue string
			if providerAlias == "default" {
				issue = "The aws provider is missing the `default_tags` block"
			} else {
				issue = fmt.Sprintf("The aws provider with alias `%s` is missing the `default_tags` block", providerAlias)
			}
			if err := runner.EmitIssue(r, issue, provider.DefRange); err != nil {
				return err
			}
			continue
		}

		for _, block := range provider.Body.Blocks {
			var providerTags []string
			attr, ok := block.Body.Attributes[providerTagsAttributeName]
			if !ok {
				if err := r.emitIssue(runner, providerTags, config, provider.DefRange); err != nil {
					return err
				}
				continue
			}

			err := runner.EvaluateExpr(attr.Expr, func(val cty.Value) error {
				keys, known := getKeysForValue(val)

				if !known {
					return tflint.ErrUnknownValue
				}

				logger.Debug("Walk `%s` provider with tags `%v`", providerAlias, keys)
				providerTags = keys
				return nil
			}, nil)

			if err != nil {
				logger.Warn("Could not evaluate tags, skipping %s.%s.%s.%s: %s", provider.Labels[0], providerAlias, providerDefaultTagsBlockName, providerTagsAttributeName, err)
				continue
			}

			// Check tags
			if err := r.emitIssue(runner, providerTags, config, attr.Range); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *AwsProviderMissingDefaultTagsRule) emitIssue(runner tflint.Runner, tags []string, config awsProviderTagsRuleConfig, location hcl.Range) error {
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
		if err := runner.EmitIssue(r, issue, location); err != nil {
			return err
		}
	}

	return nil
}
