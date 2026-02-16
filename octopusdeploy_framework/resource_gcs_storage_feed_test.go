package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccOctopusDeployGcsStorageFeed(t *testing.T) {
	t.Skip("Skipping until server is updated with GCS storage feed support")

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_gcs_storage_feed." + localName

	downloadAttempts := acctest.RandIntRange(1, 10)
	downloadRetryBackoffSeconds := acctest.RandIntRange(0, 60)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	project := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	serviceAccountKey := acctest.RandStringFromCharSet(300, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testGcsStorageFeedCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testGcsStorageFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "download_attempts", strconv.Itoa(downloadAttempts)),
					resource.TestCheckResourceAttr(prefix, "download_retry_backoff_seconds", strconv.Itoa(downloadRetryBackoffSeconds)),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "project", project),
					resource.TestCheckResourceAttr(prefix, "use_service_account_key", "true"),
				),
				Config: testGcsStorageFeedBasic(localName, downloadAttempts, downloadRetryBackoffSeconds, name, project, serviceAccountKey),
			},
		},
	})
}

func TestAccOctopusDeployGcsStorageFeedWithOIDC(t *testing.T) {
	t.Skip("Skipping until server is updated with GCS storage feed support")

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_gcs_storage_feed." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	project := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	audience := "https://gcs.googleapis.com"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testGcsStorageFeedCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testGcsStorageFeedExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "project", project),
					resource.TestCheckResourceAttr(prefix, "use_service_account_key", "false"),
					resource.TestCheckResourceAttr(prefix, "oidc_authentication.audience", audience),
					resource.TestCheckTypeSetElemAttr(prefix, "oidc_authentication.subject_keys.*", "feed"),
					resource.TestCheckTypeSetElemAttr(prefix, "oidc_authentication.subject_keys.*", "space"),
				),
				Config: testGcsStorageFeedWithOIDC(localName, name, project, audience),
			},
		},
	})
}

func testGcsStorageFeedBasic(localName string, downloadAttempts int, downloadRetryBackoffSeconds int, name string, project string, serviceAccountKey string) string {
	return fmt.Sprintf(`resource "octopusdeploy_gcs_storage_feed" "%s" {
		download_attempts = "%v"
		download_retry_backoff_seconds = "%v"
		name = "%s"
		project = "%s"
		use_service_account_key = true
		service_account_json_key = "%s"
	}`, localName, downloadAttempts, downloadRetryBackoffSeconds, name, project, serviceAccountKey)
}

func testGcsStorageFeedWithOIDC(localName string, name string, project string, audience string) string {
	return fmt.Sprintf(`resource "octopusdeploy_gcs_storage_feed" "%s" {
		name = "%s"
		project = "%s"
		use_service_account_key = false
		oidc_authentication = {
			audience = "%s"
			subject_keys = ["feed", "space"]
		}
	}`, localName, name, project, audience)
}

func testGcsStorageFeedExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		feedID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Feeds.GetByID(feedID); err != nil {
			return err
		}

		return nil
	}
}

func testGcsStorageFeedCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_gcs_storage_feed" {
			continue
		}

		feed, err := octoClient.Feeds.GetByID(rs.Primary.ID)
		if err == nil && feed != nil {
			return fmt.Errorf("GCS storage feed (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
