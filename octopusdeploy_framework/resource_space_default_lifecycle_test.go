package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeploySpaceDefaultLifecycle_Adopt(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_lifecycle." + localName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testSpaceCheckDestroy,
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testSpaceDefaultLifecycleEmpty(localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "space_id", "octopusdeploy_space."+localName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "Default Lifecycle"),
					resource.TestCheckNoResourceAttr(resourceName, "phases"),
					resource.TestCheckNoResourceAttr(resourceName, "release_retention_with_strategy"),
					resource.TestCheckNoResourceAttr(resourceName, "tentacle_retention_with_strategy"),
				),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultLifecycle_WithPhases(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_lifecycle." + localName
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testSpaceCheckDestroy,
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testSpaceDefaultLifecycleWithPhases(localName, lifecycleName, description, false),
				Check: resource.ComposeTestCheckFunc(
					testSpaceDefaultLifecycleMatches(resourceName, lifecycleName, 1),
					resource.TestCheckResourceAttr(resourceName, "name", lifecycleName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "phases.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "phases.0.name", "Development"),
					resource.TestCheckResourceAttrSet(resourceName, "phases.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "phases.0.automatic_deployment_targets.0", "octopusdeploy_environment."+localName+"_dev", "id"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.strategy", "Count"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.quantity_to_keep", "15"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.unit", "Items"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy.strategy", "Forever"),
				),
			},
			{
				Config: testSpaceDefaultLifecycleWithPhases(localName, lifecycleName, description, true),
				Check: resource.ComposeTestCheckFunc(
					testSpaceDefaultLifecycleMatches(resourceName, lifecycleName, 2),
					resource.TestCheckResourceAttr(resourceName, "phases.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "phases.1.name", "Production"),
					resource.TestCheckResourceAttr(resourceName, "phases.1.is_priority_phase", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "phases.1.optional_deployment_targets.0", "octopusdeploy_environment."+localName+"_prod", "id"),
					resource.TestCheckResourceAttr(resourceName, "phases.1.release_retention_with_strategy.strategy", "Forever"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("Not found: %s", resourceName)
					}
					return rs.Primary.Attributes["space_id"], nil
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("Not found: %s", resourceName)
					}
					return rs.Primary.Attributes["space_id"] + ":" + rs.Primary.ID, nil
				},
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultLifecycle_AdoptRenamed(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_lifecycle." + localName
	renamedName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	otherName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	var spaceID, defaultLifecycleID string

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testSpaceCheckDestroy,
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testSpaceDefaultLifecycleSpace(localName),
				Check: func(s *terraform.State) error {
					spaceID = s.RootModule().Resources["octopusdeploy_space."+localName].Primary.ID
					return nil
				},
			},
			{
				// Rename the default lifecycle and add another lifecycle, so it can no longer be found by name
				PreConfig: func() {
					all, err := lifecycles.GetAll(octoClient, spaceID)
					if err != nil || len(all) != 1 {
						t.Fatalf("expected exactly one lifecycle in space %s: %v", spaceID, err)
					}
					defaultLifecycle := all[0]
					defaultLifecycleID = defaultLifecycle.GetID()
					defaultLifecycle.Name = renamedName
					if _, err := lifecycles.Update(octoClient, defaultLifecycle); err != nil {
						t.Fatalf("unable to rename default lifecycle: %v", err)
					}

					other := lifecycles.NewLifecycle(otherName)
					other.SpaceID = spaceID
					if _, err := lifecycles.Add(octoClient, other); err != nil {
						t.Fatalf("unable to add lifecycle: %v", err)
					}
				},
				Config: testSpaceDefaultLifecycleEmpty(localName),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						return resource.TestCheckResourceAttr(resourceName, "id", defaultLifecycleID)(s)
					},
					resource.TestCheckResourceAttr(resourceName, "name", renamedName),
				),
			},
		},
	})
}

func testSpaceDefaultLifecycleSpace(localName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_space" "%[1]s" {
		name                 = "%[1]s"
		space_managers_teams = ["teams-administrators"]
	}`, localName)
}

func testSpaceDefaultLifecycleEmpty(localName string) string {
	return testSpaceDefaultLifecycleSpace(localName) + fmt.Sprintf(`

	resource "octopusdeploy_space_default_lifecycle" "%[1]s" {
		space_id = octopusdeploy_space.%[1]s.id
	}`, localName)
}

func testSpaceDefaultLifecycleWithPhases(localName string, lifecycleName string, description string, includeProduction bool) string {
	productionPhase := ""
	if includeProduction {
		productionPhase = fmt.Sprintf(`
			{
				name                        = "Production"
				optional_deployment_targets = [octopusdeploy_environment.%[1]s_prod.id]
				is_priority_phase           = true
				release_retention_with_strategy = {
					strategy = "Forever"
				}
			},`, localName)
	}

	return testSpaceDefaultLifecycleSpace(localName) + fmt.Sprintf(`

	resource "octopusdeploy_environment" "%[1]s_dev" {
		name     = "Development"
		space_id = octopusdeploy_space.%[1]s.id
	}

	resource "octopusdeploy_environment" "%[1]s_prod" {
		name     = "Production"
		space_id = octopusdeploy_space.%[1]s.id
	}

	resource "octopusdeploy_space_default_lifecycle" "%[1]s" {
		space_id    = octopusdeploy_space.%[1]s.id
		name        = "%[2]s"
		description = "%[3]s"

		phases = [
			{
				name                         = "Development"
				automatic_deployment_targets = [octopusdeploy_environment.%[1]s_dev.id]
			},%[4]s
		]

		release_retention_with_strategy = {
			strategy         = "Count"
			quantity_to_keep = 15
			unit             = "Items"
		}

		tentacle_retention_with_strategy = {
			strategy = "Forever"
		}
	}`, localName, lifecycleName, description, productionPhase)
}

func testSpaceDefaultLifecycleMatches(resourceName string, expectedName string, expectedPhases int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		lifecycle, err := lifecycles.GetByID(octoClient, rs.Primary.Attributes["space_id"], rs.Primary.ID)
		if err != nil {
			return err
		}
		if lifecycle.Name != expectedName {
			return fmt.Errorf("expected default lifecycle name %q but was %q", expectedName, lifecycle.Name)
		}
		if len(lifecycle.Phases) != expectedPhases {
			return fmt.Errorf("expected %d phases but found %d", expectedPhases, len(lifecycle.Phases))
		}

		return nil
	}
}
