package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubaccounts"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccPlatformHubGcpAccountCreate(t *testing.T) {
	t.Skip("Skipping test - Platform Hub GCP account API not available on test server")

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_gcp_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	jsonKey := `{"type":"service_account","project_id":"test-project"}`

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubGcpAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubGcpAccountBasic(localName, name, description, jsonKey),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGcpAccountExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
		},
	})
}

func TestAccPlatformHubGcpAccountUpdate(t *testing.T) {
	t.Skip("Skipping test - Platform Hub GCP account API not available on test server")

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_gcp_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	jsonKey := `{"type":"service_account","project_id":"test-project"}`

	updatedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedJsonKey := `{"type":"service_account","project_id":"updated-project"}`

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubGcpAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubGcpAccountBasic(localName, name, description, jsonKey),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGcpAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testPlatformHubGcpAccountBasic(localName, updatedName, updatedDescription, updatedJsonKey),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGcpAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccPlatformHubGcpAccountImport(t *testing.T) {
	t.Skip("Skipping test - Platform Hub GCP account API not available on test server")

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_gcp_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	jsonKey := `{"type":"service_account","project_id":"test-project"}`

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubGcpAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubGcpAccountBasic(localName, name, description, jsonKey),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGcpAccountExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"json_key"},
			},
		},
	})
}

func testPlatformHubGcpAccountBasic(localName string, name string, description string, jsonKey string) string {
	return fmt.Sprintf(`resource "octopusdeploy_platform_hub_gcp_account" "%s" {
		name        = "%s"
		description = "%s"
		json_key    = %q
	}`, localName, name, description, jsonKey)
}

func testPlatformHubGcpAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Platform Hub GCP Account ID is set")
		}

		client := octoClient
		account, err := platformhubaccounts.GetByID(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error retrieving Platform Hub GCP Account: %s", err)
		}

		if account.GetID() != rs.Primary.ID {
			return fmt.Errorf("Platform Hub GCP Account not found")
		}

		return nil
	}
}

func testPlatformHubGcpAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_platform_hub_gcp_account" {
			continue
		}

		client := octoClient
		account, err := platformhubaccounts.GetByID(client, rs.Primary.ID)
		if err == nil && account != nil {
			return fmt.Errorf("Platform Hub GCP Account (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
