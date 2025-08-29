package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/userroles"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployUserRoleBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_user_role." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccUserRoleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccUserRoleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "granted_space_permissions.#", "0"),
					resource.TestCheckResourceAttr(prefix, "granted_system_permissions.#", "0"),
				),
				Config: testAccUserRoleBasic(localName, name),
			},
		},
	})
}

func TestAccOctopusDeployUserRoleUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_user_role." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccUserRoleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccUserRoleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
				),
				Config: testAccUserRoleWithDescription(localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccUserRoleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
				),
				Config: testAccUserRoleWithDescription(localName, name, newDescription),
			},
		},
	})
}

func TestAccOctopusDeployUserRoleWithPermissions(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_user_role." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccUserRoleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccUserRoleExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "granted_space_permissions.#", "2"),
					resource.TestCheckResourceAttr(prefix, "granted_system_permissions.#", "1"),
				),
				Config: testAccUserRoleWithPermissions(localName, name, description),
			},
		},
	})
}

func TestAccOctopusDeployUserRoleImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_user_role." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccUserRoleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserRoleBasic(localName, name),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccUserRoleImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccUserRoleBasic(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user_role" "%s" {
		name = "%s"
	}`, localName, name)
}

func testAccUserRoleWithDescription(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user_role" "%s" {
		name        = "%s"
		description = "%s"
	}`, localName, name, description)
}

func testAccUserRoleWithPermissions(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user_role" "%s" {
		name                        = "%s"
		description                 = "%s"
		granted_space_permissions   = ["AccountCreate", "AccountView"]
		granted_system_permissions  = ["SpaceView"]
	}`, localName, name, description)
}

func testAccUserRoleExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		userRoleID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := userroles.GetByID(octoClient, userRoleID); err != nil {
			return err
		}

		return nil
	}
}

func testAccUserRoleCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_user_role" {
			continue
		}

		if userRole, err := userroles.GetByID(octoClient, rs.Primary.ID); err == nil {
			return fmt.Errorf("user role (%s) still exists", userRole.GetID())
		}
	}

	return nil
}

func testAccUserRoleImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}