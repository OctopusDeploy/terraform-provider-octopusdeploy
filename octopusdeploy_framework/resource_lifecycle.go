package octopusdeploy_framework

import (
	"context"
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type lifecycleTypeResource struct {
	*Config
}

var _ resource.Resource = &lifecycleTypeResource{}
var _ resource.ResourceWithImportState = &lifecycleTypeResource{}

type lifecycleTypeResourceModel struct {
	SpaceID                   types.String `tfsdk:"space_id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	Phase                     types.List   `tfsdk:"phase"`
	ReleaseRetentionPolicy    types.List   `tfsdk:"release_retention_policy"`
	TentacleRetentionPolicy   types.List   `tfsdk:"tentacle_retention_policy"`
	ReleaseRetentionStrategy  types.List   `tfsdk:"release_retention_strategy"`
	TentacleRetentionStrategy types.List   `tfsdk:"tentacle_retention_strategy"`

	schemas.ResourceModel
}

func NewLifecycleResource() resource.Resource {
	return &lifecycleTypeResource{}
}

func (r *lifecycleTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")
	lifecycleID := idParts[len(idParts)-1]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), lifecycleID)...)
	// Note: This implementation assumes that space_id is set at the provider level
}

func (r *lifecycleTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("lifecycle")
}

func (r *lifecycleTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.LifecycleSchema{}.GetResourceSchema()
}

func (r *lifecycleTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *lifecycleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	releaseRetentionPolicySet, tentacleRetentionPolicySet, defaultPolicy := setDefaultRetentionPolicy(data)
	releaseRetentionStrategySet, tentacleRetentionStrategySet, defaultStrategy := setDefaultRetentionStrategy(data)

	newLifecycle := expandLifecycle(data)
	lifecycle, err := lifecycles.Add(r.Config.Client, newLifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
		return
	}

	handleUnitCasing(lifecycle, newLifecycle)

	data = flattenLifecycleResource(lifecycle)

	removeDefaultRetentionPolicy(releaseRetentionPolicySet, data, defaultPolicy, tentacleRetentionPolicySet)
	removeDefaultRetentionStrategy(releaseRetentionStrategySet, data, defaultStrategy, tentacleRetentionStrategySet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	releaseRetentionPolicySet, tentacleRetentionPolicySet, defaultPolicy := setDefaultRetentionPolicy(data)
	releaseRetentionStrategySet, tentacleRetentionStrategySet, defaultStrategy := setDefaultRetentionStrategy(data)

	lifecycle, err := lifecycles.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "lifecycle"); err != nil {
			resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
		}
		return
	}

	handleUnitCasing(lifecycle, expandLifecycle(data))

	data = flattenLifecycleResource(lifecycle)

	removeDefaultRetentionPolicy(releaseRetentionPolicySet, data, defaultPolicy, tentacleRetentionPolicySet)
	removeDefaultRetentionStrategy(releaseRetentionStrategySet, data, defaultStrategy, tentacleRetentionStrategySet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *lifecycleTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	releaseRetentionPolicySet, tentacleRetentionPolicySet, defaultPolicy := setDefaultRetentionPolicy(data)
	releaseRetentionStrategySet, tentacleRetentionStrategySet, defaultStrategy := setDefaultRetentionStrategy(data)

	tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", data.ID.ValueString()))

	lifecycle := expandLifecycle(data)
	lifecycle.ID = state.ID.ValueString()

	updatedLifecycle, err := lifecycles.Update(r.Config.Client, lifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
		return
	}

	handleUnitCasing(updatedLifecycle, lifecycle)

	data = flattenLifecycleResource(updatedLifecycle)

	removeDefaultRetentionPolicy(releaseRetentionPolicySet, data, defaultPolicy, tentacleRetentionPolicySet)
	removeDefaultRetentionStrategy(releaseRetentionStrategySet, data, defaultStrategy, tentacleRetentionStrategySet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func handleUnitCasing(resource *lifecycles.Lifecycle, data *lifecycles.Lifecycle) {
	// Set state to the casing provided in the desired state, as the Api will always return capitalised units
	resource.ReleaseRetentionPolicy = updateRetentionPeriodUnit(resource.ReleaseRetentionPolicy, data.ReleaseRetentionPolicy.Unit)
	resource.TentacleRetentionPolicy = updateRetentionPeriodUnit(resource.TentacleRetentionPolicy, data.TentacleRetentionPolicy.Unit)

	if len(data.Phases) == 0 {
		return
	}

	for i, phase := range resource.Phases {
		if phase.ReleaseRetentionPolicy != nil && phase.ReleaseRetentionPolicy.Unit != "" {
			phase.ReleaseRetentionPolicy = updateRetentionPeriodUnit(phase.ReleaseRetentionPolicy, data.Phases[i].ReleaseRetentionPolicy.Unit)
		}
		if phase.TentacleRetentionPolicy != nil && phase.TentacleRetentionPolicy.Unit != "" {
			phase.TentacleRetentionPolicy = updateRetentionPeriodUnit(phase.TentacleRetentionPolicy, data.Phases[i].TentacleRetentionPolicy.Unit)
		}
	}
}

func updateRetentionPeriodUnit(retentionPeriodResource *core.RetentionPeriod, dataUnit string) *core.RetentionPeriod {
	if strings.EqualFold(retentionPeriodResource.Unit, dataUnit) {
		period := core.RetentionPeriod{
			QuantityToKeep:    retentionPeriodResource.QuantityToKeep,
			ShouldKeepForever: retentionPeriodResource.ShouldKeepForever,
			Unit:              dataUnit,
			Strategy:          retentionPeriodResource.Strategy,
		}

		return &period
	}

	return retentionPeriodResource
}

func removeDefaultRetentionPolicy(releaseRetentionPolicySet bool, data *lifecycleTypeResourceModel, defaultPolicy types.List, tentacleRetentionPolicySet bool) {
	// Remove default policies from data before setting state, but only if we added them
	if !releaseRetentionPolicySet && data.ReleaseRetentionPolicy.Equal(defaultPolicy) {
		data.ReleaseRetentionPolicy = types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})
	}
	if !tentacleRetentionPolicySet && data.TentacleRetentionPolicy.Equal(defaultPolicy) {
		data.TentacleRetentionPolicy = types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})
	}
}
func removeDefaultRetentionStrategy(releaseRetentionStrategySet bool, data *lifecycleTypeResourceModel, defaultStrategy types.List, tentacleRetentionStrategySet bool) {
	// Remove default policies from data before setting state, but only if we added them
	if !releaseRetentionStrategySet && data.ReleaseRetentionStrategy.Equal(defaultStrategy) {
		data.ReleaseRetentionStrategy = types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})
	}
	if !tentacleRetentionStrategySet && data.TentacleRetentionStrategy.Equal(defaultStrategy) {
		data.TentacleRetentionStrategy = types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})
	}
}

func setDefaultRetentionPolicy(data *lifecycleTypeResourceModel) (bool, bool, types.List) {
	releaseRetentionPolicySet := !data.ReleaseRetentionPolicy.IsNull() && len(data.ReleaseRetentionPolicy.Elements()) > 0
	tentacleRetentionPolicySet := !data.TentacleRetentionPolicy.IsNull() && len(data.TentacleRetentionPolicy.Elements()) > 0

	// Set default policies only if they're not in the plan
	defaultPolicy := flattenRetentionPeriod(core.NewRetentionPeriod(30, "Days", false))
	if !releaseRetentionPolicySet {
		data.ReleaseRetentionPolicy = defaultPolicy
	}
	if !tentacleRetentionPolicySet {
		data.TentacleRetentionPolicy = defaultPolicy
	}
	return releaseRetentionPolicySet, tentacleRetentionPolicySet, defaultPolicy
}

func setDefaultRetentionStrategy(data *lifecycleTypeResourceModel) (bool, bool, types.List) {
	releaseRetentionStrategySet := !data.ReleaseRetentionStrategy.IsNull() && len(data.ReleaseRetentionStrategy.Elements()) > 0
	tentacleRetentionStrategySet := !data.TentacleRetentionStrategy.IsNull() && len(data.TentacleRetentionStrategy.Elements()) > 0

	// Set default Strategy only if they're not in the plan
	defaultStrategy := flattenRetentionPeriod(core.SpaceDefaultRetentionPeriod())
	if !releaseRetentionStrategySet {
		data.ReleaseRetentionStrategy = defaultStrategy
	}
	if !tentacleRetentionStrategySet {
		data.TentacleRetentionStrategy = defaultStrategy
	}
	return releaseRetentionStrategySet, tentacleRetentionStrategySet, defaultStrategy
}

func (r *lifecycleTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data lifecycleTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := lifecycles.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete lifecycle", err.Error())
		return
	}
}

func resourceConfiguration(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *Config {
	if req.ProviderData == nil {
		return nil
	}

	p, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return p
}

func flattenLifecycleResource(lifecycle *lifecycles.Lifecycle) *lifecycleTypeResourceModel {
	flattenedLifecycle := &lifecycleTypeResourceModel{
		SpaceID:                   types.StringValue(lifecycle.SpaceID),
		Name:                      types.StringValue(lifecycle.Name),
		Description:               types.StringValue(lifecycle.Description),
		Phase:                     flattenPhases(lifecycle.Phases),
		ReleaseRetentionPolicy:    flattenRetentionPeriod(lifecycle.ReleaseRetentionPolicy),
		TentacleRetentionPolicy:   flattenRetentionPeriod(lifecycle.TentacleRetentionPolicy),
		ReleaseRetentionStrategy:  flattenRetentionStrategy(lifecycle.ReleaseRetentionStrategy),
		TentacleRetentionStrategy: flattenRetentionStrategy(lifecycle.TentacleRetentionStrategy),
	}
	flattenedLifecycle.ID = types.StringValue(lifecycle.GetID())

	return flattenedLifecycle
}

func flattenPhases(phases []*lifecycles.Phase) types.List {
	if phases == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getPhaseAttrTypes()})
	}
	phasesList := make([]attr.Value, 0, len(phases))

	for _, phase := range phases {
		attrs := map[string]attr.Value{
			"id":                                    types.StringValue(phase.ID),
			"name":                                  types.StringValue(phase.Name),
			"automatic_deployment_targets":          util.FlattenStringList(phase.AutomaticDeploymentTargets),
			"optional_deployment_targets":           util.FlattenStringList(phase.OptionalDeploymentTargets),
			"minimum_environments_before_promotion": types.Int64Value(int64(phase.MinimumEnvironmentsBeforePromotion)),
			"is_optional_phase":                     types.BoolValue(phase.IsOptionalPhase),
			"is_priority_phase":                     types.BoolValue(phase.IsPriorityPhase),
			"release_retention_policy":              util.Ternary(phase.ReleaseRetentionPolicy != nil, flattenRetentionPeriod(phase.ReleaseRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})),
			"tentacle_retention_policy":             util.Ternary(phase.TentacleRetentionPolicy != nil, flattenRetentionPeriod(phase.TentacleRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})),
			"release_retention_strategy":            util.Ternary(phase.ReleaseRetentionStrategy != nil, flattenRetentionStrategy(phase.ReleaseRetentionStrategy), types.ListNull(types.ObjectType{AttrTypes: getRetentionStrategyAttrTypes()})),
			"tentacle_retention_strategy":           util.Ternary(phase.TentacleRetentionStrategy != nil, flattenRetentionStrategy(phase.TentacleRetentionStrategy), types.ListNull(types.ObjectType{AttrTypes: getRetentionStrategyAttrTypes()})),
		}
		phasesList = append(phasesList, types.ObjectValueMust(getPhaseAttrTypes(), attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: getPhaseAttrTypes()}, phasesList)
}

func flattenRetentionPeriod(retentionPeriod *core.RetentionPeriod) types.List {
	if retentionPeriod == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
		[]attr.Value{
			types.ObjectValueMust(
				getRetentionPeriodAttrTypes(),
				map[string]attr.Value{
					"quantity_to_keep":    types.Int64Value(int64(retentionPeriod.QuantityToKeep)),
					"should_keep_forever": types.BoolValue(retentionPeriod.ShouldKeepForever),
					"unit":                types.StringValue(retentionPeriod.Unit),
				},
			),
		},
	)
}

func flattenRetentionStrategy(retentionPeriod *core.RetentionPeriod) types.List {
	if retentionPeriod == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
		[]attr.Value{
			types.ObjectValueMust(
				getRetentionPeriodAttrTypes(),
				map[string]attr.Value{
					"strategy":         types.StringValue(retentionPeriod.Strategy),
					"quantity_to_keep": types.Int64Value(int64(retentionPeriod.QuantityToKeep)),
					"unit":             types.StringValue(retentionPeriod.Unit),
				},
			),
		},
	)
}

func expandLifecycle(data *lifecycleTypeResourceModel) *lifecycles.Lifecycle {
	if data == nil {
		return nil
	}

	lifecycle := lifecycles.NewLifecycle(data.Name.ValueString())
	lifecycle.Description = data.Description.ValueString()
	lifecycle.SpaceID = data.SpaceID.ValueString()
	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		lifecycle.ID = data.ID.ValueString()
	}

	lifecycle.Phases = expandPhases(data.Phase)
	lifecycle.ReleaseRetentionPolicy = expandRetentionPeriod(data.ReleaseRetentionPolicy)
	lifecycle.TentacleRetentionPolicy = expandRetentionPeriod(data.TentacleRetentionPolicy)
	lifecycle.ReleaseRetentionStrategy = expandRetentionStrategy(data.ReleaseRetentionStrategy)
	lifecycle.TentacleRetentionStrategy = expandRetentionStrategy(data.TentacleRetentionStrategy)

	return lifecycle
}

func expandPhases(phases types.List) []*lifecycles.Phase {
	if phases.IsNull() || phases.IsUnknown() || len(phases.Elements()) == 0 {
		return nil
	}

	result := make([]*lifecycles.Phase, 0, len(phases.Elements()))

	for _, phaseElem := range phases.Elements() {
		phaseObj := phaseElem.(types.Object)
		phaseAttrs := phaseObj.Attributes()

		phase := &lifecycles.Phase{}

		if v, ok := phaseAttrs["id"].(types.String); ok && !v.IsNull() {
			phase.ID = v.ValueString()
		}

		if v, ok := phaseAttrs["name"].(types.String); ok && !v.IsNull() {
			phase.Name = v.ValueString()
		}

		if v, ok := phaseAttrs["automatic_deployment_targets"].(types.List); ok && !v.IsNull() {
			phase.AutomaticDeploymentTargets = util.ExpandStringList(v)
		}

		if v, ok := phaseAttrs["optional_deployment_targets"].(types.List); ok && !v.IsNull() {
			phase.OptionalDeploymentTargets = util.ExpandStringList(v)
		}

		if v, ok := phaseAttrs["minimum_environments_before_promotion"].(types.Int64); ok && !v.IsNull() {
			phase.MinimumEnvironmentsBeforePromotion = int32(v.ValueInt64())
		}

		if v, ok := phaseAttrs["is_optional_phase"].(types.Bool); ok && !v.IsNull() {
			phase.IsOptionalPhase = v.ValueBool()
		}

		if v, ok := phaseAttrs["is_priority_phase"].(types.Bool); ok && !v.IsNull() {
			phase.IsPriorityPhase = v.ValueBool()
		}

		if v, ok := phaseAttrs["release_retention_policy"].(types.List); ok && !v.IsNull() {
			phase.ReleaseRetentionPolicy = expandRetentionPeriod(v)
		}

		if v, ok := phaseAttrs["tentacle_retention_policy"].(types.List); ok && !v.IsNull() {
			phase.TentacleRetentionPolicy = expandRetentionPeriod(v)
		}

		if v, ok := phaseAttrs["release_retention_strategy"].(types.List); ok && !v.IsNull() {
			phase.ReleaseRetentionStrategy = expandRetentionStrategy(v)
		}

		if v, ok := phaseAttrs["tentacle_retention_strategy"].(types.List); ok && !v.IsNull() {
			phase.TentacleRetentionStrategy = expandRetentionStrategy(v)
		}

		result = append(result, phase)
	}

	return result
}

func expandRetentionPeriod(v types.List) *core.RetentionPeriod {
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

func expandRetentionStrategy(v types.List) *core.RetentionPeriod {
	if v.IsNull() || v.IsUnknown() || len(v.Elements()) == 0 {
		return nil
	}

	obj := v.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	var strategy string
	if s, ok := attrs["strategy"].(types.String); ok && !s.IsNull() {
		strategy = s.ValueString()
	}

	if strategy == core.RetentionStrategyCount {
		var unit string
		if u, ok := attrs["unit"].(types.String); ok && !u.IsNull() {
			unit = u.ValueString()
		}

		var quantityToKeep int32
		if qty, ok := attrs["quantity_to_keep"].(types.Int64); ok && !qty.IsNull() {
			quantityToKeep = int32(qty.ValueInt64())
		}
		return core.CountBasedRetentionPeriod(quantityToKeep, unit)
	}

	if strategy == core.RetentionStrategyForever {
		return core.KeepForeverRetentionPeriod()
	}
	return core.SpaceDefaultRetentionPeriod()

}

func getRetentionPeriodAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}

func getRetentionStrategyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"strategy":         types.StringType,
		"quantity_to_keep": types.Int64Type,
		"unit":             types.StringType,
	}
}

func getPhaseAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"is_priority_phase":                     types.BoolType,
		"release_retention_policy":              types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()}},
		"tentacle_retention_policy":             types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()}},
		"release_retention_strategy":            types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionStrategyAttrTypes()}},
		"tentacle_retention_strategy":           types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionStrategyAttrTypes()}},
	}
}
