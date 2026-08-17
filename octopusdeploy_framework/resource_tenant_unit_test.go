package octopusdeploy_framework

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenTenantTags(t *testing.T) {
	emptySet := types.SetValueMust(types.StringType, []attr.Value{})
	populatedSet := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Region/US")})

	cases := []struct {
		name     string
		tags     []string
		current  types.Set
		expected types.Set
	}{
		{
			name:     "tags from the api win over state",
			tags:     []string{"Region/EU"},
			current:  populatedSet,
			expected: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Region/EU")}),
		},
		{
			name:     "tags from the api populate an empty set",
			tags:     []string{"Region/US"},
			current:  emptySet,
			expected: populatedSet,
		},
		{
			name:     "no tags keeps an empty set empty",
			tags:     nil,
			current:  emptySet,
			expected: emptySet,
		},
		{
			name:     "no tags clears a populated set",
			tags:     nil,
			current:  populatedSet,
			expected: emptySet,
		},
		{
			name:     "no tags keeps a null set null",
			tags:     nil,
			current:  types.SetNull(types.StringType),
			expected: types.SetNull(types.StringType),
		},
		{
			name:     "no tags resolves an unknown set to empty",
			tags:     []string{},
			current:  types.SetUnknown(types.StringType),
			expected: emptySet,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := flattenTenantTags(tc.tags, tc.current)
			if !actual.Equal(tc.expected) {
				t.Errorf("expected %s, got %s", tc.expected, actual)
			}
			if actual.IsUnknown() {
				t.Error("flattenTenantTags must never return an unknown value")
			}
		})
	}
}
