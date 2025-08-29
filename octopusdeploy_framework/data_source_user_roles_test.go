package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployDataSourceUserRoles(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_user_roles.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserRolesConfig(localName, take),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserRolesDataSourceID(name),
					resource.TestCheckResourceAttr(name, "take", fmt.Sprintf("%d", take)),
					resource.TestCheckResourceAttrSet(name, "user_roles.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDataSourceUserRolesWithFilters(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_user_roles.%s", localName)
	partialName := "admin"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserRolesWithFiltersConfig(localName, partialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserRolesDataSourceID(name),
					resource.TestCheckResourceAttr(name, "partial_name", partialName),
					resource.TestCheckResourceAttrSet(name, "user_roles.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDataSourceUserRolesWithSpaceId(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_user_roles.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserRolesWithSpaceIdConfig(localName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserRolesDataSourceID(name),
					resource.TestCheckResourceAttr(name, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttrSet(name, "user_roles.#"),
				),
			},
		},
	})
}

func testAccCheckUserRolesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find UserRoles data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot UserRoles source ID not set")
		}
		return nil
	}
}

func testAccDataSourceUserRolesConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_user_roles" "%s" {
		take = %v
	}`, localName, take)
}

func testAccDataSourceUserRolesWithFiltersConfig(localName, partialName string) string {
	return fmt.Sprintf(`data "octopusdeploy_user_roles" "%s" {
		partial_name = "%s"
		take = 20
	}`, localName, partialName)
}

func testAccDataSourceUserRolesWithSpaceIdConfig(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_user_roles" "%s" {
		space_id = "Spaces-1"
		take = 10
	}`, localName)
}