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
	space := NewTestSpace(t)

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
				Config: testAccRunbookBasic(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

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
				Config: testAccRunbookBasic(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccRunbookExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
				),
				Config: testAccRunbookBasic(localName, projectLocalName, name, newDescription, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

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
				Config: testAccRunbookWithConnectivityPolicy(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

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
				Config: testAccRunbookWithLegacyCountRetentionPolicy(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

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
				Config: testAccRunbookWithLegacyForeverRetentionPolicy(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

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
				Config: testAccRunbookWithCountStrategyRetentionPolicy(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile("Missing Required Field"),
				Config:      testAccRunbookWithCountStrategyRetentionPolicyWithoutUnit(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile("Missing Required Field"),
				Config:      testAccRunbookWithCountStrategyRetentionPolicyWithoutQuantity(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile("Invalid Field"),
				Config:      testAccRunbookWithForeverStrategyRetentionPolicyWithQuantity(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

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
				Config: testAccRunbookWithForeverStrategyRetentionPolicy(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

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
				Config: testAccRunbookWithDefaultStrategyRetentionPolicy(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccRunbookCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookBasic(localName, projectLocalName, name, description, projectName, space.ID, space.LifecycleID, space.ProjectGroupID),
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

func TestAccOctopusDeployRunbookWithTags(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_runbook." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagSetName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tagName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testAccRunbookCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookWithTags(localName, projectLocalName, projectName, name, tagSetName, tagName, space.ID, space.LifecycleID, space.ProjectGroupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "runbook_tags.#", "1"),
				),
			},
		},
	})
}

func testAccRunbookBasic(localName string, projectLocalName string, name string, description string, projectName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "%s"
		project_group_id                     = "%s"
		space_id                             = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "%s"
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, localName, projectLocalName, name, description, spaceID)
}

func testAccRunbookWithConnectivityPolicy(localName string, projectLocalName string, name string, description string, projectName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "%s"
		project_group_id                     = "%s"
		space_id                             = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "%s"

		connectivity_policy {
			allow_deployments_to_no_targets = true
			exclude_unhealthy_targets        = true
		}
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, localName, projectLocalName, name, description, spaceID)
}

func testAccRunbookWithLegacyCountRetentionPolicy(localName string, projectLocalName string, name string, description string, projectName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "%s"
		project_group_id                     = "%s"
		space_id                             = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "%s"

		retention_policy {
			quantity_to_keep = 10
		}
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, localName, projectLocalName, name, description, spaceID)
}

func testAccRunbookWithLegacyForeverRetentionPolicy(localName string, projectLocalName string, name string, description string, projectName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "%s"
		project_group_id                     = "%s"
		space_id                             = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "%s"

		retention_policy {
			should_keep_forever = true
		}
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, localName, projectLocalName, name, description, spaceID)
}

func testAccRunbookWithCountStrategyRetentionPolicy(localName string, projectLocalName string, name string, description string, projectName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "%s"
		project_group_id                     = "%s"
		space_id                             = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "%s"

		retention_policy_with_strategy {
			strategy		= "Count"
			quantity_to_keep = 10
			unit 		 = "Items"
		}
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, localName, projectLocalName, name, description, spaceID)
}

func testAccRunbookWithCountStrategyRetentionPolicyWithoutUnit(localName string, projectLocalName string, name string, description string, projectName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "%s"
		project_group_id                     = "%s"
		space_id                             = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "%s"

		retention_policy_with_strategy {
			strategy		= "Count"
			quantity_to_keep = 10
		}
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, localName, projectLocalName, name, description, spaceID)
}

func testAccRunbookWithCountStrategyRetentionPolicyWithoutQuantity(localName string, projectLocalName string, name string, description string, projectName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "%s"
		project_group_id                     = "%s"
		space_id                             = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "%s"

		retention_policy_with_strategy {
			strategy		= "Count"
			unit 		 	= "Items"
		}
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, localName, projectLocalName, name, description, spaceID)
}

func testAccRunbookWithForeverStrategyRetentionPolicyWithQuantity(localName string, projectLocalName string, name string, description string, projectName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "%s"
		project_group_id                     = "%s"
		space_id                             = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "%s"

		retention_policy_with_strategy {
			strategy		 = "Forever"
			quantity_to_keep = 10
		}
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, localName, projectLocalName, name, description, spaceID)
}

func testAccRunbookWithForeverStrategyRetentionPolicy(localName string, projectLocalName string, name string, description string, projectName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "%s"
		project_group_id                     = "%s"
		space_id                             = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "%s"

		retention_policy_with_strategy {
			strategy		= "Forever"
		}
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, localName, projectLocalName, name, description, spaceID)
}

func testAccRunbookWithDefaultStrategyRetentionPolicy(localName string, projectLocalName string, name string, description string, projectName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name                                 = "%s"
		lifecycle_id                         = "%s"
		project_group_id                     = "%s"
		space_id                             = "%s"
	}

	resource "octopusdeploy_runbook" "%s" {
		project_id   = octopusdeploy_project.%s.id
		name         = "%s"
		description  = "%s"
		space_id     = "%s"

		retention_policy_with_strategy {
			strategy		= "Default"
		}
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, localName, projectLocalName, name, description, spaceID)
}

func testAccRunbookExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[prefix]
		runbookID := rs.Primary.ID
		if _, err := runbooks.GetByID(octoClient, rs.Primary.Attributes["space_id"], runbookID); err != nil {
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

		if runbook, err := runbooks.GetByID(octoClient, rs.Primary.Attributes["space_id"], rs.Primary.ID); err == nil {
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

func testAccRunbookWithTags(localName string, projectLocalName string, projectName string, name string, tagSetName string, tagName string, spaceID string, lifecycleID string, projectGroupID string) string {
	return providerSpaceConfig(spaceID) + fmt.Sprintf(`
	resource "octopusdeploy_project" "%s" {
		name             = "%s"
		lifecycle_id     = "%s"
		project_group_id = "%s"
		space_id         = "%s"
	}

	resource "octopusdeploy_tag_set" "%s" {
		name     = "%s"
		space_id = "%s"
		scopes   = ["Runbook"]
	}

	resource "octopusdeploy_tag" "%s" {
		name       = "%s"
		color      = "#ff0000"
		tag_set_id = octopusdeploy_tag_set.%s.id
	}

	resource "octopusdeploy_runbook" "%s" {
		name         = "%s"
		project_id   = octopusdeploy_project.%s.id
		space_id     = "%s"
		runbook_tags = [octopusdeploy_tag.%s.canonical_tag_name]
	}`, projectLocalName, projectName, lifecycleID, projectGroupID, spaceID, tagSetName, tagSetName, spaceID, tagName, tagName, tagSetName, localName, name, projectLocalName, spaceID, tagName)
}
