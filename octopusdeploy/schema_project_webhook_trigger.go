package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
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
	return flattenedProjectWebhookTrigger
}

func expandProjectWebhookTrigger(projectWebhookTrigger *schema.ResourceData, project *projects.Project) (*triggers.ProjectTrigger, error) {
	name := projectWebhookTrigger.Get("name").(string)
	description := projectWebhookTrigger.Get("description").(string)
	isDisabled := projectWebhookTrigger.Get("is_disabled").(bool)

	var action actions.ITriggerAction = nil
	var filter filters.ITriggerFilter = nil

	if attr, ok := projectWebhookTrigger.GetOk("run_runbook_action"); ok {
		runRunbookActionList := attr.(*schema.Set).List()
		runRunbookActionMap := runRunbookActionList[0].(map[string]interface{})
		deploymentAction := actions.NewRunRunbookAction()

		deploymentAction.Runbook = runRunbookActionMap["runbook_id"].(string)
		deploymentAction.Environments = expandArray(runRunbookActionMap["target_environment_ids"].([]interface{}))
		deploymentAction.Tenants = expandArray(projectWebhookTrigger.Get("tenant_ids").([]interface{}))
		action = deploymentAction
	}

	// Filter configuration
	//password configuration

	// if NewProjectTrigger doesn't actually use the description value
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
	}
}
