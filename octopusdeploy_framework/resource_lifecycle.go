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

type lifecycleTypeResourceModelDEPRECATED struct {
	SpaceID                          types.String `tfsdk:"space_id"`
	Name                             types.String `tfsdk:"name"`
	Description                      types.String `tfsdk:"description"`
	Phase                            types.List   `tfsdk:"phase"`
	ReleaseRetentionWithoutStrategy  types.List   `tfsdk:"release_retention_policy"`
	TentacleRetentionWithoutStrategy types.List   `tfsdk:"tentacle_retention_policy"`

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

func (r *lifecycleTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *lifecycleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data *lifecycleTypeResourceModelDEPRECATED
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	isReleaseRetentionWithoutStrategySet := attributeIsUsed(data.ReleaseRetentionWithoutStrategy)
	isTentacleRetentionWithoutStrategySet := attributeIsUsed(data.TentacleRetentionWithoutStrategy)
	initialRetentionWithoutStrategySetting := setInitialRetentionDEPRECATED(data)

	lifecycle := expandLifecycleDEPRECATED(data)
	newLifecycle, err := lifecycles.Add(r.Config.Client, lifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
		return
	}

	data = flattenResourceLifecycleDEPRECATED(newLifecycle)
	removeInitialRetentionDEPRECATED(data, isReleaseRetentionWithoutStrategySet, isTentacleRetentionWithoutStrategySet, initialRetentionWithoutStrategySetting)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
func (r *lifecycleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *lifecycleTypeResourceModelDEPRECATED
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	isReleaseRetentionWithoutStrategySet := attributeIsUsed(data.ReleaseRetentionWithoutStrategy)
	isTentacleRetentionWithoutStrategySet := attributeIsUsed(data.TentacleRetentionWithoutStrategy)
	initialRetentionWithoutStrategySetting := setInitialRetentionDEPRECATED(data)

	lifecycle, err := lifecycles.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "lifecycle"); err != nil {
			resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
		}
		return
	}
	data = flattenResourceLifecycleDEPRECATED(lifecycle)

	removeInitialRetentionDEPRECATED(data, isReleaseRetentionWithoutStrategySet, isTentacleRetentionWithoutStrategySet, initialRetentionWithoutStrategySetting)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data, state *lifecycleTypeResourceModelDEPRECATED
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	isReleaseRetentionWithoutStrategySet := attributeIsUsed(data.ReleaseRetentionWithoutStrategy)
	isTentacleRetentionWithoutStrategySet := attributeIsUsed(data.TentacleRetentionWithoutStrategy)
	initialRetentionWithoutStrategySetting := setInitialRetentionDEPRECATED(data)

	tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", data.ID.ValueString()))
	lifecycle := expandLifecycleDEPRECATED(data)
	lifecycle.ID = state.ID.ValueString()
	updatedLifecycle, err := lifecycles.Update(r.Config.Client, lifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
		return
	}
	data = flattenResourceLifecycleDEPRECATED(updatedLifecycle)

	removeInitialRetentionDEPRECATED(data, isReleaseRetentionWithoutStrategySet, isTentacleRetentionWithoutStrategySet, initialRetentionWithoutStrategySetting)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data lifecycleTypeResourceModelDEPRECATED
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := lifecycles.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete lifecycle", err.Error())
		return
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

func setInitialRetentionDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED) types.List {
	var initialRetentionWithoutStrategySetting types.List
	initialRetentionWithoutStrategySetting = flattenResourceRetentionDEPRECATED(core.CountBasedRetentionPeriod(30, "Days"))
	setInitialRetentionWithoutStrategyBlockDEPRECATED(data, initialRetentionWithoutStrategySetting)
	return initialRetentionWithoutStrategySetting
}
func removeInitialRetentionDEPRECATED(data *lifecycleTypeResourceModelDEPRECATED, isReleaseRetentionWithoutStrategySet bool, isTentacleRetentionWithoutStrategySet bool, initialRetentionWithoutStrategySetting types.List) {
	if !isReleaseRetentionWithoutStrategySet && (data.ReleaseRetentionWithoutStrategy.Equal(initialRetentionWithoutStrategySetting) || data.ReleaseRetentionWithoutStrategy.IsNull()) {
		data.ReleaseRetentionWithoutStrategy = ListNullRetentionWithoutStrategyDEPRECATED
	}
	if !isTentacleRetentionWithoutStrategySet && (data.TentacleRetentionWithoutStrategy.Equal(initialRetentionWithoutStrategySetting) || data.TentacleRetentionWithoutStrategy.IsNull()) {
		data.TentacleRetentionWithoutStrategy = ListNullRetentionWithoutStrategyDEPRECATED
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

func flattenResourceLifecycleDEPRECATED(lifecycle *lifecycles.Lifecycle) *lifecycleTypeResourceModelDEPRECATED {
	var flattenedLifecycle *lifecycleTypeResourceModelDEPRECATED
	flattenedLifecycle = &lifecycleTypeResourceModelDEPRECATED{
		SpaceID:     types.StringValue(lifecycle.SpaceID),
		Name:        types.StringValue(lifecycle.Name),
		Description: types.StringValue(lifecycle.Description),
		Phase:       flattenResourcePhasesDEPRECATED(lifecycle.Phases),
	}

	flattenedLifecycle.ReleaseRetentionWithoutStrategy = flattenResourceRetentionDEPRECATED(lifecycle.ReleaseRetentionPolicy)
	flattenedLifecycle.TentacleRetentionWithoutStrategy = flattenResourceRetentionDEPRECATED(lifecycle.TentacleRetentionPolicy)

	flattenedLifecycle.ID = types.StringValue(lifecycle.GetID())
	return flattenedLifecycle
}

func flattenResourcePhasesDEPRECATED(phasesFromGo []*lifecycles.Phase) types.List {
	var phaseAttrTypes = getResourcePhaseAttrTypesDEPRECATED()
	if phasesFromGo == nil {
		return types.ListNull(types.ObjectType{AttrTypes: phaseAttrTypes})
	}
	phasesList := make([]attr.Value, 0, len(phasesFromGo))
	for _, phaseFromGo := range phasesFromGo {
		attrs := map[string]attr.Value{
			"id":                                    types.StringValue(phaseFromGo.ID),
			"name":                                  types.StringValue(phaseFromGo.Name),
			"automatic_deployment_targets":          util.FlattenStringList(phaseFromGo.AutomaticDeploymentTargets),
			"optional_deployment_targets":           util.FlattenStringList(phaseFromGo.OptionalDeploymentTargets),
			"minimum_environments_before_promotion": types.Int64Value(int64(phaseFromGo.MinimumEnvironmentsBeforePromotion)),
			"is_optional_phase":                     types.BoolValue(phaseFromGo.IsOptionalPhase),
			"is_priority_phase":                     types.BoolValue(phaseFromGo.IsPriorityPhase),
		}

		attrs["release_retention_policy"] = util.Ternary(phaseFromGo.ReleaseRetentionPolicy != nil, flattenResourceRetentionDEPRECATED(phaseFromGo.ReleaseRetentionPolicy), ListNullRetentionWithoutStrategyDEPRECATED)
		attrs["tentacle_retention_policy"] = util.Ternary(phaseFromGo.TentacleRetentionPolicy != nil, flattenResourceRetentionDEPRECATED(phaseFromGo.TentacleRetentionPolicy), ListNullRetentionWithoutStrategyDEPRECATED)

		phasesList = append(phasesList, types.ObjectValueMust(phaseAttrTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: phaseAttrTypes}, phasesList)

}

func flattenResourceRetentionDEPRECATED(retentionFromGo *core.RetentionPeriod) types.List {
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

func expandLifecycleDEPRECATED(lifecycleUserInput *lifecycleTypeResourceModelDEPRECATED) *lifecycles.Lifecycle {
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

	lifecycleSentToGo.ReleaseRetentionPolicy = expandRetentionDEPRECATED(lifecycleUserInput.ReleaseRetentionWithoutStrategy)
	lifecycleSentToGo.TentacleRetentionPolicy = expandRetentionDEPRECATED(lifecycleUserInput.TentacleRetentionWithoutStrategy)

	return lifecycleSentToGo
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

var ListNullRetentionWithoutStrategyDEPRECATED = types.ListNull(types.ObjectType{AttrTypes: getResourceRetentionWithoutStrategyAttrTypesDEPRECATED()})

func getResourceRetentionWithoutStrategyAttrTypesDEPRECATED() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}

func getResourcePhaseAttrTypesDEPRECATED() map[string]attr.Type {
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
	}
}
