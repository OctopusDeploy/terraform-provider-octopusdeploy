package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeploySpaceDefaultLifecycle_WithPhases(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_lifecycle." + localName
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	var developmentPhaseID string

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testSpaceCheckDestroy,
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testSpaceDefaultLifecycleWithPhases(localName, lifecycleName, description, false),
				Check: func(s *terraform.State) error {
					developmentPhaseID = s.RootModule().Resources[resourceName].Primary.Attributes["phases.0.id"]
					return nil
				},
			},
			{
				// Inserting Production before Development must keep the ID of Development
				Config: testSpaceDefaultLifecycleWithPhases(localName, lifecycleName, description, true),
				Check: func(s *terraform.State) error {
					attributes := s.RootModule().Resources[resourceName].Primary.Attributes
					if attributes["phases.1.id"] != developmentPhaseID {
						return fmt.Errorf("expected Development to keep phase ID %q but it is %q", developmentPhaseID, attributes["phases.1.id"])
					}
					if attributes["phases.0.id"] == "" || attributes["phases.0.id"] == developmentPhaseID {
						return fmt.Errorf("expected Production to get a new phase ID but it is %q", attributes["phases.0.id"])
					}
					return nil
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
				// Rename the default lifecycle and give its name to another lifecycle, so that it cannot be found by name
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

					other := lifecycles.NewLifecycle(schemas.SpaceDefaultLifecycleName)
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

	developmentPhase := fmt.Sprintf(`
			{
				name                         = "Development"
				automatic_deployment_targets = [octopusdeploy_environment.%[1]s_dev.id]
			},`, localName)
	phases := productionPhase + developmentPhase

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

		phases = [%[4]s
		]

		release_retention_with_strategy = {
			strategy         = "Count"
			quantity_to_keep = 15
			unit             = "Items"
		}

		tentacle_retention_with_strategy = {
			strategy = "Forever"
		}
	}`, localName, lifecycleName, description, phases)
}
