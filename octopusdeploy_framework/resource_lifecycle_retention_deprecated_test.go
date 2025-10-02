package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccDeprecatedLifecycleRetentionUpdates(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleResource := "octopusdeploy_lifecycle." + lifecycleName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1 create lifecycle without retention settings
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
				Config: deprecatedNewRetentionLifecycle(lifecycleName, "Default", "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.strategy", "Default"),
				),
			},
			//	3 update with Forever retention policies
			{
				Config: deprecatedNewRetentionLifecycle(lifecycleName, "Forever", "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Forever"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.strategy", "Forever"),
				)},
			// 4 update with Count retention policies using days
			{
				Config: deprecatedNewRetentionLifecycle(lifecycleName, "Count", "1", "Days", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.unit", "Days"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.unit", "Days"),
				),
			},
			// 5 update with Count retention policies using items
			{
				Config: deprecatedNewRetentionLifecycle(lifecycleName, "Count", "1", "Items", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.unit", "Items"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.unit", "Items"),
				),
			},
			// 6 only set release retention
			{
				Config: deprecatedNewRetentionLifeycleWithOnlyRelease(lifecycleName, "Default", "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			// 7 change only release retention to count
			{
				Config: deprecatedNewRetentionLifeycleWithOnlyRelease(lifecycleName, "Count", "3", "Items", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.quantity_to_keep", "3"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.unit", "Items"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			//8 set old style retention Block to Forever
			{
				Config: deprecatedOldRetentionLifecycle(lifecycleName, "", "", "true"),
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
			// 9 update with Count retention policies using days
			{
				Config: deprecatedOldRetentionLifecycle(lifecycleName, "1", "Days", "false"),
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
			// 10 update with Count retention policies using items
			{
				Config: deprecatedOldRetentionLifecycle(lifecycleName, "1", "Items", "false"),
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
		},
	})
}

func TestAccDeprecatedRetentionAttributeValidation(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      deprecatedNewRetentionLifecycle(lifecycleName, "Default", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must not be set when strategy is Forever or Default.*quantity_to_keep must not be set when strategy is Forever or Default`),
			},
			{
				Config:      deprecatedNewRetentionLifecycle(lifecycleName, "Forever", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must not be set when strategy is Forever or Default.*quantity_to_keep must not be set when strategy is Forever or Default`),
			},
			{
				Config:      deprecatedNewRetentionLifecycle(lifecycleName, "Default", "", "days", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must not be set when strategy is Forever or Default.*unit must not be set when strategy is Forever or Default`),
			},
			{
				Config:      deprecatedNewRetentionLifecycle(lifecycleName, "Forever", "", "items", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must not be set when strategy is Forever or Default.*unit must not be set when strategy is Forever or Default`),
			},
			{
				Config:      deprecatedNewRetentionLifecycle(lifecycleName, "", "", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)The argument "strategy" is required, but no definition was found.*The argument "strategy" is required, but no definition was found.`),
			},
			{
				Config:      deprecatedNewRetentionLifecycle(lifecycleName, "Count", "1", "Days", "true"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)An argument named "should_keep_forever" is not expected here.*An argument named "should_keep_forever" is not expected here.`),
			},
			{
				Config:      deprecatedNewRetentionLifecycle(lifecycleName, "Count", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must be set when strategy is set to Count.*unit must be set when strategy is set to Count.`),
			},
			{
				Config:      deprecatedNewRetentionLifecycle(lifecycleName, "Count", "", "Days", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must be set when strategy is set to Count.*quantity_to_keep must be set when strategy is set to Count`),
			},
			//Using Old retention Blocks
			// when quantity_to_keep is > 0 should_keep_forever shouldn't be true
			{
				Config:      deprecatedOldRetentionLifecycle(lifecycleName, "1", "", "true"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`should_keep_forever must be false when quantity_to_keep is not 0`),
			},
			{
				Config:      deprecatedOldRetentionLifecycle(lifecycleName, "1", "", "true"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`should_keep_forever must be false when quantity_to_keep is not 0`),
			},
			// when quantity_to_keep is 0, should_keep_forever shouldn't be false
			{
				Config:      deprecatedOldRetentionLifecycle(lifecycleName, "0", "", "false"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`should_keep_forever must be true when quantity_to_keep is 0`),
			},
			{
				Config:      deprecatedOldRetentionLifecycle(lifecycleName, "", "", "false"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The non-refresh plan was not empty`),
			},
			{
				Config:      deprecatedOldRetentionLifecycle(lifecycleName, "", "Items", "false"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The non-refresh plan was not empty`),
			},
			{
				Config:      deprecatedOldRetentionLifecycle(lifecycleName, "", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The non-refresh plan was not empty`),
			},
		},
	})
}

func TestAccDeprecatedLifecycle_WithPhase_InheritingRetentions(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	phaseName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleResource := "octopusdeploy_lifecycle." + lifecycleName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: deprecatedRetentionLifecycle_withBasicPhase(lifecycleName, phaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "phase.#", "1"),
					//check that the phase retention policies remain empty so will inherit their policies from elsewhere
					resource.TestCheckResourceAttr(lifecycleResource, "phase.0.release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "phase.0.tentacle_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "phase.0.release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "phase.0.tentacle_retention_with_strategy.#", "0"),
				),
			},
		},
	})
}

func deprecatedNewRetentionLifecycle(lifecycleName string, strategy string, quantityToKeep string, unit string, shouldKeepForever string) string {
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
	resource := fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
    	release_retention_with_strategy {
			%s
    		%s
			%s
			%s
  		}
		tentacle_retention_with_strategy {
			%s
    		%s
			%s	
			%s
  		}

	}`, lifecycleName, lifecycleName, strategyAttribute, quantityToKeepAttribute, unitAttribute, shouldKeepForeverAttribute, strategyAttribute, quantityToKeepAttribute, unitAttribute, shouldKeepForeverAttribute)

	return resource
}

func deprecatedOldRetentionLifecycle(lifecycleName string, quantityToKeep string, unit string, shouldKeepForever string) string {
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
	resource := fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
    	release_retention_policy {
			%s
    		%s
			%s
  		}
		tentacle_retention_policy {
			%s
    		%s
			%s
  		}

	}`, lifecycleName, lifecycleName, quantityToKeepAttribute, unitAttribute, shouldKeepForeverAttribute, quantityToKeepAttribute, unitAttribute, shouldKeepForeverAttribute)

	return resource
}

func deprecatedRetentionLifecycle_withBasicPhase(lifecycleName string, phaseName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
  		phase {
    		name = "%s"
  		}
	}`, lifecycleName, lifecycleName, phaseName)
}

func deprecatedNewRetentionLifeycleWithOnlyRelease(lifecycleName string, strategy string, quantityToKeep string, unit string, shouldKeepForever string) string {
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
	resource := fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
    	release_retention_with_strategy {
			%s
    		%s
			%s
			%s
  		}

	}`, lifecycleName, lifecycleName, strategyAttribute, quantityToKeepAttribute, unitAttribute, shouldKeepForeverAttribute)

	return resource
}
