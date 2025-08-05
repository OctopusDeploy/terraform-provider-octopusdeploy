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

func TestAzureContainerRegistryResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	// override the path to check for terraformrc file and test against the real 0.39.0 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAzureContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.39.0",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: azureContainerRegistryConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   azureContainerRegistryConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   azureContainerRegistryUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAzureContainerRegistryUpdated(t),
				),
			},
		},
	})
}

const azureContainerRegistryConfig = `resource "octopusdeploy_azure_container_registry" "feed_azure_container_registry_migration" {
						  name                                 = "Azure Container Registry"
						  feed_uri                             = "https://test-azure.azurecr.io"
						  username                             = "username"
						  password                             = "password"
						  registry_path                        = "test-registry"
					   }`

const azureContainerRegistryUpdatedConfig = `resource "octopusdeploy_azure_container_registry" "feed_azure_container_registry_migration" {
						  name                                 = "Updated_Azure_Container_Registry"
						  feed_uri                             = "https://test-azure-updated.azurecr.io"
						  registry_path                        = "test-registry-updated"
						  oidc_authentication = {
						  client_id     = "00000000-0000-0000-0000-000000000000"
						  tenant_id     = "00000000-0000-0000-0000-000000000000"
						  audience      = "audience"
						  subject_keys = ["feed", "space"]
						  } 
					   }`

func testAzureContainerRegistryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_azure_container_registry" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAzureContainerRegistryUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedId := s.RootModule().Resources["octopusdeploy_azure_container_registry"+".feed_azure_container_registry_migration"].Primary.ID
		feed, err := octoClient.Feeds.GetByID(feedId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve feed by ID: %s", err)
		}

		azureContainerRegistryFeed := feed.(*feeds.AzureContainerRegistry)

		assert.Regexp(t, "^Feeds\\-\\d+$", azureContainerRegistryFeed.GetID(), "Feed ID did not match expected value")
		assert.Equal(t, "Updated_Azure_Container_Registry", azureContainerRegistryFeed.Name, "Feed name did not match expected value")
		assert.Equal(t, "https://test-azure-updated.azurecr.io", azureContainerRegistryFeed.FeedURI, "Feed URI did not match expected value")
		assert.Equal(t, "test-registry-updated", azureContainerRegistryFeed.RegistryPath, "Feed registry path did not match expected value")

		assert.NotNil(t, azureContainerRegistryFeed.OidcAuthentication, "OIDCAuthentication should not be nil")
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", azureContainerRegistryFeed.OidcAuthentication.ClientId, "OIDC ClientID did not match expected value")
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", azureContainerRegistryFeed.OidcAuthentication.TenantId, "OIDC TenantID did not match expected value")
		assert.Equal(t, "audience", azureContainerRegistryFeed.OidcAuthentication.Audience, "OIDC Audience did not match expected value")
		assert.ElementsMatch(t, []string{"feed", "space"}, azureContainerRegistryFeed.OidcAuthentication.SubjectKeys, "OIDC SubjectKeys did not match expected value")

		return nil
	}
}
