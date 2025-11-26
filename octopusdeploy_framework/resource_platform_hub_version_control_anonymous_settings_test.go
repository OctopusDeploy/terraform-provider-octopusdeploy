package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlatformHubVersionControlAnonymousSettingsBasic(t *testing.T) {
	resourceName := "octopusdeploy_platform_hub_version_control_anonymous_settings.test"

	url := "https://github.com/acme/hello-world.git"
	defaultBranch := "main"
	basePath := ".octopus"

	updatedURL := "https://github.com/acme/hello-world.git"
	updatedBranch := "develop"
	updatedBasePath := ".octopus/updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubVersionControlAnonymousSettings(url, defaultBranch, basePath),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubVersionControlSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "default_branch", defaultBranch),
					resource.TestCheckResourceAttr(resourceName, "base_path", basePath),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testPlatformHubVersionControlAnonymousSettings(updatedURL, updatedBranch, updatedBasePath),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubVersionControlSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "url", updatedURL),
					resource.TestCheckResourceAttr(resourceName, "default_branch", updatedBranch),
					resource.TestCheckResourceAttr(resourceName, "base_path", updatedBasePath),
				),
			},
		},
	})
}

func testPlatformHubVersionControlAnonymousSettings(url, defaultBranch, basePath string) string {
	return fmt.Sprintf(`
resource "octopusdeploy_platform_hub_version_control_anonymous_settings" "test" {
	url            = "%s"
	default_branch = "%s"
	base_path      = "%s"
}`, url, defaultBranch, basePath)
}
