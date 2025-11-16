package octopusdeploy_framework

import (
	"fmt"
	"strings"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccTenantProjectVariableBasic tests V1 API (environment_id, no scope block)
// This test works on all Octopus server versions, including those without V2 feature flag support
func TestAccTenantProjectVariableBasic(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	primaryEnvironmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	primaryEnvironmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	primaryLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	primaryValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secondaryEnvironmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secondaryEnvironmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secondaryLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	secondaryValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	primaryResourceName := "octopusdeploy_tenant_project_variable." + primaryLocalName
	secondaryResourceName := "octopusdeploy_tenant_project_variable." + secondaryLocalName

	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTenantProjectVariableCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExists(primaryResourceName),
					testTenantProjectVariableExists(secondaryResourceName),
					resource.TestCheckResourceAttr(primaryResourceName, "value", primaryValue),
					resource.TestCheckResourceAttr(secondaryResourceName, "value", secondaryValue),
				),
				Config: testAccTenantProjectVariable(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, primaryEnvironmentLocalName, primaryEnvironmentName, secondaryEnvironmentLocalName, secondaryEnvironmentName, tenantLocalName, tenantName, tenantDescription, primaryLocalName, primaryValue, secondaryLocalName, secondaryValue),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExists(primaryResourceName),
					testTenantProjectVariableExists(secondaryResourceName),
					resource.TestCheckResourceAttr(primaryResourceName, "value", primaryValue),
					resource.TestCheckResourceAttr(secondaryResourceName, "value", newValue),
				),
				Config: testAccTenantProjectVariable(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, primaryEnvironmentLocalName, primaryEnvironmentName, secondaryEnvironmentLocalName, secondaryEnvironmentName, tenantLocalName, tenantName, tenantDescription, primaryLocalName, primaryValue, secondaryLocalName, newValue),
			},
		},
	})
}

func testAccTenantProjectVariable(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, projectDescription string, primaryEnvironmentLocalName string, primaryEnvironmentName string, secondaryEnvironmentLocalName string, secondaryEnvironmentName string, tenantLocalName string, tenantName string, tenantDescription string, primaryLocalName string, primaryValue string, secondaryLocalName string, secondaryValue string) string {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	var tfconfig = testAccLifecycle(lifecycleLocalName, lifecycleName) + "\n" +
		testAccProjectGroup(projectGroupLocalName, projectGroupName) + "\n" +
		testAccProjectWithTemplate(projectLocalName, projectName, lifecycleLocalName, projectGroupLocalName) + "\n" +
		testAccEnvironment(primaryEnvironmentLocalName, primaryEnvironmentName, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure) + "\n" +
		testAccEnvironment(secondaryEnvironmentLocalName, secondaryEnvironmentName, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure) + "\n" +
		testAccTenantWithProjectEnvironment(tenantLocalName, tenantName, projectLocalName, primaryEnvironmentLocalName, secondaryEnvironmentLocalName) + "\n" +
		testTenantProjectVariable(primaryLocalName, primaryEnvironmentLocalName, projectLocalName, tenantLocalName, projectLocalName, primaryValue) + "\n" +
		testTenantProjectVariable(secondaryLocalName, secondaryEnvironmentLocalName, projectLocalName, tenantLocalName, projectLocalName, secondaryValue)
	return tfconfig
}

func testAccProjectWithTemplate(localName string, name string, lifecycleLocalName string, projectGroupLocalName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project" "%s" {
		lifecycle_id     = octopusdeploy_lifecycle.%s.id
		name             = "%s"
		project_group_id = octopusdeploy_project_group.%s.id

		template {
			name  = "project variable template name"
			label = "project variable template label"

			display_settings = {
				"Octopus.ControlType" = "Sensitive"
			}
		}
	}`, localName, lifecycleLocalName, name, projectGroupLocalName)
}

func testAccProjectGroup(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_project_group" "%s" {
		name = "%s"
	}`, localName, name)
}

func testAccTenantWithProjectEnvironment(localName string, name string, projectLocalName string, primaryEnvironmentLocalName string, secondaryEnvironmentLocalName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_tenant" "%s" {
		name = "%s"
	}

	resource "octopusdeploy_tenant_project" "project_environment" {
		tenant_id = octopusdeploy_tenant.%s.id
		project_id   = "${octopusdeploy_project.%s.id}"
		environment_ids = [octopusdeploy_environment.%s.id, octopusdeploy_environment.%s.id]
	}`, localName, name, localName, projectLocalName, primaryEnvironmentLocalName, secondaryEnvironmentLocalName)
}

func testTenantProjectVariable(localName string, environmentLocalName string, projectLocalName string, tenantLocalName string, templateLocalName string, value string) string {
	return fmt.Sprintf(`resource "octopusdeploy_tenant_project_variable" "%s" {
		environment_id = octopusdeploy_environment.%s.id
		project_id     = octopusdeploy_project.%s.id
		tenant_id      = octopusdeploy_tenant.%s.id
		template_id    = octopusdeploy_project.%s.template[0].id
		value          = "%s"
 		depends_on     = [
            octopusdeploy_project.%s,
            octopusdeploy_tenant_project.project_environment
        ]
	}`, localName, environmentLocalName, projectLocalName, tenantLocalName, templateLocalName, value, projectLocalName)
}

func testTenantProjectVariableExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var resourceState *terraform.ResourceState
		for _, r := range s.RootModule().Resources {
			if r.Type == "octopusdeploy_tenant_project_variable" {
				resourceState = r
				break
			}
		}

		if resourceState == nil {
			return fmt.Errorf("tenant project variable resource not found")
		}

		if len(resourceState.Primary.ID) == 0 {
			return fmt.Errorf("tenant project variable ID is not set")
		}

		// Check if this is a V2 ID (no colons) or V1 ID (has colons)
		if !strings.Contains(resourceState.Primary.ID, ":") {
			// V2 API - use real ID
			tenantID := resourceState.Primary.Attributes["tenant_id"]
			spaceID := resourceState.Primary.Attributes["space_id"]

			client := octoClient
			query := variables.GetTenantProjectVariablesQuery{
				TenantID:                tenantID,
				SpaceID:                 spaceID,
				IncludeMissingVariables: false,
			}

			getResp, err := tenants.GetProjectVariables(client, query)
			if err != nil {
				return fmt.Errorf("Error retrieving tenant project variables: %s", err.Error())
			}

			// Search for the variable by ID
			for _, v := range getResp.Variables {
				if v.GetID() == resourceState.Primary.ID {
					return nil
				}
			}

			return fmt.Errorf("Tenant project variable with ID %s not found via V2 API", resourceState.Primary.ID)
		}

		// V1 API - use composite ID
		environmentID := resourceState.Primary.Attributes["environment_id"]
		projectID := resourceState.Primary.Attributes["project_id"]
		templateID := resourceState.Primary.Attributes["template_id"]
		tenantID := resourceState.Primary.Attributes["tenant_id"]

		tenant, err := octoClient.Tenants.GetByID(tenantID)
		if err != nil {
			return err
		}

		tenantVariables, err := octoClient.Tenants.GetVariables(tenant)
		if err != nil {
			return err
		}

		if projectVariable, ok := tenantVariables.ProjectVariables[projectID]; ok {
			if _, ok := projectVariable.Variables[environmentID]; ok {
				if _, ok := projectVariable.Variables[environmentID][templateID]; ok {
					return nil
				}
			}
		}

		return fmt.Errorf("tenant project variable not found")
	}
}

func testAccTenantProjectVariableCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tenant_project_variable" {
			continue
		}

		// Check if this is a V2 ID (no colons) or V1 ID (has colons)
		if !strings.Contains(rs.Primary.ID, ":") {
			// V2 API - use real ID
			tenantID := rs.Primary.Attributes["tenant_id"]
			spaceID := rs.Primary.Attributes["space_id"]

			client := octoClient
			query := variables.GetTenantProjectVariablesQuery{
				TenantID:                tenantID,
				SpaceID:                 spaceID,
				IncludeMissingVariables: false,
			}

			getResp, err := tenants.GetProjectVariables(client, query)
			if err != nil {
				// If we can't get the variables, assume they're gone
				return nil
			}

			// Search for the variable by ID
			for _, v := range getResp.Variables {
				if v.GetID() == rs.Primary.ID {
					return fmt.Errorf("Tenant project variable (%s) still exists", rs.Primary.ID)
				}
			}

			continue
		}

		// V1 API - use composite ID
		importStrings := strings.Split(rs.Primary.ID, ":")
		if len(importStrings) != 4 {
			return fmt.Errorf("octopusdeploy_tenant_project_variable import must be in the form of TenantID:ProjectID:EnvironmentID:TemplateID (e.g. Tenants-123:Projects-456:Environments-789:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c")
		}

		tenantID := importStrings[0]
		projectID := importStrings[1]
		environmentID := importStrings[2]
		templateID := importStrings[3]

		tenant, err := octoClient.Tenants.GetByID(tenantID)
		if err != nil {
			return nil
		}

		tenantVariables, err := octoClient.Tenants.GetVariables(tenant)
		if err != nil {
			return nil
		}

		if projectVariable, ok := tenantVariables.ProjectVariables[projectID]; ok {
			if _, ok := projectVariable.Variables[environmentID]; ok {
				if _, ok := projectVariable.Variables[environmentID][templateID]; ok {
					return fmt.Errorf("tenant project variable (%s) still exists", rs.Primary.ID)
				}
			}
		}
	}

	return nil
}

// TestAccTenantProjectVariableMigration tests transitioning from environment_id to scope block
// When V2 is available, variables are created with V2 API regardless of whether environment_id or scope is used
func TestAccTenantProjectVariableMigration(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	env1LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	env1Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	env2LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	env2Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	variableLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := "octopusdeploy_tenant_project_variable." + variableLocalName

	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	finalValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTenantProjectVariableCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with environment_id (auto-converted to scope when V2 available)
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExists(""),
					resource.TestCheckResourceAttr(resourceName, "value", value),
				),
				Config: testAccTenantProjectVariableMigrationV1(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, variableLocalName, value),
			},
			{
				// Step 2: Replace environment_id with scope block containing multiple environments
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExists(""),
					resource.TestCheckResourceAttr(resourceName, "value", newValue),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.environment_ids.#", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "environment_id"),
				),
				Config: testAccTenantProjectVariableMigrationV2(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, variableLocalName, newValue),
			},
			{
				// Step 3: Update value again
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExists(""),
					resource.TestCheckResourceAttr(resourceName, "value", finalValue),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.environment_ids.#", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "environment_id"),
				),
				Config: testAccTenantProjectVariableMigrationV2(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, variableLocalName, finalValue),
			},
		},
	})
}

func testAccTenantProjectVariableMigrationV1(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, localName, value string) string {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(1, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		testAccProjectGroup(projectGroupLocalName, projectGroupName)+"\n"+
		testAccEnvironment(env1LocalName, env1Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(env2LocalName, env2Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
        resource "octopusdeploy_project" "%[1]s" {
            lifecycle_id     = octopusdeploy_lifecycle.%[2]s.id
            name             = "%[3]s"
            project_group_id = octopusdeploy_project_group.%[4]s.id

            template {
                default_value = "Default Value"
                help_text     = "This is help text"
                label         = "Test Label"
                name          = "Test Template Migration"

                display_settings = {
                    "Octopus.ControlType" = "Sensitive"
                }
            }
        }

        resource "octopusdeploy_tenant" "%[5]s" {
            name = "%[6]s"
        }

        resource "octopusdeploy_tenant_project" "project_environment_migration" {
            tenant_id        = octopusdeploy_tenant.%[5]s.id
            project_id       = octopusdeploy_project.%[1]s.id
            environment_ids  = [octopusdeploy_environment.%[7]s.id, octopusdeploy_environment.%[8]s.id]
        }

        resource "octopusdeploy_tenant_project_variable" "%[9]s" {
            tenant_id      = octopusdeploy_tenant.%[5]s.id
            project_id     = octopusdeploy_project.%[1]s.id
            environment_id = octopusdeploy_environment.%[7]s.id
            template_id    = octopusdeploy_project.%[1]s.template[0].id
            value          = "%[10]s"

            depends_on = [octopusdeploy_tenant_project.project_environment_migration]
        }`, projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName, tenantLocalName, tenantName, env1LocalName, env2LocalName, localName, value)
}

func testAccTenantProjectVariableMigrationV2(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, localName, value string) string {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(1, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		testAccProjectGroup(projectGroupLocalName, projectGroupName)+"\n"+
		testAccEnvironment(env1LocalName, env1Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(env2LocalName, env2Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
        resource "octopusdeploy_project" "%[1]s" {
            lifecycle_id     = octopusdeploy_lifecycle.%[2]s.id
            name             = "%[3]s"
            project_group_id = octopusdeploy_project_group.%[4]s.id

            template {
                default_value = "Default Value"
                help_text     = "This is help text"
                label         = "Test Label"
                name          = "Test Template Migration"

                display_settings = {
                    "Octopus.ControlType" = "Sensitive"
                }
            }
        }

        resource "octopusdeploy_tenant" "%[5]s" {
            name = "%[6]s"
        }

        resource "octopusdeploy_tenant_project" "project_environment_migration" {
            tenant_id        = octopusdeploy_tenant.%[5]s.id
            project_id       = octopusdeploy_project.%[1]s.id
            environment_ids  = [octopusdeploy_environment.%[7]s.id, octopusdeploy_environment.%[8]s.id]
        }

        resource "octopusdeploy_tenant_project_variable" "%[9]s" {
            tenant_id   = octopusdeploy_tenant.%[5]s.id
            project_id  = octopusdeploy_project.%[1]s.id
            template_id = octopusdeploy_project.%[1]s.template[0].id
            value       = "%[10]s"

            scope {
                environment_ids = [octopusdeploy_environment.%[7]s.id, octopusdeploy_environment.%[8]s.id]
            }

            depends_on = [octopusdeploy_tenant_project.project_environment_migration]
        }`, projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName, tenantLocalName, tenantName, env1LocalName, env2LocalName, localName, value)
}

// TestAccTenantProjectVariableWithScope tests V2 API with multi-environment scoping
func TestAccTenantProjectVariableWithScope(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	env1LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	env1Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	env2LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	env2Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	variableLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := "octopusdeploy_tenant_project_variable." + variableLocalName

	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTenantProjectVariableCheckDestroyV2,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Create with scope
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExistsV2(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.environment_ids.#", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "environment_id"),
				),
				Config: testAccTenantProjectVariableWithScope(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, variableLocalName, value),
			},
			{
				// Update value
				Check: resource.ComposeTestCheckFunc(
					testTenantProjectVariableExistsV2(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", newValue),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.environment_ids.#", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "environment_id"),
				),
				Config: testAccTenantProjectVariableWithScope(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, variableLocalName, newValue),
			},
		},
	})
}

func testAccTenantProjectVariableWithScope(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, variableLocalName, value string) string {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(1, 10)
	useGuidedFailure := false

	return testAccLifecycle(lifecycleLocalName, lifecycleName) + "\n" +
		testAccProjectGroup(projectGroupLocalName, projectGroupName) + "\n" +
		testAccProjectWithTemplate(projectLocalName, projectName, lifecycleLocalName, projectGroupLocalName) + "\n" +
		testAccEnvironment(env1LocalName, env1Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure) + "\n" +
		testAccEnvironment(env2LocalName, env2Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure) + "\n" +
		testAccTenantWithProjectEnvironment(tenantLocalName, tenantName, projectLocalName, env1LocalName, env2LocalName) + "\n" +
		testTenantProjectVariableWithScope(variableLocalName, projectLocalName, tenantLocalName, env1LocalName, env2LocalName, value)
}

func testTenantProjectVariableWithScope(localName, projectLocalName, tenantLocalName, env1LocalName, env2LocalName, value string) string {
	return fmt.Sprintf(`resource "octopusdeploy_tenant_project_variable" "%s" {
		project_id  = octopusdeploy_project.%s.id
		tenant_id   = octopusdeploy_tenant.%s.id
		template_id = octopusdeploy_project.%s.template[0].id
		value       = "%s"

		scope {
			environment_ids = [octopusdeploy_environment.%s.id, octopusdeploy_environment.%s.id]
		}

		depends_on = [
			octopusdeploy_project.%s,
			octopusdeploy_tenant_project.project_environment
		]
	}`, localName, projectLocalName, tenantLocalName, projectLocalName, value, env1LocalName, env2LocalName, projectLocalName)
}

// testTenantProjectVariableExistsV2 checks if a V2 tenant project variable exists (uses real ID, not composite)
func testTenantProjectVariableExistsV2(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if len(rs.Primary.ID) == 0 {
			return fmt.Errorf("Tenant project variable ID is not set")
		}

		// V2 uses real IDs (e.g., TenantVariables-123), not composite IDs with colons
		if strings.Contains(rs.Primary.ID, ":") {
			return fmt.Errorf("Expected V2 ID (e.g., TenantVariables-123) but got V1 composite ID: %s", rs.Primary.ID)
		}

		// Use V2 API to verify the variable exists
		tenantID := rs.Primary.Attributes["tenant_id"]
		spaceID := rs.Primary.Attributes["space_id"]

		client := octoClient
		query := variables.GetTenantProjectVariablesQuery{
			TenantID:                tenantID,
			SpaceID:                 spaceID,
			IncludeMissingVariables: false,
		}

		getResp, err := tenants.GetProjectVariables(client, query)
		if err != nil {
			return fmt.Errorf("Error retrieving tenant project variables: %s", err.Error())
		}

		// Search for the variable by ID
		for _, v := range getResp.Variables {
			if v.GetID() == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("Tenant project variable with ID %s not found via V2 API", rs.Primary.ID)
	}
}

// testAccTenantProjectVariableCheckDestroyV2 checks that V2 tenant project variables are destroyed
func testAccTenantProjectVariableCheckDestroyV2(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tenant_project_variable" {
			continue
		}

		// V2 uses real IDs, not composite IDs
		if strings.Contains(rs.Primary.ID, ":") {
			// This is a V1 ID, use V1 destroy check
			continue
		}

		tenantID := rs.Primary.Attributes["tenant_id"]
		spaceID := rs.Primary.Attributes["space_id"]

		client := octoClient
		query := variables.GetTenantProjectVariablesQuery{
			TenantID:                tenantID,
			SpaceID:                 spaceID,
			IncludeMissingVariables: false,
		}

		getResp, err := tenants.GetProjectVariables(client, query)
		if err != nil {
			// If we can't get the variables, assume they're gone
			return nil
		}

		// Search for the variable by ID
		for _, v := range getResp.Variables {
			if v.GetID() == rs.Primary.ID {
				return fmt.Errorf("Tenant project variable (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
