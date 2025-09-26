package octopusdeploy_framework

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetDeprecatedDefaultRetention(data *lifecycleTypeResourceModel, onlyDeprecatedRetentionIsUsed bool) (bool, bool) {

	hasUserDefinedReleaseRetention := attributeIsUsed(data.ReleaseRetention)
	hasUserDefinedTentacleRetention := attributeIsUsed(data.TentacleRetention)

	var initialRetentionSetting types.List
	if onlyDeprecatedRetentionIsUsed {
		initialRetentionSetting = DeprecatedFlattenRetention(core.NewRetentionPeriod(30, "Days", false))
	} else {
		initialRetentionSetting = types.ListNull(types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()})
	}

	if !hasUserDefinedReleaseRetention {
		data.ReleaseRetention = initialRetentionSetting
	}
	if !hasUserDefinedTentacleRetention {
		data.TentacleRetention = initialRetentionSetting
	}

	return hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention
}
func RemoveDeprecatedDefaultRetentionFromUnsetBlocks(data *lifecycleTypeResourceModel, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention bool) {
	// Remove retention policies from data before setting state, but only if they were not included in the plan
	if !hasUserDefinedReleaseRetention {
		data.ReleaseRetention = types.ListNull(types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()})
	}
	if !hasUserDefinedTentacleRetention {
		data.TentacleRetention = types.ListNull(types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()})
	}
}
func IsDeprecatedRetentionInPlan(data *lifecycleTypeResourceModel) bool {
	if attributeIsUsed(data.ReleaseRetention) || attributeIsUsed(data.TentacleRetention) {
		return true
	}
	for _, phase := range data.Phase.Elements() {
		phaseAttributes := phase.(types.Object).Attributes()
		releaseRetention := phaseAttributes["release_retention_policy"].(types.List)
		tentacleRetention := phaseAttributes["tentacle_retention_policy"].(types.List)
		if attributeIsUsed(releaseRetention) || attributeIsUsed(tentacleRetention) {
			return true
		}
	}
	return false
}
func IsOnlyDeprecatedRetentionUsed(data *lifecycleTypeResourceModel, onlyDeprecatedRetentionIsSupported bool) bool {
	deprecatedRetentionIsInPlan := IsDeprecatedRetentionInPlan(data)
	if deprecatedRetentionIsInPlan {
		return true
	}
	if onlyDeprecatedRetentionIsSupported {
		return true
	}
	return false
}
func DeprecatedFlattenPhases(goPhases []*lifecycles.Phase) types.List {
	var deprecatedAttributeTypes = DeprecatedGetAttributeTypes()
	if goPhases == nil {
		return types.ListNull(types.ObjectType{AttrTypes: deprecatedAttributeTypes})
	}
	phasesList := make([]attr.Value, 0, len(goPhases))

	for _, goPhase := range goPhases {
		attrs := map[string]attr.Value{
			"id":                                    types.StringValue(goPhase.ID),
			"name":                                  types.StringValue(goPhase.Name),
			"automatic_deployment_targets":          util.FlattenStringList(goPhase.AutomaticDeploymentTargets),
			"optional_deployment_targets":           util.FlattenStringList(goPhase.OptionalDeploymentTargets),
			"minimum_environments_before_promotion": types.Int64Value(int64(goPhase.MinimumEnvironmentsBeforePromotion)),
			"is_optional_phase":                     types.BoolValue(goPhase.IsOptionalPhase),
			"is_priority_phase":                     types.BoolValue(goPhase.IsPriorityPhase),
			"release_retention_policy":              util.Ternary(goPhase.ReleaseRetentionPolicy != nil, DeprecatedFlattenRetention(goPhase.ReleaseRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()})),
			"tentacle_retention_policy":             util.Ternary(goPhase.TentacleRetentionPolicy != nil, DeprecatedFlattenRetention(goPhase.TentacleRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()})),
			"release_retention_with_strategy":       types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}),
			"tentacle_retention_with_strategy":      types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}),
		}
		phasesList = append(phasesList, types.ObjectValueMust(deprecatedAttributeTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: deprecatedAttributeTypes}, phasesList)
}
func DeprecatedFlattenRetention(goRetention *core.RetentionPeriod) types.List {
	var deprecatedAttributeTypes = DeprecatedGetRetentionAttTypes()
	if goRetention == nil {
		return types.ListNull(types.ObjectType{AttrTypes: deprecatedAttributeTypes})
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: deprecatedAttributeTypes},
		[]attr.Value{
			types.ObjectValueMust(
				deprecatedAttributeTypes,
				map[string]attr.Value{
					"quantity_to_keep":    types.Int64Value(int64(goRetention.QuantityToKeep)),
					"should_keep_forever": types.BoolValue(goRetention.ShouldKeepForever),
					"unit":                types.StringValue(goRetention.Unit),
				},
			),
		},
	)
}
func DeprecatedExpandPhases(phasesUserInput types.List) []*lifecycles.Phase {
	if phasesUserInput.IsNull() || phasesUserInput.IsUnknown() || len(phasesUserInput.Elements()) == 0 {
		return nil
	}

	allPhasesPlan := make([]*lifecycles.Phase, 0, len(phasesUserInput.Elements()))

	for _, singlePhaseUserInput := range phasesUserInput.Elements() {
		phaseAttrsUserInput := singlePhaseUserInput.(types.Object).Attributes()

		phasePlan := &lifecycles.Phase{}

		if v, ok := phaseAttrsUserInput["id"].(types.String); ok && !v.IsNull() {
			phasePlan.ID = v.ValueString()
		}

		if v, ok := phaseAttrsUserInput["name"].(types.String); ok && !v.IsNull() {
			phasePlan.Name = v.ValueString()
		}

		if v, ok := phaseAttrsUserInput["automatic_deployment_targets"].(types.List); ok && !v.IsNull() {
			phasePlan.AutomaticDeploymentTargets = util.ExpandStringList(v)
		}

		if v, ok := phaseAttrsUserInput["optional_deployment_targets"].(types.List); ok && !v.IsNull() {
			phasePlan.OptionalDeploymentTargets = util.ExpandStringList(v)
		}

		if v, ok := phaseAttrsUserInput["minimum_environments_before_promotion"].(types.Int64); ok && !v.IsNull() {
			phasePlan.MinimumEnvironmentsBeforePromotion = int32(v.ValueInt64())
		}

		if v, ok := phaseAttrsUserInput["is_optional_phase"].(types.Bool); ok && !v.IsNull() {
			phasePlan.IsOptionalPhase = v.ValueBool()
		}

		if v, ok := phaseAttrsUserInput["is_priority_phase"].(types.Bool); ok && !v.IsNull() {
			phasePlan.IsPriorityPhase = v.ValueBool()
		}

		if v, ok := phaseAttrsUserInput["release_retention_policy"].(types.List); ok && !v.IsNull() {
			phasePlan.ReleaseRetentionPolicy = DeprecatedExpandRetention(v)
		}
		if v, ok := phaseAttrsUserInput["tentacle_retention_policy"].(types.List); ok && !v.IsNull() {
			phasePlan.TentacleRetentionPolicy = DeprecatedExpandRetention(v)
		}
		allPhasesPlan = append(allPhasesPlan, phasePlan)
	}

	return allPhasesPlan
}
func DeprecatedExpandRetention(v types.List) *core.RetentionPeriod {
	if v.IsNull() || v.IsUnknown() || len(v.Elements()) == 0 {
		return nil
	}

	obj := v.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	var quantityToKeep int32
	if qty, ok := attrs["quantity_to_keep"].(types.Int64); ok && !qty.IsNull() {
		quantityToKeep = int32(qty.ValueInt64())
	}

	var shouldKeepForever bool
	if keep, ok := attrs["should_keep_forever"].(types.Bool); ok && !keep.IsNull() {
		shouldKeepForever = keep.ValueBool()
	}

	var unit string
	if u, ok := attrs["unit"].(types.String); ok && !u.IsNull() {
		unit = u.ValueString()
	}

	return core.NewRetentionPeriod(quantityToKeep, unit, shouldKeepForever)
}
func DeprecatedGetRetentionAttTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}
func DeprecatedGetAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"is_priority_phase":                     types.BoolType,
		"release_retention_policy":              types.ListType{ElemType: types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()}},
		"tentacle_retention_policy":             types.ListType{ElemType: types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()}},
		"release_retention_with_strategy":       types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}},
		"tentacle_retention_with_strategy":      types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}},
	}
}
