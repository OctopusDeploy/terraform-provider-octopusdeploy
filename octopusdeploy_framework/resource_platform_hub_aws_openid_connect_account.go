package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubaccounts"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &platformHubAwsOpenIDConnectAccountResource{}
var _ resource.ResourceWithImportState = &platformHubAwsOpenIDConnectAccountResource{}

type platformHubAwsOpenIDConnectAccountResource struct {
	*Config
}

func NewPlatformHubAwsOpenIDConnectAccountResource() resource.Resource {
	return &platformHubAwsOpenIDConnectAccountResource{}
}

func (a *platformHubAwsOpenIDConnectAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.PlatformHubAwsOpenIDConnectAccountResourceName)
}

func (a *platformHubAwsOpenIDConnectAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PlatformHubAwsOpenIDConnectAccountSchema{}.GetResourceSchema()
}

func (a *platformHubAwsOpenIDConnectAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	a.Config = ResourceConfiguration(req, resp)
}

func (a *platformHubAwsOpenIDConnectAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (a *platformHubAwsOpenIDConnectAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.PlatformHubAwsOpenIDConnectAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Platform Hub AWS OpenID Connect account", map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	account := expandPlatformHubAwsOIDCAccount(ctx, &plan)
	if account == nil {
		resp.Diagnostics.AddError("Error creating Platform Hub AWS OpenID Connect account", "Failed to expand account model")
		return
	}

	createdAccount, err := platformhubaccounts.Add(a.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Platform Hub AWS OpenID Connect account", err.Error())
		return
	}

	setPlatformHubAwsOIDCAccount(ctx, &plan, createdAccount.(*platformhubaccounts.PlatformHubAwsOIDCAccount))

	tflog.Debug(ctx, "Platform Hub AWS OpenID Connect account created", map[string]interface{}{
		"id":          plan.ID.ValueString(),
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (a *platformHubAwsOpenIDConnectAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.PlatformHubAwsOpenIDConnectAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Platform Hub AWS OpenID Connect account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	account, err := platformhubaccounts.GetByID(a.Client, state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "platform hub aws openid connect account"); err != nil {
			resp.Diagnostics.AddError("Error reading Platform Hub AWS OpenID Connect account", err.Error())
		}
		return
	}

	setPlatformHubAwsOIDCAccount(ctx, &state, account.(*platformhubaccounts.PlatformHubAwsOIDCAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (a *platformHubAwsOpenIDConnectAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.PlatformHubAwsOpenIDConnectAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Platform Hub AWS OpenID Connect account", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	account := expandPlatformHubAwsOIDCAccount(ctx, &plan)
	if account == nil {
		resp.Diagnostics.AddError("Error updating Platform Hub AWS OpenID Connect account", "Failed to expand account model")
		return
	}

	updatedAccount, err := platformhubaccounts.Update(a.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Platform Hub AWS OpenID Connect account", err.Error())
		return
	}

	setPlatformHubAwsOIDCAccount(ctx, &plan, updatedAccount.(*platformhubaccounts.PlatformHubAwsOIDCAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *platformHubAwsOpenIDConnectAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.PlatformHubAwsOpenIDConnectAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Platform Hub AWS OpenID Connect account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := platformhubaccounts.Delete(a.Client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Platform Hub AWS OpenID Connect account", err.Error())
		return
	}
}

func expandPlatformHubAwsOIDCAccount(ctx context.Context, model *schemas.PlatformHubAwsOpenIDConnectAccountModel) *platformhubaccounts.PlatformHubAwsOIDCAccount {
	if model == nil {
		tflog.Error(ctx, "Model is nil in expandPlatformHubAwsOIDCAccount")
		return nil
	}

	name := model.Name.ValueString()
	roleArn := model.RoleArn.ValueString()

	account, err := platformhubaccounts.NewPlatformHubAwsOIDCAccount(name, roleArn)
	if err != nil {
		tflog.Error(ctx, "Failed to create Platform Hub AWS OpenID Connect account", map[string]interface{}{
			"error": err.Error(),
		})
		return nil
	}

	// Only set these if they're not empty
	if !model.ID.IsNull() {
		account.SetID(model.ID.ValueString())
	}
	if !model.Description.IsNull() {
		account.SetDescription(model.Description.ValueString())
	}
	if !model.SessionDuration.IsNull() {
		account.SessionDuration = model.SessionDuration.ValueString()
	}

	// Convert sets to slices
	if !model.ExecutionSubjectKeys.IsNull() {
		var executionKeys []string
		model.ExecutionSubjectKeys.ElementsAs(ctx, &executionKeys, false)
		account.DeploymentSubjectKeys = executionKeys
	}

	if !model.HealthSubjectKeys.IsNull() {
		var healthKeys []string
		model.HealthSubjectKeys.ElementsAs(ctx, &healthKeys, false)
		account.HealthCheckSubjectKeys = healthKeys
	}

	if !model.AccountTestSubjectKeys.IsNull() {
		var accountTestKeys []string
		model.AccountTestSubjectKeys.ElementsAs(ctx, &accountTestKeys, false)
		account.AccountTestSubjectKeys = accountTestKeys
	}

	tflog.Debug(ctx, "Expanded Platform Hub AWS OpenID Connect account", map[string]interface{}{
		"id":   account.GetID(),
		"name": account.GetName(),
	})

	return account
}

func setPlatformHubAwsOIDCAccount(ctx context.Context, model *schemas.PlatformHubAwsOpenIDConnectAccountModel, account *platformhubaccounts.PlatformHubAwsOIDCAccount) {
	if account == nil {
		tflog.Warn(ctx, "Account is nil in setPlatformHubAwsOIDCAccount")
		return
	}

	model.ID = types.StringValue(account.GetID())
	model.Name = types.StringValue(account.GetName())
	model.Description = types.StringValue(account.GetDescription())
	model.RoleArn = types.StringValue(account.RoleArn)

	if account.SessionDuration == "" {
		model.SessionDuration = types.StringNull()
	} else {
		model.SessionDuration = types.StringValue(account.SessionDuration)
	}

	// Use FlattenStringSet to preserve user intent: null stays null, [] stays []
	model.ExecutionSubjectKeys = util.FlattenStringSet(account.DeploymentSubjectKeys, model.ExecutionSubjectKeys)
	model.HealthSubjectKeys = util.FlattenStringSet(account.HealthCheckSubjectKeys, model.HealthSubjectKeys)
	model.AccountTestSubjectKeys = util.FlattenStringSet(account.AccountTestSubjectKeys, model.AccountTestSubjectKeys)

	tflog.Debug(ctx, "Platform Hub AWS OpenID Connect account state set", map[string]interface{}{
		"id":          model.ID.ValueString(),
		"name":        model.Name.ValueString(),
		"description": model.Description.ValueString(),
		"role_arn":    model.RoleArn.ValueString(),
	})
}
