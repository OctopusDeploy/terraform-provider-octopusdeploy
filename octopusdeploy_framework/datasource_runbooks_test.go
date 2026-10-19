package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceRunbooks(t *testing.T) {
	prefix := acctest.RandStringFromCharSet(12, acctest.CharSetAlpha)
	config := testAccDataSourceRunbooksConfig(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// Space wide query, narrowed by name so the assertion is deterministic.
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.space_wide", "runbooks.#", "2"),
					resource.TestCheckResourceAttrSet("data.octopusdeploy_runbooks.space_wide", "runbooks.0.id"),
					resource.TestCheckResourceAttrSet("data.octopusdeploy_runbooks.space_wide", "runbooks.0.project_id"),
					resource.TestCheckResourceAttrSet("data.octopusdeploy_runbooks.space_wide", "runbooks.0.space_id"),

					// Project scoped query.
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.by_project", "runbooks.#", "2"),

					// Project scoped, narrowed further by partial name.
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.by_partial_name", "runbooks.#", "1"),
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.by_partial_name", "runbooks.0.name", prefix+"Alpha"),
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.by_partial_name", "runbooks.0.description", "first probe runbook"),
					resource.TestCheckResourceAttrSet("data.octopusdeploy_runbooks.by_partial_name", "runbooks.0.runbook_process_id"),
					resource.TestCheckResourceAttrSet("data.octopusdeploy_runbooks.by_partial_name", "runbooks.0.multi_tenancy_mode"),
					resource.TestCheckResourceAttrSet("data.octopusdeploy_runbooks.by_partial_name", "runbooks.0.environment_scope"),
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.by_partial_name", "runbooks.0.retention_policy.#", "1"),

					// ids, on the space wide endpoint and on the project scoped one, where
					// the endpoint has no ids parameter and the filtering is done locally.
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.by_ids", "runbooks.#", "1"),
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.by_ids", "runbooks.0.name", prefix+"Beta"),
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.by_ids_in_project", "runbooks.#", "1"),
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.by_ids_in_project", "runbooks.0.name", prefix+"Beta"),

					// A project with no runbooks returns an empty list rather than an error.
					resource.TestCheckResourceAttr("data.octopusdeploy_runbooks.empty_project", "runbooks.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceRunbooksConfig(prefix string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_lifecycle" "rbds" {
		name = "%[1]slifecycle"
	}

	resource "octopusdeploy_project_group" "rbds" {
		name = "%[1]sgroup"
	}

	resource "octopusdeploy_project" "rbds" {
		name             = "%[1]sproject"
		lifecycle_id     = octopusdeploy_lifecycle.rbds.id
		project_group_id = octopusdeploy_project_group.rbds.id
	}

	resource "octopusdeploy_project" "rbds_empty" {
		name             = "%[1]sempty"
		lifecycle_id     = octopusdeploy_lifecycle.rbds.id
		project_group_id = octopusdeploy_project_group.rbds.id
	}

	resource "octopusdeploy_runbook" "alpha" {
		project_id  = octopusdeploy_project.rbds.id
		name        = "%[1]sAlpha"
		description = "first probe runbook"
	}

	resource "octopusdeploy_runbook" "beta" {
		project_id  = octopusdeploy_project.rbds.id
		name        = "%[1]sBeta"
		description = "second probe runbook"
	}

	data "octopusdeploy_runbooks" "space_wide" {
		partial_name = "%[1]s"
		depends_on   = [octopusdeploy_runbook.alpha, octopusdeploy_runbook.beta]
	}

	data "octopusdeploy_runbooks" "by_project" {
		project_id = octopusdeploy_project.rbds.id
		depends_on = [octopusdeploy_runbook.alpha, octopusdeploy_runbook.beta]
	}

	data "octopusdeploy_runbooks" "by_partial_name" {
		project_id   = octopusdeploy_project.rbds.id
		partial_name = "%[1]sAlpha"
		depends_on   = [octopusdeploy_runbook.alpha, octopusdeploy_runbook.beta]
	}

	data "octopusdeploy_runbooks" "by_ids" {
		ids = [octopusdeploy_runbook.beta.id]
	}

	data "octopusdeploy_runbooks" "by_ids_in_project" {
		project_id = octopusdeploy_project.rbds.id
		ids        = [octopusdeploy_runbook.beta.id]
	}

	data "octopusdeploy_runbooks" "empty_project" {
		project_id = octopusdeploy_project.rbds_empty.id
	}
	`, prefix)
}
