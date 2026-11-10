package octopusdeploy_framework

import (
	"context"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/serviceaccounts"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestMapServiceAccountOIDCModelToRequest_SendsNullClaimConditionsWhenUnset(t *testing.T) {
	for name, list := range map[string]types.List{
		"null":    types.ListNull(schemas.OIDCClaimConditionObjectType()),
		"unknown": types.ListUnknown(schemas.OIDCClaimConditionObjectType()),
	} {
		t.Run(name, func(t *testing.T) {
			model := &schemas.OIDCServiceAccountSchemaModel{
				Name:             types.StringValue("identity"),
				ServiceAccountID: types.StringValue("Users-1"),
				Issuer:           types.StringValue("https://token.actions.githubusercontent.com"),
				Subject:          types.StringValue("repo:org/repo:*"),
				Audience:         types.StringNull(),
				ClaimConditions:  list,
			}

			identity, diags := mapServiceAccountOIDCModelToRequest(context.Background(), model)

			require.False(t, diags.HasError())
			require.Nil(t, identity.ClaimConditions, "a nil slice is sent as null and Octopus keeps the existing conditions")
		})
	}
}

func TestCheckClaimConditionsApplied_ErrorsWhenServerDropsConditions(t *testing.T) {
	request := serviceaccounts.NewOIDCIdentity("Users-1", "identity", "https://issuer", "subject")
	request.ClaimConditions = []serviceaccounts.OIDCIdentityClaimCondition{{Claim: "ref", Values: []string{"main"}}}
	response := serviceaccounts.NewOIDCIdentity("Users-1", "identity", "https://issuer", "subject")

	require.True(t, checkClaimConditionsApplied(request, response).HasError())
	require.False(t, checkClaimConditionsApplied(response, response).HasError())
}
