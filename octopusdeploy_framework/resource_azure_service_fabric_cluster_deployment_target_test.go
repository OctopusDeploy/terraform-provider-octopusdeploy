package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployAzureServiceFabricClusterDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_service_fabric_cluster_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	connectionEndpoint := "tcp://example-cluster.eastus.cloudapp.azure.com:19000"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureServiceFabricClusterDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureServiceFabricClusterDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "connection_endpoint", connectionEndpoint),
					resource.TestCheckResourceAttr(prefix, "roles.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.0", "azure-service-fabric"),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccAzureServiceFabricClusterDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, connectionEndpoint),
			},
		},
	})
}

func TestAccOctopusDeployAzureServiceFabricClusterDeploymentTargetUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_service_fabric_cluster_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	connectionEndpoint := "tcp://example-cluster.eastus.cloudapp.azure.com:19000"
	newConnectionEndpoint := "tcp://example-cluster-updated.westus.cloudapp.azure.com:19000"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureServiceFabricClusterDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureServiceFabricClusterDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "connection_endpoint", connectionEndpoint),
				),
				Config: testAccAzureServiceFabricClusterDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, connectionEndpoint),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureServiceFabricClusterDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "connection_endpoint", newConnectionEndpoint),
				),
				Config: testAccAzureServiceFabricClusterDeploymentTargetBasic(localName, environmentLocalName, newName, environmentName, newConnectionEndpoint),
			},
		},
	})
}

func TestAccOctopusDeployAzureServiceFabricClusterDeploymentTargetMinimal(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_azure_service_fabric_cluster_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	connectionEndpoint := "tcp://minimal-cluster.eastus.cloudapp.azure.com:19000"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureServiceFabricClusterDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccAzureServiceFabricClusterDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "connection_endpoint", connectionEndpoint),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccAzureServiceFabricClusterDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, connectionEndpoint),
			},
		},
	})
}

func TestAccOctopusDeployAzureServiceFabricClusterDeploymentTargetImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_azure_service_fabric_cluster_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	connectionEndpoint := "tcp://example-cluster.eastus.cloudapp.azure.com:19000"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccAzureServiceFabricClusterDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureServiceFabricClusterDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, connectionEndpoint),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"has_latest_calamari", "health_status", "status_summary", "aad_user_credential_password"},
				ImportStateIdFunc:       testAccAzureServiceFabricClusterDeploymentTargetImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAzureServiceFabricClusterDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, connectionEndpoint string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "%s" {
			name        = "%s"
			description = "Test environment for Azure Service Fabric cluster deployment target"
		}

		resource "octopusdeploy_azure_service_fabric_cluster_deployment_target" "%s" {
			name               = "%s"
			environments       = [octopusdeploy_environment.%s.id]
			roles              = ["azure-service-fabric"]
			connection_endpoint = "%s"
		}`, environmentLocalName, environmentName, localName, name, environmentLocalName, connectionEndpoint)
}


func testAccAzureServiceFabricClusterDeploymentTargetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		deploymentTargetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), deploymentTargetID); err != nil {
			return err
		}

		return nil
	}
}

func testAccAzureServiceFabricClusterDeploymentTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_azure_service_fabric_cluster_deployment_target" {
			continue
		}

		if deploymentTarget, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("Azure Service Fabric cluster deployment target (%s) still exists", deploymentTarget.GetID())
		}
	}

	return nil
}

func testAccAzureServiceFabricClusterDeploymentTargetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}