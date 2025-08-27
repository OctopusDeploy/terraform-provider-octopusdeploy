package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeploySpaceBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_space." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	slug := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:                 testSpaceCheckDestroy,
		ProtoV6ProviderFactories:     ProtoV6ProviderFactories(),
		PreCheck:                     func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testSpaceExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "slug", slug),
					resource.TestCheckResourceAttr(prefix, "space_managers_teams.#", "1"),
					resource.TestCheckResourceAttrSet(prefix, "space_managers_teams.0"),
				),
				Config: testSpaceBasic(localName, name, slug),
			},
		},
	})
}

func TestAccOctopusDeploySpaceUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_space." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	slug := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:                 testSpaceCheckDestroy,
		ProtoV6ProviderFactories:     ProtoV6ProviderFactories(),
		PreCheck:                     func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testSpaceExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "slug", slug),
					resource.TestCheckResourceAttr(prefix, "description", description),
				),
				Config: testSpaceWithDescription(localName, name, slug, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testSpaceExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "slug", slug),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
				),
				Config: testSpaceWithDescription(localName, newName, slug, newDescription),
			},
		},
	})
}

func TestAccOctopusDeploySpaceWithTaskQueueStopped(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_space." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	slug := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:                 testSpaceCheckDestroy,
		ProtoV6ProviderFactories:     ProtoV6ProviderFactories(),
		PreCheck:                     func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testSpaceExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "slug", slug),
					resource.TestCheckResourceAttr(prefix, "is_task_queue_stopped", "true"),
				),
				Config: testSpaceWithTaskQueueStopped(localName, name, slug),
			},
		},
	})
}

func TestAccOctopusDeploySpaceImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	slug := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:                 testSpaceCheckDestroy,
		ProtoV6ProviderFactories:     ProtoV6ProviderFactories(),
		PreCheck:                     func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testSpaceBasic(localName, name, slug),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"space_managers_teams",
				},
			},
		},
	})
}

func testSpaceBasic(localName string, name string, slug string) string {
	userLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userDisplayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userEmailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	userPassword := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userUsername := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	return fmt.Sprintf(testUserBasic(userLocalName, userDisplayName, true, false, userPassword, userUsername, userEmailAddress)+"\n"+
		`resource "octopusdeploy_space" "%s" {
			name = "%s"
			slug = "%s"
			space_managers_teams = ["teams-managers"]
			lifecycle {
				ignore_changes = [space_managers_teams]
			}
		}`, localName, name, slug)
}

func testSpaceWithDescription(localName string, name string, slug string, description string) string {
	userLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userDisplayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userEmailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	userPassword := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userUsername := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	return fmt.Sprintf(testUserBasic(userLocalName, userDisplayName, true, false, userPassword, userUsername, userEmailAddress)+"\n"+
		`resource "octopusdeploy_space" "%s" {
			name = "%s"
			slug = "%s"
			description = "%s"
			space_managers_teams = ["teams-managers"]
			lifecycle {
				ignore_changes = [space_managers_teams]
			}
		}`, localName, name, slug, description)
}

func testSpaceWithTaskQueueStopped(localName string, name string, slug string) string {
	userLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userDisplayName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userEmailAddress := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "@example.com"
	userPassword := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userUsername := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	return fmt.Sprintf(testUserBasic(userLocalName, userDisplayName, true, false, userPassword, userUsername, userEmailAddress)+"\n"+
		`resource "octopusdeploy_space" "%s" {
			name = "%s"
			slug = "%s"
			is_task_queue_stopped = true
			space_managers_teams = ["teams-managers"]
			lifecycle {
				ignore_changes = [space_managers_teams]
			}
		}`, localName, name, slug)
}

func testUserBasic(localName string, displayName string, isActive bool, isService bool, password string, username string, emailAddress string) string {
	return fmt.Sprintf(`resource "octopusdeploy_user" "%s" {
		display_name  = "%s"
		email_address = "%s"
		is_active     = %v
		is_service    = %v
		password      = "%s"
		username      = "%s"

		identity {
			provider = "Octopus ID"
			claim {
				name = "email"
				is_identifying_claim = true
				value = "%s"
			}
			claim {
				name = "dn"
				is_identifying_claim = false
				value = "%s"
			}
		}
	}`, localName, displayName, emailAddress, isActive, isService, password, username, emailAddress, displayName)
}

func testSpaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if _, err := spaces.GetByID(octoClient, rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testSpaceCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_space" {
			continue
		}

		spaceID := rs.Primary.ID
		space, err := spaces.GetByID(octoClient, spaceID)
		if err == nil {
			if space != nil {
				return fmt.Errorf("space (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}