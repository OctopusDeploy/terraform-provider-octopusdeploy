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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type lifecycleTypeResource struct {
	*Config
	newRetentionNotSupported    bool
	allowDeprecatedRetention bool
}

var _ resource.Resource = &lifecycleTypeResource{}
var _ resource.ResourceWithImportState = &lifecycleTypeResource{}

type lifecycleTypeResourceModel struct {
	SpaceID           types.String `tfsdk:"space_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Phase             types.List   `tfsdk:"phase"`
	ReleaseRetention  types.List   `tfsdk:"release_retention_with_strategy"`
	TentacleRetention types.List   `tfsdk:"tentacle_retention_with_strategy"`

	schemas.ResourceModel
}

type lifecycleTypeResourceModelDEPRECATED struct {
	SpaceID              types.String `tfsdk:"space_id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Phase                types.List   `tfsdk:"phase"`
	ReleaseRetentionWithoutStrategy  types.List   `tfsdk:"release_retention_policy"`
	TentacleRetentionWithoutStrategy types.List   `tfsdk:"tentacle_retention_policy"`
	ReleaseRetention     types.List   `tfsdk:"release_retention_with_strategy"`
	TentacleRetention    types.List   `tfsdk:"tentacle_retention_with_strategy"`

	schemas.ResourceModel
}

func NewLifecycleResource() resource.Resource {
	allowDeprecatedRetention := schemas.AllowDeprecatedRetention()
	return &lifecycleTypeResource{allowDeprecatedRetention: allowDeprecatedRetention}
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

func (r *lifecycleTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
	if r.Config != nil {
		if !r.Config.IsVersionSameOrGreaterThan("2025.3") {
			r.newRetentionNotSupported = true
		}
	}
}

func (r *lifecycleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.allowDeprecatedRetention {
		var stateData *lifecycleTypeResourceModelDEPRECATED
		resp.Diagnostics.Append(req.Config.Get(ctx, &stateData)...)
		if resp.Diagnostics.HasError() {
			return
		}

		validateRetentionBlocksInConfigDEPRECATED(stateData, &resp.Diagnostics, r.newRetentionNotSupported)
		onlyRetentionWithoutStrategyWillBeUsed := isRetentionWithoutStrategyToBeUsedForLifecycleDEPRECATED(stateData)
		isNewReleaseRetentionSet, isNewTentacleRetentionSet, isReleaseRetentionWithoutStrategySet, isTentacleRetentionWithoutStrategySet := determineWhichRetentionBlocksAreInConfigDEPRECATED(stateData)
		initialRetentionSettingForNewBlock, initialRetentionSettingForRetentionWithoutStrategyBlock := setInitialRetentionDEPRECATED(stateData, onlyRetentionWithoutStrategyWillBeUsed)

		lifecycleSentToGo := expandLifecycleDEPRECATED(stateData, onlyRetentionWithoutStrategyWillBeUsed)
		lifecycleFromGo, err := lifecycles.Add(r.Config.Client, lifecycleSentToGo)
		if err != nil {
			resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
			return
		}

		handleUnitCasing(lifecycleFromGo, lifecycleSentToGo)
		stateData = flattenResourceLifecycleDEPRECATED(lifecycleFromGo, onlyRetentionWithoutStrategyWillBeUsed)
		removeInitialRetentionDEPRECATED(stateData, isNewReleaseRetentionSet, isNewTentacleRetentionSet, isReleaseRetentionWithoutStrategySet, isTentacleRetentionWithoutStrategySet, initialRetentionSettingForNewBlock, initialRetentionSettingForRetentionWithoutStrategyBlock)

		resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)

	} else {
		var stateData *lifecycleTypeResourceModel
		resp.Diagnostics.Append(req.Config.Get(ctx, &stateData)...)
		if resp.Diagnostics.HasError() {
			return
		}

		initialRetentionSetting := flattenResourceRetention(core.SpaceDefaultRetentionPeriod())
		isReleaseRetentionDefined, isTentacleRetentionDefined := setInitialRetention(stateData, initialRetentionSetting)

		lifecycle := expandLifecycle(stateData)
		newLifecycle, err := lifecycles.Add(r.Config.Client, lifecycle)
		if err != nil {
			resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
			return
		}
		stateData = flattenResourceLifecycle(newLifecycle)

		removeDefaultRetentionFromUnsetBlocks(stateData, isReleaseRetentionDefined, isTentacleRetentionDefined, initialRetentionSetting)
		resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	}
}
func (r *lifecycleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.allowDeprecatedRetention {
		var stateData *lifecycleTypeResourceModelDEPRECATED
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
		if resp.Diagnostics.HasError() {
			return
		}

		validateRetentionBlocksInConfigDEPRECATED(stateData, &resp.Diagnostics, r.newRetentionNotSupported)
		onlyRetentionWithoutStrategyWillBeUsed := isRetentionWithoutStrategyToBeUsedForLifecycleDEPRECATED(stateData)
		isNewReleaseRetentionSet, isNewTentacleRetentionSet, isReleaseRetentionWithoutStrategySet, isTentacleRetentionWithoutStrategySet := determineWhichRetentionBlocksAreInConfigDEPRECATED(stateData)
		initialRetentionSettingForNewBlock, initialRetentionSettingForWithoutStrategyBlock := setInitialRetentionDEPRECATED(stateData, onlyRetentionWithoutStrategyWillBeUsed)

		lifecycleFromGo, err := lifecycles.GetByID(r.Config.Client, stateData.SpaceID.ValueString(), stateData.ID.ValueString())
		if err != nil {
			if err := errors.ProcessApiErrorV2(ctx, resp, stateData, err, "lifecycle"); err != nil {
				resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
			}
			return
		}
	
		handleUnitCasing(lifecycleFromGo, expandLifecycleDEPRECATED(stateData, onlyRetentionWithoutStrategyWillBeUsed))
		stateData = flattenResourceLifecycleDEPRECATED(lifecycleFromGo, onlyRetentionWithoutStrategyWillBeUsed)

		removeInitialRetentionDEPRECATED(stateData, isNewReleaseRetentionSet, isNewTentacleRetentionSet, isReleaseRetentionWithoutStrategySet, isTentacleRetentionWithoutStrategySet, initialRetentionSettingForNewBlock, initialRetentionSettingForWithoutStrategyBlock)
		resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)

	} else {
		var stateData *lifecycleTypeResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
		if resp.Diagnostics.HasError() {
			return
		}

		initialRetentionSetting := flattenResourceRetention(core.SpaceDefaultRetentionPeriod())
		isReleaseRetentionDefined, isTentacleRetentionDefined := setInitialRetention(stateData, initialRetentionSetting)

		lifecycleFromGo, err := lifecycles.GetByID(r.Config.Client, stateData.SpaceID.ValueString(), stateData.ID.ValueString())
		if err != nil {
			if err := errors.ProcessApiErrorV2(ctx, resp, stateData, err, "lifecycle"); err != nil {
				resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
			}
			return
		}

		handleUnitCasing(lifecycleFromGo, expandLifecycle(stateData))
		stateData = flattenResourceLifecycle(lifecycleFromGo)

		removeDefaultRetentionFromUnsetBlocks(stateData, isReleaseRetentionDefined, isTentacleRetentionDefined, initialRetentionSetting)
		resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	}
}

func (r *lifecycleTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.allowDeprecatedRetention {
		var stateData, state *lifecycleTypeResourceModelDEPRECATED
		resp.Diagnostics.Append(req.Plan.Get(ctx, &stateData)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}

		validateRetentionBlocksInConfigDEPRECATED(stateData, &resp.Diagnostics, r.newRetentionNotSupported)
		onlyRetentionWithoutStrategyWillBeUsed := isRetentionWithoutStrategyToBeUsedForLifecycleDEPRECATED(stateData)
		isNewReleaseRetentionSet, isNewTentacleRetentionSet, isReleaseRetentionWithoutStrategySet, isTentacleRetentionWithoutStrategySet := determineWhichRetentionBlocksAreInConfigDEPRECATED(stateData)
		initialRetentionSettingForNewBlock, initialRetentionSettingForWithoutStrategyBlock := setInitialRetentionDEPRECATED(stateData, onlyRetentionWithoutStrategyWillBeUsed)

		tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", stateData.ID.ValueString()))
		lifecycleSentToGo := expandLifecycleDEPRECATED(stateData, onlyRetentionWithoutStrategyWillBeUsed)
		lifecycleSentToGo.ID = state.ID.ValueString()
		lifecycleFromGo, err := lifecycles.Update(r.Config.Client, lifecycleSentToGo)
		if err != nil {
			resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
			return
		}
	
		handleUnitCasing(lifecycleFromGo, lifecycleSentToGo)
		stateData = flattenResourceLifecycleDEPRECATED(lifecycleFromGo, onlyRetentionWithoutStrategyWillBeUsed)

		removeInitialRetentionDEPRECATED(stateData, isNewReleaseRetentionSet, isNewTentacleRetentionSet, isReleaseRetentionWithoutStrategySet, isTentacleRetentionWithoutStrategySet, initialRetentionSettingForNewBlock, initialRetentionSettingForWithoutStrategyBlock)
		resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)

	} else {
		var stateData, state *lifecycleTypeResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &stateData)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}

		initialRetentionSetting := flattenResourceRetention(core.SpaceDefaultRetentionPeriod())
		isReleaseRetentionDefined, isTentacleRetentionDefined := setInitialRetention(stateData, initialRetentionSetting)

		tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", stateData.ID.ValueString()))
		lifecycle := expandLifecycle(stateData)
		lifecycle.ID = state.ID.ValueString()
		updatedLifecycle, err := lifecycles.Update(r.Config.Client, lifecycle)
		if err != nil {
			resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
			return
		}

		handleUnitCasing(updatedLifecycle, expandLifecycle(stateData))
		stateData = flattenResourceLifecycle(updatedLifecycle)

		removeDefaultRetentionFromUnsetBlocks(stateData, isReleaseRetentionDefined, isTentacleRetentionDefined, initialRetentionSetting)
		resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	}
}

func (r *lifecycleTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.allowDeprecatedRetention {
		var stateData lifecycleTypeResourceModelDEPRECATED
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if err := lifecycles.DeleteByID(r.Config.Client, stateData.SpaceID.ValueString(), stateData.ID.ValueString()); err != nil {
			resp.Diagnostics.AddError("unable to delete lifecycle", err.Error())
			return
		}

	} else {
		var stateData lifecycleTypeResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if err := lifecycles.DeleteByID(r.Config.Client, stateData.SpaceID.ValueString(), stateData.ID.ValueString()); err != nil {
			resp.Diagnostics.AddError("unable to delete lifecycle", err.Error())
			return
		}
	}
}

func handleUnitCasing(lifecycleFromGo *lifecycles.Lifecycle, lifecycleInState *lifecycles.Lifecycle) {
	// Set state to the casing provided in the desired state, as the Api will always return capitalised units
	lifecycleFromGo.ReleaseRetentionPolicy = updateUnitsFromGoToMatchStateCasing(lifecycleFromGo.ReleaseRetentionPolicy, lifecycleInState.ReleaseRetentionPolicy.Unit)
	lifecycleFromGo.TentacleRetentionPolicy = updateUnitsFromGoToMatchStateCasing(lifecycleFromGo.TentacleRetentionPolicy, lifecycleInState.TentacleRetentionPolicy.Unit)

	if len(lifecycleInState.Phases) == 0 {
		return
	}

	for i, phaseFromGo := range lifecycleFromGo.Phases {
		if phaseFromGo.ReleaseRetentionPolicy != nil && phaseFromGo.ReleaseRetentionPolicy.Unit != "" {
			phaseFromGo.ReleaseRetentionPolicy = updateUnitsFromGoToMatchStateCasing(phaseFromGo.ReleaseRetentionPolicy, lifecycleInState.Phases[i].ReleaseRetentionPolicy.Unit)
		}
		if phaseFromGo.TentacleRetentionPolicy != nil && phaseFromGo.TentacleRetentionPolicy.Unit != "" {
			phaseFromGo.TentacleRetentionPolicy = updateUnitsFromGoToMatchStateCasing(phaseFromGo.TentacleRetentionPolicy, lifecycleInState.Phases[i].TentacleRetentionPolicy.Unit)
		}
	}
}

func updateUnitsFromGoToMatchStateCasing(retentionPeriodFromGo *core.RetentionPeriod, unitInState string) *core.RetentionPeriod {
	if strings.EqualFold(retentionPeriodFromGo.Unit, unitInState) {
		replacementForRetentionFromGo := core.RetentionPeriod{
			QuantityToKeep:    retentionPeriodFromGo.QuantityToKeep,
			ShouldKeepForever: retentionPeriodFromGo.ShouldKeepForever,
			Unit:              unitInState,
			Strategy:          retentionPeriodFromGo.Strategy,
		}

		return &replacementForRetentionFromGo
	}

	return retentionPeriodFromGo
}

func removeInitialRetentionDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED, isNewReleaseRetentionSet bool, isTentacleRetentionDefined bool, isReleaseRetentionWithoutStrategySet bool, isTentacleRetentionWithoutStrategySet bool, initialRetentionSettingForNewBlock types.List, initialRetentionSettingForWithoutStrategyBlock types.List) {
	if !isNewReleaseRetentionSet && (data.ReleaseRetention.Equal(initialRetentionSettingForNewBlock) || data.ReleaseRetention.IsNull()) {
		data.ReleaseRetention = ListNullRetention
	}
	if !isTentacleRetentionDefined && (data.TentacleRetention.Equal(initialRetentionSettingForNewBlock) || data.TentacleRetention.IsNull()) {
		data.TentacleRetention = ListNullRetention
	}

	if !isReleaseRetentionWithoutStrategySet && (data.ReleaseRetentionWithoutStrategy.Equal(initialRetentionSettingForWithoutStrategyBlock) || data.ReleaseRetentionWithoutStrategy.IsNull()) {
		data.ReleaseRetentionWithoutStrategy = ListNullRetentionWithoutStrategyDEPRECATED
	}

	if !isTentacleRetentionWithoutStrategySet && (data.TentacleRetentionWithoutStrategy.Equal(initialRetentionSettingForWithoutStrategyBlock) || data.TentacleRetentionWithoutStrategy.IsNull()) {
		data.TentacleRetentionWithoutStrategy = ListNullRetentionWithoutStrategyDEPRECATED
	}
}

func setInitialRetentionWithoutStrategyBlockDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED, initialRetentionSettingForWithoutStrategyBlock types.List) (bool, bool) {
	hasUserDefinedReleaseRetentionWithoutStrategy := attributeIsUsed(data.ReleaseRetentionWithoutStrategy)
	hasUserDefinedTentacleRetentionWithoutStrategy := attributeIsUsed(data.TentacleRetentionWithoutStrategy)
	if !hasUserDefinedReleaseRetentionWithoutStrategy {
		data.ReleaseRetentionWithoutStrategy = initialRetentionSettingForWithoutStrategyBlock
	}
	if !hasUserDefinedTentacleRetentionWithoutStrategy {
		data.TentacleRetentionWithoutStrategy = initialRetentionSettingForWithoutStrategyBlock
	}

	return hasUserDefinedReleaseRetentionWithoutStrategy, hasUserDefinedTentacleRetentionWithoutStrategy
}
func setInitialRetentionNewBlockDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED, initialRetentionSettingForNewBlock types.List) (bool, bool) {
	hasUserDefinedReleaseRetention := attributeIsUsed(data.ReleaseRetention)
	hasUserDefinedTentacleRetention := attributeIsUsed(data.TentacleRetention)
	if !hasUserDefinedReleaseRetention {
		data.ReleaseRetention = initialRetentionSettingForNewBlock
	}
	if !hasUserDefinedTentacleRetention {
		data.TentacleRetention = initialRetentionSettingForNewBlock
	}
	return hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention
}
func setInitialRetentionDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED, onlyRetentionWithoutStrategyWillBeUsed bool) (types.List, types.List) {
	var initialRetentionSettingForWithoutStrategyBlock types.List
	var initialRetentionSettingForNewBlock types.List

	if onlyRetentionWithoutStrategyWillBeUsed {
		initialRetentionSettingForWithoutStrategyBlock = flattenResourceRetentionWithoutStrategyDEPRECATED(core.CountBasedRetentionPeriod(30, "Days"))
		initialRetentionSettingForNewBlock = ListNullRetention
		setInitialRetentionWithoutStrategyBlockDEPRECATED(data, initialRetentionSettingForWithoutStrategyBlock)
	} else {
		initialRetentionSettingForWithoutStrategyBlock = ListNullRetentionWithoutStrategyDEPRECATED
		initialRetentionSettingForNewBlock = flattenResourceRetention(core.SpaceDefaultRetentionPeriod())
		setInitialRetentionNewBlockDEPRECATED(data, initialRetentionSettingForNewBlock)
	}
	return initialRetentionSettingForNewBlock, initialRetentionSettingForWithoutStrategyBlock
}

func setInitialRetention(data *lifecycleTypeResourceModel, initialRetentionSetting types.List) (bool, bool) {
	isReleaseRetentionDefined := attributeIsUsed(data.ReleaseRetention)
	isTentacleRetentionDefined := attributeIsUsed(data.TentacleRetention)

	if !isReleaseRetentionDefined {
		data.ReleaseRetention = initialRetentionSetting
	}
	if !isTentacleRetentionDefined {
		data.TentacleRetention = initialRetentionSetting
	}

	return isReleaseRetentionDefined, isTentacleRetentionDefined
}
func removeDefaultRetentionFromUnsetBlocks(data *lifecycleTypeResourceModel, isReleaseRetentionDefined bool, isTentacleRetentionDefined bool, initialRetentionSetting types.List) {
	if !isReleaseRetentionDefined && data.ReleaseRetention.Equal(initialRetentionSetting) {
		data.ReleaseRetention = ListNullRetention
	}
	if !isTentacleRetentionDefined && data.TentacleRetention.Equal(initialRetentionSetting) {
		data.TentacleRetention = ListNullRetention
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

func flattenResourceLifecycle(lifecycle *lifecycles.Lifecycle) *lifecycleTypeResourceModel {
	var flattenedLifecycle *lifecycleTypeResourceModel
	flattenedLifecycle = &lifecycleTypeResourceModel{
		SpaceID:           types.StringValue(lifecycle.SpaceID),
		Name:              types.StringValue(lifecycle.Name),
		Description:       types.StringValue(lifecycle.Description),
		Phase:             flattenResourcePhases(lifecycle.Phases),
		ReleaseRetention:  flattenResourceRetention(lifecycle.ReleaseRetentionPolicy),
		TentacleRetention: flattenResourceRetention(lifecycle.TentacleRetentionPolicy),
	}
	flattenedLifecycle.ID = types.StringValue(lifecycle.GetID())

	return flattenedLifecycle
}
func flattenResourceLifecycleDEPRECATED(lifecycle *lifecycles.Lifecycle, onlyRetentionWithoutStrategyWillBeUsed bool) *lifecycleTypeResourceModelDEPRECATED {
	var flattenedLifecycle *lifecycleTypeResourceModelDEPRECATED
	flattenedLifecycle = &lifecycleTypeResourceModelDEPRECATED{
		SpaceID:     types.StringValue(lifecycle.SpaceID),
		Name:        types.StringValue(lifecycle.Name),
		Description: types.StringValue(lifecycle.Description),
		Phase:       flattenResourcePhasesDEPRECATED(lifecycle.Phases, onlyRetentionWithoutStrategyWillBeUsed),
	}
	if onlyRetentionWithoutStrategyWillBeUsed {
		flattenedLifecycle.ReleaseRetentionWithoutStrategy = flattenResourceRetentionWithoutStrategyDEPRECATED(lifecycle.ReleaseRetentionPolicy)
		flattenedLifecycle.TentacleRetentionWithoutStrategy = flattenResourceRetentionWithoutStrategyDEPRECATED(lifecycle.TentacleRetentionPolicy)
	} else {
		flattenedLifecycle.ReleaseRetention = flattenResourceRetention(lifecycle.ReleaseRetentionPolicy)
		flattenedLifecycle.TentacleRetention = flattenResourceRetention(lifecycle.TentacleRetentionPolicy)
	}

	flattenedLifecycle.ID = types.StringValue(lifecycle.GetID())
	return flattenedLifecycle
}
func flattenResourcePhases(goPhases []*lifecycles.Phase) types.List {
	var phaseAttrTypes = getResourcePhaseAttrTypes()
	if goPhases == nil {
		return types.ListNull(types.ObjectType{AttrTypes: phaseAttrTypes})
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
			"release_retention_with_strategy":       util.Ternary(goPhase.ReleaseRetentionPolicy != nil, flattenResourceRetention(goPhase.ReleaseRetentionPolicy), ListNullRetention),
			"tentacle_retention_with_strategy":      util.Ternary(goPhase.TentacleRetentionPolicy != nil, flattenResourceRetention(goPhase.TentacleRetentionPolicy), ListNullRetention),
		}
		phasesList = append(phasesList, types.ObjectValueMust(phaseAttrTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: phaseAttrTypes}, phasesList)
}

func flattenResourcePhasesDEPRECATED(phasesFromGo []*lifecycles.Phase, onlyRetentionWithoutStrategyWillBeUsed bool) types.List {
	var phaseAttrTypes = getResourcePhaseAttrTypesDEPRECATED()
	if phasesFromGo == nil {
		return types.ListNull(types.ObjectType{AttrTypes: phaseAttrTypes})
	}
	phasesList := make([]attr.Value, 0, len(phasesFromGo))
	for _, goPhase := range phasesFromGo {
		attrs := map[string]attr.Value{
			"id":                                    types.StringValue(goPhase.ID),
			"name":                                  types.StringValue(goPhase.Name),
			"automatic_deployment_targets":          util.FlattenStringList(goPhase.AutomaticDeploymentTargets),
			"optional_deployment_targets":           util.FlattenStringList(goPhase.OptionalDeploymentTargets),
			"minimum_environments_before_promotion": types.Int64Value(int64(goPhase.MinimumEnvironmentsBeforePromotion)),
			"is_optional_phase":                     types.BoolValue(goPhase.IsOptionalPhase),
			"is_priority_phase":                     types.BoolValue(goPhase.IsPriorityPhase),
		}
		if onlyRetentionWithoutStrategyWillBeUsed {
			attrs["release_retention_policy"] = util.Ternary(goPhase.ReleaseRetentionPolicy != nil, flattenResourceRetentionWithoutStrategyDEPRECATED(goPhase.ReleaseRetentionPolicy), ListNullRetentionWithoutStrategyDEPRECATED)
			attrs["tentacle_retention_policy"] = util.Ternary(goPhase.TentacleRetentionPolicy != nil, flattenResourceRetentionWithoutStrategyDEPRECATED(goPhase.TentacleRetentionPolicy), ListNullRetentionWithoutStrategyDEPRECATED)
			attrs["release_retention_with_strategy"] = ListNullRetention
			attrs["tentacle_retention_with_strategy"] = ListNullRetention
		} else {
			attrs["release_retention_policy"] = ListNullRetentionWithoutStrategyDEPRECATED
			attrs["tentacle_retention_policy"] = ListNullRetentionWithoutStrategyDEPRECATED
			attrs["release_retention_with_strategy"] = util.Ternary(goPhase.ReleaseRetentionPolicy != nil, flattenResourceRetention(goPhase.ReleaseRetentionPolicy), ListNullRetention)
			attrs["tentacle_retention_with_strategy"] = util.Ternary(goPhase.TentacleRetentionPolicy != nil, flattenResourceRetention(goPhase.TentacleRetentionPolicy), ListNullRetention)
		}
		phasesList = append(phasesList, types.ObjectValueMust(phaseAttrTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: phaseAttrTypes}, phasesList)

}

func flattenResourceRetention(retentionFromGo *core.RetentionPeriod) types.List {
	var retentionAttrTypes = getResourceRetentionAttrTypes()
	if retentionFromGo == nil {
		return ListNullRetention
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: retentionAttrTypes},
		[]attr.Value{
			types.ObjectValueMust(
				retentionAttrTypes,
				map[string]attr.Value{
					"strategy":         types.StringValue(retentionFromGo.Strategy),
					"unit":             types.StringValue(retentionFromGo.Unit),
					"quantity_to_keep": types.Int64Value(int64(retentionFromGo.QuantityToKeep)),
				},
			),
		},
	)
}
func flattenResourceRetentionWithoutStrategyDEPRECATED(retentionFromGo *core.RetentionPeriod) types.List {
	var retentionAttrTypes = getResourceRetentionWithoutStrategyAttrTypesDEPRECATED()
	if retentionFromGo == nil {
		return ListNullRetentionWithoutStrategyDEPRECATED
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: retentionAttrTypes},
		[]attr.Value{
			types.ObjectValueMust(
				retentionAttrTypes,
				map[string]attr.Value{
					"quantity_to_keep":    types.Int64Value(int64(retentionFromGo.QuantityToKeep)),
					"should_keep_forever": types.BoolValue(retentionFromGo.ShouldKeepForever),
					"unit":                types.StringValue(retentionFromGo.Unit),
				},
			),
		})
}

func expandLifecycle(lifecycleUserInput *lifecycleTypeResourceModel) *lifecycles.Lifecycle {
	if lifecycleUserInput == nil {
		return nil
	}

	lifecycleSentToGo := lifecycles.NewLifecycle(lifecycleUserInput.Name.ValueString())
	lifecycleSentToGo.Description = lifecycleUserInput.Description.ValueString()
	lifecycleSentToGo.SpaceID = lifecycleUserInput.SpaceID.ValueString()
	if !lifecycleUserInput.ID.IsNull() && lifecycleUserInput.ID.ValueString() != "" {
		lifecycleSentToGo.ID = lifecycleUserInput.ID.ValueString()
	}
	lifecycleSentToGo.Phases = expandPhases(lifecycleUserInput.Phase)
	lifecycleSentToGo.ReleaseRetentionPolicy = expandRetention(lifecycleUserInput.ReleaseRetention)
	lifecycleSentToGo.TentacleRetentionPolicy = expandRetention(lifecycleUserInput.TentacleRetention)

	return lifecycleSentToGo
}
func expandLifecycleDEPRECATED(lifecycleUserInput *lifecycleTypeResourceModelDEPRECATED, onlyRetentionWithoutStrategyWillBeUsed bool) *lifecycles.Lifecycle {
	if lifecycleUserInput == nil {
		return nil
	}
	lifecycleSentToGo := lifecycles.NewLifecycle(lifecycleUserInput.Name.ValueString())
	lifecycleSentToGo.Description = lifecycleUserInput.Description.ValueString()
	lifecycleSentToGo.SpaceID = lifecycleUserInput.SpaceID.ValueString()
	if !lifecycleUserInput.ID.IsNull() && lifecycleUserInput.ID.ValueString() != "" {
		lifecycleSentToGo.ID = lifecycleUserInput.ID.ValueString()
	}
	lifecycleSentToGo.Phases = expandPhasesDEPRECATED(lifecycleUserInput.Phase)

	if onlyRetentionWithoutStrategyWillBeUsed {
		lifecycleSentToGo.ReleaseRetentionPolicy = expandRetentionDEPRECATED(lifecycleUserInput.ReleaseRetentionWithoutStrategy)
		lifecycleSentToGo.TentacleRetentionPolicy = expandRetentionDEPRECATED(lifecycleUserInput.TentacleRetentionWithoutStrategy)
	} else {
		lifecycleSentToGo.ReleaseRetentionPolicy = expandRetention(lifecycleUserInput.ReleaseRetention)
		lifecycleSentToGo.TentacleRetentionPolicy = expandRetention(lifecycleUserInput.TentacleRetention)
	}

	return lifecycleSentToGo
}

func expandPhases(phasesUserInput types.List) []*lifecycles.Phase {
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
			phasePlan.ReleaseRetentionPolicy = expandRetention(v)
		}

		if v, ok := phaseAttrsUserInput["tentacle_retention_with_strategy"].(types.List); ok && !v.IsNull() {
			phasePlan.TentacleRetentionPolicy = expandRetention(v)
		}

		allPhasesPlan = append(allPhasesPlan, phasePlan)
	}

	return allPhasesPlan
}

func expandPhasesDEPRECATED(phasesUserInput types.List) []*lifecycles.Phase {
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
			phasePlan.ReleaseRetentionPolicy = expandRetention(v)
		}

		if v, ok := phaseAttrsUserInput["tentacle_retention_with_strategy"].(types.List); ok && !v.IsNull() {
			phasePlan.TentacleRetentionPolicy = expandRetention(v)
		}

		if v, ok := phaseAttrsUserInput["release_retention_policy"].(types.List); ok && !v.IsNull() {
			phasePlan.ReleaseRetentionPolicy = expandRetentionDEPRECATED(v)
		}

		if v, ok := phaseAttrsUserInput["tentacle_retention_policy"].(types.List); ok && !v.IsNull() {
			phasePlan.TentacleRetentionPolicy = expandRetentionDEPRECATED(v)
		}

		allPhasesPlan = append(allPhasesPlan, phasePlan)
	}

	return allPhasesPlan
}

func expandRetention(retentionUserInput types.List) *core.RetentionPeriod {
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
func expandRetentionDEPRECATED(v types.List) *core.RetentionPeriod {
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

var ListNullRetention = types.ListNull(types.ObjectType{AttrTypes: getResourceRetentionAttrTypes()})
var ListNullRetentionWithoutStrategyDEPRECATED = types.ListNull(types.ObjectType{AttrTypes: getResourceRetentionWithoutStrategyAttrTypesDEPRECATED()})

func getResourceRetentionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"strategy":         types.StringType,
		"quantity_to_keep": types.Int64Type,
		"unit":             types.StringType,
	}
}
func getResourceRetentionWithoutStrategyAttrTypesDEPRECATED() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}

func getResourcePhaseAttrTypes() map[string]attr.Type {
	var resourceAttrTypes = getResourceRetentionAttrTypes()
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"is_priority_phase":                     types.BoolType,
		"release_retention_with_strategy":       types.ListType{ElemType: types.ObjectType{AttrTypes: resourceAttrTypes}},
		"tentacle_retention_with_strategy":      types.ListType{ElemType: types.ObjectType{AttrTypes: resourceAttrTypes}},
	}
}

func getResourcePhaseAttrTypesDEPRECATED() map[string]attr.Type {
	var retentionAttrTypes = getResourceRetentionAttrTypes()
	var AttrTypesWithoutStrategy = getResourceRetentionWithoutStrategyAttrTypesDEPRECATED()
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"is_priority_phase":                     types.BoolType,
		"release_retention_policy":              types.ListType{ElemType: types.ObjectType{AttrTypes: AttrTypesWithoutStrategy}},
		"tentacle_retention_policy":             types.ListType{ElemType: types.ObjectType{AttrTypes: AttrTypesWithoutStrategy}},
		"release_retention_with_strategy":       types.ListType{ElemType: types.ObjectType{AttrTypes: retentionAttrTypes}},
		"tentacle_retention_with_strategy":      types.ListType{ElemType: types.ObjectType{AttrTypes: retentionAttrTypes}},
	}
}

func validateRetentionBlocksInConfigDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED, diag *diag.Diagnostics, retentionWithStrategyNotSupported bool) {
	isNewRetentionBlockInConfig := isNewBlockInRetentionConfigDEPRECATED(data)
	isRetentionWithoutStrategyBlockInConfig := isWithoutStrategyBlockInRetentionConfigDEPRECATED(data)
	if isNewRetentionBlockInConfig && isRetentionWithoutStrategyBlockInConfig {
		diag.AddError("Retention blocks conflict", "Both release_retention_with_strategy and release_retention_policy are used. Only one can be used at a time.")
	}
	if isNewRetentionBlockInConfig && retentionWithStrategyNotSupported {
		diag.AddError("Octopus Server Upgrade Required.", "retention_with_strategy is not supported on this Octopus Server version. Please upgrade to Octopus Server 2025.3 or later or use retention_policy.")
	}
}
func determineWhichRetentionBlocksAreInConfigDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED) (bool, bool, bool, bool) {
	isNewReleaseRetentionSet := attributeIsUsed(data.ReleaseRetention)
	isNewTentacleRetentionSet := attributeIsUsed(data.TentacleRetention)
	isReleaseRetentionWithoutStrategySet := attributeIsUsed(data.ReleaseRetentionWithoutStrategy)
	isTentacleRetentionWithoutStrategySet := attributeIsUsed(data.TentacleRetentionWithoutStrategy)

	return isNewReleaseRetentionSet, isNewTentacleRetentionSet, isReleaseRetentionWithoutStrategySet, isTentacleRetentionWithoutStrategySet
}
func isNewBlockInRetentionConfigDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED) bool {
	if attributeIsUsed(data.ReleaseRetention) || attributeIsUsed(data.TentacleRetention) {
		return true
	}
	for _, phase := range data.Phase.Elements() {
		phaseAttributes := phase.(types.Object).Attributes()
		releaseRetention := phaseAttributes["release_retention_with_strategy"].(types.List)
		tentacleRetention := phaseAttributes["tentacle_retention_with_strategy"].(types.List)

		if attributeIsUsed(releaseRetention) || attributeIsUsed(tentacleRetention) {
			return true
		}
	}
	return false
}
func isWithoutStrategyBlockInRetentionConfigDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED) bool {
	if attributeIsUsed(data.ReleaseRetentionWithoutStrategy) || attributeIsUsed(data.TentacleRetentionWithoutStrategy) {
		return true
	}
	for _, phase := range data.Phase.Elements() {
		phaseAttributes := phase.(types.Object).Attributes()
		releaseRetentionWithoutStrategy := phaseAttributes["release_retention_policy"].(types.List)
		tentacleRetentionWithoutStrategy := phaseAttributes["tentacle_retention_policy"].(types.List)
		if attributeIsUsed(releaseRetentionWithoutStrategy) || attributeIsUsed(tentacleRetentionWithoutStrategy) {
			return true
		}
	}
	return false
}
func isRetentionWithoutStrategyToBeUsedForLifecycleDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED) bool {
	releaseRetentionIsInConfig := isNewBlockInRetentionConfigDEPRECATED(data)
	if releaseRetentionIsInConfig {
		return false
	}
	return true
}
