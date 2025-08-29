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

func TestAccOctopusDeployAzureOIDCAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_openid_connect." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	applicationID := uuid.New()
	tenantID := uuid.New()
	subscriptionID := uuid.New()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureOIDCAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureOIDCAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "application_id", applicationID.String()),
					resource.TestCheckResourceAttr(prefix, "tenant_id", tenantID.String()),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttr(prefix, "audience", "api://AzureADTokenExchange"),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(core.TenantedDeploymentModeUntenanted)),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccAzureOIDCAccountBasic(localName, name, description, applicationID, tenantID, subscriptionID),
			},
		},
	})
}

func TestAccOctopusDeployAzureOIDCAccountUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_openid_connect." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	applicationID := uuid.New()
	tenantID := uuid.New()
	subscriptionID := uuid.New()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureOIDCAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureOIDCAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "application_id", applicationID.String()),
				),
				Config: testAccAzureOIDCAccountBasic(localName, name, description, applicationID, tenantID, subscriptionID),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureOIDCAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
					resource.TestCheckResourceAttr(prefix, "application_id", applicationID.String()),
				),
				Config: testAccAzureOIDCAccountBasic(localName, newName, newDescription, applicationID, tenantID, subscriptionID),
			},
		},
	})
}

func TestAccOctopusDeployAzureOIDCAccountWithSubjectKeys(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_openid_connect." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	applicationID := uuid.New()
	tenantID := uuid.New()
	subscriptionID := uuid.New()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureOIDCAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureOIDCAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "execution_subject_keys.#", "1"),
					resource.TestCheckResourceAttr(prefix, "execution_subject_keys.0", "space"),
					resource.TestCheckResourceAttr(prefix, "health_subject_keys.#", "1"),
					resource.TestCheckResourceAttr(prefix, "health_subject_keys.0", "target"),
					resource.TestCheckResourceAttr(prefix, "account_test_subject_keys.#", "1"),
					resource.TestCheckResourceAttr(prefix, "account_test_subject_keys.0", "type"),
				),
				Config: testAccAzureOIDCAccountWithSubjectKeys(localName, name, description, applicationID, tenantID, subscriptionID),
			},
		},
	})
}

func TestAccOctopusDeployAzureOIDCAccountImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_azure_openid_connect." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	applicationID := uuid.New()
	tenantID := uuid.New()
	subscriptionID := uuid.New()

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureOIDCAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureOIDCAccountBasic(localName, name, description, applicationID, tenantID, subscriptionID),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"audience"},
				ImportStateIdFunc:       testAccAzureOIDCAccountImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAzureOIDCAccountBasic(localName, name, description string, applicationID, tenantID, subscriptionID uuid.UUID) string {
	return fmt.Sprintf(`resource "octopusdeploy_azure_openid_connect" "%s" {
		name            = "%s"
		description     = "%s"
		application_id  = "%s"
		tenant_id       = "%s"
		subscription_id = "%s"
		space_id        = "Spaces-1"
		audience        = "api://AzureADTokenExchange"
	}`, localName, name, description, applicationID.String(), tenantID.String(), subscriptionID.String())
}

func testAccAzureOIDCAccountWithSubjectKeys(localName, name, description string, applicationID, tenantID, subscriptionID uuid.UUID) string {
	return fmt.Sprintf(`resource "octopusdeploy_azure_openid_connect" "%s" {
		name                      = "%s"
		description               = "%s"
		application_id            = "%s"
		tenant_id                 = "%s"
		subscription_id           = "%s"
		space_id                  = "Spaces-1"
		audience                  = "api://AzureADTokenExchange"
		execution_subject_keys    = ["space"]
		health_subject_keys       = ["target"]
		account_test_subject_keys = ["type"]
	}`, localName, name, description, applicationID.String(), tenantID.String(), subscriptionID.String())
}

func testAccAzureOIDCAccountExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), accountID); err != nil {
			return err
		}

		return nil
	}
}

func testAccAzureOIDCAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_azure_openid_connect" {
			continue
		}

		if account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("Azure OIDC account (%s) still exists", account.GetID())
		}
	}

	return nil
}

func testAccAzureOIDCAccountImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}