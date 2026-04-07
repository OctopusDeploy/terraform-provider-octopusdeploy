package octopusdeploy_framework

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccProjectCaCUpdate verifies that updating a CaC (Config as Code) project
// (e.g. adding a library variable set or changing deployment settings) does not
// fail with the "Cannot update deployment settings" error.
// Requires GIT_URL, GIT_USERNAME, GIT_PASSWORD env vars.
func TestAccProjectCaCUpdate(t *testing.T) {
	gitURL := os.Getenv("GIT_URL")
	gitUsername := os.Getenv("GIT_USERNAME")
	gitPassword := os.Getenv("GIT_PASSWORD")
	if gitURL == "" || gitUsername == "" || gitPassword == "" {
		t.Skip("Skipping CaC project update test: GIT_URL, GIT_USERNAME, GIT_PASSWORD must be set")
	}

	localName := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	basePath := ".octopus/" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	projectPrefix := "octopusdeploy_project." + localName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCaCProjectConfig(localName, basePath, gitURL, gitUsername, gitPassword, "Off", false),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(projectPrefix, "is_version_controlled", "true"),
					resource.TestCheckResourceAttr(projectPrefix, "default_guided_failure_mode", "Off"),
					resource.TestCheckResourceAttr(projectPrefix, "included_library_variable_sets.#", "0"),
				),
			},
			{
				Config: testAccCaCProjectConfig(localName, basePath, gitURL, gitUsername, gitPassword, "On", true),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(projectPrefix, "default_guided_failure_mode", "On"),
					resource.TestCheckResourceAttr(projectPrefix, "included_library_variable_sets.#", "1"),
				),
			},
		},
	})
}

func testAccCaCProjectConfig(localName, basePath, gitURL, gitUsername, gitPassword, guidedFailureMode string, includeLibVarSet bool) string {
	includedSets := "included_library_variable_sets = []"
	if includeLibVarSet {
		includedSets = fmt.Sprintf("included_library_variable_sets = [octopusdeploy_library_variable_set.%s.id]", localName)
	}

	return fmt.Sprintf(`
		data "octopusdeploy_lifecycles" "default" {
		  ids          = null
		  partial_name = "Default Lifecycle"
		  skip         = 0
		  take         = 1
		}

		resource "octopusdeploy_project_group" "%[1]s" {
		  name = "%[1]s"
		}

		resource "octopusdeploy_git_credential" "%[1]s" {
		  name     = "%[1]s"
		  username = "%[4]s"
		  password = "%[5]s"
		}

		resource "octopusdeploy_library_variable_set" "%[1]s" {
		  name = "%[1]s-lvs"
		}

		resource "octopusdeploy_project" "%[1]s" {
		  name                        = "%[1]s"
		  default_guided_failure_mode = "%[6]s"
		  is_version_controlled       = true
		  project_group_id            = octopusdeploy_project_group.%[1]s.id
		  lifecycle_id                = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		  %[7]s

		  git_library_persistence_settings {
		    git_credential_id = octopusdeploy_git_credential.%[1]s.id
		    url               = "%[3]s"
		    base_path         = "%[2]s"
		    default_branch    = "main"
		  }
		}
		`,
		localName,  // 1
		basePath,   // 2
		gitURL,     // 3
		gitUsername, // 4
		gitPassword, // 5
		guidedFailureMode, // 6
		includedSets, // 7
	)
}
