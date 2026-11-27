package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Covers https://github.com/OctopusDeploy/terraform-provider-octopusdeploy/issues/47.
// A process that lives outside the provider's space cannot be imported by its bare
// identifier, because the lookup resolves against the client's space. Prefixing the
// identifier with the space makes it importable, and the failure without one says so.
func TestAccProcessImportFromAnotherSpace(t *testing.T) {
	spaceName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	config := testAccProcessInSpaceConfiguration(spaceName, localName)

	orderResource := "octopusdeploy_process_steps_order." + localName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(orderResource, "process_id"),
					resource.TestCheckResourceAttrSet(orderResource, "space_id"),
				),
			},
			{
				// Without the space the process cannot be found, and the diagnostic
				// has to explain which space was searched.
				Config:            config,
				ResourceName:      orderResource,
				ImportState:       true,
				ImportStateIdFunc: processIdFromState(orderResource, false),
				ExpectError:       regexp.MustCompile(`(?s)Could not find.*If it belongs to another space`),
			},
			{
				Config:            config,
				ResourceName:      orderResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: processIdFromState(orderResource, true),
			},
		},
	})
}

// processIdFromState builds the import identifier for the steps order resource,
// optionally prefixed with the space it lives in.
func processIdFromState(resourceAddress string, qualifyWithSpace bool) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceAddress]
		if !ok {
			return "", fmt.Errorf("resource %q not found in state", resourceAddress)
		}

		processId := rs.Primary.Attributes["process_id"]
		if !qualifyWithSpace {
			return processId, nil
		}

		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["space_id"], processId), nil
	}
}

func testAccProcessInSpaceConfiguration(spaceName string, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_space" "%[2]s" {
		  name                  = "%[1]s"
		  is_default            = false
		  is_task_queue_stopped = true
		  description           = "Space for the process import test"
		  space_managers_teams  = ["teams-administrators"]
		}

		data "octopusdeploy_lifecycles" "%[2]s" {
		  ids          = null
		  partial_name = "Default Lifecycle"
		  space_id     = octopusdeploy_space.%[2]s.id
		  skip         = 0
		  take         = 1
		}

		resource "octopusdeploy_project_group" "%[2]s" {
		  space_id = octopusdeploy_space.%[2]s.id
		  name     = "%[1]s"
		}

		resource "octopusdeploy_project" "%[2]s" {
		  space_id         = octopusdeploy_space.%[2]s.id
		  name             = "%[1]s"
		  lifecycle_id     = data.octopusdeploy_lifecycles.%[2]s.lifecycles[0].id
		  project_group_id = octopusdeploy_project_group.%[2]s.id
		}

		resource "octopusdeploy_process" "%[2]s" {
		  space_id   = octopusdeploy_space.%[2]s.id
		  project_id = octopusdeploy_project.%[2]s.id
		}

		resource "octopusdeploy_process_step" "%[2]s" {
		  space_id   = octopusdeploy_space.%[2]s.id
		  process_id = octopusdeploy_process.%[2]s.id
		  name       = "%[1]s"
		  type       = "Octopus.Script"
		  properties = {
		    "Octopus.Action.TargetRoles" = "role-one"
		  }
		  execution_properties = {
		    "Octopus.Action.Script.ScriptBody" = "."
		  }
		}

		resource "octopusdeploy_process_steps_order" "%[2]s" {
		  space_id   = octopusdeploy_space.%[2]s.id
		  process_id = octopusdeploy_process.%[2]s.id
		  steps      = [octopusdeploy_process_step.%[2]s.id]
		}
	`, spaceName, localName)
}
