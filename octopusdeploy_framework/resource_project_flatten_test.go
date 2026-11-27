package octopusdeploy_framework

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func templateWithDisplaySettings(name string, displaySettings map[string]string) actiontemplates.ActionTemplateParameter {
	defaultValue := core.NewPropertyValue("", false)
	return actiontemplates.ActionTemplateParameter{
		Name:            name,
		DefaultValue:    &defaultValue,
		DisplaySettings: displaySettings,
	}
}

func priorTemplateList(t *testing.T, displaySettings ...attr.Value) types.List {
	t.Helper()

	elements := make([]attr.Value, 0, len(displaySettings))
	for i, ds := range displaySettings {
		elements = append(elements, types.ObjectValueMust(getTemplateAttrTypes(), map[string]attr.Value{
			"id":               types.StringNull(),
			"name":             types.StringValue(string(rune('a' + i))),
			"label":            types.StringNull(),
			"help_text":        types.StringNull(),
			"default_value":    types.StringNull(),
			"display_settings": ds,
		}))
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getTemplateAttrTypes()}, elements)
}

func displaySettingsAt(t *testing.T, list types.List, index int) types.Map {
	t.Helper()

	require.Greater(t, len(list.Elements()), index)
	object, ok := list.Elements()[index].(types.Object)
	require.True(t, ok)

	value, ok := object.Attributes()["display_settings"].(types.Map)
	require.True(t, ok)

	return value
}

func TestFlattenTemplatesDisplaySettings(t *testing.T) {
	emptyMap := types.MapValueMust(types.StringType, map[string]attr.Value{})
	populatedMap := types.MapValueMust(types.StringType, map[string]attr.Value{
		"Octopus.ControlType": types.StringValue("SingleLineText"),
	})

	t.Run("server returns no display settings and prior is null", func(t *testing.T) {
		templates := []actiontemplates.ActionTemplateParameter{templateWithDisplaySettings("a", nil)}

		flattened := flattenTemplates(templates, priorTemplateList(t, types.MapNull(types.StringType)))

		require.True(t, displaySettingsAt(t, flattened, 0).IsNull())
	})

	t.Run("server returns no display settings and prior is an empty map", func(t *testing.T) {
		templates := []actiontemplates.ActionTemplateParameter{templateWithDisplaySettings("a", map[string]string{})}

		flattened := flattenTemplates(templates, priorTemplateList(t, emptyMap))

		require.Equal(t, emptyMap, displaySettingsAt(t, flattened, 0))
	})

	t.Run("server returns display settings", func(t *testing.T) {
		templates := []actiontemplates.ActionTemplateParameter{
			templateWithDisplaySettings("a", map[string]string{"Octopus.ControlType": "SingleLineText"}),
		}

		flattened := flattenTemplates(templates, priorTemplateList(t, types.MapNull(types.StringType)))

		require.Equal(t, populatedMap, displaySettingsAt(t, flattened, 0))
	})

	t.Run("prior list is null, as it is on import", func(t *testing.T) {
		templates := []actiontemplates.ActionTemplateParameter{templateWithDisplaySettings("a", nil)}

		flattened := flattenTemplates(templates, types.ListNull(types.ObjectType{AttrTypes: getTemplateAttrTypes()}))

		require.True(t, displaySettingsAt(t, flattened, 0).IsNull())
	})

	t.Run("prior list length does not match the server response", func(t *testing.T) {
		templates := []actiontemplates.ActionTemplateParameter{
			templateWithDisplaySettings("a", nil),
			templateWithDisplaySettings("b", nil),
		}

		flattened := flattenTemplates(templates, priorTemplateList(t, emptyMap))

		require.True(t, displaySettingsAt(t, flattened, 0).IsNull())
		require.True(t, displaySettingsAt(t, flattened, 1).IsNull())
	})

	t.Run("prior values are mirrored per element", func(t *testing.T) {
		templates := []actiontemplates.ActionTemplateParameter{
			templateWithDisplaySettings("a", nil),
			templateWithDisplaySettings("b", map[string]string{}),
		}

		flattened := flattenTemplates(templates, priorTemplateList(t, types.MapNull(types.StringType), emptyMap))

		require.True(t, displaySettingsAt(t, flattened, 0).IsNull())
		require.Equal(t, emptyMap, displaySettingsAt(t, flattened, 1))
	})

	t.Run("no templates", func(t *testing.T) {
		flattened := flattenTemplates(nil, types.ListNull(types.ObjectType{AttrTypes: getTemplateAttrTypes()}))

		require.True(t, flattened.IsNull())
	})
}
