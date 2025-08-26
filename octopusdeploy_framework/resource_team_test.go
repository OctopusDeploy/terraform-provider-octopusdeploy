package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployTeamBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_team." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTeamCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccTeamExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "space_id"),
				),
				Config: testAccTeamBasic(localName, name, description),
			},
		},
	})
}

func TestAccOctopusDeployTeamUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_team." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTeamCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccTeamExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
				),
				Config: testAccTeamBasic(localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccTeamExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
				),
				Config: testAccTeamBasic(localName, name, newDescription),
			},
		},
	})
}

func TestAccOctopusDeployTeamWithUserRoleBlocks(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_team." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTeamCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccTeamExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "user_role.#", "1"),
				),
				Config: testAccTeamWithUserRole(localName, name, description, userRoleName),
			},
		},
	})
}

func TestAccOctopusDeployTeamWithExternalSecurityGroupBlocks(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_team." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	externalGroupId := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	externalGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTeamCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccTeamExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "external_security_group.#", "1"),
				),
				Config: testAccTeamWithExternalSecurityGroup(localName, name, description, externalGroupId, externalGroupName),
			},
		},
	})
}

func TestAccOctopusDeployTeamImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_team." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTeamCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamBasic(localName, name, description),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccTeamImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccTeamBasic(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_team" "%s" {
		name        = "%s"
		description = "%s"
	}`, localName, name, description)
}

func testAccTeamWithUserRole(localName string, name string, description string, userRoleName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_user_role" "%s" {
		granted_space_permissions = ["AccountCreate"]
		name = "%s"
	}

	resource "octopusdeploy_team" "%s" {
		name        = "%s"
		description = "%s"
		
		user_role {
			space_id     = "Spaces-1"
			user_role_id = octopusdeploy_user_role.%s.id
		}
	}`, userRoleName, userRoleName, localName, name, description, userRoleName)
}

func testAccTeamWithExternalSecurityGroup(localName string, name string, description string, externalGroupId string, externalGroupName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_team" "%s" {
		name        = "%s"
		description = "%s"
		
		external_security_group {
			id           = "%s"
			display_name = "%s"
		}
	}`, localName, name, description, externalGroupId, externalGroupName)
}

func testAccTeamExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Teams.GetByID(teamID); err != nil {
			return err
		}

		return nil
	}
}

func testAccTeamCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_team" {
			continue
		}

		if team, err := octoClient.Teams.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("team (%s) still exists", team.GetID())
		}
	}

	return nil
}

func testAccTeamImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

// Helper function to validate team properties for CRUD test
func testAccTeamValidateProperties(expectedName string, expectedDescription string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamResource := s.RootModule().Resources["octopusdeploy_team.team1"]
		teamID := teamResource.Primary.ID

		team, err := octoClient.Teams.GetByID(teamID)
		if err != nil {
			return fmt.Errorf("Received an error retrieving team %s", err)
		}

		if team.Name != expectedName {
			return fmt.Errorf("Expected team name to be %s but was %s", expectedName, team.Name)
		}

		if team.Description != expectedDescription {
			return fmt.Errorf("Expected team description to be %s but was %s", expectedDescription, team.Description)
		}

		return nil
	}
}