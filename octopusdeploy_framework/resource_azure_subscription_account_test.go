package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployAzureSubscriptionAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_subscription_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	subscriptionID := uuid.New()
	azureEnvironment := "AzureChinaCloud"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureSubscriptionAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureSubscriptionAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, "azure_environment", azureEnvironment),
					resource.TestCheckResourceAttr(prefix, "management_endpoint", "https://management.core.chinacloudapi.cn/"),
					resource.TestCheckResourceAttr(prefix, "storage_endpoint_suffix", "core.chinacloudapi.cn"),
					resource.TestCheckResourceAttr(prefix, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(core.TenantedDeploymentModeUntenanted)),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccAzureSubscriptionAccountBasic(localName, name, description, subscriptionID, azureEnvironment),
			},
		},
	})
}

func TestAccOctopusDeployAzureSubscriptionAccountUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_subscription_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	subscriptionID := uuid.New()
	azureEnvironment := "AzureChinaCloud"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureSubscriptionAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureSubscriptionAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
				),
				Config: testAccAzureSubscriptionAccountBasic(localName, name, description, subscriptionID, azureEnvironment),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureSubscriptionAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
				),
				Config: testAccAzureSubscriptionAccountBasic(localName, newName, newDescription, subscriptionID, azureEnvironment),
			},
		},
	})
}

func TestAccOctopusDeployAzureSubscriptionAccountMinimal(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_subscription_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	subscriptionID := uuid.New()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureSubscriptionAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureSubscriptionAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccAzureSubscriptionAccountMinimal(localName, name, subscriptionID),
			},
		},
	})
}

func TestAccOctopusDeployAzureSubscriptionAccountTenanted(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_subscription_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	subscriptionID := uuid.New()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureSubscriptionAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureSubscriptionAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(core.TenantedDeploymentModeTenantedOrUntenanted)),
				),
				Config: testAccAzureSubscriptionAccountTenanted(localName, name, subscriptionID),
			},
		},
	})
}

func TestAccOctopusDeployAzureSubscriptionAccountImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_azure_subscription_account." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	subscriptionID := uuid.New()
	azureEnvironment := "AzureChinaCloud"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureSubscriptionAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureSubscriptionAccountBasic(localName, name, description, subscriptionID, azureEnvironment),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "certificate_thumbprint"},
				ImportStateIdFunc:       testAccAzureSubscriptionAccountImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAzureSubscriptionAccountBasic(localName, name, description string, subscriptionID uuid.UUID, azureEnvironment string) string {
	return fmt.Sprintf(`resource "octopusdeploy_azure_subscription_account" "%s" {
		name                    = "%s"
		description             = "%s"
		subscription_id         = "%s"
		azure_environment       = "%s"
		management_endpoint     = "https://management.core.chinacloudapi.cn/"
		storage_endpoint_suffix = "core.chinacloudapi.cn"
		space_id                = "Spaces-1"
	}`, localName, name, description, subscriptionID.String(), azureEnvironment)
}

func testAccAzureSubscriptionAccountMinimal(localName, name string, subscriptionID uuid.UUID) string {
	return fmt.Sprintf(`resource "octopusdeploy_azure_subscription_account" "%s" {
		name            = "%s"
		subscription_id = "%s"
		space_id        = "Spaces-1"
	}`, localName, name, subscriptionID.String())
}

func testAccAzureSubscriptionAccountTenanted(localName, name string, subscriptionID uuid.UUID) string {
	return fmt.Sprintf(`resource "octopusdeploy_azure_subscription_account" "%s" {
		name                              = "%s"
		subscription_id                   = "%s"
		space_id                          = "Spaces-1"
		tenanted_deployment_participation = "%s"
	}`, localName, name, subscriptionID.String(), string(core.TenantedDeploymentModeTenantedOrUntenanted))
}

func testAccAzureSubscriptionAccountExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), accountID); err != nil {
			return err
		}

		return nil
	}
}

func testAccAzureSubscriptionAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_azure_subscription_account" {
			continue
		}

		if account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("Azure Subscription account (%s) still exists", account.GetID())
		}
	}

	return nil
}

func testAccAzureSubscriptionAccountImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}