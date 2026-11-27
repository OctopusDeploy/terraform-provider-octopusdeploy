package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceStepTemplatesQuery(t *testing.T) {
	prefix := fmt.Sprintf("zz-%s", acctest.RandStringFromCharSet(20, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStepTemplatesConfig(prefix),
				Check: resource.ComposeTestCheckFunc(
					// partial_name narrows to the three templates sharing the prefix.
					resource.TestCheckResourceAttr("data.octopusdeploy_step_templates.by_partial_name", "step_templates.#", "3"),
					resource.TestCheckResourceAttr("data.octopusdeploy_step_templates.by_partial_name", "step_templates.0.name",
						fmt.Sprintf("%s-00", prefix)),
					// take bounds the result set.
					resource.TestCheckResourceAttr("data.octopusdeploy_step_templates.limited", "step_templates.#", "1"),
					// ids selects a specific template.
					resource.TestCheckResourceAttr("data.octopusdeploy_step_templates.by_id", "step_templates.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.octopusdeploy_step_templates.by_id", "step_templates.0.id",
						"octopusdeploy_step_template.templates.1", "id"),
					// A query matching nothing is an empty list, not an error.
					resource.TestCheckResourceAttr("data.octopusdeploy_step_templates.none", "step_templates.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceStepTemplatesConfig(prefix string) string {
	return fmt.Sprintf(`resource "octopusdeploy_step_template" "templates" {
  count           = 3
  action_type     = "Octopus.Script"
  name            = format("%s-%%02d", count.index)
  step_package_id = "Octopus.Script"
  packages        = []
  parameters      = []
  properties      = {
    "Octopus.Action.Script.ScriptBody"   = "echo \"hello\"",
    "Octopus.Action.Script.ScriptSource" = "Inline",
    "Octopus.Action.Script.Syntax"       = "PowerShell"
  }
}

data "octopusdeploy_step_templates" "by_partial_name" {
  partial_name = "%s"
  depends_on   = [octopusdeploy_step_template.templates]
}

data "octopusdeploy_step_templates" "limited" {
  partial_name = "%s"
  take         = 1
  depends_on   = [octopusdeploy_step_template.templates]
}

data "octopusdeploy_step_templates" "by_id" {
  ids = [octopusdeploy_step_template.templates[1].id]
}

data "octopusdeploy_step_templates" "none" {
  partial_name = "%s-matches-nothing"
  depends_on   = [octopusdeploy_step_template.templates]
}`, prefix, prefix, prefix, prefix)
}
