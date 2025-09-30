package octopusdeploy_framework

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetDeprecatedDefaultRetention(data *lifecycleTypeResourceModel, initialDeprecatedRetentionSetting types.List) (bool, bool) {
	hasUserDefinedReleaseRetention := attributeIsUsed(data.DeprecatedReleaseRetention)
	hasUserDefinedTentacleRetention := attributeIsUsed(data.DeprecatedTentacleRetention)
	if !hasUserDefinedReleaseRetention {
		data.DeprecatedReleaseRetention = initialDeprecatedRetentionSetting
	}
	if !hasUserDefinedTentacleRetention {
		data.DeprecatedTentacleRetention = initialDeprecatedRetentionSetting
	}

	return hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention
}
func RemoveDeprecatedDefaultRetentionFromUnsetBlocks(data *lifecycleTypeResourceModel, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention bool) {
	// Remove retention policies from data before setting state, but only if we added the initial value to them in the first place
	if !hasUserDefinedReleaseRetention {
		data.DeprecatedReleaseRetention = ListNullDeprecatedRetention
	}
	if !hasUserDefinedTentacleRetention {
		data.DeprecatedTentacleRetention = ListNullDeprecatedRetention
	}
}

var ListNullDeprecatedRetention = types.ListNull(types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()})

var ListNullRetentionWithStrategy = types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})

func IsDeprecatedRetentionInPlan(data *lifecycleTypeResourceModel) bool {
	if attributeIsUsed(data.DeprecatedReleaseRetention) || attributeIsUsed(data.DeprecatedTentacleRetention) {
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
func IsRetentionWithStrategyInPlan(data *lifecycleTypeResourceModel) bool {
	if attributeIsUsed(data.ReleaseRetentionWithStrategy) || attributeIsUsed(data.TentacleRetentionWithStrategy) {
		return true
	}
	for _, phase := range data.Phase.Elements() {
		phaseAttributes := phase.(types.Object).Attributes()
		releaseRetentionWithStrategy := phaseAttributes["release_retention_with_strategy"].(types.List)
		tentacleRetentionWithStrategy := phaseAttributes["tentacle_retention_with_strategy"].(types.List)

		if attributeIsUsed(releaseRetentionWithStrategy) || attributeIsUsed(tentacleRetentionWithStrategy) {
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
	var deprecatedAttributeTypes = GetPhaseAttributeTypes()
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
			"release_retention_policy":              util.Ternary(goPhase.ReleaseRetentionPolicy != nil, FlattenDeprecatedRetention(goPhase.ReleaseRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()})),
			"tentacle_retention_policy":             util.Ternary(goPhase.TentacleRetentionPolicy != nil, FlattenDeprecatedRetention(goPhase.TentacleRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()})),
			"release_retention_with_strategy":       types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}),
			"tentacle_retention_with_strategy":      types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}),
		}
		phasesList = append(phasesList, types.ObjectValueMust(deprecatedAttributeTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: deprecatedAttributeTypes}, phasesList)
}
func FlattenDeprecatedRetention(goRetention *core.RetentionPeriod) types.List {
	var deprecatedAttributeTypes = DeprecatedGetRetentionAttTypes()
	if goRetention == nil {
		return ListNullDeprecatedRetention
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
func ValidateRetentionBlocksUsed(data *lifecycleTypeResourceModel, diag *diag.Diagnostics, onlyDeprecatedRetentionIsSupported bool) {
	retentionWithStrategyIsInPlan := IsRetentionWithStrategyInPlan(data)
	deprecatedRetentionIsInPlan := IsDeprecatedRetentionInPlan(data)
	if retentionWithStrategyIsInPlan && deprecatedRetentionIsInPlan {
		diag.AddError("Retention blocks conflict", "Both release_retention_with_strategy and release_retention_policy are used. Only one can be used at a time.")
	}
	if retentionWithStrategyIsInPlan && onlyDeprecatedRetentionIsSupported {
		diag.AddError("retention with strategy is not supported on this Octopus Server version. Please upgrade to Octopus Server 2025.3 or later.", "")
	}
}
func DeprecatedRetentionObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}
