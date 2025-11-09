package octopusdeploy_framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeploySpaceDefaultLifecycleReleaseRetentionPolicy_CreateForever(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_lifecycle_release_retention_policy.policy_" + localName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: spaceDefaultLifecycleRetentionPolicy_Forever(localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "strategy", "Forever"),
				),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultLifecycleReleaseRetentionPolicy_CreateCountAndModifyToForever(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_lifecycle_release_retention_policy.policy_" + localName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testSpaceCheckDestroy,
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: spaceDefaultLifecycleRetentionPolicy_Count(localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "strategy", "Count"),
					resource.TestCheckResourceAttr(resourceName, "quantity_to_keep", "5"),
					resource.TestCheckResourceAttr(resourceName, "unit", "Days"),
				),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultLifecycleReleaseRetentionPolicy_ModifyStrategyFromCountToForever(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_lifecycle_release_retention_policy.policy_" + localName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: spaceDefaultLifecycleRetentionPolicy_Count(localName),
			},
			{
				Config: spaceDefaultLifecycleRetentionPolicy_Forever(localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "strategy", "Forever"),
				),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultLifecycleReleaseRetentionPolicy_ModifyStrategyFromForeverToCount(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_lifecycle_release_retention_policy.policy_" + localName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: spaceDefaultLifecycleRetentionPolicy_Forever(localName),
			},
			{
				Config: spaceDefaultLifecycleRetentionPolicy_Count(localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "strategy", "Count"),
					resource.TestCheckResourceAttr(resourceName, "quantity_to_keep", "5"),
					resource.TestCheckResourceAttr(resourceName, "unit", "Days"),
				),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultLifecycleReleaseRetentionPolicy_Delete(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_lifecycle_release_retention_policy.policy_" + localName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: spaceDefaultLifecycleRetentionPolicy_Forever(localName),
			},
			{
				Config: noSpaceDefaultLifecycleReleaseRetentionPolicy(localName),
				Check: resource.ComposeTestCheckFunc(
					testNoResource(resourceName),
				),
			},
		},
	})
}

func spaceDefaultLifecycleRetentionPolicy_Forever(localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_space" "space_%[1]s" {
			name                  = "%[1]s"
			is_default            = false
			is_task_queue_stopped = false
			description           = "Test space for lifecycles datasource"
			space_managers_teams  = ["teams-administrators"]
		}

		resource "octopusdeploy_space_default_lifecycle_release_retention_policy" "policy_%[1]s" {
			space_id = octopusdeploy_space.space_%[1]s.id
			strategy = "Forever"
		}
	`, localName)
}

func spaceDefaultLifecycleRetentionPolicy_Count(localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_space" "space_%[1]s" {
			name                  = "%[1]s"
			is_default            = false
			is_task_queue_stopped = false
			description           = "Test space for lifecycles datasource"
			space_managers_teams  = ["teams-administrators"]
		}

		resource "octopusdeploy_space_default_lifecycle_release_retention_policy" "policy_%[1]s" {
			space_id = octopusdeploy_space.space_%[1]s.id
			strategy = "Count"
			quantity_to_keep = 5
			unit = "Days"
		}
	`, localName)
}

func noSpaceDefaultLifecycleReleaseRetentionPolicy(localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_space" "space_%[1]s" {
			name                  = "%[1]s"
			is_default            = false
			is_task_queue_stopped = false
			description           = "Test space for lifecycles datasource"
			space_managers_teams  = ["teams-administrators"]
		}
	`, localName)
}

func testNoResource(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if ok {
			return fmt.Errorf("resource %s still exists", resourceName)
		}
		return nil
	}
}
