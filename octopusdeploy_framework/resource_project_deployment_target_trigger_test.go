package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployProjectDeploymentTargetTriggerBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project_deployment_target_trigger." + localName

	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	triggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectDeploymentTargetTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectDeploymentTargetTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", triggerName),
					resource.TestCheckResourceAttrSet(prefix, "project_id"),
					resource.TestCheckResourceAttr(prefix, "should_redeploy", "false"),
					resource.TestCheckResourceAttr(prefix, "event_groups.#", "1"),
					resource.TestCheckResourceAttr(prefix, "event_groups.0", "Machine"),
					resource.TestCheckResourceAttr(prefix, "event_categories.#", "0"),
					resource.TestCheckResourceAttr(prefix, "roles.#", "0"),
					resource.TestCheckResourceAttr(prefix, "environment_ids.#", "0"),
				),
				Config: testAccProjectDeploymentTargetTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, triggerName),
			},
		},
	})
}

func TestAccOctopusDeployProjectDeploymentTargetTriggerUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project_deployment_target_trigger." + localName

	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	triggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newTriggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectDeploymentTargetTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectDeploymentTargetTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", triggerName),
					resource.TestCheckResourceAttr(prefix, "should_redeploy", "false"),
				),
				Config: testAccProjectDeploymentTargetTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, triggerName),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectDeploymentTargetTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newTriggerName),
					resource.TestCheckResourceAttr(prefix, "should_redeploy", "true"),
				),
				Config: testAccProjectDeploymentTargetTriggerUpdate(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, newTriggerName),
			},
		},
	})
}

func TestAccOctopusDeployProjectDeploymentTargetTriggerWithFilters(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project_deployment_target_trigger." + localName

	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	triggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectDeploymentTargetTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectDeploymentTargetTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", triggerName),
					resource.TestCheckResourceAttr(prefix, "event_groups.#", "2"),
					resource.TestCheckResourceAttr(prefix, "event_groups.0", "Machine"),
					resource.TestCheckResourceAttr(prefix, "event_groups.1", "MachineHealthChanged"),
					resource.TestCheckResourceAttr(prefix, "event_categories.#", "2"),
					resource.TestCheckResourceAttr(prefix, "event_categories.0", "MachineAdded"),
					resource.TestCheckResourceAttr(prefix, "event_categories.1", "MachineHealthy"),
					resource.TestCheckResourceAttr(prefix, "roles.#", "1"),
					resource.TestCheckResourceAttr(prefix, "roles.0", "web-server"),
					resource.TestCheckResourceAttr(prefix, "environment_ids.#", "1"),
				),
				Config: testAccProjectDeploymentTargetTriggerWithFilters(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription, triggerName),
			},
		},
	})
}

func TestAccOctopusDeployProjectDeploymentTargetTriggerImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_project_deployment_target_trigger." + localName

	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	triggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectDeploymentTargetTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDeploymentTargetTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, triggerName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
				ImportStateIdFunc:       testAccProjectDeploymentTargetTriggerImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccProjectDeploymentTargetTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, triggerName string) string {
	return testAccProjectDeploymentTargetTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription) + fmt.Sprintf(`
	resource "octopusdeploy_project_deployment_target_trigger" "%s" {
		name       = "%s"
		project_id = octopusdeploy_project.%s.id
		event_groups = ["Machine"]
	}`, localName, triggerName, projectLocalName)
}

func testAccProjectDeploymentTargetTriggerUpdate(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, triggerName string) string {
	return testAccProjectDeploymentTargetTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription) + fmt.Sprintf(`
	resource "octopusdeploy_project_deployment_target_trigger" "%s" {
		name           = "%s"
		project_id     = octopusdeploy_project.%s.id
		should_redeploy = true
		event_groups = ["Machine"]
	}`, localName, triggerName, projectLocalName)
}

func testAccProjectDeploymentTargetTriggerWithFilters(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription, triggerName string) string {
	return testAccProjectDeploymentTargetTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription) + fmt.Sprintf(`
	resource "octopusdeploy_environment" "%s" {
		name        = "%s"
		description = "%s"
	}

	resource "octopusdeploy_project_deployment_target_trigger" "%s" {
		name           = "%s"
		project_id     = octopusdeploy_project.%s.id
		should_redeploy = false
		
		event_groups = ["Machine", "MachineHealthChanged"]
		event_categories = ["MachineAdded", "MachineHealthy"]
		roles = ["web-server"]
		environment_ids = [octopusdeploy_environment.%s.id]
	}`, environmentLocalName, environmentName, environmentDescription, localName, triggerName, projectLocalName, environmentLocalName)
}

func testAccProjectDeploymentTargetTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_project" "%s" {
		description      = "%s"
		lifecycle_id     = octopusdeploy_lifecycle.%s.id
		name             = "%s"
		project_group_id = octopusdeploy_project_group.%s.id
	}`, 
		lifecycleLocalName, lifecycleName,
		projectGroupLocalName, projectGroupName,
		projectLocalName, projectDescription, lifecycleLocalName, projectName, projectGroupLocalName)
}

func testAccProjectDeploymentTargetTriggerExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		triggerID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.ProjectTriggers.GetByID(triggerID); err != nil {
			return err
		}

		return nil
	}
}

func testAccProjectDeploymentTargetTriggerCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project_deployment_target_trigger" {
			continue
		}

		if trigger, err := octoClient.ProjectTriggers.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("project deployment target trigger (%s) still exists", trigger.GetID())
		}
	}

	return nil
}

func testAccProjectDeploymentTargetTriggerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}