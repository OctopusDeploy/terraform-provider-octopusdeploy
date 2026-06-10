package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployMachinePolicyBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckResourceAttrSet(prefix, "id"),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "space_id", space.ID),
					resource.TestCheckResourceAttr(prefix, "is_default", "false"),
				),
				Config: testAccMachinePolicyBasic(localName, name, space.ID),
			},
		},
	})
}

func TestAccOctopusDeployMachinePolicyUpdate(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "description", description),
				),
				Config: testAccMachinePolicyWithDescription(localName, name, description, space.ID),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
				),
				Config: testAccMachinePolicyWithDescription(localName, newName, newDescription, space.ID),
			},
		},
	})
}

func TestAccOctopusDeployMachinePolicyWithConnectivitySettings(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "connection_connect_timeout", "60000000000"),
					resource.TestCheckResourceAttr(prefix, "connection_retry_count_limit", "5"),
					resource.TestCheckResourceAttr(prefix, "connection_retry_sleep_interval", "1000000000"),
					resource.TestCheckResourceAttr(prefix, "connection_retry_time_limit", "300000000000"),
				),
				Config: testAccMachinePolicyWithConnectivitySettings(localName, name, space.ID),
			},
		},
	})
}

func TestAccOctopusDeployMachinePolicyImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	space := NewTestSpace(t)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccMachinePolicyBasic(localName, name, space.ID),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMachinePolicyImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccMachinePolicyBasic(localName, name, spaceID string) string {
	return fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name     = "%s"
		space_id = "%s"
	}`, localName, name, spaceID)
}

func testAccMachinePolicyWithDescription(localName, name, description, spaceID string) string {
	return fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name        = "%s"
		description = "%s"
		space_id    = "%s"
	}`, localName, name, description, spaceID)
}

func testAccMachinePolicyWithConnectivitySettings(localName, name, spaceID string) string {
	return fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name                           = "%s"
		space_id                       = "%s"
		connection_connect_timeout     = 60000000000
		connection_retry_count_limit   = 5
		connection_retry_sleep_interval = 1000000000
		connection_retry_time_limit     = 300000000000
	}`, localName, name, spaceID)
}

func testAccMachinePolicyExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[prefix]
		machinePolicyID := rs.Primary.ID
		if _, err := machinepolicies.GetByID(octoClient, rs.Primary.Attributes["space_id"], machinePolicyID); err != nil {
			return err
		}

		return nil
	}
}

func testAccMachinePolicyCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_machine_policy" {
			continue
		}

		if machinePolicy, err := machinepolicies.GetByID(octoClient, rs.Primary.Attributes["space_id"], rs.Primary.ID); err == nil {
			return fmt.Errorf("machine policy (%s) still exists", machinePolicy.GetID())
		}
	}

	return nil
}

func testAccMachinePolicyImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
