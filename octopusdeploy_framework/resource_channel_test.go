package octopusdeploy_framework

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccChannelBasic(t *testing.T) {
	options := test.NewChannelTestOptions()

	// Create test dependencies
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	channel := channels.Channel{
		Name:        acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		Description: acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	}

	resourceName := fmt.Sprintf("octopusdeploy_channel.%s", options.LocalName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testChannelBasic(options.LocalName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channel),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", channel.Name),
					resource.TestCheckResourceAttr(resourceName, "description", channel.Description),
				),
			},
		},
	})
}

func TestBuildChannelUpdateRequest_ClearsEmptyCollections(t *testing.T) {
	channel := channels.NewChannel("channel", "Projects-1")
	channel.CustomFieldDefinitions = []channels.ChannelCustomFieldDefinition{{
		FieldName:   "field",
		Description: "description",
	}}
	channel.Rules = []channels.ChannelRule{{
		ID:           "ChannelRules-1",
		Tag:          "beta",
		VersionRange: "[1.0.0,2.0.0)",
	}}
	channel.TenantTags = []string{"TagSets-1/tag-a"}
	channel.GitReferenceRules = []string{"refs/heads/main"}
	channel.GitResourceRules = []channels.ChannelGitResourceRule{{
		Id:    "ChannelGitResourceRules-1",
		Rules: []string{"refs/tags/*"},
		GitDependencyActions: []gitdependencies.DeploymentActionGitDependency{{
			DeploymentActionSlug: "deploy-package",
			GitDependencyName:    "app-config",
		}},
	}}

	plan := schemas.ChannelModel{
		CustomFieldDefinitions: types.ListNull(types.ObjectType{AttrTypes: getChannelCustomFieldDefinitionAttrTypes()}),
		GitReferenceRules:      types.ListNull(types.StringType),
		GitResourceRules:       types.ListNull(types.ObjectType{AttrTypes: getChannelGitResourceRuleAttrTypes()}),
		Rule:                   types.ListNull(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()}),
		TenantTags:             types.SetNull(types.StringType),
	}

	updateReq := buildChannelUpdateRequest(channel, plan)
	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	var payload map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &payload))

	assert.JSONEq(t, `[]`, string(payload["CustomFieldDefinitions"]))
	assert.JSONEq(t, `[]`, string(payload["GitReferenceRules"]))
	assert.JSONEq(t, `[]`, string(payload["GitResourceRules"]))
	assert.JSONEq(t, `[]`, string(payload["Rules"]))
	assert.JSONEq(t, `[]`, string(payload["TenantTags"]))
}

func TestBuildChannelUpdateRequest_ClearsExplicitEmptyCollections(t *testing.T) {
	channel := channels.NewChannel("channel", "Projects-1")
	channel.CustomFieldDefinitions = []channels.ChannelCustomFieldDefinition{{
		FieldName:   "field",
		Description: "description",
	}}
	channel.Rules = []channels.ChannelRule{{
		ID:           "ChannelRules-1",
		Tag:          "beta",
		VersionRange: "[1.0.0,2.0.0)",
	}}
	channel.TenantTags = []string{"TagSets-1/tag-a"}
	channel.GitReferenceRules = []string{"refs/heads/main"}
	channel.GitResourceRules = []channels.ChannelGitResourceRule{{
		Id:    "ChannelGitResourceRules-1",
		Rules: []string{"refs/tags/*"},
		GitDependencyActions: []gitdependencies.DeploymentActionGitDependency{{
			DeploymentActionSlug: "deploy-package",
			GitDependencyName:    "app-config",
		}},
	}}

	plan := schemas.ChannelModel{
		CustomFieldDefinitions: types.ListValueMust(types.ObjectType{AttrTypes: getChannelCustomFieldDefinitionAttrTypes()}, []attr.Value{}),
		GitReferenceRules:      types.ListValueMust(types.StringType, []attr.Value{}),
		GitResourceRules:       types.ListValueMust(types.ObjectType{AttrTypes: getChannelGitResourceRuleAttrTypes()}, []attr.Value{}),
		Rule:                   types.ListValueMust(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()}, []attr.Value{}),
		TenantTags:             types.SetValueMust(types.StringType, []attr.Value{}),
	}

	updateReq := buildChannelUpdateRequest(channel, plan)
	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	var payload map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &payload))

	assert.JSONEq(t, `[]`, string(payload["CustomFieldDefinitions"]))
	assert.JSONEq(t, `[]`, string(payload["GitReferenceRules"]))
	assert.JSONEq(t, `[]`, string(payload["GitResourceRules"]))
	assert.JSONEq(t, `[]`, string(payload["Rules"]))
	assert.JSONEq(t, `[]`, string(payload["TenantTags"]))
}

func TestExpandAndFlattenChannelGitRules(t *testing.T) {
	model := schemas.ChannelModel{
		GitReferenceRules: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("refs/heads/main"),
			types.StringValue("refs/tags/release-*"),
		}),
		GitResourceRules: types.ListValueMust(types.ObjectType{AttrTypes: getChannelGitResourceRuleAttrTypes()}, []attr.Value{
			types.ObjectValueMust(getChannelGitResourceRuleAttrTypes(), map[string]attr.Value{
				"id":    types.StringValue("ChannelGitResourceRules-1"),
				"rules": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("refs/heads/release/*")}),
				"git_dependency_actions": types.ListValueMust(types.ObjectType{AttrTypes: getDeploymentActionGitDependencyAttrTypes()}, []attr.Value{
					types.ObjectValueMust(getDeploymentActionGitDependencyAttrTypes(), map[string]attr.Value{
						"deployment_action_slug": types.StringValue("deploy-package"),
						"git_dependency_name":    types.StringValue("app-config"),
					}),
				}),
			}),
		}),
	}

	channel := expandChannel(t.Context(), model)
	require.Equal(t, []string{"refs/heads/main", "refs/tags/release-*"}, channel.GitReferenceRules)
	require.Len(t, channel.GitResourceRules, 1)
	assert.Equal(t, "ChannelGitResourceRules-1", channel.GitResourceRules[0].Id)
	assert.Equal(t, []string{"refs/heads/release/*"}, channel.GitResourceRules[0].Rules)
	require.Len(t, channel.GitResourceRules[0].GitDependencyActions, 1)
	assert.Equal(t, "deploy-package", channel.GitResourceRules[0].GitDependencyActions[0].DeploymentActionSlug)
	assert.Equal(t, "app-config", channel.GitResourceRules[0].GitDependencyActions[0].GitDependencyName)

	flattened := flattenChannel(t.Context(), channel, schemas.ChannelModel{
		GitReferenceRules: types.ListNull(types.StringType),
		GitResourceRules:  types.ListNull(types.ObjectType{AttrTypes: getChannelGitResourceRuleAttrTypes()}),
	})

	assert.Equal(t, 2, len(flattened.GitReferenceRules.Elements()))
	assert.Equal(t, 1, len(flattened.GitResourceRules.Elements()))
}

func TestAccChannelRuleRemoval(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := fmt.Sprintf("octopusdeploy_channel.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testChannelWithRule(localName, true),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					testChannelRuleCount(resourceName, 1),
				),
			},
			{
				Config: testChannelWithRule(localName, false),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "0"),
					testChannelRuleCount(resourceName, 0),
				),
			},
		},
	})
}

func TestAccChannelTenantTagRemoval(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := fmt.Sprintf("octopusdeploy_channel.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testChannelWithTenantTags(localName, true),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tenant_tags.#", "1"),
					testChannelTenantTagCount(resourceName, 1),
				),
			},
			{
				Config: testChannelWithTenantTags(localName, false),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tenant_tags.#", "0"),
					testChannelTenantTagCount(resourceName, 0),
				),
			},
		},
	})
}

func TestAccChannelGitReferenceRules(t *testing.T) {
	gitURL, gitUsername, gitPassword := testAccGitSettings()
	if gitURL == "" || gitUsername == "" || gitPassword == "" {
		t.Skip("Skipping Git reference rules test: GIT_URL, GIT_USERNAME, and GIT_PASSWORD or GIT_CREDENTIAL must be set")
	}

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	basePath := ".octopus/" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resourceName := fmt.Sprintf("octopusdeploy_channel.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testChannelWithGitReferenceRules(localName, basePath, gitURL, gitUsername, gitPassword, []string{"refs/heads/main", "refs/tags/release-*"}),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "git_reference_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "git_reference_rules.0", "refs/heads/main"),
					resource.TestCheckResourceAttr(resourceName, "git_reference_rules.1", "refs/tags/release-*"),
					testChannelGitReferenceRules(resourceName, []string{"refs/heads/main", "refs/tags/release-*"}),
				),
			},
			{
				Config: testChannelWithGitReferenceRules(localName, basePath, gitURL, gitUsername, gitPassword, []string{"refs/heads/release/*"}),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "git_reference_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "git_reference_rules.0", "refs/heads/release/*"),
					testChannelGitReferenceRules(resourceName, []string{"refs/heads/release/*"}),
				),
			},
			{
				Config: testChannelWithGitReferenceRules(localName, basePath, gitURL, gitUsername, gitPassword, nil),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "git_reference_rules.#", "0"),
					testChannelGitReferenceRules(resourceName, []string{}),
				),
			},
		},
	})
}

func TestAccChannelGitResourceRules(t *testing.T) {
	gitURL, gitUsername, gitPassword := testAccGitSettings()
	basePath, actionSlug, dependencyName, guidedFailureMode, ok := testAccGitResourceRuleSettings()
	if gitURL == "" || gitUsername == "" || gitPassword == "" || !ok {
		t.Skip("Skipping Git resource rules test: GIT_URL, GIT_USERNAME, GIT_PASSWORD or GIT_CREDENTIAL, GIT_RESOURCE_RULE_BASE_PATH, GIT_RESOURCE_RULE_ACTION_SLUG, and GIT_RESOURCE_RULE_DEPENDENCY_NAME must be set")
	}

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := fmt.Sprintf("octopusdeploy_channel.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testChannelWithGitResourceRules(localName, basePath, gitURL, gitUsername, gitPassword, guidedFailureMode, actionSlug, dependencyName, []string{"refs/heads/main"}),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "git_resource_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "git_resource_rules.0.rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "git_resource_rules.0.rules.0", "refs/heads/main"),
					resource.TestCheckResourceAttr(resourceName, "git_resource_rules.0.git_dependency_actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "git_resource_rules.0.git_dependency_actions.0.deployment_action_slug", actionSlug),
					resource.TestCheckResourceAttr(resourceName, "git_resource_rules.0.git_dependency_actions.0.git_dependency_name", dependencyName),
					testChannelGitResourceRules(resourceName, []string{"refs/heads/main"}, actionSlug, dependencyName),
				),
			},
			{
				Config: testChannelWithGitResourceRules(localName, basePath, gitURL, gitUsername, gitPassword, guidedFailureMode, actionSlug, dependencyName, []string{"refs/heads/release/*"}),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "git_resource_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "git_resource_rules.0.rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "git_resource_rules.0.rules.0", "refs/heads/release/*"),
					testChannelGitResourceRules(resourceName, []string{"refs/heads/release/*"}, actionSlug, dependencyName),
				),
			},
			{
				Config: testChannelWithGitResourceRules(localName, basePath, gitURL, gitUsername, gitPassword, guidedFailureMode, actionSlug, dependencyName, nil),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "git_resource_rules.#", "0"),
					testChannelGitResourceRuleCount(resourceName, 0),
				),
			},
		},
	})
}

func testChannelExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		channelID, err := getChannelID(s, resourceName)
		if err != nil {
			return err
		}

		if _, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), channelID); err != nil {
			return fmt.Errorf("channel %s not found", channelID)
		}

		return nil
	}
}

func testChannelDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_channel" {
			continue
		}

		if _, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("channel %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testChannelRuleCount(resourceName string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		channelID, err := getChannelID(s, resourceName)
		if err != nil {
			return err
		}

		channel, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), channelID)
		if err != nil {
			return fmt.Errorf("channel %s not found", channelID)
		}

		if len(channel.Rules) != expected {
			return fmt.Errorf("expected %d channel rules, got %d", expected, len(channel.Rules))
		}

		return nil
	}
}

func testChannelTenantTagCount(resourceName string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		channelID, err := getChannelID(s, resourceName)
		if err != nil {
			return err
		}

		channel, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), channelID)
		if err != nil {
			return fmt.Errorf("channel %s not found", channelID)
		}

		if len(channel.TenantTags) != expected {
			return fmt.Errorf("expected %d tenant tags, got %d", expected, len(channel.TenantTags))
		}

		return nil
	}
}

func testChannelGitReferenceRules(resourceName string, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		channelID, err := getChannelID(s, resourceName)
		if err != nil {
			return err
		}

		channel, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), channelID)
		if err != nil {
			return fmt.Errorf("channel %s not found", channelID)
		}

		return compareStringSlices(expected, channel.GitReferenceRules)
	}
}

func testChannelGitResourceRules(resourceName string, expectedRules []string, expectedActionSlug string, expectedDependencyName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		channelID, err := getChannelID(s, resourceName)
		if err != nil {
			return err
		}

		channel, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), channelID)
		if err != nil {
			return fmt.Errorf("channel %s not found", channelID)
		}

		if len(channel.GitResourceRules) != 1 {
			return fmt.Errorf("expected 1 Git resource rule, got %d", len(channel.GitResourceRules))
		}

		rule := channel.GitResourceRules[0]
		if err := compareStringSlices(expectedRules, rule.Rules); err != nil {
			return err
		}

		if len(rule.GitDependencyActions) != 1 {
			return fmt.Errorf("expected 1 Git dependency action, got %d", len(rule.GitDependencyActions))
		}

		action := rule.GitDependencyActions[0]
		if action.DeploymentActionSlug != expectedActionSlug {
			return fmt.Errorf("expected deployment action slug %q, got %q", expectedActionSlug, action.DeploymentActionSlug)
		}
		if action.GitDependencyName != expectedDependencyName {
			return fmt.Errorf("expected Git dependency name %q, got %q", expectedDependencyName, action.GitDependencyName)
		}

		return nil
	}
}

func testChannelGitResourceRuleCount(resourceName string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		channelID, err := getChannelID(s, resourceName)
		if err != nil {
			return err
		}

		channel, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), channelID)
		if err != nil {
			return fmt.Errorf("channel %s not found", channelID)
		}

		if len(channel.GitResourceRules) != expected {
			return fmt.Errorf("expected %d Git resource rules, got %d", expected, len(channel.GitResourceRules))
		}

		return nil
	}
}

func getChannelID(s *terraform.State, resourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("Not found: %s", resourceName)
	}

	return rs.Primary.ID, nil
}

func testChannelBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName string, channel channels.Channel) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_lifecycle" "%s" {
			name = "%s"
		}

		resource "octopusdeploy_project_group" "%s" {
			name = "%s"
		}

		resource "octopusdeploy_project" "%s" {
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			name             = "%s"
			project_group_id = octopusdeploy_project_group.%s.id
		}

		resource "octopusdeploy_channel" "%s" {
			name        = "%s"
			description = "%s"
			project_id  = octopusdeploy_project.%s.id
		}
	`,
		lifecycleLocalName, lifecycleName,
		projectGroupLocalName, projectGroupName,
		projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName,
		localName, channel.Name, channel.Description, projectLocalName,
	)
}

func testChannelWithRule(localName string, includeRule bool) string {
	ruleBlock := ""
	if includeRule {
		ruleBlock = `
			rule {
				version_range = "[1.0.0,2.0.0)"

				action_package {
					deployment_action = octopusdeploy_process_step.package_step.name
				}
			}
		`
	}

	return fmt.Sprintf(`
		data "octopusdeploy_lifecycles" "default" {
		  ids          = null
		  partial_name = "Default Lifecycle"
		  skip         = 0
		  take         = 1
		}

		data "octopusdeploy_feeds" "built_in_feed" {
		  feed_type    = "BuiltIn"
		  ids          = null
		  partial_name = ""
		  skip         = 0
		  take         = 1
		}

		resource "octopusdeploy_project_group" "%[1]s" {
		  name = "%[1]s"
		}

		resource "octopusdeploy_project" "%[1]s" {
		  lifecycle_id     = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		  name             = "%[1]s"
		  project_group_id = octopusdeploy_project_group.%[1]s.id
		}

		resource "octopusdeploy_process" "%[1]s" {
		  project_id = octopusdeploy_project.%[1]s.id
		}

		resource "octopusdeploy_process_step" "package_step" {
		  process_id = octopusdeploy_process.%[1]s.id
		  name       = "Package deployment"
		  type       = "Octopus.TentaclePackage"
		  properties = {
			"Octopus.Action.TargetRoles" = "Webserver"
		  }
		  primary_package = {
			feed_id    = data.octopusdeploy_feeds.built_in_feed.feeds[0].id
			package_id = "MyPackage"
		  }
		  execution_properties = {
			"Octopus.Action.RunOnServer" = "True"
		  }
		}

		resource "octopusdeploy_channel" "%[1]s" {
		  name       = "%[1]s"
		  project_id = octopusdeploy_project.%[1]s.id

		  %[2]s

		  depends_on = [octopusdeploy_process_step.package_step]
		}
	`, localName, ruleBlock)
}

func testChannelWithTenantTags(localName string, includeTenantTags bool) string {
	tenantTags := ""
	if includeTenantTags {
		tenantTags = "tenant_tags = [octopusdeploy_tag.channel_tag.canonical_tag_name]"
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

		resource "octopusdeploy_project" "%[1]s" {
		  lifecycle_id     = data.octopusdeploy_lifecycles.default.lifecycles[0].id
		  name             = "%[1]s"
		  project_group_id = octopusdeploy_project_group.%[1]s.id
		}

		resource "octopusdeploy_tag_set" "%[1]s" {
		  name = "%[1]s"
		}

		resource "octopusdeploy_tag" "channel_tag" {
		  name        = "%[1]s"
		  color       = "#6e6e6e"
		  description = "Channel tag"
		  tag_set_id  = octopusdeploy_tag_set.%[1]s.id
		}

		resource "octopusdeploy_channel" "%[1]s" {
		  name       = "%[1]s"
		  project_id = octopusdeploy_project.%[1]s.id

		  %[2]s
		}
	`, localName, tenantTags)
}

func testChannelWithGitReferenceRules(localName, basePath, gitURL, gitUsername, gitPassword string, gitReferenceRules []string) string {
	rules := ""
	if gitReferenceRules != nil {
		rules = fmt.Sprintf("git_reference_rules = %s", quoteStringList(gitReferenceRules))
	}

	return testAccCaCProjectConfig(localName, basePath, gitURL, gitUsername, gitPassword, "Off", false) + fmt.Sprintf(`
		resource "octopusdeploy_channel" "%[1]s" {
		  name       = "%[1]s"
		  project_id = octopusdeploy_project.%[1]s.id

		  %[2]s
		}
	`, localName, rules)
}

func testChannelWithGitResourceRules(localName, basePath, gitURL, gitUsername, gitPassword, guidedFailureMode, actionSlug, dependencyName string, gitResourceRules []string) string {
	ruleBlock := ""
	if gitResourceRules != nil {
		ruleBlock = fmt.Sprintf(`
		  git_resource_rules = [{
		    rules = %[1]s

		    git_dependency_actions = [{
		      deployment_action_slug = %[2]q
		      git_dependency_name    = %[3]q
		    }]
		  }]
		`, quoteStringList(gitResourceRules), actionSlug, dependencyName)
	}

	return testAccCaCProjectConfig(localName, basePath, gitURL, gitUsername, gitPassword, guidedFailureMode, false) + fmt.Sprintf(`
		resource "octopusdeploy_channel" "%[1]s" {
		  name       = "%[1]s"
		  project_id = octopusdeploy_project.%[1]s.id

		  %[2]s
		}
	`, localName, ruleBlock)
}

func quoteStringList(values []string) string {
	if len(values) == 0 {
		return "[]"
	}

	result := "["
	for i, value := range values {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%q", value)
	}
	return result + "]"
}

func compareStringSlices(expected []string, actual []string) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("expected %d values, got %d: %v", len(expected), len(actual), actual)
	}

	counts := map[string]int{}
	for _, value := range expected {
		counts[value]++
	}
	for _, value := range actual {
		counts[value]--
	}
	for value, count := range counts {
		if count != 0 {
			return fmt.Errorf("expected values %v, got %v; mismatch at %q", expected, actual, value)
		}
	}

	return nil
}
