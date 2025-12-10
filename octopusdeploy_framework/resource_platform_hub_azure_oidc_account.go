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

var _ resource.Resource = &platformHubAzureOidcAccountResource{}
var _ resource.ResourceWithImportState = &platformHubAzureOidcAccountResource{}

type platformHubAzureOidcAccountResource struct {
	*Config
}

func NewPlatformHubAzureOidcAccountResource() resource.Resource {
	return &platformHubAzureOidcAccountResource{}
}

func (a *platformHubAzureOidcAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.PlatformHubAzureOidcAccountResourceName)
}

func (a *platformHubAzureOidcAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PlatformHubAzureOidcAccountSchema{}.GetResourceSchema()
}

func (a *platformHubAzureOidcAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	a.Config = ResourceConfiguration(req, resp)
}

func (a *platformHubAzureOidcAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (a *platformHubAzureOidcAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.PlatformHubAzureOidcAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Platform Hub Azure OIDC account", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	account, diags := expandPlatformHubAzureOidcAccount(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdAccount, err := platformhubaccounts.Add(a.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Platform Hub Azure OIDC account", err.Error())
		return
	}

	setPlatformHubAzureOidcAccount(ctx, &plan, createdAccount.(*platformhubaccounts.PlatformHubAzureOidcAccount))

	tflog.Debug(ctx, "Platform Hub Azure OIDC account created", map[string]interface{}{
		"id":   plan.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (a *platformHubAzureOidcAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.PlatformHubAzureOidcAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Platform Hub Azure OIDC account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	account, err := platformhubaccounts.GetByID(a.Client, state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "platform hub azure oidc account"); err != nil {
			resp.Diagnostics.AddError("Error reading Platform Hub Azure OIDC account", err.Error())
		}
		return
	}

	setPlatformHubAzureOidcAccount(ctx, &state, account.(*platformhubaccounts.PlatformHubAzureOidcAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (a *platformHubAzureOidcAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.PlatformHubAzureOidcAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Platform Hub Azure OIDC account", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	account, diags := expandPlatformHubAzureOidcAccount(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedAccount, err := platformhubaccounts.Update(a.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Platform Hub Azure OIDC account", err.Error())
		return
	}

	setPlatformHubAzureOidcAccount(ctx, &plan, updatedAccount.(*platformhubaccounts.PlatformHubAzureOidcAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *platformHubAzureOidcAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.PlatformHubAzureOidcAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Platform Hub Azure OIDC account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := platformhubaccounts.Delete(a.Client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Platform Hub Azure OIDC account", err.Error())
		return
	}
}

func expandPlatformHubAzureOidcAccount(ctx context.Context, model *schemas.PlatformHubAzureOidcAccountModel) (*platformhubaccounts.PlatformHubAzureOidcAccount, diag.Diagnostics) {
	var diags diag.Diagnostics

	if model == nil {
		tflog.Error(ctx, "Model is nil in expandPlatformHubAzureOidcAccount")
		diags.AddError("Invalid input", "Model is nil")
		return nil, diags
	}

	name := model.Name.ValueString()
	subscriptionID := model.SubscriptionID.ValueString()
	applicationID := model.ApplicationID.ValueString()
	tenantID := model.TenantID.ValueString()

	account, err := platformhubaccounts.NewPlatformHubAzureOidcAccount(name, subscriptionID, applicationID, tenantID)
	if err != nil {
		tflog.Error(ctx, "Failed to create Platform Hub Azure OIDC account", map[string]interface{}{
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
	if !model.HealthSubjectKeys.IsNull() {
		var healthSubjectKeys []string
		diags.Append(model.HealthSubjectKeys.ElementsAs(ctx, &healthSubjectKeys, false)...)
		if !diags.HasError() {
			account.HealthSubjectKeys = healthSubjectKeys
		}
	}
	if !model.AccountTestSubjectKeys.IsNull() {
		var accountTestSubjectKeys []string
		diags.Append(model.AccountTestSubjectKeys.ElementsAs(ctx, &accountTestSubjectKeys, false)...)
		if !diags.HasError() {
			account.AccountTestSubjectKeys = accountTestSubjectKeys
		}
	}
	if !model.Audience.IsNull() {
		account.Audience = model.Audience.ValueString()
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

	return account, diags
}

func setPlatformHubAzureOidcAccount(ctx context.Context, model *schemas.PlatformHubAzureOidcAccountModel, account *platformhubaccounts.PlatformHubAzureOidcAccount) {
	if account == nil {
		tflog.Warn(ctx, "Account is nil in setPlatformHubAzureOidcAccount")
		return
	}

	model.ID = types.StringValue(account.GetID())
	model.Name = types.StringValue(account.GetName())
	model.Description = types.StringValue(account.GetDescription())
	model.SubscriptionID = types.StringValue(account.SubscriptionID)
	model.ApplicationID = types.StringValue(account.ApplicationID)
	model.TenantID = types.StringValue(account.TenantID)

	if account.ExecutionSubjectKeys != nil {
		executionSubjectKeys, _ := types.ListValueFrom(ctx, types.StringType, account.ExecutionSubjectKeys)
		model.ExecutionSubjectKeys = executionSubjectKeys
	}
	if account.HealthSubjectKeys != nil {
		healthSubjectKeys, _ := types.ListValueFrom(ctx, types.StringType, account.HealthSubjectKeys)
		model.HealthSubjectKeys = healthSubjectKeys
	}
	if account.AccountTestSubjectKeys != nil {
		accountTestSubjectKeys, _ := types.ListValueFrom(ctx, types.StringType, account.AccountTestSubjectKeys)
		model.AccountTestSubjectKeys = accountTestSubjectKeys
	}

	model.Audience = types.StringValue(account.Audience)
	model.AzureEnvironment = types.StringValue(account.AzureEnvironment)
	model.AuthenticationEndpoint = types.StringValue(account.AuthenticationEndpoint)
	model.ResourceManagementEndpoint = types.StringValue(account.ResourceManagementEndpoint)

	tflog.Debug(ctx, "Platform Hub Azure OIDC account state set", map[string]interface{}{
		"id":   model.ID.ValueString(),
		"name": model.Name.ValueString(),
	})
}
