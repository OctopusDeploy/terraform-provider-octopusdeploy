package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectWebhookTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectWebhookTriggerCreate,
		DeleteContext: resourceProjectWebhookTriggerDelete,
		Description:   "This resource manages a webhook trigger for a project or runbook in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceProjectWebhookTriggerRead,
		Schema:        getProjectWebhookTriggerSchema(),
		UpdateContext: resourceProjectWebhookTriggerUpdate,
	}
}

func resourceProjectWebhookTriggerRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	octopusClient := m.(*client.Client)
	spaceId := data.Get("space_id").(string)
	spaceId = util.Ternary(len(spaceId) > 0, spaceId, octopusClient.GetSpaceID())

	webhookTrigger, err := triggers.GetById(octopusClient, spaceId, data.Id())

	if webhookTrigger == nil {
		data.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		return nil
	}

	flattenedWebhookTrigger := flattenProjectWebhookTrigger(webhookTrigger)
	for key, value := range flattenedWebhookTrigger {
		err := data.Set(key, value)
		if err != nil {
			return nil
		}
	}

	return nil
}

func resourceProjectWebhookTriggerCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	octopusClient := m.(*client.Client)
	projectId := data.Get("project_id").(string)
	spaceId := data.Get("space_id").(string)
	project, err := projects.GetByID(octopusClient, spaceId, projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	expandedWebhookTrigger, err := expandProjectWebhookTrigger(data, project)

	if err != nil {
		return diag.FromErr(err)
	}

	webhookTrigger, err := triggers.Add(octopusClient, expandedWebhookTrigger)

	if err != nil {
		return diag.FromErr(err)
	}

	if isEmpty(webhookTrigger.GetID()) {
		log.Println("ID is nil")
	} else {
		data.SetId(webhookTrigger.GetID())
	}

	return nil
}

func resourceProjectWebhookTriggerUpdate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	octopusClient := m.(*client.Client)
	projectId := data.Get("project_id").(string)
	spaceId := data.Get("space_id").(string)
	project, err := projects.GetByID(octopusClient, spaceId, projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	expandedWebhookTrigger, err := expandProjectWebhookTrigger(data, project)

	if err != nil {
		return diag.FromErr(err)
	}

	expandedWebhookTrigger.ID = data.Id()

	if err != nil {
		return diag.FromErr(err)
	}

	webhookTrigger, err := triggers.Update(octopusClient, expandedWebhookTrigger)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(webhookTrigger.GetID())

	return nil
}

func resourceProjectWebhookTriggerDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	octopusClient := m.(*client.Client)
	spaceId := data.Get("space_id").(string)
	err := triggers.DeleteById(octopusClient, spaceId, data.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}
