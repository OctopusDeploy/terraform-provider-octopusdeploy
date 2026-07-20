package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/approvalpolicies"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccApprovalPolicyBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_approval_policy." + localName

	policyName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedPolicyName := policyName + "-updated"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testApprovalPolicyDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testApprovalPolicyBasic(localName, policyName, 2),
				Check: resource.ComposeTestCheckFunc(
					testApprovalPolicyExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", policyName),
					resource.TestCheckResourceAttr(prefix, "minimum_approvers_required", "2"),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "false"),
					resource.TestCheckResourceAttr(prefix, "scoping_strategy", "Id"),
					resource.TestCheckResourceAttrSet(prefix, "approving_team_ids.0"),
					resource.TestCheckResourceAttr(prefix, "id_scopes.#", "1"),
					resource.TestCheckResourceAttrSet(prefix, "id_scopes.0.project_id"),
				),
			},
			{
				Config: testApprovalPolicyBasic(localName, updatedPolicyName, 1),
				Check: resource.ComposeTestCheckFunc(
					testApprovalPolicyExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", updatedPolicyName),
					resource.TestCheckResourceAttr(prefix, "minimum_approvers_required", "1"),
					resource.TestCheckResourceAttr(prefix, "scoping_strategy", "Id"),
					resource.TestCheckResourceAttr(prefix, "is_disabled", "false"),
				),
			},
		},
	})
}

func testApprovalPolicyBasic(localName string, policyName string, minimumApproversRequired int) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project_group" "%s" {
		name        = "Test Project Group %s"
		description = "Project group for approval policy acceptance test"
	}

	resource "octopusdeploy_lifecycle" "%s" {
		name = "Test Lifecycle %s"
	}

	resource "octopusdeploy_project" "%s" {
		name              = "Test Project %s"
		project_group_id  = octopusdeploy_project_group.%s.id
		lifecycle_id      = octopusdeploy_lifecycle.%s.id
	}

	resource "octopusdeploy_environment" "%s" {
		name        = "Test Environment %s"
		description = "Environment for approval policy acceptance test"
	}

	resource "octopusdeploy_team" "%s" {
		name        = "Test Team %s"
		description = "Team for approval policy acceptance test"
	}

	resource "octopusdeploy_approval_policy" "%s" {
		name                        = "%s"
		minimum_approvers_required  = %d
		scoping_strategy            = "Id"
		approving_team_ids          = [octopusdeploy_team.%s.id]

		id_scopes = [
			{
				project_id      = octopusdeploy_project.%s.id
				environment_ids = [octopusdeploy_environment.%s.id]
			}
		]
	}
	`,
		localName, localName,
		localName, localName,
		localName, localName, localName, localName,
		localName, localName,
		localName, localName,
		localName, policyName, minimumApproversRequired, localName,
		localName, localName,
	)
}

// TestAccApprovalPolicyScopeValidation verifies the config-time validation that
// the provided scope matches the scoping strategy and that both scopes cannot be
// set at once. All steps are plan-only and expect an error, so nothing is created.
func TestAccApprovalPolicyScopeValidation(t *testing.T) {
	policyName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testApprovalPolicyBothScopes(policyName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`[Oo]nly one of`),
			},
			{
				Config:      testApprovalPolicyScopeMismatch(policyName, "Id", true),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`cannot be set when scoping_strategy is`),
			},
			{
				Config:      testApprovalPolicyScopeMismatch(policyName, "Tag", false),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`cannot be set when scoping_strategy is`),
			},
		},
	})
}

func testApprovalPolicyBothScopes(policyName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_approval_policy" "test" {
		name               = "%s"
		approving_team_ids = ["Teams-1"]

		tag_scopes = [
			{
				project_tags     = ["projects/example"]
				environment_tags = ["environments/example"]
			}
		]

		id_scopes = [
			{
				project_id      = "Projects-1"
				environment_ids = ["Environments-1"]
			}
		]
	}
	`, policyName)
}

// testApprovalPolicyScopeMismatch builds a config whose scope does not match the
// scoping strategy: strategy "Id" with tag_scopes (useTagScope=true) or strategy
// "Tag" with id_scopes (useTagScope=false).
func testApprovalPolicyScopeMismatch(policyName string, strategy string, useTagScope bool) string {
	scopeBlock := `
		id_scopes = [
			{
				project_id      = "Projects-1"
				environment_ids = ["Environments-1"]
			}
		]`
	if useTagScope {
		scopeBlock = `
		tag_scopes = [
			{
				project_tags     = ["projects/example"]
				environment_tags = ["environments/example"]
			}
		]`
	}

	return fmt.Sprintf(`
	resource "octopusdeploy_approval_policy" "test" {
		name               = "%s"
		scoping_strategy   = "%s"
		approving_team_ids = ["Teams-1"]
%s
	}
	`, policyName, strategy, scopeBlock)
}

// TestAccApprovalPolicyTagScopesResolveNames verifies that tag_scopes hold stable
// tag IDs regardless of whether they are configured by ID or by canonical name.
// Step 1 references same-apply tags by their id attribute (unknown at plan,
// resolved at apply). Step 2 references the now-existing tags by canonical name;
// the provider resolves those names to the same IDs, so there is no diff. The
// TypeSetElemAttrPair checks assert the stored values equal the tags' IDs, which
// only holds if resolution actually happened.
func TestAccApprovalPolicyTagScopesResolveNames(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	policyName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_approval_policy." + localName
	projectTagResource := "octopusdeploy_tag." + localName + "_project"
	envTagResource := "octopusdeploy_tag." + localName + "_env"

	checks := resource.ComposeTestCheckFunc(
		testApprovalPolicyExists(prefix),
		resource.TestCheckResourceAttr(prefix, "scoping_strategy", "Tag"),
		resource.TestCheckResourceAttr(prefix, "tag_scopes.#", "1"),
		resource.TestCheckResourceAttr(prefix, "tag_scopes.0.project_tags.#", "1"),
		resource.TestCheckResourceAttr(prefix, "tag_scopes.0.environment_tags.#", "1"),
		resource.TestCheckTypeSetElemAttrPair(prefix, "tag_scopes.0.project_tags.*", projectTagResource, "id"),
		resource.TestCheckTypeSetElemAttrPair(prefix, "tag_scopes.0.environment_tags.*", envTagResource, "id"),
	)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testApprovalPolicyDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testApprovalPolicyTagScopes(localName, policyName, true),
				Check:  checks,
			},
			{
				Config: testApprovalPolicyTagScopes(localName, policyName, false),
				Check:  checks,
			},
		},
	})
}

// testApprovalPolicyTagScopes builds a Tag-scoped approval policy. When byID is
// true the scopes reference the tags by their id attribute; otherwise they use
// canonical tag names, which the provider must resolve to the tags' IDs.
func testApprovalPolicyTagScopes(localName string, policyName string, byID bool) string {
	projectTagRef := fmt.Sprintf(`"${octopusdeploy_tag_set.%[1]s.name}/${octopusdeploy_tag.%[1]s_project.name}"`, localName)
	envTagRef := fmt.Sprintf(`"${octopusdeploy_tag_set.%[1]s.name}/${octopusdeploy_tag.%[1]s_env.name}"`, localName)
	if byID {
		projectTagRef = fmt.Sprintf("octopusdeploy_tag.%s_project.id", localName)
		envTagRef = fmt.Sprintf("octopusdeploy_tag.%s_env.id", localName)
	}

	return fmt.Sprintf(`
	resource "octopusdeploy_tag_set" "%[1]s" {
		name        = "TagSet %[1]s"
		description = "Tag set for approval policy acceptance test"
		scopes      = ["Project", "Environment"]
	}

	resource "octopusdeploy_tag" "%[1]s_project" {
		name        = "proj-%[1]s"
		color       = "#111111"
		tag_set_id  = octopusdeploy_tag_set.%[1]s.id
	}

	resource "octopusdeploy_tag" "%[1]s_env" {
		name        = "env-%[1]s"
		color       = "#222222"
		tag_set_id  = octopusdeploy_tag_set.%[1]s.id
	}

	resource "octopusdeploy_team" "%[1]s" {
		name        = "Test Team %[1]s"
		description = "Team for approval policy acceptance test"
	}

	resource "octopusdeploy_approval_policy" "%[1]s" {
		name                       = "%[2]s"
		scoping_strategy           = "Tag"
		minimum_approvers_required = 2
		approving_team_ids         = [octopusdeploy_team.%[1]s.id]

		tag_scopes = [
			{
				project_tags     = [%[3]s]
				environment_tags = [%[4]s]
			}
		]
	}
	`, localName, policyName, projectTagRef, envTagRef)
}

func testApprovalPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policyResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		_, err := approvalpolicies.GetByID(octoClient, octoClient.GetSpaceID(), policyResource.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve approval policy (%s): %s", policyResource.Primary.ID, err.Error())
		}

		return nil
	}
}

func testApprovalPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_approval_policy" {
			continue
		}

		approvalPolicy, err := approvalpolicies.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID)
		if err == nil {
			if approvalPolicy != nil {
				return fmt.Errorf("approval policy (%s) still exists", approvalPolicy.Name)
			}
		}
	}

	return nil
}
