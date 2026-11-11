package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/serviceaccounts"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ServiceAccountOIDCIdentity{}

type ServiceAccountOIDCIdentity struct {
	*Config
}

func NewServiceAccountOIDCIdentity() resource.Resource {
	return &ServiceAccountOIDCIdentity{}
}

func (s *ServiceAccountOIDCIdentity) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ServiceAccountOIDCIdentityResourceName)
}

func (s *ServiceAccountOIDCIdentity) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ServiceAccountOIDCIdentitySchema{}.GetResourceSchema()
}

func (s *ServiceAccountOIDCIdentity) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	s.Config = ResourceConfiguration(req, resp)
}
func (s *ServiceAccountOIDCIdentity) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.OIDCServiceAccountSchemaModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	identityRequest, diags := mapServiceAccountOIDCModelToRequest(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Octopus Server versions before 2026.4.9767 ignore claim conditions.
	if len(identityRequest.ClaimConditions) > 0 && !s.Config.IsVersionSameOrGreaterThan("2026.4.9767") {
		resp.Diagnostics.AddAttributeError(path.Root("claim_conditions"), "Octopus Server doesn't support claim conditions",
			"Claim conditions need Octopus Server 2026.4.9767 or later. Upgrade Octopus Server or remove claim_conditions from the configuration.")
		return
	}
	identityCreateResponse, err := serviceaccounts.AddOIDCIdentity(s.Client, identityRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating OIDC identity", err.Error())
		return
	}
	identityResponse, err := serviceaccounts.GetOIDCIdentityByID(s.Client, identityRequest.ServiceAccountID, identityCreateResponse.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error creating OIDC identity", err.Error())
		return
	}

	resp.Diagnostics.Append(updateServiceAccountOIDCModel(ctx, identityResponse, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (s *ServiceAccountOIDCIdentity) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.OIDCServiceAccountSchemaModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityResponse, err := serviceaccounts.GetOIDCIdentityByID(s.Client, state.ServiceAccountID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "service account OIDC identity"); err != nil {
			resp.Diagnostics.AddError("Error reading service account OIDC identity", err.Error())
		}
		return
	}

	resp.Diagnostics.Append(updateServiceAccountOIDCModel(ctx, identityResponse, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (s *ServiceAccountOIDCIdentity) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.OIDCServiceAccountSchemaModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityRequest, diags := mapServiceAccountOIDCModelToRequest(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The plan copies claim conditions from state when the config leaves them out.
	// Send null instead so Octopus keeps the conditions it has now.
	var configClaimConditions types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("claim_conditions"), &configClaimConditions)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if configClaimConditions.IsNull() {
		identityRequest.ClaimConditions = nil
	}

	// Octopus Server versions before 2026.4.9767 ignore claim conditions.
	if len(identityRequest.ClaimConditions) > 0 && !s.Config.IsVersionSameOrGreaterThan("2026.4.9767") {
		resp.Diagnostics.AddAttributeError(path.Root("claim_conditions"), "Octopus Server doesn't support claim conditions",
			"Claim conditions need Octopus Server 2026.4.9767 or later. Upgrade Octopus Server or remove claim_conditions from the configuration.")
		return
	}

	err := serviceaccounts.UpdateOIDCIdentity(s.Client, identityRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error updating service account OIDC identity", err.Error())
		return
	}
	identityResponse, err := serviceaccounts.GetOIDCIdentityByID(s.Client, identityRequest.ServiceAccountID, identityRequest.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error creating OIDC identity", err.Error())
		return
	}

	resp.Diagnostics.Append(updateServiceAccountOIDCModel(ctx, identityResponse, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (s *ServiceAccountOIDCIdentity) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.OIDCServiceAccountSchemaModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := serviceaccounts.DeleteOIDCIdentityByID(s.Client, state.ServiceAccountID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting service account OIDC identity", err.Error())
		return
	}
}

func mapServiceAccountOIDCModelToRequest(ctx context.Context, model *schemas.OIDCServiceAccountSchemaModel) (*serviceaccounts.OIDCIdentity, diag.Diagnostics) {
	identity := serviceaccounts.NewOIDCIdentity(model.ServiceAccountID.ValueString(), model.Name.ValueString(), model.Issuer.ValueString(), model.Subject.ValueString())
	identity.ID = model.ID.ValueString()
	identity.Name = model.Name.ValueString()

	if !model.Audience.IsNull() {
		identity.Audience = model.Audience.ValueString()
	}

	conditions, diags := mapClaimConditionsToRequest(ctx, model.ClaimConditions)
	identity.ClaimConditions = conditions

	return identity, diags
}

// mapClaimConditionsToRequest returns nil when the list is null or unknown.
// Octopus keeps the existing conditions when the field is null. An empty list clears them.
func mapClaimConditionsToRequest(ctx context.Context, list types.List) ([]serviceaccounts.OIDCIdentityClaimCondition, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []schemas.OIDCClaimConditionModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]serviceaccounts.OIDCIdentityClaimCondition, 0, len(models))
	for _, model := range models {
		values := make([]string, 0, len(model.Values.Elements()))
		diags.Append(model.Values.ElementsAs(ctx, &values, false)...)
		result = append(result, serviceaccounts.OIDCIdentityClaimCondition{
			Claim:  model.Claim.ValueString(),
			Values: values,
		})
	}
	return result, diags
}

// mapClaimConditionsToState returns an empty list rather than null when the server has no conditions.
// A computed attribute must be known after apply.
func mapClaimConditionsToState(ctx context.Context, conditions []serviceaccounts.OIDCIdentityClaimCondition) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	models := make([]schemas.OIDCClaimConditionModel, 0, len(conditions))
	for _, condition := range conditions {
		values, valueDiags := types.ListValueFrom(ctx, types.StringType, condition.Values)
		diags.Append(valueDiags...)
		models = append(models, schemas.OIDCClaimConditionModel{
			Claim:  types.StringValue(condition.Claim),
			Values: values,
		})
	}
	list, listDiags := types.ListValueFrom(ctx, schemas.OIDCClaimConditionObjectType(), models)
	diags.Append(listDiags...)
	return list, diags
}

func updateServiceAccountOIDCModel(ctx context.Context, request *serviceaccounts.OIDCIdentity, model *schemas.OIDCServiceAccountSchemaModel) diag.Diagnostics {
	model.Name = types.StringValue(request.Name)
	model.Issuer = types.StringValue(request.Issuer)
	model.Subject = types.StringValue(request.Subject)
	model.ID = types.StringValue(request.ID)
	model.ServiceAccountID = types.StringValue(request.ServiceAccountID)

	if request.Audience == "" {
		model.Audience = types.StringNull()
	} else {
		model.Audience = types.StringValue(request.Audience)
	}

	conditions, diags := mapClaimConditionsToState(ctx, request.ClaimConditions)
	model.ClaimConditions = conditions
	return diags
}
