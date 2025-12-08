package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployRunbookBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_runbook." + localName
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccRunbookExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "project_id"),
					resource.TestCheckResourceAttrSet(prefix, "space_id"),
				),
				Config: testAccRunbookBasic(localName, projectLocalName, name, description, projectName),
			},
		},
	})
}

func TestAccOctopusDeployRunbookUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_runbook." + localName
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccRunbookExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
				),
				Config: testAccRunbookBasic(localName, projectLocalName, name, description, projectName),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccRunbookExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
				),
				Config: testAccRunbookBasic(localName, projectLocalName, name, newDescription, projectName),
			},
		},
	})
}

func TestAccOctopusDeployRunbookWithConnectivityPolicy(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_runbook." + localName
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccRunbookExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "connectivity_policy.#", "1"),
					resource.TestCheckResourceAttr(prefix, "connectivity_policy.0.allow_deployments_to_no_targets", "true"),
					resource.TestCheckResourceAttr(prefix, "connectivity_policy.0.exclude_unhealthy_targets", "true"),
				),
				Config: testAccRunbookWithConnectivityPolicy(localName, projectLocalName, name, description, projectName),
			},
		},
	})
}

func TestAccOctopusDeployRunbookWithLegacyCountRetentionPolicy(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_runbook." + localName
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccRunbookExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(prefix, "retention_policy.0.quantity_to_keep", "10"),
				),
				Config: testAccRunbookWithLegacyCountRetentionPolicy(localName, projectLocalName, name, description, projectName),
			},
		},
	})
}

func TestAccOctopusDeployRunbookWithLegacyForeverRetentionPolicy(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_runbook." + localName
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccRunbookExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(prefix, "retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(prefix, "retention_policy.0.should_keep_forever", "true"),
				),
				Config: testAccRunbookWithLegacyForeverRetentionPolicy(localName, projectLocalName, name, description, projectName),
			},
		},
	})
}

func TestAccOctopusDeployRunbookWithCountStrategyRetentionPolicy(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_runbook." + localName
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccRunbookExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "retention_policy_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(prefix, "retention_policy_with_strategy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(prefix, "retention_policy_with_strategy.0.quantity_to_keep", "10"),
					resource.TestCheckResourceAttr(prefix, "retention_policy_with_strategy.0.unit", "Items"),
				),
				Config: testAccRunbookWithCountStrategyRetentionPolicy(localName, projectLocalName, name, description, projectName),
			},
		},
	})
}

func TestAccOctopusDeployRunbookWithCountStrategyRetentionPolicyWithoutUnit(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile("Missing Required Field"),
				Config:      testAccRunbookWithCountStrategyRetentionPolicyWithoutUnit(localName, projectLocalName, name, description, projectName),
			},
		},
	})
}

func TestAccOctopusDeployRunbookWithCountStrategyRetentionPolicyWithoutQuantity(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile("Missing Required Field"),
				Config:      testAccRunbookWithCountStrategyRetentionPolicyWithoutQuantity(localName, projectLocalName, name, description, projectName),
			},
		},
	})
}

func TestAccOctopusDeployRunbookWithForeverStrategyRetentionPolicyWithQuantity(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile("Invalid Field"),
				Config:      testAccRunbookWithForeverStrategyRetentionPolicyWithQuantity(localName, projectLocalName, name, description, projectName),
			},
		},
	})
}

func TestAccOctopusDeployRunbookWithForeverStrategyRetentionPolicy(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_runbook." + localName
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccRunbookExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "retention_policy_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(prefix, "retention_policy_with_strategy.0.strategy", "Forever"),
				),
				Config: testAccRunbookWithForeverStrategyRetentionPolicy(localName, projectLocalName, name, description, projectName),
			},
		},
	})
}

func TestAccOctopusDeployRunbookWithDefaultStrategyRetentionPolicy(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_runbook." + localName
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccRunbookExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
					resource.TestCheckResourceAttr(prefix, "retention_policy_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(prefix, "retention_policy_with_strategy.0.strategy", "Default"),
				),
				Config: testAccRunbookWithDefaultStrategyRetentionPolicy(localName, projectLocalName, name, description, projectName),
			},
		},
	})
}
func TestAccOctopusDeployRunbookImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_runbook." + localName
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookBasic(localName, projectLocalName, name, description, projectName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccRunbookImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccRunbookBasic(localName string, projectLocalName string, name string, description string, projectName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "Lifecycles-1"
		project_group_id                     = "ProjectGroups-1"
		space_id                             = "Spaces-1"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "Spaces-1"
	}`, projectLocalName, projectName, localName, projectLocalName, name, description)
}

func testAccRunbookWithConnectivityPolicy(localName string, projectLocalName string, name string, description string, projectName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "Lifecycles-1"
		project_group_id                     = "ProjectGroups-1"
		space_id                             = "Spaces-1"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "Spaces-1"
		
		connectivity_policy {
			allow_deployments_to_no_targets = true
			exclude_unhealthy_targets        = true
		}
	}`, projectLocalName, projectName, localName, projectLocalName, name, description)
}

func testAccRunbookWithLegacyCountRetentionPolicy(localName string, projectLocalName string, name string, description string, projectName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "Lifecycles-1"
		project_group_id                     = "ProjectGroups-1"
		space_id                             = "Spaces-1"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "Spaces-1"
		
		retention_policy {
			quantity_to_keep = 10
		}
	}`, projectLocalName, projectName, localName, projectLocalName, name, description)
}

func testAccRunbookWithLegacyForeverRetentionPolicy(localName string, projectLocalName string, name string, description string, projectName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "Lifecycles-1"
		project_group_id                     = "ProjectGroups-1"
		space_id                             = "Spaces-1"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "Spaces-1"
		
		retention_policy {
			should_keep_forever = true
		}
	}`, projectLocalName, projectName, localName, projectLocalName, name, description)
}

func testAccRunbookWithCountStrategyRetentionPolicy(localName string, projectLocalName string, name string, description string, projectName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "Lifecycles-1"
		project_group_id                     = "ProjectGroups-1"
		space_id                             = "Spaces-1"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "Spaces-1"
		
		retention_policy_with_strategy {
			strategy		= "Count"
			quantity_to_keep = 10
			unit 		 = "Items"
		}
	}`, projectLocalName, projectName, localName, projectLocalName, name, description)
}

func testAccRunbookWithCountStrategyRetentionPolicyWithoutUnit(localName string, projectLocalName string, name string, description string, projectName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "Lifecycles-1"
		project_group_id                     = "ProjectGroups-1"
		space_id                             = "Spaces-1"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "Spaces-1"
		
		retention_policy_with_strategy {
			strategy		= "Count"
			quantity_to_keep = 10
		}
	}`, projectLocalName, projectName, localName, projectLocalName, name, description)
}

func testAccRunbookWithCountStrategyRetentionPolicyWithoutQuantity(localName string, projectLocalName string, name string, description string, projectName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "Lifecycles-1"
		project_group_id                     = "ProjectGroups-1"
		space_id                             = "Spaces-1"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "Spaces-1"
		
		retention_policy_with_strategy {
			strategy		= "Count"
			unit 		 	= "Items"
		}
	}`, projectLocalName, projectName, localName, projectLocalName, name, description)
}

func testAccRunbookWithForeverStrategyRetentionPolicyWithQuantity(localName string, projectLocalName string, name string, description string, projectName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "Lifecycles-1"
		project_group_id                     = "ProjectGroups-1"
		space_id                             = "Spaces-1"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "Spaces-1"
		
		retention_policy_with_strategy {
			strategy		 = "Forever"
			quantity_to_keep = 10
		}
	}`, projectLocalName, projectName, localName, projectLocalName, name, description)
}

func testAccRunbookWithForeverStrategyRetentionPolicy(localName string, projectLocalName string, name string, description string, projectName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "Lifecycles-1"
		project_group_id                     = "ProjectGroups-1"
		space_id                             = "Spaces-1"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "Spaces-1"
		
		retention_policy_with_strategy {
			strategy		= "Forever"
		}
	}`, projectLocalName, projectName, localName, projectLocalName, name, description)
}

func testAccRunbookWithDefaultStrategyRetentionPolicy(localName string, projectLocalName string, name string, description string, projectName string) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "Lifecycles-1"
		project_group_id                     = "ProjectGroups-1"
		space_id                             = "Spaces-1"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "Spaces-1"
		
		retention_policy_with_strategy {
			strategy		= "Default"
		}
	}`, projectLocalName, projectName, localName, projectLocalName, name, description)
}

func testAccRunbookExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		runbookID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := runbooks.GetByID(octoClient, octoClient.GetSpaceID(), runbookID); err != nil {
			return err
		}

		return nil
	}
}

func testAccRunbookCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_runbook" {
			continue
		}

		if runbook, err := runbooks.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("runbook (%s) still exists", runbook.GetID())
		}
	}

	return nil
}

func testAccRunbookImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
