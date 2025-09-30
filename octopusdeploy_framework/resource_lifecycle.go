package octopusdeploy_framework

import (
	"context"
	"fmt"
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
	"strings"
)

type lifecycleTypeResource struct {
	*Config
	onlyDeprecatedRetentionIsSupportedByServer bool
	allowDeprecatedRetention                   bool
}

var _ resource.Resource = &lifecycleTypeResource{}
var _ resource.ResourceWithImportState = &lifecycleTypeResource{}

type lifecycleTypeResourceModel struct {
	SpaceID                       types.String `tfsdk:"space_id"`
	Name                          types.String `tfsdk:"name"`
	Description                   types.String `tfsdk:"description"`
	Phase                         types.List   `tfsdk:"phase"`
	DeprecatedReleaseRetention    types.List   `tfsdk:"release_retention_policy"`
	DeprecatedTentacleRetention   types.List   `tfsdk:"tentacle_retention_policy"`
	ReleaseRetentionWithStrategy  types.List   `tfsdk:"release_retention_with_strategy"`
	TentacleRetentionWithStrategy types.List   `tfsdk:"tentacle_retention_with_strategy"`

	schemas.ResourceModel
}

func NewLifecycleResource() resource.Resource {
	return &lifecycleTypeResource{allowDeprecatedRetention: schemas.AllowDeprecatedRetentionFeatureFlag}
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
	resp.Schema = schemas.LifecycleSchema{AllowDeprecatedRetention: r.allowDeprecatedRetention}.GetResourceSchema()
}

func (r *lifecycleTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
	if r.Config != nil {
		r.onlyDeprecatedRetentionIsSupportedByServer = !r.Config.IsVersionSameOrGreaterThan("2025.3") // this method always returns true if running on the local
	}
}

func (r *lifecycleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var onlyDeprecatedRetentionIsUsed = false
	var initialRetentionWithStrategySetting = flattenRetentionWithStrategy(core.SpaceDefaultRetentionPeriod())
	var initialDeprecatedRetentionSetting = ListNullDeprecatedRetention
	if r.allowDeprecatedRetention {
		ValidateRetentionBlocksUsed(data, &resp.Diagnostics, r.onlyDeprecatedRetentionIsSupportedByServer)
		onlyDeprecatedRetentionIsUsed = IsOnlyDeprecatedRetentionUsed(data, r.onlyDeprecatedRetentionIsSupportedByServer)
		if onlyDeprecatedRetentionIsUsed {
			initialRetentionWithStrategySetting = ListNullRetentionWithStrategy
			initialDeprecatedRetentionSetting = FlattenDeprecatedRetention(core.CountBasedRetentionPeriod(30, "Days"))
		}
	}
	hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := SetDeprecatedDefaultRetention(data, initialDeprecatedRetentionSetting)
	hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy := setDefaultRetentionWithStrategy(data, initialRetentionWithStrategySetting)

	newLifecycle := expandLifecycle(data, onlyDeprecatedRetentionIsUsed)
	lifecycle, err := lifecycles.Add(r.Config.Client, newLifecycle)

	if err != nil {
		resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
		return
	}

	for _, goPhase := range lifecycle.Phases {
		tflog.Debug(ctx, fmt.Sprintf("ran add rose retention'%v'tentacle `%v`", goPhase.ReleaseRetentionPolicy, goPhase.TentacleRetentionPolicy))
	}

	data = flattenLifecycleResource(lifecycle, onlyDeprecatedRetentionIsUsed)

	RemoveDeprecatedDefaultRetentionFromUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention)
	removeDefaultRetentionWithStrategyFromUnsetBlocks(data, hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy)

	tflog.Debug(ctx, fmt.Sprintf("removed deprecated retention '%v' and '%v'", data.DeprecatedTentacleRetention, data.DeprecatedReleaseRetention))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, fmt.Sprintf("set more data '%s'", data.ReleaseRetentionWithStrategy.Type(ctx)))
}

func (r *lifecycleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var initialRetentionWithStrategySetting = flattenRetentionWithStrategy(core.SpaceDefaultRetentionPeriod())
	var onlyDeprecatedRetentionIsUsed bool

	//create one function which sets deprecated default retention and outputs
	var initialDeprecatedRetentionSetting = ListNullDeprecatedRetention
	if r.allowDeprecatedRetention {
		ValidateRetentionBlocksUsed(data, &resp.Diagnostics, r.onlyDeprecatedRetentionIsSupportedByServer)
		onlyDeprecatedRetentionIsUsed = IsOnlyDeprecatedRetentionUsed(data, r.onlyDeprecatedRetentionIsSupportedByServer)
		if onlyDeprecatedRetentionIsUsed {
			initialRetentionWithStrategySetting = ListNullRetentionWithStrategy
			initialDeprecatedRetentionSetting = FlattenDeprecatedRetention(core.CountBasedRetentionPeriod(30, "Days"))
		}
	}
	hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := SetDeprecatedDefaultRetention(data, initialDeprecatedRetentionSetting)

	hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy := setDefaultRetentionWithStrategy(data, initialRetentionWithStrategySetting)

	lifecycle, err := lifecycles.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "lifecycle"); err != nil {
			resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
		}
		return
	}

	data = flattenLifecycleResource(lifecycle, onlyDeprecatedRetentionIsUsed)

	RemoveDeprecatedDefaultRetentionFromUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention)
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

	var onlyDeprecatedRetentionIsUsed bool
	var initialRetentionWithStrategySetting = flattenRetentionWithStrategy(core.SpaceDefaultRetentionPeriod())
	var initialDeprecatedRetentionSetting = ListNullDeprecatedRetention
	if r.allowDeprecatedRetention {
		ValidateRetentionBlocksUsed(data, &resp.Diagnostics, r.onlyDeprecatedRetentionIsSupportedByServer)
		onlyDeprecatedRetentionIsUsed = IsOnlyDeprecatedRetentionUsed(data, r.onlyDeprecatedRetentionIsSupportedByServer)
		if onlyDeprecatedRetentionIsUsed {
			initialRetentionWithStrategySetting = ListNullRetentionWithStrategy
			initialDeprecatedRetentionSetting = FlattenDeprecatedRetention(core.CountBasedRetentionPeriod(30, "Days"))
		}
	}
	hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := SetDeprecatedDefaultRetention(data, initialDeprecatedRetentionSetting)
	hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy := setDefaultRetentionWithStrategy(data, initialRetentionWithStrategySetting)

	tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", data.ID.ValueString()))

	lifecycle := expandLifecycle(data, onlyDeprecatedRetentionIsUsed)
	lifecycle.ID = state.ID.ValueString()

	updatedLifecycle, err := lifecycles.Update(r.Config.Client, lifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
		return
	}

	data = flattenLifecycleResource(updatedLifecycle, onlyDeprecatedRetentionIsUsed)

	RemoveDeprecatedDefaultRetentionFromUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention)
	removeDefaultRetentionWithStrategyFromUnsetBlocks(data, hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

func setDefaultRetentionWithStrategy(data *lifecycleTypeResourceModel, initialDeprecatedRetentionSetting types.List) (bool, bool) {
	isReleaseRetentionWithStrategyDefined := attributeIsUsed(data.ReleaseRetentionWithStrategy)
	isTentacleRetentionWithStrategyDefined := attributeIsUsed(data.TentacleRetentionWithStrategy)

	if !isReleaseRetentionWithStrategyDefined {
		data.ReleaseRetentionWithStrategy = initialDeprecatedRetentionSetting
	}
	if !isTentacleRetentionWithStrategyDefined {
		data.TentacleRetentionWithStrategy = initialDeprecatedRetentionSetting
	}

	return isReleaseRetentionWithStrategyDefined, isTentacleRetentionWithStrategyDefined
}

func removeDefaultRetentionWithStrategyFromUnsetBlocks(data *lifecycleTypeResourceModel, hasUserDefinedReleaseRetentionWithStrategy bool, hasUserDefinedTentacleRetentionWithStrategy bool) {
	// Remove retention policies from data before setting state, but only if we added the initial value to them in the first place
	if !hasUserDefinedReleaseRetentionWithStrategy {
		data.ReleaseRetentionWithStrategy = ListNullRetentionWithStrategy
	}
	if !hasUserDefinedTentacleRetentionWithStrategy {
		data.TentacleRetentionWithStrategy = ListNullRetentionWithStrategy
	}
}

func attributeIsUsed(attribute types.List) bool {
	if !attribute.IsNull() && len(attribute.Elements()) > 0 {
		return true
	}
	return false
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

func flattenLifecycleResource(lifecycle *lifecycles.Lifecycle, onlyDeprecatedRetentionIsUsed bool) *lifecycleTypeResourceModel {
	var flattenedLifecycle *lifecycleTypeResourceModel
	if onlyDeprecatedRetentionIsUsed {
		flattenedLifecycle = &lifecycleTypeResourceModel{
			SpaceID:                     types.StringValue(lifecycle.SpaceID),
			Name:                        types.StringValue(lifecycle.Name),
			Description:                 types.StringValue(lifecycle.Description),
			Phase:                       DeprecatedFlattenPhases(lifecycle.Phases),
			DeprecatedReleaseRetention:  FlattenDeprecatedRetention(lifecycle.ReleaseRetentionPolicy),
			DeprecatedTentacleRetention: FlattenDeprecatedRetention(lifecycle.TentacleRetentionPolicy),
		}
	} else {
		flattenedLifecycle = &lifecycleTypeResourceModel{
			SpaceID:                       types.StringValue(lifecycle.SpaceID),
			Name:                          types.StringValue(lifecycle.Name),
			Description:                   types.StringValue(lifecycle.Description),
			Phase:                         flattenPhases(lifecycle.Phases),
			ReleaseRetentionWithStrategy:  flattenRetentionWithStrategy(lifecycle.ReleaseRetentionPolicy),
			TentacleRetentionWithStrategy: flattenRetentionWithStrategy(lifecycle.TentacleRetentionPolicy),
		}
	}

	flattenedLifecycle.ID = types.StringValue(lifecycle.GetID())

	return flattenedLifecycle
}

func flattenPhases(goPhases []*lifecycles.Phase) types.List {
	var attributeTypes = getPhaseWithStrategyAttrTypes()
	if goPhases == nil {
		return types.ListNull(types.ObjectType{AttrTypes: attributeTypes})
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
			"release_retention_policy":              types.ListNull(types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()}),
			"tentacle_retention_policy":             types.ListNull(types.ObjectType{AttrTypes: DeprecatedGetRetentionAttTypes()}),
			"release_retention_with_strategy":       util.Ternary(goPhase.ReleaseRetentionPolicy != nil, flattenRetentionWithStrategy(goPhase.ReleaseRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})),
			"tentacle_retention_with_strategy":      util.Ternary(goPhase.TentacleRetentionPolicy != nil, flattenRetentionWithStrategy(goPhase.TentacleRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})),
		}
		phasesList = append(phasesList, types.ObjectValueMust(attributeTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: attributeTypes}, phasesList)
}

func flattenRetentionWithStrategy(goRetention *core.RetentionPeriod) types.List {
	var attributeTypes = getRetentionWithStrategyAttrTypes()
	if goRetention == nil {
		return ListNullRetentionWithStrategy
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: attributeTypes},
		[]attr.Value{
			types.ObjectValueMust(
				attributeTypes,
				map[string]attr.Value{
					"strategy":         types.StringValue(goRetention.Strategy),
					"unit":             types.StringValue(goRetention.Unit),
					"quantity_to_keep": types.Int64Value(int64(goRetention.QuantityToKeep)),
				},
			),
		},
	)
}

func expandLifecycle(lifecycleUserInput *lifecycleTypeResourceModel, onlyDeprecatedRetentionIsUsed bool) *lifecycles.Lifecycle {
	if lifecycleUserInput == nil {
		return nil
	}

	lifecyclePlan := lifecycles.NewLifecycle(lifecycleUserInput.Name.ValueString())
	lifecyclePlan.Description = lifecycleUserInput.Description.ValueString()
	lifecyclePlan.SpaceID = lifecycleUserInput.SpaceID.ValueString()
	if !lifecycleUserInput.ID.IsNull() && lifecycleUserInput.ID.ValueString() != "" {
		lifecyclePlan.ID = lifecycleUserInput.ID.ValueString()
	}

	if onlyDeprecatedRetentionIsUsed {
		lifecyclePlan.Phases = DeprecatedExpandPhases(lifecycleUserInput.Phase)
		lifecyclePlan.ReleaseRetentionPolicy = DeprecatedExpandRetention(lifecycleUserInput.DeprecatedReleaseRetention)
		lifecyclePlan.TentacleRetentionPolicy = DeprecatedExpandRetention(lifecycleUserInput.DeprecatedTentacleRetention)
	} else {
		lifecyclePlan.Phases = expandPhasesWithStrategy(lifecycleUserInput.Phase)
		lifecyclePlan.ReleaseRetentionPolicy = expandRetentionWithStrategy(lifecycleUserInput.ReleaseRetentionWithStrategy)
		lifecyclePlan.TentacleRetentionPolicy = expandRetentionWithStrategy(lifecycleUserInput.TentacleRetentionWithStrategy)
	}
	return lifecyclePlan
}

func expandPhasesWithStrategy(phasesUserInput types.List) []*lifecycles.Phase {
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

		if v, ok := phaseAttrsUserInput["release_retention_with_strategy"].(types.List); ok && !v.IsNull() {
			phasePlan.ReleaseRetentionPolicy = expandRetentionWithStrategy(v)
		}
		if v, ok := phaseAttrsUserInput["tentacle_retention_with_strategy"].(types.List); ok && !v.IsNull() {
			phasePlan.TentacleRetentionPolicy = expandRetentionWithStrategy(v)
		}
		allPhasesPlan = append(allPhasesPlan, phasePlan)
	}

	return allPhasesPlan
}

func expandRetentionWithStrategy(retentionUserInput types.List) *core.RetentionPeriod {
	if retentionUserInput.IsNull() || retentionUserInput.IsUnknown() || len(retentionUserInput.Elements()) == 0 {
		return nil
	}

	retentionAttrsInput := retentionUserInput.Elements()[0].(types.Object).Attributes()

	var strategy string
	if s, ok := retentionAttrsInput["strategy"].(types.String); ok && !s.IsNull() {
		strategy = s.ValueString()
	}

	if strategy == core.RetentionStrategyCount {
		var unit string
		if u, ok := retentionAttrsInput["unit"].(types.String); ok && !u.IsNull() {
			unit = u.ValueString()
		}

		var quantityToKeep int32
		if qty, ok := retentionAttrsInput["quantity_to_keep"].(types.Int64); ok && !qty.IsNull() {
			quantityToKeep = int32(qty.ValueInt64())
		}
		return core.CountBasedRetentionPeriod(quantityToKeep, unit)
	}

	if strategy == core.RetentionStrategyForever {
		return core.KeepForeverRetentionPeriod()
	}
	return core.SpaceDefaultRetentionPeriod()

}

func getRetentionWithStrategyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"strategy":         types.StringType,
		"quantity_to_keep": types.Int64Type,
		"unit":             types.StringType,
	}
}

func getPhaseWithStrategyAttrTypes() map[string]attr.Type {
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

func GetPhaseAttributeTypes() map[string]attr.Type {
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
