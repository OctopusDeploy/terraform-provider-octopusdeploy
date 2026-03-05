package octopusdeploy_framework

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	internaltest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccProjectBasic(t *testing.T) {
	lifecycleTestOptions := internaltest.NewLifecycleTestOptions()
	projectGroupTestOptions := internaltest.NewProjectGroupTestOptions()
	projectTestOptions := internaltest.NewProjectTestOptions(lifecycleTestOptions, projectGroupTestOptions)
	projectTestOptions.Resource.IsDisabled = true

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleTestOptions.Resource.Name),
					testProjectGroupExists(projectGroupTestOptions.QualifiedName),
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(projectTestOptions.QualifiedName, "description", projectTestOptions.Resource.Description),
					resource.TestCheckResourceAttr(projectTestOptions.QualifiedName, "name", projectTestOptions.Resource.Name),
				),
				Config: internaltest.GetConfiguration([]string{
					internaltest.LifecycleConfiguration(lifecycleTestOptions),
					internaltest.ProjectGroupConfiguration(projectGroupTestOptions),
					internaltest.ProjectConfiguration(projectTestOptions),
				}),
			},
		},
	})
}

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

func testProjectGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if _, err := octoClient.ProjectGroups.GetByID(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func TestAccProjectWithUpdate(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project." + localName

	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
				),
				Config: testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description, 2),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.windows_service.0.step_name"),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.windows_service.1.step_name"),
					resource.TestCheckNoResourceAttr(prefix, "deployment_step.0.iis_website.0.step_name"),
				),
				Config: testAccProjectBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, localName, name, description, 3),
			},
		},
	})
}

func testAccProjectBasic(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, localName string, name string, description string, templateCount int) string {
	projectGroup := internaltest.NewProjectGroupTestOptions()
	projectGroup.LocalName = projectGroupLocalName
	projectGroup.Resource.Name = projectGroupName

	var templates string
	for i := 0; i < templateCount; i++ {
		templates += fmt.Sprintf("\ntemplate {\n\t\t\t\tname          = \"%d\"\n\t\t\t\tdisplay_settings = {\n\t\t\t\t\t\"Octopus.ControlType\": \"SingleLineText\"\n\t\t\t\t}\n\t\t\t}\n", i)
	}

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		internaltest.ProjectGroupConfiguration(projectGroup)+"\n"+
		`resource "octopusdeploy_project" "%s" {
			description      = "%s"
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			name             = "%s"
			project_group_id = octopusdeploy_project_group.%s.id

			%s
			
			versioning_strategy {
				template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.LastPatch}.#{Octopus.Version.NextRevision}"
			}

			connectivity_policy {
				allow_deployments_to_no_targets = true
				skip_machine_behavior           = "None"
			}

		}`, localName, description, lifecycleLocalName, name, projectGroupLocalName, templates)
}

func testAccProjectCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_project" {
			continue
		}

		if project, err := projects.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("project (%s) still exists", project.GetID())
		}
	}

	return nil
}

func testAccProjectCheckExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_project" {
				if _, err := projects.GetByID(octoClient, octoClient.GetSpaceID(), r.Primary.ID); err != nil {
					return fmt.Errorf("error retrieving project with ID %s: %s", r.Primary.ID, err)
				}
			}
		}
		return nil
	}
}

func TestAccProjectWithTags(t *testing.T) {
	t.Skip("Skipping - canonical tag name handling needs investigation")

	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSetLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSetName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_project." + projectLocalName

	resource.Test(t, resource.TestCase{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccProjectCheckDestroy,
			testAccProjectGroupCheckDestroy,
			testAccLifecycleCheckDestroy,
		),
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists(),
					resource.TestCheckResourceAttr(prefix, "name", projectName),
					resource.TestCheckResourceAttr(prefix, "description", projectDescription),
					resource.TestCheckResourceAttr(prefix, "project_tags.#", "1"),
				),
				Config: testAccProjectWithTags(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName,
					tagSetLocalName, tagSetName, tagLocalName, tagName, projectLocalName, projectName, projectDescription),
			},
		},
	})
}

func testAccProjectWithTags(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName,
	tagSetLocalName, tagSetName, tagLocalName, tagName, projectLocalName, projectName, projectDescription string) string {
	projectGroup := internaltest.NewProjectGroupTestOptions()
	projectGroup.LocalName = projectGroupLocalName
	projectGroup.Resource.Name = projectGroupName

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		internaltest.ProjectGroupConfiguration(projectGroup)+"\n"+
		`resource "octopusdeploy_tag_set" "%s" {
			name        = "%s"
			description = "Tag set for project testing"
			scopes      = ["Tenant"]
		}

		resource "octopusdeploy_tag" "%s" {
			name        = "%s"
			color       = "#6e6e6e"
			description = "Test tag for project"
			tag_set_id  = octopusdeploy_tag_set.%s.id
		}

		resource "octopusdeploy_project" "%s" {
			name             = "%s"
			description      = "%s"
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			project_group_id = octopusdeploy_project_group.%s.id
			project_tags     = [octopusdeploy_tag.%s.canonical_tag_name]
		}`,
		tagSetLocalName, tagSetName,
		tagLocalName, tagName, tagSetLocalName,
		projectLocalName, projectName, projectDescription, lifecycleLocalName, projectGroupLocalName, tagLocalName)
}

func TestProcessPersistenceSettings_GitHubApp(t *testing.T) {
	ctx := context.Background()

	connectionID := "GitHubAppConnections-1"
	repoURL, _ := url.Parse("https://github.com/example/repo")

	project := projects.NewProject("test", "lifecycle-1", "group-1")
	project.PersistenceSettings = projects.NewGitPersistenceSettings(
		".octopus",
		credentials.NewGitHubApp(connectionID),
		"main",
		[]string{},
		repoURL,
	)

	model := &projectResourceModel{}
	diags := processPersistenceSettings(ctx, project, model)

	require.False(t, diags.HasError())
	assert.True(t, model.IsVersionControlled.ValueBool())

	assert.False(t, model.GitGitHubAppPersistenceSettings.IsNull())
	assert.True(t, model.GitLibraryPersistenceSettings.IsNull())
	assert.True(t, model.GitUsernamePasswordPersistenceSettings.IsNull())
	assert.True(t, model.GitAnonymousPersistenceSettings.IsNull())

	var appSettings []gitGitHubAppPersistenceSettingsModel
	diags = model.GitGitHubAppPersistenceSettings.ElementsAs(ctx, &appSettings, false)
	require.False(t, diags.HasError())
	require.Len(t, appSettings, 1)
	assert.Equal(t, connectionID, appSettings[0].GitHubConnectionID.ValueString())
	assert.Equal(t, "https://github.com/example/repo", appSettings[0].URL.ValueString())
	assert.Equal(t, ".octopus", appSettings[0].BasePath.ValueString())
	assert.Equal(t, "main", appSettings[0].DefaultBranch.ValueString())
}

func TestExpandGitGitHubAppPersistenceSettings(t *testing.T) {
	ctx := context.Background()

	connectionID := "GitHubAppConnections-1"

	model := gitGitHubAppPersistenceSettingsModel{
		GitHubConnectionID: types.StringValue(connectionID),
		URL:                types.StringValue("https://github.com/example/repo"),
		BasePath:           types.StringValue(".octopus"),
		DefaultBranch:      types.StringValue("main"),
		ProtectedBranches:  types.SetValueMust(types.StringType, []attr.Value{}),
	}

	settings := expandGitGitHubAppPersistenceSettings(ctx, model)

	assert.Equal(t, ".octopus", settings.BasePath())
	assert.Equal(t, "main", settings.DefaultBranch())
	assert.Equal(t, "https://github.com/example/repo", settings.URL().String())

	cred := settings.Credential()
	require.NotNil(t, cred)
	assert.Equal(t, credentials.GitCredentialTypeGitHubApp, cred.Type())

	gitHubAppCred, ok := cred.(*credentials.GitHubApp)
	require.True(t, ok)
	assert.Equal(t, connectionID, gitHubAppCred.ID)
}

func TestProcessPersistenceSettings_UsernamePassword(t *testing.T) {
	ctx := context.Background()

	repoURL, _ := url.Parse("https://github.com/example/repo")

	project := projects.NewProject("test", "lifecycle-1", "group-1")
	project.PersistenceSettings = projects.NewGitPersistenceSettings(
		".octopus",
		credentials.NewUsernamePassword("git-user", core.NewSensitiveValue("secret")),
		"main",
		[]string{},
		repoURL,
	)

	model := &projectResourceModel{}
	diags := processPersistenceSettings(ctx, project, model)

	require.False(t, diags.HasError())
	assert.True(t, model.IsVersionControlled.ValueBool())

	assert.True(t, model.GitGitHubAppPersistenceSettings.IsNull())
	assert.True(t, model.GitLibraryPersistenceSettings.IsNull())
	assert.False(t, model.GitUsernamePasswordPersistenceSettings.IsNull())
	assert.True(t, model.GitAnonymousPersistenceSettings.IsNull())

	var settings []gitUsernamePasswordPersistenceSettingsModel
	diags = model.GitUsernamePasswordPersistenceSettings.ElementsAs(ctx, &settings, false)
	require.False(t, diags.HasError())
	require.Len(t, settings, 1)
	assert.Equal(t, "git-user", settings[0].Username.ValueString())
	assert.Equal(t, "secret", settings[0].Password.ValueString())
	assert.Equal(t, "https://github.com/example/repo", settings[0].URL.ValueString())
	assert.Equal(t, ".octopus", settings[0].BasePath.ValueString())
	assert.Equal(t, "main", settings[0].DefaultBranch.ValueString())
}

func TestExpandGitUsernamePasswordPersistenceSettings(t *testing.T) {
	ctx := context.Background()

	model := gitUsernamePasswordPersistenceSettingsModel{
		URL:               types.StringValue("https://github.com/example/repo"),
		Username:          types.StringValue("git-user"),
		Password:          types.StringValue("secret"),
		BasePath:          types.StringValue(".octopus"),
		DefaultBranch:     types.StringValue("main"),
		ProtectedBranches: types.SetValueMust(types.StringType, []attr.Value{}),
	}

	settings := expandGitUsernamePasswordPersistenceSettings(ctx, model)

	assert.Equal(t, ".octopus", settings.BasePath())
	assert.Equal(t, "main", settings.DefaultBranch())
	assert.Equal(t, "https://github.com/example/repo", settings.URL().String())

	cred := settings.Credential()
	require.NotNil(t, cred)
	assert.Equal(t, credentials.GitCredentialTypeUsernamePassword, cred.Type())

	upCred, ok := cred.(*credentials.UsernamePassword)
	require.True(t, ok)
	assert.Equal(t, "git-user", upCred.Username)
	require.NotNil(t, upCred.Password)
	require.NotNil(t, upCred.Password.NewValue)
	assert.Equal(t, "secret", *upCred.Password.NewValue)
}
