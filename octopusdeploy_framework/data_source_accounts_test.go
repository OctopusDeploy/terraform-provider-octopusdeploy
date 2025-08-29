package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployDataSourceAccounts(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_accounts.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsConfig(localName, take),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccountsDataSourceID(name),
					resource.TestCheckResourceAttr(name, "take", fmt.Sprintf("%d", take)),
					resource.TestCheckResourceAttrSet(name, "accounts.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDataSourceAccountsWithFilters(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_accounts.%s", localName)
	partialName := "test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsWithFiltersConfig(localName, partialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccountsDataSourceID(name),
					resource.TestCheckResourceAttr(name, "partial_name", partialName),
					resource.TestCheckResourceAttrSet(name, "accounts.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDataSourceAccountsWithSpaceId(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_accounts.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountsWithSpaceIdConfig(localName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccountsDataSourceID(name),
					resource.TestCheckResourceAttr(name, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttrSet(name, "accounts.#"),
				),
			},
		},
	})
}

func testAccCheckAccountsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find Accounts data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot Accounts source ID not set")
		}
		return nil
	}
}

func testAccDataSourceAccountsConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_accounts" "%s" {
		take = %v
	}`, localName, take)
}

func testAccDataSourceAccountsWithFiltersConfig(localName, partialName string) string {
	return fmt.Sprintf(`data "octopusdeploy_accounts" "%s" {
		partial_name = "%s"
		take = 20
	}`, localName, partialName)
}

func testAccDataSourceAccountsWithSpaceIdConfig(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_accounts" "%s" {
		space_id = "Spaces-1"
		take = 10
	}`, localName)
}