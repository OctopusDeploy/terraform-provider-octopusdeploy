package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeploySSHConnectionDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_ssh_connection_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	host := "192.168.1.100"
	port := 22
	fingerprint := "00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccSSHConnectionDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccSSHConnectionDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "host", host),
					resource.TestCheckResourceAttr(prefix, "port", fmt.Sprintf("%d", port)),
					resource.TestCheckResourceAttr(prefix, "fingerprint", fingerprint),
					resource.TestCheckResourceAttr(prefix, "roles.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.0", "ssh"),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccSSHConnectionDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, host, port, fingerprint),
			},
		},
	})
}

func TestAccOctopusDeploySSHConnectionDeploymentTargetUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_ssh_connection_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	host := "192.168.1.100"
	newHost := "192.168.1.200"
	port := 22
	fingerprint := "00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff"
	newFingerprint := "ff:ee:dd:cc:bb:aa:99:88:77:66:55:44:33:22:11:00"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccSSHConnectionDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccSSHConnectionDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "host", host),
					resource.TestCheckResourceAttr(prefix, "fingerprint", fingerprint),
				),
				Config: testAccSSHConnectionDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, host, port, fingerprint),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccSSHConnectionDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "host", newHost),
					resource.TestCheckResourceAttr(prefix, "fingerprint", newFingerprint),
				),
				Config: testAccSSHConnectionDeploymentTargetBasic(localName, environmentLocalName, newName, environmentName, newHost, port, newFingerprint),
			},
		},
	})
}

func TestAccOctopusDeploySSHConnectionDeploymentTargetWithAccount(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_ssh_connection_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	host := "192.168.1.100"
	port := 22
	fingerprint := "00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccSSHConnectionDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccSSHConnectionDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "host", host),
					resource.TestCheckResourceAttrPair(prefix, "account_id", "octopusdeploy_ssh_key_account."+accountLocalName, "id"),
				),
				Config: testAccSSHConnectionDeploymentTargetWithAccount(localName, environmentLocalName, accountLocalName, name, environmentName, accountName, host, port, fingerprint),
			},
		},
	})
}

func TestAccOctopusDeploySSHConnectionDeploymentTargetImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_ssh_connection_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	host := "192.168.1.100"
	port := 22
	fingerprint := "00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccSSHConnectionDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSSHConnectionDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, host, port, fingerprint),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"has_latest_calamari", "health_status", "status_summary"},
				ImportStateIdFunc:       testAccSSHConnectionDeploymentTargetImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccSSHConnectionDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, host string, port int, fingerprint string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "%s" {
			name        = "%s"
			description = "Test environment for SSH connection deployment target"
		}

		resource "octopusdeploy_ssh_key_account" "%s_account" {
			name         = "%s-account"
			username     = "testuser"
			private_key_file = "-----BEGIN OPENSSH PRIVATE KEY-----\\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAFwAAAAdzc2gtcn\\nNhAAAAAwEAAQAAAQEA0VktClrpGk8ijTETjc7+Nzgu+rzVNAzRYRbOXZw4/rZAoYgEXA\\nJgKa4KWAkp6Kn++vZJ7Uk8l1XrzWcTfKG4+KxNdGEPqe+n5Sxv4zWoE5n7GQZJ+hYn\\n7Q5Z0Sv1Z1M2Z2F2N+/ZqZ7Y9J3YXX5TfjD0oV8fAJ0xN7DQAA5S/4DQAA5S/4DQ\\nYQ==\\n-----END OPENSSH PRIVATE KEY-----"
			space_id     = "Spaces-1"
		}

		resource "octopusdeploy_ssh_connection_deployment_target" "%s" {
			name         = "%s"
			environments = [octopusdeploy_environment.%s.id]
			roles        = ["ssh"]
			host         = "%s"
			port         = %d
			fingerprint  = "%s"
			account_id   = octopusdeploy_ssh_key_account.%s_account.id
		}`, environmentLocalName, environmentName, localName, name, localName, name, environmentLocalName, host, port, fingerprint, localName)
}

func testAccSSHConnectionDeploymentTargetWithAccount(localName, environmentLocalName, accountLocalName, name, environmentName, accountName, host string, port int, fingerprint string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "%s" {
			name        = "%s"
			description = "Test environment for SSH connection deployment target"
		}

		resource "octopusdeploy_ssh_key_account" "%s" {
			name        = "%s"
			username    = "testuser"
			private_key_file = "-----BEGIN OPENSSH PRIVATE KEY-----\\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAFwAAAAdzc2gtcn\\nNhAAAAAwEAAQAAAQEA0VktClrpGk8ijTETjc7+Nzgu+rzVNAzRYRbOXZw4/rZAoYgEXA\\nJgKa4KWAkp6Kn++vZJ7Uk8l1XrzWcTfKG4+KxNdGEPqe+n5Sxv4zWoE5n7GQZJ+hYn\\n7Q5Z0Sv1Z1M2Z2F2N+/ZqZ7Y9J3YXX5TfjD0oV8fAJ0xN7DQAA5S/4DQAA5S/4DQ\\nYQ==\\n-----END OPENSSH PRIVATE KEY-----"
			space_id    = "Spaces-1"
		}

		resource "octopusdeploy_ssh_connection_deployment_target" "%s" {
			name         = "%s"
			environments = [octopusdeploy_environment.%s.id]
			roles        = ["ssh"]
			host         = "%s"
			port         = %d
			fingerprint  = "%s"
			account_id   = octopusdeploy_ssh_key_account.%s.id
		}`, environmentLocalName, environmentName, accountLocalName, accountName, localName, name, environmentLocalName, host, port, fingerprint, accountLocalName)
}

func testAccSSHConnectionDeploymentTargetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		deploymentTargetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), deploymentTargetID); err != nil {
			return err
		}

		return nil
	}
}

func testAccSSHConnectionDeploymentTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_ssh_connection_deployment_target" {
			continue
		}

		if deploymentTarget, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("SSH connection deployment target (%s) still exists", deploymentTarget.GetID())
		}
	}

	return nil
}

func testAccSSHConnectionDeploymentTargetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}