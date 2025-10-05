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
		if !r.Config.IsVersionSameOrGreaterThan("2025.3") {
			resp.Diagnostics.AddError("retention_with_strategy is not supported in this Octopus Deploy version", "Please upgrade your Octopus Deploy version to 2025.3 or later")
			return
		}
	}
}

func (r *lifecycleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

func (r *lifecycleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

func (r *lifecycleTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func flattenResourcePhases(goPhases []*lifecycles.Phase) types.List {
	var attributeTypes = getResourcePhaseAttrTypes()
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
			"release_retention_with_strategy":       util.Ternary(goPhase.ReleaseRetentionPolicy != nil, flattenResourceRetention(goPhase.ReleaseRetentionPolicy), ListNullRetention),
			"tentacle_retention_with_strategy":      util.Ternary(goPhase.TentacleRetentionPolicy != nil, flattenResourceRetention(goPhase.TentacleRetentionPolicy), ListNullRetention),
		}
		phasesList = append(phasesList, types.ObjectValueMust(attributeTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: attributeTypes}, phasesList)
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

func expandLifecycle(lifecycleUserInput *lifecycleTypeResourceModel) *lifecycles.Lifecycle {
	if lifecycleUserInput == nil {
		return nil
	}

	lifecyclePlan := lifecycles.NewLifecycle(lifecycleUserInput.Name.ValueString())
	lifecyclePlan.Description = lifecycleUserInput.Description.ValueString()
	lifecyclePlan.SpaceID = lifecycleUserInput.SpaceID.ValueString()
	if !lifecycleUserInput.ID.IsNull() && lifecycleUserInput.ID.ValueString() != "" {
		lifecyclePlan.ID = lifecycleUserInput.ID.ValueString()
	}
	lifecyclePlan.Phases = expandPhases(lifecycleUserInput.Phase)
	lifecyclePlan.ReleaseRetentionPolicy = expandRetention(lifecycleUserInput.ReleaseRetention)
	lifecyclePlan.TentacleRetentionPolicy = expandRetention(lifecycleUserInput.TentacleRetention)

	return lifecyclePlan
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

var ListNullRetention = types.ListNull(types.ObjectType{AttrTypes: getResourceRetentionAttrTypes()})

func getResourceRetentionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"strategy":         types.StringType,
		"quantity_to_keep": types.Int64Type,
		"unit":             types.StringType,
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
