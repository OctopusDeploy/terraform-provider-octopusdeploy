package octopusdeploy_framework

import (
	"fmt"
	"strings"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
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

			for _, v := range getResp.Variables {
				if v.GetID() == rs.Primary.ID {
					return nil
				}
			}

			return fmt.Errorf("Tenant common variable with ID %s not found via V2 API", rs.Primary.ID)
		}

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
				return nil
			}

			for _, v := range getResp.Variables {
				if v.GetID() == rs.Primary.ID {
					return fmt.Errorf("Tenant common variable (%s) still exists", rs.Primary.ID)
				}
			}

			continue
		}

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
				ProtoV6ProviderFactories: ProtoV6ProviderFactoriesWithFeatureToggleOverrides(map[string]bool{
					"CommonVariableScopingFeatureToggle": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckNoResourceAttr(resourceName, "scope.#"),
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[resourceName]
						if !strings.Contains(rs.Primary.ID, ":") {
							return fmt.Errorf("Expected V1 composite ID with colons, got: %s", rs.Primary.ID)
						}
						return nil
					},
				),
				Config: testAccTenantCommonVariableMigrationV1(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, tenantVariablesLocalName, value),
			},
			{
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
						if strings.Contains(rs.Primary.ID, ":") {
							return fmt.Errorf("Expected V2 real ID without colons, got: %s", rs.Primary.ID)
						}
						return nil
					},
				),
				Config: testAccTenantCommonVariableMigrationV2(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, tenantVariablesLocalName, newValue),
			},
			{
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
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExistsV2(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttr(resourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope.0.environment_ids.#", "2"),
				),
				Config: testAccTenantCommonVariableWithScope(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, env1LocalName, env1Name, env2LocalName, env2Name, tenantLocalName, tenantName, tenantVariablesLocalName, value),
			},
			{
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

		if strings.Contains(rs.Primary.ID, ":") {
			return fmt.Errorf("Expected V2 ID (e.g., TenantVariables-123) but got V1 composite ID: %s", rs.Primary.ID)
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
			return fmt.Errorf("Error retrieving tenant common variables: %s", err.Error())
		}

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

		if strings.Contains(rs.Primary.ID, ":") {
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
			return nil
		}

		for _, v := range getResp.Variables {
			if v.GetID() == rs.Primary.ID {
				return fmt.Errorf("Tenant common variable (%s) still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

// TestAccTenantCommonVariableImportV1
// Covers the existing v1 import not having a space id, failing v2 updates after the import.
func TestAccTenantCommonVariableImportV1(t *testing.T) {
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	librarySetLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	librarySetName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	variableLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resourceName := "octopusdeploy_tenant_common_variable." + variableLocalName
	importedValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	updatedValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	var ids tenantCommonVariableIDs

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccTenantCommonVariableCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTenantCommonVariableImportV1Dependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, environmentLocalName, environmentName, librarySetLocalName, librarySetName, tenantLocalName, tenantName),
				Check: testAccCaptureTenantCommonVariableIDs(
					"octopusdeploy_tenant."+tenantLocalName,
					"octopusdeploy_library_variable_set."+librarySetLocalName,
					&ids,
				),
			},
			{
				// Create the variable V1, and import
				PreConfig: func() {
					if err := createTenantCommonVariableV1(ids, importedValue); err != nil {
						t.Fatalf("create tenant common variable outside Terraform: %s", err)
					}
				},
				Config:             testAccTenantCommonVariableImportV1(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, environmentLocalName, environmentName, librarySetLocalName, librarySetName, tenantLocalName, tenantName, variableLocalName, importedValue),
				ResourceName:       resourceName,
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return ids.v1ImportID(), nil
				},
				ImportStateCheck: testAccCheckImportedTenantCommonVariableHasSpaceID(&ids),
			},
			{
				// The genuine update that used to fail.
				Config: testAccTenantCommonVariableImportV1(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, environmentLocalName, environmentName, librarySetLocalName, librarySetName, tenantLocalName, tenantName, variableLocalName, updatedValue),
				Check: resource.ComposeTestCheckFunc(
					testTenantCommonVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", updatedValue),
					testAccCheckTenantCommonVariableSpaceID(resourceName, &ids),
				),
			},
		},
	})
}

type tenantCommonVariableIDs struct {
	spaceID              string
	tenantID             string
	libraryVariableSetID string
	templateID           string
}

func (s tenantCommonVariableIDs) v1ImportID() string {
	return strings.Join([]string{s.tenantID, s.libraryVariableSetID, s.templateID}, ":")
}

func testAccCaptureTenantCommonVariableIDs(tenantResourceName string, librarySetResourceName string, ids *tenantCommonVariableIDs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attribute := func(resourceName string, key string) (string, error) {
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return "", fmt.Errorf("resource %q was not found in Terraform state", resourceName)
			}
			value := rs.Primary.Attributes[key]
			if value == "" && key == "id" {
				value = rs.Primary.ID
			}
			if value == "" {
				return "", fmt.Errorf("resource %q has no %s", resourceName, key)
			}
			return value, nil
		}

		var err error
		if ids.tenantID, err = attribute(tenantResourceName, "id"); err != nil {
			return err
		}
		if ids.libraryVariableSetID, err = attribute(librarySetResourceName, "id"); err != nil {
			return err
		}
		if ids.templateID, err = attribute(librarySetResourceName, "template.0.id"); err != nil {
			return err
		}
		if ids.spaceID, err = attribute(tenantResourceName, "space_id"); err != nil {
			return err
		}
		return nil
	}
}

func createTenantCommonVariableV1(ids tenantCommonVariableIDs, value string) error {
	tenant, err := tenants.GetByID(octoClient, ids.spaceID, ids.tenantID)
	if err != nil {
		return err
	}

	tenantVariables, err := octoClient.Tenants.GetVariables(tenant)
	if err != nil {
		return err
	}

	libraryVariable, ok := tenantVariables.LibraryVariables[ids.libraryVariableSetID]
	if !ok {
		return fmt.Errorf("tenant %s is not connected to library variable set %s", ids.tenantID, ids.libraryVariableSetID)
	}

	libraryVariable.Variables[ids.templateID] = core.NewPropertyValue(value, false)

	_, err = octoClient.Tenants.UpdateVariables(tenant, tenantVariables)
	return err
}

func testAccCheckImportedTenantCommonVariableHasSpaceID(ids *tenantCommonVariableIDs) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		for _, state := range states {
			if state.Ephemeral.Type != "octopusdeploy_tenant_common_variable" {
				continue
			}

			if spaceID := state.Attributes["space_id"]; spaceID != ids.spaceID {
				return fmt.Errorf("imported tenant common variable has space_id %q, want %q", spaceID, ids.spaceID)
			}
			return nil
		}

		return fmt.Errorf("imported state does not contain a tenant common variable")
	}
}

func testAccCheckTenantCommonVariableSpaceID(resourceName string, ids *tenantCommonVariableIDs) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestCheckResourceAttr(resourceName, "space_id", ids.spaceID)(s)
	}
}

func testTenantCommonVariable(localName string, librarySetLocalName string, tenantLocalName string, value string) string {
	return fmt.Sprintf(`resource "octopusdeploy_tenant_common_variable" "%[1]s" {
		library_variable_set_id = octopusdeploy_library_variable_set.%[2]s.id
		template_id             = octopusdeploy_library_variable_set.%[2]s.template[0].id
		tenant_id               = octopusdeploy_tenant.%[3]s.id
		value                   = "%[4]s"
		depends_on              = [octopusdeploy_tenant_project.project_environment]
	}`, localName, librarySetLocalName, tenantLocalName, value)
}

func testAccTenantCommonVariableImportV1(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, environmentLocalName string, environmentName string, librarySetLocalName string, librarySetName string, tenantLocalName string, tenantName string, variableLocalName string, value string) string {
	return testAccTenantCommonVariableImportV1Dependencies(lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, environmentLocalName, environmentName, librarySetLocalName, librarySetName, tenantLocalName, tenantName) + "\n" +
		testTenantCommonVariable(variableLocalName, librarySetLocalName, tenantLocalName, value)
}

func testAccTenantCommonVariableImportV1Dependencies(lifecycleLocalName string, lifecycleName string, projectGroupLocalName string, projectGroupName string, projectLocalName string, projectName string, environmentLocalName string, environmentName string, librarySetLocalName string, librarySetName string, tenantLocalName string, tenantName string) string {
	return testAccLifecycle(lifecycleLocalName, lifecycleName) + "\n" +
		testAccProjectGroup(projectGroupLocalName, projectGroupName) + "\n" +
		testAccEnvironment(environmentLocalName, environmentName, acctest.RandStringFromCharSet(20, acctest.CharSetAlpha), false, acctest.RandIntRange(1, 10), false) + "\n" +
		fmt.Sprintf(`resource "octopusdeploy_library_variable_set" "%[1]s" {
		name = "%[2]s"

		template {
			name          = "common variable template name"
			label         = "common variable template label"
			default_value = "default"

			display_settings = {
				"Octopus.ControlType" = "SingleLineText"
			}
		}
	}

	resource "octopusdeploy_project" "%[3]s" {
		included_library_variable_sets = [octopusdeploy_library_variable_set.%[1]s.id]
		lifecycle_id                   = octopusdeploy_lifecycle.%[4]s.id
		name                           = "%[5]s"
		project_group_id               = octopusdeploy_project_group.%[6]s.id
	}

	resource "octopusdeploy_tenant" "%[7]s" {
		name = "%[8]s"
	}

	resource "octopusdeploy_tenant_project" "project_environment" {
		tenant_id       = octopusdeploy_tenant.%[7]s.id
		project_id      = octopusdeploy_project.%[3]s.id
		environment_ids = [octopusdeploy_environment.%[9]s.id]
	}`, librarySetLocalName, librarySetName, projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName, tenantLocalName, tenantName, environmentLocalName)
}
