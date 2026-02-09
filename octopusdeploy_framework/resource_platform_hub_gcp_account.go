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

var _ resource.Resource = &platformHubGcpAccountResource{}
var _ resource.ResourceWithImportState = &platformHubGcpAccountResource{}

type platformHubGcpAccountResource struct {
	*Config
}

func NewPlatformHubGcpAccountResource() resource.Resource {
	return &platformHubGcpAccountResource{}
}

func (g *platformHubGcpAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.PlatformHubGcpAccountResourceName)
}

func (g *platformHubGcpAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PlatformHubGcpAccountSchema{}.GetResourceSchema()
}

func (g *platformHubGcpAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	g.Config = ResourceConfiguration(req, resp)
}

func (g *platformHubGcpAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (g *platformHubGcpAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.PlatformHubGcpAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Platform Hub GCP account", map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	account := expandPlatformHubGcpAccount(&plan)
	if account == nil {
		resp.Diagnostics.AddError("Error creating Platform Hub GCP account", "Failed to expand account model")
		return
	}

	createdAccount, err := platformhubaccounts.Add(g.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Platform Hub GCP account", err.Error())
		return
	}

	setPlatformHubGcpAccount(ctx, &plan, createdAccount.(*platformhubaccounts.PlatformHubGcpAccount))

	tflog.Debug(ctx, "Platform Hub GCP account created", map[string]interface{}{
		"id":          plan.ID.ValueString(),
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (g *platformHubGcpAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.PlatformHubGcpAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Platform Hub GCP account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	account, err := platformhubaccounts.GetByID(g.Client, state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "platform hub gcp account"); err != nil {
			resp.Diagnostics.AddError("Error reading Platform Hub GCP account", err.Error())
		}
		return
	}

	setPlatformHubGcpAccount(ctx, &state, account.(*platformhubaccounts.PlatformHubGcpAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (g *platformHubGcpAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.PlatformHubGcpAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Platform Hub GCP account", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	account := expandPlatformHubGcpAccount(&plan)
	if account == nil {
		resp.Diagnostics.AddError("Error updating Platform Hub GCP account", "Failed to expand account model")
		return
	}

	updatedAccount, err := platformhubaccounts.Update(g.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Platform Hub GCP account", err.Error())
		return
	}

	setPlatformHubGcpAccount(ctx, &plan, updatedAccount.(*platformhubaccounts.PlatformHubGcpAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (g *platformHubGcpAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.PlatformHubGcpAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Platform Hub GCP account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := platformhubaccounts.Delete(g.Client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Platform Hub GCP account", err.Error())
		return
	}
}

func expandPlatformHubGcpAccount(model *schemas.PlatformHubGcpAccountModel) *platformhubaccounts.PlatformHubGcpAccount {
	if model == nil {
		tflog.Error(context.Background(), "Model is nil in expandPlatformHubGcpAccount")
		return nil
	}

	name := model.Name.ValueString()
	jsonKey := core.NewSensitiveValue(model.JsonKey.ValueString())

	account, err := platformhubaccounts.NewPlatformHubGcpAccount(name, jsonKey)
	if err != nil {
		tflog.Error(context.Background(), "Failed to create Platform Hub GCP account", map[string]interface{}{
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

	tflog.Debug(context.Background(), "Expanded Platform Hub GCP account", map[string]interface{}{
		"id":   account.GetID(),
		"name": account.GetName(),
	})

	return account
}

func setPlatformHubGcpAccount(ctx context.Context, model *schemas.PlatformHubGcpAccountModel, account *platformhubaccounts.PlatformHubGcpAccount) {
	if account == nil {
		tflog.Warn(ctx, "Account is nil in setPlatformHubGcpAccount")
		return
	}

	model.ID = types.StringValue(account.GetID())
	model.Name = types.StringValue(account.GetName())
	model.Description = types.StringValue(account.GetDescription())

	tflog.Debug(ctx, "Platform Hub GCP account state set", map[string]interface{}{
		"id":          model.ID.ValueString(),
		"name":        model.Name.ValueString(),
		"description": model.Description.ValueString(),
	})
}
