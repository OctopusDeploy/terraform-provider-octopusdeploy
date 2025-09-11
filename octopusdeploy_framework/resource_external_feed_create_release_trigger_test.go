package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployExternalFeedCreateReleaseTriggerBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_external_feed_create_release_trigger." + localName

	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	triggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccExternalFeedCreateReleaseTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccExternalFeedCreateReleaseTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", triggerName),
					resource.TestCheckResourceAttrSet(prefix, "project_id"),
					resource.TestCheckResourceAttrSet(prefix, "channel_id"),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "false"),
					resource.TestCheckResourceAttr(prefix, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttr(prefix, "package.#", "0"),
				),
				Config: testAccExternalFeedCreateReleaseTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription, triggerName),
			},
		},
	})
}

func TestAccOctopusDeployExternalFeedCreateReleaseTriggerUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_external_feed_create_release_trigger." + localName

	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	triggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newTriggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccExternalFeedCreateReleaseTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccExternalFeedCreateReleaseTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", triggerName),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "false"),
				),
				Config: testAccExternalFeedCreateReleaseTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription, triggerName),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccExternalFeedCreateReleaseTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newTriggerName),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "true"),
				),
				Config: testAccExternalFeedCreateReleaseTriggerUpdate(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription, newTriggerName),
			},
		},
	})
}

func testAccOctopusDeployExternalFeedCreateReleaseTriggerWithPrimaryPackage(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_external_feed_create_release_trigger." + localName

	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	triggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccExternalFeedCreateReleaseTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccExternalFeedCreateReleaseTriggerExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", triggerName),
					resource.TestCheckResourceAttr(prefix, "primary_package.#", "1"),
					resource.TestCheckResourceAttr(prefix, "primary_package.0.deployment_action_slug", "test-action"),
					resource.TestCheckResourceAttr(prefix, "package.#", "0"),
				),
				Config: testAccExternalFeedCreateReleaseTriggerWithPrimaryPackage(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription, triggerName),
			},
		},
	})
}

func TestAccOctopusDeployExternalFeedCreateReleaseTriggerImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_external_feed_create_release_trigger." + localName

	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	triggerName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccExternalFeedCreateReleaseTriggerCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccExternalFeedCreateReleaseTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription, triggerName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccExternalFeedCreateReleaseTriggerImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccExternalFeedCreateReleaseTriggerBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription, triggerName string) string {
	return testAccExternalFeedCreateReleaseTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription) + fmt.Sprintf(`
	resource "octopusdeploy_external_feed_create_release_trigger" "%s" {
		name       = "%s"
		space_id   = "Spaces-1"
		project_id = octopusdeploy_project.%s.id
		channel_id = octopusdeploy_channel.%s.id
	}`, localName, triggerName, projectLocalName, channelLocalName)
}

func testAccExternalFeedCreateReleaseTriggerUpdate(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription, triggerName string) string {
	return testAccExternalFeedCreateReleaseTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription) + fmt.Sprintf(`
	resource "octopusdeploy_external_feed_create_release_trigger" "%s" {
		name        = "%s"
		space_id    = "Spaces-1"
		project_id  = octopusdeploy_project.%s.id
		channel_id  = octopusdeploy_channel.%s.id
		is_disabled = true
	}`, localName, triggerName, projectLocalName, channelLocalName)
}

func testAccExternalFeedCreateReleaseTriggerWithPrimaryPackage(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription, triggerName string) string {
	return testAccExternalFeedCreateReleaseTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription) + fmt.Sprintf(`
	resource "octopusdeploy_external_feed_create_release_trigger" "%s" {
		name       = "%s"
		space_id   = "Spaces-1"
		project_id = octopusdeploy_project.%s.id
		channel_id = octopusdeploy_channel.%s.id
		
		primary_package {
			deployment_action_slug = "test-action"
		}
	}`, localName, triggerName, projectLocalName, channelLocalName)
}

func testAccExternalFeedCreateReleaseTriggerDependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, channelLocalName, channelName, channelDescription string) string {
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
	}


	resource "octopusdeploy_channel" "%s" {
		description = "%s"
		name        = "%s"
		project_id  = octopusdeploy_project.%s.id
	}`, 
		lifecycleLocalName, lifecycleName,
		projectGroupLocalName, projectGroupName,
		projectLocalName, projectDescription, lifecycleLocalName, projectName, projectGroupLocalName,
		channelLocalName, channelDescription, channelName, projectLocalName)
}

func testAccExternalFeedCreateReleaseTriggerExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		triggerID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.ProjectTriggers.GetByID(triggerID); err != nil {
			return err
		}

		return nil
	}
}

func testAccExternalFeedCreateReleaseTriggerCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_external_feed_create_release_trigger" {
			continue
		}

		if trigger, err := octoClient.ProjectTriggers.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("external feed create release trigger (%s) still exists", trigger.GetID())
		}
	}

	return nil
}

func testAccExternalFeedCreateReleaseTriggerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}