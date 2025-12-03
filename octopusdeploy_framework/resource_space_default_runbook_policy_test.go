package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOctopusDeploySpaceDefaultRunbookRetentionPolicy_CreateForever(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_runbook_retention_policy.policy_" + localName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: spaceDefaultRunbookRetentionPolicy_Forever(localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "strategy", "Forever"),
				),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultRunbookRetentionPolicy_ModifiesBetweenStrategies(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_runbook_retention_policy.policy_" + localName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testSpaceCheckDestroy,
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: spaceDefaultRunbookRetentionPolicy_Count(localName, 5, "Days"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "strategy", "Count"),
					resource.TestCheckResourceAttr(resourceName, "quantity_to_keep", "5"),
					resource.TestCheckResourceAttr(resourceName, "unit", "Days"),
				)},
			{
				Config: spaceDefaultRunbookRetentionPolicy_Forever(localName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "strategy", "Forever"),
				),
			},
			{
				Config: spaceDefaultRunbookRetentionPolicy_Count(localName, 5, "Days"),
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

func TestAccOctopusDeploySpaceDefaultRunbookRetentionPolicy_CreateCountWithMissingFields_HasException(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      invalidSpaceDefaultRunbookRetentionPolicy_CountWithoutUnit(localName),
				ExpectError: regexp.MustCompile(`Missing Required Field`),
			},
			{
				Config:      invalidSpaceDefaultRunbookRetentionPolicy_CountWithoutQuantityToKeep(localName),
				ExpectError: regexp.MustCompile(`Missing Required Field`),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultRunbookRetentionPolicy_CreateForeverWithCountAttributes_ThrowsException(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      invalidSpaceDefaultRunbookRetentionPolicy_ForeverWithUnit(localName),
				ExpectError: regexp.MustCompile(`Invalid Field`),
			},
			{
				Config:      invalidSpaceDefaultRunbookRetentionPolicy_ForeverWithQuantityToKeep(localName),
				ExpectError: regexp.MustCompile(`Invalid Field`),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultRunbookRetentionPolicy_CreateCountWithInvalidQuantities_ThrowsException(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      spaceDefaultRunbookRetentionPolicy_Count(localName, -1, "Days"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
			{
				Config:      spaceDefaultRunbookRetentionPolicy_Count(localName, 0, "Days"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultRunbookRetentionPolicy_CreateCountWithInvalidUnits_ThrowsException(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      spaceDefaultRunbookRetentionPolicy_Count(localName, 1, "Years"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultRunbookRetentionPolicy_CreateWithInvalidStrategy_ThrowsException(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      spaceDefaultRunbookRetentionPolicy(localName, "deleteAtRandom", nil, nil),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
		},
	})
}

func TestAccOctopusDeploySpaceDefaultRunbookRetentionPolicy_Delete(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_space_default_runbook_retention_policy.policy_" + localName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		PreCheck:                 func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: spaceDefaultRunbookRetentionPolicy_Forever(localName),
			},
			{
				Config: noSpaceDefaultRunbookRetentionPolicy(localName),
				Check: resource.ComposeTestCheckFunc(
					testNoSpaceDefaultRunbookRetentionPolicyResource(resourceName),
				),
			},
		},
	})
}

func spaceDefaultRunbookRetentionPolicy(localName string, strategy string, quantityToKeep *int64, unit *string) string {
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
			description           = "Test space for runbook resource"
			space_managers_teams  = ["teams-administrators"]
		}
			
		resource "octopusdeploy_space_default_runbook_retention_policy" "policy_%[1]s" {
			space_id = octopusdeploy_space.space_%[1]s.id
			strategy = "%[2]s"
			%[3]s
			%[4]s
		}
	`, localName, strategy, quantityToKeepStr, unitStr)
}

func spaceDefaultRunbookRetentionPolicy_Forever(localName string) string {
	return spaceDefaultRunbookRetentionPolicy(localName, "Forever", nil, nil)
}

func spaceDefaultRunbookRetentionPolicy_Count(localName string, quantityToKeep int64, unit string) string {
	return spaceDefaultRunbookRetentionPolicy(localName, "Count", &quantityToKeep, &unit)
}

func noSpaceDefaultRunbookRetentionPolicy(localName string) string {
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

func invalidSpaceDefaultRunbookRetentionPolicy_CountWithoutUnit(localName string) string {
	return spaceDefaultRunbookRetentionPolicy(localName, "Count", nil, nil)
}

func invalidSpaceDefaultRunbookRetentionPolicy_CountWithoutQuantityToKeep(localName string) string {
	return spaceDefaultRunbookRetentionPolicy(localName, "Count", nil, nil)
}

func invalidSpaceDefaultRunbookRetentionPolicy_ForeverWithUnit(localName string) string {
	var unit string = "Days"
	return spaceDefaultRunbookRetentionPolicy(localName, "Forever", nil, &unit)
}

func invalidSpaceDefaultRunbookRetentionPolicy_ForeverWithQuantityToKeep(localName string) string {
	var quantityToKeep int64 = 5
	return spaceDefaultRunbookRetentionPolicy(localName, "Forever", &quantityToKeep, nil)
}

func testNoSpaceDefaultRunbookRetentionPolicyResource(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if ok {
			return fmt.Errorf("resource %s still exists", resourceName)
		}
		return nil
	}
}
