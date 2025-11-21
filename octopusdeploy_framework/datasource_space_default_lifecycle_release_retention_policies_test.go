package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSpaceDefaultLifecycleReleaseRetentionPolicy(t *testing.T) {
	spaceName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "data.octopusdeploy_space_default_lifecycle_release_retention_policy.policy_test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpaceDefaultLifecycleReleaseRetentionPolicyConfig(spaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttrSet(resourceName, "strategy"),
				),
			},
		},
	})
}

func testAccDataSourceSpaceDefaultLifecycleReleaseRetentionPolicyConfig(spaceName string) string {
	return fmt.Sprintf(`
resource "octopusdeploy_space" "octopus_project_space_test" {
  name                  = "%s"
  is_default            = false
  is_task_queue_stopped = false
  description           = "Test space for lifecycles datasource"
  space_managers_teams  = ["teams-administrators"]
}

data "octopusdeploy_space_default_lifecycle_release_retention_policy" "policy_test" {
  space_id       = octopusdeploy_space.octopus_project_space_test.id
}
`, spaceName)
}
