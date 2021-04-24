package rules

import (
	"fmt"
	"log"
	"sort"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-aws/project"
	"github.com/terraform-linters/tflint-ruleset-aws/rules/tags"
	"github.com/zclconf/go-cty/cty"
)

// AwsResourceMissingTagsRule checks whether resources are tagged correctly
type AwsResourceMissingTagsRule struct{}

type awsResourceTagsRuleConfig struct {
	Tags    []string `hcl:"tags"`
	Exclude []string `hcl:"exclude,optional"`

	Remain hcl.Body `hcl:",remain"`
}

const (
	tagsAttributeName = "tags"
	tagBlockName      = "tag"
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
func (r *AwsResourceMissingTagsRule) Severity() string {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *AwsResourceMissingTagsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks resources for missing tags
func (r *AwsResourceMissingTagsRule) Check(runner tflint.Runner) error {
	config := awsResourceTagsRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}

	for _, resourceType := range tags.Resources {
		// Skip this resource if its type is excluded in configuration
		if stringInSlice(resourceType, config.Exclude) {
			continue
		}

		// Special handling for tags on aws_autoscaling_group resources
		if resourceType == "aws_autoscaling_group" {
			err := r.checkAwsAutoScalingGroups(runner, config)
			err = runner.EnsureNoError(err, func() error {
				return nil
			})
			if err != nil {
				return err
			}
			continue
		}

		err := runner.WalkResources(resourceType, func(resource *configs.Resource) error {
			body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
				Attributes: []hcl.AttributeSchema{
					{
						Name: tagsAttributeName,
					},
				},
			})
			if diags.HasErrors() {
				return diags
			}

			if attribute, ok := body.Attributes[tagsAttributeName]; ok {
				log.Printf("[DEBUG] Walk `%s` attribute", resource.Type+"."+resource.Name+"."+tagsAttributeName)
				resourceTags := make(map[string]string)
				wantType := cty.Map(cty.String)
				err := runner.EvaluateExpr(attribute.Expr, &resourceTags, &wantType)
				err = runner.EnsureNoError(err, func() error {
					r.emitIssueOnExpr(runner, resourceTags, config, attribute.Expr)
					return nil
				})
				if err != nil {
					return err
				}
			} else {
				log.Printf("[DEBUG] Walk `%s` resource", resource.Type+"."+resource.Name)
				r.emitIssue(runner, map[string]string{}, config, resource.DeclRange)
			}

			return nil
		})

		if err != nil {
			return err
		}
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

	return runner.WalkResources(resourceType, func(resource *configs.Resource) error {
		asgTagBlockTags, tagBlockLocation, err := r.checkAwsAutoScalingGroupsTag(runner, config, resource)
		if err != nil {
			return err
		}

		asgTagsAttributeTags, tagsAttributeLocation, err := r.checkAwsAutoScalingGroupsTags(runner, config, resource)
		if err != nil {
			return err
		}

		var location hcl.Range
		tags := make(map[string]string)
		switch {
		case len(asgTagBlockTags) > 0 && len(asgTagsAttributeTags) > 0:
			issue := fmt.Sprintf("Only tag block or tags attribute may be present, but found both")
			runner.EmitIssue(r, issue, resource.DeclRange)
			return nil
		case len(asgTagBlockTags) == 0 && len(asgTagsAttributeTags) == 0:
			r.emitIssue(runner, map[string]string{}, config, resource.DeclRange)
			return nil
		case len(asgTagBlockTags) > 0 && len(asgTagsAttributeTags) == 0:
			tags = asgTagBlockTags
			location = tagBlockLocation
		case len(asgTagBlockTags) == 0 && len(asgTagsAttributeTags) > 0:
			tags = asgTagsAttributeTags
			location = tagsAttributeLocation
		}

		return runner.EnsureNoError(err, func() error {
			r.emitIssue(runner, tags, config, location)
			return nil
		})
	})
}

// checkAwsAutoScalingGroupsTag checks tag{} blocks on aws_autoscaling_group resources
func (r *AwsResourceMissingTagsRule) checkAwsAutoScalingGroupsTag(runner tflint.Runner, config awsResourceTagsRuleConfig, resource *configs.Resource) (map[string]string, hcl.Range, error) {
	tags := make(map[string]string)

	err := runner.WalkResourceBlocks(resource.Type, tagBlockName, func(tagBlock *hcl.Block) error {
		attributes, diags := tagBlock.Body.JustAttributes()
		if diags.HasErrors() {
			return diags
		}

		if _, ok := attributes["key"]; !ok {
			err := &tflint.Error{
				Code:  tflint.UnevaluableError,
				Level: tflint.WarningLevel,
				Message: fmt.Sprintf("Did not find expected field \"key\" in aws_autoscaling_group \"%s\" starting at line %d",
					resource.Name,
					resource.DeclRange.Start.Line,
				),
			}
			return err
		}

		var key string
		err := runner.EvaluateExpr(attributes["key"].Expr, &key, nil)
		if err != nil {
			return err
		}
		tags[key] = ""

		return nil
	})
	if err != nil {
		return tags, (hcl.Range{}), err
	}

	return tags, resource.DeclRange, nil
}

// checkAwsAutoScalingGroupsTag checks the tags attribute on aws_autoscaling_group resources
func (r *AwsResourceMissingTagsRule) checkAwsAutoScalingGroupsTags(runner tflint.Runner, config awsResourceTagsRuleConfig, resource *configs.Resource) (map[string]string, hcl.Range, error) {
	tags := make(map[string]string)
	body, _, diags := resource.Config.PartialContent(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: tagsAttributeName,
			},
		},
	})
	if diags.HasErrors() {
		return tags, (hcl.Range{}), diags
	}

	attribute, ok := body.Attributes[tagsAttributeName]
	if ok {
		wantType := cty.List(cty.Object(map[string]cty.Type{
			"key":                 cty.String,
			"value":               cty.String,
			"propagate_at_launch": cty.Bool,
		}))
		var asgTags []awsAutoscalingGroupTag
		err := runner.EvaluateExpr(attribute.Expr, &asgTags, &wantType)
		if err != nil {
			return tags, attribute.Expr.Range(), err
		}
		for _, tag := range asgTags {
			tags[tag.Key] = tag.Value
		}
		return tags, attribute.Expr.Range(), nil
	}

	return tags, resource.DeclRange, nil
}

func (r *AwsResourceMissingTagsRule) emitIssue(runner tflint.Runner, tags map[string]string, config awsResourceTagsRuleConfig, location hcl.Range) {
	var missing []string
	for _, tag := range config.Tags {
		if _, ok := tags[tag]; !ok {
			missing = append(missing, fmt.Sprintf("\"%s\"", tag))
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		wanted := strings.Join(missing, ", ")
		issue := fmt.Sprintf("The resource is missing the following tags: %s.", wanted)
		runner.EmitIssue(r, issue, location)
	}
}

func (r *AwsResourceMissingTagsRule) emitIssueOnExpr(runner tflint.Runner, tags map[string]string, config awsResourceTagsRuleConfig, expr hcl.Expression) {
	var missing []string
	for _, tag := range config.Tags {
		if _, ok := tags[tag]; !ok {
			missing = append(missing, fmt.Sprintf("\"%s\"", tag))
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		wanted := strings.Join(missing, ", ")
		issue := fmt.Sprintf("The resource is missing the following tags: %s.", wanted)
		runner.EmitIssueOnExpr(r, issue, expr)
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
