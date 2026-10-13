package octopusdeploy_framework

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const spaceDefaultLifecycleResourceName = "space_default_lifecycle"

var spaceDefaultLifecyclePhaseObjectType = types.ObjectType{AttrTypes: schemas.SpaceDefaultLifecyclePhaseAttrTypes()}
var spaceDefaultLifecycleRetentionObjectType = types.ObjectType{AttrTypes: schemas.SpaceDefaultLifecycleRetentionAttrTypes()}

type spaceDefaultLifecycleResource struct {
	*Config
}

var _ resource.Resource = &spaceDefaultLifecycleResource{}
var _ resource.ResourceWithImportState = &spaceDefaultLifecycleResource{}

func NewSpaceDefaultLifecycleResource() resource.Resource {
	return &spaceDefaultLifecycleResource{}
}

func (r *spaceDefaultLifecycleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(spaceDefaultLifecycleResourceName)
}

func (r *spaceDefaultLifecycleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.SpaceDefaultLifecycleSchema{}.GetResourceSchema()
}

func (r *spaceDefaultLifecycleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
	if r.Config != nil {
		resp.Diagnostics.Append(r.Config.EnsureResourceCompatibilityByVersion(spaceDefaultLifecycleResourceName, "2025.3")...)
	}
}

// The default lifecycle cannot be created via the API; Octopus creates it with the space.
// Create adopts the existing lifecycle and updates it to match the plan.
func (r *spaceDefaultLifecycleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.SpaceDefaultLifecycleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	util.Create(ctx, spaceDefaultLifecycleResourceName, plan)

	existing, err := r.findSpaceDefaultLifecycle(plan.SpaceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load default lifecycle of space", err.Error())
		return
	}

	resp.Diagnostics.Append(r.update(ctx, existing, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	util.Created(ctx, spaceDefaultLifecycleResourceName)
}

func (r *spaceDefaultLifecycleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.SpaceDefaultLifecycleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var lifecycle *lifecycles.Lifecycle
	var err error
	if state.ID.IsNull() || state.ID.ValueString() == "" {
		// imported by space ID
		lifecycle, err = r.findSpaceDefaultLifecycle(state.SpaceID.ValueString())
	} else {
		lifecycle, err = lifecycles.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	}
	if err != nil {
		if isNotFound(err) {
			tflog.Info(ctx, fmt.Sprintf("default lifecycle (%s) of space (%s) not found; removing from state", state.ID.ValueString(), state.SpaceID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to load default lifecycle of space", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenSpaceDefaultLifecycle(ctx, lifecycle, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *spaceDefaultLifecycleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state schemas.SpaceDefaultLifecycleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	existing, err := lifecycles.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load default lifecycle of space", err.Error())
		return
	}

	resp.Diagnostics.Append(r.update(ctx, existing, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// The default lifecycle cannot be deleted via the API; it is deleted together with its space.
// Delete clears the phases so that environments referenced by them can be deleted afterwards.
func (r *spaceDefaultLifecycleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.SpaceDefaultLifecycleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	existing, err := lifecycles.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("unable to load default lifecycle of space", err.Error())
		return
	}

	if len(existing.Phases) == 0 {
		return
	}

	existing.Phases = []*lifecycles.Phase{}
	if _, err := lifecycles.Update(r.Client, existing); err != nil {
		resp.Diagnostics.AddError("unable to clear phases of default lifecycle of space", err.Error())
	}
}

// ImportState accepts either "<space-id>", in which case the default lifecycle of that space is looked up on read,
// or "<space-id>:<lifecycle-id>" to import a specific lifecycle, e.g. one that has been renamed.
func (r *spaceDefaultLifecycleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	spaceID, lifecycleID, hasLifecycleID := strings.Cut(req.ID, ":")
	if spaceID == "" || (hasLifecycleID && lifecycleID == "") {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected an import ID in the format '<space-id>' or '<space-id>:<lifecycle-id>', got '%s'.", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), spaceID)...)
	if hasLifecycleID {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), lifecycleID)...)
	}
}

func (r *spaceDefaultLifecycleResource) update(ctx context.Context, existing *lifecycles.Lifecycle, plan *schemas.SpaceDefaultLifecycleResourceModel) diag.Diagnostics {
	toUpdate, diags := expandSpaceDefaultLifecycle(ctx, existing, plan)
	if diags.HasError() {
		return diags
	}

	tflog.Debug(ctx, fmt.Sprintf("updating default lifecycle '%s' of space '%s'", toUpdate.ID, plan.SpaceID.ValueString()))
	updated, err := lifecycles.Update(r.Client, toUpdate)
	if err != nil {
		diags.AddError("unable to update default lifecycle of space", err.Error())
		return diags
	}

	diags.Append(flattenSpaceDefaultLifecycle(ctx, updated, plan)...)
	return diags
}

func (r *spaceDefaultLifecycleResource) findSpaceDefaultLifecycle(spaceID string) (*lifecycles.Lifecycle, error) {
	all, err := lifecycles.GetAll(r.Client, spaceID)
	if err != nil {
		return nil, err
	}
	for _, l := range all {
		if l.Name == schemas.SpaceDefaultLifecycleName {
			return l, nil
		}
	}

	// The default lifecycle may have been renamed. It is created together with the space, so fall back to the oldest lifecycle.
	var oldest *lifecycles.Lifecycle
	oldestNumber := 0
	for _, l := range all {
		number, ok := lifecycleIDNumber(l.GetID())
		if ok && (oldest == nil || number < oldestNumber) {
			oldest, oldestNumber = l, number
		}
	}
	if oldest == nil {
		return nil, fmt.Errorf("could not find the default lifecycle in space '%s'", spaceID)
	}
	return oldest, nil
}

// lifecycleIDNumber returns the numeric part of a lifecycle ID such as "Lifecycles-42".
func lifecycleIDNumber(id string) (int, bool) {
	index := strings.LastIndex(id, "-")
	if index < 0 {
		return 0, false
	}
	number, err := strconv.Atoi(id[index+1:])
	if err != nil {
		return 0, false
	}
	return number, true
}

func isNotFound(err error) bool {
	if apiError, ok := err.(*core.APIError); ok {
		return apiError.StatusCode == http.StatusNotFound
	}
	return false
}

func expandSpaceDefaultLifecycle(ctx context.Context, existing *lifecycles.Lifecycle, model *schemas.SpaceDefaultLifecycleResourceModel) (*lifecycles.Lifecycle, diag.Diagnostics) {
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

// flattenSpaceDefaultLifecycle writes the lifecycle into model. The values already in model (the plan or prior state) are used
// to keep optional values null when they were not configured and the server returned the implicit default.
func flattenSpaceDefaultLifecycle(ctx context.Context, lifecycle *lifecycles.Lifecycle, model *schemas.SpaceDefaultLifecycleResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.ID = types.StringValue(lifecycle.GetID())
	if lifecycle.SpaceID != "" {
		model.SpaceID = types.StringValue(lifecycle.SpaceID)
	}
	model.Name = types.StringValue(lifecycle.Name)
	model.Description = types.StringValue(lifecycle.Description)
	model.ReleaseRetention = flattenSpaceDefaultLifecycleRetention(ctx, lifecycle.ReleaseRetentionPolicy, model.ReleaseRetention, true, &diags)
	model.TentacleRetention = flattenSpaceDefaultLifecycleRetention(ctx, lifecycle.TentacleRetentionPolicy, model.TentacleRetention, true, &diags)

	if len(lifecycle.Phases) == 0 && model.Phases.IsNull() {
		model.Phases = types.ListNull(spaceDefaultLifecyclePhaseObjectType)
		return diags
	}

	var priorPhases []schemas.SpaceDefaultLifecyclePhaseModel
	if !model.Phases.IsNull() && !model.Phases.IsUnknown() {
		diags.Append(model.Phases.ElementsAs(ctx, &priorPhases, false)...)
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

	return diags
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

func isKnown(value types.String) bool {
	return !value.IsNull() && !value.IsUnknown()
}
