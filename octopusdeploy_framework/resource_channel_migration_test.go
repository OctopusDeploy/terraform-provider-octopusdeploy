package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestChannelResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testChannelDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "1.3.1",
						Source:            "OctopusDeploy/octopusdeploy",
					},
				},
				Config: channelConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   channelConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedChannelConfig,
				Check: resource.ComposeTestCheckFunc(
					testChannelUpdated(t),
				),
			},
		},
	})
}

const channelConfig = `
	resource "octopusdeploy_lifecycle" "test_lifecycle" {
	  name = "Test Lifecycle"
	  description = "Test lifecycle for channel migration"
	}

	resource "octopusdeploy_lifecycle" "custom_lifecycle" {
	  name = "Custom Lifecycle"
	  description = "Custom lifecycle for channel"
	}

	resource "octopusdeploy_project_group" "test_project_group" {
	  name = "Test Project Group"
	  description = "Test project group for channel migration"
	}

	resource "octopusdeploy_tag_set" "test_tagset" {
	  name = "test-tagset"
	  description = "Test tagset for channel"
	}

	resource "octopusdeploy_tag" "test_tag1" {
	  name = "tag1"
	  color = "#ff0000"
	  tag_set_id = octopusdeploy_tag_set.test_tagset.id
	}

	resource "octopusdeploy_tag" "test_tag2" {
	  name = "tag2"
	  color = "#00ff00"
	  tag_set_id = octopusdeploy_tag_set.test_tagset.id
	}

	resource "octopusdeploy_project" "test_project" {
	  lifecycle_id     = octopusdeploy_lifecycle.test_lifecycle.id
	  name             = "Test Project"
	  description      = "Test project for channel migration"
	  project_group_id = octopusdeploy_project_group.test_project_group.id
	}

	resource "octopusdeploy_channel" "test_channel" {
	  name         = "Test Channel"
	  description  = "Test channel for migration"
	  project_id   = octopusdeploy_project.test_project.id
	  lifecycle_id = octopusdeploy_lifecycle.custom_lifecycle.id
	  is_default   = false
	  tenant_tags  = [octopusdeploy_tag.test_tag1.canonical_tag_name, octopusdeploy_tag.test_tag2.canonical_tag_name]
	}`

const updatedChannelConfig = `
	resource "octopusdeploy_lifecycle" "test_lifecycle" {
	  name = "Test Lifecycle"
	  description = "Test lifecycle for channel migration"
	}

	resource "octopusdeploy_lifecycle" "custom_lifecycle" {
	  name = "Custom Lifecycle"
	  description = "Custom lifecycle for channel"
	}

	resource "octopusdeploy_project_group" "test_project_group" {
	  name = "Test Project Group"
	  description = "Test project group for channel migration"
	}

	resource "octopusdeploy_tag_set" "test_tagset" {
	  name = "test-tagset"
	  description = "Test tagset for channel"
	}

	resource "octopusdeploy_tag" "test_tag1" {
	  name = "tag1"
	  color = "#ff0000"
	  tag_set_id = octopusdeploy_tag_set.test_tagset.id
	}

	resource "octopusdeploy_tag" "test_tag2" {
	  name = "tag2"
	  color = "#00ff00"
	  tag_set_id = octopusdeploy_tag_set.test_tagset.id
	}

	resource "octopusdeploy_tag" "test_tag3" {
	  name = "tag3"
	  color = "#0000ff"
	  tag_set_id = octopusdeploy_tag_set.test_tagset.id
	}

	resource "octopusdeploy_project" "test_project" {
	  lifecycle_id     = octopusdeploy_lifecycle.test_lifecycle.id
	  name             = "Test Project"
	  description      = "Test project for channel migration"
	  project_group_id = octopusdeploy_project_group.test_project_group.id
	}

	resource "octopusdeploy_channel" "test_channel" {
	  name         = "Updated Test Channel"
	  description  = "Updated test channel for migration"
	  project_id   = octopusdeploy_project.test_project.id
	  lifecycle_id = octopusdeploy_lifecycle.custom_lifecycle.id
	  is_default   = false
	  tenant_tags  = [octopusdeploy_tag.test_tag1.canonical_tag_name, octopusdeploy_tag.test_tag3.canonical_tag_name]
	}`

func testChannelUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		channelId := s.RootModule().Resources["octopusdeploy_channel.test_channel"].Primary.ID
		channel, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), channelId)
		if err != nil {
			return fmt.Errorf("failed to retrieve channel by ID: %s", err)
		}

		projectID := s.RootModule().Resources["octopusdeploy_project.test_project"].Primary.ID
		lifecycleID := s.RootModule().Resources["octopusdeploy_lifecycle.custom_lifecycle"].Primary.ID

		assert.NotEmpty(t, channel.GetID(), "Channel ID should not be empty")
		assert.Equal(t, "Updated Test Channel", channel.Name, "Channel name did not match expected value")
		assert.Equal(t, "Updated test channel for migration", channel.Description, "Channel description did not match expected value")
		assert.Equal(t, projectID, channel.ProjectID, "Project ID did not match expected value")
		assert.Equal(t, lifecycleID, channel.LifecycleID, "Lifecycle ID did not match expected value")
		assert.True(t, channel.IsDefault, "Channel should be default")

		expectedTenantTags := []string{"test-tagset/tag1", "test-tagset/tag3"}
		assert.ElementsMatch(t, expectedTenantTags, channel.TenantTags, "Tenant tags should match expected values")

		return nil
	}
}
