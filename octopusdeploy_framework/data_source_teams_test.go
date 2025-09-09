package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployDataSourceTeams(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_teams.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeamsConfig(localName, take),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamsDataSourceID(name),
					resource.TestCheckResourceAttr(name, "take", fmt.Sprintf("%d", take)),
					resource.TestCheckResourceAttrSet(name, "teams.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDataSourceTeamsWithFilters(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_teams.%s", localName)
	partialName := "test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeamsWithFiltersConfig(localName, partialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamsDataSourceID(name),
					resource.TestCheckResourceAttr(name, "partial_name", partialName),
					resource.TestCheckResourceAttrSet(name, "teams.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDataSourceTeamsIncludeSystem(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_teams.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeamsIncludeSystemConfig(localName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamsDataSourceID(name),
					resource.TestCheckResourceAttr(name, "include_system", "true"),
					resource.TestCheckResourceAttrSet(name, "teams.#"),
				),
			},
		},
	})
}

func testAccCheckTeamsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find Teams data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot Teams source ID not set")
		}
		return nil
	}
}

func testAccDataSourceTeamsConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_teams" "%s" {
		take = %v
	}`, localName, take)
}

func testAccDataSourceTeamsWithFiltersConfig(localName, partialName string) string {
	return fmt.Sprintf(`data "octopusdeploy_teams" "%s" {
		partial_name = "%s"
		take = 20
	}`, localName, partialName)
}

func testAccDataSourceTeamsIncludeSystemConfig(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_teams" "%s" {
		include_system = true
		take = 10
	}`, localName)
}