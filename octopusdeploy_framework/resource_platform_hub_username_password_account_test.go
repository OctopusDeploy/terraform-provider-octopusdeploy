package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubaccounts"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccPlatformHubUsernamePasswordAccountCreate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_username_password_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubUsernamePasswordAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubUsernamePasswordAccountBasic(localName, name, description, username),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubUsernamePasswordAccountExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "username", username),
				),
			},
		},
	})
}

func TestAccPlatformHubUsernamePasswordAccountUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_username_password_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	updatedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedUsername := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubUsernamePasswordAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubUsernamePasswordAccountBasic(localName, name, description, username),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubUsernamePasswordAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "username", username),
				),
			},
			{
				Config: testPlatformHubUsernamePasswordAccountBasic(localName, updatedName, updatedDescription, updatedUsername),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubUsernamePasswordAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "username", updatedUsername),
				),
			},
		},
	})
}

func TestAccPlatformHubUsernamePasswordAccountImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_username_password_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	username := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubUsernamePasswordAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubUsernamePasswordAccountBasic(localName, name, description, username),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubUsernamePasswordAccountExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testPlatformHubUsernamePasswordAccountBasic(localName string, name string, description string, username string) string {
	return fmt.Sprintf(`resource "octopusdeploy_platform_hub_username_password_account" "%s" {
		name        = "%s"
		description = "%s"
		username    = "%s"
		password    = "test_password_123456789"
	}`, localName, name, description, username)
}

func testPlatformHubUsernamePasswordAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Platform Hub Username-Password Account ID is set")
		}

		client := octoClient
		account, err := platformhubaccounts.GetByID(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error retrieving Platform Hub Username-Password Account: %s", err)
		}

		if account.GetID() != rs.Primary.ID {
			return fmt.Errorf("Platform Hub Username-Password Account not found")
		}

		return nil
	}
}

func testPlatformHubUsernamePasswordAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_platform_hub_username_password_account" {
			continue
		}

		client := octoClient
		account, err := platformhubaccounts.GetByID(client, rs.Primary.ID)
		if err == nil && account != nil {
			return fmt.Errorf("Platform Hub Username-Password Account (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
