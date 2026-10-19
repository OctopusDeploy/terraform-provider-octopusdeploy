package octopusdeploy_framework

import (
	"fmt"
	"regexp"
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
				Check: resource.TestCheckResourceAttr(prefix, "step_template.name", "Hello World"),
				// The template stays in the config: dropping it here destroys the very
				// template the data source reads.
				Config: createTestAccDataSourceStepTemplateConfig() + "\n" + testAccDataSourceStepTemplateConfig(localName),
			},
		},
	})
}

// The API returns 30 step templates per page by default. Looking up by name used to filter
// the first page in memory, so a template that sorted past it was invisible; looking up by ID
// ignored the ID and returned whichever template happened to come back first.
func TestAccDataSourceStepTemplateLookupBeyondFirstPage(t *testing.T) {
	prefix := fmt.Sprintf("zz-%s", acctest.RandStringFromCharSet(20, acctest.CharSetAlpha))
	lastName := fmt.Sprintf("%s-%02d", prefix, stepTemplatePageTestCount-1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStepTemplatePaginationConfig(prefix),
				Check: resource.ComposeTestCheckFunc(
					// An ID lookup returns the template that was asked for, not the first one.
					resource.TestCheckResourceAttrPair(
						"data.octopusdeploy_step_template.by_id", "step_template.id",
						"octopusdeploy_step_template.page.5", "id"),
					resource.TestCheckResourceAttr("data.octopusdeploy_step_template.by_id", "step_template.name",
						fmt.Sprintf("%s-05", prefix)),
					// A name that sorts onto the second page still resolves.
					resource.TestCheckResourceAttr("data.octopusdeploy_step_template.by_name", "step_template.name", lastName),
				),
			},
		},
	})
}

// One more than the server's default page size of 30.
const stepTemplatePageTestCount = 31

func testAccDataSourceStepTemplatePaginationConfig(prefix string) string {
	return fmt.Sprintf(`resource "octopusdeploy_step_template" "page" {
  count           = %d
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

data "octopusdeploy_step_template" "by_name" {
  name       = format("%s-%%02d", %d)
  depends_on = [octopusdeploy_step_template.page]
}

data "octopusdeploy_step_template" "by_id" {
  id = octopusdeploy_step_template.page[5].id
}`, stepTemplatePageTestCount, prefix, prefix, stepTemplatePageTestCount-1)
}

func TestAccDataSourceStepTemplateNotFound(t *testing.T) {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`data "octopusdeploy_step_template" "missing" {
					name = "%s"
				}`, name),
				ExpectError: regexp.MustCompile("Step Template not found"),
			},
		},
	})
}

func testAccDataSourceStepTemplateConfig(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_step_template" "%s" {
		name       = "Hello World"
		depends_on = [octopusdeploy_step_template.steptemplate_hello_world]
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
