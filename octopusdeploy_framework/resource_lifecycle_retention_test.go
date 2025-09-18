package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccBestPracticeLifecycleRetentionValidation(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	shouldKeepForeverAttribute := "should_keep_forever = false"
	quantityToKeepAttribute := "quantity_to_keep = 0"
	unitAttribute := `unit = "Items"`

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Default strategy validation
			{
				Config:      lifecycleRetentionPolicyByStrategy_withUnwantedAttribute(lifecycleName, "Default", quantityToKeepAttribute),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep should not be supplied when strategy is set to Default.*quantity_to_keep should not be supplied when strategy is set to Default`),
			},
			{
				Config:      lifecycleRetentionPolicyByStrategy_withUnwantedAttribute(lifecycleName, "Default", shouldKeepForeverAttribute),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)should_keep_forever should not be supplied when strategy is set to Default.*should_keep_forever should not be supplied when strategy is set to Default`),
			},
			{
				Config:   lifecycleRetentionPolicyByStrategy_withUnwantedAttribute(lifecycleName, "Default", unitAttribute),
				PlanOnly: true,
				ExpectError: regexp.MustCompile(
					`(?s)unit should not be supplied when strategy is set to Default.*unit should not be supplied when strategy is set to Default`,
				),
			},
			//Forever Strategy Validation
			{
				Config:      lifecycleRetentionPolicyByStrategy_withUnwantedAttribute(lifecycleName, "Forever", quantityToKeepAttribute),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep should not be supplied when strategy is set to Forever.*quantity_to_keep should not be supplied when strategy is set to Forever`),
			},
			{
				Config:      lifecycleRetentionPolicyByStrategy_withUnwantedAttribute(lifecycleName, "Forever", shouldKeepForeverAttribute),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)should_keep_forever should not be supplied when strategy is set to Forever.*should_keep_forever should not be supplied when strategy is set to Forever`),
			},
			{
				Config:   lifecycleRetentionPolicyByStrategy_withUnwantedAttribute(lifecycleName, "Forever", unitAttribute),
				PlanOnly: true,
				ExpectError: regexp.MustCompile(
					`(?s)unit should not be supplied when strategy is set to Forever.*unit should not be supplied when strategy is set to Forever`,
				),
			},
			//Count Strategy Validation
			{
				Config:      lifecycleRetentionPolicyByStrategy_withUnwantedAttribute(lifecycleName, "Count", shouldKeepForeverAttribute),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)should_keep_forever should not be supplied when strategy is set to Count.*should_keep_forever should not be supplied when strategy is set to Count`),
			},
		},
	})
}

func TestAccBestPracticeLifecycleRetentionUpdates(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleResource := "octopusdeploy_lifecycle." + lifecycleName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1 create lifecycle without a retention policy
			{
				Config: lifecycleRetentionPolicyByStrategy(lifecycleName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),
				),
			},
			// 2 update with default retention policies
			{
				Config: lifecycleRetentionPolicyByStrategy(lifecycleName, "Default"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),
				),
			},
			// 3 update with Count retention policies
			{
				Config: lifecycleRetentionPolicyByStrategy(lifecycleName, "Count"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "10"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "10"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),
				),
			},
			// 4 update with forever retention policies
			{
				Config: lifecycleRetentionPolicyByStrategy(lifecycleName, "Forever"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.strategy", "Forever"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.strategy", "Forever"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),
				),
			},
			// 5 remove policies to return to default behaviour
			{Config: lifecycleRetentionPolicyByStrategy(lifecycleName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),
				),
			},
		},
	})
}

func TestAccOldPracticeLifecycleRetentionPolicyUpdates(t *testing.T) {
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
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),
				),
			},
			// 2 update with default retention policies
			{
				Config: defaultRetentionLifecycle_usingQuantityToKeep(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),
				),
			},
			// 3 update with Count retention policies using days
			{
				Config: countRetentionLifecycle(lifecycleName, "Days"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Days"),
				),
			},
			// 3 update with Count retention policies using items
			{
				Config: countRetentionLifecycle(lifecycleName, "Items"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),
				),
			},
			// 4 update with Default retention policies
			{
				Config: defaultRetentionLifecycle_notUsingQuantityToKeep(lifecycleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Items"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items"),
				),
			},
		},
	})
}

func TestAccOldPracticeRetentionAttributeValidation(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	commonError := regexp.MustCompile(`Incorrect use of retention attributes\.\nFor best practice use strategy attribute\.`)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// when quantity_to_keep is > 0 should_keep_forever shouldn't be true
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "1", "Items", "true"),
				PlanOnly:    true,
				ExpectError: commonError,
			},
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "1", "", "true"),
				PlanOnly:    true,
				ExpectError: commonError,
			},
			// when quantity_to_keep is 0, should_keep_forever shouldn't be false
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "0", "", "false"),
				PlanOnly:    true,
				ExpectError: commonError,
			},
			//Error message intentionally vague as it occurs at api.
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "", "", "false"),
				PlanOnly:    true,
				ExpectError: commonError,
			},
			//Error message intentionally vague as it occurs at api.
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "", "Items", "false"),
				PlanOnly:    true,
				ExpectError: commonError,
			},
			// when there is an empty block. //Error message intentionally vague as it occurs at api.
			{
				Config:      lifecycleGivenRetentionAttributes(lifecycleName, "", "", ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("please either add retention policy attributes or remove the entire block"),
			},
		},
	})
}

func TestAccLifecyclePhaseRetentionPolicyInheritance(t *testing.T) {
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

//best practice lifecycles

func lifecycleRetentionPolicyByStrategy(lifecycleName string, strategy string) string {
	if strategy == "Default" || strategy == "Forever" {
		return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			strategy = "%s"
		}
		tentacle_retention_policy {
			strategy = "%s"
		}
    }`, lifecycleName, lifecycleName, strategy, strategy)
	}
	if strategy == "Count" {
		return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			strategy = "%s"
			unit = "Days"
			quantity_to_keep = 10
		}
		tentacle_retention_policy {
			strategy = "%s"
			unit = "Items"
			quantity_to_keep = 10
		}
    }`, lifecycleName, lifecycleName, strategy, strategy)
	}
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
	}`, lifecycleName, lifecycleName)

}

func lifecycleRetentionPolicyByStrategy_withUnwantedAttribute(lifecycleName string, strategy string, unwantedAttribute string) string {
	var countAttributes string
	if strategy == "Count" {
		countAttributes = `quantity_to_keep = "1"
							unit = "Items"`
	}

	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
    	release_retention_policy {
    		strategy = "%s"      
    		%s
			%s
  		}
 		tentacle_retention_policy {
    		strategy = "%s"      
    		%s
			%s
  		}
	}`, lifecycleName, lifecycleName, strategy, unwantedAttribute, countAttributes, strategy, unwantedAttribute, countAttributes)
}

func lifecycle_withBasicPhase(lifecycleName string, phaseName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
  		phase {
    		name = "%s"
  		}
	}`, lifecycleName, lifecycleName, phaseName)
}

// old use of attributes
func countRetentionLifecycle(lifecycleName string, unit string) string {
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

func defaultRetentionLifecycle_usingQuantityToKeep(lifecycleName string) string {
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

func defaultRetentionLifecycle_notUsingQuantityToKeep(lifecycleName string) string {
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

func lifecycleGivenRetentionAttributes(lifecycleName string, quantityToKeep string, unit string, shouldKeepForever string) string {
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
	}`, lifecycleName, lifecycleName, quantityToKeepAttribute, shouldKeepForeverAttribute, unitAttribute, quantityToKeepAttribute, shouldKeepForeverAttribute, unitAttribute)

}

func lifecycleBasic(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
}`, lifecycleName, lifecycleName)
}
