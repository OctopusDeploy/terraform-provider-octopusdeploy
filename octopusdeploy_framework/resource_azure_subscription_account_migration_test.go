package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAzureSubscriptionAccountResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAzureSubscriptionAccountDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "1.2.1",
						Source:            "OctopusDeploy/octopusdeploy",
					},
				},
				Config: azureSubscriptionAccountConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   azureSubscriptionAccountConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   updatedAzureSubscriptionAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					testAzureSubscriptionAccountUpdated(t),
				),
			},
		},
	})
}

const azureSubscriptionAccountConfig = `
	resource "octopusdeploy_environment" "test_env1" {
	  name = "Test Environment 1"
	  description = "Test environment for azure account migration"
	}

	resource "octopusdeploy_environment" "test_env2" {
	  name = "Test Environment 2"
	  description = "Test environment for azure account migration"
	}

	resource "octopusdeploy_tag_set" "test_tagset" {
	  name = "test-tagset"
	  description = "Test tagset"
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

	resource "octopusdeploy_tenant" "test_tenant1" {
	  name = "Test Tenant 1"
	  description = "Test tenant for azure account migration"
	}

	resource "octopusdeploy_tenant" "test_tenant2" {
	  name = "Test Tenant 2"
	  description = "Test tenant for azure account migration"
	}

	resource "octopusdeploy_azure_subscription_account" "azure_account" {
	  name                               = "Azure Subscription Account"
	  description                        = "Test Azure subscription account for migration"
	  subscription_id                    = "00000000-0000-0000-0000-000000000000"
	  azure_environment                  = "AzureCloud"
	  management_endpoint                = "https://management.azure.com/"
	  storage_endpoint_suffix            = "azure.com"
	  tenanted_deployment_participation  = "TenantedOrUntenanted"
	  environments                       = [octopusdeploy_environment.test_env1.id, octopusdeploy_environment.test_env2.id]
	  tenants                            = [octopusdeploy_tenant.test_tenant1.id]
	  tenant_tags                        = ["test-tagset/tag1", "test-tagset/tag2"]
	  depends_on                         = [octopusdeploy_tag.test_tag1, octopusdeploy_tag.test_tag2]
	}`

const updatedAzureSubscriptionAccountConfig = `
	resource "octopusdeploy_environment" "test_env1" {
	  name = "Test Environment 1"
	  description = "Test environment for azure account migration"
	}

	resource "octopusdeploy_environment" "test_env2" {
	  name = "Test Environment 2"
	  description = "Test environment for azure account migration"
	}

	resource "octopusdeploy_tag_set" "test_tagset" {
	  name = "test-tagset"
	  description = "Test tagset"
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

	resource "octopusdeploy_tenant" "test_tenant1" {
	  name = "Test Tenant 1"
	  description = "Test tenant for azure account migration"
	}

	resource "octopusdeploy_tenant" "test_tenant2" {
	  name = "Test Tenant 2"
	  description = "Test tenant for azure account migration"
	}

	resource "octopusdeploy_azure_subscription_account" "azure_account" {
	  name                               = "Updated Azure Subscription Account"
	  description                        = "Updated test Azure subscription account"
	  subscription_id                    = "00000000-0000-0000-0000-000000000000"
	  azure_environment                  = "AzureCloud"
	  management_endpoint                = "https://management.core.chinacloudapi.cn/"
	  storage_endpoint_suffix            = "core.chinacloudapi.cn"
	  tenanted_deployment_participation  = "TenantedOrUntenanted"
	  environments                       = [octopusdeploy_environment.test_env1.id, octopusdeploy_environment.test_env2.id]
	  tenants                            = [octopusdeploy_tenant.test_tenant1.id]
	  tenant_tags                        = ["test-tagset/tag1", "test-tagset/tag2"]
	  depends_on                         = [octopusdeploy_tag.test_tag1, octopusdeploy_tag.test_tag2]
	}`

func testAzureSubscriptionAccountDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_azure_subscription_account" {
			account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
			if err == nil && account != nil {
				return fmt.Errorf("account (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAzureSubscriptionAccountUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		accountId := s.RootModule().Resources["octopusdeploy_azure_subscription_account.azure_account"].Primary.ID
		account, err := accounts.GetByID(octoClient, octoClient.GetSpaceID(), accountId)
		if err != nil {
			return fmt.Errorf("failed to retrieve account by ID: %s", err)
		}

		azureAccount := account.(*accounts.AzureSubscriptionAccount)

		env1ID := s.RootModule().Resources["octopusdeploy_environment.test_env1"].Primary.ID
		env2ID := s.RootModule().Resources["octopusdeploy_environment.test_env2"].Primary.ID
		tenant1ID := s.RootModule().Resources["octopusdeploy_tenant.test_tenant1"].Primary.ID

		assert.NotEmpty(t, azureAccount.GetID(), "Account ID should not be empty")
		assert.Equal(t, "Updated Azure Subscription Account", azureAccount.Name, "Account name did not match expected value")
		assert.Equal(t, "Updated test Azure subscription account", azureAccount.GetDescription(), "Account description did not match expected value")
		assert.Equal(t, "AzureCloud", azureAccount.AzureEnvironment, "Azure environment did not match expected value")
		assert.Equal(t, "https://management.core.chinacloudapi.cn/", azureAccount.ManagementEndpoint, "Management endpoint did not match expected value")
		assert.Equal(t, "core.chinacloudapi.cn", azureAccount.StorageEndpointSuffix, "Storage endpoint suffix did not match expected value")
		assert.Equal(t, "TenantedOrUntenanted", string(azureAccount.GetTenantedDeploymentMode()), "Tenanted deployment participation did not match expected value")

		expectedEnvironmentIDs := []string{env1ID, env2ID}
		assert.ElementsMatch(t, expectedEnvironmentIDs, azureAccount.GetEnvironmentIDs(), "Environment IDs should match expected values")

		expectedTenantIDs := []string{tenant1ID}
		assert.ElementsMatch(t, expectedTenantIDs, azureAccount.GetTenantIDs(), "Tenant IDs should match expected values")

		expectedTenantTags := []string{"test-tagset/tag1", "test-tagset/tag2"}
		assert.ElementsMatch(t, expectedTenantTags, azureAccount.TenantTags, "Tenant tags should match expected values")

		return nil
	}
}
