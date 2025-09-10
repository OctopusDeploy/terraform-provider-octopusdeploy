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
				Config: scopedUserRoleConfig(userRoleName),
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   scopedUserRoleConfig(userRoleName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   scopedUserRoleConfig(userRoleName),
			},
		},
	})
}

func TestScopedUserRoleResource_UpgradeFromSDK_ToPluginFramework_WithScopes(t *testing.T) {
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

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
				Config: scopedUserRoleConfigWithScopes(userRoleName, environmentName),
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   scopedUserRoleConfigWithScopes(userRoleName, environmentName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   scopedUserRoleConfigWithScopes(userRoleName, environmentName),
			},
		},
	})
}

func scopedUserRoleConfig(userRoleName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user_role" "user_role1" {
		name = "%s"
		description = "Test user role with environment permissions"
		granted_space_permissions = ["EnvironmentView"]
	}

	resource "octopusdeploy_scoped_user_role" "scoped_user_role1" {
		space_id     = "Spaces-1"
		team_id      = "teams-everyone"
		user_role_id = octopusdeploy_user_role.user_role1.id
	}`, userRoleName)
}

func scopedUserRoleConfigWithScopes(userRoleName, environmentName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user_role" "user_role1" {
		name = "%s"
		description = "Test user role with environment permissions"
		granted_space_permissions = ["EnvironmentView"]
	}

	resource "octopusdeploy_environment" "environment1" {
		name = "%s"
	}

	resource "octopusdeploy_scoped_user_role" "scoped_user_role1" {
		space_id        = "Spaces-1"
		team_id         = "teams-everyone"
		user_role_id    = octopusdeploy_user_role.user_role1.id
		environment_ids = [octopusdeploy_environment.environment1.id]
	}`, userRoleName, environmentName)
}
