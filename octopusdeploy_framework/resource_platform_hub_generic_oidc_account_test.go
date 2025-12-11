package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubaccounts"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccPlatformHubGenericOidcAccountCreate(t *testing.T) {
	t.Skip("Skipping test - Platform Hub Generic OIDC account API not available on test server")

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_generic_oidc_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	audience := "api://default"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubGenericOidcAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubGenericOidcAccountBasic(localName, name, description, audience),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGenericOidcAccountExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "audience", audience),
				),
			},
		},
	})
}

func TestAccPlatformHubGenericOidcAccountUpdate(t *testing.T) {
	t.Skip("Skipping test - Platform Hub Generic OIDC account API not available on test server")

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_generic_oidc_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	audience := "api://default"

	updatedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedAudience := "api://updated"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubGenericOidcAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubGenericOidcAccountBasic(localName, name, description, audience),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGenericOidcAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "audience", audience),
				),
			},
			{
				Config: testPlatformHubGenericOidcAccountBasic(localName, updatedName, updatedDescription, updatedAudience),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGenericOidcAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "audience", updatedAudience),
				),
			},
		},
	})
}

func TestAccPlatformHubGenericOidcAccountImport(t *testing.T) {
	t.Skip("Skipping test - Platform Hub Generic OIDC account API not available on test server")

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_generic_oidc_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	audience := "api://default"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubGenericOidcAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubGenericOidcAccountBasic(localName, name, description, audience),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGenericOidcAccountExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testPlatformHubGenericOidcAccountBasic(localName string, name string, description string, audience string) string {
	return fmt.Sprintf(`resource "octopusdeploy_platform_hub_generic_oidc_account" "%s" {
		name                   = "%s"
		description            = "%s"
		execution_subject_keys = ["space", "project"]
		audience               = "%s"
	}`, localName, name, description, audience)
}

func testPlatformHubGenericOidcAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Platform Hub Generic OIDC Account ID is set")
		}

		client := octoClient
		account, err := platformhubaccounts.GetByID(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error retrieving Platform Hub Generic OIDC Account: %s", err)
		}

		if account.GetID() != rs.Primary.ID {
			return fmt.Errorf("Platform Hub Generic OIDC Account not found")
		}

		return nil
	}
}

func testPlatformHubGenericOidcAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_platform_hub_generic_oidc_account" {
			continue
		}

		client := octoClient
		account, err := platformhubaccounts.GetByID(client, rs.Primary.ID)
		if err == nil && account != nil {
			return fmt.Errorf("Platform Hub Generic OIDC Account (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
