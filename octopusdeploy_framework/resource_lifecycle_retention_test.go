package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLifecycleRetentionUpdates(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks() {
		t.Skip("Skipping test because deprecated retention is used")
	}
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
				Config: defaultRetentionLifecycle(lifecycleName),
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
				Config: foreverRetentionLifecycle(lifecycleName),
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
				Config: countRetentionLifecycle(lifecycleName, "Days"),
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
				Config: countRetentionLifecycle(lifecycleName, "Items"),
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
				Config: lifecycleWithOnlyReleaseRetentionGivenAttributes(lifecycleName, "Default", "", "", ""),
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
				Config: lifecycleWithOnlyReleaseRetentionGivenAttributes(lifecycleName, "Count", "3", "Items", ""),
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
		},
	})
}

func TestAccRetentionAttributeValidation(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks() {
		t.Skip("Skipping test because deprecated retention is used")
	}
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "Default", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must not be set when strategy is Forever or Default.*quantity_to_keep must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "Forever", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must not be set when strategy is Forever or Default.*quantity_to_keep must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "Default", "", "days", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must not be set when strategy is Forever or Default.*unit must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "Forever", "", "items", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must not be set when strategy is Forever or Default.*unit must not be set when strategy is Forever or Default`),
			},
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "", "", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)The argument "strategy" is required, but no definition was found.*The argument "strategy" is required, but no definition was found.`),
			},
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "Count", "1", "Days", "true"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)An argument named "should_keep_forever" is not expected here.*An argument named "should_keep_forever" is not expected here.`),
			},
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "Count", "1", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)unit must be set when strategy is set to Count.*unit must be set when strategy is set to Count.`),
			},
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "Count", "", "Days", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep must be set when strategy is set to Count.*quantity_to_keep must be set when strategy is set to Count`),
			},
		},
	})
}

func TestAccLifecycle_WithPhase_InheritingRetentions(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks() {
		t.Skip("Skipping test because deprecated retention is used")
	}
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	phaseName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleResource := "octopusdeploy_lifecycle." + lifecycleName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: lifecycle_withBasicPhase(lifecycleName, phaseName),
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

func countRetentionLifecycle(lifecycleName string, unit string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_with_strategy {
			strategy 			= "Count"
			quantity_to_keep    = "1"
			unit                = "%s"
		}
		tentacle_retention_with_strategy {
			strategy 			= "Count"
			quantity_to_keep 	= "1"
			unit            	 = "%s"
		}
    }`, lifecycleName, lifecycleName, unit, unit)
}

func defaultRetentionLifecycle(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_with_strategy{
			strategy    = "Default"
		}
		tentacle_retention_with_strategy {
			strategy    = "Default"
		}
    }`, lifecycleName, lifecycleName)
}
func foreverRetentionLifecycle(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_with_strategy{
			strategy    = "Forever"
		}
		tentacle_retention_with_strategy {
			strategy    = "Forever"
		}
	}`, lifecycleName, lifecycleName)
}

func lifecycleGivenRetentionAttributes(lifecycleName string, strategy string, quantityToKeep string, unit string, shouldKeepForever string) string {
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

func lifecycle_withBasicPhase(lifecycleName string, phaseName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
  		phase {
    		name = "%s"
  		}
	}`, lifecycleName, lifecycleName, phaseName)
}

func lifecycleBasic(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
}`, lifecycleName, lifecycleName)
}

func lifecycleWithOnlyReleaseRetentionGivenAttributes(lifecycleName string, strategy string, quantityToKeep string, unit string, shouldKeepForever string) string {
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
