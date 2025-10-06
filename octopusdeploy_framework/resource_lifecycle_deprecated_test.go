package octopusdeploy_framework

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestExpandLifecycleWithNilUsingNewRetentionBlockDEPRECATED(t *testing.T) {
	lifecycle := expandLifecycleDEPRECATED(nil, false)
	require.Nil(t, lifecycle)
}

func TestExpanDLifecycleWithNilUsingRetentionWithoutStrategyBlockDEPRECATED(t *testing.T) {
	lifecycle := expandLifecycleDEPRECATED(nil, true)
	require.Nil(t, lifecycle)
}

func TestExpandLifecycleUsingNewRetentionBlockDEPRECATED(t *testing.T) {
	description := "test-description"
	name := "test-name"
	spaceID := "test-space-id"
	Id := "test-id"
	releaseRetention := core.KeepForeverRetentionPeriod()
	tentacleRetention := core.KeepForeverRetentionPeriod()
	retentionAttributeTypes := getResourceRetentionAttrTypesDEPRECATED()

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
						"strategy":         types.StringValue(releaseRetention.Strategy),
						"quantity_to_keep": types.Int64Value(int64(releaseRetention.QuantityToKeep)),
						"unit":             types.StringValue(releaseRetention.Unit),
					},
				),
			},
		),
	}
	data.ID = types.StringValue(Id)

	lifecycle := expandLifecycleDEPRECATED(data, false)
	require.NotNil(t, lifecycle)
	require.Equal(t, description, lifecycle.Description)
	require.NotEmpty(t, lifecycle.ID)
	require.NotNil(t, lifecycle.Links)
	require.Empty(t, lifecycle.Links)
	require.Equal(t, name, lifecycle.Name)
	require.Empty(t, lifecycle.Phases)
	require.Equal(t, releaseRetention, lifecycle.ReleaseRetentionPolicy)
	require.Equal(t, tentacleRetention, lifecycle.TentacleRetentionPolicy)
}

func TestExpandLifecycleUsingRetentionWithoutStrategyBlockDEPRECATED(t *testing.T) {
	description := "test-description"
	name := "test-name"
	spaceID := "test-space-id"
	Id := "test-id"
	releaseRetention := core.KeepForeverRetentionPeriod()
	tentacleRetention := core.KeepForeverRetentionPeriod()
	retentionAttributeTypes := getResourceRetentionAttrTypesDEPRECATED()

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
						"quantity_to_keep":    types.Int64Value(int64(releaseRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(releaseRetention.ShouldKeepForever),
						"unit":                types.StringValue(releaseRetention.Unit),
					},
				),
			},
		),
	}
	data.ID = types.StringValue(Id)

	lifecycle := expandLifecycleDEPRECATED(data, true)
	require.NotNil(t, lifecycle)
	require.Equal(t, description, lifecycle.Description)
	require.NotEmpty(t, lifecycle.ID)
	require.NotNil(t, lifecycle.Links)
	require.Empty(t, lifecycle.Links)
	require.Equal(t, name, lifecycle.Name)
	require.Empty(t, lifecycle.Phases)
	require.Equal(t, releaseRetention, lifecycle.ReleaseRetentionPolicy)
	require.Equal(t, tentacleRetention, lifecycle.TentacleRetentionPolicy)
}

func TestExpandAndFlattenPhasesWithSensibleDefaultsDEPRECATED(t *testing.T) {
	phase := createTestPhaseDEPRECATED("TestPhase", []string{"AutoTarget1", "AutoTarget2"}, true, 5)

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

func TestExpandAndFlattenMultiplePhasesWithSensibleDefaultsDEPRECATED(t *testing.T) {
	phase1 := createTestPhaseDEPRECATED("Phase1", []string{"AutoTarget1", "AutoTarget2"}, true, 5)
	phase2 := createTestPhaseDEPRECATED("Phase2", []string{"AutoTarget3", "AutoTarget4"}, false, 3)

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

func createTestPhaseDEPRECATED(name string, autoTargets []string, isOptional bool, minEnvs int32) *lifecycles.Phase {
	phase := lifecycles.NewPhase(name)
	phase.AutomaticDeploymentTargets = autoTargets
	phase.IsOptionalPhase = isOptional
	phase.MinimumEnvironmentsBeforePromotion = minEnvs
	phase.ReleaseRetentionPolicy = core.CountBasedRetentionPeriod(15, "Items")
	phase.TentacleRetentionPolicy = core.KeepForeverRetentionPeriod()
	phase.ID = name + "-Id"
	return phase
}
