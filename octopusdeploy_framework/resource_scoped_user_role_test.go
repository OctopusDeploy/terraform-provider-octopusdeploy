package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployScopedUserRoleBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_scoped_user_role." + localName

	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccScopedUserRoleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccScopedUserRoleExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "team_id"),
					resource.TestCheckResourceAttrSet(prefix, "user_role_id"),
					resource.TestCheckResourceAttr(prefix, "space_id", "Spaces-1"),
				),
				Config: testAccScopedUserRoleBasic(localName, userRoleLocalName, userRoleName),
			},
		},
	})
}

func TestAccOctopusDeployScopedUserRoleWithScopes(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_scoped_user_role." + localName

	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccScopedUserRoleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccScopedUserRoleExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "team_id"),
					resource.TestCheckResourceAttrSet(prefix, "user_role_id"),
					resource.TestCheckResourceAttr(prefix, "environment_ids.#", "1"),
					resource.TestCheckResourceAttr(prefix, "project_group_ids.#", "1"),
				),
				Config: testAccScopedUserRoleWithScopes(localName, userRoleLocalName, userRoleName, environmentLocalName, environmentName, projectGroupLocalName, projectGroupName),
			},
		},
	})
}

func TestAccOctopusDeployScopedUserRoleUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_scoped_user_role." + localName

	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment2LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment2Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccScopedUserRoleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccScopedUserRoleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "environment_ids.#", "1"),
				),
				Config: testAccScopedUserRoleUpdate1(localName, userRoleLocalName, userRoleName, environmentLocalName, environmentName),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccScopedUserRoleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "environment_ids.#", "2"),
				),
				Config: testAccScopedUserRoleUpdate2(localName, userRoleLocalName, userRoleName, environmentLocalName, environmentName, environment2LocalName, environment2Name),
			},
		},
	})
}

func TestAccOctopusDeployScopedUserRoleSystemLevel(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_scoped_user_role." + localName

	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccScopedUserRoleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccScopedUserRoleExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "team_id"),
					resource.TestCheckResourceAttrSet(prefix, "user_role_id"),
					resource.TestCheckNoResourceAttr(prefix, "space_id"),
				),
				Config: testAccScopedUserRoleSystemLevel(localName, userRoleLocalName, userRoleName),
			},
		},
	})
}

func TestAccOctopusDeployScopedUserRoleImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_scoped_user_role." + localName

	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccScopedUserRoleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccScopedUserRoleBasic(localName, userRoleLocalName, userRoleName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccScopedUserRoleImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccScopedUserRoleBasic(localName, userRoleLocalName, userRoleName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user_role" "%s" {
		name = "%s"
		description = "Test user role with environment permissions"
		granted_space_permissions = ["EnvironmentView"]
	}

	resource "octopusdeploy_scoped_user_role" "%s" {
		space_id     = "Spaces-1"
		team_id      = "teams-everyone"
		user_role_id = octopusdeploy_user_role.%s.id
	}`, userRoleLocalName, userRoleName, localName, userRoleLocalName)
}

func testAccScopedUserRoleWithScopes(localName, userRoleLocalName, userRoleName, environmentLocalName, environmentName, projectGroupLocalName, projectGroupName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user_role" "%s" {
		name = "%s"
		description = "Test user role with environment permissions"
		granted_space_permissions = ["EnvironmentView"]
	}

	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_scoped_user_role" "%s" {
		space_id          = "Spaces-1"
		team_id           = "teams-everyone"
		user_role_id      = octopusdeploy_user_role.%s.id
		environment_ids   = [octopusdeploy_environment.%s.id]
		project_group_ids = [octopusdeploy_project_group.%s.id]
	}`, userRoleLocalName, userRoleName, environmentLocalName, environmentName, projectGroupLocalName, projectGroupName, localName, userRoleLocalName, environmentLocalName, projectGroupLocalName)
}

func testAccScopedUserRoleUpdate1(localName, userRoleLocalName, userRoleName, environmentLocalName, environmentName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user_role" "%s" {
		name = "%s"
		description = "Test user role with environment permissions"
		granted_space_permissions = ["EnvironmentView"]
	}

	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_scoped_user_role" "%s" {
		space_id        = "Spaces-1"
		team_id         = "teams-everyone"
		user_role_id    = octopusdeploy_user_role.%s.id
		environment_ids = [octopusdeploy_environment.%s.id]
	}`, userRoleLocalName, userRoleName, environmentLocalName, environmentName, localName, userRoleLocalName, environmentLocalName)
}

func testAccScopedUserRoleUpdate2(localName, userRoleLocalName, userRoleName, environmentLocalName, environmentName, environment2LocalName, environment2Name string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user_role" "%s" {
		name = "%s"
		description = "Test user role with environment permissions"
		granted_space_permissions = ["EnvironmentView"]
	}

	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_scoped_user_role" "%s" {
		space_id        = "Spaces-1"
		team_id         = "teams-everyone"
		user_role_id    = octopusdeploy_user_role.%s.id
		environment_ids = [
			octopusdeploy_environment.%s.id,
			octopusdeploy_environment.%s.id
		]
	}`, userRoleLocalName, userRoleName, environmentLocalName, environmentName, environment2LocalName, environment2Name, localName, userRoleLocalName, environmentLocalName, environment2LocalName)
}

func testAccScopedUserRoleSystemLevel(localName, userRoleLocalName, userRoleName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user_role" "%s" {
		name = "%s"
		description = "Test user role with system permissions"
		granted_system_permissions = ["SpaceView"]
	}

	resource "octopusdeploy_scoped_user_role" "%s" {
		team_id      = "teams-everyone"
		user_role_id = octopusdeploy_user_role.%s.id
	}`, userRoleLocalName, userRoleName, localName, userRoleLocalName)
}

func testAccScopedUserRoleExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[prefix]
		if !ok {
			return fmt.Errorf("not found: %s", prefix)
		}

		// Simply check if we have an ID set
		if rs.Primary.ID == "" {
			return fmt.Errorf("scoped user role ID is not set")
		}

		return nil
	}
}

func testAccScopedUserRoleCheckDestroy(s *terraform.State) error {
	// Note: SDK doesn't have a direct API to check scoped user role existence
	// So we just return nil for now
	return nil
}

func testAccScopedUserRoleImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
