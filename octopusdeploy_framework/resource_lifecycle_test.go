package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

func TestExpandLifecycleWithNil(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks {
		t.Skip("Skipping test because users may still use the deprecated retention blocks")
	}
	lifecycle := expandLifecycle(nil)
	require.Nil(t, lifecycle)
}

func TestExpandLifecycle(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks {
		t.Skip("Skipping test because users may still use the deprecated retention blocks")
	}
	description := "test-description"
	name := "test-name"
	spaceID := "test-space-id"
	Id := "test-id"
	releaseRetention := core.KeepForeverRetentionPeriod()
	tentacleRetention := core.KeepForeverRetentionPeriod()
	retentionAttributeTypes := getResourceRetentionAttrTypes()

	data := &lifecycleTypeResourceModel{
		Description: types.StringValue(description),
		Name:        types.StringValue(name),
		SpaceID:     types.StringValue(spaceID),
		ReleaseRetention: types.ListValueMust(
			types.ObjectType{AttrTypes: retentionAttributeTypes},
			[]attr.Value{
				types.ObjectValueMust(
					retentionAttributeTypes,
					map[string]attr.Value{
						"strategy":         types.StringValue(releaseRetention.Strategy),
						"quantity_to_keep": types.Int64Value(int64(releaseRetention.QuantityToKeep)),
						"unit":             types.StringValue(releaseRetention.Unit),
					},
				),
			},
		),
		TentacleRetention: types.ListValueMust(
			types.ObjectType{AttrTypes: retentionAttributeTypes},
			[]attr.Value{
				types.ObjectValueMust(
					retentionAttributeTypes,
					map[string]attr.Value{
						"strategy":         types.StringValue(releaseRetention.Strategy),
						"quantity_to_keep": types.Int64Value(int64(releaseRetention.QuantityToKeep)),
						"unit":             types.StringValue(releaseRetention.Unit),
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

func TestExpandPhasesWithEmptyInput(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks {
		t.Skip("Skipping test because users may still use the deprecated retention blocks")
	}
	emptyList := types.ListValueMust(types.ObjectType{AttrTypes: getResourcePhaseAttrTypes()}, []attr.Value{})
	phases := expandPhases(emptyList)
	require.Nil(t, phases)
}

func TestExpandPhasesWithNullInput(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks {
		t.Skip("Skipping test because users may still use the deprecated retention blocks")
	}
	nullList := types.ListNull(types.ObjectType{AttrTypes: getResourcePhaseAttrTypes()})
	phases := expandPhases(nullList)
	require.Nil(t, phases)
}

func TestExpandPhasesWithUnknownInput(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks {
		t.Skip("Skipping test because users may still use the deprecated retention blocks")
	}
	unknownList := types.ListUnknown(types.ObjectType{AttrTypes: getResourcePhaseAttrTypes()})
	phases := expandPhases(unknownList)
	require.Nil(t, phases)
}

func TestExpandAndFlattenPhasesWithSensibleDefaults(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks {
		t.Skip("Skipping test because users may still use the deprecated retention blocks")
	}
	phase := createTestPhase("TestPhase", []string{"AutoTarget1", "AutoTarget2"}, true, 5)

	flattenedPhases := flattenResourcePhases([]*lifecycles.Phase{phase})
	require.NotNil(t, flattenedPhases)
	require.Equal(t, 1, len(flattenedPhases.Elements()))

	expandedPhases := expandPhases(flattenedPhases)
	require.NotNil(t, expandedPhases)
	require.Len(t, expandedPhases, 1)

	expandedPhase := expandedPhases[0]
	require.NotEmpty(t, expandedPhase.ID)
	require.Equal(t, phase.AutomaticDeploymentTargets, expandedPhase.AutomaticDeploymentTargets)
	require.Equal(t, phase.IsOptionalPhase, expandedPhase.IsOptionalPhase)
	require.EqualValues(t, phase.MinimumEnvironmentsBeforePromotion, expandedPhase.MinimumEnvironmentsBeforePromotion)
	require.Equal(t, phase.Name, expandedPhase.Name)
	require.Equal(t, phase.ReleaseRetentionPolicy, expandedPhase.ReleaseRetentionPolicy)
	require.Equal(t, phase.TentacleRetentionPolicy, expandedPhase.TentacleRetentionPolicy)
}

func TestExpandAndFlattenMultiplePhasesWithSensibleDefaults(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks {
		t.Skip("Skipping test because users may still use the deprecated retention blocks")
	}
	phase1 := createTestPhase("Phase1", []string{"AutoTarget1", "AutoTarget2"}, true, 5)
	phase2 := createTestPhase("Phase2", []string{"AutoTarget3", "AutoTarget4"}, false, 3)

	flattenedPhases := flattenResourcePhases([]*lifecycles.Phase{phase1, phase2})
	require.NotNil(t, flattenedPhases)
	require.Equal(t, 2, len(flattenedPhases.Elements()))

	expandedPhases := expandPhases(flattenedPhases)
	require.NotNil(t, expandedPhases)
	require.Len(t, expandedPhases, 2)

	require.NotEmpty(t, expandedPhases[0].ID)
	require.Equal(t, phase1.AutomaticDeploymentTargets, expandedPhases[0].AutomaticDeploymentTargets)
	require.Equal(t, phase1.IsOptionalPhase, expandedPhases[0].IsOptionalPhase)
	require.EqualValues(t, phase1.MinimumEnvironmentsBeforePromotion, expandedPhases[0].MinimumEnvironmentsBeforePromotion)
	require.Equal(t, phase1.Name, expandedPhases[0].Name)
	require.Equal(t, phase1.ReleaseRetentionPolicy, expandedPhases[0].ReleaseRetentionPolicy)
	require.Equal(t, phase1.TentacleRetentionPolicy, expandedPhases[0].TentacleRetentionPolicy)

	require.NotEmpty(t, expandedPhases[1].ID)
	require.Equal(t, phase2.AutomaticDeploymentTargets, expandedPhases[1].AutomaticDeploymentTargets)
	require.Equal(t, phase2.IsOptionalPhase, expandedPhases[1].IsOptionalPhase)
	require.EqualValues(t, phase2.MinimumEnvironmentsBeforePromotion, expandedPhases[1].MinimumEnvironmentsBeforePromotion)
	require.Equal(t, phase2.Name, expandedPhases[1].Name)
	require.Equal(t, phase2.ReleaseRetentionPolicy, expandedPhases[1].ReleaseRetentionPolicy)
	require.Equal(t, phase2.TentacleRetentionPolicy, expandedPhases[1].TentacleRetentionPolicy)
}

func createTestPhase(name string, autoTargets []string, isOptional bool, minEnvs int32) *lifecycles.Phase {
	phase := lifecycles.NewPhase(name)
	phase.AutomaticDeploymentTargets = autoTargets
	phase.IsOptionalPhase = isOptional
	phase.MinimumEnvironmentsBeforePromotion = minEnvs
	phase.ReleaseRetentionPolicy = core.NewRetentionPeriod(15, "Items", false)
	phase.TentacleRetentionPolicy = core.NewRetentionPeriod(0, "Days", true)
	phase.ID = name + "-Id"
	return phase
}

//Integration test under here

func TestAccLifecycleBasic(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks {
		t.Skip("Skipping test because users may still use the deprecated retention blocks")
	}
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_lifecycle." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy .#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy .#", "0"),
				),
				Config: testAccLifecycle(localName, name),
			},
		},
	})
}

func TestAccLifecycleWithUpdate(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks {
		t.Skip("Skipping test because users may still use the deprecated retention blocks")
	}
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_lifecycle." + localName
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	phaseName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// create lifecycle with no description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "release_retention_policy"),
				),
				Config: testAccLifecycle(localName, name),
			},
			// update lifecycle with a description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "release_retention_policy"),
				),
				Config: testAccLifecycleWithDescription(localName, name, description),
			},
			// update lifecycle by removing its description
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "release_retention_policy"),
				),
				Config: testAccLifecycle(localName, name),
			},
			// update lifecycle by adding a phase
			{
				Config: testAccLifecycleWithPhase(localName, name, description, phaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "phase.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.name", phaseName),
					resource.TestCheckNoResourceAttr(resourceName, "release_retention_policy"),
				),
			},
			// update lifecycle by modifying its phase
			{
				Config: testAccLifecycleWithPhaseAndPhaseRetention(localName, name, description, phaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "phase.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.name", phaseName),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_with_strategy.0.strategy", "Default"),
					resource.TestCheckNoResourceAttr(resourceName, "release_retention_policy"),
				),
			},
		},
	})
}

func TestAccLifecycleComplex(t *testing.T) {
	if schemas.AllowDeprecatedAndNewRetentionBlocks {
		t.Skip("Skipping test because users may still use the deprecated retention blocks")
	}
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_lifecycle." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.0.quantity_to_keep", "2"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy.0.strategy", "Count"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy.0.unit", "Days"),
					resource.TestCheckNoResourceAttr(resourceName, "release_retention_policy"),
					testAccCheckLifecyclePhaseCount(name, 2),
				),
				Config: testAccLifecycleComplex(localName, name),
			},
		},
	})
}

func testAccLifecycleWithPhase(localName string, name string, description string, phaseName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
        description = "%s"
		phase {
			name = "%s"
		}
	}`, localName, name, description, phaseName)
}
func testAccLifecycleWithPhaseAndPhaseRetention(localName string, name string, description string, phaseName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
        description = "%s"
		phase {
			name = "%s"
			release_retention_with_strategy {
				strategy         = "Default"
			}
		}
	}`, localName, name, description, phaseName)
}

func testAccLifecycle(localName string, name string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
        description = ""
	}`, localName, name)
}

func testAccLifecycleWithDescription(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
       description = "%s"
    }`, localName, name, description)
}

func testAccLifecycleWithRetention(localName string, name string, description string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
       name        = "%s"
       description = "%s"
		release_retention_with_strategy {
			unit             = "Days"
			quantity_to_keep = 60
			should_keep_forever = false
		}

		tentacle_retention_with_strategy {
			unit             = "Items"
			quantity_to_keep = 0
			should_keep_forever = true
		}
    }`, localName, name, description)
}

func testAccLifecycleComplex(localName string, name string) string {
	environment1LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment1Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment2LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment2Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment3LocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environment3Name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAccEnvironment(environment1LocalName, environment1Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(environment2LocalName, environment2Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testAccEnvironment(environment3LocalName, environment3Name, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		`resource "octopusdeploy_lifecycle" "%s" {
			name        = "%s"
			description = "Funky Lifecycle description"

			release_retention_with_strategy {
				strategy         = "Count"
				unit             = "Days"
				quantity_to_keep = 2
			}

			tentacle_retention_with_strategy {
				strategy         = "Count"
				unit             = "Days"
				quantity_to_keep = 1
			}

			phase {
				automatic_deployment_targets          = ["${octopusdeploy_environment.%s.id}"]
				is_optional_phase                     = true
				minimum_environments_before_promotion = 2
				name                                  = "P1"
				optional_deployment_targets           = ["${octopusdeploy_environment.%s.id}"]
			}

			phase {
				name = "P2"
			}
	}`, localName, name, environment2LocalName, environment3LocalName)
}

func testAccCheckLifecycleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := existsHelperLifecycle(s, octoClient); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckLifecyclePhaseCount(name string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceList, err := octoClient.Lifecycles.GetByPartialName(name)
		if err != nil {
			return err
		}

		resource := resourceList[0]

		if len(resource.Phases) != expected {
			return fmt.Errorf("lifecycle has %d phases instead of the expected %d", len(resource.Phases), expected)
		}

		return nil
	}
}

func existsHelperLifecycle(s *terraform.State, client *client.Client) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "octopusdeploy_lifecycle" {
			if _, err := client.Lifecycles.GetByID(r.Primary.ID); err != nil {
				return fmt.Errorf("error retrieving lifecycle %s", err)
			}
		}
	}
	return nil
}

func testAccLifecycleCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_lifecycle" {
			continue
		}

		lifecycle, err := octoClient.Lifecycles.GetByID(rs.Primary.ID)
		if err == nil && lifecycle != nil {
			return fmt.Errorf("lifecycle (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}
