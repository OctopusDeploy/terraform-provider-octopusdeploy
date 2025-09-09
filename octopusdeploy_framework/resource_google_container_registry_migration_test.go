package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGoogleContainerRegistryResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	// override the path to check for terraformrc file and test against the real 0.39.0 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testGoogleContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.39.0",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: googleContainerRegistryConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   googleContainerRegistryConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   googleContainerRegistryUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testGoogleContainerRegistryUpdated(t),
				),
			},
		},
	})
}

const googleContainerRegistryConfig = `resource "octopusdeploy_google_container_registry" "feed_google_container_registry_migration" {
  name          = "Google Container Registry"
  feed_uri      = "https://gcr.io/test-project"
  username      = "username"
  password      = "password"
  registry_path = "test-registry"
}`

const googleContainerRegistryUpdatedConfig = `resource "octopusdeploy_google_container_registry" "feed_google_container_registry_migration" {
  name          = "Updated_Google_Container_Registry"
  feed_uri      = "https://gcr.io/test-project-updated"
  registry_path = "test-registry-updated"
  oidc_authentication = {
    audience      = "audience"
    subject_keys  = ["feed", "space"]
  }
}`

func testGoogleContainerRegistryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_google_container_registry" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testGoogleContainerRegistryUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedId := s.RootModule().Resources["octopusdeploy_google_container_registry"+".feed_google_container_registry_migration"].Primary.ID
		feed, err := octoClient.Feeds.GetByID(feedId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve feed by ID: %s", err)
		}

		gcrFeed := feed.(*feeds.GoogleContainerRegistry)

		assert.Regexp(t, "^Feeds\\-\\d+$", gcrFeed.GetID(), "Feed ID did not match expected value")
		assert.Equal(t, "Updated_Google_Container_Registry", gcrFeed.Name, "Feed name did not match expected value")
		assert.Equal(t, "https://gcr.io/test-project-updated", gcrFeed.FeedURI, "Feed URI did not match expected value")
		assert.Equal(t, "test-registry-updated", gcrFeed.RegistryPath, "Feed registry path did not match expected value")

		assert.NotNil(t, gcrFeed.OidcAuthentication, "OIDCAuthentication should not be nil")
		assert.Equal(t, "audience", gcrFeed.OidcAuthentication.Audience, "OIDC Audience did not match expected value")
		assert.ElementsMatch(t, []string{"feed", "space"}, gcrFeed.OidcAuthentication.SubjectKeys, "OIDC SubjectKeys did not match expected value")

		return nil
	}
}
