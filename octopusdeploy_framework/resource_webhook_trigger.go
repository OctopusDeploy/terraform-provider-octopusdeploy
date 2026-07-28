package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type webhookTriggerResource struct {
	*Config
}

func NewWebhookTriggerResource() resource.Resource {
	return &webhookTriggerResource{}
}

var (
	_ resource.ResourceWithImportState    = &webhookTriggerResource{}
	_ resource.ResourceWithValidateConfig = &webhookTriggerResource{}
)

func (r *webhookTriggerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("webhook_trigger")
}

func (r *webhookTriggerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.WebhookTriggerSchema{}.GetResourceSchema()
}

func (r *webhookTriggerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *webhookTriggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *webhookTriggerResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data schemas.WebhookTriggerResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Values only known at apply time can't be checked here; the server validates them anyway.
	if data.Secret.IsUnknown() || data.RequireApiKey.IsUnknown() {
		return
	}

	requireApiKey := data.RequireApiKey.ValueBool()
	hasSecret := !data.Secret.IsNull()

	if requireApiKey && hasSecret {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret"),
			"Invalid webhook trigger configuration",
			"A webhook trigger cannot have both 'secret' and 'require_api_key' set to true.",
		)
		return
	}

	if !requireApiKey && !hasSecret {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret"),
			"missing webhook trigger authentication",
			"A webhook trigger requires either a 'secret' or 'require_api_key' to be true.",
		)
	}
}

func (r *webhookTriggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.WebhookTriggerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.Client

	project, err := projects.GetByID(client, data.SpaceId.ValueString(), data.ProjectId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("could not find project", err.Error())
		return
	}

	action := expandRunRunbookAction(data)
	filter := expandWebhookTriggerFilter(data)

	tflog.Info(ctx, fmt.Sprintf("creating webhook trigger: %s", data.Name.ValueString()))

	projectTrigger := triggers.NewProjectTrigger(
		data.Name.ValueString(),
		data.Description.ValueString(),
		data.IsDisabled.ValueBool(),
		project,
		action,
		filter,
	)

	createdWebhookTrigger, err := triggers.Add(client, projectTrigger)
	if err != nil {
		resp.Diagnostics.AddError("unable to create webhook trigger", err.Error())
		return
	}

	createdFilter, ok := createdWebhookTrigger.Filter.(*filters.WebhookTriggerFilter)
	if !ok {
		resp.Diagnostics.AddError("unexpected trigger filter type", fmt.Sprintf("expected a webhook trigger filter but got %T", createdWebhookTrigger.Filter))
		return
	}

	createdAction, ok := createdWebhookTrigger.Action.(*actions.RunRunbookAction)
	if !ok {
		resp.Diagnostics.AddError("unexpected trigger action type", fmt.Sprintf("expected a run runbook action but got %T", createdWebhookTrigger.Action))
		return
	}

	data.ID = types.StringValue(createdWebhookTrigger.GetID())
	data.Name = types.StringValue(createdWebhookTrigger.Name)
	data.ProjectId = types.StringValue(createdWebhookTrigger.ProjectID)
	data.SpaceId = types.StringValue(createdWebhookTrigger.SpaceID)
	data.IsDisabled = types.BoolValue(createdWebhookTrigger.IsDisabled)
	data.WebhookId = types.StringValue(createdFilter.WebhookID)
	data.RequireApiKey = types.BoolValue(createdFilter.RequireAPIKey)
	data.RunRunbookAction = flattenRunRunbookAction(createdAction, data.RunRunbookAction)
	data.TenantIds = flattenStringList(createdAction.Tenants, data.TenantIds)

	tflog.Info(ctx, fmt.Sprintf("webhook trigger created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookTriggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.WebhookTriggerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading webhook trigger (%s)", data.ID))

	client := r.Config.Client

	webhookTrigger, err := triggers.GetById(client, data.SpaceId.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "webhook trigger"); err != nil {
			resp.Diagnostics.AddError("unable to load webhook trigger", err.Error())
		}
		return
	}

	webhookFilter, ok := webhookTrigger.Filter.(*filters.WebhookTriggerFilter)
	if !ok {
		resp.Diagnostics.AddError("unexpected trigger filter type", fmt.Sprintf("expected a webhook trigger filter but got %T", webhookTrigger.Filter))
		return
	}

	runRunbookAction, ok := webhookTrigger.Action.(*actions.RunRunbookAction)
	if !ok {
		resp.Diagnostics.AddError("unexpected trigger action type", fmt.Sprintf("expected a run runbook action but got %T", webhookTrigger.Action))
		return
	}

	data.ID = types.StringValue(webhookTrigger.GetID())
	data.Name = types.StringValue(webhookTrigger.Name)
	data.ProjectId = types.StringValue(webhookTrigger.ProjectID)
	data.SpaceId = types.StringValue(webhookTrigger.SpaceID)
	data.IsDisabled = types.BoolValue(webhookTrigger.IsDisabled)
	data.WebhookId = types.StringValue(webhookFilter.WebhookID)
	data.RequireApiKey = types.BoolValue(webhookFilter.RequireAPIKey)
	data.RunRunbookAction = flattenRunRunbookAction(runRunbookAction, data.RunRunbookAction)
	data.TenantIds = flattenStringList(runRunbookAction.Tenants, data.TenantIds)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookTriggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.WebhookTriggerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating webhook trigger '%s'", data.ID.ValueString()))

	client := r.Config.Client

	webhookTrigger, err := triggers.GetById(client, state.SpaceId.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load webhook trigger", err.Error())
		return
	}

	action := expandRunRunbookAction(data)
	filter := expandWebhookTriggerFilter(data)
	project, err := projects.GetByID(client, data.SpaceId.ValueString(), data.ProjectId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("could not find project", err.Error())
		return
	}

	updatedWebhookTrigger := triggers.NewProjectTrigger(
		data.Name.ValueString(),
		data.Description.ValueString(),
		data.IsDisabled.ValueBool(),
		project,
		action,
		filter,
	)
	updatedWebhookTrigger.ID = webhookTrigger.ID

	updatedWebhookTrigger, err = triggers.Update(client, updatedWebhookTrigger)
	if err != nil {
		resp.Diagnostics.AddError("unable to update webhook trigger", err.Error())
		return
	}

	data.WebhookId = state.WebhookId

	tflog.Info(ctx, fmt.Sprintf("webhook trigger updated (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookTriggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.WebhookTriggerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.Client

	if err := triggers.DeleteById(client, data.SpaceId.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete webhook trigger", err.Error())
		return
	}
}

func expandWebhookTriggerFilter(data *schemas.WebhookTriggerResourceModel) *filters.WebhookTriggerFilter {
	if data.RequireApiKey.ValueBool() {
		return filters.NewWebhookTriggerFilter(true, nil)
	}

	return filters.NewWebhookTriggerFilter(false, core.NewSensitiveValue(data.Secret.ValueString()))
}

func expandRunRunbookAction(data *schemas.WebhookTriggerResourceModel) *actions.RunRunbookAction {
	action := actions.NewRunRunbookAction()

	if data.RunRunbookAction != nil {
		action.Runbook = data.RunRunbookAction.RunbookId.ValueString()

		if !data.RunRunbookAction.TargetEnvironmentIds.IsNull() && !data.RunRunbookAction.TargetEnvironmentIds.IsUnknown() {
			action.Environments = util.ExpandStringList(data.RunRunbookAction.TargetEnvironmentIds)
		}
	}
	if !data.TenantIds.IsNull() && !data.TenantIds.IsUnknown() {
		action.Tenants = util.ExpandStringList(data.TenantIds)
	}

	return action
}

func flattenRunRunbookAction(action *actions.RunRunbookAction, current *schemas.RunRunbookActionModel) *schemas.RunRunbookActionModel {
	currentEnvironmentIds := types.ListNull(types.StringType)
	if current != nil {
		currentEnvironmentIds = current.TargetEnvironmentIds
	}

	return &schemas.RunRunbookActionModel{
		RunbookId:            types.StringValue(action.Runbook),
		TargetEnvironmentIds: flattenStringList(action.Environments, currentEnvironmentIds),
	}
}
