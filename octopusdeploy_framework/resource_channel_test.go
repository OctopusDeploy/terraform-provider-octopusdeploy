package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

// TestAccChannelWithMRPRule covers the "Most Recently Published" ordering
// strategy and the version-tag regex on a channel package version rule. It also
// asserts that version_range, tag, and version_tag_regex coexist with an MRP
// strategy — versioning_strategy only changes ordering, not which filters apply.
// These fields require the `non-semver-ordering` feature toggle to be enabled on
// the target Octopus instance — without it the server silently ignores MRP rules
// and the round-trip assertions below will fail.
func TestAccChannelWithMRPRule(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	channelName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	feedLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	feedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := fmt.Sprintf("octopusdeploy_channel.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testChannelDestroy,
		Steps: []resource.TestStep{
			{
				// Create an MRP rule with a catch-all regex plus the SemVer
				// range/tag filters, which apply alongside the regex.
				Config: testChannelWithMRPRule(localName, channelName, feedLocalName, feedName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, ".*"),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", channelName),
					resource.TestCheckResourceAttr(resourceName, "rule.0.versioning_strategy", "MostRecentlyPublished"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.version_tag_regex", ".*"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.version_range", "[1.0.0,)"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.tag", "^$"),
				),
			},
			{
				// Update the regex in place; the strategy stays MRP.
				Config: testChannelWithMRPRule(localName, channelName, feedLocalName, feedName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, "^[0-9]+"),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.0.versioning_strategy", "MostRecentlyPublished"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.version_tag_regex", "^[0-9]+"),
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

// testChannelWithMRPRule builds a channel whose single package version rule uses
// the MostRecentlyPublished strategy and a version tag regex. The rule references
// a deployment step's package via an action_package block, which the server's
// channel validator requires whenever a version constraint is set on the rule.
func testChannelWithMRPRule(localName, channelName, feedLocalName, feedName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, versionTagRegex string) string {
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

		resource "octopusdeploy_nuget_feed" "%s" {
			name             = "%s"
			feed_uri         = "https://api.nuget.org/v3/index.json"
			is_enhanced_mode = true
		}

		resource "octopusdeploy_process" "process" {
			project_id = octopusdeploy_project.%s.id
		}

		resource "octopusdeploy_process_step" "step" {
			process_id = octopusdeploy_process.process.id
			name       = "Reference NuGet package"
			type       = "Octopus.Script"

			packages = {
				"nuget-pkg" = {
					package_id           = "Newtonsoft.Json"
					feed_id              = octopusdeploy_nuget_feed.%s.id
					acquisition_location = "Server"
				}
			}

			execution_properties = {
				"Octopus.Action.RunOnServer"         = "True"
				"Octopus.Action.Script.ScriptSource" = "Inline"
				"Octopus.Action.Script.Syntax"       = "Bash"
				"Octopus.Action.Script.ScriptBody"   = "echo 'mrp test step'"
			}
		}

		resource "octopusdeploy_channel" "%s" {
			name       = "%s"
			project_id = octopusdeploy_project.%s.id

			rule {
				versioning_strategy = "MostRecentlyPublished"
				version_tag_regex   = "%s"
				version_range       = "[1.0.0,)"
				tag                 = "^$"
				action_package {
					deployment_action = octopusdeploy_process_step.step.name
					package_reference = "nuget-pkg"
				}
			}
		}
	`,
		lifecycleLocalName, lifecycleName,
		projectGroupLocalName, projectGroupName,
		projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName,
		feedLocalName, feedName,
		projectLocalName,
		feedLocalName,
		localName, channelName, projectLocalName, versionTagRegex,
	)
}
