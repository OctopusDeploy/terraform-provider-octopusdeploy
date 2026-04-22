package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func flattenProjectWebhookTrigger(projectWebhookTrigger *triggers.ProjectTrigger) map[string]interface{} {
	flattenedProjectWebhookTrigger := map[string]interface{}{}
	flattenedProjectWebhookTrigger["id"] = projectWebhookTrigger.GetID()
	flattenedProjectWebhookTrigger["name"] = projectWebhookTrigger.Name
	flattenedProjectWebhookTrigger["description"] = projectWebhookTrigger.Description
	flattenedProjectWebhookTrigger["project_id"] = projectWebhookTrigger.ProjectID
	flattenedProjectWebhookTrigger["space_id"] = projectWebhookTrigger.SpaceID
	flattenedProjectWebhookTrigger["is_disabled"] = projectWebhookTrigger.IsDisabled

	runRunbookAction := projectWebhookTrigger.Action.(*actions.RunRunbookAction)
	flattenedProjectWebhookTrigger["tenant_ids"] = runRunbookAction.Tenants
	flattenedProjectWebhookTrigger["run_runbook_action"] = []map[string]interface{}{
		{
			"runbook_id":             runRunbookAction.Runbook,
			"target_environment_ids": runRunbookAction.Environments,
		},
	}
	webhookFilter := projectWebhookTrigger.Filter.(*filters.WebhookTriggerFilter)
	secret := ""
	if webhookFilter.Secret.NewValue != nil {
		secret = *webhookFilter.Secret.NewValue
	}
	flattenedProjectWebhookTrigger["webhook_trigger_filter"] = []map[string]interface{}{
		{
			"webhook_id": webhookFilter.WebhookId,
			"secret":     secret,
		},
	}

	return flattenedProjectWebhookTrigger
}

func expandProjectWebhookTrigger(projectWebhookTrigger *schema.ResourceData, project *projects.Project) (*triggers.ProjectTrigger, error) {
	name := projectWebhookTrigger.Get("name").(string)
	description := projectWebhookTrigger.Get("description").(string)
	isDisabled := projectWebhookTrigger.Get("is_disabled").(bool)

	var action actions.ITriggerAction = nil
	var filter filters.ITriggerFilter = nil

	if attributes, ok := projectWebhookTrigger.GetOk("run_runbook_action"); ok {
		runRunbookActionList := attributes.(*schema.Set).List()
		runRunbookActionMap := runRunbookActionList[0].(map[string]interface{})
		deploymentAction := actions.NewRunRunbookAction()

		deploymentAction.Runbook = runRunbookActionMap["runbook_id"].(string)
		deploymentAction.Environments = expandArray(runRunbookActionMap["target_environment_ids"].([]interface{}))
		deploymentAction.Tenants = expandArray(projectWebhookTrigger.Get("tenant_ids").([]interface{}))
		action = deploymentAction
	}

	if attributes, ok := projectWebhookTrigger.GetOk("webhook_trigger_filter"); ok {
		webhookFilterList := attributes.(*schema.Set).List()
		webhookFilterMap := webhookFilterList[0].(map[string]interface{})
		secret := *core.NewSensitiveValue(webhookFilterMap["secret"].(string))

		webhookFilter := filters.NewWebhookTriggerFilter(secret)
		webhookFilter.WebhookId = webhookFilterMap["webhook_id"].(string)
		filter = webhookFilter

	}
	projectTriggerToCreate := triggers.NewProjectTrigger(name, description, isDisabled, project, action, filter)
	projectTriggerToCreate.Description = description

	return projectTriggerToCreate, nil
}

func getProjectWebhookTriggerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": getNameSchema(true),
		"description": {
			Description: "A description of the trigger.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"project_id": {
			Description:      "The ID of the project to attach the trigger.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			ForceNew:         true,
		},
		"space_id": {
			Required:         true,
			Description:      "The space ID where this trigger's project exists.",
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			ForceNew:         true,
		},
		"run_runbook_action": {
			Description: "Configuration for running a runbook. Can not be used with 'deploy_latest_release_action' or 'deploy_new_release_action'.",
			Optional:    true,
			Type:        schema.TypeSet,
			Elem:        &schema.Resource{Schema: getRunRunbookActionSchema()},
			MaxItems:    1,
		},
		"tenant_ids": {
			Description: "The IDs of the tenants to deploy to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"is_disabled": {
			Description: "Indicates whether the trigger is disabled.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"webhook_trigger_filter": {
			Description: "Allows access to the webhook",
			Required:    true,
			Type:        schema.TypeSet,
			Elem:        &schema.Resource{Schema: getWebhookDetailsSchema()},
			MaxItems:    1,
		},
	}
}

func getWebhookDetailsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"webhook_id": {
			Description: "The ID of the webhook.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"secret": {
			Description: "The password used to gain access to the webhook.",
			Type:        schema.TypeString,
			Sensitive:   true,
			Required:    true,
		},
	}
}
