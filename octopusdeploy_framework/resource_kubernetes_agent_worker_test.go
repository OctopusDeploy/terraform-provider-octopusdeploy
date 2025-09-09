package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployKubernetesAgentWorkerBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_kubernetes_agent_worker." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	uri := "poll://abcdef0123456789/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	workerPoolLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerPoolName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesAgentWorkerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesAgentWorkerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "uri", uri),
					resource.TestCheckResourceAttr(prefix, "thumbprint", thumbprint),
					resource.TestCheckResourceAttr(prefix, "worker_pool_ids.#", "1"),
					resource.TestCheckResourceAttr(prefix, "communication_mode", "Polling"),
				),
				Config: testAccKubernetesAgentWorkerBasic(localName, workerPoolLocalName, workerPoolName, name, uri, thumbprint),
			},
		},
	})
}

func TestAccOctopusDeployKubernetesAgentWorkerUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_kubernetes_agent_worker." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	uri := "poll://abcdef0123456789/"
	newUri := "poll://fedcba9876543210/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	newThumbprint := "ABCDEF1234567890ABCDEF1234567890ABCDEF12"

	workerPoolLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerPoolName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerPool2LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerPool2Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesAgentWorkerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesAgentWorkerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "uri", uri),
					resource.TestCheckResourceAttr(prefix, "thumbprint", thumbprint),
					resource.TestCheckResourceAttr(prefix, "worker_pool_ids.#", "1"),
				),
				Config: testAccKubernetesAgentWorkerBasic(localName, workerPoolLocalName, workerPoolName, name, uri, thumbprint),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesAgentWorkerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "uri", newUri),
					resource.TestCheckResourceAttr(prefix, "thumbprint", newThumbprint),
					resource.TestCheckResourceAttr(prefix, "worker_pool_ids.#", "2"),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "true"),
				),
				Config: testAccKubernetesAgentWorkerUpdate(localName, workerPoolLocalName, workerPoolName, workerPool2LocalName, workerPool2Name, newName, newUri, newThumbprint),
			},
		},
	})
}

func TestAccOctopusDeployKubernetesAgentWorkerWithUpgradeLocked(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_kubernetes_agent_worker." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	uri := "poll://abcdef0123456789/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	workerPoolLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerPoolName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesAgentWorkerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesAgentWorkerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "upgrade_locked", "true"),
				),
				Config: testAccKubernetesAgentWorkerWithUpgradeLocked(localName, workerPoolLocalName, workerPoolName, name, uri, thumbprint),
			},
		},
	})
}

func TestAccOctopusDeployKubernetesAgentWorkerImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_kubernetes_agent_worker." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	uri := "poll://abcdef0123456789/"
	thumbprint := "1234567890ABCDEF1234567890ABCDEF12345678"
	
	workerPoolLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	workerPoolName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccKubernetesAgentWorkerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAgentWorkerBasic(localName, workerPoolLocalName, workerPoolName, name, uri, thumbprint),
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
					"machine_policy_id",
					"space_id",
				},
				ImportStateIdFunc: testAccKubernetesAgentWorkerImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccKubernetesAgentWorkerBasic(localName, workerPoolLocalName, workerPoolName, name, uri, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_static_worker_pool" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_kubernetes_agent_worker" "%s" {
		name            = "%s"
		uri             = "%s"
		thumbprint      = "%s"
		worker_pool_ids = [octopusdeploy_static_worker_pool.%s.id]
		communication_mode = "Polling"
	}`, workerPoolLocalName, workerPoolName, localName, name, uri, thumbprint, workerPoolLocalName)
}

func testAccKubernetesAgentWorkerUpdate(localName, workerPoolLocalName, workerPoolName, workerPool2LocalName, workerPool2Name, name, uri, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_static_worker_pool" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_static_worker_pool" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_kubernetes_agent_worker" "%s" {
		name            = "%s"
		uri             = "%s"
		thumbprint      = "%s"
		worker_pool_ids = [octopusdeploy_static_worker_pool.%s.id, octopusdeploy_static_worker_pool.%s.id]
		is_disabled     = true
		machine_policy_id = "MachinePolicies-1"
		communication_mode = "Polling"
	}`, workerPoolLocalName, workerPoolName, workerPool2LocalName, workerPool2Name, localName, name, uri, thumbprint, workerPoolLocalName, workerPool2LocalName)
}

func testAccKubernetesAgentWorkerWithUpgradeLocked(localName, workerPoolLocalName, workerPoolName, name, uri, thumbprint string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_static_worker_pool" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_kubernetes_agent_worker" "%s" {
		name            = "%s"
		uri             = "%s"
		thumbprint      = "%s"
		worker_pool_ids = [octopusdeploy_static_worker_pool.%s.id]
		upgrade_locked  = true
		communication_mode = "Polling"
	}`, workerPoolLocalName, workerPoolName, localName, name, uri, thumbprint, workerPoolLocalName)
}

func testAccKubernetesAgentWorkerExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		workerID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := workers.GetByID(octoClient, octoClient.GetSpaceID(), workerID); err != nil {
			return err
		}

		return nil
	}
}

func testAccKubernetesAgentWorkerCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_kubernetes_agent_worker" {
			continue
		}

		if worker, err := workers.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("kubernetes agent worker (%s) still exists", worker.GetID())
		}
	}

	return nil
}

func testAccKubernetesAgentWorkerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}