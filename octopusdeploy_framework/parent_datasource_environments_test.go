package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceParentEnvironments(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := fmt.Sprintf("data.octopusdeploy_parent_environments.%s", localName)

	spaceName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	take := 10

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParentEnvironmentsDataSourceID(prefix),
				),
				Config: testAccDataSourceParentEnvironmentsEmpty(localName),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParentEnvironmentsDataSourceID(prefix),
					resource.TestCheckResourceAttr(prefix, "parent_environments.#", "3"),
				),
				Config: fmt.Sprintf(`%s
				
				%s`,
					createTestAccDataSourceParentEnvironmentsConfig(spaceName, environmentLocalName, environmentName),
					testAccDataSourceParentEnvironmentsConfig(localName, take, spaceName, environmentLocalName),
				),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParentEnvironmentsDataSourceID(prefix),
					resource.TestCheckResourceAttr(prefix, "name", environmentName),
					resource.TestCheckResourceAttr(prefix, "parent_environments.#", "1"),
					resource.TestCheckResourceAttrSet(prefix, "parent_environments.0.id"),
					resource.TestCheckResourceAttr(prefix, "parent_environments.0.name", environmentName),
				),
				Config: fmt.Sprintf(`%s

			%s`,
					createTestAccDataSourceParentEnvironmentsConfig(spaceName, environmentLocalName, environmentName),
					testAccDataSourceParentEnvironmentByNameConfig(localName, environmentName, spaceName, environmentLocalName),
				),
			},
		},
	})
}

func testAccCheckParentEnvironmentsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find Environments data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot Environments source ID not set")
		}
		return nil
	}
}

func createTestAccDataSourceParentEnvironmentsConfig(spaceName string, localName string, name string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_space" "%[1]s" {
			name                  = "%[1]s"
			is_default            = false
			is_task_queue_stopped = true
			description           = "Test space for environments datasource"
			space_managers_teams  = ["teams-administrators"]
		}

		resource "octopusdeploy_parent_environment" "%[2]s" {
			name = "%[3]s"
			space_id = octopusdeploy_space.%[1]s.id
		}
		
		resource "octopusdeploy_parent_environment" "%[2]s-1" {
			name = "%[3]s-1"
			space_id = octopusdeploy_space.%[1]s.id
		}

		resource "octopusdeploy_parent_environment" "%[2]s-2" {
			name = "%[3]s-2"
			space_id = octopusdeploy_space.%[1]s.id
		}
	`, spaceName, localName, name)
}

func testAccDataSourceParentEnvironmentsEmpty(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_parent_environments" "%s" {}`, localName)
}

func testAccDataSourceParentEnvironmentsConfig(localName string, take int, spaceName string, environmentLocalName string) string {
	return fmt.Sprintf(`data "octopusdeploy_parent_environments" "%s" {
		take = %v
		space_id = octopusdeploy_space.%s.id
		depends_on = [octopusdeploy_parent_environment.%s, octopusdeploy_parent_environment.%[4]s-1, octopusdeploy_parent_environment.%[4]s-2]
	}`, localName, take, spaceName, environmentLocalName)
}

func testAccDataSourceParentEnvironmentByNameConfig(localName string, name string, spaceName string, environmentLocalName string) string {
	return fmt.Sprintf(`data "octopusdeploy_parent_environments" "%s" {
		name = "%s"	
		space_id = octopusdeploy_space.%s.id
		depends_on = [octopusdeploy_parent_environment.%s, octopusdeploy_parent_environment.%[4]s-1, octopusdeploy_parent_environment.%[4]s-2]
	}`, localName, name, spaceName, environmentLocalName)
}
