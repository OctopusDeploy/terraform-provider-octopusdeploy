package octopusdeploy

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/require"
)

func TestMachineUpdatePolicyKubernetesAgentUpdateBehaviorValidation(t *testing.T) {
	validate := getMachineUpdatePolicySchema()["kubernetes_agent_update_behavior"].ValidateDiagFunc
	path := cty.GetAttrPath("kubernetes_agent_update_behavior")

	for _, value := range []string{"NeverUpdate", "Update", "Block"} {
		diags := validate(value, path)
		require.False(t, diags.HasError(), "expected %s to be valid", value)
	}

	diags := validate("Blocked", path)
	require.True(t, diags.HasError(), "expected Blocked to be invalid")
}
