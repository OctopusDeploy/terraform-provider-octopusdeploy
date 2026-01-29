package schemas

import (
	"flag"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

var createSharedContainer = flag.Bool("createSharedContainer", false, "Set to true to run integration tests in containers")

func TestExpandVariableScope(t *testing.T) {
	// scope := ExpandVariableScopes(nil)
	scope := MapToVariableScope(types.ListNull(types.ObjectType{AttrTypes: VariableScopeObjectType()}))
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, variables.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	// flattenedVariableScope := []interface{}{}
	flattenedVariableScope := types.ListValueMust(types.ObjectType{AttrTypes: VariableScopeObjectType()}, []attr.Value{})
	scope = MapToVariableScope(flattenedVariableScope)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, variables.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	// flattenedVariableScope = []interface{}{nil}
	flattenedVariableScope = types.ListValueMust(types.ObjectType{AttrTypes: VariableScopeObjectType()}, []attr.Value{types.ObjectNull(VariableScopeObjectType())})
	scope = MapToVariableScope(flattenedVariableScope)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, variables.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	//flattenedVariableScope = []interface{}{"foo"}
	// scope = MapToVariableScope(flattenedVariableScope)
	// assert.True(t, scope.IsEmpty())
	// assert.Equal(t, variables.VariableScope{}, scope)
	// assert.Empty(t, scope.Channels)

	// flattenedVariableScope = []interface{}{map[string]interface{}{}}
	flattenedVariableScope = types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{}}, []attr.Value{})
	scope = MapToVariableScope(flattenedVariableScope)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, variables.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	// flattenedVariableScope = []interface{}{map[string]interface{}{
	// 	"actions": []interface{}{"foo"},
	// }}
	flattenedVariableScope = types.ListValueMust(
		types.ObjectType{AttrTypes: map[string]attr.Type{"actions": types.ListType{ElemType: types.StringType}}},
		[]attr.Value{basetypes.NewObjectValueMust(
			map[string]attr.Type{"actions": types.ListType{ElemType: types.StringType}},
			map[string]attr.Value{"actions": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("foo")})},
		)})
	expectedScope := variables.VariableScope{
		Actions: []string{"foo"},
	}
	scope = MapToVariableScope(flattenedVariableScope)
	assert.False(t, scope.IsEmpty())
	assert.Equal(t, expectedScope, scope)
	assert.NotEmpty(t, scope.Actions)
	assert.Empty(t, scope.Channels)
}

func TestFlattenVariableScope(t *testing.T) {
	scopes := basetypes.NewObjectValueMust(
		map[string]attr.Type{
			"actions": types.ListType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"actions": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("foo")}),
		},
	)

	flattenedVariableScope := types.ListValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"actions": types.ListType{ElemType: types.StringType},
			},
		},
		[]attr.Value{scopes},
	)
	expectedScope := variables.VariableScope{
		Actions: []string{"foo"},
	}

	scope := MapToVariableScope(flattenedVariableScope)
	assert.False(t, scope.IsEmpty())
	assert.Equal(t, expectedScope, scope)
	assert.NotEmpty(t, scope.Actions)
	assert.Empty(t, scope.Channels)

	flattenedVariableScope = types.ListValueMust(types.ObjectType{AttrTypes: VariableScopeObjectType()}, []attr.Value{MapFromVariableScope(scope)})

	assert.NotNil(t, flattenedVariableScope)
	assert.Len(t, flattenedVariableScope.Elements(), 1)
	actionScope := flattenedVariableScope.Elements()[0].(types.Object).Attributes()["actions"].(types.List)
	t.Logf("Action scope: %#v", actionScope)
	assert.Len(t, actionScope.Elements(), 1)
}

func TestVariableScopeOrderInsensitive(t *testing.T) {
	originalScope := variables.VariableScope{
		Environments: []string{"Environments-62", "Environments-61", "Environments-63"},
	}

	reorderedScope := variables.VariableScope{
		Environments: []string{"Environments-61", "Environments-63", "Environments-62"},
	}

	originalTerraform := MapFromVariableScope(originalScope)
	reorderedTerraform := MapFromVariableScope(reorderedScope)

	// These are not equal because the order is different, but they are equivalent
	assert.NotEqual(t, originalTerraform, reorderedTerraform)

	originalEnvs := originalTerraform.(types.Object).Attributes()["environments"].(types.List)
	reorderedEnvs := reorderedTerraform.(types.Object).Attributes()["environments"].(types.List)

	originalStrings := util.ExpandStringList(originalEnvs)
	reorderedStrings := util.ExpandStringList(reorderedEnvs)

	assert.ElementsMatch(t, originalStrings, reorderedStrings)
}
