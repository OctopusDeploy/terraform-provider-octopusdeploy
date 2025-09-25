package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	onlyDeprecatedRetentionIsSupported bool
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

func (r *lifecycleTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.LifecycleSchema{}.GetResourceSchema()
}

func (r *lifecycleTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
	if r.Config != nil {
		r.onlyDeprecatedRetentionIsSupported = !r.Config.IsVersionSameOrGreaterThan("2025.3") // this always returns true if running on the local
	}

}

func (r *lifecycleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	validateRetentionBlocksUsed(data, &resp.Diagnostics, r.onlyDeprecatedRetentionIsSupported)
	onlyDeprecatedRetentionIsUsed := isOnlyDeprecatedRetentionUsed(data, r.onlyDeprecatedRetentionIsSupported)
	if resp.Diagnostics.HasError() {
		return
	}

	hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := setDefaultRetentionDeprecated(data, onlyDeprecatedRetentionIsUsed)
	hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy := setDefaultRetentionWithStrategy(data, onlyDeprecatedRetentionIsUsed)
	tflog.Debug(ctx, fmt.Sprintf("before expand rose'%v'", data.Phase))
	newLifecycle := expandLifecycle(data, onlyDeprecatedRetentionIsUsed)
	tflog.Debug(ctx, fmt.Sprintf("after expand rose'%v'", data.Phase))
	lifecycle, err := lifecycles.Add(r.Config.Client, newLifecycle)

	if err != nil {
		resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
		return
	}

	for _, goPhase := range lifecycle.Phases {
		tflog.Debug(ctx, fmt.Sprintf("ran add rose retention'%v'tentacle `%v`", goPhase.ReleaseRetentionPolicy, goPhase.TentacleRetentionPolicy))
	}

	data = flattenLifecycleResource(lifecycle, onlyDeprecatedRetentionIsUsed)
	tflog.Debug(ctx, fmt.Sprintf("after flatten rose retention'%v'", data.Phase))
	removeDefaultRetentionFromUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention)
	removeDefaultRetentionWithStrategyFromUnsetBlocks(data, hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy)
	tflog.Debug(ctx, fmt.Sprintf("after flatten  and removal rose'%v'", data.Phase))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	onlyDeprecatedRetentionIsUsed := isOnlyDeprecatedRetentionUsed(data, r.onlyDeprecatedRetentionIsSupported)
	hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := setDefaultRetentionDeprecated(data, onlyDeprecatedRetentionIsUsed)
	hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy := setDefaultRetentionWithStrategy(data, onlyDeprecatedRetentionIsUsed)

	lifecycle, err := lifecycles.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "lifecycle"); err != nil {
			resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
		}
		return
	}

	data = flattenLifecycleResource(lifecycle, onlyDeprecatedRetentionIsUsed)

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

	validateRetentionBlocksUsed(data, &resp.Diagnostics, r.onlyDeprecatedRetentionIsSupported)
	onlyDeprecatedRetentionIsUsed := isOnlyDeprecatedRetentionUsed(data, r.onlyDeprecatedRetentionIsSupported)
	hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := setDefaultRetentionDeprecated(data, onlyDeprecatedRetentionIsUsed)
	hasUserDefinedReleaseRetentionWithStrategy, hasUserDefinedTentacleRetentionWithStrategy := setDefaultRetentionWithStrategy(data, onlyDeprecatedRetentionIsUsed)

	tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", data.ID.ValueString()))

	lifecycle := expandLifecycle(data, onlyDeprecatedRetentionIsUsed)
	lifecycle.ID = state.ID.ValueString()

	updatedLifecycle, err := lifecycles.Update(r.Config.Client, lifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
		return
	}

	data = flattenLifecycleResource(updatedLifecycle, onlyDeprecatedRetentionIsUsed)

	removeDefaultRetentionFromUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention)
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

func setDefaultRetentionDeprecated(data *lifecycleTypeResourceModel, onlyDeprecatedRetentionIsUsed bool) (bool, bool) {

	hasUserDefinedReleaseRetention := attributeIsUsed(data.ReleaseRetention)
	hasUserDefinedTentacleRetention := attributeIsUsed(data.TentacleRetention)

	var initialRetentionSetting types.List
	if onlyDeprecatedRetentionIsUsed {
		initialRetentionSetting = deprecatedFlattenRetention(core.NewRetentionPeriod(30, "Days", false))
	} else {
		initialRetentionSetting = types.ListNull(types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()})
	}

	if !hasUserDefinedReleaseRetention {
		data.ReleaseRetention = initialRetentionSetting
	}
	if !hasUserDefinedTentacleRetention {
		data.TentacleRetention = initialRetentionSetting
	}

	return hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention
}
func setDefaultRetentionWithStrategy(data *lifecycleTypeResourceModel, onlyDeprecatedRetentionIsUsed bool) (bool, bool) {
	isReleaseRetentionWithStrategyDefined := attributeIsUsed(data.ReleaseRetentionWithStrategy)
	isTentacleRetentionWithStrategyDefined := attributeIsUsed(data.TentacleRetentionWithStrategy)

	var initialRetentionSetting types.List
	if onlyDeprecatedRetentionIsUsed {
		initialRetentionSetting = types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})
	} else {
		initialRetentionSetting = flattenRetentionWithStrategy(core.SpaceDefaultRetentionPeriod())
	}

	if !isReleaseRetentionWithStrategyDefined {
		data.ReleaseRetentionWithStrategy = initialRetentionSetting
	}
	if !isTentacleRetentionWithStrategyDefined {
		data.TentacleRetentionWithStrategy = initialRetentionSetting
	}

	return isReleaseRetentionWithStrategyDefined, isTentacleRetentionWithStrategyDefined
}
func removeDefaultRetentionFromUnsetBlocks(data *lifecycleTypeResourceModel, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention bool) {
	// Remove retention policies from data before setting state, but only if they were not included in the plan
	if !hasUserDefinedReleaseRetention {
		data.ReleaseRetention = types.ListNull(types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()})
	}
	if !hasUserDefinedTentacleRetention {
		data.TentacleRetention = types.ListNull(types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()})
	}
}
func removeDefaultRetentionWithStrategyFromUnsetBlocks(data *lifecycleTypeResourceModel, hasUserDefinedReleaseRetentionWithStrategy bool, hasUserDefinedTentacleRetentionWithStrategy bool) {
	// Remove retention policies from data before setting state, but only if they were not included in the plan
	if !hasUserDefinedReleaseRetentionWithStrategy {
		data.ReleaseRetentionWithStrategy = types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})
	}
	if !hasUserDefinedTentacleRetentionWithStrategy {
		data.TentacleRetentionWithStrategy = types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})
	}
}

func isRetentionWithStrategyInPlan(data *lifecycleTypeResourceModel) bool {
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

func isDeprecatedRetentionInPlan(data *lifecycleTypeResourceModel) bool {
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

func isOnlyDeprecatedRetentionUsed(data *lifecycleTypeResourceModel, onlyDeprecatedRetentionIsSupported bool) bool {
	deprecatedRetentionIsInPlan := isDeprecatedRetentionInPlan(data)
	if deprecatedRetentionIsInPlan {
		return true
	}
	if onlyDeprecatedRetentionIsSupported {
		return true
	}
	return false
}

func validateRetentionBlocksUsed(data *lifecycleTypeResourceModel, diag *diag.Diagnostics, onlyDeprecatedRetentionIsSupported bool) {
	retentionWithStrategyIsInPlan := isRetentionWithStrategyInPlan(data)
	deprecatedRetentionIsInPlan := isDeprecatedRetentionInPlan(data)
	if retentionWithStrategyIsInPlan && deprecatedRetentionIsInPlan {
		diag.AddError("Retention blocks conflict", "Both release_retention_with_strategy and release_retention_policy are used. Only one can be used at a time.")
	}
	if retentionWithStrategyIsInPlan && onlyDeprecatedRetentionIsSupported {
		diag.AddError("retention with strategy is not supported on this Octopus Server version. Please upgrade to Octopus Server 2025.3 or later.", "")
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
			SpaceID:           types.StringValue(lifecycle.SpaceID),
			Name:              types.StringValue(lifecycle.Name),
			Description:       types.StringValue(lifecycle.Description),
			Phase:             deprecatedFlattenPhases(lifecycle.Phases),
			ReleaseRetention:  deprecatedFlattenRetention(lifecycle.ReleaseRetentionPolicy),
			TentacleRetention: deprecatedFlattenRetention(lifecycle.TentacleRetentionPolicy),
		}
	} else {
		flattenedLifecycle = &lifecycleTypeResourceModel{
			SpaceID:                       types.StringValue(lifecycle.SpaceID),
			Name:                          types.StringValue(lifecycle.Name),
			Description:                   types.StringValue(lifecycle.Description),
			Phase:                         flattenPhasesWithStrategy(lifecycle.Phases),
			ReleaseRetentionWithStrategy:  flattenRetentionWithStrategy(lifecycle.ReleaseRetentionPolicy),
			TentacleRetentionWithStrategy: flattenRetentionWithStrategy(lifecycle.TentacleRetentionPolicy),
		}
	}

	flattenedLifecycle.ID = types.StringValue(lifecycle.GetID())

	return flattenedLifecycle
}

func deprecatedFlattenPhases(goPhases []*lifecycles.Phase) types.List {
	var deprecatedAttributeTypes = deprecatedGetAttributeTypes()
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
			"release_retention_policy":              util.Ternary(goPhase.ReleaseRetentionPolicy != nil, deprecatedFlattenRetention(goPhase.ReleaseRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()})),
			"tentacle_retention_policy":             util.Ternary(goPhase.TentacleRetentionPolicy != nil, deprecatedFlattenRetention(goPhase.TentacleRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()})),
			"release_retention_with_strategy":       types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}),
			"tentacle_retention_with_strategy":      types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}),
		}
		phasesList = append(phasesList, types.ObjectValueMust(deprecatedAttributeTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: deprecatedAttributeTypes}, phasesList)
}

func flattenPhasesWithStrategy(goPhases []*lifecycles.Phase) types.List {
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
			"release_retention_policy":              types.ListNull(types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()}),
			"tentacle_retention_policy":             types.ListNull(types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()}),
			"release_retention_with_strategy":       util.Ternary(goPhase.ReleaseRetentionPolicy != nil, flattenRetentionWithStrategy(goPhase.ReleaseRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})),
			"tentacle_retention_with_strategy":      util.Ternary(goPhase.TentacleRetentionPolicy != nil, flattenRetentionWithStrategy(goPhase.TentacleRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()})),
		}
		phasesList = append(phasesList, types.ObjectValueMust(attributeTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: attributeTypes}, phasesList)
}

func deprecatedFlattenRetention(goRetention *core.RetentionPeriod) types.List {
	var deprecatedAttributeTypes = deprecatedGetRetentionAttTypes()
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

func flattenRetentionWithStrategy(goRetention *core.RetentionPeriod) types.List {
	var attributeTypes = getRetentionWithStrategyAttrTypes()
	if goRetention == nil {
		return types.ListNull(types.ObjectType{AttrTypes: attributeTypes})
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
		lifecyclePlan.Phases = deprecatedExpandPhases(lifecycleUserInput.Phase)
		lifecyclePlan.ReleaseRetentionPolicy = deprecatedExpandRetention(lifecycleUserInput.ReleaseRetention)
		lifecyclePlan.TentacleRetentionPolicy = deprecatedExpandRetention(lifecycleUserInput.TentacleRetention)
	} else {
		lifecyclePlan.Phases = expandPhasesWithStrategy(lifecycleUserInput.Phase)
		lifecyclePlan.ReleaseRetentionPolicy = expandRetentionWithStrategy(lifecycleUserInput.ReleaseRetentionWithStrategy)
		lifecyclePlan.TentacleRetentionPolicy = expandRetentionWithStrategy(lifecycleUserInput.TentacleRetentionWithStrategy)
	}
	return lifecyclePlan
}

func deprecatedExpandPhases(phasesUserInput types.List) []*lifecycles.Phase {
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
			phasePlan.ReleaseRetentionPolicy = deprecatedExpandRetention(v)
		}
		if v, ok := phaseAttrsUserInput["tentacle_retention_policy"].(types.List); ok && !v.IsNull() {
			phasePlan.TentacleRetentionPolicy = deprecatedExpandRetention(v)
		}
		allPhasesPlan = append(allPhasesPlan, phasePlan)
	}

	return allPhasesPlan
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

func deprecatedExpandRetention(v types.List) *core.RetentionPeriod {
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

func deprecatedGetRetentionAttTypes() map[string]attr.Type {
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

func deprecatedGetAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"is_priority_phase":                     types.BoolType,
		"release_retention_policy":              types.ListType{ElemType: types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()}},
		"tentacle_retention_policy":             types.ListType{ElemType: types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()}},
		"release_retention_with_strategy":       types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}},
		"tentacle_retention_with_strategy":      types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}},
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
		"release_retention_policy":              types.ListType{ElemType: types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()}},
		"tentacle_retention_policy":             types.ListType{ElemType: types.ObjectType{AttrTypes: deprecatedGetRetentionAttTypes()}},
		"release_retention_with_strategy":       types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}},
		"tentacle_retention_with_strategy":      types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionWithStrategyAttrTypes()}},
	}
}
