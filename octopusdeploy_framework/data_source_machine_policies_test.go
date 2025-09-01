package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployDataSourceMachinePolicies(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_machine_policies.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMachinePoliciesConfig(localName, take),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMachinePoliciesDataSourceID(name),
					resource.TestCheckResourceAttr(name, "take", fmt.Sprintf("%d", take)),
					resource.TestCheckResourceAttrSet(name, "machine_policies.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDataSourceMachinePoliciesWithFilters(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_machine_policies.%s", localName)
	partialName := "default"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMachinePoliciesWithFiltersConfig(localName, partialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMachinePoliciesDataSourceID(name),
					resource.TestCheckResourceAttr(name, "partial_name", partialName),
					resource.TestCheckResourceAttrSet(name, "machine_policies.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDataSourceMachinePoliciesWithSpaceId(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_machine_policies.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMachinePoliciesWithSpaceIdConfig(localName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMachinePoliciesDataSourceID(name),
					resource.TestCheckResourceAttr(name, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttrSet(name, "machine_policies.#"),
				),
			},
		},
	})
}

func testAccCheckMachinePoliciesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find MachinePolicies data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot MachinePolicies source ID not set")
		}
		return nil
	}
}

func testAccDataSourceMachinePoliciesConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_machine_policies" "%s" {
		take = %v
	}`, localName, take)
}

func testAccDataSourceMachinePoliciesWithFiltersConfig(localName, partialName string) string {
	return fmt.Sprintf(`data "octopusdeploy_machine_policies" "%s" {
		partial_name = "%s"
		take = 20
	}`, localName, partialName)
}

func testAccDataSourceMachinePoliciesWithSpaceIdConfig(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_machine_policies" "%s" {
		space_id = "Spaces-1"
		take = 10
	}`, localName)
}