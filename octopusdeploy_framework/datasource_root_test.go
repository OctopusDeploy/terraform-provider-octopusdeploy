package octopusdeploy_framework

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceRoot(t *testing.T) {

	resourceName := "data.octopusdeploy_root.root"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRootConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "server_version", "2026.1.123"),
				),
			},
		},
	})
}

func testAccDataSourceRootConfig() string {
	return `
data "octopusdeploy_root" "root" {
}
`
}
