package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"os"
	"testing"
)

func TestScopedUserRoleResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	space := NewTestSpace(t)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccScopedUserRoleCheckDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "1.2.1",
						Source:            "OctopusDeploy/octopusdeploy",
					},
				},
				Config: scopedUserRoleConfig(space.ID, userRoleName),
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   scopedUserRoleConfig(space.ID, userRoleName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   scopedUserRoleConfig(space.ID, userRoleName),
			},
		},
	})
}

func TestScopedUserRoleResource_UpgradeFromSDK_ToPluginFramework_WithScopes(t *testing.T) {
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	space := NewTestSpace(t)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccScopedUserRoleCheckDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "1.2.1",
						Source:            "OctopusDeploy/octopusdeploy",
					},
				},
				Config: scopedUserRoleConfigWithScopes(space.ID, userRoleName, environmentName),
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   scopedUserRoleConfigWithScopes(space.ID, userRoleName, environmentName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   scopedUserRoleConfigWithScopes(space.ID, userRoleName, environmentName),
			},
		},
	})
}

func scopedUserRoleConfig(spaceID, userRoleName string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_user_role" "user_role1" {
		name = "%s"
		description = "Test user role with environment permissions"
		granted_space_permissions = ["EnvironmentView"]
	}

	resource "octopusdeploy_scoped_user_role" "scoped_user_role1" {
		space_id     = "%s"
		team_id      = "teams-everyone"
		user_role_id = octopusdeploy_user_role.user_role1.id
	}`, userRoleName, spaceID)
}

func scopedUserRoleConfigWithScopes(spaceID, userRoleName, environmentName string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_user_role" "user_role1" {
		name = "%s"
		description = "Test user role with environment permissions"
		granted_space_permissions = ["EnvironmentView"]
	}

	resource "octopusdeploy_environment" "environment1" {
		space_id = "%s"
		name = "%s"
	}

	resource "octopusdeploy_scoped_user_role" "scoped_user_role1" {
		space_id        = "%s"
		team_id         = "teams-everyone"
		user_role_id    = octopusdeploy_user_role.user_role1.id
		environment_ids = [octopusdeploy_environment.environment1.id]
	}`, userRoleName, spaceID, environmentName, spaceID)
}
