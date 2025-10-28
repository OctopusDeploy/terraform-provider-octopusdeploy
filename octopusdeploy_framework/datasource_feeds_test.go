package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceFeeds(t *testing.T) {
	localName := acctest.RandStringFromCharSet(50, acctest.CharSetAlpha)
	prefix := fmt.Sprintf("data.octopusdeploy_feeds.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: createTestAccDataSourceFeedsConfig(),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeedsDataSourceID(prefix),
					resource.TestCheckResourceAttrSet(prefix, "feeds.#"),
					resource.TestCheckResourceAttrSet(prefix, "feeds.0.feed_type"),
				),
				Config: testAccDataSourceFeedsConfig(localName, take),
			},
			{
				Check: func(s *terraform.State) error {
					nugetPrefix := fmt.Sprintf("data.octopusdeploy_feeds.%s_nuget", localName)
					return resource.ComposeTestCheckFunc(
						testAccCheckFeedsDataSourceID(nugetPrefix),
						testAccCheckFeedsFilteredByType(nugetPrefix, "NuGet"),
					)(s)
				},
				Config: testAccDataSourceFeedsWithFilterConfig(localName+"_nuget", "NuGet"),
			},
			{
				Check: func(s *terraform.State) error {
					mavenPrefix := fmt.Sprintf("data.octopusdeploy_feeds.%s_maven", localName)
					return resource.ComposeTestCheckFunc(
						testAccCheckFeedsDataSourceID(mavenPrefix),
						testAccCheckFeedsFilteredByType(mavenPrefix, "Maven"),
					)(s)
				},
				Config: testAccDataSourceFeedsWithFilterConfig(localName+"_maven", "Maven"),
			},
			{
				Check: func(s *terraform.State) error {
					helmPrefix := fmt.Sprintf("data.octopusdeploy_feeds.%s_helm", localName)
					return resource.ComposeTestCheckFunc(
						testAccCheckFeedsDataSourceID(helmPrefix),
						testAccCheckFeedsFilteredByType(helmPrefix, "Helm"),
					)(s)
				},
				Config: testAccDataSourceFeedsWithFilterConfig(localName+"_helm", "Helm"),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeedsDataSourceID(prefix),
				),
				Config: testAccDataSourceFeedsEmpty(localName),
			},
		},
	})
}

func testAccCheckFeedsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find Feeds data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot Feeds source ID not set")
		}
		return nil
	}
}

func testAccDataSourceFeedsConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_feeds" "%s" {
		take = %v
	}`, localName, take)
}

func testAccDataSourceFeedsEmpty(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_feeds" "%s" {}`, localName)
}

func testAccCheckFeedsFilteredByType(n string, expectedType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cannot find Feeds data source: %s", n)
		}

		count := rs.Primary.Attributes["feeds.#"]
		if count == "" {
			return fmt.Errorf("feeds count not set")
		}

		if count == "0" {
			return nil
		}

		feedCount := 0
		for key, value := range rs.Primary.Attributes {
			if key == fmt.Sprintf("feeds.%d.feed_type", feedCount) {
				if value != expectedType {
					return fmt.Errorf("feed %d has unexpected type: expected %s, got %s", feedCount, expectedType, value)
				}
				feedCount++
			}
		}

		return nil
	}
}

func testAccDataSourceFeedsWithFilterConfig(localName string, feedType string) string {
	return fmt.Sprintf(`%s

data "octopusdeploy_feeds" "%s" {
	feed_type = "%s"
}`, createTestAccDataSourceFeedsConfig(), localName, feedType)
}

func createTestAccDataSourceFeedsConfig() string {
	return `resource "octopusdeploy_nuget_feed" "nuget_feed" {
    feed_uri                       = "https://api.nuget.org/v3/index.json"
    is_enhanced_mode               = true
    password                       = "test-password"
    name                           = "Test NuGet Feed"
    username                       = "test-username"
  }

  resource "octopusdeploy_maven_feed" "maven_feed" {
    download_attempts              = 10
    download_retry_backoff_seconds = 20
    feed_uri                       = "https://repo.maven.apache.org/maven2/"
    password                       = "test-password"
    name                           = "Test Maven Feed"
    username                       = "test-username"
  }

  resource "octopusdeploy_github_repository_feed" "ghr_feed" {
    download_attempts              = 1
    download_retry_backoff_seconds = 30
    feed_uri                       = "https://api.github.com"
    password                       = "test-password"
    name                           = "Test GitHub Repository Feed"
    username                       = "test-username"
  }

  resource "octopusdeploy_helm_feed" "helm_feed" {
    feed_uri = "https://charts.helm.sh/stable"
    password = "test-password"
    name     = "Test Helm Feed"
    username = "test-username"
  }

  resource "octopusdeploy_artifactory_generic_feed" "artifactory_generic_feed" {
    feed_uri                       = "https://example.jfrog.io"
    password                       = "test-password"
    name                           = "Test Artifactory Generic Feed"
    username                       = "test-username"
    repository                     = "repo"
    layout_regex                   = "this is regex"
  }`
}
