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

// AwsResourceMissingTagsRule checks whether resources are tagged correctly
type AwsResourceMissingTagsRule struct {
	tflint.DefaultRule
}

type awsResourceTagsRuleConfig struct {
	Tags    []string `hclext:"tags"`
	Exclude []string `hclext:"exclude,optional"`
}

const (
	defaultTagsBlockName  = "default_tags"
	tagsAttributeName     = "tags"
	tagBlockName          = "tag"
	providerAttributeName = "provider"
)

// NewAwsResourceMissingTagsRule returns new rules for all resources that support tags
func NewAwsResourceMissingTagsRule() *AwsResourceMissingTagsRule {
	return &AwsResourceMissingTagsRule{}
}

// Name returns the rule name
func (r *AwsResourceMissingTagsRule) Name() string {
	return "aws_resource_missing_tags"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsResourceMissingTagsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsResourceMissingTagsRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsResourceMissingTagsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

func (r *AwsResourceMissingTagsRule) getProviderLevelTags(runner tflint.Runner) (map[string][]string, error) {
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
	allProviderTags := make(map[string][]string)
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
			var providerTags []string
			attr, ok := block.Body.Attributes[tagsAttributeName]
			if !ok {
				continue
			}

			err := runner.EvaluateExpr(attr.Expr, func(val cty.Value) error {
				keys, known := getKeysForValue(val)

				if !known {
					logger.Warn("The missing aws tags rule can only evaluate provided variables, skipping %s.", provider.Labels[0]+"."+providerAlias+"."+defaultTagsBlockName+"."+tagsAttributeName)
					return nil
				}

				logger.Debug("Walk `%s` provider with tags `%v`", providerAlias, keys)
				providerTags = keys
				return nil
			}, nil)

			if err != nil {
				return nil, err
			}

			allProviderTags[providerAlias] = providerTags
		}
	}
	return allProviderTags, nil
}

// Check checks resources for missing tags
func (r *AwsResourceMissingTagsRule) Check(runner tflint.Runner) error {
	config := awsResourceTagsRuleConfig{}
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

			// If the resource has a tags attribute
			if attribute, okResource := resource.Body.Attributes[tagsAttributeName]; okResource {
				logger.Debug(
					"Walk `%s` attribute",
					resource.Labels[0]+"."+resource.Labels[1]+"."+tagsAttributeName,
				)

				err := runner.EvaluateExpr(attribute.Expr, func(val cty.Value) error {
					keys, known := getKeysForValue(val)
					if !known {
						logger.Warn("The missing aws tags rule can only evaluate provided variables, skipping %s.", resource.Labels[0]+"."+resource.Labels[1]+"."+tagsAttributeName)
						return nil
					}

					r.emitIssue(runner, append(providerTagsMap[providerAlias], keys...), config, attribute.Expr.Range())
					return nil
				}, nil)

				if err != nil {
					return err
				}
			} else {
				logger.Debug("Walk `%s` resource", resource.Labels[0]+"."+resource.Labels[1])
				r.emitIssue(runner, providerTagsMap[providerAlias], config, resource.DefRange)
			}
		}
	}

	// Special handling for tags on aws_autoscaling_group resources
	if err := r.checkAwsAutoScalingGroups(runner, config); err != nil {
		return err
	}

	return nil
}

// awsAutoscalingGroupTag is used by go-cty to evaluate tags in aws_autoscaling_group resources
// The type does not need to be public, but its fields do
// https://github.com/zclconf/go-cty/blob/master/docs/gocty.md#converting-to-and-from-structs
type awsAutoscalingGroupTag struct {
	Key               string `cty:"key"`
	Value             string `cty:"value"`
	PropagateAtLaunch bool   `cty:"propagate_at_launch"`
}

// checkAwsAutoScalingGroups handles the special case for tags on AutoScaling Groups
// See: https://github.com/terraform-providers/terraform-provider-aws/blob/master/aws/autoscaling_tags.go
func (r *AwsResourceMissingTagsRule) checkAwsAutoScalingGroups(runner tflint.Runner, config awsResourceTagsRuleConfig) error {
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
		asgTagBlockTags, tagBlockLocation, err := r.checkAwsAutoScalingGroupsTag(runner, config, resource)
		if err != nil {
			return err
		}

		asgTagsAttributeTags, tagsAttributeLocation, err := r.checkAwsAutoScalingGroupsTags(runner, config, resource)
		if err != nil {
			return err
		}

		switch {
		case len(asgTagBlockTags) > 0 && len(asgTagsAttributeTags) > 0:
			runner.EmitIssue(r, "Only tag block or tags attribute may be present, but found both", resource.DefRange)
		case len(asgTagBlockTags) == 0 && len(asgTagsAttributeTags) == 0:
			r.emitIssue(runner, []string{}, config, resource.DefRange)
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
func (r *AwsResourceMissingTagsRule) checkAwsAutoScalingGroupsTag(runner tflint.Runner, config awsResourceTagsRuleConfig, resourceBlock *hclext.Block) ([]string, hcl.Range, error) {
	tags := make([]string, 0)

	resources, err := runner.GetResourceContent("aws_autoscaling_group", &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: tagBlockName,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "key"},
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
			attribute, exists := tag.Body.Attributes["key"]
			if !exists {
				return tags, hcl.Range{}, fmt.Errorf(`Did not find expected field "key" in aws_autoscaling_group "%s" starting at line %d`, resource.Labels[0], resource.DefRange.Start.Line)
			}

			err := runner.EvaluateExpr(attribute.Expr, func(key string) error {
				tags = append(tags, key)
				return nil
			}, nil)
			if err != nil {
				return tags, hcl.Range{}, err
			}
		}
	}

	return tags, resourceBlock.DefRange, nil
}

// checkAwsAutoScalingGroupsTag checks the tags attribute on aws_autoscaling_group resources
func (r *AwsResourceMissingTagsRule) checkAwsAutoScalingGroupsTags(runner tflint.Runner, config awsResourceTagsRuleConfig, resourceBlock *hclext.Block) ([]string, hcl.Range, error) {
	tags := make([]string, 0)

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
					tags = append(tags, tag.Key)
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

func (r *AwsResourceMissingTagsRule) emitIssue(runner tflint.Runner, tags []string, config awsResourceTagsRuleConfig, location hcl.Range) {
	var missing []string
	for _, tag := range config.Tags {
		if !slices.Contains(tags, tag) {
			missing = append(missing, fmt.Sprintf("%q", tag))
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		wanted := strings.Join(missing, ", ")
		issue := fmt.Sprintf("The resource is missing the following tags: %s.", wanted)
		runner.EmitIssue(r, issue, location)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// getKeysForValue returns a list of keys from a cty.Value, which is assumed to be a map (or unknown).
// It returns a boolean indicating whether the keys were known.
// If _any_ key is unknown, the entire value is considered unknown, since we can't know if a required tag might be matched by the unknown key.
// Values are entirely ignored and can be unknown.
func getKeysForValue(value cty.Value) (keys []string, known bool) {
	if !value.CanIterateElements() || !value.IsKnown() {
		return nil, false
	}
	if value.IsNull() {
		return keys, true
	}

	return keys, !value.ForEachElement(func(key, _ cty.Value) bool {
		// If any key is unknown or sensitive, return early as any missing tag could be this unknown key.
		if !key.IsKnown() || key.IsNull() || key.IsMarked() {
			return true
		}

		keys = append(keys, key.AsString())

		return false
	})
}
