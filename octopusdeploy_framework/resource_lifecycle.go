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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

type lifecycleTypeResource struct {
	*Config
	retentionWithStrategyNotSupported bool
	allowDeprecatedRetention          bool
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

type deprecatedLifecycleTypeResourceModel struct {
	SpaceID                     types.String `tfsdk:"space_id"`
	Name                        types.String `tfsdk:"name"`
	Description                 types.String `tfsdk:"description"`
	Phase                       types.List   `tfsdk:"phase"`
	DeprecatedReleaseRetention  types.List   `tfsdk:"release_retention_policy"`
	DeprecatedTentacleRetention types.List   `tfsdk:"tentacle_retention_policy"`
	ReleaseRetention            types.List   `tfsdk:"release_retention_with_strategy"`
	TentacleRetention           types.List   `tfsdk:"tentacle_retention_with_strategy"`

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
	resp.Schema = schemas.LifecycleSchema{r.allowDeprecatedRetention}.GetResourceSchema()
}

func (r *lifecycleTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
	if r.Config != nil {
		if !r.Config.IsVersionSameOrGreaterThan("2025.3") {
			r.retentionWithStrategyNotSupported = true
		}
	}
}

func (r *lifecycleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	if r.allowDeprecatedRetention {

		var data *deprecatedLifecycleTypeResourceModel
		resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		ValidateRetentionBlocksInPlan(data, &resp.Diagnostics, r.retentionWithStrategyNotSupported)
		deprecatedRetentionUsed := IsDeprecatedRetentionUsed(data)

		if deprecatedRetentionUsed {
			tflog.Debug(ctx, fmt.Sprintf("deprecatedRetentionUsed: `%v`", deprecatedRetentionUsed))
			initialRetentionSetting := core.CountBasedRetentionPeriod(30, "Days")
			initialRetentionAsNewBlock := flattenResourceRetention(initialRetentionSetting)
			initialRetentionAsDeprecatedBlock := flattenDeprecatedResourceRetention(initialRetentionSetting)
			hasUserDefinedDeprecatedReleaseRetention, hasUserDefinedDeprecatedTentacleRetention := setInitialRetentionForOldBlock(data, initialRetentionAsDeprecatedBlock)

			lifecycle := expandDeprecatedLifecycle(data, deprecatedRetentionUsed)
			newLifecycle, err := lifecycles.Add(r.Config.Client, lifecycle)
			if err != nil {
				resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
				return
			}

			data = flattenDeprecatedResourceLifecycle(newLifecycle, deprecatedRetentionUsed)
			removeDefaultRetentionFromDeprecatedAndNewUnsetBlocks(data, false, false, hasUserDefinedDeprecatedReleaseRetention, hasUserDefinedDeprecatedTentacleRetention, initialRetentionAsNewBlock, initialRetentionAsDeprecatedBlock, ctx)
		} else {
			tflog.Debug(ctx, fmt.Sprintf("deprecatedRetentionUsed: `%v`", deprecatedRetentionUsed))
			initialRetentionSetting := core.SpaceDefaultRetentionPeriod()
			initialRetentionAsNewBlock := flattenResourceRetention(initialRetentionSetting)
			initialRetentionAsDeprecatedBlock := flattenDeprecatedResourceRetention(initialRetentionSetting)
			hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := setInitialRetentionForNewBlock(data, initialRetentionAsDeprecatedBlock)

			lifecycle := expandDeprecatedLifecycle(data, deprecatedRetentionUsed)
			newLifecycle, err := lifecycles.Add(r.Config.Client, lifecycle)
			if err != nil {
				resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
				return
			}

			data = flattenDeprecatedResourceLifecycle(newLifecycle, deprecatedRetentionUsed)
			removeDefaultRetentionFromDeprecatedAndNewUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention, false, false, initialRetentionAsNewBlock, initialRetentionAsDeprecatedBlock, ctx)
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	} else {
		var data *lifecycleTypeResourceModel
		resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		initialRetentionSetting := flattenResourceRetention(core.SpaceDefaultRetentionPeriod())
		isReleaseRetentionDefined, isTentacleRetentionDefined := setDefaultRetention(data, initialRetentionSetting)

		lifecycle := expandLifecycle(data)
		newLifecycle, err := lifecycles.Add(r.Config.Client, lifecycle)
		if err != nil {
			resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
			return
		}
		data = flattenResourceLifecycle(newLifecycle)

		removeDefaultRetentionFromUnsetBlocks(data, isReleaseRetentionDefined, isTentacleRetentionDefined, initialRetentionSetting)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}
func (r *lifecycleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.allowDeprecatedRetention {
		var data *deprecatedLifecycleTypeResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		ValidateRetentionBlocksInPlan(data, &resp.Diagnostics, r.retentionWithStrategyNotSupported)
		deprecatedRetentionUsed := IsDeprecatedRetentionUsed(data)

		if deprecatedRetentionUsed {
			initialRetentionSetting := core.CountBasedRetentionPeriod(30, "Days")
			initialRetentionAsNewBlock := flattenResourceRetention(initialRetentionSetting)
			initialRetentionAsDeprecatedBlock := flattenDeprecatedResourceRetention(initialRetentionSetting)
			hasUserDefinedDeprecatedReleaseRetention, hasUserDefinedDeprecatedTentacleRetention := setInitialRetentionForOldBlock(data, initialRetentionAsDeprecatedBlock)

			lifecycle, err := lifecycles.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
			if err != nil {
				if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "lifecycle"); err != nil {
					resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
				}
				return
			}
			data = flattenDeprecatedResourceLifecycle(lifecycle, deprecatedRetentionUsed)

			removeDefaultRetentionFromDeprecatedAndNewUnsetBlocks(data, false, false, hasUserDefinedDeprecatedReleaseRetention, hasUserDefinedDeprecatedTentacleRetention, initialRetentionAsNewBlock, initialRetentionAsDeprecatedBlock, ctx)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		} else {
			initialRetentionSetting := core.SpaceDefaultRetentionPeriod()
			initialRetentionAsNewBlock := flattenResourceRetention(initialRetentionSetting)
			initialRetentionAsDeprecatedBlock := flattenDeprecatedResourceRetention(initialRetentionSetting)
			hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := setInitialRetentionForNewBlock(data, initialRetentionAsDeprecatedBlock)

			lifecycle, err := lifecycles.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
			if err != nil {
				if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "lifecycle"); err != nil {
					resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
				}
				return
			}
			data = flattenDeprecatedResourceLifecycle(lifecycle, deprecatedRetentionUsed)

			removeDefaultRetentionFromDeprecatedAndNewUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention, false, false, initialRetentionAsNewBlock, initialRetentionAsDeprecatedBlock, ctx)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}

	} else {
		var data *lifecycleTypeResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		initialRetentionSetting := flattenResourceRetention(core.SpaceDefaultRetentionPeriod())
		isReleaseRetentionDefined, isTentacleRetentionDefined := setDefaultRetention(data, initialRetentionSetting)

		lifecycle, err := lifecycles.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
		if err != nil {
			if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "lifecycle"); err != nil {
				resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
			}
			return
		}
		data = flattenResourceLifecycle(lifecycle)

		removeDefaultRetentionFromUnsetBlocks(data, isReleaseRetentionDefined, isTentacleRetentionDefined, initialRetentionSetting)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *lifecycleTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.allowDeprecatedRetention {
		var data, state *deprecatedLifecycleTypeResourceModel
		resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}

		ValidateRetentionBlocksInPlan(data, &resp.Diagnostics, r.retentionWithStrategyNotSupported)
		deprecatedRetentionUsed := IsDeprecatedRetentionUsed(data)

		if deprecatedRetentionUsed {
			initialRetentionSetting := core.CountBasedRetentionPeriod(30, "Days")
			initialRetentionAsNewBlock := flattenResourceRetention(initialRetentionSetting)
			initialRetentionAsDeprecatedBlock := flattenDeprecatedResourceRetention(initialRetentionSetting)
			hasUserDefinedDeprecatedReleaseRetention, hasUserDefinedDeprecatedTentacleRetention := setInitialRetentionForOldBlock(data, initialRetentionAsDeprecatedBlock)

			tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", data.ID.ValueString()))
			lifecycle := expandDeprecatedLifecycle(data, deprecatedRetentionUsed)
			lifecycle.ID = state.ID.ValueString()
			updatedLifecycle, err := lifecycles.Update(r.Config.Client, lifecycle)
			if err != nil {
				resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
				return
			}
			data = flattenDeprecatedResourceLifecycle(updatedLifecycle, deprecatedRetentionUsed)

			removeDefaultRetentionFromDeprecatedAndNewUnsetBlocks(data, false, false, hasUserDefinedDeprecatedReleaseRetention, hasUserDefinedDeprecatedTentacleRetention, initialRetentionAsNewBlock, initialRetentionAsDeprecatedBlock, ctx)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		} else {
			initialRetentionSetting := core.SpaceDefaultRetentionPeriod()
			initialRetentionAsNewBlock := flattenResourceRetention(initialRetentionSetting)
			initialRetentionAsDeprecatedBlock := flattenDeprecatedResourceRetention(initialRetentionSetting)
			hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention := setInitialRetentionForNewBlock(data, initialRetentionAsDeprecatedBlock)

			tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", data.ID.ValueString()))
			lifecycle := expandDeprecatedLifecycle(data, deprecatedRetentionUsed)
			lifecycle.ID = state.ID.ValueString()
			updatedLifecycle, err := lifecycles.Update(r.Config.Client, lifecycle)
			if err != nil {
				resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
				return
			}
			data = flattenDeprecatedResourceLifecycle(updatedLifecycle, deprecatedRetentionUsed)

			removeDefaultRetentionFromDeprecatedAndNewUnsetBlocks(data, hasUserDefinedReleaseRetention, hasUserDefinedTentacleRetention, false, false, initialRetentionAsNewBlock, initialRetentionAsDeprecatedBlock, ctx)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}

	} else {
		var data, state *lifecycleTypeResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}

		initialRetentionSetting := flattenResourceRetention(core.SpaceDefaultRetentionPeriod())
		isReleaseRetentionDefined, isTentacleRetentionDefined := setDefaultRetention(data, initialRetentionSetting)

		tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", data.ID.ValueString()))
		lifecycle := expandLifecycle(data)
		lifecycle.ID = state.ID.ValueString()
		updatedLifecycle, err := lifecycles.Update(r.Config.Client, lifecycle)
		if err != nil {
			resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
			return
		}
		data = flattenResourceLifecycle(updatedLifecycle)

		removeDefaultRetentionFromUnsetBlocks(data, isReleaseRetentionDefined, isTentacleRetentionDefined, initialRetentionSetting)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *lifecycleTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.allowDeprecatedRetention {
		var data deprecatedLifecycleTypeResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if err := lifecycles.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
			resp.Diagnostics.AddError("unable to delete lifecycle", err.Error())
			return
		}
	} else {
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

}
func ValidateRetentionBlocksInPlan(data *deprecatedLifecycleTypeResourceModel, diag *diag.Diagnostics, onlyDeprecatedRetentionIsSupported bool) {
	latestRetentionIsInPlan := IsLatestRetentionInPlan(data)
	deprecatedRetentionIsInPlan := IsDeprecatedRetentionInPlan(data)
	if latestRetentionIsInPlan && deprecatedRetentionIsInPlan {
		diag.AddError("Retention blocks conflict", "Both release_retention_with_strategy and release_retention_policy are used. Only one can be used at a time.")
	}
	if latestRetentionIsInPlan && onlyDeprecatedRetentionIsSupported {
		diag.AddError("retention with strategy is not supported on this Octopus Server version. Please upgrade to Octopus Server 2025.3 or later.", "")
	}
}

func IsLatestRetentionInPlan(data *deprecatedLifecycleTypeResourceModel) bool {
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

func IsDeprecatedRetentionInPlan(data *deprecatedLifecycleTypeResourceModel) bool {
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

func IsDeprecatedRetentionUsed(data *deprecatedLifecycleTypeResourceModel) bool {
	releaseRetentionIsInPlan := IsLatestRetentionInPlan(data)
	if releaseRetentionIsInPlan {
		return false
	}
	return true
}

func setInitialRetentionForOldBlock(data *deprecatedLifecycleTypeResourceModel, initialRetentionSettingForOldBlock types.List) (bool, bool) {
	hasUserDefinedDeprecatedReleaseRetention := attributeIsUsed(data.DeprecatedReleaseRetention)
	hasUserDefinedDeprecatedTentacleRetention := attributeIsUsed(data.DeprecatedTentacleRetention)
	if !hasUserDefinedDeprecatedReleaseRetention {
		data.DeprecatedReleaseRetention = initialRetentionSettingForOldBlock
	}
	if !hasUserDefinedDeprecatedTentacleRetention {
		data.DeprecatedTentacleRetention = initialRetentionSettingForOldBlock
	}

	return hasUserDefinedDeprecatedReleaseRetention, hasUserDefinedDeprecatedTentacleRetention
}

func setInitialRetentionForNewBlock(data *deprecatedLifecycleTypeResourceModel, initialRetentionSettingForNewBlock types.List) (bool, bool) {
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

func setDefaultRetention(data *lifecycleTypeResourceModel, initialRetentionSetting types.List) (bool, bool) {
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

func removeDefaultRetentionFromDeprecatedAndNewUnsetBlocks(data *deprecatedLifecycleTypeResourceModel, isReleaseRetentionDefined bool, isTentacleRetentionDefined bool, isDeprecatedReleaseRetentionDefined bool, isDeprecatedTentacleRetentionDefined bool, initialRetentionSetting types.List, initialDeprecatedRetentionSetting types.List, ctx context.Context) {
	if !isReleaseRetentionDefined && (data.ReleaseRetention.Equal(initialRetentionSetting) || data.ReleaseRetention.IsNull()) {
		data.ReleaseRetention = ListNullRetention
	}
	if !isTentacleRetentionDefined && (data.TentacleRetention.Equal(initialRetentionSetting) || data.TentacleRetention.IsNull()) {
		data.TentacleRetention = ListNullRetention
	}

	if !isDeprecatedReleaseRetentionDefined && (data.DeprecatedReleaseRetention.Equal(initialDeprecatedRetentionSetting) || data.DeprecatedReleaseRetention.IsNull()) {
		data.DeprecatedReleaseRetention = ListNullDeprecatedRetention
	}

	if !isDeprecatedTentacleRetentionDefined && (data.DeprecatedTentacleRetention.Equal(initialDeprecatedRetentionSetting) || data.DeprecatedTentacleRetention.IsNull()) {
		data.DeprecatedTentacleRetention = ListNullDeprecatedRetention
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
func flattenDeprecatedResourceLifecycle(lifecycle *lifecycles.Lifecycle, deprecatedRetentionUsed bool) *deprecatedLifecycleTypeResourceModel {
	var flattenedLifecycle *deprecatedLifecycleTypeResourceModel
	flattenedLifecycle = &deprecatedLifecycleTypeResourceModel{
		SpaceID:     types.StringValue(lifecycle.SpaceID),
		Name:        types.StringValue(lifecycle.Name),
		Description: types.StringValue(lifecycle.Description),
		Phase:       flattenDeprecatedResourcePhases(lifecycle.Phases, deprecatedRetentionUsed),
	}
	if deprecatedRetentionUsed {
		flattenedLifecycle.DeprecatedReleaseRetention = flattenDeprecatedResourceRetention(lifecycle.ReleaseRetentionPolicy)
		flattenedLifecycle.DeprecatedTentacleRetention = flattenDeprecatedResourceRetention(lifecycle.TentacleRetentionPolicy)
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
func flattenDeprecatedResourcePhases(goPhases []*lifecycles.Phase, deprecatedRetentionUsed bool) types.List {
	var phaseAttrTypes = getDeprecatedResourcePhaseAttrTypes()
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
		}
		if deprecatedRetentionUsed {
			attrs["release_retention_policy"] = util.Ternary(goPhase.ReleaseRetentionPolicy != nil, flattenDeprecatedResourceRetention(goPhase.ReleaseRetentionPolicy), ListNullDeprecatedRetention)
			attrs["tentacle_retention_policy"] = util.Ternary(goPhase.TentacleRetentionPolicy != nil, flattenDeprecatedResourceRetention(goPhase.TentacleRetentionPolicy), ListNullDeprecatedRetention)
			attrs["release_retention_with_strategy"] = ListNullRetention
			attrs["tentacle_retention_with_strategy"] = ListNullRetention
		} else {
			attrs["release_retention_policy"] = ListNullDeprecatedRetention
			attrs["tentacle_retention_policy"] = ListNullDeprecatedRetention
			attrs["release_retention_with_strategy"] = util.Ternary(goPhase.ReleaseRetentionPolicy != nil, flattenResourceRetention(goPhase.ReleaseRetentionPolicy), ListNullRetention)
			attrs["tentacle_retention_with_strategy"] = util.Ternary(goPhase.TentacleRetentionPolicy != nil, flattenResourceRetention(goPhase.TentacleRetentionPolicy), ListNullRetention)
		}
		phasesList = append(phasesList, types.ObjectValueMust(phaseAttrTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: phaseAttrTypes}, phasesList)

}

func flattenResourceRetention(goRetention *core.RetentionPeriod) types.List {
	var retentionAttrTypes = getResourceRetentionAttrTypes()
	if goRetention == nil {
		return ListNullRetention
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: retentionAttrTypes},
		[]attr.Value{
			types.ObjectValueMust(
				retentionAttrTypes,
				map[string]attr.Value{
					"strategy":         types.StringValue(goRetention.Strategy),
					"unit":             types.StringValue(goRetention.Unit),
					"quantity_to_keep": types.Int64Value(int64(goRetention.QuantityToKeep)),
				},
			),
		},
	)
}

func flattenDeprecatedResourceRetention(goRetention *core.RetentionPeriod) types.List {
	var retentionAttrTypes = getDeprecatedResourceRetentionAttTypes()
	if goRetention == nil {
		return ListNullDeprecatedRetention
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: retentionAttrTypes},
		[]attr.Value{
			types.ObjectValueMust(
				retentionAttrTypes,
				map[string]attr.Value{
					"quantity_to_keep":    types.Int64Value(int64(goRetention.QuantityToKeep)),
					"should_keep_forever": types.BoolValue(goRetention.ShouldKeepForever),
					"unit":                types.StringValue(goRetention.Unit),
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

func expandDeprecatedLifecycle(lifecycleUserInput *deprecatedLifecycleTypeResourceModel, deprecatedRetentionUsed bool) *lifecycles.Lifecycle {
	if lifecycleUserInput == nil {
		return nil
	}
	lifecycleSentToGo := lifecycles.NewLifecycle(lifecycleUserInput.Name.ValueString())
	lifecycleSentToGo.Description = lifecycleUserInput.Description.ValueString()
	lifecycleSentToGo.SpaceID = lifecycleUserInput.SpaceID.ValueString()
	if !lifecycleUserInput.ID.IsNull() && lifecycleUserInput.ID.ValueString() != "" {
		lifecycleSentToGo.ID = lifecycleUserInput.ID.ValueString()
	}
	lifecycleSentToGo.Phases = expandDeprecatedPhases(lifecycleUserInput.Phase)

	if deprecatedRetentionUsed {
		lifecycleSentToGo.ReleaseRetentionPolicy = expandDeprecatedRetention(lifecycleUserInput.DeprecatedReleaseRetention)
		lifecycleSentToGo.TentacleRetentionPolicy = expandDeprecatedRetention(lifecycleUserInput.DeprecatedTentacleRetention)
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

func expandDeprecatedPhases(phasesUserInput types.List) []*lifecycles.Phase {
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
			phasePlan.ReleaseRetentionPolicy = expandDeprecatedRetention(v)
		}

		if v, ok := phaseAttrsUserInput["tentacle_retention_policy"].(types.List); ok && !v.IsNull() {
			phasePlan.TentacleRetentionPolicy = expandDeprecatedRetention(v)
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

func expandDeprecatedRetention(v types.List) *core.RetentionPeriod {
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
var ListNullDeprecatedRetention = types.ListNull(types.ObjectType{AttrTypes: getDeprecatedResourceRetentionAttTypes()})

func getResourceRetentionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"strategy":         types.StringType,
		"quantity_to_keep": types.Int64Type,
		"unit":             types.StringType,
	}
}

func getDeprecatedResourceRetentionAttTypes() map[string]attr.Type {
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

func getDeprecatedResourcePhaseAttrTypes() map[string]attr.Type {
	var retentionAttrTypes = getResourceRetentionAttrTypes()
	var deprecatedRetentionAttrTypes = getDeprecatedResourceRetentionAttTypes()
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"is_priority_phase":                     types.BoolType,
		"release_retention_policy":              types.ListType{ElemType: types.ObjectType{AttrTypes: deprecatedRetentionAttrTypes}},
		"tentacle_retention_policy":             types.ListType{ElemType: types.ObjectType{AttrTypes: deprecatedRetentionAttrTypes}},
		"release_retention_with_strategy":       types.ListType{ElemType: types.ObjectType{AttrTypes: retentionAttrTypes}},
		"tentacle_retention_with_strategy":      types.ListType{ElemType: types.ObjectType{AttrTypes: retentionAttrTypes}},
	}
}
