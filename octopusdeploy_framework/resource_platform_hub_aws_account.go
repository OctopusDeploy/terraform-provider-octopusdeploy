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

var _ resource.Resource = &platformHubAwsAccountResource{}
var _ resource.ResourceWithImportState = &platformHubAwsAccountResource{}

type platformHubAwsAccountResource struct {
	*Config
}

func NewPlatformHubAwsAccountResource() resource.Resource {
	return &platformHubAwsAccountResource{}
}

func (a *platformHubAwsAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.PlatformHubAwsAccountResourceName)
}

func (a *platformHubAwsAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PlatformHubAwsAccountSchema{}.GetResourceSchema()
}

func (a *platformHubAwsAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	a.Config = ResourceConfiguration(req, resp)
}

func (a *platformHubAwsAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (a *platformHubAwsAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.PlatformHubAwsAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Platform Hub AWS account", map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	account := expandPlatformHubAwsAccount(&plan)
	if account == nil {
		resp.Diagnostics.AddError("Error creating Platform Hub AWS account", "Failed to expand account model")
		return
	}

	createdAccount, err := platformhubaccounts.Add(a.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Platform Hub AWS account", err.Error())
		return
	}

	setPlatformHubAwsAccount(ctx, &plan, createdAccount.(*platformhubaccounts.PlatformHubAwsAccount))

	tflog.Debug(ctx, "Platform Hub AWS account created", map[string]interface{}{
		"id":          plan.ID.ValueString(),
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (a *platformHubAwsAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.PlatformHubAwsAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Platform Hub AWS account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	account, err := platformhubaccounts.GetByID(a.Client, state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "platform hub aws account"); err != nil {
			resp.Diagnostics.AddError("Error reading Platform Hub AWS account", err.Error())
		}
		return
	}

	setPlatformHubAwsAccount(ctx, &state, account.(*platformhubaccounts.PlatformHubAwsAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (a *platformHubAwsAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.PlatformHubAwsAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Platform Hub AWS account", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	account := expandPlatformHubAwsAccount(&plan)
	if account == nil {
		resp.Diagnostics.AddError("Error updating Platform Hub AWS account", "Failed to expand account model")
		return
	}

	updatedAccount, err := platformhubaccounts.Update(a.Client, account)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Platform Hub AWS account", err.Error())
		return
	}

	setPlatformHubAwsAccount(ctx, &plan, updatedAccount.(*platformhubaccounts.PlatformHubAwsAccount))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *platformHubAwsAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.PlatformHubAwsAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Platform Hub AWS account", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := platformhubaccounts.Delete(a.Client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Platform Hub AWS account", err.Error())
		return
	}
}

func expandPlatformHubAwsAccount(model *schemas.PlatformHubAwsAccountModel) *platformhubaccounts.PlatformHubAwsAccount {
	if model == nil {
		tflog.Error(context.Background(), "Model is nil in expandPlatformHubAwsAccount")
		return nil
	}

	name := model.Name.ValueString()
	accessKey := model.AccessKey.ValueString()
	secretKey := core.NewSensitiveValue(model.SecretKey.ValueString())

	account, err := platformhubaccounts.NewPlatformHubAwsAccount(name, accessKey, secretKey)
	if err != nil {
		tflog.Error(context.Background(), "Failed to create Platform Hub AWS account", map[string]interface{}{
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

	tflog.Debug(context.Background(), "Expanded Platform Hub AWS account", map[string]interface{}{
		"id":   account.GetID(),
		"name": account.GetName(),
	})

	return account
}

func setPlatformHubAwsAccount(ctx context.Context, model *schemas.PlatformHubAwsAccountModel, account *platformhubaccounts.PlatformHubAwsAccount) {
	if account == nil {
		tflog.Warn(ctx, "Account is nil in setPlatformHubAwsAccount")
		return
	}

	model.ID = types.StringValue(account.GetID())
	model.Name = types.StringValue(account.GetName())
	model.Description = types.StringValue(account.GetDescription())
	model.AccessKey = types.StringValue(account.AccessKey)

	tflog.Debug(ctx, "Platform Hub AWS account state set", map[string]interface{}{
		"id":          model.ID.ValueString(),
		"name":        model.Name.ValueString(),
		"description": model.Description.ValueString(),
		"access_key":  model.AccessKey.ValueString(),
	})
}
