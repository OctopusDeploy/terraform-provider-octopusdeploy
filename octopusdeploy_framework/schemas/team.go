package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const TeamResourceDescription = "team"

type TeamSchema struct{}

type TeamModel struct {
	CanBeDeleted          types.Bool   `tfsdk:"can_be_deleted"`
	CanBeRenamed          types.Bool   `tfsdk:"can_be_renamed"`
	CanChangeMembers      types.Bool   `tfsdk:"can_change_members"`
	CanChangeRoles        types.Bool   `tfsdk:"can_change_roles"`
	Description           types.String `tfsdk:"description"`
	ExternalSecurityGroup types.List   `tfsdk:"external_security_group"`
	Name                  types.String `tfsdk:"name"`
	SpaceID               types.String `tfsdk:"space_id"`
	Users                 types.Set    `tfsdk:"users"`
	UserRole              types.Set    `tfsdk:"user_role"`

	ResourceModel
}

func (t TeamSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(TeamResourceDescription),
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        GetNameResourceSchema(true),
			"description": GetDescriptionResourceSchema(TeamResourceDescription),
			"space_id":    GetSpaceIdResourceSchema(TeamResourceDescription),
			"users": resourceSchema.SetAttribute{
				Description: "A list of user IDs designated to be members of this team.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"can_be_deleted": resourceSchema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"can_be_renamed": resourceSchema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"can_change_members": resourceSchema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"can_change_roles": resourceSchema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
		},
		Blocks: map[string]resourceSchema.Block{
			"user_role": resourceSchema.SetNestedBlock{
				Description:  "User roles associated with this team.",
				NestedObject: GetUserRoleSchema(),
			},
			"external_security_group": GetSecurityGroupSchema(),
		},
	}
}

func GetUserRoleSchema() resourceSchema.NestedBlockObject {
	return resourceSchema.NestedBlockObject{
		Attributes: map[string]resourceSchema.Attribute{
			"environment_ids": resourceSchema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"id": resourceSchema.StringAttribute{
				Computed: true,
			},
			"project_group_ids": resourceSchema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"project_ids": resourceSchema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"space_id": resourceSchema.StringAttribute{
				Required: true,
			},
			"team_id": resourceSchema.StringAttribute{
				Computed: true,
			},
			"tenant_ids": resourceSchema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"user_role_id": resourceSchema.StringAttribute{
				Required: true,
			},
		},
	}
}

func GetSecurityGroupSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"display_id_and_name": resourceSchema.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				"display_name": resourceSchema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"id": resourceSchema.StringAttribute{
					Description: "The unique ID of this external security group.",
					Optional:    true,
					Computed:    true,
				},
			},
		},
	}
}
