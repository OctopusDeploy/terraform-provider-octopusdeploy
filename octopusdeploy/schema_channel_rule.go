package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandChannelRule(channelRule map[string]interface{}) channels.ChannelRule {
	return channels.ChannelRule{
		ActionPackages:     expandDeploymentActionPackages(channelRule["action_package"]),
		ID:                 channelRule["id"].(string),
		Tag:                channelRule["tag"].(string),
		VersionRange:       channelRule["version_range"].(string),
		VersioningStrategy: channelRule["versioning_strategy"].(string),
		VersionTagRegex:    channelRule["version_tag_regex"].(string),
	}
}

func flattenChannelRules(channelRules []channels.ChannelRule) []map[string]interface{} {
	var flattenedRules = make([]map[string]interface{}, len(channelRules))
	for key, channelRule := range channelRules {
		flattenedRules[key] = map[string]interface{}{
			"action_package":      flattenDeploymentActionPackages(channelRule.ActionPackages),
			"id":                  channelRule.ID,
			"tag":                 channelRule.Tag,
			"version_range":       channelRule.VersionRange,
			"versioning_strategy": channelRule.VersioningStrategy,
			"version_tag_regex":   channelRule.VersionTagRegex,
		}
	}

	return flattenedRules
}

func getChannelRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action_package": {
			Elem:     &schema.Resource{Schema: getDeploymentActionPackageSchema()},
			Required: true,
			Type:     schema.TypeList,
		},
		"id": getIDSchema(),
		"tag": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"version_range": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"versioning_strategy": {
			Description: "The ordering strategy used to determine the latest package version. Valid values are `\"SemVer\"` (default) and `\"MostRecentlyPublished\"`. When `MostRecentlyPublished`, the channel ranks candidate package versions by publish date rather than by Semantic Versioning comparison; use this with non-SemVer schemes such as date-stamped or feature-branch tags. Requires the `non-semver-ordering` feature toggle on the Octopus instance.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"version_tag_regex": {
			Description: "A regular expression matched against the full package version string. Used in place of `version_range` and `tag` filtering when `versioning_strategy` is `\"MostRecentlyPublished\"`.",
			Optional:    true,
			Type:        schema.TypeString,
		},
	}
}
