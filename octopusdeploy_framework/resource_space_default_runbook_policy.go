package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/retention"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type spaceDefaultRunbookRetentionPolicyResource struct {
	*Config
}

// NewSpaceDefaultRunbookRetentionPolicyResource creates a new resource for space default lifecycle tentacle retention policies.
func NewSpaceDefaultLifecycleRunbookPolicyResource() resource.Resource {
	return &spaceDefaultRunbookRetentionPolicyResource{}
}

var _ resource.Resource = &spaceDefaultRunbookRetentionPolicyResource{}

func (s *spaceDefaultRunbookRetentionPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.SpaceDefaultRunbookRetentionPolicySchema{}.GetResourceSchema()
}

func (s *spaceDefaultRunbookRetentionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	s.Config = ResourceConfiguration(req, resp)
}

// Metadata implements resource.Resource.
func (s *spaceDefaultRunbookRetentionPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("space_default_lifecycle_tentacle_retention_policy")
}

// We cannot create via the API; they are created automatically when a space is created and deleted when a space is deleted.
// This create is a reads and updates.
func (s *spaceDefaultRunbookRetentionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.SpaceDefaultRunbookRetentionPoliciesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	util.Create(ctx, resourceDescription, data)

	// Read existing policy
	query := retention.SpaceDefaultRetentionPolicyQuery{
		SpaceID:       data.SpaceID.ValueString(),
		RetentionType: retention.RunbookRetentionType,
	}
	existingPolicy, err := retention.Get(s.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read existing retention policy", err.Error())
		return
	}
	var newPolicy retention.ISpaceDefaultRetentionPolicy

	switch strategy := data.Strategy.ValueString(); {
	case strategy == "Forever":
		newPolicy = retention.NewKeepForeverRunbookRetentionPolicy(data.SpaceID.ValueString(), existingPolicy.GetID())
	case strategy == "Count":
		newPolicy = retention.NewCountBasedRunbookRetentionPolicy(int(data.QuantityToKeep.ValueInt64()), data.Unit.ValueString(), data.SpaceID.ValueString(), existingPolicy.GetID())
	default:
		resp.Diagnostics.AddError("Invalid strategy", "The strategy must be either 'Forever' or 'Count'.")
		return
	}

	updatedPolicy, err := retention.Update(s.Client, newPolicy)

	if err != nil {
		resp.Diagnostics.AddError("Failed to update retention policy", err.Error())
		return
	}

	updateRunbookPolicyModelFromResource(&data, updatedPolicy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements resource.Resource.
func (s *spaceDefaultRunbookRetentionPolicyResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	// Deletion is not supported so is no-op.
}

// Read implements resource.Resource.
func (s *spaceDefaultRunbookRetentionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.SpaceDefaultRunbookRetentionPoliciesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	query := retention.SpaceDefaultRetentionPolicyQuery{
		SpaceID:       data.SpaceID.ValueString(),
		RetentionType: retention.RunbookRetentionType,
	}
	policy, err := retention.Get(s.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read retention policy", err.Error())
		return
	}

	updateRunbookPolicyModelFromResource(&data, policy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.Resource.
func (s *spaceDefaultRunbookRetentionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state schemas.SpaceDefaultRunbookRetentionPoliciesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var newPolicy retention.ISpaceDefaultRetentionPolicy

	switch strategy := data.Strategy.ValueString(); {
	case strategy == "Forever":
		newPolicy = retention.NewKeepForeverRunbookRetentionPolicy(data.SpaceID.ValueString(), state.ID.ValueString())
	case strategy == "Count":
		newPolicy = retention.NewCountBasedRunbookRetentionPolicy(int(data.QuantityToKeep.ValueInt64()), data.Unit.ValueString(), data.SpaceID.ValueString(), state.ID.ValueString())
	default:
		resp.Diagnostics.AddError("Invalid strategy", "The strategy must be either 'Forever' or 'Count'.")
		return
	}

	updatedPolicy, err := retention.Update(s.Client, newPolicy)

	if err != nil {
		resp.Diagnostics.AddError("Failed to update retention policy", err.Error())
		return
	}
	updateRunbookPolicyModelFromResource(&data, updatedPolicy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func updateRunbookPolicyModelFromResource(data *schemas.SpaceDefaultRunbookRetentionPoliciesResourceModel, resource *retention.SpaceDefaultRetentionPolicyResource) {
	data.ID = types.StringValue(resource.GetID())
	data.Strategy = types.StringValue(resource.Strategy)
	data.QuantityToKeep = util.Ternary(resource.QuantityToKeep == 0, types.Int64Null(), types.Int64Value(int64(resource.QuantityToKeep)))
	data.Unit = util.Ternary(resource.Unit == "", types.StringNull(), types.StringValue(resource.Unit))
}
