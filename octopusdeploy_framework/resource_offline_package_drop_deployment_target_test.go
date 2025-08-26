package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployOfflinePackageDropDeploymentTargetBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_offline_package_drop_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	applicationsDirectory := "/tmp/applications"
	workingDirectory := "/tmp/working"
	dropFolderPath := "/tmp/packages"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccOfflinePackageDropDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccOfflinePackageDropDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "applications_directory", applicationsDirectory),
					resource.TestCheckResourceAttr(prefix, "working_directory", workingDirectory),
					resource.TestCheckResourceAttr(prefix, "destination.#", "1"),
					resource.TestCheckResourceAttr(prefix, "destination.0.destination_type", "FileSystem"),
					resource.TestCheckResourceAttr(prefix, "destination.0.drop_folder_path", dropFolderPath),
					resource.TestCheckResourceAttr(prefix, "roles.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.0", "offline"),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccOfflinePackageDropDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, applicationsDirectory, workingDirectory, dropFolderPath),
			},
		},
	})
}

func TestAccOctopusDeployOfflinePackageDropDeploymentTargetUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_offline_package_drop_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	applicationsDirectory := "/tmp/applications"
	workingDirectory := "/tmp/working"
	dropFolderPath := "/tmp/packages"
	newDropFolderPath := "/tmp/packages-updated"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccOfflinePackageDropDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccOfflinePackageDropDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "destination.0.drop_folder_path", dropFolderPath),
				),
				Config: testAccOfflinePackageDropDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, applicationsDirectory, workingDirectory, dropFolderPath),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccOfflinePackageDropDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "destination.0.drop_folder_path", newDropFolderPath),
				),
				Config: testAccOfflinePackageDropDeploymentTargetBasic(localName, environmentLocalName, newName, environmentName, applicationsDirectory, workingDirectory, newDropFolderPath),
			},
		},
	})
}

func TestAccOctopusDeployOfflinePackageDropDeploymentTargetMinimal(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_offline_package_drop_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccOfflinePackageDropDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccOfflinePackageDropDeploymentTargetExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttrSet(prefix, "id"),
				),
				Config: testAccOfflinePackageDropDeploymentTargetMinimal(localName, environmentLocalName, name, environmentName),
			},
		},
	})
}

func TestAccOctopusDeployOfflinePackageDropDeploymentTargetImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_offline_package_drop_deployment_target." + localName

	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	applicationsDirectory := "/tmp/applications"
	workingDirectory := "/tmp/working"
	dropFolderPath := "/tmp/packages"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccOfflinePackageDropDeploymentTargetCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOfflinePackageDropDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, applicationsDirectory, workingDirectory, dropFolderPath),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"has_latest_calamari", "health_status", "status_summary"},
				ImportStateIdFunc:       testAccOfflinePackageDropDeploymentTargetImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccOfflinePackageDropDeploymentTargetBasic(localName, environmentLocalName, name, environmentName, applicationsDirectory, workingDirectory, dropFolderPath string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "%s" {
			name        = "%s"
			description = "Test environment for offline package drop deployment target"
		}

		resource "octopusdeploy_offline_package_drop_deployment_target" "%s" {
			name                   = "%s"
			environments           = [octopusdeploy_environment.%s.id]
			roles                  = ["offline"]
			applications_directory = "%s"
			working_directory      = "%s"
			
			destination {
				destination_type  = "FileSystem"
				drop_folder_path  = "%s"
			}
		}`, environmentLocalName, environmentName, localName, name, environmentLocalName, applicationsDirectory, workingDirectory, dropFolderPath)
}

func testAccOfflinePackageDropDeploymentTargetMinimal(localName, environmentLocalName, name, environmentName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_environment" "%s" {
			name        = "%s"
			description = "Test environment for offline package drop deployment target"
		}

		resource "octopusdeploy_offline_package_drop_deployment_target" "%s" {
			name                   = "%s"
			environments           = [octopusdeploy_environment.%s.id]
			roles                  = ["offline"]
			applications_directory = "/tmp/apps"
			working_directory      = "/tmp/work"
		}`, environmentLocalName, environmentName, localName, name, environmentLocalName)
}

func testAccOfflinePackageDropDeploymentTargetExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		deploymentTargetID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), deploymentTargetID); err != nil {
			return err
		}

		return nil
	}
}

func testAccOfflinePackageDropDeploymentTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_offline_package_drop_deployment_target" {
			continue
		}

		if deploymentTarget, err := machines.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("offline package drop deployment target (%s) still exists", deploymentTarget.GetID())
		}
	}

	return nil
}

func testAccOfflinePackageDropDeploymentTargetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}