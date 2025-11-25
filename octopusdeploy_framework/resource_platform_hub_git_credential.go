package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubgitcredential"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &platformHubGitCredentialResource{}
var _ resource.ResourceWithImportState = &platformHubGitCredentialResource{}

type platformHubGitCredentialResource struct {
	*Config
}

type platformHubGitCredentialResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`

	RepositoryRestrictions *gitCredentialRepositoryRestrictionResourceModel `tfsdk:"repository_restrictions"`

	schemas.ResourceModel
}

func NewPlatformHubGitCredentialResource() resource.Resource {
	return &platformHubGitCredentialResource{}
}

func (g *platformHubGitCredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.PlatformHubGitCredentialResourceName)
}

func (g *platformHubGitCredentialResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PlatformHubGitCredentialSchema{}.GetResourceSchema()
}

func (g *platformHubGitCredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	g.Config = ResourceConfiguration(req, resp)
}

func (g *platformHubGitCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (g *platformHubGitCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan platformHubGitCredentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Git credential", map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	gitCredential := expandPlatformHubGitCredential(&plan)
	createdResponse, err := platformhubgitcredential.Add(g.Client, gitCredential)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Platform Hub Git credential", err.Error())
		return
	}

	createdGitCredential, err := platformhubgitcredential.GetByID(g.Client, createdResponse.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving created Platform Hub Git credential", err.Error())
		return
	}

	setPlatformHubGitCredential(ctx, &plan, createdGitCredential)

	tflog.Debug(ctx, "Platform Hub Git credential created", map[string]interface{}{
		"id":          plan.ID.ValueString(),
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (g *platformHubGitCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state platformHubGitCredentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Git credential", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	gitCredential, err := platformhubgitcredential.GetByID(g.Client, state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "git credential"); err != nil {
			resp.Diagnostics.AddError("Error reading Platform Hub Git credential", err.Error())
		}
		return
	}

	setPlatformHubGitCredential(ctx, &state, gitCredential)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (g *platformHubGitCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan platformHubGitCredentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gitCredential := expandPlatformHubGitCredential(&plan)
	_, err := platformhubgitcredential.Update(g.Client, gitCredential)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Platform Hub Git credential", err.Error())
		return
	}

	updatedGitCredential, err := platformhubgitcredential.GetByID(g.Client, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving updated Platform Hub Git credential", err.Error())
		return
	}

	setPlatformHubGitCredential(ctx, &plan, updatedGitCredential)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (g *platformHubGitCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state platformHubGitCredentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := platformhubgitcredential.Delete(g.Client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Platform Hub Git credential", err.Error())
		return
	}
}

func expandPlatformHubGitCredential(model *platformHubGitCredentialResourceModel) *platformhubgitcredential.PlatformHubGitCredential {
	if model == nil {
		tflog.Error(context.Background(), "Model is nil in expandPlatformHubGitCredential")
		return nil
	}

	password := core.NewSensitiveValue(model.Password.ValueString())
	name := model.Name.ValueString()
	username := model.Username.ValueString()

	usernamePassword := credentials.NewUsernamePassword(username, password)

	gitCredential := platformhubgitcredential.NewPlatformHubGitCredential(name, usernamePassword)

	// Only set these if they're not empty
	if !model.ID.IsNull() {
		gitCredential.ID = model.ID.ValueString()
	}
	if !model.Description.IsNull() {
		gitCredential.Description = model.Description.ValueString()
	}

	if model.RepositoryRestrictions != nil {
		var allowedRepositories = make([]string, 0, len(model.RepositoryRestrictions.AllowedRepositories.Elements()))
		for _, url := range model.RepositoryRestrictions.AllowedRepositories.Elements() {
			allowedRepositories = append(allowedRepositories, url.(types.String).ValueString())
		}
		gitCredential.RepositoryRestrictions = &credentials.RepositoryRestrictions{Enabled: model.RepositoryRestrictions.Enabled.ValueBool(), AllowedRepositories: allowedRepositories}
	} else {
		//Default to disabled if state doesn't have it
		gitCredential.RepositoryRestrictions = &credentials.RepositoryRestrictions{Enabled: false, AllowedRepositories: []string{}}
	}

	tflog.Debug(context.Background(), "Expanded Git credential", map[string]interface{}{
		"id":          gitCredential.ID,
		"name":        gitCredential.Name,
		"description": gitCredential.Description,
		"username":    username,
		// Don't log the password
		"repository_restrictions": gitCredential.RepositoryRestrictions,
	})

	return gitCredential
}

func setPlatformHubGitCredential(ctx context.Context, model *platformHubGitCredentialResourceModel, resource *platformhubgitcredential.PlatformHubGitCredential) {
	if resource == nil {
		tflog.Warn(ctx, "Resource is nil in setPlatformHubGitCredential")
		return
	}

	model.ID = types.StringValue(resource.GetID())
	model.Name = types.StringValue(resource.GetName())
	model.Description = types.StringValue(resource.Description)

	tflog.Debug(ctx, "Setting Platform Hub Git credential state", map[string]interface{}{
		"id":                      resource.GetID(),
		"name":                    resource.GetName(),
		"description":             resource.Description,
		"repository_restrictions": resource.RepositoryRestrictions,
	})

	if resource.Details != nil {
		if usernamePassword, ok := resource.Details.(*credentials.UsernamePassword); ok && usernamePassword != nil {
			model.Username = types.StringValue(usernamePassword.Username)
			// Note: We don't set the password here as it's sensitive and not returned by the API
		} else {
			tflog.Debug(ctx, "Git credential is not of type UsernamePassword", map[string]interface{}{
				"type": resource.Details.Type(),
			})
		}
	} else {
		tflog.Warn(ctx, "Resource Details is nil")
	}

	if resource.RepositoryRestrictions != nil {
		var allowedRepositories = make([]string, 0, len(resource.RepositoryRestrictions.AllowedRepositories))
		for _, id := range resource.RepositoryRestrictions.AllowedRepositories {
			allowedRepositories = append(allowedRepositories, id)
		}
		repositoriesList, _ := types.SetValueFrom(
			ctx,
			types.StringType,
			allowedRepositories,
		)

		model.RepositoryRestrictions = &gitCredentialRepositoryRestrictionResourceModel{
			Enabled:             types.BoolValue(resource.RepositoryRestrictions.Enabled),
			AllowedRepositories: repositoriesList,
		}
	} else { //Default to disabled if resource doesn't have it
		model.RepositoryRestrictions = &gitCredentialRepositoryRestrictionResourceModel{
			Enabled:             types.BoolValue(false),
			AllowedRepositories: types.SetValueMust(types.StringType, make([]attr.Value, 0)),
		}
	}

	tflog.Debug(ctx, "Git credential state set", map[string]interface{}{
		"id":                      model.ID.ValueString(),
		"name":                    model.Name.ValueString(),
		"description":             model.Description.ValueString(),
		"username":                model.Username.ValueString(),
		"repository_restrictions": model.RepositoryRestrictions,
	})
}
