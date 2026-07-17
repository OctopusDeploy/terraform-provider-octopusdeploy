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
