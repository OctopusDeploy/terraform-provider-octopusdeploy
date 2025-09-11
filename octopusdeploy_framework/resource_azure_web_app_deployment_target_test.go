package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployAzureWebAppDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_web_app_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	resourceGroupName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	webAppName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureWebAppDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureWebAppDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(prefix, "web_app_name", webAppName),
					resource.TestCheckResourceAttrPair(prefix, "account_id", "octopusdeploy_azure_service_principal."+accountLocalName, "id"),
					resource.TestCheckResourceAttr(prefix, "roles.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.0", "azure-web-app"),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccAzureWebAppDeploymentTargetBasic(localName, environmentLocalName, accountLocalName, name, environmentName, accountName, resourceGroupName, webAppName),
			},
		},
	})
}

func TestAccOctopusDeployAzureWebAppDeploymentTargetUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_web_app_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	resourceGroupName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	newResourceGroupName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	webAppName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	newWebAppName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureWebAppDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureWebAppDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(prefix, "web_app_name", webAppName),
				),
				Config: testAccAzureWebAppDeploymentTargetBasic(localName, environmentLocalName, accountLocalName, name, environmentName, accountName, resourceGroupName, webAppName),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureWebAppDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "resource_group_name", newResourceGroupName),
					resource.TestCheckResourceAttr(prefix, "web_app_name", newWebAppName),
				),
				Config: testAccAzureWebAppDeploymentTargetBasic(localName, environmentLocalName, accountLocalName, newName, environmentName, accountName, newResourceGroupName, newWebAppName),
			},
		},
	})
}

func TestAccOctopusDeployAzureWebAppDeploymentTargetWithSlot(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_web_app_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	resourceGroupName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	webAppName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	webAppSlotName := "staging"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureWebAppDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureWebAppDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(prefix, "web_app_name", webAppName),
					resource.TestCheckResourceAttr(prefix, "web_app_slot_name", webAppSlotName),
				),
				Config: testAccAzureWebAppDeploymentTargetWithSlot(localName, environmentLocalName, accountLocalName, name, environmentName, accountName, resourceGroupName, webAppName, webAppSlotName),
			},
		},
	})
}

func TestAccOctopusDeployAzureWebAppDeploymentTargetImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_azure_web_app_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	resourceGroupName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	webAppName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureWebAppDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureWebAppDeploymentTargetBasic(localName, environmentLocalName, accountLocalName, name, environmentName, accountName, resourceGroupName, webAppName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"has_latest_calamari", "health_status", "status_summary"},
				ImportStateIdFunc:       testAccAzureWebAppDeploymentTargetImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAzureWebAppDeploymentTargetBasic(localName, environmentLocalName, accountLocalName, name, environmentName, accountName, resourceGroupName, webAppName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "%s" {
			name        = "%s"
			description = "Test environment for Azure Web App deployment target"
		}

		resource "octopusdeploy_azure_service_principal" "%s" {
			application_id = "00000000-1111-2222-3333-444444444444"
			name           = "%s"
			password       = "test-client-secret"
			subscription_id = "00000000-1111-2222-3333-444444444444"
			tenant_id      = "00000000-1111-2222-3333-444444444444"
		}

		resource "octopusdeploy_azure_web_app_deployment_target" "%s" {
			name                = "%s"
			environments        = [octopusdeploy_environment.%s.id]
			roles               = ["azure-web-app"]
			account_id          = octopusdeploy_azure_service_principal.%s.id
			resource_group_name = "%s"
			web_app_name        = "%s"
		}`, environmentLocalName, environmentName, accountLocalName, accountName, localName, name, environmentLocalName, accountLocalName, resourceGroupName, webAppName)
}

func testAccAzureWebAppDeploymentTargetWithSlot(localName, environmentLocalName, accountLocalName, name, environmentName, accountName, resourceGroupName, webAppName, webAppSlotName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "%s" {
			name        = "%s"
			description = "Test environment for Azure Web App deployment target"
		}

		resource "octopusdeploy_azure_service_principal" "%s" {
			application_id = "00000000-1111-2222-3333-444444444444"
			name           = "%s"
			password       = "test-client-secret"
			subscription_id = "00000000-1111-2222-3333-444444444444"
			tenant_id      = "00000000-1111-2222-3333-444444444444"
		}

		resource "octopusdeploy_azure_web_app_deployment_target" "%s" {
			name                = "%s"
			environments        = [octopusdeploy_environment.%s.id]
			roles               = ["azure-web-app"]
			account_id          = octopusdeploy_azure_service_principal.%s.id
			resource_group_name = "%s"
			web_app_name        = "%s"
			web_app_slot_name   = "%s"
		}`, environmentLocalName, environmentName, accountLocalName, accountName, localName, name, environmentLocalName, accountLocalName, resourceGroupName, webAppName, webAppSlotName)
}

func testAccAzureWebAppDeploymentTargetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		deploymentTargetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), deploymentTargetID); err != nil {
			return err
		}

		return nil
	}
}

func testAccAzureWebAppDeploymentTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_azure_web_app_deployment_target" {
			continue
		}

		if deploymentTarget, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("Azure Web App deployment target (%s) still exists", deploymentTarget.GetID())
		}
	}

	return nil
}

func testAccAzureWebAppDeploymentTargetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}