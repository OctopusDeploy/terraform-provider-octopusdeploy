package octopusdeploy_framework

import (
	"fmt"
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

// unit tests
func TestExpandLifecycleWithNil_usingNewRetentionBlockDEPRECATED(t *testing.T) {
	lifecycle := expandLifecycleDEPRECATED(nil, false)
	require.Nil(t, lifecycle)
}
func TestExpanDLifecycleWithNil_usingRetentionWithoutStrategyBlockDEPRECATED(t *testing.T) {
	lifecycle := expandLifecycleDEPRECATED(nil, true)
	require.Nil(t, lifecycle)
}

func TestExpandLifecycle_usingNewRetentionBlockDEPRECATED(t *testing.T) {
	description := "test-description"
	name := "test-name"
	spaceID := "test-space-id"
	Id := "test-id"
	releaseRetention := core.KeepForeverRetentionPeriod()
	tentacleRetention := core.CountBasedRetentionPeriod(2, "Items")
	retentionAttributeTypes := getResourceRetentionAttrTypes()

	data := &lifecycleTypeResourceModelDEPRECATED{
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
						"strategy":         types.StringValue(tentacleRetention.Strategy),
						"quantity_to_keep": types.Int64Value(int64(tentacleRetention.QuantityToKeep)),
						"unit":             types.StringValue(tentacleRetention.Unit),
					},
				),
			},
		),
	}
	data.ID = types.StringValue(Id)

	lifecycle := expandLifecycleDEPRECATED(data, false)
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
func TestExpandLifecycle_usingRetentionWithoutStrategyBlockDEPRECATED(t *testing.T) {
	description := "test-description"
	name := "test-name"
	spaceID := "test-space-id"
	Id := "test-id"
	releaseRetention := core.KeepForeverRetentionPeriod()
	tentacleRetention := core.CountBasedRetentionPeriod(2, "Items")
	retentionAttributeTypes := getResourceRetentionWithoutStrategyAttrTypesDEPRECATED()

	data := &lifecycleTypeResourceModelDEPRECATED{
		Description: types.StringValue(description),
		Name:        types.StringValue(name),
		SpaceID:     types.StringValue(spaceID),
		ReleaseRetentionWithoutStrategy: types.ListValueMust(
			types.ObjectType{AttrTypes: retentionAttributeTypes},
			[]attr.Value{
				types.ObjectValueMust(
					retentionAttributeTypes,
					map[string]attr.Value{
						"quantity_to_keep":    types.Int64Value(int64(releaseRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(releaseRetention.ShouldKeepForever),
						"unit":                types.StringValue(releaseRetention.Unit),
					},
				),
			},
		),
		TentacleRetentionWithoutStrategy: types.ListValueMust(
			types.ObjectType{AttrTypes: retentionAttributeTypes},
			[]attr.Value{
				types.ObjectValueMust(
					retentionAttributeTypes,
					map[string]attr.Value{
						"quantity_to_keep":    types.Int64Value(int64(tentacleRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(tentacleRetention.ShouldKeepForever),
						"unit":                types.StringValue(tentacleRetention.Unit),
					},
				),
			},
		),
	}
	data.ID = types.StringValue(Id)

	lifecycle := expandLifecycleDEPRECATED(data, true)
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

func TestExpandPhasesWithEmptyInput_DEPRECATED(t *testing.T) {
	emptyList := types.ListValueMust(types.ObjectType{AttrTypes: getResourcePhaseAttrTypesDEPRECATED()}, []attr.Value{})
	phases := expandPhasesDEPRECATED(emptyList)
	require.Nil(t, phases)
}
func TestExpandPhasesWithNullInput_DEPRECATED(t *testing.T) {
	nullList := types.ListNull(types.ObjectType{AttrTypes: getResourcePhaseAttrTypesDEPRECATED()})
	phases := expandPhasesDEPRECATED(nullList)
	require.Nil(t, phases)
}
func TestExpandPhasesWithUnknownInput_DEPRECATED(t *testing.T) {
	unknownList := types.ListUnknown(types.ObjectType{AttrTypes: getResourcePhaseAttrTypesDEPRECATED()})
	phases := expandPhasesDEPRECATED(unknownList)
	require.Nil(t, phases)
}

func TestExpandAndFlattenPhasesWithSensibleDefaults_UsingNewRetentionBlockDEPRECATED(t *testing.T) {
	phase := createTestPhase("TestPhase", []string{"AutoTarget1", "AutoTarget2"}, true, 5)

	flattenedPhases := flattenResourcePhasesDEPRECATED([]*lifecycles.Phase{phase}, false)
	require.NotNil(t, flattenedPhases)
	require.Equal(t, 1, len(flattenedPhases.Elements()))

	expandedPhases := expandPhasesDEPRECATED(flattenedPhases)
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

func TestExpandAndFlattenPhasesWithSensibleDefaults_UsingRetentionWithoutStrategyBlockDEPRECATED(t *testing.T) {
	phase := createTestPhase("TestPhase", []string{"AutoTarget1", "AutoTarget2"}, true, 5)

	flattenedPhases := flattenResourcePhasesDEPRECATED([]*lifecycles.Phase{phase}, true)
	require.NotNil(t, flattenedPhases)
	require.Equal(t, 1, len(flattenedPhases.Elements()))

	expandedPhases := expandPhasesDEPRECATED(flattenedPhases)
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

func TestExpandAndFlattenMultiplePhasesWithSensibleDefaults_UsingNewRetentionBlockDEPRECATED(t *testing.T) {
	phase1 := createTestPhase("Phase1", []string{"AutoTarget1", "AutoTarget2"}, true, 5)
	phase2 := createTestPhase("Phase2", []string{"AutoTarget3", "AutoTarget4"}, false, 3)

	flattenedPhases := flattenResourcePhasesDEPRECATED([]*lifecycles.Phase{phase1, phase2}, false)
	require.NotNil(t, flattenedPhases)
	require.Equal(t, 2, len(flattenedPhases.Elements()))

	expandedPhases := expandPhasesDEPRECATED(flattenedPhases)
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

func TestExpandAndFlattenMultiplePhasesWithSensibleDefaults_UsingRetentionWithoutStrategyBlockDEPRECATED(t *testing.T) {
	phase1 := createTestPhase("Phase1", []string{"AutoTarget1", "AutoTarget2"}, true, 5)
	phase2 := createTestPhase("Phase2", []string{"AutoTarget3", "AutoTarget4"}, false, 3)

	flattenedPhases := flattenResourcePhasesDEPRECATED([]*lifecycles.Phase{phase1, phase2}, true)
	require.NotNil(t, flattenedPhases)
	require.Equal(t, 2, len(flattenedPhases.Elements()))

	expandedPhases := expandPhasesDEPRECATED(flattenedPhases)
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

//integration tests

func TestAccLifecycleBasicDEPRECATED(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy .#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "0"),
				),
				Config: testAccLifecycle(localName, name),
			},
		},
	})
}

func TestAccLifecycleWithUpdateDEPRECATED(t *testing.T) {
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
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
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
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy.#", "0"),
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
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy.#", "0"),
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
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "phase.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.name", phaseName),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.tentacle_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.tentacle_retention_with_strategy.#", "0"),
				),
			},
			// update lifecycle by modifying its phase retention
			{
				Config: testAccLifecycleWithPhase_AndPhaseRetention(localName, name, description, phaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "phase.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.name", phaseName),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_with_strategy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_with_strategy.0.strategy", "Default"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_policy.#", "0"),
				),
			},
			// update lifecycle by switching its phase to use retention without strategy
			{
				Config: testAccLifecycleWithPhase_AndPhaseRetentionWithoutStrategyDEPRECATED(localName, name, description, phaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "phase.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.name", phaseName),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_policy.0.should_keep_forever", "true"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "phase.0.tentacle_retention_policy.#", "0"),
				),
			},
		},
	})
}

func TestAccLifecycleComplex_usingNewRetentionDEPRECATED(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "0"),
					testAccCheckLifecyclePhaseCount(name, 2),
				),
				Config: testAccLifecycleComplex(localName, name),
			},
		},
	})
}

func TestAccLifecycleComplex_usingRetentionWithoutStrategyDEPRECATED(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_lifecycle." + localName

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testAccLifecycleCheckDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLifecycleComplexDEPRECATED(localName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.quantity_to_keep", "2"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttrSet(resourceName, "space_id"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.should_keep_forever", "false"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.quantity_to_keep", "1"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_policy.0.unit", "Days"),
					resource.TestCheckResourceAttr(resourceName, "release_retention_with_strategy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tentacle_retention_with_strategy.#", "0"),
					testAccCheckLifecyclePhaseCount(name, 2),
				),
			},
		},
	})
}

// Setup for testing
func createTestPhase(name string, autoTargets []string, isOptional bool, minEnvs int32) *lifecycles.Phase {
	phase := lifecycles.NewPhase(name)
	phase.AutomaticDeploymentTargets = autoTargets
	phase.IsOptionalPhase = isOptional
	phase.MinimumEnvironmentsBeforePromotion = minEnvs
	phase.ReleaseRetentionPolicy = core.CountBasedRetentionPeriod(15, "Items")
	phase.TentacleRetentionPolicy = core.KeepForeverRetentionPeriod()
	phase.ID = name + "-Id"
	return phase
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

func testAccLifecycleWithPhase_AndPhaseRetention(localName string, name string, description string, phaseName string) string {
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

func testAccLifecycleWithPhase_AndPhaseRetentionWithoutStrategyDEPRECATED(localName string, name string, description string, phaseName string) string {
	return fmt.Sprintf(`resource "octopusdeploy_lifecycle" "%s" {
		name = "%s"
        description = "%s"
		phase {
			name = "%s"
			release_retention_policy {
				should_keep_forever = true
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

func testAccLifecycleComplexDEPRECATED(localName string, name string) string {
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

			release_retention_policy {
				quantity_to_keep = 2
				unit             = "Days"
				should_keep_forever = false
			}

			tentacle_retention_policy {
				quantity_to_keep = 1
				unit             = "Days"
				should_keep_forever = false
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
