package octopusdeploy_framework

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/uritemplates"
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
	expand := func(t *testing.T, gitRef string, partialName string, skip int64, take int64) string {
		template, params := runbooksByProjectTemplate("Spaces-1", "Projects-1", gitRef, partialName, skip, take)
		uri, err := uritemplates.NewUriTemplateCache().Expand(template, params)
		require.NoError(t, err)
		return uri
	}

	t.Run("database stored runbooks without a git ref", func(t *testing.T) {
		require.Equal(t, "/api/Spaces-1/projects/Projects-1/runbooks", expand(t, "", "", 0, 0))
	})

	t.Run("git stored runbooks with a branch name", func(t *testing.T) {
		require.Equal(t, "/api/Spaces-1/projects/Projects-1/main/runbooks", expand(t, "main", "", 0, 0))
	})

	t.Run("full git ref is encoded as a single path segment", func(t *testing.T) {
		require.Equal(t, "/api/Spaces-1/projects/Projects-1/refs%2Fheads%2Fmain/runbooks", expand(t, "refs/heads/main", "", 0, 0))
	})

	t.Run("query parameters apply to the database endpoint", func(t *testing.T) {
		require.Equal(t, "/api/Spaces-1/projects/Projects-1/runbooks?skip=5&take=10&partialName=Test", expand(t, "", "Test", 5, 10))
	})

	t.Run("query parameters apply to the git ref endpoint", func(t *testing.T) {
		require.Equal(t, "/api/Spaces-1/projects/Projects-1/main/runbooks?skip=5&take=10&partialName=Test", expand(t, "main", "Test", 5, 10))
	})
}

func TestPageRunbooks(t *testing.T) {
	all := []*runbooks.Runbook{runbookWithID("Runbooks-1"), runbookWithID("Runbooks-2"), runbookWithID("Runbooks-3")}
	ids := func(found []*runbooks.Runbook) []string {
		result := []string{}
		for _, runbook := range found {
			result = append(result, runbook.GetID())
		}
		return result
	}

	t.Run("no skip or take returns everything", func(t *testing.T) {
		require.Equal(t, []string{"Runbooks-1", "Runbooks-2", "Runbooks-3"}, ids(pageRunbooks(all, 0, 0)))
	})

	t.Run("skip drops leading runbooks", func(t *testing.T) {
		require.Equal(t, []string{"Runbooks-2", "Runbooks-3"}, ids(pageRunbooks(all, 1, 0)))
	})

	t.Run("take limits the result", func(t *testing.T) {
		require.Equal(t, []string{"Runbooks-1", "Runbooks-2"}, ids(pageRunbooks(all, 0, 2)))
	})

	t.Run("skip and take together", func(t *testing.T) {
		require.Equal(t, []string{"Runbooks-2"}, ids(pageRunbooks(all, 1, 1)))
	})

	t.Run("take larger than the result returns everything after skip", func(t *testing.T) {
		require.Equal(t, []string{"Runbooks-3"}, ids(pageRunbooks(all, 2, 10)))
	})

	t.Run("skip past the end returns an empty list", func(t *testing.T) {
		require.Empty(t, pageRunbooks(all, 3, 0))
		require.Empty(t, pageRunbooks(all, 10, 1))
	})

	t.Run("an empty result stays empty", func(t *testing.T) {
		require.Empty(t, pageRunbooks([]*runbooks.Runbook{}, 0, 0))
	})
}

func TestProjectStoresRunbooksInGit(t *testing.T) {
	gitProject := func(runbooksInGit bool) *projects.Project {
		repoURL, _ := url.Parse("https://github.com/example/repo")
		settings := projects.NewGitPersistenceSettings(".octopus", credentials.NewGitHubApp("GitHubAppConnections-1"), "main", []string{}, repoURL)
		if runbooksInGit {
			settings.SetRunbooksAreInGit()
		}

		project := projects.NewProject("test", "Lifecycles-1", "ProjectGroups-1")
		project.IsVersionControlled = true
		project.PersistenceSettings = settings
		return project
	}

	t.Run("nil project", func(t *testing.T) {
		require.False(t, projectStoresRunbooksInGit(nil))
	})

	t.Run("database project", func(t *testing.T) {
		project := projects.NewProject("test", "Lifecycles-1", "ProjectGroups-1")
		project.PersistenceSettings = projects.NewDatabasePersistenceSettings()
		require.False(t, projectStoresRunbooksInGit(project))
	})

	t.Run("version controlled project with runbooks in the database", func(t *testing.T) {
		require.False(t, projectStoresRunbooksInGit(gitProject(false)))
	})

	t.Run("version controlled project with runbooks in Git", func(t *testing.T) {
		require.True(t, projectStoresRunbooksInGit(gitProject(true)))
	})
}

func TestRunbooksTake(t *testing.T) {
	t.Run("unset take reads every runbook", func(t *testing.T) {
		require.Equal(t, int64(math.MaxInt32), runbooksTake(0))
	})

	t.Run("negative take reads every runbook", func(t *testing.T) {
		require.Equal(t, int64(math.MaxInt32), runbooksTake(-1))
	})

	t.Run("configured take is kept", func(t *testing.T) {
		require.Equal(t, int64(10), runbooksTake(10))
	})
}
