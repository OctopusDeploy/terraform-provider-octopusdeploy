package octopusdeploy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getProjectDeploymentTargetTriggerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": getNameSchema(true),
		"space_id": {
			Description: "The space ID associated with the project to attach the trigger. Defaults to the space the provider is configured for.",
			Optional:    true,
			// Computed so that omitting the attribute takes the value the
			// server reports rather than producing a diff on every plan, and
			// ForceNew because a trigger cannot be moved between spaces.
			Computed:         true,
			ForceNew:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"project_id": {
			Description: "The ID of the project to attach the trigger.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"should_redeploy": {
			Default:     false,
			Description: "Enable to re-deploy to the deployment targets even if they are already up-to-date with the current deployment.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"event_groups": {
			Description: "Apply event group filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"event_categories": {
			Description: "Apply event category filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"roles": {
			Description: "Apply event role filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"environment_ids": {
			Description: "Apply environment id filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}
