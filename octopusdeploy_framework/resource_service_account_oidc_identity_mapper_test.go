package octopusdeploy_framework

import (
	"context"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/serviceaccounts"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func newOIDCModel(conditions types.List) *schemas.OIDCServiceAccountSchemaModel {
	return &schemas.OIDCServiceAccountSchemaModel{
		Name:             types.StringValue("identity"),
		ServiceAccountID: types.StringValue("Users-1"),
		Issuer:           types.StringValue("https://token.actions.githubusercontent.com"),
		Subject:          types.StringValue("repo:org/repo:*"),
		Audience:         types.StringNull(),
		ClaimConditions:  conditions,
	}
}

func TestMapServiceAccountOIDCModelToRequest_SendsNullClaimConditionsWhenUnset(t *testing.T) {
	for name, list := range map[string]types.List{
		"null":    types.ListNull(schemas.OIDCClaimConditionObjectType()),
		"unknown": types.ListUnknown(schemas.OIDCClaimConditionObjectType()),
	} {
		t.Run(name, func(t *testing.T) {
			identity, diags := mapServiceAccountOIDCModelToRequest(context.Background(), newOIDCModel(list))

			require.False(t, diags.HasError())
			require.Nil(t, identity.ClaimConditions, "a nil slice is sent as null and Octopus keeps the existing conditions")
		})
	}
}

func TestMapServiceAccountOIDCModelToRequest_SendsEmptyClaimConditionsWhenEmpty(t *testing.T) {
	list, _ := types.ListValueFrom(context.Background(), schemas.OIDCClaimConditionObjectType(), []schemas.OIDCClaimConditionModel{})

	identity, diags := mapServiceAccountOIDCModelToRequest(context.Background(), newOIDCModel(list))

	require.False(t, diags.HasError())
	require.NotNil(t, identity.ClaimConditions, "an empty list clears the conditions on Octopus")
	require.Empty(t, identity.ClaimConditions)
}

func TestClaimConditionsRoundTrip(t *testing.T) {
	ctx := context.Background()
	conditions := []serviceaccounts.OIDCIdentityClaimCondition{
		{Claim: "repository_owner", Values: []string{"OctopusDeploy"}},
		{Claim: "ref", Values: []string{"refs/heads/main", "refs/tags/*"}},
	}

	model := &schemas.OIDCServiceAccountSchemaModel{}
	identity := serviceaccounts.NewOIDCIdentity("Users-1", "identity", "https://issuer", "subject")
	identity.ClaimConditions = conditions
	require.False(t, updateServiceAccountOIDCModel(ctx, identity, model).HasError())

	request, diags := mapServiceAccountOIDCModelToRequest(ctx, model)

	require.False(t, diags.HasError())
	require.Equal(t, conditions, request.ClaimConditions)
}

func TestUpdateServiceAccountOIDCModel_StoresEmptyListWhenServerHasNoClaimConditions(t *testing.T) {
	identity := serviceaccounts.NewOIDCIdentity("Users-1", "identity", "https://issuer", "subject")
	model := &schemas.OIDCServiceAccountSchemaModel{}

	diags := updateServiceAccountOIDCModel(context.Background(), identity, model)

	require.False(t, diags.HasError())
	require.False(t, model.ClaimConditions.IsNull())
	require.Empty(t, model.ClaimConditions.Elements())
}

func TestCheckClaimConditionsApplied_ErrorsWhenServerDropsConditions(t *testing.T) {
	request := serviceaccounts.NewOIDCIdentity("Users-1", "identity", "https://issuer", "subject")
	request.ClaimConditions = []serviceaccounts.OIDCIdentityClaimCondition{{Claim: "ref", Values: []string{"main"}}}
	response := serviceaccounts.NewOIDCIdentity("Users-1", "identity", "https://issuer", "subject")

	require.True(t, checkClaimConditionsApplied(request, response).HasError())
	require.False(t, checkClaimConditionsApplied(response, response).HasError())
}
