package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployKubernetesAgentDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_kubernetes_agent_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	uri := "poll://abcdef0123456789/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesAgentDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesAgentDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "uri", uri),
					resource.TestCheckResourceAttr(prefix, "thumbprint", thumbprint),
					resource.TestCheckResourceAttr(prefix, "environments.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.0", "k8s-agent"),
					resource.TestCheckResourceAttr(prefix, "communication_mode", "Polling"),
				),
				Config: testAccKubernetesAgentDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, uri, thumbprint),
			},
		},
	})
}

func TestAccOctopusDeployKubernetesAgentDeploymentTargetUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_kubernetes_agent_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	uri := "poll://abcdef0123456789/"
	newUri := "poll://fedcba9876543210/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	newThumbprint := "ABCDEF1234567890ABCDEF1234567890ABCDEF12"

	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesAgentDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesAgentDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "uri", uri),
					resource.TestCheckResourceAttr(prefix, "thumbprint", thumbprint),
				),
				Config: testAccKubernetesAgentDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, uri, thumbprint),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesAgentDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "uri", newUri),
					resource.TestCheckResourceAttr(prefix, "thumbprint", newThumbprint),
					resource.TestCheckResourceAttr(prefix, "roles.#", "2"),
					resource.TestCheckResourceAttr(prefix, "default_namespace", "production"),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "true"),
				),
				Config: testAccKubernetesAgentDeploymentTargetUpdate(localName, environmentLocalName, environmentName, newName, newUri, newThumbprint),
			},
		},
	})
}

func TestAccOctopusDeployKubernetesAgentDeploymentTargetWithTenants(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_kubernetes_agent_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	uri := "poll://abcdef0123456789/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesAgentDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesAgentDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "tenanted_deployment_participation", "TenantedOrUntenanted"),
					resource.TestCheckResourceAttr(prefix, "tenants.#", "1"),
				),
				Config: testAccKubernetesAgentDeploymentTargetWithTenants(localName, environmentLocalName, environmentName, tenantLocalName, tenantName, name, uri, thumbprint),
			},
		},
	})
}

func TestAccOctopusDeployKubernetesAgentDeploymentTargetImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_kubernetes_agent_deployment_target." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	uri := "poll://abcdef0123456789/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesAgentDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAgentDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, uri, thumbprint),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"agent_version",
					"agent_tentacle_version",
					"agent_upgrade_status",
					"agent_helm_release_name",
					"agent_kubernetes_namespace",
					"default_namespace",
					"machine_policy_id",
					"space_id",
				},
				ImportStateIdFunc: testAccKubernetesAgentDeploymentTargetImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccKubernetesAgentDeploymentTargetBasic(localName, environmentLocalName, environmentName, name, uri, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_kubernetes_agent_deployment_target" "%s" {
		name         = "%s"
		uri          = "%s"
		thumbprint   = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["k8s-agent"]
		communication_mode = "Polling"
		tenanted_deployment_participation = "Untenanted"
		tenants = []
		tenant_tags = []
	}`, environmentLocalName, environmentName, localName, name, uri, thumbprint, environmentLocalName)
}

func testAccKubernetesAgentDeploymentTargetUpdate(localName, environmentLocalName, environmentName, name, uri, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_kubernetes_agent_deployment_target" "%s" {
		name              = "%s"
		uri               = "%s"
		thumbprint        = "%s"
		environments      = [octopusdeploy_environment.%s.id]
		roles             = ["app-server", "k8s-agent"]
		default_namespace = "production"
		is_disabled       = true
		machine_policy_id = "MachinePolicies-1"
		communication_mode = "Polling"
		tenanted_deployment_participation = "Untenanted"
		tenants = []
		tenant_tags = []
	}`, environmentLocalName, environmentName, localName, name, uri, thumbprint, environmentLocalName)
}

func testAccKubernetesAgentDeploymentTargetWithTenants(localName, environmentLocalName, environmentName, tenantLocalName, tenantName, name, uri, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_tenant" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_kubernetes_agent_deployment_target" "%s" {
		name         = "%s"
		uri          = "%s"
		thumbprint   = "%s"
		environments = [octopusdeploy_environment.%s.id]
		roles        = ["k8s-agent"]
		tenanted_deployment_participation = "TenantedOrUntenanted"
		tenants      = [octopusdeploy_tenant.%s.id]
		tenant_tags  = []
		communication_mode = "Polling"
	}`, environmentLocalName, environmentName, tenantLocalName, tenantName, localName, name, uri, thumbprint, environmentLocalName, tenantLocalName)
}

func testAccKubernetesAgentDeploymentTargetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		targetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), targetID); err != nil {
			return err
		}

		return nil
	}
}

func testAccKubernetesAgentDeploymentTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_kubernetes_agent_deployment_target" {
			continue
		}

		if target, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("kubernetes agent deployment target (%s) still exists", target.GetID())
		}
	}

	return nil
}

func testAccKubernetesAgentDeploymentTargetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}