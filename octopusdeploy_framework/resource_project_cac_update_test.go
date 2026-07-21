package octopusdeploy_framework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccProjectCaCUpdate verifies that updating a CaC (Config as Code) project
// (e.g. adding a library variable set or changing deployment settings) does not
// fail with the "Cannot update deployment settings" error.
// Requires GIT_URL, GIT_USERNAME, and GIT_PASSWORD or GIT_CREDENTIAL env vars.
func TestAccProjectCaCUpdate(t *testing.T) {
	gitURL, gitUsername, gitPassword := testAccGitSettings()
	if gitURL == "" || gitUsername == "" || gitPassword == "" {
		t.Skip("Skipping CaC project update test: GIT_URL, GIT_USERNAME, and GIT_PASSWORD or GIT_CREDENTIAL must be set")
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

// TestAccProjectCaCProtectedMainAllowsDatabaseBackedUpdates exercises the
// compatibility contract for the current octopusdeploy_project schema.
// Database-backed and repository-metadata changes must not implicitly write CaC
// DeploymentSettings, while an explicitly configured OCL-backed change retains
// the existing schema behavior and is rejected by Octopus protection.
//
// Protection in this test is exclusively Octopus ProtectedBranchNamePatterns;
// the Git repository does not need GitHub branch protection.
func TestAccProjectCaCProtectedMainAllowsDatabaseBackedUpdates(t *testing.T) {
	gitURL, gitUsername, gitPassword := testAccGitSettings()
	if gitURL == "" || gitUsername == "" || gitPassword == "" {
		t.Skip("Skipping CaC protected-main database-backed update test: GIT_URL, GIT_USERNAME, and GIT_PASSWORD or GIT_CREDENTIAL must be set")
	}

	localName := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	basePath := ".octopus/" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	projectPrefix := "octopusdeploy_project." + localName
	baselineDescription := "CaC database baseline"
	updatedDescription := "CaC database update"
	baselineReleaseNotes := "CaC release notes baseline"
	mutatedReleaseNotes := "CaC release notes mutation"
	var projectID string

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCaCProjectBehaviorConfig(cacProjectBehaviorConfig{
					LocalName:            localName,
					BasePath:             basePath,
					GitURL:               gitURL,
					GitUsername:          gitUsername,
					GitPassword:          gitPassword,
					DefaultBranch:        "main",
					ProtectedBranches:    nil,
					Description:          baselineDescription,
					ReleaseNotesTemplate: baselineReleaseNotes,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCaptureProjectID(projectPrefix, &projectID),
					resource.TestCheckResourceAttr(projectPrefix, "description", baselineDescription),
					resource.TestCheckResourceAttr(projectPrefix, "release_notes_template", baselineReleaseNotes),
				),
			},
			{
				// Repository protection metadata is Project-owned. Applying it
				// must not echo DeploymentSettings after main becomes protected.
				Config: testAccCaCProjectBehaviorConfig(cacProjectBehaviorConfig{
					LocalName:            localName,
					BasePath:             basePath,
					GitURL:               gitURL,
					GitUsername:          gitUsername,
					GitPassword:          gitPassword,
					DefaultBranch:        "main",
					ProtectedBranches:    []string{"main"},
					Description:          baselineDescription,
					ReleaseNotesTemplate: baselineReleaseNotes,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(projectPrefix, "git_library_persistence_settings.0.protected_branches.#", "1"),
					testAccCheckCaCProject(&projectID, "main", baselineDescription, []string{"main"}),
				),
			},
			{
				// Description is database-backed. Protected main must not make
				// this update depend on a DeploymentSettings write.
				Config: testAccCaCProjectBehaviorConfig(cacProjectBehaviorConfig{
					LocalName:            localName,
					BasePath:             basePath,
					GitURL:               gitURL,
					GitUsername:          gitUsername,
					GitPassword:          gitPassword,
					DefaultBranch:        "main",
					ProtectedBranches:    []string{"main"},
					Description:          updatedDescription,
					ReleaseNotesTemplate: baselineReleaseNotes,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(projectPrefix, "description", updatedDescription),
					testAccCheckCaCProject(&projectID, "main", updatedDescription, []string{"main"}),
				),
			},
			{
				// Included library variable sets are also database-backed.
				Config: testAccCaCProjectBehaviorConfig(cacProjectBehaviorConfig{
					LocalName:            localName,
					BasePath:             basePath,
					GitURL:               gitURL,
					GitUsername:          gitUsername,
					GitPassword:          gitPassword,
					DefaultBranch:        "main",
					ProtectedBranches:    []string{"main"},
					Description:          updatedDescription,
					ReleaseNotesTemplate: baselineReleaseNotes,
					IncludeLibrarySet:    true,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(projectPrefix, "included_library_variable_sets.#", "1"),
					testAccCheckCaCProject(&projectID, "main", updatedDescription, []string{"main"}),
				),
			},
			{
				// The current schema still permits explicit OCL ownership. The
				// tactical fix must retain that behavior until a future schema
				// change: Octopus rejects this write because main is protected.
				Config: testAccCaCProjectBehaviorConfig(cacProjectBehaviorConfig{
					LocalName:            localName,
					BasePath:             basePath,
					GitURL:               gitURL,
					GitUsername:          gitUsername,
					GitPassword:          gitPassword,
					DefaultBranch:        "main",
					ProtectedBranches:    []string{"main"},
					Description:          updatedDescription,
					ReleaseNotesTemplate: mutatedReleaseNotes,
					IncludeLibrarySet:    true,
				}),
				ExpectError: regexp.MustCompile(`(?i)protected branch|branch.*protected`),
			},
		},
	})
}

// TestAccProjectCaCDefaultBranchFeatureToProtectedMainWithValidOCL verifies
// that moving the default from a writable feature branch to protected main is a
// Project-pointer operation once main contains valid OCL for the project. It
// also verifies that a target missing OCL is rejected before the remote pointer
// changes.
//
// GitHub is used only to create disposable refs with exact or missing content;
// Octopus ProtectedBranchNamePatterns is the only protection under test.
func TestAccProjectCaCDefaultBranchFeatureToProtectedMainWithValidOCL(t *testing.T) {
	gitURL, gitUsername, gitPassword := testAccGitSettings()
	if gitURL == "" || gitUsername == "" || gitPassword == "" {
		t.Skip("Skipping CaC feature-to-protected-main default-branch test: GIT_URL, GIT_USERNAME, and GIT_PASSWORD or GIT_CREDENTIAL must be set")
	}

	github, err := newGitHubRefFixture(gitURL, gitPassword)
	if err != nil {
		t.Skipf("Skipping GitHub-specific CaC transition fixture: %s", err)
	}

	mainBeforeCreate, err := github.branchSHA("main")
	if err != nil {
		t.Fatalf("failed to capture main before CaC creation: %s", err)
	}

	localName := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	basePath := ".octopus/" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	projectPrefix := "octopusdeploy_project." + localName
	exactBranch := "octopus-test/exact-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	missingBranch := "octopus-test/missing-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	description := "CaC branch transition"
	releaseNotes := "CaC transition release notes"
	var projectID string
	refsCreated := false

	t.Cleanup(func() {
		if !refsCreated {
			return
		}
		if err := github.deleteBranch(exactBranch); err != nil {
			t.Logf("failed to delete exact-content fixture branch %q: %s", exactBranch, err)
		}
		if err := github.deleteBranch(missingBranch); err != nil {
			t.Logf("failed to delete missing-content fixture branch %q: %s", missingBranch, err)
		}
	})

	config := func(defaultBranch string, protected []string) string {
		return testAccCaCProjectBehaviorConfig(cacProjectBehaviorConfig{
			LocalName:            localName,
			BasePath:             basePath,
			GitURL:               gitURL,
			GitUsername:          gitUsername,
			GitPassword:          gitPassword,
			DefaultBranch:        defaultBranch,
			ProtectedBranches:    protected,
			Description:          description,
			ReleaseNotesTemplate: releaseNotes,
		})
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccProjectCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config("main", nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCaptureProjectID(projectPrefix, &projectID),
					func(_ *terraform.State) error {
						mainWithOCL, err := github.branchSHA("main")
						if err != nil {
							return fmt.Errorf("read main after CaC creation: %w", err)
						}
						if err := github.createBranch(exactBranch, mainWithOCL); err != nil {
							return fmt.Errorf("create exact-content branch: %w", err)
						}
						if err := github.createBranch(missingBranch, mainBeforeCreate); err != nil {
							_ = github.deleteBranch(exactBranch)
							return fmt.Errorf("create missing-content branch: %w", err)
						}
						refsCreated = true
						return nil
					},
				),
			},
			{
				// Establish a writable feature default while main is already
				// protected in Octopus. The feature ref is an exact copy of main.
				Config: config(exactBranch, []string{"main"}),
				Check:  testAccCheckCaCProject(&projectID, exactBranch, description, []string{"main"}),
			},
			{
				// Selecting protected main with exact OCL must be pointer-only.
				// Any DeploymentSettings PUT would be rejected by Octopus.
				Config: config("main", []string{"main"}),
				Check:  testAccCheckCaCProject(&projectID, "main", description, []string{"main"}),
			},
			{
				// The missing branch predates this project's OCL. The provider
				// must preflight it and fail before changing the Project pointer.
				Config:      config(missingBranch, []string{"main"}),
				ExpectError: regexp.MustCompile(`(?i)configuration.*not found|schema_version\.ocl|does not exist|missing`),
			},
			{
				PreConfig: func() {
					project, err := projects.GetByID(octoClient, octoClient.GetSpaceID(), projectID)
					if err != nil {
						t.Fatalf("read project after rejected missing-OCL transition: %s", err)
					}
					settings, ok := project.PersistenceSettings.(projects.GitPersistenceSettings)
					if !ok {
						t.Fatalf("project %s does not have Git persistence settings", projectID)
					}
					if settings.DefaultBranch() != "main" {
						t.Fatalf("missing-OCL transition partially changed default branch: got %q, want main", settings.DefaultBranch())
					}
				},
				Config: config("main", []string{"main"}),
				Check:  testAccCheckCaCProject(&projectID, "main", description, []string{"main"}),
			},
		},
	})
}

func testAccGitSettings() (string, string, string) {
	gitPassword := os.Getenv("GIT_PASSWORD")
	if gitPassword == "" {
		gitPassword = os.Getenv("GIT_CREDENTIAL")
	}

	return os.Getenv("GIT_URL"), os.Getenv("GIT_USERNAME"), gitPassword
}

func testAccGitResourceRuleSettings() (string, string, string, string, bool) {
	basePath, basePathSet := os.LookupEnv("GIT_RESOURCE_RULE_BASE_PATH")
	actionSlug, actionSlugSet := os.LookupEnv("GIT_RESOURCE_RULE_ACTION_SLUG")
	dependencyName, dependencyNameSet := os.LookupEnv("GIT_RESOURCE_RULE_DEPENDENCY_NAME")
	guidedFailureMode := os.Getenv("GIT_RESOURCE_RULE_GUIDED_FAILURE_MODE")
	if guidedFailureMode == "" {
		guidedFailureMode = "EnvironmentDefault"
	}

	return basePath, actionSlug, dependencyName, guidedFailureMode, basePathSet && actionSlugSet && dependencyNameSet
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
		localName,         // 1
		basePath,          // 2
		gitURL,            // 3
		gitUsername,       // 4
		gitPassword,       // 5
		guidedFailureMode, // 6
		includedSets,      // 7
	)
}

type cacProjectBehaviorConfig struct {
	LocalName            string
	BasePath             string
	GitURL               string
	GitUsername          string
	GitPassword          string
	DefaultBranch        string
	ProtectedBranches    []string
	Description          string
	ReleaseNotesTemplate string
	IncludeLibrarySet    bool
}

func testAccCaCProjectBehaviorConfig(config cacProjectBehaviorConfig) string {
	protectedBranches, _ := json.Marshal(config.ProtectedBranches)
	if config.ProtectedBranches == nil {
		protectedBranches = []byte("[]")
	}

	includedSets := "included_library_variable_sets = []"
	if config.IncludeLibrarySet {
		includedSets = fmt.Sprintf("included_library_variable_sets = [octopusdeploy_library_variable_set.%s.id]", config.LocalName)
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
		  name                          = "%[1]s"
		  description                   = "%[8]s"
		  default_guided_failure_mode   = "Off"
		  default_to_skip_if_already_installed = false
		  release_notes_template        = "%[9]s"
		  is_version_controlled         = true
		  project_group_id              = octopusdeploy_project_group.%[1]s.id
		  lifecycle_id                  = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		  %[10]s

		  git_library_persistence_settings {
		    git_credential_id  = octopusdeploy_git_credential.%[1]s.id
		    url                = "%[3]s"
		    base_path          = "%[2]s"
		    default_branch     = "%[6]s"
		    protected_branches = %[7]s
		  }

		  connectivity_policy {
		    allow_deployments_to_no_targets = true
		    skip_machine_behavior           = "None"
		  }
		}
		`,
		config.LocalName,            // 1
		config.BasePath,             // 2
		config.GitURL,               // 3
		config.GitUsername,          // 4
		config.GitPassword,          // 5
		config.DefaultBranch,        // 6
		string(protectedBranches),   // 7
		config.Description,          // 8
		config.ReleaseNotesTemplate, // 9
		includedSets,                // 10
	)
}

func testAccCaptureProjectID(resourceName string, projectID *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resourceState, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("project resource %q was not found in Terraform state", resourceName)
		}
		*projectID = resourceState.Primary.ID
		if *projectID == "" {
			return fmt.Errorf("project resource %q has an empty ID", resourceName)
		}
		return nil
	}
}

func testAccCheckCaCProject(projectID *string, expectedDefault, expectedDescription string, expectedProtected []string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if projectID == nil || *projectID == "" {
			return fmt.Errorf("project ID was not captured")
		}
		project, err := projects.GetByID(octoClient, octoClient.GetSpaceID(), *projectID)
		if err != nil {
			return fmt.Errorf("read project %s: %w", *projectID, err)
		}
		if project.Description != expectedDescription {
			return fmt.Errorf("project %s description = %q, want %q", *projectID, project.Description, expectedDescription)
		}

		settings, ok := project.PersistenceSettings.(projects.GitPersistenceSettings)
		if !ok {
			return fmt.Errorf("project %s does not have Git persistence settings", *projectID)
		}
		if settings.DefaultBranch() != expectedDefault {
			return fmt.Errorf("project %s default branch = %q, want %q", *projectID, settings.DefaultBranch(), expectedDefault)
		}

		actualProtected := settings.ProtectedBranchNamePatterns()
		for _, expected := range expectedProtected {
			found := false
			for _, actual := range actualProtected {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("project %s protected branches = %v, want branch %q", *projectID, actualProtected, expected)
			}
		}
		return nil
	}
}

type githubRefFixture struct {
	owner  string
	repo   string
	token  string
	client *http.Client
}

func newGitHubRefFixture(repositoryURL, token string) (*githubRefFixture, error) {
	parsed, err := url.Parse(repositoryURL)
	if err != nil {
		return nil, fmt.Errorf("parse repository URL: %w", err)
	}
	if !strings.EqualFold(parsed.Hostname(), "github.com") {
		return nil, fmt.Errorf("repository host %q is not github.com", parsed.Hostname())
	}
	path := strings.TrimSuffix(strings.Trim(parsed.Path, "/"), ".git")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("repository URL does not identify owner/repo: %s", repositoryURL)
	}
	return &githubRefFixture{
		owner:  parts[0],
		repo:   parts[1],
		token:  token,
		client: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (g *githubRefFixture) request(method, path string, body any, expected ...int) ([]byte, error) {
	var requestBody io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(encoded)
	}

	request, err := http.NewRequest(method, "https://api.github.com"+path, requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+g.token)
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("User-Agent", "terraform-provider-octopusdeploy-acceptance-tests")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := g.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	for _, status := range expected {
		if response.StatusCode == status {
			return responseBody, nil
		}
	}
	return nil, fmt.Errorf("GitHub %s %s returned HTTP %d: %s", method, path, response.StatusCode, string(responseBody))
}

func (g *githubRefFixture) branchSHA(branch string) (string, error) {
	path := fmt.Sprintf("/repos/%s/%s/git/ref/heads/%s", g.owner, g.repo, url.PathEscape(branch))
	body, err := g.request(http.MethodGet, path, nil, http.StatusOK)
	if err != nil {
		return "", err
	}
	var response struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	if response.Object.SHA == "" {
		return "", fmt.Errorf("GitHub returned an empty SHA for branch %q", branch)
	}
	return response.Object.SHA, nil
}

func (g *githubRefFixture) createBranch(branch, sha string) error {
	path := fmt.Sprintf("/repos/%s/%s/git/refs", g.owner, g.repo)
	_, err := g.request(http.MethodPost, path, map[string]string{
		"ref": "refs/heads/" + branch,
		"sha": sha,
	}, http.StatusCreated)
	return err
}

func (g *githubRefFixture) deleteBranch(branch string) error {
	path := fmt.Sprintf("/repos/%s/%s/git/refs/heads/%s", g.owner, g.repo, url.PathEscape(branch))
	_, err := g.request(http.MethodDelete, path, nil, http.StatusNoContent, http.StatusNotFound)
	return err
}
