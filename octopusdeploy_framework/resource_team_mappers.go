package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/userroles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func mapTeamStateToResource(ctx context.Context, model schemas.TeamModel) *teams.Team {
	name := model.Name.ValueString()
	team := teams.NewTeam(name)
	team.ID = model.ID.ValueString()

	team.Description = model.Description.ValueString()
	team.SpaceID = model.SpaceID.ValueString()
	team.MemberUserIDs = expandStringSet(model.Users)

	if !model.ExternalSecurityGroup.IsNull() && !model.ExternalSecurityGroup.IsUnknown() {
		team.ExternalSecurityGroups = mapExternalSecurityGroupsStateToResource(ctx, model.ExternalSecurityGroup)
	} else {
		team.ExternalSecurityGroups = []core.NamedReferenceItem{}
	}

	if !model.CanBeDeleted.IsNull() && !model.CanBeDeleted.IsUnknown() {
		team.CanBeDeleted = model.CanBeDeleted.ValueBool()
	}

	if !model.CanBeRenamed.IsNull() && !model.CanBeRenamed.IsUnknown() {
		team.CanBeRenamed = model.CanBeRenamed.ValueBool()
	}

	if !model.CanChangeMembers.IsNull() && !model.CanChangeMembers.IsUnknown() {
		team.CanChangeMembers = model.CanChangeMembers.ValueBool()
	}

	if !model.CanChangeRoles.IsNull() && !model.CanChangeRoles.IsUnknown() {
		team.CanChangeRoles = model.CanChangeRoles.ValueBool()
	}

	return team
}

func mapTeamResourceToState(ctx context.Context, team *teams.Team, model schemas.TeamModel, client *client.Client) schemas.TeamModel {
	model.ID = types.StringValue(team.GetID())
	model.Name = types.StringValue(team.Name)
	model.Description = types.StringValue(team.Description)
	model.SpaceID = types.StringValue(team.SpaceID)

	if len(team.MemberUserIDs) > 0 {
		model.Users = flattenStringSet(team.MemberUserIDs, model.Users)
	} else {
		model.Users = types.SetNull(types.StringType)
	}

	if len(team.ExternalSecurityGroups) > 0 {
		model.ExternalSecurityGroup = mapExternalSecurityGroupsResourceToState(ctx, team.ExternalSecurityGroups)
	} else {
		objectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":                  types.StringType,
				"display_name":        types.StringType,
				"display_id_and_name": types.BoolType,
			},
		}
		model.ExternalSecurityGroup = types.ListNull(objectType)
	}

	model.CanBeDeleted = types.BoolValue(team.CanBeDeleted)
	model.CanBeRenamed = types.BoolValue(team.CanBeRenamed)
	model.CanChangeMembers = types.BoolValue(team.CanChangeMembers)
	model.CanChangeRoles = types.BoolValue(team.CanChangeRoles)

	userRoles, err := client.Teams.GetScopedUserRoles(*team, core.SkipTakeQuery{})
	if err == nil {
		model.UserRole = mapUserRoleSetResourceToState(ctx, userRoles.Items, model.UserRole)
	}

	return model
}

func mapExternalSecurityGroupsStateToResource(ctx context.Context, securityGroups types.List) []core.NamedReferenceItem {
	if securityGroups.IsNull() || securityGroups.IsUnknown() {
		return nil
	}

	securityGroupElements := securityGroups.Elements()
	externalGroups := make([]core.NamedReferenceItem, len(securityGroupElements))

	for i, securityGroupElement := range securityGroupElements {
		securityGroupAttributes := securityGroupElement.(types.Object).Attributes()

		externalGroups[i] = core.NamedReferenceItem{
			ID:               securityGroupAttributes["id"].(types.String).ValueString(),
			DisplayName:      securityGroupAttributes["display_name"].(types.String).ValueString(),
			DisplayIDAndName: securityGroupAttributes["display_id_and_name"].(types.Bool).ValueBool(),
		}
	}

	return externalGroups
}

func mapExternalSecurityGroupsResourceToState(ctx context.Context, groups []core.NamedReferenceItem) types.List {
	if len(groups) == 0 {
		objectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":                  types.StringType,
				"display_name":        types.StringType,
				"display_id_and_name": types.BoolType,
			},
		}
		return types.ListNull(objectType)
	}

	groupList := make([]attr.Value, len(groups))
	for i, group := range groups {
		objectValue := types.ObjectValueMust(
			map[string]attr.Type{
				"id":                  types.StringType,
				"display_name":        types.StringType,
				"display_id_and_name": types.BoolType,
			},
			map[string]attr.Value{
				"id":                  types.StringValue(group.ID),
				"display_name":        types.StringValue(group.DisplayName),
				"display_id_and_name": types.BoolValue(group.DisplayIDAndName),
			},
		)
		groupList[i] = objectValue
	}

	objectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                  types.StringType,
			"display_name":        types.StringType,
			"display_id_and_name": types.BoolType,
		},
	}
	return types.ListValueMust(objectType, groupList)
}

func mapUserRoleSetStateToResource(ctx context.Context, team *teams.Team, userRoles types.Set) []*userroles.ScopedUserRole {
	if userRoles.IsNull() || userRoles.IsUnknown() {
		return nil
	}

	userRoleSetElements := userRoles.Elements()
	scopedUserRoles := make([]*userroles.ScopedUserRole, len(userRoleSetElements))

	for i, userRoleSetElement := range userRoleSetElements {
		userRoleSetElementAttributes := userRoleSetElement.(types.Object).Attributes()

		userRoleId := userRoleSetElementAttributes["user_role_id"].(types.String).ValueString()
		spaceId := userRoleSetElementAttributes["space_id"].(types.String).ValueString()

		scopedUserRole := userroles.NewScopedUserRole(userRoleId)
		scopedUserRole.TeamID = team.ID
		scopedUserRole.SpaceID = spaceId

		if id, ok := userRoleSetElementAttributes["id"]; ok && !id.IsNull() && !id.IsUnknown() {
			scopedUserRole.ID = id.(types.String).ValueString()
		}

		if envIds, ok := userRoleSetElementAttributes["environment_ids"]; ok && !envIds.IsNull() && !envIds.IsUnknown() {
			scopedUserRole.EnvironmentIDs = expandStringSet(envIds.(types.Set))
		}

		if projGroupIds, ok := userRoleSetElementAttributes["project_group_ids"]; ok && !projGroupIds.IsNull() && !projGroupIds.IsUnknown() {
			scopedUserRole.ProjectGroupIDs = expandStringSet(projGroupIds.(types.Set))
		}

		if projIds, ok := userRoleSetElementAttributes["project_ids"]; ok && !projIds.IsNull() && !projIds.IsUnknown() {
			scopedUserRole.ProjectIDs = expandStringSet(projIds.(types.Set))
		}

		if tenantIds, ok := userRoleSetElementAttributes["tenant_ids"]; ok && !tenantIds.IsNull() && !tenantIds.IsUnknown() {
			scopedUserRole.TenantIDs = expandStringSet(tenantIds.(types.Set))
		}

		scopedUserRoles[i] = scopedUserRole
	}

	return scopedUserRoles
}

func mapUserRoleSetResourceToState(ctx context.Context, userRoles []*userroles.ScopedUserRole, model types.Set) types.Set {
	if len(userRoles) == 0 {
		objectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":                types.StringType,
				"user_role_id":      types.StringType,
				"space_id":          types.StringType,
				"team_id":           types.StringType,
				"environment_ids":   types.SetType{ElemType: types.StringType},
				"project_group_ids": types.SetType{ElemType: types.StringType},
				"project_ids":       types.SetType{ElemType: types.StringType},
				"tenant_ids":        types.SetType{ElemType: types.StringType},
			},
		}
		return types.SetNull(objectType)
	}

	roleList := make([]attr.Value, len(userRoles))
	for i, role := range userRoles {
		objectValue := types.ObjectValueMust(
			map[string]attr.Type{
				"id":                types.StringType,
				"user_role_id":      types.StringType,
				"space_id":          types.StringType,
				"team_id":           types.StringType,
				"environment_ids":   types.SetType{ElemType: types.StringType},
				"project_group_ids": types.SetType{ElemType: types.StringType},
				"project_ids":       types.SetType{ElemType: types.StringType},
				"tenant_ids":        types.SetType{ElemType: types.StringType},
			},
			map[string]attr.Value{
				"id":                types.StringValue(role.ID),
				"user_role_id":      types.StringValue(role.UserRoleID),
				"space_id":          types.StringValue(role.SpaceID),
				"team_id":           types.StringValue(role.TeamID),
				"environment_ids":   flattenStringSet(role.EnvironmentIDs, types.SetNull(types.StringType)),
				"project_group_ids": flattenStringSet(role.ProjectGroupIDs, types.SetNull(types.StringType)),
				"project_ids":       flattenStringSet(role.ProjectIDs, types.SetNull(types.StringType)),
				"tenant_ids":        flattenStringSet(role.TenantIDs, types.SetNull(types.StringType)),
			},
		)
		roleList[i] = objectValue
	}

	objectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                types.StringType,
			"user_role_id":      types.StringType,
			"space_id":          types.StringType,
			"team_id":           types.StringType,
			"environment_ids":   types.SetType{ElemType: types.StringType},
			"project_group_ids": types.SetType{ElemType: types.StringType},
			"project_ids":       types.SetType{ElemType: types.StringType},
			"tenant_ids":        types.SetType{ElemType: types.StringType},
		},
	}
	return types.SetValueMust(objectType, roleList)
}
