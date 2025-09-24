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
	serverSupportsRetentionWithStrategy bool
}

var _ resource.Resource = &lifecycleTypeResource{}
var _ resource.ResourceWithImportState = &lifecycleTypeResource{}

type lifecycleTypeResourceModel struct {
	SpaceID                       types.String `tfsdk:"space_id"`
	Name                          types.String `tfsdk:"name"`
	Description                   types.String `tfsdk:"description"`
	Phase                         types.List   `tfsdk:"phase"`
	ReleaseRetention              types.List   `tfsdk:"release_retention_policy"`
	TentacleRetention             types.List   `tfsdk:"tentacle_retention_policy"`
	ReleaseRetentionWithStrategy  types.List   `tfsdk:"release_retention_with_strategy"`
	TentacleRetentionWithStrategy types.List   `tfsdk:"tentacle_retention_with_strategy"`

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

func (r *lifecycleTypeResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.LifecycleSchema{}.GetResourceSchema()
}

func (r *lifecycleTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
	if r.Config != nil {
		r.serverSupportsRetentionWithStrategy = r.Config.IsVersionSameOrGreaterThan("2025.3")
	}

}

func (r *lifecycleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	isRetentionWithStrategyUsed := determineIfRetentionWithStrategyIsUsed(data, r.serverSupportsRetentionWithStrategy)
	hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := setDefaultRetention(data, isRetentionWithStrategyUsed)
	hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy := setDefaultRetentionWithStrategy(data, isRetentionWithStrategyUsed)
	//if retentionWithStrategy is used then set the default retention usining new, otherwise set it using old
	tflog.Debug(ctx, fmt.Sprintf("updating lifecycle Rose '%+v'", data))
	newLifecycle := expandLifecycle(data)
	lifecycle, err := lifecycles.Add(r.Config.Client, newLifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
		return
	}

	handleUnitCasing(lifecycle, newLifecycle)

	data = flattenLifecycleResource(lifecycle)
	//if retentionWithStrategy is used then remove all of the old ones and reset the retentionWithStrategy to the default
	removeDefaultRetentionFromUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention)
	removeDefaultRetentionWithStrategyFromUnsetBlocks(data, hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	isRetentionWithStrategyUsed := determineIfRetentionWithStrategyIsUsed(data, r.serverSupportsRetentionWithStrategy)
	hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := setDefaultRetention(data, isRetentionWithStrategyUsed)
	hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy := setDefaultRetentionWithStrategy(data, isRetentionWithStrategyUsed)

	lifecycle, err := lifecycles.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "lifecycle"); err != nil {
			resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
		}
		return
	}

	handleUnitCasing(lifecycle, expandLifecycle(data))

	data = flattenLifecycleResource(lifecycle)

	removeDefaultRetentionFromUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention)
	removeDefaultRetentionWithStrategyFromUnsetBlocks(data, hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *lifecycleTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	isRetentionWithStrategyUsed := determineIfRetentionWithStrategyIsUsed(data, r.serverSupportsRetentionWithStrategy)
	hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := setDefaultRetention(data, isRetentionWithStrategyUsed)
	hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy := setDefaultRetentionWithStrategy(data, isRetentionWithStrategyUsed)

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

	removeDefaultRetentionFromUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention)
	removeDefaultRetentionWithStrategyFromUnsetBlocks(data, hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func handleUnitCasing(resource *lifecycles.Lifecycle, data *lifecycles.Lifecycle) {
	// Set state to the casing provided in the desired state, as the Api will always return capitalised units
	resource.ReleaseRetentionPolicy = updateRetentionUnit(resource.ReleaseRetentionPolicy, data.ReleaseRetentionPolicy.Unit)
	resource.TentacleRetentionPolicy = updateRetentionUnit(resource.TentacleRetentionPolicy, data.TentacleRetentionPolicy.Unit)

	if len(data.Phases) == 0 {
		return
	}

	for i, phase := range resource.Phases {
		if phase.ReleaseRetentionPolicy != nil && phase.ReleaseRetentionPolicy.Unit != "" {
			phase.ReleaseRetentionPolicy = updateRetentionUnit(phase.ReleaseRetentionPolicy, data.Phases[i].ReleaseRetentionPolicy.Unit)
		}
		if phase.TentacleRetentionPolicy != nil && phase.TentacleRetentionPolicy.Unit != "" {
			phase.TentacleRetentionPolicy = updateRetentionUnit(phase.TentacleRetentionPolicy, data.Phases[i].TentacleRetentionPolicy.Unit)
		}
	}
}

func updateRetentionUnit(retentionResource *core.RetentionPeriod, dataUnit string) *core.RetentionPeriod {
	if strings.EqualFold(retentionResource.Unit, dataUnit) {
		retention := core.RetentionPeriod{
			QuantityToKeep:    retentionResource.QuantityToKeep,
			ShouldKeepForever: retentionResource.ShouldKeepForever,
			Unit:              dataUnit,
			Strategy:          retentionResource.Strategy,
		}

		return &retention
	}

	return retentionResource
}

func removeDefaultRetentionFromUnsetBlocks(data *lifecycleTypeResourceModel, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention bool) {
	// Remove default policies from data before setting state, but only if we added them
	if !hasUserDefinedReleaseRetention {
		data.ReleaseRetention = types.ListNull(types.ObjectType{AttrTypes: getRetentionAttrTypes()})
	}
	if !hasUserDefinedTentacleRetention {
		data.TentacleRetention = types.ListNull(types.ObjectType{AttrTypes: getRetentionAttrTypes()})
	}
}
func removeDefaultRetentionWithStrategyFromUnsetBlocks(data *lifecycleTypeResourceModel, hasUserDefinedReleaseRetentionWithStrategy bool, hasUserDefinedTentacleRetentionWithStrategy bool) {
	// Remove default policies from data before setting state, but only if we added them
	if !hasUserDefinedReleaseRetentionWithStrategy {
		data.ReleaseRetentionWithStrategy = types.ListNull(types.ObjectType{AttrTypes: getRetentionAttrTypes()})
	}
	if !hasUserDefinedTentacleRetentionWithStrategy {
		data.TentacleRetentionWithStrategy = types.ListNull(types.ObjectType{AttrTypes: getRetentionAttrTypes()})
	}
}

func setDefaultRetention(data *lifecycleTypeResourceModel, isRetentionWithStrategyUsed bool) (bool, bool) {
	hasUserDefinedReleaseRetention := !data.ReleaseRetention.IsNull() && len(data.ReleaseRetention.Elements()) > 0
	hasUserDefinedTentacleRetention := !data.TentacleRetention.IsNull() && len(data.TentacleRetention.Elements()) > 0

	var initialRetentionSetting types.List
	if isRetentionWithStrategyUsed {
		initialRetentionSetting = types.ListNull(types.ObjectType{AttrTypes: getRetentionAttrTypes()})
	} else {
		initialRetentionSetting = flattenRetention(core.NewRetentionPeriod(30, "Days", false))
	}

	if !hasUserDefinedReleaseRetention {
		data.ReleaseRetention = initialRetentionSetting
	}
	if !hasUserDefinedTentacleRetention {
		data.TentacleRetention = initialRetentionSetting
	}

	return hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention
}

func setDefaultRetentionWithStrategy(data *lifecycleTypeResourceModel, isRetentionWithStrategyUsed bool) (bool, bool) {
	hasUserDefinedReleaseRetentionWithStrategy := !data.ReleaseRetentionWithStrategy.IsNull() && len(data.ReleaseRetentionWithStrategy.Elements()) > 0
	hasUserDefinedTentacleRetentionWithStrategy := !data.TentacleRetentionWithStrategy.IsNull() && len(data.TentacleRetentionWithStrategy.Elements()) > 0

	var initialRetentionSetting types.List
	if !isRetentionWithStrategyUsed {
		initialRetentionSetting = types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})
	} else {
		initialRetentionSetting = flattenRetention(core.SpaceDefaultRetentionPeriod())
	}
	//initialRetentionSetting = flattenRetention(core.NewRetentionPeriod(30, "Days", false))

	if !hasUserDefinedReleaseRetentionWithStrategy {
		data.ReleaseRetentionWithStrategy = initialRetentionSetting
	}
	if !hasUserDefinedTentacleRetentionWithStrategy {
		data.TentacleRetentionWithStrategy = initialRetentionSetting
	}

	return hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy
}

func determineIfRetentionWithStrategyIsUsed(data *lifecycleTypeResourceModel, serverSupportsRetentionWithStrategy bool) bool {
	if !data.ReleaseRetentionWithStrategy.IsNull() && len(data.ReleaseRetentionWithStrategy.Elements()) > 0 {
		return true
	}
	if !data.TentacleRetentionWithStrategy.IsNull() && len(data.TentacleRetentionWithStrategy.Elements()) > 0 {
		return true
	}
	if !data.ReleaseRetention.IsNull() && len(data.ReleaseRetention.Elements()) > 0 {
		return false
	}
	if !data.TentacleRetention.IsNull() && len(data.TentacleRetention.Elements()) > 0 {
		return false
	}
	phases := data.Phase.Elements()
	for _, phase := range phases {
		phaseObj := phase.(types.Object)
		phaseAttrs := phaseObj.Attributes()
		releaseRetentionWithStrategy := phaseAttrs["release_retention_with_strategy"].(types.List)
		tentacleRetentionWithStrategy := phaseAttrs["tentacle_retention_with_strategy"].(types.List)
		releaseRetention := phaseAttrs["release_retention_policy"].(types.List)
		tentacleRetention := phaseAttrs["tentacle_retention_policy"].(types.List)
		if !releaseRetentionWithStrategy.IsNull() && len(releaseRetentionWithStrategy.Elements()) > 0 {
			return true
		}
		if !tentacleRetentionWithStrategy.IsNull() && len(tentacleRetentionWithStrategy.Elements()) > 0 {
			return true
		}
		if !releaseRetention.IsNull() && len(data.ReleaseRetention.Elements()) > 0 {
			return false
		}
		if !tentacleRetention.IsNull() && len(data.TentacleRetention.Elements()) > 0 {
			return false
		}
	}
	if !serverSupportsRetentionWithStrategy {
		return true
	}
	return false

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
	//TODO: add an if to only change the retention attributes that are being altered in the plan
	//TODO: get access to the plan here
	flattenedLifecycle := &lifecycleTypeResourceModel{
		SpaceID:                       types.StringValue(lifecycle.SpaceID),
		Name:                          types.StringValue(lifecycle.Name),
		Description:                   types.StringValue(lifecycle.Description),
		Phase:                         flattenPhases(lifecycle.Phases),
		ReleaseRetention:              flattenRetention(lifecycle.ReleaseRetentionPolicy),
		TentacleRetention:             flattenRetention(lifecycle.TentacleRetentionPolicy),
		ReleaseRetentionWithStrategy:  flattenRetentionWithStrategy(lifecycle.ReleaseRetentionPolicy),
		TentacleRetentionWithStrategy: flattenRetentionWithStrategy(lifecycle.TentacleRetentionPolicy),
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
			"release_retention":                     util.Ternary(phase.ReleaseRetentionPolicy != nil, flattenRetention(phase.ReleaseRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionAttrTypes()})),
			"tentacle_retention":                    util.Ternary(phase.TentacleRetentionPolicy != nil, flattenRetention(phase.TentacleRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionAttrTypes()})),
			"release_retention_with_strategy":       util.Ternary(phase.ReleaseRetentionPolicy != nil, flattenRetentionWithStrategy(phase.ReleaseRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})),
			"tentacle_retention_with_strategy":      util.Ternary(phase.TentacleRetentionPolicy != nil, flattenRetentionWithStrategy(phase.TentacleRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})),
		}
		phasesList = append(phasesList, types.ObjectValueMust(getPhaseAttrTypes(), attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: getPhaseAttrTypes()}, phasesList)
}

func flattenRetention(retention *core.RetentionPeriod) types.List {
	if retention == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getRetentionAttrTypes()})
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: getRetentionAttrTypes()},
		[]attr.Value{
			types.ObjectValueMust(
				getRetentionAttrTypes(),
				map[string]attr.Value{
					"quantity_to_keep":    types.Int64Value(int64(retention.QuantityToKeep)),
					"should_keep_forever": types.BoolValue(retention.ShouldKeepForever),
					"unit":                types.StringValue(retention.Unit),
				},
			),
		},
	)
}

func flattenRetentionWithStrategy(retentionWithStrategy *core.RetentionPeriod) types.List {
	if retentionWithStrategy == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()},
		[]attr.Value{
			types.ObjectValueMust(
				getRetentionWithStrategyAttrTypes(),
				map[string]attr.Value{
					"strategy":         types.StringValue(retentionWithStrategy.Strategy),
					"unit":             types.StringValue(retentionWithStrategy.Unit),
					"quantity_to_keep": types.Int64Value(int64(retentionWithStrategy.QuantityToKeep)),
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
	lifecycle.ReleaseRetentionPolicy = expandRetentionWithOrWithoutStrategy(data.ReleaseRetention, data.ReleaseRetentionWithStrategy)
	lifecycle.TentacleRetentionPolicy = expandRetentionWithOrWithoutStrategy(data.TentacleRetention, data.TentacleRetentionWithStrategy)

	return lifecycle
}

func expandPhases(phases types.List) []*lifecycles.Phase {
	//TODO: don't allow with strategy to be used with the pre 2025.3 octopus
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
			phase.ReleaseRetentionPolicy = expandRetention(v)
		}
		if v, ok := phaseAttrs["release_retention_with_strategy"].(types.List); ok && !v.IsNull() {
			phase.ReleaseRetentionPolicy = expandRetentionWithStrategy(v)
		}
		if v, ok := phaseAttrs["tentacle_retention_policy"].(types.List); ok && !v.IsNull() {
			phase.TentacleRetentionPolicy = expandRetention(v)
		}
		if v, ok := phaseAttrs["release_retention_with_strategy"].(types.List); ok && !v.IsNull() {
			phase.TentacleRetentionPolicy = expandRetentionWithStrategy(v)
		}
		result = append(result, phase)
	}

	return result
}

func expandRetentionWithOrWithoutStrategy(oldRetention types.List, newRetention types.List) *core.RetentionPeriod {
	oldRetentionPresent := !oldRetention.IsNull() && !oldRetention.IsUnknown()
	newRetentionPresent := !newRetention.IsNull() && !newRetention.IsUnknown()
	if !oldRetentionPresent && !newRetentionPresent {
		return nil
	}
	if oldRetentionPresent {
		return expandRetention(oldRetention)
	}
	return expandRetentionWithStrategy(newRetention)
}
func expandRetention(v types.List) *core.RetentionPeriod {
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

func expandRetentionWithStrategy(v types.List) *core.RetentionPeriod {
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

func getRetentionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}

func getRetentionWithStrategyAttrTypes() map[string]attr.Type {
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
		"release_retention_policy":              types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionAttrTypes()}},
		"tentacle_retention_policy":             types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionAttrTypes()}},
		"release_retention_with_strategy":       types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}},
		"tentacle_retention_with_strategy":      types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}},
	}
}
