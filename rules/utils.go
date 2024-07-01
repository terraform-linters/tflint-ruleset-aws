package rules

import "github.com/zclconf/go-cty/cty"

var validElastiCacheNodeTypes = map[string]bool{
	// https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/CacheNodes.SupportedTypes.html
	"cache.t2.micro":      true,
	"cache.t2.small":      true,
	"cache.t2.medium":     true,
	"cache.t3.micro":      true,
	"cache.t3.small":      true,
	"cache.t3.medium":     true,
	"cache.t4g.micro":     true,
	"cache.t4g.small":     true,
	"cache.t4g.medium":    true,
	"cache.m3.medium":     true,
	"cache.m3.large":      true,
	"cache.m3.xlarge":     true,
	"cache.m3.2xlarge":    true,
	"cache.m4.large":      true,
	"cache.m4.xlarge":     true,
	"cache.m4.2xlarge":    true,
	"cache.m4.4xlarge":    true,
	"cache.m4.10xlarge":   true,
	"cache.m5.large":      true,
	"cache.m5.xlarge":     true,
	"cache.m5.2xlarge":    true,
	"cache.m5.4xlarge":    true,
	"cache.m5.12xlarge":   true,
	"cache.m5.24xlarge":   true,
	"cache.m6g.large":     true,
	"cache.m6g.xlarge":    true,
	"cache.m6g.2xlarge":   true,
	"cache.m6g.4xlarge":   true,
	"cache.m6g.8xlarge":   true,
	"cache.m6g.12xlarge":  true,
	"cache.m6g.16xlarge":  true,
	"cache.m7g.large":     true,
	"cache.m7g.xlarge":    true,
	"cache.m7g.2xlarge":   true,
	"cache.m7g.4xlarge":   true,
	"cache.m7g.8xlarge":   true,
	"cache.m7g.12xlarge":  true,
	"cache.m7g.16xlarge":  true,
	"cache.r3.large":      true,
	"cache.r3.xlarge":     true,
	"cache.r3.2xlarge":    true,
	"cache.r3.4xlarge":    true,
	"cache.r3.8xlarge":    true,
	"cache.r4.large":      true,
	"cache.r4.xlarge":     true,
	"cache.r4.2xlarge":    true,
	"cache.r4.4xlarge":    true,
	"cache.r4.8xlarge":    true,
	"cache.r4.16xlarge":   true,
	"cache.r5.large":      true,
	"cache.r5.xlarge":     true,
	"cache.r5.2xlarge":    true,
	"cache.r5.4xlarge":    true,
	"cache.r5.12xlarge":   true,
	"cache.r5.24xlarge":   true,
	"cache.r6g.large":     true,
	"cache.r6g.xlarge":    true,
	"cache.r6g.2xlarge":   true,
	"cache.r6g.4xlarge":   true,
	"cache.r6g.8xlarge":   true,
	"cache.r6g.12xlarge":  true,
	"cache.r6g.16xlarge":  true,
	"cache.r6gd.xlarge":   true,
	"cache.r6gd.2xlarge":  true,
	"cache.r6gd.4xlarge":  true,
	"cache.r6gd.8xlarge":  true,
	"cache.r6gd.12xlarge": true,
	"cache.r6gd.16xlarge": true,
	"cache.r7g.large":     true,
	"cache.r7g.xlarge":    true,
	"cache.r7g.2xlarge":   true,
	"cache.r7g.4xlarge":   true,
	"cache.r7g.8xlarge":   true,
	"cache.r7g.12xlarge":  true,
	"cache.r7g.16xlarge":  true,
	"cache.c7gn.large":    true,
	"cache.c7gn.xlarge":   true,
	"cache.m1.small":      true,
	"cache.m1.medium":     true,
	"cache.m1.large":      true,
	"cache.m1.xlarge":     true,
	"cache.m2.xlarge":     true,
	"cache.m2.2xlarge":    true,
	"cache.m2.4xlarge":    true,
	"cache.c1.xlarge":     true,
	"cache.t1.micro":      true,
}

var previousElastiCacheNodeTypes = map[string]bool{
	// https://aws.amazon.com/elasticache/previous-generation/?nc1=h_ls
	"c1": true,
	"m1": true,
	"m2": true,
	"m3": true,
	"r3": true,
	"t1": true,
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
