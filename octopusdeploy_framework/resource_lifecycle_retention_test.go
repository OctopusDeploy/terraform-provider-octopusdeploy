package octopusdeploy_framework

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestExpandLifecycleWithDefaultRetentionStrategies(t *testing.T) {
	description := "test-description"
	name := "test-name"
	spaceID := "test-space-id"
	Id := "test-id"
	releaseRetention := core.SpaceDefaultRetentionPeriod()
	tentacleRetention := core.SpaceDefaultRetentionPeriod()

	data := &lifecycleTypeResourceModel{
		Description: types.StringValue(description),
		Name:        types.StringValue(name),
		SpaceID:     types.StringValue(spaceID),
		ReleaseRetentionPolicy: types.ListValueMust(
			types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
			[]attr.Value{
				types.ObjectValueMust(
					getRetentionPeriodAttrTypes(),
					map[string]attr.Value{
						"strategy":            types.StringValue(releaseRetention.Strategy),
						"quantity_to_keep":    types.Int64Value(int64(releaseRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(releaseRetention.ShouldKeepForever),
						"unit":                types.StringValue(releaseRetention.Unit),
					},
				),
			},
		),
		TentacleRetentionPolicy: types.ListValueMust(
			types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
			[]attr.Value{
				types.ObjectValueMust(
					getRetentionPeriodAttrTypes(),
					map[string]attr.Value{
						"strategy":            types.StringValue(tentacleRetention.Strategy),
						"quantity_to_keep":    types.Int64Value(int64(tentacleRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(tentacleRetention.ShouldKeepForever),
						"unit":                types.StringValue(tentacleRetention.Unit),
					},
				),
			},
		),
	}
	data.ID = types.StringValue(Id)

	lifecycle := expandLifecycle(data)

	require.Equal(t, description, lifecycle.Description)
	require.NotEmpty(t, lifecycle.ID)
	require.NotNil(t, lifecycle.Links)
	require.Empty(t, lifecycle.Links)
	require.Equal(t, name, lifecycle.Name)
	require.Empty(t, lifecycle.Phases)
	require.Equal(t, releaseRetention, lifecycle.ReleaseRetentionPolicy)
	require.Equal(t, tentacleRetention, lifecycle.TentacleRetentionPolicy)
	require.Equal(t, spaceID, lifecycle.SpaceID)
}

func TestDefaultRetentionValidation_withStrategyAttribute(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	shouldKeepForeverAttribute := "should_keep_forever = false"
	quantityToKeepAttribute := "quantity_to_keep = 0"
	unitAttribute := `unit = "Items"`

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      lifecycle_withDefaultRetention_usingStrategy_andUnwantedAttribute(lifecycleName, quantityToKeepAttribute),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep should not be supplied when strategy is set to Default.*quantity_to_keep should not be supplied when strategy is set to Default`),
			},
			{
				Config:      lifecycle_withDefaultRetention_usingStrategy_andUnwantedAttribute(lifecycleName, shouldKeepForeverAttribute),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)should_keep_forever should not be supplied when strategy is set to Default.*should_keep_forever should not be supplied when strategy is set to Default`),
			},
			{
				Config:   lifecycle_withDefaultRetention_usingStrategy_andUnwantedAttribute(lifecycleName, unitAttribute),
				PlanOnly: true,
				ExpectError: regexp.MustCompile(
					`(?s)unit should not be supplied when strategy is set to Default.*unit should not be supplied when strategy is set to Default`,
				),
			},
		},
	})
}

//func TestCountRetentionValidation_withStrategyAttribute(t *testing.T) {
//	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
//	shouldKeepForeverAttribute := "should_keep_forever = false"
//	quantityToKeepAttribute := "quantity_to_keep = 0"
//	unitAttribute := `unit = "Items"`
//
//	resource.Test(t, resource.TestCase{
//		CheckDestroy:             testAccLifecycleCheckDestroy,
//		PreCheck:                 func() { TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			//should error when should_keep_forever is also supplied
//			{
//				Config:      lifecycle_withDefaultRetention_usingStrategy_andUnwantedAttribute(lifecycleName, quantityToKeepAttribute),
//				PlanOnly:    true,
//				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep should not be supplied when strategy is set to Default.*quantity_to_keep should not be supplied when strategy is set to Default`),
//			},
//			{
//				Config:      lifecycle_withDefaultRetention_usingStrategy_andUnwantedAttribute(lifecycleName, shouldKeepForeverAttribute),
//				PlanOnly:    true,
//				ExpectError: regexp.MustCompile(`(?s)should_keep_forever should not be supplied when strategy is set to Default.*should_keep_forever should not be supplied when strategy is set to Default`),
//			},
//			{
//				Config:   lifecycle_withDefaultRetention_usingStrategy_andUnwantedAttribute(lifecycleName, unitAttribute),
//				PlanOnly: true,
//				ExpectError: regexp.MustCompile(
//					`(?s)unit should not be supplied when strategy is set to Default.*unit should not be supplied when strategy is set to Default`,
//				),
//			},
//		},
//	})
//}
//func TestForeverRetentionValidation_withStrategyAttribute(t *testing.T) {
//	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
//	shouldKeepForeverAttribute := "should_keep_forever = false"
//	quantityToKeepAttribute := "quantity_to_keep = 0"
//	unitAttribute := `unit = "Items"`
//
//	resource.Test(t, resource.TestCase{
//		CheckDestroy:             testAccLifecycleCheckDestroy,
//		PreCheck:                 func() { TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			//should error when should_keep_forever is also supplied
//			{
//				Config:      lifecycle_withDefaultRetention_usingStrategy_andUnwantedAttribute(lifecycleName, quantityToKeepAttribute),
//				PlanOnly:    true,
//				ExpectError: regexp.MustCompile(`(?s)quantity_to_keep should not be supplied when strategy is set to Default.*quantity_to_keep should not be supplied when strategy is set to Default`),
//			},
//			{
//				Config:      lifecycle_withDefaultRetention_usingStrategy_andUnwantedAttribute(lifecycleName, shouldKeepForeverAttribute),
//				PlanOnly:    true,
//				ExpectError: regexp.MustCompile(`(?s)should_keep_forever should not be supplied when strategy is set to Default.*should_keep_forever should not be supplied when strategy is set to Default`),
//			},
//			{
//				Config:   lifecycle_withDefaultRetention_usingStrategy_andUnwantedAttribute(lifecycleName, unitAttribute),
//				PlanOnly: true,
//				ExpectError: regexp.MustCompile(
//					`(?s)unit should not be supplied when strategy is set to Default.*unit should not be supplied when strategy is set to Default`,
//				),
//			},
//		},
//	})
//}

func TestLifecycleRetentionPolicyUpdates_withStrategyAttribute(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleResource := "octopusdeploy_lifecycle." + lifecycleName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1 create lifecycle without a retention policy
			{
				Config: lifecycleWithoutRetentionPolicies(lifecycleName),
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
				Config: lifecycle_withDefaultRetention_usingStrategy(lifecycleName),
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
				Config: lifecycle_withCountRetention_usingStrategy(lifecycleName),
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
				Config: lifecycle_withForeverRetention_usingStrategy(lifecycleName),
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
			{Config: lifecycleWithoutRetentionPolicies(lifecycleName),
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
func TestLifecycleRetentionPolicyUpdates_withoutStrategyAttribute(t *testing.T) {
	lifecycleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	lifecycleResource := "octopusdeploy_lifecycle." + lifecycleName

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 1 create lifecycle without a retention policy
			{
				Config: lifecycleWithoutRetentionPolicies(lifecycleName),
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
				Config: lifecycle_withDefaultRetention_withoutStrategyAttribute_usingQuantityToKeep(lifecycleName),
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
				Config: lifecycle_withCountRetention_withoutStrategyAttribute_usingDays(lifecycleName),
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
			// 3 update with Count retention policies
			{
				Config: lifecycle_withCountRetention_withoutStrategyAttribute_usingItems(lifecycleName),
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
				Config: lifecycle_withDefaultRetention_withoutStrategyAttribute_notUsingQuantityToKeep(lifecycleName),
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
			// 5 remove policies to return to default behaviour
			{Config: lifecycleWithoutRetentionPolicies(lifecycleName),
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

func TestLifecycle_WithBasicPhase(t *testing.T) {
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

				Config: testAccLifecycle_withBasicPhase(lifecycleName, phaseName),
			},
		},
	})
}

//do a few cases for old attribute type errors
//make sure an empty retention is recorded as "nill" in the state

// create lifecycles for testing
func lifecycleWithoutRetentionPolicies(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
	}`, lifecycleName, lifecycleName)
}

func lifecycle_withDefaultRetention_usingStrategy(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			strategy = "Default"
		}

		tentacle_retention_policy {
			strategy = "Default"
		}
    }`, lifecycleName, lifecycleName)
}

func lifecycle_withForeverRetention_usingStrategy(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			strategy = "Forever"
		}

		tentacle_retention_policy {
			strategy = "Forever"
		}
    }`, lifecycleName, lifecycleName)
}

func lifecycle_withCountRetention_usingStrategy(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			strategy = "Count"
			unit = "Days"
			quantity_to_keep = 10
		}
		tentacle_retention_policy {
			strategy = "Count"
			unit = "Items"
			quantity_to_keep = 10
		}
    }`, lifecycleName, lifecycleName)
}

func lifecycle_withDefaultRetention_usingStrategy_andUnwantedAttribute(lifecycleName string, unwantedAttribute string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
        description = ""
    	release_retention_policy {
    		strategy = "Default"      
    		%s
  		}
 		tentacle_retention_policy {
    		strategy = "Default"      
    		%s
  		}
	}`, lifecycleName, lifecycleName, unwantedAttribute, unwantedAttribute)
}

func lifecycle_withCountRetention_withoutStrategyAttribute_usingDays(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			should_keep_forever = "false"
			quantity_to_keep    = "1"
			unit                = "Days"
		}
		tentacle_retention_policy {
			quantity_to_keep = "1"
			unit             = "Days"
		}
    }`, lifecycleName, lifecycleName)
}

func lifecycle_withCountRetention_withoutStrategyAttribute_usingItems(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			should_keep_forever = "false"
			quantity_to_keep    = "1"
			unit                = "Items"
		}
		tentacle_retention_policy {
			quantity_to_keep = "1"
			unit             = "Items"
		}
    }`, lifecycleName, lifecycleName)
}

func lifecycle_withDefaultRetention_withoutStrategyAttribute_usingQuantityToKeep(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			quantity_to_keep    = "0"
			should_keep_forever = "true"
		}
		tentacle_retention_policy {
			quantity_to_keep = "0"
		}
    }`, lifecycleName, lifecycleName)
}

func lifecycle_withDefaultRetention_withoutStrategyAttribute_notUsingQuantityToKeep(lifecycleName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
		release_retention_policy {
			should_keep_forever = "true"
			unit             	= "Items"
		}
		tentacle_retention_policy {
			should_keep_forever = "true"
		}
    }`, lifecycleName, lifecycleName)
}

func testAccLifecycle_withBasicPhase(lifecycleName string, phaseName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
  		phase {
    		name = "%s"
  		}
	}`, lifecycleName, lifecycleName, phaseName)
}
