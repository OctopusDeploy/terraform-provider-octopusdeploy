package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeployDataSourceWorkerPools(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_worker_pools.%s", localName)
	take := 10

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkerPoolsConfig(localName, take),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkerPoolsDataSourceID(name),
					resource.TestCheckResourceAttr(name, "take", fmt.Sprintf("%d", take)),
					resource.TestCheckResourceAttrSet(name, "worker_pools.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDataSourceWorkerPoolsWithFilters(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_worker_pools.%s", localName)
	partialName := "default"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkerPoolsWithFiltersConfig(localName, partialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkerPoolsDataSourceID(name),
					resource.TestCheckResourceAttr(name, "partial_name", partialName),
					resource.TestCheckResourceAttrSet(name, "worker_pools.#"),
				),
			},
		},
	})
}

func TestAccOctopusDeployDataSourceWorkerPoolsWithSpaceId(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("data.octopusdeploy_worker_pools.%s", localName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkerPoolsWithSpaceIdConfig(localName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkerPoolsDataSourceID(name),
					resource.TestCheckResourceAttr(name, "space_id", "Spaces-1"),
					resource.TestCheckResourceAttrSet(name, "worker_pools.#"),
				),
			},
		},
	})
}

func testAccCheckWorkerPoolsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("cannot find WorkerPools data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("snapshot WorkerPools source ID not set")
		}
		return nil
	}
}

func testAccDataSourceWorkerPoolsConfig(localName string, take int) string {
	return fmt.Sprintf(`data "octopusdeploy_worker_pools" "%s" {
		take = %v
	}`, localName, take)
}

func testAccDataSourceWorkerPoolsWithFiltersConfig(localName, partialName string) string {
	return fmt.Sprintf(`data "octopusdeploy_worker_pools" "%s" {
		partial_name = "%s"
		take = 20
	}`, localName, partialName)
}

func testAccDataSourceWorkerPoolsWithSpaceIdConfig(localName string) string {
	return fmt.Sprintf(`data "octopusdeploy_worker_pools" "%s" {
		space_id = "Spaces-1"
		take = 10
	}`, localName)
}