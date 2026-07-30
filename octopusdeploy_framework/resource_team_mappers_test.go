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

// TestUpdateUserRolesRemovalOwnership guards scopedUserRolesToRemove, the exact removal-candidate selection
// updateUserRoles uses. Regression test for #180: updating a team must not delete externally-managed scoped
// user roles, and must still delete roles it previously owned that are dropped from config.
func TestUpdateUserRolesRemovalOwnership(t *testing.T) {
	newRole := func(id, userRoleID string) *userroles.ScopedUserRole {
		r := userroles.NewScopedUserRole(userRoleID)
		r.ID = id
		return r
	}

	previousStateWithIDs := func(ids ...string) types.Set {
		elems := make([]attr.Value, 0, len(ids))
		for _, id := range ids {
			elems = append(elems, types.ObjectValueMust(userRoleObjectType.AttrTypes, map[string]attr.Value{
				"id":                types.StringValue(id),
				"user_role_id":      types.StringValue("user-role-for-" + id),
				"space_id":          types.StringValue("Spaces-1"),
				"team_id":           types.StringValue("Teams-1"),
				"environment_ids":   types.SetNull(types.StringType),
				"project_group_ids": types.SetNull(types.StringType),
				"project_ids":       types.SetNull(types.StringType),
				"tenant_ids":        types.SetNull(types.StringType),
			}))
		}
		return types.SetValueMust(userRoleObjectType, elems)
	}

	removalIDs := func(newUserRoles, serverRoles []*userroles.ScopedUserRole, previous types.Set) map[string]bool {
		toRemove := scopedUserRolesToRemove(newUserRoles, serverRoles, previous)
		ids := make(map[string]bool)
		for _, r := range toRemove {
			ids[r.ID] = true
		}
		return ids
	}

	t.Run("ShouldNotRemoveExternallyManagedRolesWhenTeamHasNoInlineRoles", func(t *testing.T) {
		// Team resource declares no user_role blocks (newUserRoles empty) and never managed any
		// (previous state null). All roles on the server belong to standalone resources.
		external := []*userroles.ScopedUserRole{newRole("role-1", "ur-1"), newRole("role-2", "ur-2")}

		ids := removalIDs(nil, external, types.SetNull(userRoleObjectType))
		assert.Empty(t, ids, "must not delete scoped user roles managed outside the team resource")
	})

	t.Run("ShouldRemoveOnlyPreviouslyOwnedRolesDroppedFromConfig", func(t *testing.T) {
		// Inline user_role user previously owned role-a and role-b, now keeps only role-a.
		// role-ext is a standalone-managed role that also lives on the team.
		roleA := newRole("role-a", "ur-a")
		roleB := newRole("role-b", "ur-b")
		roleExt := newRole("role-ext", "ur-ext")

		newUserRoles := []*userroles.ScopedUserRole{roleA}
		serverRoles := []*userroles.ScopedUserRole{roleA, roleB, roleExt}

		ids := removalIDs(newUserRoles, serverRoles, previousStateWithIDs("role-a", "role-b"))

		assert.True(t, ids["role-b"], "should remove role-b (previously owned, dropped from config)")
		assert.False(t, ids["role-ext"], "must not remove externally-managed role-ext")
		assert.False(t, ids["role-a"], "must not remove role-a (still in config)")
		assert.Len(t, ids, 1)
	})
}
