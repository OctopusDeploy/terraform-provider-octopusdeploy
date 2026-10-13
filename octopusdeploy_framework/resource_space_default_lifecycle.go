package octopusdeploy_framework

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// The default lifecycle of a space cannot be created or deleted through the API; Octopus creates it with the space
// and deletes it with the space. These helpers locate it and apply the configured settings on top of it.

var spaceDefaultLifecycleObjectType = types.ObjectType{AttrTypes: schemas.SpaceDefaultLifecycleAttrTypes()}
var spaceDefaultLifecyclePhaseObjectType = types.ObjectType{AttrTypes: schemas.SpaceDefaultLifecyclePhaseAttrTypes()}
var spaceDefaultLifecycleRetentionObjectType = types.ObjectType{AttrTypes: schemas.SpaceDefaultLifecycleRetentionAttrTypes()}

// applySpaceDefaultLifecycle updates the default lifecycle of the space to match the plan and returns the resulting state value.
// lifecycleID may be empty, in which case the default lifecycle is looked up in the space.
func (s *spaceResource) applySpaceDefaultLifecycle(ctx context.Context, spaceID string, lifecycleID string, plan types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	if plan.IsNull() {
		return types.ObjectNull(spaceDefaultLifecycleObjectType.AttrTypes), diags
	}

	if !s.Config.IsVersionSameOrGreaterThan("2025.3") {
		diags.AddError("Octopus Server Upgrade Required.", "default_lifecycle is not supported on this Octopus Server version. Please upgrade to Octopus Server 2025.3 or later.")
		return plan, diags
	}

	var model schemas.SpaceDefaultLifecycleModel
	diags.Append(plan.As(ctx, &model, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return plan, diags
	}

	var existing *lifecycles.Lifecycle
	var err error
	if lifecycleID != "" {
		existing, err = lifecycles.GetByID(s.Client, spaceID, lifecycleID)
		if isNotFound(err) {
			existing, err = s.findSpaceDefaultLifecycle(spaceID)
		}
	} else {
		existing, err = s.findSpaceDefaultLifecycle(spaceID)
	}
	if err != nil {
		diags.AddError("unable to load default lifecycle of space", err.Error())
		return plan, diags
	}

	toUpdate, d := expandSpaceDefaultLifecycle(ctx, existing, model)
	diags.Append(d...)
	if diags.HasError() {
		return plan, diags
	}

	tflog.Debug(ctx, fmt.Sprintf("updating default lifecycle '%s' of space '%s'", toUpdate.ID, spaceID))
	updated, err := lifecycles.Update(s.Client, toUpdate)
	if err != nil {
		diags.AddError("unable to update default lifecycle of space", err.Error())
		return plan, diags
	}

	return flattenSpaceDefaultLifecycle(ctx, updated, model)
}

// readSpaceDefaultLifecycle refreshes the default lifecycle state. It returns a null value when the lifecycle is not managed or no longer exists.
func (s *spaceResource) readSpaceDefaultLifecycle(ctx context.Context, spaceID string, state types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	if state.IsNull() || state.IsUnknown() {
		return types.ObjectNull(spaceDefaultLifecycleObjectType.AttrTypes), diags
	}

	var model schemas.SpaceDefaultLifecycleModel
	diags.Append(state.As(ctx, &model, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return state, diags
	}

	lifecycle, err := lifecycles.GetByID(s.Client, spaceID, model.ID.ValueString())
	if err != nil {
		if isNotFound(err) {
			tflog.Info(ctx, fmt.Sprintf("default lifecycle (%s) of space (%s) not found", model.ID.ValueString(), spaceID))
			return types.ObjectNull(spaceDefaultLifecycleObjectType.AttrTypes), diags
		}
		diags.AddError("unable to load default lifecycle of space", err.Error())
		return state, diags
	}

	return flattenSpaceDefaultLifecycle(ctx, lifecycle, model)
}

func (s *spaceResource) findSpaceDefaultLifecycle(spaceID string) (*lifecycles.Lifecycle, error) {
	all, err := lifecycles.GetAll(s.Client, spaceID)
	if err != nil {
		return nil, err
	}
	for _, l := range all {
		if l.Name == schemas.SpaceDefaultLifecycleName {
			return l, nil
		}
	}
	if len(all) == 1 {
		return all[0], nil
	}
	return nil, fmt.Errorf("could not find a lifecycle named '%s' in space '%s'", schemas.SpaceDefaultLifecycleName, spaceID)
}

func isNotFound(err error) bool {
	if apiError, ok := err.(*core.APIError); ok {
		return apiError.StatusCode == http.StatusNotFound
	}
	return false
}

func expandSpaceDefaultLifecycle(ctx context.Context, existing *lifecycles.Lifecycle, model schemas.SpaceDefaultLifecycleModel) (*lifecycles.Lifecycle, diag.Diagnostics) {
	var diags diag.Diagnostics

	if isKnown(model.Name) {
		existing.Name = model.Name.ValueString()
	}
	if isKnown(model.Description) {
		existing.Description = model.Description.ValueString()
	}

	existing.ReleaseRetentionPolicy = expandSpaceDefaultLifecycleRetention(ctx, model.ReleaseRetention, &diags)
	if existing.ReleaseRetentionPolicy == nil {
		existing.ReleaseRetentionPolicy = core.SpaceDefaultRetentionPeriod()
	}
	existing.TentacleRetentionPolicy = expandSpaceDefaultLifecycleRetention(ctx, model.TentacleRetention, &diags)
	if existing.TentacleRetentionPolicy == nil {
		existing.TentacleRetentionPolicy = core.SpaceDefaultRetentionPeriod()
	}

	existing.Phases = []*lifecycles.Phase{}
	if !model.Phases.IsNull() && !model.Phases.IsUnknown() {
		var phaseModels []schemas.SpaceDefaultLifecyclePhaseModel
		diags.Append(model.Phases.ElementsAs(ctx, &phaseModels, false)...)
		for _, p := range phaseModels {
			phase := lifecycles.NewPhase(p.Name.ValueString())
			if isKnown(p.ID) {
				phase.ID = p.ID.ValueString()
			}
			phase.AutomaticDeploymentTargets = util.ExpandStringList(p.AutomaticDeploymentTargets)
			phase.OptionalDeploymentTargets = util.ExpandStringList(p.OptionalDeploymentTargets)
			phase.MinimumEnvironmentsBeforePromotion = int32(p.MinimumEnvironmentsBeforePromotion.ValueInt64())
			phase.IsOptionalPhase = p.IsOptionalPhase.ValueBool()
			phase.IsPriorityPhase = p.IsPriorityPhase.ValueBool()
			phase.ReleaseRetentionPolicy = expandSpaceDefaultLifecycleRetention(ctx, p.ReleaseRetention, &diags)
			phase.TentacleRetentionPolicy = expandSpaceDefaultLifecycleRetention(ctx, p.TentacleRetention, &diags)
			existing.Phases = append(existing.Phases, phase)
		}
	}

	return existing, diags
}

func expandSpaceDefaultLifecycleRetention(ctx context.Context, value types.Object, diags *diag.Diagnostics) *core.RetentionPeriod {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var retention schemas.SpaceDefaultLifecycleRetentionModel
	diags.Append(value.As(ctx, &retention, basetypes.ObjectAsOptions{})...)

	switch retention.Strategy.ValueString() {
	case core.RetentionStrategyCount:
		return core.CountBasedRetentionPeriod(int32(retention.QuantityToKeep.ValueInt64()), retention.Unit.ValueString())
	case core.RetentionStrategyForever:
		return core.KeepForeverRetentionPeriod()
	default:
		return core.SpaceDefaultRetentionPeriod()
	}
}

// flattenSpaceDefaultLifecycle converts the lifecycle to its state value. prior is the planned or previous state and is used
// to keep optional values null when they were not configured and the server returned the implicit default.
func flattenSpaceDefaultLifecycle(ctx context.Context, lifecycle *lifecycles.Lifecycle, prior schemas.SpaceDefaultLifecycleModel) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := schemas.SpaceDefaultLifecycleModel{
		ID:                types.StringValue(lifecycle.GetID()),
		Name:              types.StringValue(lifecycle.Name),
		Description:       types.StringValue(lifecycle.Description),
		ReleaseRetention:  flattenSpaceDefaultLifecycleRetention(ctx, lifecycle.ReleaseRetentionPolicy, prior.ReleaseRetention, true, &diags),
		TentacleRetention: flattenSpaceDefaultLifecycleRetention(ctx, lifecycle.TentacleRetentionPolicy, prior.TentacleRetention, true, &diags),
	}

	if len(lifecycle.Phases) == 0 && prior.Phases.IsNull() {
		model.Phases = types.ListNull(spaceDefaultLifecyclePhaseObjectType)
	} else {
		var priorPhases []schemas.SpaceDefaultLifecyclePhaseModel
		if !prior.Phases.IsNull() && !prior.Phases.IsUnknown() {
			diags.Append(prior.Phases.ElementsAs(ctx, &priorPhases, false)...)
		}

		phases := make([]schemas.SpaceDefaultLifecyclePhaseModel, 0, len(lifecycle.Phases))
		for i, p := range lifecycle.Phases {
			priorPhase := schemas.SpaceDefaultLifecyclePhaseModel{
				ReleaseRetention:  types.ObjectNull(spaceDefaultLifecycleRetentionObjectType.AttrTypes),
				TentacleRetention: types.ObjectNull(spaceDefaultLifecycleRetentionObjectType.AttrTypes),
			}
			if i < len(priorPhases) {
				priorPhase = priorPhases[i]
			}
			phases = append(phases, schemas.SpaceDefaultLifecyclePhaseModel{
				ID:                                 types.StringValue(p.ID),
				Name:                               types.StringValue(p.Name),
				AutomaticDeploymentTargets:         util.FlattenStringList(p.AutomaticDeploymentTargets),
				OptionalDeploymentTargets:          util.FlattenStringList(p.OptionalDeploymentTargets),
				MinimumEnvironmentsBeforePromotion: types.Int64Value(int64(p.MinimumEnvironmentsBeforePromotion)),
				IsOptionalPhase:                    types.BoolValue(p.IsOptionalPhase),
				IsPriorityPhase:                    types.BoolValue(p.IsPriorityPhase),
				ReleaseRetention:                   flattenSpaceDefaultLifecycleRetention(ctx, p.ReleaseRetentionPolicy, priorPhase.ReleaseRetention, false, &diags),
				TentacleRetention:                  flattenSpaceDefaultLifecycleRetention(ctx, p.TentacleRetentionPolicy, priorPhase.TentacleRetention, false, &diags),
			})
		}
		list, d := types.ListValueFrom(ctx, spaceDefaultLifecyclePhaseObjectType, phases)
		diags.Append(d...)
		model.Phases = list
	}

	value, d := types.ObjectValueFrom(ctx, spaceDefaultLifecycleObjectType.AttrTypes, model)
	diags.Append(d...)
	return value, diags
}

func flattenSpaceDefaultLifecycleRetention(ctx context.Context, retention *core.RetentionPeriod, prior types.Object, nullWhenSpaceDefault bool, diags *diag.Diagnostics) types.Object {
	if retention == nil {
		return types.ObjectNull(spaceDefaultLifecycleRetentionObjectType.AttrTypes)
	}
	if nullWhenSpaceDefault && prior.IsNull() && retention.Strategy == core.RetentionStrategyDefault {
		return types.ObjectNull(spaceDefaultLifecycleRetentionObjectType.AttrTypes)
	}

	model := schemas.SpaceDefaultLifecycleRetentionModel{
		Strategy:       types.StringValue(retention.Strategy),
		QuantityToKeep: types.Int64Null(),
		Unit:           types.StringNull(),
	}
	if retention.Strategy == core.RetentionStrategyCount {
		model.QuantityToKeep = types.Int64Value(int64(retention.QuantityToKeep))
		model.Unit = types.StringValue(retention.Unit)

		// The API always returns capitalised units, keep the casing used in the configuration
		if !prior.IsNull() && !prior.IsUnknown() {
			var priorModel schemas.SpaceDefaultLifecycleRetentionModel
			diags.Append(prior.As(ctx, &priorModel, basetypes.ObjectAsOptions{})...)
			if strings.EqualFold(priorModel.Unit.ValueString(), retention.Unit) {
				model.Unit = priorModel.Unit
			}
		}
	}

	value, d := types.ObjectValueFrom(ctx, spaceDefaultLifecycleRetentionObjectType.AttrTypes, model)
	diags.Append(d...)
	return value
}

// getSpaceDefaultLifecycleID returns the lifecycle ID held in state, or an empty string when the lifecycle is not managed yet.
func getSpaceDefaultLifecycleID(ctx context.Context, state types.Object) string {
	if state.IsNull() || state.IsUnknown() {
		return ""
	}
	var model schemas.SpaceDefaultLifecycleModel
	if diags := state.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
		return ""
	}
	return model.ID.ValueString()
}

func isKnown(value types.String) bool {
	return !value.IsNull() && !value.IsUnknown()
}
