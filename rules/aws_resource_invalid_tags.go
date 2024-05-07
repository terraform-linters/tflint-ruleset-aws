package rules

import (
	"fmt"
	"sort"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/aws"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/tags"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/slices"
)

// AwsResourceInvalidTagsRule checks whether resources are tagged with valid values
type AwsResourceInvalidTagsRule struct {
	tflint.DefaultRule
}

type awsResourceInvalidTagsRuleConfig struct {
	Tags    map[string][]string `hclext:"tags"`
	Exclude []string            `hclext:"exclude,optional"`
}

// NewAwsResourceInvalidTagsRule returns new rules for all resources that support tags
func NewAwsResourceInvalidTagsRule() *AwsResourceInvalidTagsRule {
	return &AwsResourceInvalidTagsRule{}
}

// Name returns the rule name
func (r *AwsResourceInvalidTagsRule) Name() string {
	return "aws_resource_invalid_tags"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsResourceInvalidTagsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsResourceInvalidTagsRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsResourceInvalidTagsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks resources for invalid tags
func (r *AwsResourceInvalidTagsRule) Check(runner tflint.Runner) error {
	config := awsResourceInvalidTagsRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	providerTagsMap, err := r.getProviderLevelTags(runner)

	if err != nil {
		return err
	}

	for _, resourceType := range tags.Resources {
		// Skip this resource if its type is excluded in configuration
		if stringInSlice(resourceType, config.Exclude) {
			continue
		}

		resources, err := runner.GetResourceContent(resourceType, &hclext.BodySchema{
			Attributes: []hclext.AttributeSchema{
				{Name: tagsAttributeName},
				{Name: providerAttributeName},
			},
		}, nil)

		if err != nil {
			return err
		}

		if resources.IsEmpty() {
			continue
		}

		for _, resource := range resources.Blocks {
			providerAlias := "default"

			// Override the provider alias if defined
			if val, ok := resource.Body.Attributes[providerAttributeName]; ok {
				provider, diagnostics := aws.DecodeProviderConfigRef(val.Expr, "provider")
				if diagnostics.HasErrors() {
					logger.Error("error decoding provider: %w", diagnostics)
					return diagnostics
				}
				providerAlias = provider.Alias
			}

			providerAliasTags := providerTagsMap[providerAlias]

			// If the resource has a tags attribute
			if attribute, okResource := resource.Body.Attributes[tagsAttributeName]; okResource {
				logger.Debug(
					"Walk `%s` attribute",
					resource.Labels[0]+"."+resource.Labels[1]+"."+tagsAttributeName,
				)

				err := runner.EvaluateExpr(attribute.Expr, func(val cty.Value) error {
					knownTags, known := getKnownForValue(val)
					if !known {
						logger.Warn("The invalid aws tags rule can only evaluate provided variables, skipping %s.", resource.Labels[0]+"."+resource.Labels[1]+"."+tagsAttributeName)
						return nil
					}

					// merge the known tags with the provider tags to comply with the implementation
					// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/resource-tagging#propagating-tags-to-all-resources
					for providerTagKey, providerTagValue := range providerAliasTags {
						if _, ok := knownTags[providerTagKey]; !ok {
							knownTags[providerTagKey] = providerTagValue
						}
					}

					r.emitIssue(runner, knownTags, config, attribute.Expr.Range())
					return nil
				}, nil)

				if err != nil {
					return err
				}
			} else {
				logger.Debug("Walk `%s` resource", resource.Labels[0]+"."+resource.Labels[1])
				r.emitIssue(runner, providerAliasTags, config, resource.DefRange)
			}
		}
	}

	// Special handling for tags on aws_autoscaling_group resources
	if err := r.checkAwsAutoScalingGroups(runner, config); err != nil {
		return err
	}

	return nil
}

func (r *AwsResourceInvalidTagsRule) getProviderLevelTags(runner tflint.Runner) (awsProvidersTags, error) {
	providerSchema := &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{
				Name:     "alias",
				Required: false,
			},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: defaultTagsBlockName,
				Body: &hclext.BodySchema{Attributes: []hclext.AttributeSchema{{Name: tagsAttributeName}}},
			},
		},
	}

	providerBody, err := runner.GetProviderContent("aws", providerSchema, nil)
	if err != nil {
		return nil, err
	}

	// Get provider default tags
	allProviderTags := make(awsProvidersTags)
	var providerAlias string

	for _, provider := range providerBody.Blocks.OfType(providerAttributeName) {
		// Get the alias attribute, in terraform when there is a single aws provider its called "default"
		providerAttr, ok := provider.Body.Attributes["alias"]
		if !ok {
			providerAlias = "default"
		} else {
			err := runner.EvaluateExpr(providerAttr.Expr, func(alias string) error {
				logger.Debug("Walk `%s` provider", providerAlias)
				providerAlias = alias
				// Init the provider reference even if it doesn't have tags
				allProviderTags[alias] = nil
				return nil
			}, nil)
			if err != nil {
				return nil, err
			}
		}

		for _, block := range provider.Body.Blocks {
			attr, ok := block.Body.Attributes[tagsAttributeName]
			if !ok {
				continue
			}

			err := runner.EvaluateExpr(attr.Expr, func(val cty.Value) error {
				tags, known := getKnownForValue(val)

				if !known {
					logger.Warn("The invalid aws tags rule can only evaluate provided variables, skipping %s.", provider.Labels[0]+"."+providerAlias+"."+defaultTagsBlockName+"."+tagsAttributeName)
					return nil
				}

				logger.Debug("Walk `%s` provider with tags `%v`", providerAlias, tags)
				allProviderTags[providerAlias] = tags
				return nil
			}, nil)

			if err != nil {
				return nil, err
			}
		}
	}
	return allProviderTags, nil
}

// checkAwsAutoScalingGroups handles the special case for tags on AutoScaling Groups
// See: https://github.com/terraform-providers/terraform-provider-aws/blob/master/aws/autoscaling_tags.go
func (r *AwsResourceInvalidTagsRule) checkAwsAutoScalingGroups(runner tflint.Runner, config awsResourceInvalidTagsRuleConfig) error {
	resourceType := "aws_autoscaling_group"

	// Skip autoscaling group check if its type is excluded in configuration
	if stringInSlice(resourceType, config.Exclude) {
		return nil
	}

	resources, err := runner.GetResourceContent(resourceType, &hclext.BodySchema{}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		asgTagBlockTags, tagBlockLocation, err := r.checkAwsAutoScalingGroupsTag(runner, resource)
		if err != nil {
			return err
		}

		asgTagsAttributeTags, tagsAttributeLocation, err := r.checkAwsAutoScalingGroupsTags(runner, resource)
		if err != nil {
			return err
		}

		switch {
		case len(asgTagBlockTags) > 0 && len(asgTagsAttributeTags) > 0:
			runner.EmitIssue(r, "Only tag block or tags attribute may be present, but found both", resource.DefRange)
		case len(asgTagBlockTags) == 0 && len(asgTagsAttributeTags) == 0:
			r.emitIssue(runner, map[string]string{}, config, resource.DefRange)
		case len(asgTagBlockTags) > 0 && len(asgTagsAttributeTags) == 0:
			tags := asgTagBlockTags
			location := tagBlockLocation
			r.emitIssue(runner, tags, config, location)
		case len(asgTagBlockTags) == 0 && len(asgTagsAttributeTags) > 0:
			tags := asgTagsAttributeTags
			location := tagsAttributeLocation
			r.emitIssue(runner, tags, config, location)
		}
	}

	return nil
}

// checkAwsAutoScalingGroupsTag checks tag{} blocks on aws_autoscaling_group resources
func (r *AwsResourceInvalidTagsRule) checkAwsAutoScalingGroupsTag(runner tflint.Runner, resourceBlock *hclext.Block) (map[string]string, hcl.Range, error) {
	tags := map[string]string{}

	resources, err := runner.GetResourceContent("aws_autoscaling_group", &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: tagBlockName,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "key"},
						{Name: "value"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return tags, hcl.Range{}, err
	}

	for _, resource := range resources.Blocks {
		if resource.Labels[0] != resourceBlock.Labels[0] {
			continue
		}

		for _, tag := range resource.Body.Blocks {
			keyAttribute, keyExists := tag.Body.Attributes["key"]
			if !keyExists {
				return tags, hcl.Range{}, fmt.Errorf(`Did not find expected field "key" in aws_autoscaling_group "%s" starting at line %d`, resource.Labels[0], resource.DefRange.Start.Line)
			}

			valueAttribute, valueExists := tag.Body.Attributes["value"]
			if !valueExists {
				return tags, hcl.Range{}, fmt.Errorf(`Did not find expected field "value" in aws_autoscaling_group "%s" starting at line %d`, resource.Labels[0], resource.DefRange.Start.Line)
			}

			err := runner.EvaluateExpr(keyAttribute.Expr, func(key string) error {
				return runner.EvaluateExpr(valueAttribute.Expr, func(value string) error {
					tags[key] = value
					return nil
				}, nil)
			}, nil)
			if err != nil {
				return tags, hcl.Range{}, err
			}
		}
	}

	return tags, resourceBlock.DefRange, nil
}

// checkAwsAutoScalingGroupsTag checks the tags attribute on aws_autoscaling_group resources
func (r *AwsResourceInvalidTagsRule) checkAwsAutoScalingGroupsTags(runner tflint.Runner, resourceBlock *hclext.Block) (map[string]string, hcl.Range, error) {
	tags := map[string]string{}

	resources, err := runner.GetResourceContent("aws_autoscaling_group", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: tagsAttributeName},
		},
	}, nil)
	if err != nil {
		return tags, hcl.Range{}, err
	}

	for _, resource := range resources.Blocks {
		if resource.Labels[0] != resourceBlock.Labels[0] {
			continue
		}

		attribute, ok := resource.Body.Attributes[tagsAttributeName]
		if ok {
			wantType := cty.List(cty.Object(map[string]cty.Type{
				"key":                 cty.String,
				"value":               cty.String,
				"propagate_at_launch": cty.Bool,
			}))
			err := runner.EvaluateExpr(attribute.Expr, func(asgTags []awsAutoscalingGroupTag) error {
				for _, tag := range asgTags {
					tags[tag.Key] = tag.Value
				}
				return nil
			}, &tflint.EvaluateExprOption{WantType: &wantType})
			if err != nil {
				return tags, attribute.Expr.Range(), err
			}
			return tags, attribute.Expr.Range(), nil
		}
	}

	return tags, resourceBlock.DefRange, nil
}

func (r *AwsResourceInvalidTagsRule) emitIssue(runner tflint.Runner, tags map[string]string, config awsResourceInvalidTagsRuleConfig, location hcl.Range) {
	// sort the tag names for deterministic output
	// only evaluate the given tags on the resource NOT the configured tags
	// the checking of tag presence should use the `aws_resource_missing_tags` lint rule
	tagsToMatch := sort.StringSlice{}
	for tagName := range tags {
		tagsToMatch = append(tagsToMatch, tagName)
	}
	tagsToMatch.Sort()

	str := ""
	for _, tagName := range tagsToMatch {
		allowedValues, ok := config.Tags[tagName]
		// if the tag has a rule configuration then check
		if ok {
			valueProvided := tags[tagName]
			if !slices.Contains(allowedValues, valueProvided) {
				str = str + fmt.Sprintf("Received '%s' for tag '%s', expected one of '%s'. ", valueProvided, tagName, strings.Join(allowedValues, ","))
			}
		}
	}

	if len(str) > 0 {
		runner.EmitIssue(r, strings.TrimSpace(str), location)
	}
}
