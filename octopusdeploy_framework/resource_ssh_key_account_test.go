package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeploySSHKeyAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_ssh_key_account." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	privateKeyFile := "-----BEGIN RSA PRIVATE KEY-----\\nMIIEpAIBAAKCAQEA" + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "\\n-----END RSA PRIVATE KEY-----"
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	passphrase := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccSSHKeyAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccSSHKeyAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "username", username),
					resource.TestCheckResourceAttr(prefix, "private_key_passphrase", passphrase),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "space_id"),
					resource.TestCheckResourceAttr(prefix, "tenants.#", "0"),
				),
				Config: testAccSSHKeyAccountBasic(localName, name, privateKeyFile, username, passphrase, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeploySSHKeyAccountUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_ssh_key_account." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	privateKeyFile := "-----BEGIN RSA PRIVATE KEY-----\\nMIIEpAIBAAKCAQEA" + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "\\n-----END RSA PRIVATE KEY-----"
	newPrivateKeyFile := "-----BEGIN RSA PRIVATE KEY-----\\nMIIEpAIBAAKCAQEA" + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "\\n-----END RSA PRIVATE KEY-----"
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newUsername := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	passphrase := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newPassphrase := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccSSHKeyAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccSSHKeyAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "username", username),
					resource.TestCheckResourceAttr(prefix, "private_key_passphrase", passphrase),
				),
				Config: testAccSSHKeyAccountBasic(localName, name, privateKeyFile, username, passphrase, tenantedDeploymentParticipation),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccSSHKeyAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "username", newUsername),
					resource.TestCheckResourceAttr(prefix, "private_key_passphrase", newPassphrase),
				),
				Config: testAccSSHKeyAccountBasic(localName, name, newPrivateKeyFile, newUsername, newPassphrase, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeploySSHKeyAccountWithDescription(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_ssh_key_account." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	privateKeyFile := "-----BEGIN RSA PRIVATE KEY-----\\nMIIEpAIBAAKCAQEA" + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "\\n-----END RSA PRIVATE KEY-----"
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	passphrase := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccSSHKeyAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccSSHKeyAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "username", username),
				),
				Config: testAccSSHKeyAccountWithDescription(localName, name, description, privateKeyFile, username, passphrase, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeploySSHKeyAccountTenanted(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_ssh_key_account." + localName
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	privateKeyFile := "-----BEGIN RSA PRIVATE KEY-----\\nMIIEpAIBAAKCAQEA" + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "\\n-----END RSA PRIVATE KEY-----"
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	passphrase := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccSSHKeyAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccSSHKeyAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", string(tenantedDeploymentParticipation)),
					resource.TestCheckResourceAttr(prefix, "tenants.#", "1"),
				),
				Config: testAccSSHKeyAccountTenanted(localName, tenantLocalName, name, tenantName, privateKeyFile, username, passphrase, tenantedDeploymentParticipation),
			},
		},
	})
}

func TestAccOctopusDeploySSHKeyAccountImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_ssh_key_account." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	privateKeyFile := "-----BEGIN RSA PRIVATE KEY-----\\nMIIEpAIBAAKCAQEA" + acctest.RandStringFromCharSet(20, acctest.CharSetAlpha) + "\\n-----END RSA PRIVATE KEY-----"
	username := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	passphrase := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantedDeploymentParticipation := core.TenantedDeploymentModeTenantedOrUntenanted

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccSSHKeyAccountCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSSHKeyAccountBasic(localName, name, privateKeyFile, username, passphrase, tenantedDeploymentParticipation),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"private_key_file", "private_key_passphrase"}, // sensitive fields
				ImportStateIdFunc: testAccSSHKeyAccountImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccSSHKeyAccountBasic(localName string, name string, privateKeyFile string, username string, passphrase string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_ssh_key_account" "%s" {
		name                              = "%s"
		private_key_file                  = "%s"
		username                          = "%s"
		private_key_passphrase            = "%s"
		tenanted_deployment_participation = "%s"
		tenants                           = []
	}`, localName, name, privateKeyFile, username, passphrase, tenantedDeploymentParticipation)
}

func testAccSSHKeyAccountWithDescription(localName string, name string, description string, privateKeyFile string, username string, passphrase string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`resource "octopusdeploy_ssh_key_account" "%s" {
		name                              = "%s"
		description                       = "%s"
		private_key_file                  = "%s"
		username                          = "%s"
		private_key_passphrase            = "%s"
		tenanted_deployment_participation = "%s"
		tenants                           = []
	}`, localName, name, description, privateKeyFile, username, passphrase, tenantedDeploymentParticipation)
}

func testAccSSHKeyAccountTenanted(localName string, tenantLocalName string, name string, tenantName string, privateKeyFile string, username string, passphrase string, tenantedDeploymentParticipation core.TenantedDeploymentMode) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_tenant" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_ssh_key_account" "%s" {
		name                              = "%s"
		private_key_file                  = "%s"
		username                          = "%s"
		private_key_passphrase            = "%s"
		tenanted_deployment_participation = "%s"
		tenants                           = [octopusdeploy_tenant.%s.id]
	}`, tenantLocalName, tenantName, localName, name, privateKeyFile, username, passphrase, tenantedDeploymentParticipation, tenantLocalName)
}

func testAccSSHKeyAccountExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		accountID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), accountID); err != nil {
			return err
		}

		return nil
	}
}

func testAccSSHKeyAccountCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_ssh_key_account" {
			continue
		}

		if account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("SSH key account (%s) still exists", account.GetID())
		}
	}

	return nil
}

func testAccSSHKeyAccountImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}