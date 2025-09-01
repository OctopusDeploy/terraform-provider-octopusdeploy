package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployListeningTentacleDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_listening_tentacle_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tentacleUrl := "https://example-tentacle.local:10933/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccListeningTentacleDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccListeningTentacleDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tentacle_url", tentacleUrl),
					resource.TestCheckResourceAttr(prefix, "thumbprint", thumbprint),
					resource.TestCheckResourceAttr(prefix, "environments.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.0", "test-role"),
				),
				Config: testAccListeningTentacleDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, tentacleUrl, thumbprint),
			},
		},
	})
}

func TestAccOctopusDeployListeningTentacleDeploymentTargetUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_listening_tentacle_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tentacleUrl := "https://example-tentacle.local:10933/"
	newTentacleUrl := "https://new-tentacle.local:10933/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	newThumbprint := "ABCDEF1234567890ABCDEF1234567890ABCDEF12"

	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccListeningTentacleDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccListeningTentacleDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tentacle_url", tentacleUrl),
					resource.TestCheckResourceAttr(prefix, "thumbprint", thumbprint),
				),
				Config: testAccListeningTentacleDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, tentacleUrl, thumbprint),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccListeningTentacleDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "tentacle_url", newTentacleUrl),
					resource.TestCheckResourceAttr(prefix, "thumbprint", newThumbprint),
					resource.TestCheckResourceAttr(prefix, "roles.#", "2"),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "true"),
				),
				Config: testAccListeningTentacleDeploymentTargetUpdate(localName, environmentLocalName, environmentName, newName, newTentacleUrl, newThumbprint),
			},
		},
	})
}

func TestAccOctopusDeployListeningTentacleDeploymentTargetWithTenants(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_listening_tentacle_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tentacleUrl := "https://example-tentacle.local:10933/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccListeningTentacleDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccListeningTentacleDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", "TenantedOrUntenanted"),
					resource.TestCheckResourceAttr(prefix, "tenants.#", "1"),
				),
				Config: testAccListeningTentacleDeploymentTargetWithTenants(localName, environmentLocalName, environmentName, tenantLocalName, tenantName, name, tentacleUrl, thumbprint),
			},
		},
	})
}

func TestAccOctopusDeployListeningTentacleDeploymentTargetImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_listening_tentacle_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tentacleUrl := "https://example-tentacle.local:10933/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccListeningTentacleDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccListeningTentacleDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, tentacleUrl, thumbprint),
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
				ImportStateIdFunc: testAccListeningTentacleDeploymentTargetImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccListeningTentacleDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, tentacleUrl, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_listening_tentacle_deployment_target" "%s" {
		name         = "%s"
		tentacle_url = "%s"
		thumbprint   = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["test-role"]
	}`, environmentLocalName, environmentName, localName, name, tentacleUrl, thumbprint, environmentLocalName)
}

func testAccListeningTentacleDeploymentTargetUpdate(localName, environmentLocalName, environmentName, name, tentacleUrl, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_listening_tentacle_deployment_target" "%s" {
		name         = "%s"
		tentacle_url = "%s"
		thumbprint   = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["test-role", "web-server"]
		is_disabled  = true
		machine_policy_id = "MachinePolicies-1"
	}`, environmentLocalName, environmentName, localName, name, tentacleUrl, thumbprint, environmentLocalName)
}

func testAccListeningTentacleDeploymentTargetWithTenants(localName, environmentLocalName, environmentName, tenantLocalName, tenantName, name, tentacleUrl, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_tenant" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_listening_tentacle_deployment_target" "%s" {
		name         = "%s"
		tentacle_url = "%s"
		thumbprint   = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["test-role"]
		tenanted_deployment_participation = "TenantedOrUntenanted"
		tenants      = [octopusdeploy_tenant.%s.id]
	}`, environmentLocalName, environmentName, tenantLocalName, tenantName, localName, name, tentacleUrl, thumbprint, environmentLocalName, tenantLocalName)
}

func testAccListeningTentacleDeploymentTargetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		targetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), targetID); err != nil {
			return err
		}

		return nil
	}
}

func testAccListeningTentacleDeploymentTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_listening_tentacle_deployment_target" {
			continue
		}

		if target, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("listening tentacle deployment target (%s) still exists", target.GetID())
		}
	}

	return nil
}

func testAccListeningTentacleDeploymentTargetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}