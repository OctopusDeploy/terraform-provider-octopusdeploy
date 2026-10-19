package octopusdeploy

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExternalFeedCreateReleaseTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExternalFeedCreateReleaseTriggerCreate,
		DeleteContext: resourceExternalFeedCreateReleaseTriggerDelete,
		Description:   "This resource manages External feed triggers (release creation type) in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceExternalFeedCreateReleaseTriggerRead,
		Schema:        getExternalFeedCreateReleaseTriggerSchema(),
		UpdateContext: resourceExternalFeedCreateReleaseTriggerUpdate,
	}
}

func buildExternalFeedCreateReleaseTriggerResource(d *schema.ResourceData, client *client.Client) (*triggers.ProjectTrigger, error) {
	name := d.Get("name").(string)
	spaceId := d.Get("space_id").(string)
	projectId := d.Get("project_id").(string)
	channelId := d.Get("channel_id").(string)

	isDisabled := false
	if v, ok := d.GetOk("is_disabled"); ok {
		isDisabled = v.(bool)
	}

	flattenedPackages := d.Get("package")
	packages := expandDeploymentActionSlugPackages(flattenedPackages)

	flattenedPrimaryPackages := d.Get("primary_package")
	primaryPackages := expandDeploymentActionSlugPrimaryPackages(flattenedPrimaryPackages)

	packages = append(packages, primaryPackages...)

	action := actions.NewCreateReleaseAction(channelId)
	filter := filters.NewFeedTriggerFilter(packages)

	project, err := projects.GetByID(client, spaceId, projectId)
	if err != nil {
		return nil, err
	}

	createReleaseTrigger := triggers.NewProjectTrigger(name, "", isDisabled, project, action, filter)

	return createReleaseTrigger, nil
}

func resourceExternalFeedCreateReleaseTriggerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	projectTrigger, err := buildExternalFeedCreateReleaseTriggerResource(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	resource, err := triggers.Add(client, projectTrigger)
	if err != nil {
		return diag.FromErr(err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceExternalFeedCreateReleaseTriggerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*client.Client)
	spaceId := d.Get("space_id").(string)
	spaceId = util.Ternary(len(spaceId) > 0, spaceId, client.GetSpaceID())

	projectTrigger, err := triggers.GetById(client, spaceId, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to read external feed create release trigger (%s) in space %s: %w", id, spaceId, err))
	}
	if projectTrigger == nil {
		d.SetId("")
		return nil
	}

	action, ok := projectTrigger.Action.(*actions.CreateReleaseAction)
	if !ok {
		return diag.Errorf("project trigger (%s) in space %s is not a create release trigger", id, spaceId)
	}
	filter, ok := projectTrigger.Filter.(*filters.FeedTriggerFilter)
	if !ok {
		return diag.Errorf("project trigger (%s) in space %s is not an external feed trigger", id, spaceId)
	}

	d.Set("name", projectTrigger.Name)
	d.Set("space_id", projectTrigger.SpaceID)
	d.Set("project_id", projectTrigger.ProjectID)
	d.Set("is_disabled", projectTrigger.IsDisabled)
	d.Set("channel_id", action.ChannelID)
	d.Set("package", flattenDeploymentActionSlugPackages(filter.Packages))
	d.Set("primary_package", flattenDeploymentActionSlugPrimaryPackages(filter.Packages))

	return nil
}

func resourceExternalFeedCreateReleaseTriggerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	projectTrigger, err := buildExternalFeedCreateReleaseTriggerResource(d, client)
	if err != nil {
		return diag.FromErr(err)
	}
	projectTrigger.ID = d.Id() // set ID so Octopus API knows which project trigger to update

	resource, err := triggers.Update(client, projectTrigger)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceExternalFeedCreateReleaseTriggerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	spaceId := d.Get("space_id").(string)
	spaceId = util.Ternary(len(spaceId) > 0, spaceId, client.GetSpaceID())

	err := triggers.DeleteById(client, spaceId, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete external feed create release trigger (%s) in space %s: %w", d.Id(), spaceId, err))
	}

	d.SetId("")
	return nil
}
