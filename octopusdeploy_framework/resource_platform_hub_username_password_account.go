package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubaccounts"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &platformHubUsernamePasswordAccountResource{}
var _ resource.ResourceWithImportState = &platformHubUsernamePasswordAccountResource{}

type platformHubUsernamePasswordAccountResource struct {
	*Config
}

func NewPlatformHubUsernamePasswordAccountResource() resource.Resource {
	return &platformHubUsernamePasswordAccountResource{}
}

func (u *platformHubUsernamePasswordAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.PlatformHubUsernamePasswordAccountResourceName)
}

func (u *platformHubUsernamePasswordAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PlatformHubUsernamePasswordAccountSchema{}.GetResourceSchema()
}

func (u *platformHubUsernamePasswordAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	u.Config = ResourceConfiguration(req, resp)
}

func (u *platformHubUsernamePasswordAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (u *platformHubUsernamePasswordAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.PlatformHubUsernamePasswordAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Platform Hub Username-Password account", map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	account := expandPlatformHubUsernamePasswordAccount(&plan)
	if account == nil {
		resp.Diagnostics.AddError("Error creating Platform Hub Username-Password account", "Failed to expand account model")
		return
	}

	createdAccount, err := platformhubaccounts.Add(u.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Platform Hub Username-Password account", err.Error())
		return
	}

	setPlatformHubUsernamePasswordAccount(ctx, &plan, createdAccount.(*platformhubaccounts.PlatformHubUsernamePasswordAccount))

	tflog.Debug(ctx, "Platform Hub Username-Password account created", map[string]interface{}{
		"id":          plan.ID.ValueString(),
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (u *platformHubUsernamePasswordAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.PlatformHubUsernamePasswordAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Platform Hub Username-Password account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	account, err := platformhubaccounts.GetByID(u.Client, state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "platform hub username-password account"); err != nil {
			resp.Diagnostics.AddError("Error reading Platform Hub Username-Password account", err.Error())
		}
		return
	}

	setPlatformHubUsernamePasswordAccount(ctx, &state, account.(*platformhubaccounts.PlatformHubUsernamePasswordAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (u *platformHubUsernamePasswordAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.PlatformHubUsernamePasswordAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Platform Hub Username-Password account", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	account := expandPlatformHubUsernamePasswordAccount(&plan)
	if account == nil {
		resp.Diagnostics.AddError("Error updating Platform Hub Username-Password account", "Failed to expand account model")
		return
	}

	updatedAccount, err := platformhubaccounts.Update(u.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Platform Hub Username-Password account", err.Error())
		return
	}

	setPlatformHubUsernamePasswordAccount(ctx, &plan, updatedAccount.(*platformhubaccounts.PlatformHubUsernamePasswordAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (u *platformHubUsernamePasswordAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.PlatformHubUsernamePasswordAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Platform Hub Username-Password account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := platformhubaccounts.Delete(u.Client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Platform Hub Username-Password account", err.Error())
		return
	}
}

func expandPlatformHubUsernamePasswordAccount(model *schemas.PlatformHubUsernamePasswordAccountModel) *platformhubaccounts.PlatformHubUsernamePasswordAccount {
	if model == nil {
		tflog.Error(context.Background(), "Model is nil in expandPlatformHubUsernamePasswordAccount")
		return nil
	}

	name := model.Name.ValueString()
	username := model.Username.ValueString()
	password := core.NewSensitiveValue(model.Password.ValueString())

	account, err := platformhubaccounts.NewPlatformHubUsernamePasswordAccount(name, username, password)
	if err != nil {
		tflog.Error(context.Background(), "Failed to create Platform Hub Username-Password account", map[string]interface{}{
			"error": err.Error(),
		})
		return nil
	}

	if !model.ID.IsNull() {
		account.SetID(model.ID.ValueString())
	}
	if !model.Description.IsNull() {
		account.SetDescription(model.Description.ValueString())
	}

	tflog.Debug(context.Background(), "Expanded Platform Hub Username-Password account", map[string]interface{}{
		"id":   account.GetID(),
		"name": account.GetName(),
	})

	return account
}

func setPlatformHubUsernamePasswordAccount(ctx context.Context, model *schemas.PlatformHubUsernamePasswordAccountModel, account *platformhubaccounts.PlatformHubUsernamePasswordAccount) {
	if account == nil {
		tflog.Warn(ctx, "Account is nil in setPlatformHubUsernamePasswordAccount")
		return
	}

	model.ID = types.StringValue(account.GetID())
	model.Name = types.StringValue(account.GetName())
	model.Description = types.StringValue(account.GetDescription())
	model.Username = types.StringValue(account.Username)

	tflog.Debug(ctx, "Platform Hub Username-Password account state set", map[string]interface{}{
		"id":          model.ID.ValueString(),
		"name":        model.Name.ValueString(),
		"description": model.Description.ValueString(),
		"username":    model.Username.ValueString(),
	})
}
