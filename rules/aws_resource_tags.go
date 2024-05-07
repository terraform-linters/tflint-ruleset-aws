package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

const (
	defaultTagsBlockName  = "default_tags"
	tagsAttributeName     = "tags"
	tagBlockName          = "tag"
	providerAttributeName = "provider"
)

type AwsResourceTagsRule struct {
	tflint.DefaultRule
}

// awsAutoscalingGroupTag is used by go-cty to evaluate tags in aws_autoscaling_group resources
// The type does not need to be public, but its fields do
// https://github.com/zclconf/go-cty/blob/master/docs/gocty.md#converting-to-and-from-structs
type awsAutoscalingGroupTag struct {
	Key               string `cty:"key"`
	Value             string `cty:"value"`
	PropagateAtLaunch bool   `cty:"propagate_at_launch"`
}

type awsTags map[string]string
type awsProvidersTags map[string]awsTags

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getKnownForValue(value cty.Value) (map[string]string, bool) {
	tags := map[string]string{}

	if !value.CanIterateElements() || !value.IsKnown() {
		return nil, false
	}
	if value.IsNull() {
		return tags, true
	}

	return tags, !value.ForEachElement(func(key, value cty.Value) bool {
		// If any key is unknown or sensitive, return early as any missing tag could be this unknown key.
		if !key.IsKnown() || key.IsNull() || key.IsMarked() {
			return true
		}

		if !value.IsKnown() || value.IsMarked() {
			return true
		}

		// We assume the value of the tag is ALWAYS a string
		tags[key.AsString()] = value.AsString()

		return false
	})
}
