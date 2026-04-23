package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccProjectWebhookTrigger_basic(t *testing.T) {
	projectWebhookTriggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectWebhookTriggerResource := "octopusdeploy_project_webhook_trigger." + projectWebhookTriggerName
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	runbookName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testWebhookTriggerDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{ // 1 create project webhook trigger
				Config: projectWebhookTrigger(projectWebhookTriggerName, "Description1", projectName, environmentName, runbookName),
				Check: resource.ComposeTestCheckFunc(
					testWebhookTriggerExists(projectWebhookTriggerResource),
					resource.TestCheckResourceAttr(projectWebhookTriggerResource, "name", projectWebhookTriggerName),
					resource.TestCheckResourceAttr(projectWebhookTriggerResource, "description", "Description1"),
					resource.TestCheckResourceAttr(projectWebhookTriggerResource, "webhook_trigger_filter.0.secret", "Password01!"),
				),
			},
			//2 update project webhook trigger description
			{
				Config: projectWebhookTrigger(projectWebhookTriggerName, "Description2", projectName, environmentName, runbookName),
				Check: resource.ComposeTestCheckFunc(
					testWebhookTriggerExists(projectWebhookTriggerResource),
					resource.TestCheckResourceAttr(projectWebhookTriggerResource, "name", projectWebhookTriggerName),
					resource.TestCheckResourceAttr(projectWebhookTriggerResource, "description", "Description2"),
					resource.TestCheckResourceAttr(projectWebhookTriggerResource, "webhook_trigger_filter.0.secret", "Password01!"),
				),
			},
		},
	})
}

func testWebhookTriggerDestroy(s *terraform.State) error {
	for _, resourceState := range s.RootModule().Resources {
		if resourceState.Type != "octopusdeploy_project_webhook_trigger" {
			continue
		}
		spaceId := resourceState.Primary.Attributes["space_id"]
		triggerId := resourceState.Primary.ID

		trigger, err := triggers.GetById(octoClient, spaceId, triggerId)
		if err == nil && trigger != nil {
			return fmt.Errorf("trigger (%s) still exists", triggerId)
		}
	}

	return nil
}

func testWebhookTriggerExists(projectWebhookTriggerName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, resourceState := range s.RootModule().Resources {
			if resourceState.Type != "octopusdeploy_project_webhook_trigger" {

				continue
			}
			if resourceState.Primary.ID != projectWebhookTriggerName {
				continue
			}
			spaceId := resourceState.Primary.Attributes["space_id"]
			triggerId := resourceState.Primary.ID
			trigger, err := triggers.GetById(octoClient, spaceId, triggerId)
			if err != nil {
				return err
			}
			if trigger == nil {
				return fmt.Errorf("project webhook trigger (%s) not found", projectWebhookTriggerName)
			}
			return nil
		}
		return nil
	}
}

func projectWebhookTrigger(projectWebhookTriggerName string, projectWebhookTriggerDescription string, projectName string, environmentName string, runbookName string) string {
	setupAboveTrigger := setupForProjectWebhookTrigger(projectName, environmentName, runbookName)
	return fmt.Sprintf(`
	%s

	resource "octopusdeploy_project_webhook_trigger" "%s" {
  		name        = "%s"
  		description = "%s"
  		project_id  = octopusdeploy_project.%s.id
        space_id    = "Spaces-1"
  		run_runbook_action {
    		target_environment_ids = [octopusdeploy_environment.%s.id]
    		runbook_id             = octopusdeploy_runbook.%s.id
  		}
  		webhook_trigger_filter {
    		secret = "Password01!"
  		}
	}
`, setupAboveTrigger, projectWebhookTriggerName, projectWebhookTriggerName, projectWebhookTriggerDescription, projectName, environmentName, runbookName)
}

func setupForProjectWebhookTrigger(projectName string, environmentName string, runbookName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name             = "%s"
		lifecycle_id     = "Lifecycles-1"
		project_group_id = "ProjectGroups-1"
	}

	resource "octopusdeploy_environment" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		name       = "%s"
		project_id = octopusdeploy_project.%s.id
	}
	`, projectName, projectName,
		environmentName, environmentName,
		runbookName, runbookName, projectName,
	)
}
