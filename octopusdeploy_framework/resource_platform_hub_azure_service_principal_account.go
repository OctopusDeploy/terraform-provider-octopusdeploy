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

var _ resource.Resource = &platformHubAzureServicePrincipalAccountResource{}
var _ resource.ResourceWithImportState = &platformHubAzureServicePrincipalAccountResource{}

type platformHubAzureServicePrincipalAccountResource struct {
	*Config
}

func NewPlatformHubAzureServicePrincipalAccountResource() resource.Resource {
	return &platformHubAzureServicePrincipalAccountResource{}
}

func (a *platformHubAzureServicePrincipalAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.PlatformHubAzureServicePrincipalAccountResourceName)
}

func (a *platformHubAzureServicePrincipalAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PlatformHubAzureServicePrincipalAccountSchema{}.GetResourceSchema()
}

func (a *platformHubAzureServicePrincipalAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	a.Config = ResourceConfiguration(req, resp)
}

func (a *platformHubAzureServicePrincipalAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (a *platformHubAzureServicePrincipalAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.PlatformHubAzureServicePrincipalAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Platform Hub Azure Service Principal account", map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	account := expandPlatformHubAzureServicePrincipalAccount(&plan)
	if account == nil {
		resp.Diagnostics.AddError("Error creating Platform Hub Azure Service Principal account", "Failed to expand account model")
		return
	}

	createdAccount, err := platformhubaccounts.Add(a.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Platform Hub Azure Service Principal account", err.Error())
		return
	}

	setPlatformHubAzureServicePrincipalAccount(ctx, &plan, createdAccount.(*platformhubaccounts.PlatformHubAzureServicePrincipalAccount))

	tflog.Debug(ctx, "Platform Hub Azure Service Principal account created", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (a *platformHubAzureServicePrincipalAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.PlatformHubAzureServicePrincipalAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Platform Hub Azure Service Principal account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	account, err := platformhubaccounts.GetByID(a.Client, state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "platform hub azure service principal account"); err != nil {
			resp.Diagnostics.AddError("Error reading Platform Hub Azure Service Principal account", err.Error())
		}
		return
	}

	setPlatformHubAzureServicePrincipalAccount(ctx, &state, account.(*platformhubaccounts.PlatformHubAzureServicePrincipalAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (a *platformHubAzureServicePrincipalAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.PlatformHubAzureServicePrincipalAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Platform Hub Azure Service Principal account", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	account := expandPlatformHubAzureServicePrincipalAccount(&plan)
	if account == nil {
		resp.Diagnostics.AddError("Error updating Platform Hub Azure Service Principal account", "Failed to expand account model")
		return
	}

	updatedAccount, err := platformhubaccounts.Update(a.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Platform Hub Azure Service Principal account", err.Error())
		return
	}

	setPlatformHubAzureServicePrincipalAccount(ctx, &plan, updatedAccount.(*platformhubaccounts.PlatformHubAzureServicePrincipalAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *platformHubAzureServicePrincipalAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.PlatformHubAzureServicePrincipalAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Platform Hub Azure Service Principal account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := platformhubaccounts.Delete(a.Client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Platform Hub Azure Service Principal account", err.Error())
		return
	}
}

func expandPlatformHubAzureServicePrincipalAccount(model *schemas.PlatformHubAzureServicePrincipalAccountModel) *platformhubaccounts.PlatformHubAzureServicePrincipalAccount {
	if model == nil {
		tflog.Error(context.Background(), "Model is nil in expandPlatformHubAzureServicePrincipalAccount")
		return nil
	}

	name := model.Name.ValueString()
	subscriptionID := model.SubscriptionID.ValueString()
	tenantID := model.TenantID.ValueString()
	applicationID := model.ApplicationID.ValueString()
	password := core.NewSensitiveValue(model.Password.ValueString())

	account, err := platformhubaccounts.NewPlatformHubAzureServicePrincipalAccount(name, subscriptionID, tenantID, applicationID, password)
	if err != nil {
		tflog.Error(context.Background(), "Failed to create Platform Hub Azure Service Principal account", map[string]interface{}{
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
	if !model.AzureEnvironment.IsNull() {
		account.AzureEnvironment = model.AzureEnvironment.ValueString()
	}
	if !model.AuthenticationEndpoint.IsNull() {
		account.AuthenticationEndpoint = model.AuthenticationEndpoint.ValueString()
	}
	if !model.ResourceManagementEndpoint.IsNull() {
		account.ResourceManagementEndpoint = model.ResourceManagementEndpoint.ValueString()
	}

	return account
}

func setPlatformHubAzureServicePrincipalAccount(ctx context.Context, model *schemas.PlatformHubAzureServicePrincipalAccountModel, account *platformhubaccounts.PlatformHubAzureServicePrincipalAccount) {
	if account == nil {
		tflog.Warn(ctx, "Account is nil in setPlatformHubAzureServicePrincipalAccount")
		return
	}

	model.ID = types.StringValue(account.GetID())
	model.Name = types.StringValue(account.GetName())
	model.Description = types.StringValue(account.GetDescription())
	model.SubscriptionID = types.StringValue(account.SubscriptionID)
	model.TenantID = types.StringValue(account.TenantID)
	model.ApplicationID = types.StringValue(account.ApplicationID)
	model.AzureEnvironment = types.StringValue(account.AzureEnvironment)
	model.AuthenticationEndpoint = types.StringValue(account.AuthenticationEndpoint)
	model.ResourceManagementEndpoint = types.StringValue(account.ResourceManagementEndpoint)

	tflog.Debug(ctx, "Platform Hub Azure Service Principal account state set", map[string]interface{}{
		"id":   model.ID.ValueString(),
		"name": model.Name.ValueString(),
	})
}
