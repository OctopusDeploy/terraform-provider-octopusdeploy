package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTag(tfTag map[string]interface{}) octopusdeploy.Tag {
	tag := octopusdeploy.Tag{
		CanonicalTagName: tfTag["canonical_tag_name"].(string),
		Color:            tfTag["color"].(string),
		Description:      tfTag["description"].(string),
		Name:             tfTag["name"].(string),
		SortOrder:        tfTag["sort_order"].(int),
	}

	return tag
}

func getTagSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"canonical_tag_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"color": {
			Required: true,
			Type:     schema.TypeString,
		},
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"sort_order": {
			Optional: true,
			Type:     schema.TypeInt,
		},
	}
}
