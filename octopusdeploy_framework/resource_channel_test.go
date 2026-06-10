package octopusdeploy_framework

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccChannelBasic(t *testing.T) {
	options := test.NewChannelTestOptions()

	// Create test dependencies
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	channel := channels.Channel{
		Name:        acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		Description: acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	}

	resourceName := fmt.Sprintf("octopusdeploy_channel.%s", options.LocalName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testChannelBasic(options.LocalName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channel),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", channel.Name),
					resource.TestCheckResourceAttr(resourceName, "description", channel.Description),
				),
			},
		},
	})
}

func TestBuildChannelUpdateRequest_ClearsEmptyCollections(t *testing.T) {
	channel := channels.NewChannel("channel", "Projects-1")
	channel.CustomFieldDefinitions = []channels.ChannelCustomFieldDefinition{{
		FieldName:   "field",
		Description: "description",
	}}
	channel.Rules = []channels.ChannelRule{{
		ID:           "ChannelRules-1",
		Tag:          "beta",
		VersionRange: "[1.0.0,2.0.0)",
	}}
	channel.TenantTags = []string{"TagSets-1/tag-a"}

	plan := schemas.ChannelModel{
		CustomFieldDefinitions: types.ListNull(types.ObjectType{AttrTypes: getChannelCustomFieldDefinitionAttrTypes()}),
		Rule:                   types.ListNull(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()}),
		TenantTags:             types.SetNull(types.StringType),
	}

	updateReq := buildChannelUpdateRequest(channel, plan)
	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	var payload map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &payload))

	assert.JSONEq(t, `[]`, string(payload["CustomFieldDefinitions"]))
	assert.JSONEq(t, `[]`, string(payload["Rules"]))
	assert.JSONEq(t, `[]`, string(payload["TenantTags"]))
}

func TestAccChannelRuleRemoval(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := fmt.Sprintf("octopusdeploy_channel.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testChannelWithRule(localName, true),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testChannelRuleCount(resourceName, 1),
				),
			},
			{
				Config: testChannelWithRule(localName, false),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "0"),
					testChannelRuleCount(resourceName, 0),
				),
			},
		},
	})
}

func testChannelExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		channelID, err := getChannelID(s, resourceName)
		if err != nil {
			return err
		}

		if _, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), channelID); err != nil {
			return fmt.Errorf("channel %s not found", channelID)
		}

		return nil
	}
}

func testChannelDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_channel" {
			continue
		}

		if _, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("channel %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testChannelRuleCount(resourceName string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		channelID, err := getChannelID(s, resourceName)
		if err != nil {
			return err
		}

		channel, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), channelID)
		if err != nil {
			return fmt.Errorf("channel %s not found", channelID)
		}

		if len(channel.Rules) != expected {
			return fmt.Errorf("expected %d channel rules, got %d", expected, len(channel.Rules))
		}

		return nil
	}
}

func getChannelID(s *terraform.State, resourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("Not found: %s", resourceName)
	}

	return rs.Primary.ID, nil
}

func testChannelBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName string, channel channels.Channel) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_lifecycle" "%s" {
			name = "%s"
		}

		resource "octopusdeploy_project_group" "%s" {
			name = "%s"
		}

		resource "octopusdeploy_project" "%s" {
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			name             = "%s"
			project_group_id = octopusdeploy_project_group.%s.id
		}

		resource "octopusdeploy_channel" "%s" {
			name        = "%s"
			description = "%s"
			project_id  = octopusdeploy_project.%s.id
		}
	`,
		lifecycleLocalName, lifecycleName,
		projectGroupLocalName, projectGroupName,
		projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName,
		localName, channel.Name, channel.Description, projectLocalName,
	)
}

func testChannelWithRule(localName string, includeRule bool) string {
	ruleBlock := ""
	if includeRule {
		ruleBlock = `
			rule {
				version_range = "[1.0.0,2.0.0)"

				action_package {
					deployment_action = octopusdeploy_process_step.package_step.name
				}
			}
		`
	}

	return fmt.Sprintf(`
		data "octopusdeploy_lifecycles" "default" {
		  ids          = null
		  partial_name = "Default Lifecycle"
		  skip         = 0
		  take         = 1
		}

		data "octopusdeploy_feeds" "built_in_feed" {
		  feed_type    = "BuiltIn"
		  ids          = null
		  partial_name = ""
		  skip         = 0
		  take         = 1
		}

		resource "octopusdeploy_project_group" "%[1]s" {
		  name = "%[1]s"
		}

		resource "octopusdeploy_project" "%[1]s" {
		  lifecycle_id     = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		  name             = "%[1]s"
		  project_group_id = octopusdeploy_project_group.%[1]s.id
		}

		resource "octopusdeploy_process" "%[1]s" {
		  project_id = octopusdeploy_project.%[1]s.id
		}

		resource "octopusdeploy_process_step" "package_step" {
		  process_id = octopusdeploy_process.%[1]s.id
		  name       = "Package deployment"
		  type       = "Octopus.TentaclePackage"
		  properties = {
			"Octopus.Action.TargetRoles" = "Webserver"
		  }
		  primary_package = {
			feed_id    = data.octopusdeploy_feeds.built_in_feed.feeds[0].id
			package_id = "MyPackage"
		  }
		  execution_properties = {
			"Octopus.Action.RunOnServer" = "True"
		  }
		}

		resource "octopusdeploy_channel" "%[1]s" {
		  name       = "%[1]s"
		  project_id = octopusdeploy_project.%[1]s.id

		  %[2]s

		  depends_on = [octopusdeploy_process_step.package_step]
		}
	`, localName, ruleBlock)
}
