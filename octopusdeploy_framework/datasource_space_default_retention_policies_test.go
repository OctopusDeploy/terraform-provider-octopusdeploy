package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSpaceDefaultRetentionPolicy(t *testing.T) {
	spaceName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	retentionType := "LifecycleRelease"
	resourceName := "data.octopusdeploy_space_default_retention_policy.space_default_retention_policy_test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceDefaultRetentionPolicyConfig(spaceName, retentionType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_type", retentionType),
					resource.TestCheckResourceAttrSet(resourceName, "strategy"),
				),
			},
		},
	})
}

func testAccDataSourceSpaceDefaultRetentionPolicyConfig(spaceName, retention_type string) string {
	return fmt.Sprintf(`
resource "octopusdeploy_space" "octopus_project_space_test" {
  name                  = "%s"
  is_default            = false
  is_task_queue_stopped = false
  description           = "Test space for lifecycles datasource"
  space_managers_teams  = ["teams-administrators"]
}

data "octopusdeploy_space_default_retention_policy" "space_default_retention_policy_test" {
  space_id       = octopusdeploy_space.octopus_project_space_test.id
  retention_type = "%s"
}
`, spaceName, retention_type)
}
