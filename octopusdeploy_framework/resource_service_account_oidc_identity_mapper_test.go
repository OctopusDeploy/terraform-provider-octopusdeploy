package octopusdeploy_framework

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/serviceaccounts"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestMapServiceAccountOIDCModelToRequest_SendsEmptyClaimConditionsWhenUnset(t *testing.T) {
	model := &schemas.OIDCServiceAccountSchemaModel{
		Name:             types.StringValue("identity"),
		ServiceAccountID: types.StringValue("Users-1"),
		Issuer:           types.StringValue("https://token.actions.githubusercontent.com"),
		Subject:          types.StringValue("repo:org/repo:*"),
		Audience:         types.StringNull(),
	}

	identity := mapServiceAccountOIDCModelToRequest(model)

	require.NotNil(t, identity.ClaimConditions, "a nil slice is sent as null and Octopus would keep the old conditions")
	require.Empty(t, identity.ClaimConditions)
}

func TestMapServiceAccountOIDCModelToRequest_MapsClaimConditions(t *testing.T) {
	model := &schemas.OIDCServiceAccountSchemaModel{
		Name:             types.StringValue("identity"),
		ServiceAccountID: types.StringValue("Users-1"),
		Issuer:           types.StringValue("https://token.actions.githubusercontent.com"),
		Subject:          types.StringValue("repo:org/repo:*"),
		Audience:         types.StringNull(),
		ClaimConditions: []schemas.OIDCClaimConditionModel{
			{
				Claim:  types.StringValue("repository_owner"),
				Values: []types.String{types.StringValue("OctopusDeploy")},
			},
			{
				Claim:  types.StringValue("ref"),
				Values: []types.String{types.StringValue("refs/heads/main"), types.StringValue("refs/tags/*")},
			},
		},
	}

	identity := mapServiceAccountOIDCModelToRequest(model)

	require.Equal(t, []serviceaccounts.OIDCIdentityClaimCondition{
		{Claim: "repository_owner", Values: []string{"OctopusDeploy"}},
		{Claim: "ref", Values: []string{"refs/heads/main", "refs/tags/*"}},
	}, identity.ClaimConditions)
}

func TestUpdateServiceAccountOIDCModel_StoresNullWhenServerHasNoClaimConditions(t *testing.T) {
	identity := serviceaccounts.NewOIDCIdentity("Users-1", "identity", "https://issuer", "subject")
	identity.ClaimConditions = []serviceaccounts.OIDCIdentityClaimCondition{}
	model := &schemas.OIDCServiceAccountSchemaModel{}

	updateServiceAccountOIDCModel(identity, model)

	require.Nil(t, model.ClaimConditions)
}

func TestUpdateServiceAccountOIDCModel_StoresClaimConditions(t *testing.T) {
	identity := serviceaccounts.NewOIDCIdentity("Users-1", "identity", "https://issuer", "subject")
	identity.ClaimConditions = []serviceaccounts.OIDCIdentityClaimCondition{
		{Claim: "environment", Values: []string{"prod", "staging"}},
	}
	model := &schemas.OIDCServiceAccountSchemaModel{}

	updateServiceAccountOIDCModel(identity, model)

	require.Equal(t, []schemas.OIDCClaimConditionModel{
		{
			Claim:  types.StringValue("environment"),
			Values: []types.String{types.StringValue("prod"), types.StringValue("staging")},
		},
	}, model.ClaimConditions)
}
