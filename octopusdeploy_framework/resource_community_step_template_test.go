package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccOctopusStepTemplateBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_step_template." + localName
	data := stepTemplateTestData{
		localName:     localName,
		prefix:        prefix,
		actionType:    "Octopus.Script",
		name:          acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
		description:   acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		stepPackageID: "Octopus.Script",
		packages: []stepTemplatePackageTestData{
			{
				packageID:          "force",
				acquisitonLocation: "Server",
				feedID:             "feeds-builtin",
				name:               "mypackage",
				properties: stepTemplatePackagePropsTestData{
					extract:       "True",
					purpose:       "",
					selectionMode: "immediate",
				},
			},
		},
		parameters: []stepTemplateParamTestData{
			{
				defaultValue: "Hello World",
				displaySettings: map[string]string{
					"Octopus.ControlType": "SingleLineText",
				},
				helpText: acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				label:    acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				name:     acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				id:       "621e1584-cdf3-4b67-9204-fc82430c908c",
			},
			{
				defaultValue: "Hello Earth",
				displaySettings: map[string]string{
					"Octopus.ControlType": "SingleLineText",
				},
				helpText: acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				label:    acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				name:     acctest.RandStringFromCharSet(10, acctest.CharSetAlpha),
				id:       "cd731d21-669a-42e1-81af-048681fd5c69",
			},
		},
		properties: map[string]string{
			"Octopus.Action.Script.ScriptBody":   "echo 'Hello World'",
			"Octopus.Action.Script.ScriptSource": "Inline",
			"Octopus.Action.Script.Syntax":       "Bash",
		},
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testStepTemplateDestroy(s, localName) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testStepTemplateRunScriptBasic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "name", data.name),
				),
			},
			{
				Config: testStepTemplateRunScriptUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "name", data.name+"-updated"),
				),
			},
		},
	})
}

func testStepTemplateRunScriptBasic(data stepTemplateTestData) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_community_step_template" "template" {
			community_action_template_id="CommunityActionTemplates-11"
		}
`,
	)
}

func testStepTemplateRunScriptUpdate(data stepTemplateTestData) string {
	data.name = data.name + "-updated"

	return testStepTemplateRunScriptBasic(data)
}

func testStepTemplateDestroy(s *terraform.State, localName string) error {
	var actionTemplateID string

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_step_template" {
			continue
		}

		actionTemplateID = rs.Primary.ID
		break
	}
	if actionTemplateID == "" {
		return fmt.Errorf("no octopusdeploy_step_template resource found")
	}

	actionTemplate, err := actiontemplates.GetByID(octoClient, octoClient.GetSpaceID(), actionTemplateID)
	if err == nil {
		if actionTemplate != nil {
			return fmt.Errorf("step template (%s) still exists", actionTemplate.Name)
		}
	}

	return nil
}
