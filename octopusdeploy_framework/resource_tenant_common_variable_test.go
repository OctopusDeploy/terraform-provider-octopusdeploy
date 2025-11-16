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

	internalTest "github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
)

// TestAccTenantCommonVariableBasic tests V1 API
func TestAccTenantCommonVariableBasic(t *testing.T) {
	//SkipCI(t, "A managed resource \"octopusdeploy_project_group\" \"ewtxiwplhaenzmhpaqyx\" has\n        not been declared in the root module.")
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantVariablesLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := "octopusdeploy_tenant_common_variable." + tenantVariablesLocalName

	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccTenantCommonVariableCheckDestroy,
		PreCheck:     func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactoriesWithFeatureToggleOverrides(map[string]bool{
			"CommonVariableScopingFeatureToggle": false,
		}),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", value),
				),
				Config: testAccTenantCommonVariableBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, tenantLocalName, tenantName, tenantDescription, tenantVariablesLocalName, value),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", newValue),
				),
				Config: testAccTenantCommonVariableBasic(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, projectDescription, environmentLocalName, environmentName, tenantLocalName, tenantName, tenantDescription, tenantVariablesLocalName, newValue),
			},
		},
	})
}

func testAccTenantCommonVariableBasic(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, projectDescription string, environmentLocalName string, environmentName string, tenantLocalName string, tenantName string, tenantDescription string, localName string, value string) string {
	projectGroup := internalTest.NewProjectGroupTestOptions()
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(1, 10)
	useGuidedFailure := false
	projectGroup.LocalName = projectGroupLocalName

	var tfConfig = fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		internalTest.ProjectGroupConfiguration(projectGroup)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
        resource "octopusdeploy_library_variable_set" "test-library-variable-set" {
            name = "test"

            template {
                default_value = "Default Value???"
                help_text     = "This is the help text"
                label         = "Test Label"
                name          = "Test Template"

                display_settings = {
                    "Octopus.ControlType" = "Sensitive"
                }
            }
        }

        resource "octopusdeploy_project" "%[1]s" {
            included_library_variable_sets = [octopusdeploy_library_variable_set.test-library-variable-set.id]
            lifecycle_id                   = octopusdeploy_lifecycle.%[2]s.id
            name                           = "%[3]s"
            project_group_id               = octopusdeploy_project_group.%[4]s.id
            depends_on                     = [octopusdeploy_library_variable_set.test-library-variable-set]
        }

        resource "octopusdeploy_tenant" "%[5]s" {
            name = "%[6]s"
        }

        resource "octopusdeploy_tenant_project" "project_environment" {
            tenant_id        = octopusdeploy_tenant.%[5]s.id
            project_id       = octopusdeploy_project.%[1]s.id
            environment_ids  = [octopusdeploy_environment.%[7]s.id]
            depends_on       = [octopusdeploy_project.%[1]s, octopusdeploy_tenant.%[5]s, octopusdeploy_environment.%[7]s]
        }

        resource "octopusdeploy_tenant_common_variable" "%[8]s" {
            library_variable_set_id = octopusdeploy_library_variable_set.test-library-variable-set.id
            template_id             = octopusdeploy_library_variable_set.test-library-variable-set.template[0].id
            tenant_id               = octopusdeploy_tenant.%[5]s.id
            value                   = "%[9]s"
            depends_on              = [octopusdeploy_library_variable_set.test-library-variable-set, octopusdeploy_tenant_project.project_environment]
        }`, projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName, tenantLocalName, tenantName, environmentLocalName, localName, value)
	return tfConfig
}

func testTenantCommonVariableExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if len(rs.Primary.ID) == 0 {
			return fmt.Errorf("Library variable ID is not set")
		}

		// Check if this is a V2 ID (no colons) or V1 ID (has colons)
		if !strings.Contains(rs.Primary.ID, ":") {
			tenantID := rs.Primary.Attributes["tenant_id"]
			spaceID := rs.Primary.Attributes["space_id"]

			client := octoClient
			query := variables.GetTenantCommonVariablesQuery{
				TenantID:                tenantID,
				SpaceID:                 spaceID,
				IncludeMissingVariables: false,
			}

			getResp, err := tenants.GetCommonVariables(client, query)
			if err != nil {
				return fmt.Errorf("Error retrieving tenant common variables: %s", err.Error())
			}

			// Search for the variable by ID
			for _, v := range getResp.Variables {
				if v.GetID() == rs.Primary.ID {
					return nil
				}
			}

			return fmt.Errorf("Tenant common variable with ID %s not found via V2 API", rs.Primary.ID)
		}

		// V1 API - use composite ID
		importStrings := strings.Split(rs.Primary.ID, ":")
		if len(importStrings) != 3 {
			return fmt.Errorf("octopusdeploy_tenant_common_variable import must be in the form of TenantID:LibraryVariableSetID:VariableID (e.g. Tenants-123:LibraryVariableSets-456:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c")
		}

		tenantID := importStrings[0]
		libraryVariableSetID := importStrings[1]
		templateID := importStrings[2]

		tenant, err := octoClient.Tenants.GetByID(tenantID)
		if err != nil {
			return err
		}

		tenantVariables, err := octoClient.Tenants.GetVariables(tenant)
		if err != nil {
			return err
		}

		if libraryVariable, ok := tenantVariables.LibraryVariables[libraryVariableSetID]; ok {
			if _, ok := libraryVariable.Variables[templateID]; ok {
				return nil
			}
		}

		return fmt.Errorf("tenant common variable not found")
	}
}

func testAccTenantCommonVariableCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tenant_common_variable" {
			continue
		}

		// Check if this is a V2 ID (no colons) or V1 ID (has colons)
		if !strings.Contains(rs.Primary.ID, ":") {
			// V2 API - use real ID
			tenantID := rs.Primary.Attributes["tenant_id"]
			spaceID := rs.Primary.Attributes["space_id"]

			client := octoClient
			query := variables.GetTenantCommonVariablesQuery{
				TenantID:                tenantID,
				SpaceID:                 spaceID,
				IncludeMissingVariables: false,
			}

			getResp, err := tenants.GetCommonVariables(client, query)
			if err != nil {
				// If we can't get the variables, assume they're gone
				return nil
			}

			// Search for the variable by ID
			for _, v := range getResp.Variables {
				if v.GetID() == rs.Primary.ID {
					return fmt.Errorf("Tenant common variable (%s) still exists", rs.Primary.ID)
				}
			}

			continue
		}

		// V1 API - use composite ID
		importStrings := strings.Split(rs.Primary.ID, ":")
		if len(importStrings) != 3 {
			return fmt.Errorf("octopusdeploy_tenant_common_variable import must be in the form of TenantID:LibraryVariableSetID:VariableID (e.g. Tenants-123:LibraryVariableSets-456:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c")
		}

		tenantID := importStrings[0]
		libraryVariableSetID := importStrings[1]
		templateID := importStrings[2]

		tenant, err := octoClient.Tenants.GetByID(tenantID)
		if err != nil {
			return nil
		}

		tenantVariables, err := octoClient.Tenants.GetVariables(tenant)
		if err != nil {
			return nil
		}

		if libraryVariable, ok := tenantVariables.LibraryVariables[libraryVariableSetID]; ok {
			if _, ok := libraryVariable.Variables[templateID]; ok {
				return fmt.Errorf("tenant common variable (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

// TestAccTenantCommonVariableMigration tests migration from V1 to V2 API
func TestAccTenantCommonVariableMigration(t *testing.T) {
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
	tenantVariablesLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := "octopusdeploy_tenant_common_variable." + tenantVariablesLocalName

	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	finalValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccTenantCommonVariableCheckDestroy,
		PreCheck:     func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				// Step 1: Create with V1 API (force V2 API off)
				ProtoV6ProviderFactories: ProtoV6ProviderFactoriesWithFeatureToggleOverrides(map[string]bool{
					"CommonVariableScopingFeatureToggle": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckNoResourceAttr(resourceName, "scope.#"),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[resourceName]
						// Verify it's a V1 composite ID (has colons)
						if !strings.Contains(rs.Primary.ID, ":") {
							return fmt.Errorf("Expected V1 composite ID with colons, got: %s", rs.Primary.ID)
						}
						return nil
					},
				),
				Config: testAccTenantCommonVariableMigrationV1(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, tenantVariablesLocalName, value),
			},
			{
				// Step 2: Migrate to V2 API by adding scope block (enable V2 API)
				ProtoV6ProviderFactories: ProtoV6ProviderFactoriesWithFeatureToggleOverrides(map[string]bool{
					"CommonVariableScopingFeatureToggle": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", newValue),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.environment_ids.#", "2"),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[resourceName]
						// Verify it's a V2 real ID (no colons)
						if strings.Contains(rs.Primary.ID, ":") {
							return fmt.Errorf("Expected V2 real ID without colons, got: %s", rs.Primary.ID)
						}
						return nil
					},
				),
				Config: testAccTenantCommonVariableMigrationV2(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, tenantVariablesLocalName, newValue),
			},
			{
				// Step 3: Update value again to verify it works after migration
				ProtoV6ProviderFactories: ProtoV6ProviderFactoriesWithFeatureToggleOverrides(map[string]bool{
					"CommonVariableScopingFeatureToggle": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", finalValue),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.environment_ids.#", "2"),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[resourceName]
						// Verify still using V2 ID
						if strings.Contains(rs.Primary.ID, ":") {
							return fmt.Errorf("Expected V2 real ID without colons, got: %s", rs.Primary.ID)
						}
						return nil
					},
				),
				Config: testAccTenantCommonVariableMigrationV2(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, tenantVariablesLocalName, finalValue),
			},
		},
	})
}

func testAccTenantCommonVariableMigrationV1(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, localName, value string) string {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(1, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		testAccProjectGroup(projectGroupLocalName, projectGroupName)+"\n"+
		testAccEnvironment(env1LocalName, env1Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(env2LocalName, env2Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
        resource "octopusdeploy_library_variable_set" "test-library-variable-set-migration" {
            name = "test-migration"

            template {
                default_value = "Default Value"
                help_text     = "This is the help text"
                label         = "Test Label"
                name          = "Test Template Migration"

                display_settings = {
                    "Octopus.ControlType" = "Sensitive"
                }
            }
        }

        resource "octopusdeploy_project" "%[1]s" {
            included_library_variable_sets = [octopusdeploy_library_variable_set.test-library-variable-set-migration.id]
            lifecycle_id                   = octopusdeploy_lifecycle.%[2]s.id
            name                           = "%[3]s"
            project_group_id               = octopusdeploy_project_group.%[4]s.id
        }

        resource "octopusdeploy_tenant" "%[5]s" {
            name = "%[6]s"
        }

        resource "octopusdeploy_tenant_project" "project_environment" {
            tenant_id        = octopusdeploy_tenant.%[5]s.id
            project_id       = octopusdeploy_project.%[1]s.id
            environment_ids  = [octopusdeploy_environment.%[7]s.id, octopusdeploy_environment.%[8]s.id]
        }

        resource "octopusdeploy_tenant_common_variable" "%[9]s" {
            library_variable_set_id = octopusdeploy_library_variable_set.test-library-variable-set-migration.id
            template_id             = octopusdeploy_library_variable_set.test-library-variable-set-migration.template[0].id
            tenant_id               = octopusdeploy_tenant.%[5]s.id
            value                   = "%[10]s"

            depends_on = [octopusdeploy_tenant_project.project_environment]
        }`, projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName, tenantLocalName, tenantName, env1LocalName, env2LocalName, localName, value)
}

func testAccTenantCommonVariableMigrationV2(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, localName, value string) string {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(1, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		testAccProjectGroup(projectGroupLocalName, projectGroupName)+"\n"+
		testAccEnvironment(env1LocalName, env1Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(env2LocalName, env2Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
        resource "octopusdeploy_library_variable_set" "test-library-variable-set-migration" {
            name = "test-migration"

            template {
                default_value = "Default Value"
                help_text     = "This is the help text"
                label         = "Test Label"
                name          = "Test Template Migration"

                display_settings = {
                    "Octopus.ControlType" = "Sensitive"
                }
            }
        }

        resource "octopusdeploy_project" "%[1]s" {
            included_library_variable_sets = [octopusdeploy_library_variable_set.test-library-variable-set-migration.id]
            lifecycle_id                   = octopusdeploy_lifecycle.%[2]s.id
            name                           = "%[3]s"
            project_group_id               = octopusdeploy_project_group.%[4]s.id
        }

        resource "octopusdeploy_tenant" "%[5]s" {
            name = "%[6]s"
        }

        resource "octopusdeploy_tenant_project" "project_environment" {
            tenant_id        = octopusdeploy_tenant.%[5]s.id
            project_id       = octopusdeploy_project.%[1]s.id
            environment_ids  = [octopusdeploy_environment.%[7]s.id, octopusdeploy_environment.%[8]s.id]
        }

        resource "octopusdeploy_tenant_common_variable" "%[9]s" {
            library_variable_set_id = octopusdeploy_library_variable_set.test-library-variable-set-migration.id
            template_id             = octopusdeploy_library_variable_set.test-library-variable-set-migration.template[0].id
            tenant_id               = octopusdeploy_tenant.%[5]s.id
            value                   = "%[10]s"

            scope {
                environment_ids = [octopusdeploy_environment.%[7]s.id, octopusdeploy_environment.%[8]s.id]
            }

            depends_on = [octopusdeploy_tenant_project.project_environment]
        }`, projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName, tenantLocalName, tenantName, env1LocalName, env2LocalName, localName, value)
}

// TestAccTenantCommonVariableWithScope tests V2 API with environment scoping
func TestAccTenantCommonVariableWithScope(t *testing.T) {
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
	tenantVariablesLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := "octopusdeploy_tenant_common_variable." + tenantVariablesLocalName

	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTenantCommonVariableCheckDestroyV2,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Create with scope
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExistsV2(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.environment_ids.#", "2"),
				),
				Config: testAccTenantCommonVariableWithScope(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, tenantVariablesLocalName, value),
			},
			{
				// Update value
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExistsV2(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", newValue),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.environment_ids.#", "2"),
				),
				Config: testAccTenantCommonVariableWithScope(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, tenantVariablesLocalName, newValue),
			},
		},
	})
}

func testAccTenantCommonVariableWithScope(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, localName, value string) string {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(1, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccLifecycle(lifecycleLocalName, lifecycleName)+"\n"+
		testAccProjectGroup(projectGroupLocalName, projectGroupName)+"\n"+
		testAccEnvironment(env1LocalName, env1Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(env2LocalName, env2Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+`
        resource "octopusdeploy_library_variable_set" "test-library-variable-set" {
            name = "test-scope"

            template {
                default_value = "Default Value"
                help_text     = "This is the help text"
                label         = "Test Label"
                name          = "Test Template Scope"

                display_settings = {
                    "Octopus.ControlType" = "Sensitive"
                }
            }
        }

        resource "octopusdeploy_project" "%[1]s" {
            included_library_variable_sets = [octopusdeploy_library_variable_set.test-library-variable-set.id]
            lifecycle_id                   = octopusdeploy_lifecycle.%[2]s.id
            name                           = "%[3]s"
            project_group_id               = octopusdeploy_project_group.%[4]s.id
        }

        resource "octopusdeploy_tenant" "%[5]s" {
            name = "%[6]s"
        }

        resource "octopusdeploy_tenant_project" "project_environment" {
            tenant_id        = octopusdeploy_tenant.%[5]s.id
            project_id       = octopusdeploy_project.%[1]s.id
            environment_ids  = [octopusdeploy_environment.%[7]s.id, octopusdeploy_environment.%[8]s.id]
        }

        resource "octopusdeploy_tenant_common_variable" "%[9]s" {
            library_variable_set_id = octopusdeploy_library_variable_set.test-library-variable-set.id
            template_id             = octopusdeploy_library_variable_set.test-library-variable-set.template[0].id
            tenant_id               = octopusdeploy_tenant.%[5]s.id
            value                   = "%[10]s"

            scope {
                environment_ids = [octopusdeploy_environment.%[7]s.id, octopusdeploy_environment.%[8]s.id]
            }

            depends_on = [octopusdeploy_tenant_project.project_environment]
        }`, projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName, tenantLocalName, tenantName, env1LocalName, env2LocalName, localName, value)
}

func testTenantCommonVariableExistsV2(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if len(rs.Primary.ID) == 0 {
			return fmt.Errorf("Tenant common variable ID is not set")
		}

		// V2 uses real IDs (e.g., TenantVariables-123), not composite IDs with colons
		if strings.Contains(rs.Primary.ID, ":") {
			return fmt.Errorf("Expected V2 ID (e.g., TenantVariables-123) but got V1 composite ID: %s", rs.Primary.ID)
		}

		// Use V2 API to verify the variable exists
		tenantID := rs.Primary.Attributes["tenant_id"]
		spaceID := rs.Primary.Attributes["space_id"]

		client := octoClient
		query := variables.GetTenantCommonVariablesQuery{
			TenantID:                tenantID,
			SpaceID:                 spaceID,
			IncludeMissingVariables: false,
		}

		getResp, err := tenants.GetCommonVariables(client, query)
		if err != nil {
			return fmt.Errorf("Error retrieving tenant common variables: %s", err.Error())
		}

		// Search for the variable by ID
		for _, v := range getResp.Variables {
			if v.GetID() == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("Tenant common variable with ID %s not found via V2 API", rs.Primary.ID)
	}
}

// testAccTenantCommonVariableCheckDestroyV2 checks that V2 tenant common variables are destroyed
func testAccTenantCommonVariableCheckDestroyV2(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_tenant_common_variable" {
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
		query := variables.GetTenantCommonVariablesQuery{
			TenantID:                tenantID,
			SpaceID:                 spaceID,
			IncludeMissingVariables: false,
		}

		getResp, err := tenants.GetCommonVariables(client, query)
		if err != nil {
			// If we can't get the variables, assume they're gone
			return nil
		}

		// Search for the variable by ID
		for _, v := range getResp.Variables {
			if v.GetID() == rs.Primary.ID {
				return fmt.Errorf("Tenant common variable (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
