package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/serviceaccounts"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
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
	identityRequest := mapServiceAccountOIDCModelToRequest(&plan)
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

	updateServiceAccountOIDCModel(identityResponse, &plan)
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

	updateServiceAccountOIDCModel(identityResponse, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (s *ServiceAccountOIDCIdentity) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.OIDCServiceAccountSchemaModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityRequest := mapServiceAccountOIDCModelToRequest(&plan)

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

	updateServiceAccountOIDCModel(identityResponse, &plan)
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

func mapServiceAccountOIDCModelToRequest(model *schemas.OIDCServiceAccountSchemaModel) *serviceaccounts.OIDCIdentity {
	identity := serviceaccounts.NewOIDCIdentity(model.ServiceAccountID.ValueString(), model.Name.ValueString(), model.Issuer.ValueString(), model.Subject.ValueString())
	identity.ID = model.ID.ValueString()
	identity.Name = model.Name.ValueString()

	if !model.Audience.IsNull() {
		identity.Audience = model.Audience.ValueString()
	}

	identity.ClaimConditions = mapClaimConditionsToRequest(model.ClaimConditions)

	return identity
}

// mapClaimConditionsToRequest always returns a non-nil slice.
// Octopus keeps the existing conditions when the field is null, so an empty list is the only way to clear them.
func mapClaimConditionsToRequest(conditions []schemas.OIDCClaimConditionModel) []serviceaccounts.OIDCIdentityClaimCondition {
	result := make([]serviceaccounts.OIDCIdentityClaimCondition, 0, len(conditions))
	for _, condition := range conditions {
		values := make([]string, 0, len(condition.Values))
		for _, value := range condition.Values {
			values = append(values, value.ValueString())
		}
		result = append(result, serviceaccounts.OIDCIdentityClaimCondition{
			Claim:  condition.Claim.ValueString(),
			Values: values,
		})
	}
	return result
}

// mapClaimConditionsToState returns nil when there are no conditions, so the attribute stays null in state.
func mapClaimConditionsToState(conditions []serviceaccounts.OIDCIdentityClaimCondition) []schemas.OIDCClaimConditionModel {
	if len(conditions) == 0 {
		return nil
	}
	result := make([]schemas.OIDCClaimConditionModel, 0, len(conditions))
	for _, condition := range conditions {
		values := make([]types.String, 0, len(condition.Values))
		for _, value := range condition.Values {
			values = append(values, types.StringValue(value))
		}
		result = append(result, schemas.OIDCClaimConditionModel{
			Claim:  types.StringValue(condition.Claim),
			Values: values,
		})
	}
	return result
}

func updateServiceAccountOIDCModel(request *serviceaccounts.OIDCIdentity, model *schemas.OIDCServiceAccountSchemaModel) {
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

	model.ClaimConditions = mapClaimConditionsToState(request.ClaimConditions)
}
