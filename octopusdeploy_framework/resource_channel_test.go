package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccChannelBasic(t *testing.T) {
	options := test.NewChannelTestOptions()
	
	// Create test dependencies  
	lifecycleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectGroupName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	channel := channels.Channel{
		Name:        acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		Description: acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	}

	resourceName := fmt.Sprintf("octopusdeploy_channel.%s", options.LocalName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		CheckDestroy:             testChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testChannelBasic(options.LocalName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName, channel),
				Check: resource.ComposeTestCheckFunc(
					testChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", channel.Name),
					resource.TestCheckResourceAttr(resourceName, "description", channel.Description),
				),
			},
		},
	})
}

func testChannelExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		channelID, err := getChannelID(s, resourceName)
		if err != nil {
			return err
		}

		if _, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), channelID); err != nil {
			return fmt.Errorf("channel %s not found", channelID)
		}

		return nil
	}
}

func testChannelDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_channel" {
			continue
		}

		if _, err := channels.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
			return fmt.Errorf("channel %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func getChannelID(s *terraform.State, resourceName string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("Not found: %s", resourceName)
	}

	return rs.Primary.ID, nil
}

func testChannelBasic(localName, lifecycleLocalName, lifecycleName, projectGroupLocalName, projectGroupName, projectLocalName, projectName string, channel channels.Channel) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_lifecycle" "%s" {
			name = "%s"
		}

		resource "octopusdeploy_project_group" "%s" {
			name = "%s"
		}

		resource "octopusdeploy_project" "%s" {
			lifecycle_id     = octopusdeploy_lifecycle.%s.id
			name             = "%s"
			project_group_id = octopusdeploy_project_group.%s.id
		}

		resource "octopusdeploy_channel" "%s" {
			name        = "%s"
			description = "%s"
			project_id  = octopusdeploy_project.%s.id
		}
	`,
		lifecycleLocalName, lifecycleName,
		projectGroupLocalName, projectGroupName,
		projectLocalName, lifecycleLocalName, projectName, projectGroupLocalName,
		localName, channel.Name, channel.Description, projectLocalName,
	)
}