package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceStepTemplates(t *testing.T) {
	localName := acctest.RandStringFromCharSet(50, acctest.CharSetAlpha)
	prefix := fmt.Sprintf("data.octopusdeploy_step_template.%s", localName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: createTestAccDataSourceStepTemplateConfig(),
			},
			{
				Check:  resource.TestCheckResourceAttr(prefix, "step_template.name", "Hello World"),
				Config: testAccDataSourceStepTemplateConfig(localName),
			},
		},
	})
}

func testAccDataSourceStepTemplateConfig(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_step_template" "%s" {
		name = "Hello World"
	}`, localName)
}

func createTestAccDataSourceStepTemplateConfig() string {
	return `resource "octopusdeploy_step_template" "steptemplate_hello_world" {
  action_type     = "Octopus.Script"
  name            = "Hello World"
  step_package_id = "Octopus.Script"
  packages        = []
  parameters      = [
    {
      default_value = "World!",
      display_settings = { "Octopus.ControlType" = "SingleLineText" },
      help_text = null,
      id = "fb95b2e8-3395-4b63-9c23-549c133841ab",
      label = null,
      name = "HelloWorld.Message"
    },
    {
      default_sensitive_value = "SecretValue",
      display_settings = { "Octopus.ControlType" = "Sensitive" },
      help_text = null,
      id = "ca5b66cc-c859-407b-b4df-d6bab42ad2f1",
      label = null,
      name = "HelloWorld.Secret"
    }
  ]
  properties      = {
    "Octopus.Action.Script.ScriptBody" = "echo \"Hello #{HelloWorld.Message}\"",
    "Octopus.Action.Script.ScriptSource" = "Inline",
    "Octopus.Action.Script.Syntax" = "PowerShell"
  }
}`
}
