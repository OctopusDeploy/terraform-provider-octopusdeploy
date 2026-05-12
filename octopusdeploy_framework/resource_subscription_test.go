package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/subscriptions"
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
		_, err := subscriptions.GetByID(octoClient, spaceID, rs.Primary.ID)
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
		_, err := subscriptions.GetByID(octoClient, spaceID, rs.Primary.ID)
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

// TestAccSubscriptionFilters exercises every filter attribute that scopes a
// subscription to specific events. It creates the supporting resources whose
// IDs are referenced by ID-based filters (users, projects, project_groups,
// environments, tenants, tags), and uses well-known string values pulled from
// /api/events/{categories,groups,agents} and Octopus document type names for
// the remaining filters.
func TestAccSubscriptionFilters(t *testing.T) {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_subscription." + name

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testAccSubscriptionCheckDestroy(resourceName),
		Steps: []resource.TestStep{
			// Step 1: every filter field populated.
			{
				Config: testAccSubscriptionAllFiltersConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					// ID-based filters: one item each, sourced from a referenced resource.
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.users.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "event_notification_subscription.filter.users.0"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "event_notification_subscription.filter.projects.0", "octopusdeploy_project."+name, "id"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.project_groups.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "event_notification_subscription.filter.project_groups.0", "octopusdeploy_project_group."+name, "id"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.environments.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "event_notification_subscription.filter.environments.0", "octopusdeploy_environment."+name, "id"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.tenants.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "event_notification_subscription.filter.tenants.0", "octopusdeploy_tenant."+name, "id"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.tags.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "event_notification_subscription.filter.tags.0", "octopusdeploy_tag."+name, "canonical_tag_name"),
					// String-valued filters: well-known values from the API.
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_groups.0", "Deployment"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.0", "Created"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.1", "Modified"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_agents.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_agents.0", "Server"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.document_types.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.document_types.0", "Machines"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.document_types.1", "Projects"),
				),
			},
			// Step 2: modify string-valued filter values to verify updates flow through.
			{
				Config: testAccSubscriptionAllFiltersConfigUpdated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.0", "Created"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.1", "Modified"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.2", "Deleted"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.document_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.document_types.0", "Deployments"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_agents.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_groups.#", "1"),
					// ID-based filter membership is unchanged from step 1.
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.projects.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.tenants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.tags.#", "1"),
				),
			},
			// Step 3: clear individual filter fields by setting them to []. The kept
			// fields confirm partial clearing works without disturbing other lists.
			{
				Config: testAccSubscriptionAllFiltersConfigCleared(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.users.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.project_groups.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.environments.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.tenants.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_agents.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_groups.#", "0"),
					// These remain populated.
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.projects.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.event_categories.0", "Modified"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.document_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "event_notification_subscription.filter.document_types.0", "Deployments"),
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

// testAccSubscriptionFiltersDependencies builds the supporting resources used
// to populate ID-based filter fields: a project group, environment, lifecycle,
// project (tenanted), tag set + tag, tenant linking project/environment, and a
// users data source resolving the local admin.
func testAccSubscriptionFiltersDependencies(name string) string {
	return fmt.Sprintf(`
data "octopusdeploy_lifecycles" "%[1]s" {
  partial_name = "Default Lifecycle"
  skip         = 0
  take         = 1
}

data "octopusdeploy_users" "%[1]s" {
  filter = "admin"
}

resource "octopusdeploy_project_group" "%[1]s" {
  name = "pg-%[1]s"
}

resource "octopusdeploy_environment" "%[1]s" {
  name = "env-%[1]s"
}

resource "octopusdeploy_project" "%[1]s" {
  name                              = "proj-%[1]s"
  lifecycle_id                      = data.octopusdeploy_lifecycles.%[1]s.lifecycles[0].id
  project_group_id                  = octopusdeploy_project_group.%[1]s.id
  tenanted_deployment_participation = "TenantedOrUntenanted"
}

resource "octopusdeploy_tag_set" "%[1]s" {
  name = "ts-%[1]s"
}

resource "octopusdeploy_tag" "%[1]s" {
  name       = "tag-%[1]s"
  tag_set_id = octopusdeploy_tag_set.%[1]s.id
  color      = "#00FF00"
}

resource "octopusdeploy_tenant" "%[1]s" {
  name = "tnt-%[1]s"

  depends_on = [octopusdeploy_tag.%[1]s]
}

resource "octopusdeploy_tenant_project" "%[1]s" {
  tenant_id       = octopusdeploy_tenant.%[1]s.id
  project_id      = octopusdeploy_project.%[1]s.id
  environment_ids = [octopusdeploy_environment.%[1]s.id]
}
`, name)
}

func testAccSubscriptionAllFiltersConfig(name string) string {
	return testAccSubscriptionFiltersDependencies(name) + fmt.Sprintf(`
resource "octopusdeploy_subscription" "%[1]s" {
  name = "%[1]s"

  event_notification_subscription = {
    webhook_uri = "https://example.com/webhook"

    filter = {
      users            = [data.octopusdeploy_users.%[1]s.users[0].id]
      projects         = [octopusdeploy_project.%[1]s.id]
      project_groups   = [octopusdeploy_project_group.%[1]s.id]
      environments     = [octopusdeploy_environment.%[1]s.id]
      tenants          = [octopusdeploy_tenant.%[1]s.id]
      tags             = [octopusdeploy_tag.%[1]s.canonical_tag_name]
      event_groups     = ["Deployment"]
      event_categories = ["Created", "Modified"]
      event_agents     = ["Server"]
      document_types   = ["Machines", "Projects"]
    }
  }

  depends_on = [octopusdeploy_tenant_project.%[1]s]
}
`, name)
}

func testAccSubscriptionAllFiltersConfigUpdated(name string) string {
	return testAccSubscriptionFiltersDependencies(name) + fmt.Sprintf(`
resource "octopusdeploy_subscription" "%[1]s" {
  name = "%[1]s"

  event_notification_subscription = {
    webhook_uri = "https://example.com/webhook"

    filter = {
      users            = [data.octopusdeploy_users.%[1]s.users[0].id]
      projects         = [octopusdeploy_project.%[1]s.id]
      project_groups   = [octopusdeploy_project_group.%[1]s.id]
      environments     = [octopusdeploy_environment.%[1]s.id]
      tenants          = [octopusdeploy_tenant.%[1]s.id]
      tags             = [octopusdeploy_tag.%[1]s.canonical_tag_name]
      event_groups     = ["Deployment"]
      event_categories = ["Created", "Modified", "Deleted"]
      event_agents     = ["Server", "Portal"]
      document_types   = ["Deployments"]
    }
  }

  depends_on = [octopusdeploy_tenant_project.%[1]s]
}
`, name)
}

func testAccSubscriptionAllFiltersConfigCleared(name string) string {
	return testAccSubscriptionFiltersDependencies(name) + fmt.Sprintf(`
resource "octopusdeploy_subscription" "%[1]s" {
  name = "%[1]s"

  event_notification_subscription = {
    webhook_uri = "https://example.com/webhook"

    filter = {
      users            = []
      projects         = [octopusdeploy_project.%[1]s.id]
      project_groups   = []
      environments     = []
      tenants          = []
      tags             = []
      event_groups     = []
      event_categories = ["Modified"]
      event_agents     = []
      document_types   = ["Deployments"]
    }
  }
}
`, name)
}
