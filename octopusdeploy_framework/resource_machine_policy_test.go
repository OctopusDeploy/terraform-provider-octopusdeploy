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
					resource.TestCheckResourceAttr(prefix, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttr(prefix, "is_default", "false"),
				),
				Config: testAccMachinePolicyBasic(localName, name),
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
				Config: testAccMachinePolicyWithDescription(localName, name, description),
			},
			{
				Check: resource.ComposeTestCheckFunc(
					testAccMachinePolicyExists(prefix),
					resource.TestCheckResourceAttr(prefix, "name", newName),
					resource.TestCheckResourceAttr(prefix, "description", newDescription),
				),
				Config: testAccMachinePolicyWithDescription(localName, newName, newDescription),
			},
		},
	})
}

func TestAccOctopusDeployMachinePolicyWithConnectivitySettings(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

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
				Config: testAccMachinePolicyWithConnectivitySettings(localName, name),
			},
		},
	})
}

func TestAccOctopusDeployMachinePolicyImport(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_machine_policy." + localName
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccMachinePolicyCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccMachinePolicyBasic(localName, name),
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

func testAccMachinePolicyBasic(localName, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name     = "%s"
		space_id = "Spaces-1"
	}`, localName, name)
}

func testAccMachinePolicyWithDescription(localName, name, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name        = "%s"
		description = "%s"
		space_id    = "Spaces-1"
	}`, localName, name, description)
}

func testAccMachinePolicyWithConnectivitySettings(localName, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_machine_policy" "%s" {
		name                           = "%s"
		space_id                       = "Spaces-1"
		connection_connect_timeout     = 60000000000
		connection_retry_count_limit   = 5
		connection_retry_sleep_interval = 1000000000
		connection_retry_time_limit     = 300000000000
	}`, localName, name)
}

func testAccMachinePolicyExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		machinePolicyID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := machinepolicies.GetByID(octoClient, octoClient.GetSpaceID(), machinePolicyID); err != nil {
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

		if machinePolicy, err := machinepolicies.GetByID(octoClient, octoClient.GetSpaceID(), rs.Primary.ID); err == nil {
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
