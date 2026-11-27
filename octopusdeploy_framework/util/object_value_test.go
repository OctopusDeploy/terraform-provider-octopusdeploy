package util

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func nestedObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"value":        types.StringType,
		"display_name": types.StringType,
	}
}

func sampleObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"control_type": types.StringType,
		"enabled":      types.BoolType,
		"select_option": types.ListType{
			ElemType: types.ObjectType{AttrTypes: nestedObjectType()},
		},
		"tags": types.SetType{ElemType: types.StringType},
	}
}

func TestObjectValueFillsMissingAttributesWithNull(t *testing.T) {
	result := ObjectValue(sampleObjectType(), map[string]attr.Value{
		"control_type": types.StringValue("Select"),
	})

	attrs := result.Attributes()
	require.Len(t, attrs, 4)
	require.Equal(t, types.StringValue("Select"), attrs["control_type"])
	require.True(t, attrs["enabled"].IsNull())
	require.True(t, attrs["select_option"].IsNull())
	require.True(t, attrs["tags"].IsNull())
}

func TestObjectValueKeepsSuppliedAttributes(t *testing.T) {
	selectOptions := types.ListValueMust(
		types.ObjectType{AttrTypes: nestedObjectType()},
		[]attr.Value{
			types.ObjectValueMust(nestedObjectType(), map[string]attr.Value{
				"value":        types.StringValue("Value-1"),
				"display_name": types.StringValue("Name-1"),
			}),
		},
	)

	result := ObjectValue(sampleObjectType(), map[string]attr.Value{
		"control_type":  types.StringValue("Select"),
		"enabled":       types.BoolValue(true),
		"select_option": selectOptions,
		"tags":          types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
	})

	attrs := result.Attributes()
	require.Equal(t, types.BoolValue(true), attrs["enabled"])
	require.Equal(t, selectOptions, attrs["select_option"])
	require.False(t, attrs["tags"].IsNull())
}

func TestObjectValueDoesNotMutateCallerMap(t *testing.T) {
	attrs := map[string]attr.Value{"control_type": types.StringValue("Select")}

	ObjectValue(sampleObjectType(), attrs)

	require.Len(t, attrs, 1, "the supplied map should be left untouched")
}

func TestObjectValueStillRejectsUndeclaredAttributes(t *testing.T) {
	require.Panics(t, func() {
		ObjectValue(sampleObjectType(), map[string]attr.Value{
			"control_type": types.StringValue("Select"),
			"not_declared": types.StringValue("boom"),
		})
	}, "an attribute absent from the object type is a mapper bug and should stay loud")
}

func TestNullValueOfNestedTypes(t *testing.T) {
	for name, attrType := range sampleObjectType() {
		t.Run(name, func(t *testing.T) {
			value := NullValueOf(attrType)
			require.True(t, value.IsNull())
			require.Equal(t, attrType, value.Type(nil))
		})
	}
}
