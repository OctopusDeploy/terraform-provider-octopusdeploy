package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployCloudRegionDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_cloud_region_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccCloudRegionDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCloudRegionDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "default_worker_pool_id", "WorkerPools-1"),
					resource.TestCheckResourceAttr(prefix, "roles.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.0", "cloud-region"),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccCloudRegionDeploymentTargetBasic(localName, environmentLocalName, name, environmentName),
			},
		},
	})
}

func TestAccOctopusDeployCloudRegionDeploymentTargetUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_cloud_region_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccCloudRegionDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCloudRegionDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testAccCloudRegionDeploymentTargetBasic(localName, environmentLocalName, name, environmentName),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCloudRegionDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
				),
				Config: testAccCloudRegionDeploymentTargetBasic(localName, environmentLocalName, newName, environmentName),
			},
		},
	})
}

func TestAccOctopusDeployCloudRegionDeploymentTargetMinimal(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_cloud_region_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccCloudRegionDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCloudRegionDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccCloudRegionDeploymentTargetMinimal(localName, environmentLocalName, name, environmentName),
			},
		},
	})
}

func TestAccOctopusDeployCloudRegionDeploymentTargetImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_cloud_region_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccCloudRegionDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegionDeploymentTargetBasic(localName, environmentLocalName, name, environmentName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudRegionDeploymentTargetImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCloudRegionDeploymentTargetBasic(localName, environmentLocalName, name, environmentName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "%s" {
			name        = "%s"
			description = "Test environment for cloud region deployment target"
		}

		resource "octopusdeploy_cloud_region_deployment_target" "%s" {
			name                   = "%s"
			environments           = [octopusdeploy_environment.%s.id]
			roles                  = ["cloud-region"]
			default_worker_pool_id = "WorkerPools-1"
		}`, environmentLocalName, environmentName, localName, name, environmentLocalName)
}

func testAccCloudRegionDeploymentTargetMinimal(localName, environmentLocalName, name, environmentName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "%s" {
			name        = "%s"
			description = "Test environment for cloud region deployment target"
		}

		resource "octopusdeploy_cloud_region_deployment_target" "%s" {
			name         = "%s"
			environments = [octopusdeploy_environment.%s.id]
			roles        = ["cloud-region"]
		}`, environmentLocalName, environmentName, localName, name, environmentLocalName)
}

func testAccCloudRegionDeploymentTargetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		deploymentTargetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), deploymentTargetID); err != nil {
			return err
		}

		return nil
	}
}

func testAccCloudRegionDeploymentTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_cloud_region_deployment_target" {
			continue
		}

		if deploymentTarget, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("cloud region deployment target (%s) still exists", deploymentTarget.GetID())
		}
	}

	return nil
}

func testAccCloudRegionDeploymentTargetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}