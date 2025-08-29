package octopusdeploy

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccProjectGroup creates a project group configuration for testing
func testAccProjectGroup(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}`, localName, name)
}

// testAccProjectGroupCheckDestroy checks that all project groups have been destroyed
func testAccProjectGroupCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project_group" {
			continue
		}

		if projectGroup, err := octoClient.ProjectGroups.GetByID(rs.Primary.ID); err == nil {
			return fmt.Errorf("project group (%s) still exists", projectGroup.GetID())
		}
	}

	return nil
}

// testAwsAccount creates an AWS account configuration for testing
func testAwsAccount(localName string, name string, accessKey string, secretKey string) string {
	return fmt.Sprintf(`resource "octopusdeploy_aws_account" "%s" {
		access_key = "%s"
		name       = "%s"
		secret_key = "%s"
	}`, localName, accessKey, name, secretKey)
}