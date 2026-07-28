package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceWebhookTrigger(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := fmt.Sprintf("octopusdeploy_webhook_trigger.%s", localName)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testWebhookTriggerCheckDestroy(s) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configTestAccWebhookTrigger(localName, "description", `secret = "secret"`, false),
				Check: resource.ComposeTestCheckFunc(
					testAssertWebhookTriggerAttributes(prefix),
					resource.TestCheckResourceAttr(prefix, "description", "description"),
					resource.TestCheckResourceAttr(prefix, "secret", "secret"),
					resource.TestCheckResourceAttr(prefix, "require_api_key", "false"),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "false"),
				),
			},
			{
				// Changing description
				Config: configTestAccWebhookTrigger(localName, "updated description", `secret = "secret"`, true),
				Check: resource.ComposeTestCheckFunc(
					testAssertWebhookTriggerAttributes(prefix),
					resource.TestCheckResourceAttr(prefix, "description", "updated description"),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "true"),
				),
			},
			{
				// Rotating secret
				Config: configTestAccWebhookTrigger(localName, "updated description", `secret = "rotated-secret"`, true),
				Check: resource.ComposeTestCheckFunc(
					testAssertWebhookTriggerAttributes(prefix),
					resource.TestCheckResourceAttr(prefix, "secret", "rotated-secret"),
				),
			},
			{
				// Switching from a secret to API key authentication
				Config: configTestAccWebhookTrigger(localName, "updated description", `require_api_key = true`, true),
				Check: resource.ComposeTestCheckFunc(
					testAssertWebhookTriggerAttributes(prefix),
					resource.TestCheckResourceAttr(prefix, "require_api_key", "true"),
					resource.TestCheckNoResourceAttr(prefix, "secret"),
				),
			},
			{
				// Switching back requires a new secret, the previous one was discarded
				Config: configTestAccWebhookTrigger(localName, "updated description", `secret = "new-secret"`, true),
				Check: resource.ComposeTestCheckFunc(
					testAssertWebhookTriggerAttributes(prefix),
					resource.TestCheckResourceAttr(prefix, "secret", "new-secret"),
					resource.TestCheckResourceAttr(prefix, "require_api_key", "false"),
				),
			},
		},
	})
}

func TestAccResourceWebhookTriggerAuthenticationValidation(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testWebhookTriggerCheckDestroy(s) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configTestAccWebhookTrigger(localName, "description", `
			  secret          = "secret"
			  require_api_key = true`, false),
				ExpectError: regexp.MustCompile("cannot have both"),
			},
			{
				Config:      configTestAccWebhookTrigger(localName, "description", "", false),
				ExpectError: regexp.MustCompile("requires either"),
			},
		},
	})
}

func testAssertWebhookTriggerAttributes(prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(prefix, "id"),
		resource.TestCheckResourceAttrSet(prefix, "name"),
		resource.TestCheckResourceAttrSet(prefix, "project_id"),
		resource.TestCheckResourceAttrSet(prefix, "space_id"),
		resource.TestCheckResourceAttrSet(prefix, "run_runbook_action.runbook_id"),
		resource.TestCheckResourceAttr(prefix, "run_runbook_action.target_environment_ids.#", "1"),
		resource.TestCheckResourceAttr(prefix, "tenant_ids.#", "1"),

		// This is sent back from the server
		resource.TestCheckResourceAttrSet(prefix, "webhook_id"),
		testWebhookTriggerExists(prefix),
	)
}

func testWebhookTriggerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		triggerResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		trigger, err := triggers.GetById(octoClient, octoClient.GetSpaceID(), triggerResource.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to read webhook trigger (%s): %s", triggerResource.Primary.ID, err.Error())
		}

		if trigger == nil {
			return fmt.Errorf("webhook trigger (%s) does not exist", triggerResource.Primary.ID)
		}

		return nil
	}
}

func testWebhookTriggerCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_webhook_trigger" {
			continue
		}

		if trigger, err := triggers.GetById(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil && trigger != nil {
			return fmt.Errorf("webhook trigger (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func configTestAccWebhookTrigger(localName string, description string, auth string, isDisabled bool) string {
	return fmt.Sprintf(`
		data "octopusdeploy_lifecycles" "default_%[1]s" {
		  ids          = null
		  partial_name = "Default Lifecycle"
		  skip         = 0
		  take         = 1
		}

		resource "octopusdeploy_project_group" "%[1]s" {
		  name        = "webhook-trigger-group-%[1]s"
		  description = "Test project group"
		}

		resource "octopusdeploy_project" "%[1]s" {
		  name                  = "webhook-trigger-project-%[1]s"
		  lifecycle_id          = data.octopusdeploy_lifecycles.default_%[1]s.lifecycles[0].id
		  project_group_id      = octopusdeploy_project_group.%[1]s.id
		  description           = "Project with a webhook trigger"
		  is_disabled           = false
		  is_version_controlled = false
		}

		resource "octopusdeploy_environment" "%[1]s" {
		  name = "webhook-trigger-env-%[1]s"
		}

		resource "octopusdeploy_runbook" "%[1]s" {
		  project_id  = octopusdeploy_project.%[1]s.id
		  name        = "webhook-trigger-runbook-%[1]s"
		  description = "Runbook"
		}

		resource "octopusdeploy_tenant" "%[1]s" {
		  name        = "webhook-trigger-tenant-%[1]s"
		  description = "Tenant"
		}

		resource "octopusdeploy_webhook_trigger" "%[1]s" {
		  name            = "webhook-trigger-%[1]s"
		  description     = "%[2]s"
		  project_id      = octopusdeploy_project.%[1]s.id
		  tenant_ids      = [octopusdeploy_tenant.%[1]s.id]
		  is_disabled     = %[4]t
		  %[3]s

		  run_runbook_action = {
		    runbook_id             = octopusdeploy_runbook.%[1]s.id
		    target_environment_ids = [octopusdeploy_environment.%[1]s.id]
		  }
		}`, localName, description, auth, isDisabled)
}
