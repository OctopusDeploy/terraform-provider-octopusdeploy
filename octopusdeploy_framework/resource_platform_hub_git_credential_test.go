package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/platformhubgitcredential"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

// TestAccPlatformHubGitCredentialCreate tests resource creation
func TestAccPlatformHubGitCredentialCreate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_git_credential." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubGitCredentialCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubGitCredentialBasic(localName, name, description),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGitCredentialExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "username", "test_user"),
					resource.TestCheckResourceAttr(resourceName, "repository_restrictions.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "repository_restrictions.allowed_repositories.#", "0"),
				),
			},
		},
	})
}

func TestAccPlatformHubGitCredentialUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_git_credential." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubGitCredentialCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubGitCredentialBasic(localName, name, description),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGitCredentialExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testPlatformHubGitCredentialBasic(localName, updatedName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGitCredentialExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccPlatformHubGitCredentialImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_git_credential." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubGitCredentialCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubGitCredentialBasic(localName, name, description),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGitCredentialExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccPlatformHubGitCredentialRepositoryRestrictions(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_platform_hub_git_credential." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testPlatformHubGitCredentialCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testPlatformHubGitCredentialWithRepositoryRestrictions(localName, name, description),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGitCredentialExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "repository_restrictions.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "repository_restrictions.allowed_repositories.#", "2"),
				),
			},
			{
				Config: testPlatformHubGitCredentialBasic(localName, name, description),
				Check: resource.ComposeTestCheckFunc(
					testPlatformHubGitCredentialExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "repository_restrictions.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "repository_restrictions.allowed_repositories.#", "0"),
				),
			},
		},
	})
}

func testPlatformHubGitCredentialBasic(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_platform_hub_git_credential" "%s" {
		name        = "%s"
		description = "%s"
		username    = "test_user"
		password    = "test_password_123"

		repository_restrictions = {
			enabled              = false
			allowed_repositories = []
		}
	}`, localName, name, description)
}

func testPlatformHubGitCredentialWithRepositoryRestrictions(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_platform_hub_git_credential" "%s" {
		name        = "%s"
		description = "%s"
		username    = "test_user"
		password    = "test_password_123"

		repository_restrictions = {
			enabled              = true
			allowed_repositories = [
				"https://github.com/example/repo1",
				"https://github.com/example/repo2"
			]
		}
	}`, localName, name, description)
}

func testPlatformHubGitCredentialExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Platform Hub Git Credential ID is set")
		}

		client := octoClient
		gitCredential, err := platformhubgitcredential.GetByID(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error retrieving Platform Hub Git Credential: %s", err)
		}

		if gitCredential.GetID() != rs.Primary.ID {
			return fmt.Errorf("Platform Hub Git Credential not found")
		}

		return nil
	}
}

func testPlatformHubGitCredentialCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_platform_hub_git_credential" {
			continue
		}

		client := octoClient
		gitCredential, err := platformhubgitcredential.GetByID(client, rs.Primary.ID)
		if err == nil && gitCredential != nil {
			return fmt.Errorf("Platform Hub Git Credential (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
