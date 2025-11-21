package octopusdeploy_framework

import (
	"fmt"
	"regexp"
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
				Config: spaceDefaultLifecycleReleaseRetentionPolicy_Forever(localName),
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
				Config: spaceDefaultLifecycleReleaseRetentionPolicy_Count(localName, 5, "Days"),
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
				Config: spaceDefaultLifecycleReleaseRetentionPolicy_Count(localName, 5, "Days"),
			},
			{
				Config: spaceDefaultLifecycleReleaseRetentionPolicy_Forever(localName),
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
				Config: spaceDefaultLifecycleReleaseRetentionPolicy_Forever(localName),
			},
			{
				Config: spaceDefaultLifecycleReleaseRetentionPolicy_Count(localName, 5, "Days"),
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

func TestAccOctopusDeploySpaceDefaultLifecycleReleaseRetentionPolicy_CreateCountWithMissingFields(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      invalidSpaceDefaultLifecycleReleaseRetentionPolicy_Count(localName),
				ExpectError: regexp.MustCompile(`Missing Required Field`),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultLifecycleReleaseRetentionPolicy_CreateForeverWithCountAttributes(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      invalidSpaceDefaultLifecycleReleaseRetentionPolicy_Forever(localName),
				ExpectError: regexp.MustCompile(`Invalid Field`),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultLifecycleReleaseRetentionPolicy_CreateCountWithInvalidQuantity(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      spaceDefaultLifecycleReleaseRetentionPolicy_Count(localName, -1, "Days"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
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
				Config: spaceDefaultLifecycleReleaseRetentionPolicy_Forever(localName),
			},
			{
				Config: noSpaceDefaultLifecycleReleaseRetentionPolicy(localName),
				Check: resource.ComposeTestCheckFunc(
					testNoSpaceDefaultLifecycleReleaseRetentionPolicyResource(resourceName),
				),
			},
		},
	})
}

func spaceDefaultLifecycleReleaseRetentionPolicy(localName string, strategy string, quantityToKeep *int64, unit *string) string {
	var quantityToKeepStr, unitStr string
	if quantityToKeep != nil {
		quantityToKeepStr = fmt.Sprintf("quantity_to_keep = %d", *quantityToKeep)
	}
	if unit != nil {
		unitStr = fmt.Sprintf(`unit = "%s"`, *unit)
	}

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
			strategy = "%[2]s"
			%[3]s
			%[4]s
		}
	`, localName, strategy, quantityToKeepStr, unitStr)
}

func spaceDefaultLifecycleReleaseRetentionPolicy_Forever(localName string) string {
	return spaceDefaultLifecycleReleaseRetentionPolicy(localName, "Forever", nil, nil)
}

func spaceDefaultLifecycleReleaseRetentionPolicy_Count(localName string, quantityToKeep int64, unit string) string {
	return spaceDefaultLifecycleReleaseRetentionPolicy(localName, "Count", &quantityToKeep, &unit)
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

func invalidSpaceDefaultLifecycleReleaseRetentionPolicy_Count(localName string) string {
	return spaceDefaultLifecycleReleaseRetentionPolicy(localName, "Count", nil, nil)
}

func invalidSpaceDefaultLifecycleReleaseRetentionPolicy_Forever(localName string) string {
	var quantityToKeep int64 = 5
	var unit string = "Days"
	return spaceDefaultLifecycleReleaseRetentionPolicy(localName, "Forever", &quantityToKeep, &unit)
}

func testNoSpaceDefaultLifecycleReleaseRetentionPolicyResource(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if ok {
			return fmt.Errorf("resource %s still exists", resourceName)
		}
		return nil
	}
}
