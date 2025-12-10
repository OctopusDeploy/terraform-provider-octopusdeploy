package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubaccounts"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &platformHubGenericOidcAccountResource{}
var _ resource.ResourceWithImportState = &platformHubGenericOidcAccountResource{}

type platformHubGenericOidcAccountResource struct {
	*Config
}

func NewPlatformHubGenericOidcAccountResource() resource.Resource {
	return &platformHubGenericOidcAccountResource{}
}

func (g *platformHubGenericOidcAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.PlatformHubGenericOidcAccountResourceName)
}

func (g *platformHubGenericOidcAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PlatformHubGenericOidcAccountSchema{}.GetResourceSchema()
}

func (g *platformHubGenericOidcAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	g.Config = ResourceConfiguration(req, resp)
}

func (g *platformHubGenericOidcAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (g *platformHubGenericOidcAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.PlatformHubGenericOidcAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Platform Hub Generic OIDC account", map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	account, diags := expandPlatformHubGenericOidcAccount(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdAccount, err := platformhubaccounts.Add(g.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Platform Hub Generic OIDC account", err.Error())
		return
	}

	setPlatformHubGenericOidcAccount(ctx, &plan, createdAccount.(*platformhubaccounts.PlatformHubGenericOidcAccount))

	tflog.Debug(ctx, "Platform Hub Generic OIDC account created", map[string]interface{}{
		"id":          plan.ID.ValueString(),
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (g *platformHubGenericOidcAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.PlatformHubGenericOidcAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Platform Hub Generic OIDC account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	account, err := platformhubaccounts.GetByID(g.Client, state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "platform hub generic oidc account"); err != nil {
			resp.Diagnostics.AddError("Error reading Platform Hub Generic OIDC account", err.Error())
		}
		return
	}

	setPlatformHubGenericOidcAccount(ctx, &state, account.(*platformhubaccounts.PlatformHubGenericOidcAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (g *platformHubGenericOidcAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.PlatformHubGenericOidcAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Platform Hub Generic OIDC account", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	account, diags := expandPlatformHubGenericOidcAccount(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedAccount, err := platformhubaccounts.Update(g.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Platform Hub Generic OIDC account", err.Error())
		return
	}

	setPlatformHubGenericOidcAccount(ctx, &plan, updatedAccount.(*platformhubaccounts.PlatformHubGenericOidcAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (g *platformHubGenericOidcAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.PlatformHubGenericOidcAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Platform Hub Generic OIDC account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := platformhubaccounts.Delete(g.Client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Platform Hub Generic OIDC account", err.Error())
		return
	}
}

func expandPlatformHubGenericOidcAccount(ctx context.Context, model *schemas.PlatformHubGenericOidcAccountModel) (*platformhubaccounts.PlatformHubGenericOidcAccount, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model == nil {
		tflog.Error(ctx, "Model is nil in expandPlatformHubGenericOidcAccount")
		diags.AddError("Invalid input", "Model is nil")
		return nil, diags
	}

	name := model.Name.ValueString()

	account, err := platformhubaccounts.NewPlatformHubGenericOidcAccount(name)
	if err != nil {
		tflog.Error(ctx, "Failed to create Platform Hub Generic OIDC account", map[string]interface{}{
			"error": err.Error(),
		})
		diags.AddError("Failed to create account", err.Error())
		return nil, diags
	}

	if !model.ID.IsNull() {
		account.SetID(model.ID.ValueString())
	}
	if !model.Description.IsNull() {
		account.SetDescription(model.Description.ValueString())
	}
	if !model.ExecutionSubjectKeys.IsNull() {
		var executionSubjectKeys []string
		diags.Append(model.ExecutionSubjectKeys.ElementsAs(ctx, &executionSubjectKeys, false)...)
		if !diags.HasError() {
			account.ExecutionSubjectKeys = executionSubjectKeys
		}
	}
	if !model.Audience.IsNull() {
		account.Audience = model.Audience.ValueString()
	}

	tflog.Debug(ctx, "Expanded Platform Hub Generic OIDC account", map[string]interface{}{
		"id":   account.GetID(),
		"name": account.GetName(),
	})

	return account, diags
}

func setPlatformHubGenericOidcAccount(ctx context.Context, model *schemas.PlatformHubGenericOidcAccountModel, account *platformhubaccounts.PlatformHubGenericOidcAccount) {
	if account == nil {
		tflog.Warn(ctx, "Account is nil in setPlatformHubGenericOidcAccount")
		return
	}

	model.ID = types.StringValue(account.GetID())
	model.Name = types.StringValue(account.GetName())
	model.Description = types.StringValue(account.GetDescription())

	if account.ExecutionSubjectKeys != nil {
		executionSubjectKeys, _ := types.ListValueFrom(ctx, types.StringType, account.ExecutionSubjectKeys)
		model.ExecutionSubjectKeys = executionSubjectKeys
	}

	if account.Audience != "" {
		model.Audience = types.StringValue(account.Audience)
	}

	tflog.Debug(ctx, "Platform Hub Generic OIDC account state set", map[string]interface{}{
		"id":          model.ID.ValueString(),
		"name":        model.Name.ValueString(),
		"description": model.Description.ValueString(),
	})
}
