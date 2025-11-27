package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubversioncontrolsettings"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type platformHubVersionControlUsernamePasswordSettingsResource struct {
	*Config
}

var _ resource.Resource = &platformHubVersionControlUsernamePasswordSettingsResource{}

func NewPlatformHubVersionControlUsernamePasswordSettingsResource() resource.Resource {
	return &platformHubVersionControlUsernamePasswordSettingsResource{}
}

func (r *platformHubVersionControlUsernamePasswordSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("platform_hub_version_control_username_password_settings")
}

func (r *platformHubVersionControlUsernamePasswordSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PlatformHubVersionControlUsernamePasswordSettingsSchema{}.GetResourceSchema()
}

func (r *platformHubVersionControlUsernamePasswordSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *platformHubVersionControlUsernamePasswordSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.PlatformHubVersionControlUsernamePasswordSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating platform hub version control username password settings")

	creds := credentials.NewUsernamePassword(
		plan.Username.ValueString(),
		core.NewSensitiveValue(plan.Password.ValueString()),
	)

	settings := platformhubversioncontrolsettings.NewResource(
		plan.URL.ValueString(),
		creds,
		plan.DefaultBranch.ValueString(),
		plan.BasePath.ValueString(),
	)

	updatedSettings, err := platformhubversioncontrolsettings.Update(r.Client, settings)
	if err != nil {
		resp.Diagnostics.AddError("Error creating platform hub version control settings", err.Error())
		return
	}

	state := flattenPlatformHubVersionControlUsernamePasswordSettings(ctx, updatedSettings, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *platformHubVersionControlUsernamePasswordSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.PlatformHubVersionControlUsernamePasswordSettingsResourceModel
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

	usernamePasswordCreds, ok := settings.Credentials.(*credentials.UsernamePassword)
	if !ok || settings.URL == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	newState := flattenPlatformHubVersionControlUsernamePasswordSettings(ctx, settings, state)
	newState.Password = state.Password
	newState.Username = types.StringValue(usernamePasswordCreds.Username)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *platformHubVersionControlUsernamePasswordSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.PlatformHubVersionControlUsernamePasswordSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating platform hub version control username password settings")

	creds := credentials.NewUsernamePassword(
		plan.Username.ValueString(),
		core.NewSensitiveValue(plan.Password.ValueString()),
	)

	settings := platformhubversioncontrolsettings.NewResource(
		plan.URL.ValueString(),
		creds,
		plan.DefaultBranch.ValueString(),
		plan.BasePath.ValueString(),
	)

	updatedSettings, err := platformhubversioncontrolsettings.Update(r.Client, settings)
	if err != nil {
		resp.Diagnostics.AddError("Error updating platform hub version control settings", err.Error())
		return
	}

	state := flattenPlatformHubVersionControlUsernamePasswordSettings(ctx, updatedSettings, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *platformHubVersionControlUsernamePasswordSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.PlatformHubVersionControlUsernamePasswordSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning(
		"Platform Hub version control settings removed from state",
		"The configuration remains on the server and must be manually removed if needed.",
	)
}

func (r *platformHubVersionControlUsernamePasswordSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	settings, err := platformhubversioncontrolsettings.Get(r.Client)
	if err != nil {
		resp.Diagnostics.AddError("Error importing platform hub version control settings", err.Error())
		return
	}

	if _, ok := settings.Credentials.(*credentials.UsernamePassword); !ok {
		resp.Diagnostics.AddError("Invalid credentials type", "Platform hub version control settings do not use username/password credentials")
		return
	}

	state := schemas.PlatformHubVersionControlUsernamePasswordSettingsResourceModel{
		URL:           types.StringValue(settings.URL),
		DefaultBranch: types.StringValue(settings.DefaultBranch),
		BasePath:      types.StringValue(settings.BasePath),
		Username:      types.StringNull(),
		Password:      types.StringNull(),
	}
	state.ID = types.StringValue("platform-hub-version-control-settings")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func flattenPlatformHubVersionControlUsernamePasswordSettings(ctx context.Context, settings *platformhubversioncontrolsettings.Resource, model schemas.PlatformHubVersionControlUsernamePasswordSettingsResourceModel) schemas.PlatformHubVersionControlUsernamePasswordSettingsResourceModel {
	model.ID = types.StringValue("platform-hub-version-control-settings")
	model.URL = types.StringValue(settings.URL)
	model.DefaultBranch = types.StringValue(settings.DefaultBranch)
	model.BasePath = types.StringValue(settings.BasePath)

	return model
}
