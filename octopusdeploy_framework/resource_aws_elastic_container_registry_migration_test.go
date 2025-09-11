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

func TestAwsElasticContainerRegistryResource_UpgradeFromSDK_ToPluginFramework(t *testing.T) {
	// override the path to check for terraformrc file and test against the real 0.39.0 version
	os.Setenv("TF_CLI_CONFIG_FILE=", "")

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAwsEcrFeedDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"octopusdeploy": {
						VersionConstraint: "0.39.0",
						Source:            "OctopusDeployLabs/octopusdeploy",
					},
				},
				Config: awsEcrFeedMigrationConfig,
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   awsEcrFeedMigrationConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
				Config:                   awsEcrFeedMigrationUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAwsEcrFeedUpdated(t),
				),
			},
		},
	})
}

const awsEcrFeedMigrationConfig = `resource "octopusdeploy_aws_elastic_container_registry" "feed_aws_ecr_migration" {
  name       = "AWS ECR Migration"
  region     = "us-east-1"
  access_key = "accesskey123"
  secret_key = "secretkey123"
}`

const awsEcrFeedMigrationUpdatedConfig = `resource "octopusdeploy_aws_elastic_container_registry" "feed_aws_ecr_migration" {
  name       = "Updated_AWS_ECR_Migration"
  region     = "us-west-2"
  oidc_authentication = {
    session_duration = "3600"
    role_arn = "arn:aws:iam::123456789012:role/example-role"
    audience = "audience"
    subject_keys = ["feed", "space"]
  }
}`

func testAwsEcrFeedDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_aws_elastic_container_registry" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("feed (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAwsEcrFeedUpdated(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedId := s.RootModule().Resources["octopusdeploy_aws_elastic_container_registry"+".feed_aws_ecr_migration"].Primary.ID
		feed, err := octoClient.Feeds.GetByID(feedId)
		if err != nil {
			return fmt.Errorf("Failed to retrieve feed by ID: %s", err)
		}

		ecrFeed := feed.(*feeds.AwsElasticContainerRegistry)

		assert.Regexp(t, "^Feeds\\-\\d+$", ecrFeed.GetID(), "Feed ID did not match expected value")
		assert.Equal(t, "Updated_AWS_ECR_Migration", ecrFeed.Name, "Feed name did not match expected value")
		assert.Equal(t, "us-west-2", ecrFeed.Region, "Feed region did not match expected value")
		assert.NotNil(t, ecrFeed.OidcAuthentication, "OIDCAuthentication should not be nil")
		assert.Equal(t, "3600", ecrFeed.OidcAuthentication.SessionDuration, "OIDC SessionDuration did not match expected value")
		assert.Equal(t, "arn:aws:iam::123456789012:role/example-role", ecrFeed.OidcAuthentication.RoleArn, "OIDC RoleArn did not match expected value")
		assert.Equal(t, "audience", ecrFeed.OidcAuthentication.Audience, "OIDC Audience did not match expected value")
		assert.ElementsMatch(t, []string{"feed", "space"}, ecrFeed.OidcAuthentication.SubjectKeys, "OIDC SubjectKeys did not match expected value")

		return nil
	}
}
