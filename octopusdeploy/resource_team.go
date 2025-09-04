package octopusdeploy

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/userroles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		DeleteContext: resourceTeamDelete,
		Description:   "This resource manages teams in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceTeamRead,
		Schema:        getTeamSchema(),
		UpdateContext: resourceTeamUpdate,
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	team := expandTeam(d)

	log.Printf("[INFO] creating team: %#v", team)

	client := m.(*client.Client)
	createdTeam, err := client.Teams.Add(team)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := resourceTeamUpdateUserRoles(ctx, d, m, createdTeam); err != nil {
		return diag.FromErr(err)
	}

	if err := setTeam(ctx, d, createdTeam); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdTeam.GetID())

	log.Printf("[INFO] team created (%s)", d.Id())
	return resourceTeamRead(ctx, d, m)
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting team (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.Teams.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] team deleted")
	return nil
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading team (%s)", d.Id())

	client := m.(*client.Client)
	team, err := client.Teams.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "team")
	}

	userRoles, err := client.Teams.GetScopedUserRoles(*team, core.SkipTakeQuery{})
	if err != nil {
		return diag.FromErr(err)
	}
	remoteUserRoles := flattenScopedUserRoles(userRoles.Items)
	d.Set("user_role", remoteUserRoles)

	if err := setTeam(ctx, d, team); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] team read (%s)", d.Id())
	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating team (%s)", d.Id())

	team := expandTeam(d)
	client := m.(*client.Client)
	updatedTeam, err := client.Teams.Update(team)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := resourceTeamUpdateUserRoles(ctx, d, m, team); err != nil {
		return diag.FromErr(err)
	}

	if err := setTeam(ctx, d, updatedTeam); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] team updated (%s)", d.Id())
	return resourceTeamRead(ctx, d, m)
}

func expandUserRoles(team *teams.Team, userRoles []interface{}) []*userroles.ScopedUserRole {
	values := make([]*userroles.ScopedUserRole, 0, len(userRoles))
	for _, rawUserRole := range userRoles {
		userRole := rawUserRole.(map[string]interface{})
		scopedUserRole := userroles.NewScopedUserRole(userRole["user_role_id"].(string))
		scopedUserRole.TeamID = team.ID

		if spaceID, ok := userRole["space_id"]; ok && spaceID.(string) != "" {
			scopedUserRole.SpaceID = spaceID.(string)
		}

		if v, ok := userRole["id"]; ok {
			scopedUserRole.ID = v.(string)
		} else {
			scopedUserRole.ID = ""
		}

		if v, ok := userRole["environment_ids"]; ok {
			scopedUserRole.EnvironmentIDs = getSliceFromTerraformTypeList(v)
		}

		if v, ok := userRole["project_group_ids"]; ok {
			scopedUserRole.ProjectGroupIDs = getSliceFromTerraformTypeList(v)
		}

		if v, ok := userRole["project_ids"]; ok {
			scopedUserRole.ProjectIDs = getSliceFromTerraformTypeList(v)
		}

		if v, ok := userRole["tenant_ids"]; ok {
			scopedUserRole.TenantIDs = getSliceFromTerraformTypeList(v)
		}
		values = append(values, scopedUserRole)
	}
	return values
}
func resourceTeamUpdateUserRoles(ctx context.Context, d *schema.ResourceData, m interface{}, team *teams.Team) error {
	log.Printf("[INFO] updating team user roles (%s)", d.Id())
	if d.HasChange("user_role") {
		log.Printf("[INFO] user role has changes (%s)", d.Id())
		o, n := d.GetChange("user_role")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := expandUserRoles(team, os.Difference(ns).List())
		add := expandUserRoles(team, ns.Difference(os).List())

		if len(remove) > 0 || len(add) > 0 {
			log.Printf("[INFO] user role found diff (%s)", d.Id())
			client := m.(*client.Client)
			if len(remove) > 0 {
				log.Printf("[INFO] removing user roles from team (%s)", d.Id())
				for _, userRole := range remove {
					if userRole.ID != "" {
						err := client.ScopedUserRoles.DeleteByID(userRole.ID)
						if err != nil {
							apiError := err.(*core.APIError)
							if apiError.StatusCode != 404 {
								// It's already been deleted, maybe mixing with the independent resource?
								return fmt.Errorf("error removing user role %s from team %s: %s", userRole.ID, team.ID, err)
							}
						}
					}
				}
			}
			if len(add) > 0 {
				log.Printf("[INFO] adding new user roles to team (%s)", d.Id())
				for _, userRole := range add {
					_, err := client.ScopedUserRoles.Add(userRole)
					if err != nil {
						return fmt.Errorf("error creating user role for team %s: %s", team.ID, err)
					}
				}
			}
		}
	}
	return nil
}

func resourceTeamUserRoleListSetHash(buf *bytes.Buffer, v interface{}) {
	vs := v.(*schema.Set).List()
	s := make([]string, len(vs))
	for i, raw := range vs {
		s[i] = raw.(string)
	}
	sort.Strings(s)
	for _, v := range s {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}
}

func resourceTeamUserRoleSetHash(v interface{}) int {
	var buf bytes.Buffer

	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["user_role_id"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["space_id"].(string)))

	if v, ok := m["environment_ids"]; ok {
		resourceTeamUserRoleListSetHash(&buf, v)
	}
	if v, ok := m["project_group_ids"]; ok {
		resourceTeamUserRoleListSetHash(&buf, v)
	}
	if v, ok := m["project_ids"]; ok {
		resourceTeamUserRoleListSetHash(&buf, v)
	}
	if v, ok := m["tenant_ids"]; ok {
		resourceTeamUserRoleListSetHash(&buf, v)
	}

	return stringHashCode(buf.String())
}

func flattenScopedUserRoles(scopedUserRoles []*userroles.ScopedUserRole) []interface{} {
	if scopedUserRoles == nil {
		return nil
	}
	var flattenedScopedUserRoles = make([]interface{}, len(scopedUserRoles))
	for key, scopedUserRole := range scopedUserRoles {
		flattenedScopedUserRoles[key] = flattenScopedUserRole(scopedUserRole)
	}

	return flattenedScopedUserRoles
}

func flattenScopedUserRole(scopedUserRole *userroles.ScopedUserRole) map[string]interface{} {
	if scopedUserRole == nil {
		return nil
	}

	result := map[string]interface{}{
		"environment_ids":   schema.NewSet(schema.HashString, flattenArray(scopedUserRole.EnvironmentIDs)),
		"id":                scopedUserRole.ID,
		"project_group_ids": schema.NewSet(schema.HashString, flattenArray(scopedUserRole.ProjectGroupIDs)),
		"project_ids":       schema.NewSet(schema.HashString, flattenArray(scopedUserRole.ProjectIDs)),
		"team_id":           scopedUserRole.TeamID,
		"tenant_ids":        schema.NewSet(schema.HashString, flattenArray(scopedUserRole.TenantIDs)),
		"user_role_id":      scopedUserRole.UserRoleID,
	}

	// Only include space_id if it's not empty
	if scopedUserRole.SpaceID != "" {
		result["space_id"] = scopedUserRole.SpaceID
	}

	return result
}
