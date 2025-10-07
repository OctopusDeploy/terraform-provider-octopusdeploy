package octopusdeploy_framework

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccLifecycleRetentionUpdates_DEPRECATED(t *testing.T) {
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
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0"),
				),
			},
			// 2 update with the default (forever) retention policies
			{
				Config: defaultRetentionLifecycle_usingQuantityToKeepDEPRECATED(lifecycleName),
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
				),
			},
			// 3 update with Count retention policies using days
			{
				Config: countRetentionLifecycle_DEPRECATED(lifecycleName, "Days"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Days"),
				),
			},
			// 4 update with Count retention policies using items
			{
				Config: countRetentionLifecycle_DEPRECATED(lifecycleName, "items"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "items"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "items"),
				),
			},
			// 5 update with Default retention policies
			{
				Config: defaultRetentionLifecycle_notUsingQuantityToKeepDEPRECATED(lifecycleName),
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
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.0.unit", "Items")),
			},
			// 6 set only release retention policy
			{
				Config: lifecycle_ReleaseRetentionWithoutStrategy_DEPRECATED(lifecycleName, "1", "Days", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(lifecycleResource, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(lifecycleResource, "space_id"),
					resource.TestCheckResourceAttr(lifecycleResource, "tentacle_retention_policy.#", "0")),
			},
		},
	})
}

func TestAccRetentionAttributeValidation_DEPRECATED(t *testing.T) {

	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
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
				Config: lifecycle_withBasicPhase(lifecycleName, phaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(lifecycleResource),
					resource.TestCheckResourceAttrSet(lifecycleResource, "id"),
					resource.TestCheckResourceAttr(lifecycleResource, "name", lifecycleName),
					resource.TestCheckResourceAttr(lifecycleResource, "phase.#", "1"),
					//check that the phase retention policies remain empty so will inherit their policies from elsewhere
					resource.TestCheckResourceAttr(lifecycleResource, "phase.0.release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(lifecycleResource, "phase.0.tentacle_retention_policy.#", "0"),
				),
			},
		},
	})
}

func countRetentionLifecycle_DEPRECATED(lifecycleName string, unit string) string {
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

func defaultRetentionLifecycle_usingQuantityToKeepDEPRECATED(lifecycleName string) string {
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

func defaultRetentionLifecycle_notUsingQuantityToKeepDEPRECATED(lifecycleName string) string {
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

func lifecycleBasic(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
}`, lifecycleName, lifecycleName)
}

func lifecycle_withBasicPhase(lifecycleName string, phaseName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
		release_retention_policy {
			should_keep_forever = "true"
		}
		tentacle_retention_policy {
			should_keep_forever = "true"
			unit = "Items"
		}
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
