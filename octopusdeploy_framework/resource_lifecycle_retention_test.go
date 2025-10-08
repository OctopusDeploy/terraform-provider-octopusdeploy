package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLifecycleRetentionUpdatesDEPRECATED(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleResource := "octopusdeploy_lifecycle." + lifecycleName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1 create lifecycle without retention settings
			{
				Config: lifecycle_noRetention(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			//2 update retentionWithoutStrategy to the default (forever) retention policies
			{
				Config: lifecycle_retentionWithoutStrategy_defaultUsingQuantityToKeepDEPRECATED(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),

					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			// 3 update retentionWithoutStrategy to Count retention policies using days
			{
				Config: lifecycle_retentionWithoutStrategy_count_DEPRECATED(lifecycleName, "Days"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Days"),

					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Days"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			// 4 update retentionWithoutStrategy to Count retention policies using items
			{
				Config: lifecycle_retentionWithoutStrategy_count_DEPRECATED(lifecycleName, "Items"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),

					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},

			// 5 update retentionWithoutStrategy to the default (forever) retention policies
			{
				Config: lifecycle_retentionWithoutStrategy_defaultNotUsingQuantityToKeepDEPRECATED(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),

					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},

			// 6 set retentionWithoutStrategy only for release
			{
				Config: lifecycle_ReleaseRetentionWithoutStrategy_DEPRECATED(lifecycleName, "1", "Days", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Days"),

					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			// 7 set new retention block to default
			{
				Config: lifecycle_newRetention_DEPRECATED(lifecycleName, "Default", "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.strategy", "Default"),
				),
			},
			// 8 set new retention block to forever
			{
				Config: lifecycle_newRetention_DEPRECATED(lifecycleName, "Forever", "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Forever"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.strategy", "Forever"),
				)},
			// 9 set new retention block to Count retention using days
			{
				Config: lifecycle_newRetention_DEPRECATED(lifecycleName, "Count", "1", "Days", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
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
			// 10 set new retention block to count using items
			{
				Config: lifecycle_newRetention_DEPRECATED(lifecycleName, "Count", "1", "Items", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
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
			// 11 set new retention block to default only for release
			{
				Config: lifecycle_newReleaseRetention_DEPRECATED(lifecycleName, "Default", "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Default"),

					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			// 12 change new retention block to count only for release
			{
				Config: lifecycle_newReleaseRetention_DEPRECATED(lifecycleName, "Count", "3", "Items", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.quantity_to_keep", "3"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.unit", "Items"),

					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
		},
	})
}

func TestAccRetentionAttributeValidationDEPRECATED(t *testing.T) {

	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      lifecycle_newRetention_DEPRECATED(lifecycleName, "Default", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must not be set when strategy is Forever or Default.*quantity_to_keep must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycle_newRetention_DEPRECATED(lifecycleName, "Forever", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must not be set when strategy is Forever or Default.*quantity_to_keep must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycle_newRetention_DEPRECATED(lifecycleName, "Default", "", "days", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must not be set when strategy is Forever or Default.*unit must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycle_newRetention_DEPRECATED(lifecycleName, "Forever", "", "items", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must not be set when strategy is Forever or Default.*unit must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycle_newRetention_DEPRECATED(lifecycleName, "", "", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)The argument "strategy" is required, but no definition was found.*The argument "strategy" is required, but no definition was found.`),
			},
			{
				Config:      lifecycle_newRetention_DEPRECATED(lifecycleName, "Count", "1", "Days", "true"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)An argument named "should_keep_forever" is not expected here.*An argument named "should_keep_forever" is not expected here.`),
			},
			{
				Config:      lifecycle_newRetention_DEPRECATED(lifecycleName, "Count", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must be set when strategy is set to Count.*unit must be set when strategy is set to Count.`),
			},
			{
				Config:      lifecycle_newRetention_DEPRECATED(lifecycleName, "Count", "", "Days", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must be set when strategy is set to Count.*quantity_to_keep must be set when strategy is set to Count`),
			},
			//Using Old retention Blocks without strategy
			// when quantity_to_keep is > 0 should_keep_forever shouldn't be true
			{
				Config:      lifecycle_retentionWithoutStrategy_DEPRECATED(lifecycleName, "1", "", "true"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`should_keep_forever must be false when quantity_to_keep is not 0`),
			},
			{
				Config:      lifecycle_retentionWithoutStrategy_DEPRECATED(lifecycleName, "1", "", "true"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`should_keep_forever must be false when quantity_to_keep is not 0`),
			},
			// when quantity_to_keep is 0, should_keep_forever shouldn't be false
			{
				Config:      lifecycle_retentionWithoutStrategy_DEPRECATED(lifecycleName, "0", "", "false"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`should_keep_forever must be true when quantity_to_keep is 0`),
			},
			{
				Config:      lifecycle_retentionWithoutStrategy_DEPRECATED(lifecycleName, "", "", "false"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The non-refresh plan was not empty`),
			},
			{
				Config:      lifecycle_retentionWithoutStrategy_DEPRECATED(lifecycleName, "", "Items", "false"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The non-refresh plan was not empty`),
			},
			{
				Config:      lifecycle_retentionWithoutStrategy_DEPRECATED(lifecycleName, "", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The non-refresh plan was not empty`),
			},
		},
	})
}

func TestAccLifecycleWithPhaseInheritingRetentionsDEPRECATED(t *testing.T) {

	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	phaseName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleResource := "octopusdeploy_lifecycle." + lifecycleName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: lifecycle_phaseAndNoRetention(lifecycleName, phaseName),
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

func lifecycle_retentionWithoutStrategy_count_DEPRECATED(lifecycleName string, unit string) string {
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

func lifecycle_retentionWithoutStrategy_defaultUsingQuantityToKeepDEPRECATED(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
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

func lifecycle_retentionWithoutStrategy_defaultNotUsingQuantityToKeepDEPRECATED(lifecycleName string) string {
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

func lifecycle_noRetention(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
}`, lifecycleName, lifecycleName)
}

func lifecycle_phaseAndNoRetention(lifecycleName string, phaseName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
  		phase {
    		name = "%s"
  		}
	}`, lifecycleName, lifecycleName, phaseName)
}

func lifecycle_retentionWithoutStrategy_DEPRECATED(lifecycleName string, quantityToKeep string, unit string, shouldKeepForever string) string {
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

func lifecycle_ReleaseRetentionWithoutStrategy_DEPRECATED(lifecycleName string, quantityToKeep string, unit string, shouldKeepForever string) string {
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

	}`, lifecycleName, lifecycleName, quantityToKeepAttribute, unitAttribute, shouldKeepForeverAttribute)

	return resource
}

func lifecycle_newRetention_DEPRECATED(lifecycleName string, strategy string, quantityToKeep string, unit string, shouldKeepForever string) string {
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

func lifecycle_newReleaseRetention_DEPRECATED(lifecycleName string, strategy string, quantityToKeep string, unit string, shouldKeepForever string) string {
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
