package schemas

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/userroles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ScopedUserRoleResourceDescription = "scoped user role"
	ScopedUserRoleResourceName        = "scoped_user_role"
)

var ScopedUserRoleSchemaAttributeNames = struct {
	ID              string
	SpaceID         string
	TeamID          string
	UserRoleID      string
	EnvironmentIDs  string
	ProjectIDs      string
	ProjectGroupIDs string
	TenantIDs       string
}{
	ID:              "id",
	SpaceID:         "space_id",
	TeamID:          "team_id",
	UserRoleID:      "user_role_id",
	EnvironmentIDs:  "environment_ids",
	ProjectIDs:      "project_ids",
	ProjectGroupIDs: "project_group_ids",
	TenantIDs:       "tenant_ids",
}

type ScopedUserRoleSchema struct{}

var _ EntitySchema = ScopedUserRoleSchema{}

func (s ScopedUserRoleSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(ScopedUserRoleResourceDescription),
		Attributes: map[string]resourceSchema.Attribute{
			ScopedUserRoleSchemaAttributeNames.ID: GetIdResourceSchema(),
			ScopedUserRoleSchemaAttributeNames.SpaceID: resourceSchema.StringAttribute{
				Description: "The space ID associated with this scoped user role. If not provided, the scoped user role will be created at the system level.",
				Optional:    true,
				Computed:    true,
			},
			ScopedUserRoleSchemaAttributeNames.TeamID: resourceSchema.StringAttribute{
				Description: "The team ID that this scoped user role belongs to.",
				Required:    true,
			},
			ScopedUserRoleSchemaAttributeNames.UserRoleID: resourceSchema.StringAttribute{
				Description: "The user role ID that defines the permissions for this scoped user role.",
				Required:    true,
			},
			ScopedUserRoleSchemaAttributeNames.EnvironmentIDs: resourceSchema.SetAttribute{
				Description: "A list of environment IDs that scope the user role. If not provided, the user role applies to all environments.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			ScopedUserRoleSchemaAttributeNames.ProjectIDs: resourceSchema.SetAttribute{
				Description: "A list of project IDs that scope the user role. If not provided, the user role applies to all projects.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			ScopedUserRoleSchemaAttributeNames.ProjectGroupIDs: resourceSchema.SetAttribute{
				Description: "A list of project group IDs that scope the user role. If not provided, the user role applies to all project groups.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			ScopedUserRoleSchemaAttributeNames.TenantIDs: resourceSchema.SetAttribute{
				Description: "A list of tenant IDs that scope the user role. If not provided, the user role applies to all tenants.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (s ScopedUserRoleSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing scoped user roles.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":                util.DataSourceString().Computed().Description("The ID of the scoped user role.").Build(),
			"space_id":          util.DataSourceString().Optional().Description("The space ID associated with this scoped user role.").Build(),
			"ids":               util.DataSourceList(types.StringType).Optional().Description("A list of scoped user role IDs to filter by.").Build(),
			"partial_name":      util.DataSourceString().Optional().Description("A partial name to filter scoped user roles by.").Build(),
			"skip":              util.DataSourceInt64().Optional().Description("A filter to specify the number of items to skip in the response.").Build(),
			"take":              util.DataSourceInt64().Optional().Description("A filter to specify the number of items to take (or return) in the response.").Build(),
			"scoped_user_roles": getScopedUserRolesAttribute(),
		},
	}
}

func getScopedUserRolesAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		Optional: false,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"id":                util.DataSourceString().Computed().Description("The ID of the scoped user role.").Build(),
				"space_id":          util.DataSourceString().Computed().Description("The space ID associated with this scoped user role.").Build(),
				"team_id":           util.DataSourceString().Computed().Description("The team ID that this scoped user role belongs to.").Build(),
				"user_role_id":      util.DataSourceString().Computed().Description("The user role ID that defines the permissions.").Build(),
				"environment_ids":   util.DataSourceSet(types.StringType).Computed().Description("The environment IDs that scope the user role.").Build(),
				"project_ids":       util.DataSourceSet(types.StringType).Computed().Description("The project IDs that scope the user role.").Build(),
				"project_group_ids": util.DataSourceSet(types.StringType).Computed().Description("The project group IDs that scope the user role.").Build(),
				"tenant_ids":        util.DataSourceSet(types.StringType).Computed().Description("The tenant IDs that scope the user role.").Build(),
			},
		},
	}
}

type ScopedUserRoleResourceModel struct {
	SpaceID         types.String `tfsdk:"space_id"`
	TeamID          types.String `tfsdk:"team_id"`
	UserRoleID      types.String `tfsdk:"user_role_id"`
	EnvironmentIDs  types.Set    `tfsdk:"environment_ids"`
	ProjectIDs      types.Set    `tfsdk:"project_ids"`
	ProjectGroupIDs types.Set    `tfsdk:"project_group_ids"`
	TenantIDs       types.Set    `tfsdk:"tenant_ids"`

	ResourceModel
}

func MapFromStateToScopedUserRole(ctx context.Context, data *ScopedUserRoleResourceModel) (*userroles.ScopedUserRole, diag.Diagnostics) {
	var diags diag.Diagnostics
	userRoleID := data.UserRoleID.ValueString()
	scopedUserRole := userroles.NewScopedUserRole(userRoleID)

	scopedUserRole.ID = data.ID.ValueString()

	if !data.SpaceID.IsNull() && !data.SpaceID.IsUnknown() {
		scopedUserRole.SpaceID = data.SpaceID.ValueString()
	}

	scopedUserRole.TeamID = data.TeamID.ValueString()

	if environmentIDs, d := util.SetToStringArray(ctx, data.EnvironmentIDs); d.HasError() {
		diags.Append(d...)
	} else {
		scopedUserRole.EnvironmentIDs = environmentIDs
	}

	if projectIDs, d := util.SetToStringArray(ctx, data.ProjectIDs); d.HasError() {
		diags.Append(d...)
	} else {
		scopedUserRole.ProjectIDs = projectIDs
	}

	if projectGroupIDs, d := util.SetToStringArray(ctx, data.ProjectGroupIDs); d.HasError() {
		diags.Append(d...)
	} else {
		scopedUserRole.ProjectGroupIDs = projectGroupIDs
	}

	if tenantIDs, d := util.SetToStringArray(ctx, data.TenantIDs); d.HasError() {
		diags.Append(d...)
	} else {
		scopedUserRole.TenantIDs = tenantIDs
	}

	return scopedUserRole, diags
}

func (data *ScopedUserRoleResourceModel) RefreshFromApiResponse(ctx context.Context, scopedUserRole *userroles.ScopedUserRole) diag.Diagnostics {
	var diags diag.Diagnostics

	if scopedUserRole == nil {
		return diags
	}

	data.ID = types.StringValue(scopedUserRole.ID)

	if scopedUserRole.SpaceID != "" {
		data.SpaceID = types.StringValue(scopedUserRole.SpaceID)
	} else {
		data.SpaceID = types.StringNull()
	}

	data.TeamID = types.StringValue(scopedUserRole.TeamID)
	data.UserRoleID = types.StringValue(scopedUserRole.UserRoleID)

	var d diag.Diagnostics
	data.EnvironmentIDs, d = types.SetValueFrom(ctx, types.StringType, scopedUserRole.EnvironmentIDs)
	diags.Append(d...)

	data.ProjectIDs, d = types.SetValueFrom(ctx, types.StringType, scopedUserRole.ProjectIDs)
	diags.Append(d...)

	data.ProjectGroupIDs, d = types.SetValueFrom(ctx, types.StringType, scopedUserRole.ProjectGroupIDs)
	diags.Append(d...)

	data.TenantIDs, d = types.SetValueFrom(ctx, types.StringType, scopedUserRole.TenantIDs)
	diags.Append(d...)

	return diags
}
