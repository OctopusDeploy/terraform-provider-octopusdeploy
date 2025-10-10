package octopusdeploy_framework

import (
	"context"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/userroles"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestFilterUserRolesByPreviousState(t *testing.T) {
	ctx := context.Background()

	t.Run("ShouldReturnEmptyListForEmptyPreviousState", func(t *testing.T) {
		role1 := userroles.NewScopedUserRole("user-role-1")
		role1.ID = "role-1"
		role2 := userroles.NewScopedUserRole("user-role-2")
		role2.ID = "role-2"

		serverRoles := []*userroles.ScopedUserRole{role1, role2}

		result := filterUserRolesByPreviousState(ctx, serverRoles, types.SetNull(userRoleObjectType))
		assert.Empty(t, result, "Should return empty list when previous state is null")

		result = filterUserRolesByPreviousState(ctx, serverRoles, types.SetUnknown(userRoleObjectType))
		assert.Equal(t, serverRoles, result, "Should return all server roles when previous state is unknown (creation scenario)")
	})

	t.Run("ShouldIncludeAllServerRolesForCreationWithUnknownIds", func(t *testing.T) {
		role1 := userroles.NewScopedUserRole("user-role-1")
		role1.ID = "role-1"

		serverRoles := []*userroles.ScopedUserRole{role1}

		previousState := types.SetValueMust(userRoleObjectType, []attr.Value{
			types.ObjectValueMust(userRoleObjectType.AttrTypes, map[string]attr.Value{
				"id":                types.StringUnknown(), // Unknown ID during creation
				"user_role_id":      types.StringValue("user-role-1"),
				"space_id":          types.StringValue("Spaces-1"),
				"team_id":           types.StringUnknown(),
				"environment_ids":   types.SetNull(types.StringType),
				"project_group_ids": types.SetNull(types.StringType),
				"project_ids":       types.SetNull(types.StringType),
				"tenant_ids":        types.SetNull(types.StringType),
			}),
		})

		result := filterUserRolesByPreviousState(ctx, serverRoles, previousState)
		assert.Equal(t, serverRoles, result, "Should return all server roles when IDs are unknown (creation)")
	})

	t.Run("ShouldFilterPreviouslyManagedRoles", func(t *testing.T) {
		role1 := userroles.NewScopedUserRole("user-role-1")
		role1.ID = "role-1"
		role2 := userroles.NewScopedUserRole("user-role-2")
		role2.ID = "role-2"
		role3 := userroles.NewScopedUserRole("user-role-3")
		role3.ID = "role-3"

		serverRoles := []*userroles.ScopedUserRole{role1, role2, role3}

		previousState := types.SetValueMust(userRoleObjectType, []attr.Value{
			types.ObjectValueMust(userRoleObjectType.AttrTypes, map[string]attr.Value{
				"id":                types.StringValue("role-1"),
				"user_role_id":      types.StringValue("user-role-1"),
				"space_id":          types.StringValue("Spaces-1"),
				"team_id":           types.StringValue("Teams-1"),
				"environment_ids":   types.SetNull(types.StringType),
				"project_group_ids": types.SetNull(types.StringType),
				"project_ids":       types.SetNull(types.StringType),
				"tenant_ids":        types.SetNull(types.StringType),
			}),
			types.ObjectValueMust(userRoleObjectType.AttrTypes, map[string]attr.Value{
				"id":                types.StringValue("role-3"),
				"user_role_id":      types.StringValue("user-role-3"),
				"space_id":          types.StringValue("Spaces-1"),
				"team_id":           types.StringValue("Teams-1"),
				"environment_ids":   types.SetNull(types.StringType),
				"project_group_ids": types.SetNull(types.StringType),
				"project_ids":       types.SetNull(types.StringType),
				"tenant_ids":        types.SetNull(types.StringType),
			}),
		})

		result := filterUserRolesByPreviousState(ctx, serverRoles, previousState)

		assert.Len(t, result, 2, "Should return exactly 2 roles")

		resultIDs := make(map[string]bool)
		for _, role := range result {
			resultIDs[role.ID] = true
		}

		assert.True(t, resultIDs["role-1"], "Should include role-1 (was in previous state)")
		assert.False(t, resultIDs["role-2"], "Should exclude role-2 (was not in previous state)")
		assert.True(t, resultIDs["role-3"], "Should include role-3 (was in previous state)")
	})

	t.Run("ShouldFilterOutNewServerRoles", func(t *testing.T) {
		externalRole := userroles.NewScopedUserRole("external-user-role")
		externalRole.ID = "external-role" // This could be from a standalone resource

		serverRoles := []*userroles.ScopedUserRole{externalRole}

		previousState := types.SetValueMust(userRoleObjectType, []attr.Value{})

		result := filterUserRolesByPreviousState(ctx, serverRoles, previousState)
		assert.Empty(t, result, "Should filter out roles that were not previously managed by the team")
	})
}
