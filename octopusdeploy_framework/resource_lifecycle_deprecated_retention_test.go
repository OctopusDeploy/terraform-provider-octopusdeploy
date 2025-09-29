package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

// replace with deprecation names
func TestAccLifecycleDeprecatedRetentionPolicyUpdates(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleResource := "octopusdeploy_lifecycle." + lifecycleName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1 create lifecycle without a retention policy
			{
				Config: lifecycleBasic(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},

			// 2 update with default retention policies
			{
				Config: defaultDepreciatedRetentionLifecycle_usingQuantityToKeep(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			// 3 update with Count retention policies using days
			{
				Config: countDepreciatedRetentionLifecycle(lifecycleName, "Days"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			// 3 update with Count retention policies using items
			{
				Config: countDepreciatedRetentionLifecycle(lifecycleName, "Items"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			// 4 update with Default retention policies
			{
				Config: defaultDepreciatedRetentionLifecycle_notUsingQuantityToKeep(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
		},
	})
}

func TestAccRetentionDepreciatedAttributeValidation(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      lifecycleGivenDepreciatedRetentionAttributes(lifecycleName, "", "1", "Items", "true"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`should_keep_forever must be false when quantity_to_keep is not 0`),
			},
			{
				Config:      lifecycleGivenDepreciatedRetentionAttributes(lifecycleName, "", "1", "", "true"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`should_keep_forever must be false when quantity_to_keep is not 0`),
			},
			{
				Config:      lifecycleGivenDepreciatedRetentionAttributes(lifecycleName, "", "0", "", "false"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`should_keep_forever must be true when quantity_to_keep is 0`),
			},
			{
				Config:      lifecycleGivenDepreciatedRetentionAttributes(lifecycleName, "", "", "", "false"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The non-refresh plan was not empty`),
			},
			{
				Config:      lifecycleGivenDepreciatedRetentionAttributes(lifecycleName, "", "", "Items", "false"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The non-refresh plan was not empty`),
			},
			{
				Config:      lifecycleGivenDepreciatedRetentionAttributes(lifecycleName, "", "", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The non-refresh plan was not empty`),
			},
			{
				Config:      lifecycleGivenDepreciatedRetentionAttributes(lifecycleName, "Default", "", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`An argument named "strategy" is not expected here.`),
			},
			{
				Config:      lifecycleWithTwoTypesOfRetention(lifecycleName),
				PlanOnly:    false,
				ExpectError: regexp.MustCompile(`Both release_retention_with_strategy and release_retention_policy are used.`),
			},
			{
				Config:      lifecycleWithTwoTypesOfRetentionIncludingPhases(lifecycleName),
				PlanOnly:    false,
				ExpectError: regexp.MustCompile(`Both release_retention_with_strategy and release_retention_policy are used.`),
			},
		},
	})
}

func TestAccLifecycle_WithPhase_InheritingDepreciatedRetentionPolicies(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	phaseName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleResource := "octopusdeploy_lifecycle." + lifecycleName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "phase.#", "1"),
					//check that the phase retention policies remain empty so will inherit their policies from elsewhere
					resource.TestCheckResourceAttr(lifecycleResource, "phase.0.release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "phase.0.tentacle_retention_policy.#", "0"),
				),

				Config: lifecycle_withBasicPhase(lifecycleName, phaseName),
			},
		},
	})
}

func countDepreciatedRetentionLifecycle(lifecycleName string, unit string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			should_keep_forever = "false"
			quantity_to_keep    = "1"
			unit                = "%s"
		}
		tentacle_retention_policy {
			quantity_to_keep = "1"
			unit             = "%s"
		}
    }`, lifecycleName, lifecycleName, unit, unit)
}

func defaultDepreciatedRetentionLifecycle_usingQuantityToKeep(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			quantity_to_keep    = "0"
			should_keep_forever = "true"
			unit                = "Items"	
		}
		tentacle_retention_policy {
			should_keep_forever = "true"
			quantity_to_keep = "0"
		}
    }`, lifecycleName, lifecycleName)
}

func defaultDepreciatedRetentionLifecycle_notUsingQuantityToKeep(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			should_keep_forever = "true"
		}
		tentacle_retention_policy {
			should_keep_forever = "true"
			unit = "Items"
		}
    }`, lifecycleName, lifecycleName)
}

func lifecycleGivenDepreciatedRetentionAttributes(lifecycleName string, strategy string, quantityToKeep string, unit string, shouldKeepForever string) string {
	var strategyAttribute string
	if strategy != "" {
		strategyAttribute = fmt.Sprintf(`strategy = "%s"`, strategy)
	}

	var quantityToKeepAttribute string
	if quantityToKeep != "" {
		quantityToKeepAttribute = fmt.Sprintf(`quantity_to_keep = "%s"`, quantityToKeep)
	}
	var shouldKeepForeverAttribute string
	if shouldKeepForever != "" {
		shouldKeepForeverAttribute = fmt.Sprintf(`should_keep_forever = "%s"`, shouldKeepForever)
	}
	var unitAttribute string
	if unit != "" {
		unitAttribute = fmt.Sprintf(`unit = "%s"`, unit)
	}

	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
    	release_retention_policy{
			%s
    		%s
			%s
			%s
  		}
		tentacle_retention_policy {
			%s
    		%s
			%s	
			%s
  		}
	}`, lifecycleName, lifecycleName, strategyAttribute, quantityToKeepAttribute, shouldKeepForeverAttribute, unitAttribute, strategyAttribute, quantityToKeepAttribute, shouldKeepForeverAttribute, unitAttribute)
}

func lifecycleWithTwoTypesOfRetention(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			should_keep_forever = "true"
		}
		tentacle_retention_with_strategy{
			strategy = "Forever"
		}
    }`, lifecycleName, lifecycleName)
}

func lifecycleWithTwoTypesOfRetentionIncludingPhases(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			should_keep_forever = "true"
		}
		phase{
   			name = "Phase1"
			release_retention_with_strategy {
				strategy = "Forever"
			}		
		}
	
    }`, lifecycleName, lifecycleName)
}
