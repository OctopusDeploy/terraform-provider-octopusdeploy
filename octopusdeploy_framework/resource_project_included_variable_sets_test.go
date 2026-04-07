package octopusdeploy_framework

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectIncludedLibraryVariableSets(t *testing.T) {
	localName := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectWithVariableSetsConfiguration(localName, []string{"set_a", "set_b", "set_c"}),
				Check:  testCheckProjectWithIncludedVariableSets(localName, []string{"set_a", "set_b", "set_c"}),
			},
			{
				Config: testAccProjectWithVariableSetsConfiguration(localName, []string{"set_d", "set_a", "set_b", "set_c"}),
				Check:  testCheckProjectWithIncludedVariableSets(localName, []string{"set_d", "set_a", "set_b", "set_c"}),
			},
			{
				Config: testAccProjectWithVariableSetsAndDataSource(localName, []string{"set_d", "set_a", "set_b", "set_c"}),
				Check:  testCheckProjectDataSourceWithVariableSets(localName, []string{"set_d", "set_a", "set_b", "set_c"}),
			},
		},
	})
}

func testAccProjectWithVariableSetsConfiguration(localName string, includedVariableSets []string) string {
	includedSetQualifiedNames := make([]string, len(includedVariableSets))
	for i, includedSet := range includedVariableSets {
		includedSetQualifiedNames[i] = fmt.Sprintf("octopusdeploy_library_variable_set.%s_%s.id", localName, includedSet)
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

		resource "octopusdeploy_library_variable_set" "%[1]s_set_a" {
		  name = "%[1]s-vs-a"
		}

		resource "octopusdeploy_library_variable_set" "%[1]s_set_b" {
		  name = "%[1]s-vs-b"
		}

		resource "octopusdeploy_library_variable_set" "%[1]s_set_c" {
		  name = "%[1]s-vs-c"
		}

		resource "octopusdeploy_library_variable_set" "%[1]s_set_d" {
		  name = "%[1]s-vs-d"
		}

		resource "octopusdeploy_project" "%[1]s" {
		  name                              = "%[1]s"
		  description                       = "Test included variable sets ordering"
		  default_guided_failure_mode       = "EnvironmentDefault"
		  tenanted_deployment_participation = "Untenanted"
		  project_group_id                  = octopusdeploy_project_group.%[1]s.id
		  lifecycle_id                      = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		  included_library_variable_sets    = [
			%[2]s
		  ]
		}
		`,
		localName,
		strings.Join(includedSetQualifiedNames, ",\n\t\t\t"),
	)
}

func testAccProjectWithVariableSetsAndDataSource(localName string, includedVariableSets []string) string {
	projectConfiguration := testAccProjectWithVariableSetsConfiguration(localName, includedVariableSets)

	return fmt.Sprintf(`
		%[1]s

		data "octopusdeploy_projects" "all_%[2]s" {
		  name = "%[2]s"
		  depends_on = [octopusdeploy_project.%[2]s]
		}
		`,
		projectConfiguration,
		localName,
	)
}

func testCheckProjectWithIncludedVariableSets(localName string, expectedVariableSets []string) resource.TestCheckFunc {
	projectQualifiedName := fmt.Sprintf("octopusdeploy_project.%s", localName)
	expectedCount := len(expectedVariableSets)

	assertions := []resource.TestCheckFunc{
		testAccProjectCheckExists(),
		resource.TestCheckResourceAttr(projectQualifiedName, "included_library_variable_sets.#", strconv.Itoa(expectedCount)),
	}

	for _, expectedSet := range expectedVariableSets {
		setQualifiedName := fmt.Sprintf("octopusdeploy_library_variable_set.%s_%s", localName, expectedSet)
		assertion := resource.TestCheckTypeSetElemAttrPair(projectQualifiedName, "included_library_variable_sets.*", setQualifiedName, "id")
		assertions = append(assertions, assertion)
	}

	return resource.ComposeTestCheckFunc(assertions...)
}

func testCheckProjectDataSourceWithVariableSets(localName string, expectedVariableSets []string) resource.TestCheckFunc {
	projectsQualifiedName := fmt.Sprintf("data.octopusdeploy_projects.all_%s", localName)
	expectedCount := len(expectedVariableSets)

	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(projectsQualifiedName, "projects.#", "1"),
		resource.TestCheckResourceAttr(projectsQualifiedName, "projects.0.included_library_variable_sets.#", strconv.Itoa(expectedCount)),
	)
}
