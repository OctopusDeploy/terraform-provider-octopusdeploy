package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/userroles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type teamResource struct {
	*Config
}

func NewTeamResource() resource.Resource {
	return &teamResource{}
}

var _ resource.ResourceWithImportState = &teamResource{}

func (r *teamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("team")
}

func (r *teamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.TeamSchema{}.GetResourceSchema()
}

func (r *teamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *teamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.TeamModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating team", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	team := mapTeamStateToResource(ctx, plan)
	createdTeam, err := r.Client.Teams.Add(team)
	if err != nil {
		resp.Diagnostics.AddError("Error creating team", err.Error())
		return
	}

	if err := r.updateUserRoles(ctx, plan, createdTeam); err != nil {
		resp.Diagnostics.AddError("Error updating user roles for team", err.Error())
		return
	}

	state := mapTeamResourceToState(ctx, createdTeam, plan, r.Client)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *teamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.TeamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := r.Client.Teams.GetByID(state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "teamResource"); err != nil {
			resp.Diagnostics.AddError("unable to load team", err.Error())
		}
		return
	}

	newState := mapTeamResourceToState(ctx, team, state, r.Client)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *teamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.TeamModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team := mapTeamStateToResource(ctx, plan)
	updatedTeam, err := r.Client.Teams.Update(team)
	if err != nil {
		resp.Diagnostics.AddError("Error updating team", err.Error())
		return
	}

	if err := r.updateUserRoles(ctx, plan, updatedTeam); err != nil {
		resp.Diagnostics.AddError("Error updating user roles for team", err.Error())
		return
	}

	state := mapTeamResourceToState(ctx, updatedTeam, plan, r.Client)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.TeamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.Client.Teams.DeleteByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting team", err.Error())
		return
	}
}

func (*teamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *teamResource) updateUserRoles(ctx context.Context, model schemas.TeamModel, team *teams.Team) error {
	newUserRoles := mapUserRoleSetStateToResource(ctx, team, model.UserRole)

	existingUserRoles, err := r.Client.Teams.GetScopedUserRoles(*team, core.SkipTakeQuery{})
	if err != nil {
		return fmt.Errorf("error getting existing user roles for team %s: %s", team.ID, err)
	}

	userRolesToAdd := findAddedScopedUserRoles(newUserRoles, existingUserRoles.Items)
	userRolesToRemove := findRemovedScopedUserRoles(newUserRoles, existingUserRoles.Items)
	userRolesToEdit := findModifiedScopedUserRoles(newUserRoles, existingUserRoles.Items)

	for _, userRole := range userRolesToAdd {
		_, err := r.Client.ScopedUserRoles.Add(userRole)
		if err != nil {
			return fmt.Errorf("error creating user role for team %s: %s", team.ID, err)
		}
	}

	for _, userRole := range userRolesToRemove {
		if err := r.Client.ScopedUserRoles.DeleteByID(userRole.ID); err != nil {
			return fmt.Errorf("error deleting user role %s from team %s: %s", userRole.ID, team.ID, err)
		}
	}

	for _, userRole := range userRolesToEdit {
		_, err := r.Client.ScopedUserRoles.Update(userRole)
		if err != nil {
			return fmt.Errorf("error updating user role for team %s: %s", team.ID, err)
		}
	}

	return nil
}

// selectScopedUserRoleIds maps a slice of scoped user roles to their corresponding IDs
func selectScopedUserRoleIds(roles []*userroles.ScopedUserRole) []string {
	return util.SliceTransform(roles, func(role *userroles.ScopedUserRole) string {
		return role.ID
	})
}

// findAddedScopedUserRoles returns roles from the first slice that don't exist in the second slice
func findAddedScopedUserRoles(newRoles, existingRoles []*userroles.ScopedUserRole) []*userroles.ScopedUserRole {
	existingUserRoleIds := selectScopedUserRoleIds(existingRoles)

	return util.SliceFilter(newRoles, func(newRole *userroles.ScopedUserRole) bool {
		return !util.SliceContains(existingUserRoleIds, newRole.ID)
	})
}

// findRemovedScopedUserRoles returns roles from the first slice that don't exist in the second slice
func findRemovedScopedUserRoles(newRoles, existingRoles []*userroles.ScopedUserRole) []*userroles.ScopedUserRole {
	newUserRoleIds := selectScopedUserRoleIds(newRoles)

	return util.SliceFilter(existingRoles, func(existingRole *userroles.ScopedUserRole) bool {
		return !util.SliceContains(newUserRoleIds, existingRole.ID)
	})
}

// findModifiedScopedUserRoles returns roles from the first slice that have matching IDs in the second slice
// but with at least one property that differs
func findModifiedScopedUserRoles(newRoles, existingRoles []*userroles.ScopedUserRole) []*userroles.ScopedUserRole {
	return util.SliceFilter(newRoles, func(newRole *userroles.ScopedUserRole) bool {
		existingRole := util.SliceFind(existingRoles, func(existing *userroles.ScopedUserRole) bool {
			return existing.ID == newRole.ID
		})
		return existingRole != nil && isScopedUserRoleModified(newRole, *existingRole)
	})
}

func isScopedUserRoleModified(a, b *userroles.ScopedUserRole) bool {
	return a.UserRoleID != b.UserRoleID ||
		a.SpaceID != b.SpaceID ||
		a.TeamID != b.TeamID ||
		!util.StringSlicesEqual(a.EnvironmentIDs, b.EnvironmentIDs) ||
		!util.StringSlicesEqual(a.ProjectGroupIDs, b.ProjectGroupIDs) ||
		!util.StringSlicesEqual(a.ProjectIDs, b.ProjectIDs) ||
		!util.StringSlicesEqual(a.TenantIDs, b.TenantIDs)
}
