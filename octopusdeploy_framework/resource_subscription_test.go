package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSubscriptionWebhook(t *testing.T) {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedName := name + "-updated"
	resourceName := "octopusdeploy_subscription." + name

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testAccSubscriptionCheckDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccSubscriptionWebhookConfig(name, name, "https://example.com/webhook", false),
				Check: resource.ComposeTestCheckFunc(
					testAccSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "is_disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.webhook_uri", "https://example.com/webhook"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.0", "Modified"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
				),
			},
			{
				Config: testAccSubscriptionWebhookConfig(name, updatedName, "https://example.com/webhook-v2", true),
				Check: resource.ComposeTestCheckFunc(
					testAccSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "is_disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.webhook_uri", "https://example.com/webhook-v2"),
				),
			},
			{
				Config: testAccSubscriptionWebhookNoFilterConfig(name, updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "is_disabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSubscriptionEmailAndWebhook(t *testing.T) {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_subscription." + name

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testAccSubscriptionCheckDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccSubscriptionEmailAndWebhookConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.webhook_uri", "https://example.com/webhook"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.email_frequency_period", "00:30:00"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.email_priority", "High"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.webhook_timeout", "00:00:30"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.document_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.document_types.0", "Machines"),
				),
			},
		},
	})
}

func testAccSubscriptionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		spaceID := rs.Primary.Attributes["space_id"]
		_, err := newclient.GetByID[subscriptionApiModel](octoClient, subscriptionURITemplate, spaceID, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error reading subscription %s: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccSubscriptionCheckDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return nil
		}

		spaceID := rs.Primary.Attributes["space_id"]
		_, err := newclient.GetByID[subscriptionApiModel](octoClient, subscriptionURITemplate, spaceID, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("subscription %s still exists", rs.Primary.ID)
		}
		return nil
	}
}

func testAccSubscriptionWebhookConfig(resourceID, name, webhookURI string, disabled bool) string {
	return fmt.Sprintf(`
resource "octopusdeploy_subscription" "%s" {
  name        = "%s"
  is_disabled = %t

  event_notification_subscription = {
    webhook_uri     = "%s"
    webhook_timeout = "00:00:10"

    filter = {
      event_categories = ["Modified"]
    }
  }
}`, resourceID, name, disabled, webhookURI)
}

func testAccSubscriptionWebhookNoFilterConfig(resourceID, name string) string {
	return fmt.Sprintf(`
resource "octopusdeploy_subscription" "%s" {
  name        = "%s"
  is_disabled = false

  event_notification_subscription = {
    webhook_uri     = "https://example.com/webhook-v2"
    webhook_timeout = "00:00:10"
  }
}`, resourceID, name)
}

func testAccSubscriptionEmailAndWebhookConfig(resourceID string) string {
	return fmt.Sprintf(`
resource "octopusdeploy_subscription" "%s" {
  name = "%s"

  event_notification_subscription = {
    webhook_uri            = "https://example.com/webhook"
    webhook_timeout        = "00:00:30"
    email_frequency_period = "00:30:00"
    email_priority         = "High"

    filter = {
      event_categories = ["Created", "Modified"]
      document_types   = ["Machines"]
    }
  }
}`, resourceID, resourceID)
}
