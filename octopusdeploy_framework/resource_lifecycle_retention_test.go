package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLifecycleRetentionUpdates(t *testing.T) {

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

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},

			// 7 set new retention block to default
			{
				Config: lifecycle_newRetention(lifecycleName, "Default", "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.strategy", "Default"),
				),
			},
			// 8 set new retention block to forever
			{
				Config: lifecycle_newRetention(lifecycleName, "Forever", "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Forever"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.0.strategy", "Forever"),
				)},
			// 9 set new retention block to Count retention using days
			{
				Config: lifecycle_newRetention(lifecycleName, "Count", "1", "Days", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),

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
				Config: lifecycle_newRetention(lifecycleName, "Count", "1", "Items", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),

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
				Config: lifecycle_newReleaseRetention(lifecycleName, "Default", "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),

					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_with_strategy.0.strategy", "Default"),

					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_with_strategy.#", "0"),
				),
			},
			// 12 change new retention block to count only for release
			{
				Config: lifecycle_newReleaseRetention(lifecycleName, "Count", "3", "Items", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),

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

func TestAccRetentionAttributeValidation(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      lifecycle_newRetention(lifecycleName, "Default", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must not be set when strategy is Forever or Default.*quantity_to_keep must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycle_newRetention(lifecycleName, "Forever", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must not be set when strategy is Forever or Default.*quantity_to_keep must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycle_newRetention(lifecycleName, "Default", "", "days", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must not be set when strategy is Forever or Default.*unit must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycle_newRetention(lifecycleName, "Forever", "", "items", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must not be set when strategy is Forever or Default.*unit must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycle_newRetention(lifecycleName, "", "", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)The argument "strategy" is required, but no definition was found.*The argument "strategy" is required, but no definition was found.`),
			},
			{
				Config:      lifecycle_newRetention(lifecycleName, "Count", "1", "Days", "true"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)An argument named "should_keep_forever" is not expected here.*An argument named "should_keep_forever" is not expected here.`),
			},
			{
				Config:      lifecycle_newRetention(lifecycleName, "Count", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must be set when strategy is set to Count.*unit must be set when strategy is set to Count.`),
			},
			{
				Config:      lifecycle_newRetention(lifecycleName, "Count", "", "Days", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must be set when strategy is set to Count.*quantity_to_keep must be set when strategy is set to Count`),
			},
			{
				Config:      lifecycle_retentionWithoutStrategy_count(lifecycleName, "Count"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)Blocks of type "release_retention_policy" are not expected here.*Blocks of type "tentacle_retention_policy" are not expected here`),
			},
		},
	})
}

func TestAccLifecycleWithPhaseInheritingRetentions(t *testing.T) {
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
					resource.TestCheckResourceAttr(lifecycleResource, "phase.0.release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "phase.0.tentacle_retention_with_strategy.#", "0"),
				),
			},
		},
	})
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

func lifecycle_newRetention(lifecycleName string, strategy string, quantityToKeep string, unit string, shouldKeepForever string) string {
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

func lifecycle_newReleaseRetention(lifecycleName string, strategy string, quantityToKeep string, unit string, shouldKeepForever string) string {
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
