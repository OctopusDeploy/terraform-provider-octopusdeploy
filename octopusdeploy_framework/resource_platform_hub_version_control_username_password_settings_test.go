package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPlatformHubVersionControlUsernamePasswordSettingsBasic(t *testing.T) {
	resourceName := "octopusdeploy_platform_hub_version_control_username_password_settings.test"

	url := "https://github.com/acme/hello-world.git"
	defaultBranch := "main"
	basePath := ".octopus"
	username := "test-user"
	password := "test-password"

	updatedURL := "https://github.com/acme/hello-world.git"
	updatedBranch := "develop"
	updatedBasePath := ".octopus/updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubVersionControlUsernamePasswordSettings(url, defaultBranch, basePath, username, password),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubVersionControlSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "default_branch", defaultBranch),
					resource.TestCheckResourceAttr(resourceName, "base_path", basePath),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testPlatformHubVersionControlUsernamePasswordSettings(updatedURL, updatedBranch, updatedBasePath, username, password),
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

func testPlatformHubVersionControlUsernamePasswordSettings(url, defaultBranch, basePath, username, password string) string {
	return fmt.Sprintf(`
resource "octopusdeploy_platform_hub_version_control_username_password_settings" "test" {
	url            = "%s"
	default_branch = "%s"
	base_path      = "%s"
	username       = "%s"
	password       = "%s"
}`, url, defaultBranch, basePath, username, password)
}

func testPlatformHubVersionControlSettingsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		return nil
	}
}
