package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployAzureServicePrincipalAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_service_principal." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	applicationID := uuid.New()
	tenantID := uuid.New()
	subscriptionID := uuid.New()
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureServicePrincipalAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureServicePrincipalAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "application_id", applicationID.String()),
					resource.TestCheckResourceAttr(prefix, "tenant_id", tenantID.String()),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "space_id"),
					resource.TestCheckResourceAttr(prefix, "tenants.#", "0"),
				),
				Config: testAccAzureServicePrincipalAccountBasic(localName, name, description, applicationID, tenantID, subscriptionID, password, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeployAzureServicePrincipalAccountUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_service_principal." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	applicationID := uuid.New()
	tenantID := uuid.New()
	subscriptionID := uuid.New()
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newPassword := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureServicePrincipalAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureServicePrincipalAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "password", password),
				),
				Config: testAccAzureServicePrincipalAccountBasic(localName, name, description, applicationID, tenantID, subscriptionID, password, tenantedDeploymentParticipation),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureServicePrincipalAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
					resource.TestCheckResourceAttr(prefix, "password", newPassword),
				),
				Config: testAccAzureServicePrincipalAccountBasic(localName, name, newDescription, applicationID, tenantID, subscriptionID, newPassword, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeployAzureServicePrincipalAccountMinimal(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_service_principal." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	applicationID := uuid.New()
	tenantID := uuid.New()
	subscriptionID := uuid.New()
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureServicePrincipalAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureServicePrincipalAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "application_id", applicationID.String()),
					resource.TestCheckResourceAttr(prefix, "tenant_id", tenantID.String()),
					resource.TestCheckResourceAttr(prefix, "subscription_id", subscriptionID.String()),
					resource.TestCheckResourceAttr(prefix, "password", password),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccAzureServicePrincipalAccountMinimal(localName, name, applicationID, tenantID, subscriptionID, password),
			},
		},
	})
}

func TestAccOctopusDeployAzureServicePrincipalAccountTenanted(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_service_principal." + localName
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	applicationID := uuid.New()
	tenantID := uuid.New()
	subscriptionID := uuid.New()
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureServicePrincipalAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureServicePrincipalAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(prefix, "tenants.#", "1"),
				),
				Config: testAccAzureServicePrincipalAccountTenanted(localName, tenantLocalName, name, description, tenantName, applicationID, tenantID, subscriptionID, password, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeployAzureServicePrincipalAccountImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_azure_service_principal." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	applicationID := uuid.New()
	tenantID := uuid.New()
	subscriptionID := uuid.New()
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureServicePrincipalAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureServicePrincipalAccountBasic(localName, name, description, applicationID, tenantID, subscriptionID, password, tenantedDeploymentParticipation),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password"}, // sensitive field
				ImportStateIdFunc: testAccAzureServicePrincipalAccountImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAzureServicePrincipalAccountBasic(localName string, name string, description string, applicationID uuid.UUID, tenantID uuid.UUID, subscriptionID uuid.UUID, password string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_azure_service_principal" "%s" {
		name                              = "%s"
		description                       = "%s"
		application_id                    = "%s"
		tenant_id                         = "%s"
		subscription_id                   = "%s"
		password                          = "%s"
		tenanted_deployment_participation = "%s"
		tenants                           = []
	}`, localName, name, description, applicationID, tenantID, subscriptionID, password, tenantedDeploymentParticipation)
}

func testAccAzureServicePrincipalAccountMinimal(localName string, name string, applicationID uuid.UUID, tenantID uuid.UUID, subscriptionID uuid.UUID, password string) string {
	return fmt.Sprintf(`resource "octopusdeploy_azure_service_principal" "%s" {
		name            = "%s"
		application_id  = "%s"
		tenant_id       = "%s"
		subscription_id = "%s"
		password        = "%s"
	}`, localName, name, applicationID, tenantID, subscriptionID, password)
}

func testAccAzureServicePrincipalAccountTenanted(localName string, tenantLocalName string, name string, description string, tenantName string, applicationID uuid.UUID, tenantID uuid.UUID, subscriptionID uuid.UUID, password string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_tenant" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_azure_service_principal" "%s" {
		name                              = "%s"
		description                       = "%s"
		application_id                    = "%s"
		tenant_id                         = "%s"
		subscription_id                   = "%s"
		password                          = "%s"
		tenanted_deployment_participation = "%s"
		tenants                           = [octopusdeploy_tenant.%s.id]
	}`, tenantLocalName, tenantName, localName, name, description, applicationID, tenantID, subscriptionID, password, tenantedDeploymentParticipation, tenantLocalName)
}

func testAccAzureServicePrincipalAccountExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), accountID); err != nil {
			return err
		}

		return nil
	}
}

func testAccAzureServicePrincipalAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_azure_service_principal" {
			continue
		}

		if account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("Azure service principal account (%s) still exists", account.GetID())
		}
	}

	return nil
}

func testAccAzureServicePrincipalAccountImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}