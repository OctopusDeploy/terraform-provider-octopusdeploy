package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployPollingTentacleDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_polling_tentacle_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tentacleUrl := "poll://abcdef0123456789/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccPollingTentacleDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccPollingTentacleDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tentacle_url", tentacleUrl),
					resource.TestCheckResourceAttr(prefix, "thumbprint", thumbprint),
					resource.TestCheckResourceAttr(prefix, "environments.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.0", "polling-role"),
				),
				Config: testAccPollingTentacleDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, tentacleUrl, thumbprint),
			},
		},
	})
}

func TestAccOctopusDeployPollingTentacleDeploymentTargetUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_polling_tentacle_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tentacleUrl := "poll://abcdef0123456789/"
	newTentacleUrl := "poll://fedcba9876543210/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	newThumbprint := "ABCDEF1234567890ABCDEF1234567890ABCDEF12"

	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccPollingTentacleDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccPollingTentacleDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tentacle_url", tentacleUrl),
					resource.TestCheckResourceAttr(prefix, "thumbprint", thumbprint),
				),
				Config: testAccPollingTentacleDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, tentacleUrl, thumbprint),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccPollingTentacleDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "tentacle_url", newTentacleUrl),
					resource.TestCheckResourceAttr(prefix, "thumbprint", newThumbprint),
					resource.TestCheckResourceAttr(prefix, "roles.#", "2"),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "true"),
				),
				Config: testAccPollingTentacleDeploymentTargetUpdate(localName, environmentLocalName, environmentName, newName, newTentacleUrl, newThumbprint),
			},
		},
	})
}

func TestAccOctopusDeployPollingTentacleDeploymentTargetWithTenants(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_polling_tentacle_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tentacleUrl := "poll://abcdef0123456789/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccPollingTentacleDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccPollingTentacleDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", "TenantedOrUntenanted"),
					resource.TestCheckResourceAttr(prefix, "tenants.#", "1"),
				),
				Config: testAccPollingTentacleDeploymentTargetWithTenants(localName, environmentLocalName, environmentName, tenantLocalName, tenantName, name, tentacleUrl, thumbprint),
			},
		},
	})
}

func TestAccOctopusDeployPollingTentacleDeploymentTargetImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_polling_tentacle_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tentacleUrl := "poll://abcdef0123456789/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccPollingTentacleDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPollingTentacleDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, tentacleUrl, thumbprint),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"has_latest_calamari",
					"health_status",
					"is_in_process",
					"operating_system",
					"shell_name",
					"shell_version",
					"status",
					"status_summary",
					"tentacle_version_details",
					"uri",
					"thumbprint", // Thumbprint may not be returned by API on read
				},
				ImportStateIdFunc: testAccPollingTentacleDeploymentTargetImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccPollingTentacleDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, tentacleUrl, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_polling_tentacle_deployment_target" "%s" {
		name         = "%s"
		tentacle_url = "%s"
		thumbprint   = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["polling-role"]
	}`, environmentLocalName, environmentName, localName, name, tentacleUrl, thumbprint, environmentLocalName)
}

func testAccPollingTentacleDeploymentTargetUpdate(localName, environmentLocalName, environmentName, name, tentacleUrl, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_polling_tentacle_deployment_target" "%s" {
		name         = "%s"
		tentacle_url = "%s"
		thumbprint   = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["app-server", "polling-role"]
		is_disabled  = true
		machine_policy_id = "MachinePolicies-1"
	}`, environmentLocalName, environmentName, localName, name, tentacleUrl, thumbprint, environmentLocalName)
}

func testAccPollingTentacleDeploymentTargetWithTenants(localName, environmentLocalName, environmentName, tenantLocalName, tenantName, name, tentacleUrl, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_tenant" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_polling_tentacle_deployment_target" "%s" {
		name         = "%s"
		tentacle_url = "%s"
		thumbprint   = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["polling-role"]
		tenanted_deployment_participation = "TenantedOrUntenanted"
		tenants      = [octopusdeploy_tenant.%s.id]
	}`, environmentLocalName, environmentName, tenantLocalName, tenantName, localName, name, tentacleUrl, thumbprint, environmentLocalName, tenantLocalName)
}

func testAccPollingTentacleDeploymentTargetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		targetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), targetID); err != nil {
			return err
		}

		return nil
	}
}

func testAccPollingTentacleDeploymentTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_polling_tentacle_deployment_target" {
			continue
		}

		if target, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("polling tentacle deployment target (%s) still exists", target.GetID())
		}
	}

	return nil
}

func testAccPollingTentacleDeploymentTargetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}