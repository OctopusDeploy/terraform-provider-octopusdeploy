package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &scopedUserRoleResource{}
var _ resource.ResourceWithImportState = &scopedUserRoleResource{}

type scopedUserRoleResource struct {
	*Config
}

func NewScopedUserRoleResource() resource.Resource {
	return &scopedUserRoleResource{}
}

func (r *scopedUserRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ScopedUserRoleResourceName)
}

func (r *scopedUserRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ScopedUserRoleSchema{}.GetResourceSchema()
}

func (r *scopedUserRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *scopedUserRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *scopedUserRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.ScopedUserRoleResourceModel
	var diags diag.Diagnostics

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scopedUserRole, diags := schemas.MapFromStateToScopedUserRole(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdScopedUserRole, err := r.Config.Client.ScopedUserRoles.Add(scopedUserRole)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create scoped user role", err.Error())
		return
	}

	diags = data.RefreshFromApiResponse(ctx, createdScopedUserRole)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *scopedUserRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.ScopedUserRoleResourceModel
	var diags diag.Diagnostics

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scopedUserRole, err := r.Config.Client.ScopedUserRoles.GetByID(data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "scoped user role"); err != nil {
			resp.Diagnostics.AddError("unable to load scoped user role", err.Error())
		}
		return
	}

	diags = data.RefreshFromApiResponse(ctx, scopedUserRole)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *scopedUserRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state schemas.ScopedUserRoleResourceModel
	var diags diag.Diagnostics

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scopedUserRole, diags := schemas.MapFromStateToScopedUserRole(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	scopedUserRole.ID = state.ID.ValueString()

	if scopedUserRole.SpaceID == "" && !state.SpaceID.IsNull() {
		scopedUserRole.SpaceID = state.SpaceID.ValueString()
	}

	updatedScopedUserRole, err := r.Config.Client.ScopedUserRoles.Update(scopedUserRole)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update scoped user role", err.Error())
		return
	}

	diags = data.RefreshFromApiResponse(ctx, updatedScopedUserRole)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *scopedUserRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.ScopedUserRoleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Config.Client.ScopedUserRoles.DeleteByID(data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete scoped user role", err.Error())
		return
	}
}
