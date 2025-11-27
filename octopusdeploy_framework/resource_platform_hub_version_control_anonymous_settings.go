package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubversioncontrolsettings"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type platformHubVersionControlAnonymousSettingsResource struct {
	*Config
}

var _ resource.Resource = &platformHubVersionControlAnonymousSettingsResource{}

func NewPlatformHubVersionControlAnonymousSettingsResource() resource.Resource {
	return &platformHubVersionControlAnonymousSettingsResource{}
}

func (r *platformHubVersionControlAnonymousSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("platform_hub_version_control_anonymous_settings")
}

func (r *platformHubVersionControlAnonymousSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PlatformHubVersionControlAnonymousSettingsSchema{}.GetResourceSchema()
}

func (r *platformHubVersionControlAnonymousSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *platformHubVersionControlAnonymousSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.PlatformHubVersionControlAnonymousSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating platform hub version control anonymous settings")

	settings := platformhubversioncontrolsettings.NewResource(
		plan.URL.ValueString(),
		credentials.NewAnonymous(),
		plan.DefaultBranch.ValueString(),
		plan.BasePath.ValueString(),
	)

	updatedSettings, err := platformhubversioncontrolsettings.Update(r.Client, settings)
	if err != nil {
		resp.Diagnostics.AddError("Error creating platform hub version control settings", err.Error())
		return
	}

	state := flattenPlatformHubVersionControlAnonymousSettings(ctx, updatedSettings)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *platformHubVersionControlAnonymousSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.PlatformHubVersionControlAnonymousSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := platformhubversioncontrolsettings.Get(r.Client)
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "platform hub version control settings"); err != nil {
			resp.Diagnostics.AddError("unable to load platform hub version control settings", err.Error())
		}
		return
	}

	_, ok := settings.Credentials.(*credentials.Anonymous)
	if !ok || settings.URL == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	newState := flattenPlatformHubVersionControlAnonymousSettings(ctx, settings)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *platformHubVersionControlAnonymousSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.PlatformHubVersionControlAnonymousSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating platform hub version control anonymous settings")

	settings := platformhubversioncontrolsettings.NewResource(
		plan.URL.ValueString(),
		credentials.NewAnonymous(),
		plan.DefaultBranch.ValueString(),
		plan.BasePath.ValueString(),
	)

	updatedSettings, err := platformhubversioncontrolsettings.Update(r.Client, settings)
	if err != nil {
		resp.Diagnostics.AddError("Error updating platform hub version control settings", err.Error())
		return
	}

	state := flattenPlatformHubVersionControlAnonymousSettings(ctx, updatedSettings)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *platformHubVersionControlAnonymousSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.PlatformHubVersionControlAnonymousSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning(
		"Platform Hub version control settings removed from state",
		"The configuration remains on the server and must be manually removed if needed.",
	)
}

func (r *platformHubVersionControlAnonymousSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	settings, err := platformhubversioncontrolsettings.Get(r.Client)
	if err != nil {
		resp.Diagnostics.AddError("Error importing platform hub version control settings", err.Error())
		return
	}

	if _, ok := settings.Credentials.(*credentials.Anonymous); !ok {
		resp.Diagnostics.AddError("Invalid credentials type", "Platform hub version control settings do not use anonymous credentials")
		return
	}

	state := flattenPlatformHubVersionControlAnonymousSettings(ctx, settings)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenPlatformHubVersionControlAnonymousSettings(ctx context.Context, settings *platformhubversioncontrolsettings.Resource) schemas.PlatformHubVersionControlAnonymousSettingsResourceModel {
	model := schemas.PlatformHubVersionControlAnonymousSettingsResourceModel{
		URL:           types.StringValue(settings.URL),
		DefaultBranch: types.StringValue(settings.DefaultBranch),
		BasePath:      types.StringValue(settings.BasePath),
	}
	model.ID = types.StringValue("platform-hub-version-control-settings")
	return model
}
