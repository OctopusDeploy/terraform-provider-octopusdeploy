package octopusdeploy_framework

import (
	"fmt"
	"testing"

	internaltest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProjectVersioningStrategyWithUpdate(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	versioningStrategyPrefix := "octopusdeploy_project_versioning_strategy." + localName
	template := "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.LastPatch}.#{Octopus.Version.NextRevision}"
	template2 := "#{Octopus.Date.Year}.#{Octopus.Date.Month}.#{Octopus.Date.Day}.i"

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(versioningStrategyPrefix, "template", template),
				),
				Config: testProjectVersioningStrategy(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description, template),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(versioningStrategyPrefix, "template", template2),
				),
				Config: testProjectVersioningStrategy(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description, template2),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(versioningStrategyPrefix, "donor_package.deployment_action", "Hello World"),
					resource.TestCheckResourceAttr(versioningStrategyPrefix, "donor_package.package_reference", ""),
				),
				Config: testProjectVersioningStrategyFromPackage(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttrWith(versioningStrategyPrefix, "donor_package_step_id", func(value string) error {
						if value == "" {
							return fmt.Errorf("donor_package_step_id should not be empty")
						}
						return nil
					}),
				),
				Config: testProjectVersioningStrategyFromStep(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description),
			},
		},
	})
}

func testProjectVersioningStrategy(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, localName string, name string, description string, template string) string {
	projectGroup := internaltest.NewProjectGroupTestOptions()
	projectGroup.LocalName = projectGroupLocalName
	projectGroup.Resource.Name = projectGroupName

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		internaltest.ProjectGroupConfiguration(projectGroup)+"\n"+
		`resource "octopusdeploy_project" "%s" {
			description      = "%s"
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			name             = "%s"
			project_group_id = octopusdeploy_project_group.%s.id

			connectivity_policy {
				allow_deployments_to_no_targets = true
				skip_machine_behavior           = "None"
			}
		}

		resource "octopusdeploy_process" "%s" {
		  project_id  = octopusdeploy_project.%s.id
		}

		resource "octopusdeploy_process_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  name = "%s"
		  properties = {
			"Octopus.Action.TargetRoles" = "role-one"
		  }
		  type = "Octopus.TentaclePackage"
		  primary_package = {
			package_id: "my.package"
			feed_id: "built-in"
		  }
		  execution_properties = {
			"Octopus.Action.RunOnServer" = "True"    
		  }
		}

		resource octopusdeploy_process_steps_order "%s" {
		  process_id = octopusdeploy_process.%s.id
		  steps = [
			octopusdeploy_process_step.%s.id
		  ]
		}

		resource "octopusdeploy_project_versioning_strategy" "%s" {
		  project_id = "${octopusdeploy_project.%s.id}"
		  template   = "%s"
		}`,
		localName,
		description,
		lifecycleLocalName,
		name,
		projectGroupLocalName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		template)
}

func testProjectVersioningStrategyFromPackage(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, localName string, name string, description string) string {
	projectGroup := internaltest.NewProjectGroupTestOptions()
	projectGroup.LocalName = projectGroupLocalName
	projectGroup.Resource.Name = projectGroupName

	configuration := fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		internaltest.ProjectGroupConfiguration(projectGroup)+"\n"+
		`resource "octopusdeploy_project" "%s" {
			description      = "%s"
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			name             = "%s"
			project_group_id = octopusdeploy_project_group.%s.id

			connectivity_policy {
				allow_deployments_to_no_targets = true
				skip_machine_behavior           = "None"
			}
		}

		resource "octopusdeploy_process" "%s" {
		  project_id  = octopusdeploy_project.%s.id
		}

		resource "octopusdeploy_process_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  name = "Hello World"
		  properties = {
			"Octopus.Action.TargetRoles" = "role-one"
		  }
		  type = "Octopus.TentaclePackage"
		  primary_package = {
			package_id: "Package"
			feed_id: "built-in"
		  }
		  execution_properties = {
			"Octopus.Action.RunOnServer" = "True"    
		  }
		}

		resource octopusdeploy_process_steps_order "%s" {
		  process_id = octopusdeploy_process.%s.id
		  steps = [
			octopusdeploy_process_step.%s.id
		  ]
		}

		resource "octopusdeploy_project_versioning_strategy" "%s" {
		  project_id = "${octopusdeploy_project.%s.id}"
		  donor_package = {
			deployment_action = "Hello World"
			package_reference = ""
		  }
          depends_on = [
			octopusdeploy_process_step.%s,
			octopusdeploy_process.%s
		  ]
		}`,
		localName,
		description,
		lifecycleLocalName,
		name,
		projectGroupLocalName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName)

	return configuration
}

func testProjectVersioningStrategyFromStep(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, localName string, name string, description string) string {
	projectGroup := internaltest.NewProjectGroupTestOptions()
	projectGroup.LocalName = projectGroupLocalName
	projectGroup.Resource.Name = projectGroupName

	configuration := fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		internaltest.ProjectGroupConfiguration(projectGroup)+"\n"+
		`resource "octopusdeploy_project" "%s" {
			description      = "%s"
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			name             = "%s"
			project_group_id = octopusdeploy_project_group.%s.id

			connectivity_policy {
				allow_deployments_to_no_targets = true
				skip_machine_behavior           = "None"
			}
		}

		resource "octopusdeploy_process" "%s" {
		  project_id  = octopusdeploy_project.%s.id
		}

		resource "octopusdeploy_process_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  name = "Hello World"
		  properties = {
			"Octopus.Action.TargetRoles" = "role-one"
		  }
		  type = "Octopus.TentaclePackage"
		  primary_package = {
			package_id: "Package"
			feed_id: "built-in"
		  }
		  execution_properties = {
			"Octopus.Action.RunOnServer" = "True"    
		  }
		}

		resource octopusdeploy_process_steps_order "%s" {
		  process_id = octopusdeploy_process.%s.id
		  steps = [
			octopusdeploy_process_step.%s.id
		  ]
		}

		resource "octopusdeploy_process_child_step" "%s" {
		  process_id  = octopusdeploy_process.%s.id
		  parent_id = octopusdeploy_process_step.%s.id
		  name = "Hello World Child"
		  type = "Octopus.TentaclePackage"
		  primary_package = {
			package_id: "Package"
			feed_id: "built-in"
		  }
		  execution_properties = {
			"Octopus.Action.RunOnServer" = "True"    
		  }
		}

		resource "octopusdeploy_process_child_steps_order" "%s" {
		  process_id = octopusdeploy_process.%s.id
		  parent_id = octopusdeploy_process_step.%s.id
		  children = [
			octopusdeploy_process_child_step.%s.id,
		  ]
		}

		resource "octopusdeploy_project_versioning_strategy" "%s" {
		  project_id = "${octopusdeploy_project.%s.id}"
		  donor_package_step_id = octopusdeploy_process_child_step.%s.id
          depends_on = [
			octopusdeploy_process.%s
		  ]
		}`,
		localName,
		description,
		lifecycleLocalName,
		name,
		projectGroupLocalName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName,
		localName)

	return configuration
}
