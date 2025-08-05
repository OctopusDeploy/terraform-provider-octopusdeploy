package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTeamResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		CheckDestroy: testTeamDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "1.2.1",
						Source:            "OctopusDeploy/octopusdeploy",
					},
				},
				Config: teamConfig(name, description),
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   teamConfig(name, description),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updateTeamConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					testTeam(t, name, description),
				),
			},
		},
	})
}

func TestTeamResource_UpgradeFromSDK_ToPluginFramework_WithUserRole(t *testing.T) {
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testTeamDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "1.2.1",
						Source:            "OctopusDeploy/octopusdeploy",
					},
				},
				Config: teamConfigWithUserRole(name, description, userRoleName),
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   teamConfigWithUserRole(name, description, userRoleName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updateTeamConfigWithUserRole(name, description, userRoleName),
				Check: resource.ComposeTestCheckFunc(
					testTeamWithUserRole(t, name, description),
				),
			},
		},
	})
}

func teamConfig(name, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_team" "team1" {
		name = "%s"
		description = "%s"
	}`, name, description)
}

func updateTeamConfig(name, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_team" "team1" {
		name = "%s"
		description = "%s - updated"
	}`, name, description)
}

func teamConfigWithUserRole(name, description, userRoleName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user_role" "user_role1" {
		granted_space_permissions = ["AccountCreate"]
		name = "%s"
	}

	resource "octopusdeploy_team" "team2" {
		name = "%s"
		description = "%s"

		user_role {
			space_id = "Spaces-1"
			user_role_id = octopusdeploy_user_role.user_role1.id
		}
	}`, userRoleName, name, description)
}

func updateTeamConfigWithUserRole(name, description, userRoleName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user_role" "user_role1" {
		granted_space_permissions = ["AccountCreate"]
		name = "%s"
	}

	resource "octopusdeploy_team" "team2" {
		name = "%s"
		description = "%s - updated"

		user_role {
			space_id = "Spaces-1"
			user_role_id = octopusdeploy_user_role.user_role1.id
		}
	}`, userRoleName, name, description)
}

func testTeamDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_team" {
			continue
		}

		team, err := octoClient.Teams.GetByID(rs.Primary.ID)
		if err == nil && team != nil {
			return fmt.Errorf("team (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testTeam(t *testing.T, name, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamId := s.RootModule().Resources["octopusdeploy_team.team1"].Primary.ID
		team, err := octoClient.Teams.GetByID(teamId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve team by ID: %s", err)
		}

		assert.NotEmpty(t, team.ID, "Team ID should not be empty")
		assert.Equal(t, name, team.Name, "Team name did not match expected value")
		assert.Equal(t, description+" - updated", team.Description, "Team description did not match expected value")

		return nil
	}
}

func testTeamWithUserRole(t *testing.T, name, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamId := s.RootModule().Resources["octopusdeploy_team.team1"].Primary.ID
		team, err := octoClient.Teams.GetByID(teamId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve team by ID: %s", err)
		}

		assert.NotEmpty(t, team.ID, "Team ID should not be empty")
		assert.Equal(t, name, team.Name, "Team name did not match expected value")
		assert.Equal(t, description+" - updated", team.Description, "Team description did not match expected value")

		userRoles, err := octoClient.Teams.GetScopedUserRoles(*team, core.SkipTakeQuery{})
		if err != nil {
			return fmt.Errorf("Failed to retrieve user roles: %s", err)
		}

		assert.NotEmpty(t, userRoles.Items, "Team should have user roles")
		assert.Len(t, userRoles.Items, 1, "Team should have exactly one user role")
		assert.Equal(t, "Spaces-1", userRoles.Items[0].SpaceID, "User role space ID should match")

		return nil
	}
}
