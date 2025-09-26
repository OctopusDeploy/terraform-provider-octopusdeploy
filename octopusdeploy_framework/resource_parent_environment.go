package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments/parentenvironments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type parentEnvironmentResource struct {
	*Config
}

func NewParentEnvironmentResource() resource.Resource {
	return &parentEnvironmentResource{}
}

var _ resource.ResourceWithImportState = &parentEnvironmentResource{}

func (r *parentEnvironmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("parent_environment")
}

func (r *parentEnvironmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ParentEnvironmentSchema{}.GetResourceSchema()
}

func (r *parentEnvironmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *parentEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.ParentEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentEnvironment := expandParentEnvironment(plan)

	createdParentEnvironment, err := parentenvironments.Add(r.Config.Client, parentEnvironment)
	if err != nil {
		resp.Diagnostics.AddError("Error creating parent environment", err.Error())
		return
	}

	// Fetch the full resource from the API after creation, since Add returns only the ID
	fullParentEnvironment, err := parentenvironments.GetByID(r.Config.Client, plan.SpaceID.ValueString(), createdParentEnvironment.GetID())
	if err != nil {
		resp.Diagnostics.AddError("Error fetching parent environment after create", err.Error())
		return
	}

	state := schemas.ParentEnvironmentModel{
		ResourceModel: schemas.ResourceModel{
			ID: types.StringValue(fullParentEnvironment.GetID()),
		},
		Name:             types.StringValue(fullParentEnvironment.Name),
		SpaceID:          types.StringValue(fullParentEnvironment.SpaceID),
		Description:      types.StringValue(fullParentEnvironment.Description),
		Slug:             types.StringValue(fullParentEnvironment.Slug),
		UseGuidedFailure: types.BoolValue(fullParentEnvironment.UseGuidedFailure),
	}
	if fullParentEnvironment.AutomaticDeprovisioningRule != nil {
		state.AutomaticDeprovisioningRule = &schemas.AutomaticDeprovisioningRuleModel{
			Days:  types.Int64Value(int64(fullParentEnvironment.AutomaticDeprovisioningRule.ExpiryDays)),
			Hours: types.Int64Value(int64(fullParentEnvironment.AutomaticDeprovisioningRule.ExpiryHours)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *parentEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.ParentEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentEnvironment, err := parentenvironments.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "parentEnvironmentResource"); err != nil {
			resp.Diagnostics.AddError("unable to load parent environment", err.Error())
		}
		return
	}

	newState := schemas.ParentEnvironmentModel{
		ResourceModel: schemas.ResourceModel{
			ID: types.StringValue(parentEnvironment.GetID()),
		},
		Name:             types.StringValue(parentEnvironment.Name),
		SpaceID:          types.StringValue(parentEnvironment.SpaceID),
		Description:      types.StringValue(parentEnvironment.Description),
		Slug:             types.StringValue(parentEnvironment.Slug),
		UseGuidedFailure: types.BoolValue(parentEnvironment.UseGuidedFailure),
	}
	if parentEnvironment.AutomaticDeprovisioningRule != nil {
		newState.AutomaticDeprovisioningRule = &schemas.AutomaticDeprovisioningRuleModel{
			Days:  types.Int64Value(int64(parentEnvironment.AutomaticDeprovisioningRule.ExpiryDays)),
			Hours: types.Int64Value(int64(parentEnvironment.AutomaticDeprovisioningRule.ExpiryHours)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *parentEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.ParentEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentEnvironment := expandParentEnvironment(plan)
	updatedParentEnvironment, err := parentenvironments.Update(r.Client, parentEnvironment)
	if err != nil {
		resp.Diagnostics.AddError("Error updating parent environment", err.Error())
		return
	}

	state := flattenParentEnvironment(updatedParentEnvironment)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *parentEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.ParentEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := parentenvironments.DeleteByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting parent environment", err.Error())
		return
	}
}

func (r *parentEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper functions to convert between Framework models and API models
func expandParentEnvironment(model schemas.ParentEnvironmentModel) *parentenvironments.ParentEnvironment {
	parentEnvironment := &parentenvironments.ParentEnvironment{
		Name:        model.Name.ValueString(),
		SpaceID:     model.SpaceID.ValueString(),
		Description: model.Description.ValueString(),
		Slug:        model.Slug.ValueString(),
	}

	// Only set UseGuidedFailure if it's explicitly set in the config
	if !model.UseGuidedFailure.IsNull() && !model.UseGuidedFailure.IsUnknown() {
		parentEnvironment.UseGuidedFailure = model.UseGuidedFailure.ValueBool()
	}

	// Set ID if it exists
	if !model.ID.IsNull() && !model.ID.IsUnknown() && model.ID.ValueString() != "" {
		parentEnvironment.ID = model.ID.ValueString()
	}

	if model.AutomaticDeprovisioningRule != nil {
		expiryDaysValue := int64(0)
		expiryHoursValue := int64(0)

		if !model.AutomaticDeprovisioningRule.Days.IsNull() && !model.AutomaticDeprovisioningRule.Days.IsUnknown() {
			expiryDaysValue = model.AutomaticDeprovisioningRule.Days.ValueInt64()
		}

		if !model.AutomaticDeprovisioningRule.Hours.IsNull() && !model.AutomaticDeprovisioningRule.Hours.IsUnknown() {
			expiryHoursValue = model.AutomaticDeprovisioningRule.Hours.ValueInt64()
		}

		parentEnvironment.AutomaticDeprovisioningRule = &parentenvironments.AutomaticDeprovisioningRule{
			ExpiryDays:  int(expiryDaysValue),
			ExpiryHours: int(expiryHoursValue),
		}
	}

	return parentEnvironment
}

func flattenParentEnvironment(parentEnvironment *parentenvironments.ParentEnvironment) schemas.ParentEnvironmentModel {
	result := schemas.ParentEnvironmentModel{
		ResourceModel: schemas.ResourceModel{
			ID: types.StringValue(parentEnvironment.GetID()),
		},
		Name:             types.StringValue(parentEnvironment.Name),
		SpaceID:          types.StringValue(parentEnvironment.SpaceID),
		Description:      types.StringValue(parentEnvironment.Description),
		Slug:             types.StringValue(parentEnvironment.Slug),
		UseGuidedFailure: types.BoolValue(parentEnvironment.UseGuidedFailure),
	}

	// Handle automatic_deprovisioning_rule
	if parentEnvironment.AutomaticDeprovisioningRule != nil {
		result.AutomaticDeprovisioningRule = &schemas.AutomaticDeprovisioningRuleModel{
			Days:  types.Int64Value(int64(parentEnvironment.AutomaticDeprovisioningRule.ExpiryDays)),
			Hours: types.Int64Value(int64(parentEnvironment.AutomaticDeprovisioningRule.ExpiryHours)),
		}
	}

	return result
}
