package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployProjectScheduledTriggerBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project_scheduled_trigger." + localName

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
	triggerDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectScheduledTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectScheduledTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", triggerName),
					resource.TestCheckResourceAttr(prefix, "description", triggerDescription),
					resource.TestCheckResourceAttr(prefix, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttrSet(prefix, "project_id"),
					// resource.TestCheckResourceAttr(prefix, "is_disabled", "false"),
					resource.TestCheckResourceAttr(prefix, "timezone", "UTC"),
					resource.TestCheckResourceAttr(prefix, "cron_expression_schedule.#", "1"),
					resource.TestCheckResourceAttr(prefix, "cron_expression_schedule.0.cron_expression", "0 0 6 * * MON-FRI"),
					resource.TestCheckResourceAttr(prefix, "deploy_latest_release_action.#", "1"),
					resource.TestCheckResourceAttrSet(prefix, "deploy_latest_release_action.0.source_environment_id"),
					resource.TestCheckResourceAttrSet(prefix, "deploy_latest_release_action.0.destination_environment_id"),
					resource.TestCheckResourceAttr(prefix, "deploy_latest_release_action.0.should_redeploy", "false"),
				),
				Config: testAccProjectScheduledTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription, triggerName, triggerDescription),
			},
		},
	})
}

func TestAccOctopusDeployProjectScheduledTriggerUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project_scheduled_trigger." + localName

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
	triggerDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newTriggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newTriggerDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectScheduledTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectScheduledTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", triggerName),
					resource.TestCheckResourceAttr(prefix, "description", triggerDescription),
					// resource.TestCheckResourceAttr(prefix, "is_disabled", "false"),
				),
				Config: testAccProjectScheduledTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription, triggerName, triggerDescription),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectScheduledTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newTriggerName),
					resource.TestCheckResourceAttr(prefix, "description", newTriggerDescription),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "true"),
					resource.TestCheckResourceAttr(prefix, "deploy_latest_release_action.0.should_redeploy", "true"),
				),
				Config: testAccProjectScheduledTriggerUpdate(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription, newTriggerName, newTriggerDescription),
			},
		},
	})
}

func TestAccOctopusDeployProjectScheduledTriggerOnceDailySchedule(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project_scheduled_trigger." + localName

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
	triggerDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectScheduledTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectScheduledTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", triggerName),
					resource.TestCheckResourceAttr(prefix, "once_daily_schedule.#", "1"),
					resource.TestCheckResourceAttr(prefix, "once_daily_schedule.0.start_time", "08:00:00"),
					resource.TestCheckResourceAttr(prefix, "once_daily_schedule.0.days_of_week.#", "2"),
					resource.TestCheckResourceAttr(prefix, "once_daily_schedule.0.days_of_week.0", "Monday"),
					resource.TestCheckResourceAttr(prefix, "once_daily_schedule.0.days_of_week.1", "Friday"),
				),
				Config: testAccProjectScheduledTriggerOnceDailySchedule(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription, triggerName, triggerDescription),
			},
		},
	})
}

func TestAccOctopusDeployProjectScheduledTriggerImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_project_scheduled_trigger." + localName

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
	triggerDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectScheduledTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectScheduledTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription, triggerName, triggerDescription),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
				ImportStateIdFunc:       testAccProjectScheduledTriggerImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccProjectScheduledTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription, triggerName, triggerDescription string) string {
	return testAccProjectScheduledTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription) + fmt.Sprintf(`
	resource "octopusdeploy_project_scheduled_trigger" "%s" {
		name        = "%s"
		description = "%s"
		space_id    = "Spaces-1"
		project_id  = octopusdeploy_project.%s.id

		deploy_latest_release_action {
			source_environment_id      = octopusdeploy_environment.%s.id
			destination_environment_id = octopusdeploy_environment.%s.id
			should_redeploy           = false
		}

		cron_expression_schedule {
			cron_expression = "0 0 6 * * MON-FRI"
		}

		timezone = "UTC"
	}`, localName, triggerName, triggerDescription, projectLocalName, environmentLocalName, environmentLocalName)
}

func testAccProjectScheduledTriggerUpdate(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription, triggerName, triggerDescription string) string {
	return testAccProjectScheduledTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription) + fmt.Sprintf(`
	resource "octopusdeploy_project_scheduled_trigger" "%s" {
		name        = "%s"
		description = "%s"
		space_id    = "Spaces-1"
		project_id  = octopusdeploy_project.%s.id
		is_disabled = true

		deploy_latest_release_action {
			source_environment_id      = octopusdeploy_environment.%s.id
			destination_environment_id = octopusdeploy_environment.%s.id
			should_redeploy           = true
		}

		cron_expression_schedule {
			cron_expression = "0 0 8 * * MON-FRI"
		}

		timezone = "UTC"
	}`, localName, triggerName, triggerDescription, projectLocalName, environmentLocalName, environmentLocalName)
}

func testAccProjectScheduledTriggerOnceDailySchedule(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription, triggerName, triggerDescription string) string {
	return testAccProjectScheduledTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription) + fmt.Sprintf(`
	resource "octopusdeploy_project_scheduled_trigger" "%s" {
		name        = "%s"
		description = "%s"
		space_id    = "Spaces-1"
		project_id  = octopusdeploy_project.%s.id

		deploy_latest_release_action {
			source_environment_id      = octopusdeploy_environment.%s.id
			destination_environment_id = octopusdeploy_environment.%s.id
			should_redeploy           = false
		}

		once_daily_schedule {
			start_time    = "08:00:00"
			days_of_week  = ["Monday", "Friday"]
		}

		timezone = "America/New_York"
	}`, localName, triggerName, triggerDescription, projectLocalName, environmentLocalName, environmentLocalName)
}

func testAccProjectScheduledTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, environmentDescription string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_environment" "%s" {
		name        = "%s"
		description = "%s"
	}

	resource "octopusdeploy_project" "%s" {
		description      = "%s"
		lifecycle_id     = octopusdeploy_lifecycle.%s.id
		name             = "%s"
		project_group_id = octopusdeploy_project_group.%s.id
	}`, 
		lifecycleLocalName, lifecycleName,
		projectGroupLocalName, projectGroupName,
		environmentLocalName, environmentName, environmentDescription,
		projectLocalName, projectDescription, lifecycleLocalName, projectName, projectGroupLocalName)
}

func testAccProjectScheduledTriggerExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		triggerID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.ProjectTriggers.GetByID(triggerID); err != nil {
			return err
		}

		return nil
	}
}

func testAccProjectScheduledTriggerCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project_scheduled_trigger" {
			continue
		}

		if trigger, err := octoClient.ProjectTriggers.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("project scheduled trigger (%s) still exists", trigger.GetID())
		}
	}

	return nil
}

func testAccProjectScheduledTriggerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}