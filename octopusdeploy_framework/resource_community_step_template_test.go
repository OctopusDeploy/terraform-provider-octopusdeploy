package octopusdeploy_framework

import (
	"fmt"
	"os"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusCommunityStepTemplateBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_community_step_template." + localName
	website := "https://library.octopus.com/step-templates/04a74a00-967d-496a-a966-1acd17fededf"
	website2 := "https://library.octopus.com/step-templates/6042d737-5902-0729-ae57-8b6650a299da"

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testCommunityStepTemplateDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testCommunityStepTemplate(website, localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "community_action_template_id"),
					resource.TestCheckResourceAttrWith(prefix, "community_action_template_id", func(value string) error {
						fmt.Fprintf(os.Stderr, "Community Action Template ID: %s\n", value)
						return nil
					}),
				),
			},
			{
				Config: testCommunityStepTemplate(website2, localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttrSet(prefix, "community_action_template_id"),
				),
			},
		},
	})
}

func testCommunityStepTemplate(website string, name string) string {
	return fmt.Sprintf(`
		data "octopusdeploy_community_step_template" "community_step_template" {
			website = "%s"
		}	
		
		resource "octopusdeploy_community_step_template" "%s" {
			community_action_template_id=data.octopusdeploy_community_step_template.community_step_template.steps[0].id
		}
`,
		website,
		name,
	)
}

func testCommunityStepTemplateDestroy(s *terraform.State) error {
	if octoClient == nil {
		return fmt.Errorf("octoClient is nil")
	}

	var actionTemplateID string

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_community_step_template" {
			continue
		}

		actionTemplateID = rs.Primary.ID
		break
	}
	if actionTemplateID == "" {
		return fmt.Errorf("no octopusdeploy_community_step_template resource found")
	}

	actionTemplate, err := actiontemplates.GetByID(octoClient, octoClient.GetSpaceID(), actionTemplateID)
	if err == nil {
		if actionTemplate != nil {
			return fmt.Errorf("step template (%s) still exists", actionTemplate.Name)
		}
	}

	return nil
}
