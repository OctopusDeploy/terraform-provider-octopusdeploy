package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/uritemplates"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
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

func TestAccDataSourceRunbooksGitRefRequiresProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
				data "octopusdeploy_runbooks" "git_ref_only" {
					git_ref = "main"
				}`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Attribute "project_id" must be specified when "git_ref" is\s+specified`),
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

func runbookWithID(id string) *runbooks.Runbook {
	runbook := runbooks.NewRunbook("runbook-"+id, "Projects-1")
	runbook.Resource = resources.Resource{ID: id}
	return runbook
}

func TestFilterRunbooksByID(t *testing.T) {
	all := []*runbooks.Runbook{runbookWithID("Runbooks-1"), runbookWithID("Runbooks-2"), runbookWithID("Runbooks-3")}

	t.Run("no ids returns everything", func(t *testing.T) {
		require.Len(t, filterRunbooksByID(all, nil), 3)
		require.Len(t, filterRunbooksByID(all, []string{}), 3)
	})

	t.Run("keeps only the requested ids", func(t *testing.T) {
		filtered := filterRunbooksByID(all, []string{"Runbooks-3", "Runbooks-1"})

		require.Len(t, filtered, 2)
		require.Equal(t, "Runbooks-1", filtered[0].GetID())
		require.Equal(t, "Runbooks-3", filtered[1].GetID())
	})

	t.Run("unknown ids match nothing", func(t *testing.T) {
		require.Empty(t, filterRunbooksByID(all, []string{"Runbooks-99"}))
	})
}

func TestRunbooksByProjectTemplate(t *testing.T) {
	expand := func(t *testing.T, gitRef string, data runbooksDataSourceModel) string {
		template, params := runbooksByProjectTemplate("Spaces-1", "Projects-1", gitRef, data)
		uri, err := uritemplates.NewUriTemplateCache().Expand(template, params)
		require.NoError(t, err)
		return uri
	}

	t.Run("database stored runbooks without a git ref", func(t *testing.T) {
		require.Equal(t, "/api/Spaces-1/projects/Projects-1/runbooks", expand(t, "", runbooksDataSourceModel{}))
	})

	t.Run("git stored runbooks with a branch name", func(t *testing.T) {
		require.Equal(t, "/api/Spaces-1/projects/Projects-1/main/runbooks", expand(t, "main", runbooksDataSourceModel{}))
	})

	t.Run("full git ref is encoded as a single path segment", func(t *testing.T) {
		require.Equal(t, "/api/Spaces-1/projects/Projects-1/refs%2Fheads%2Fmain/runbooks", expand(t, "refs/heads/main", runbooksDataSourceModel{}))
	})

	t.Run("query parameters apply to the git ref endpoint", func(t *testing.T) {
		data := runbooksDataSourceModel{
			PartialName: types.StringValue("Test"),
			Skip:        types.Int64Value(5),
			Take:        types.Int64Value(10),
		}
		require.Equal(t, "/api/Spaces-1/projects/Projects-1/main/runbooks?skip=5&take=10&partialName=Test", expand(t, "main", data))
	})
}
