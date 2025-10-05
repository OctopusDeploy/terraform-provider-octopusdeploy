package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ChannelResourceDescription = "channel"

type ChannelSchema struct{}

type ChannelModel struct {
	Description                      types.String `tfsdk:"description"`
	EphemeralEnvironmentNameTemplate types.String `tfsdk:"ephemeral_environment_name_template"`
	IsDefault                        types.Bool   `tfsdk:"is_default"`
	LifecycleId                      types.String `tfsdk:"lifecycle_id"`
	Name                             types.String `tfsdk:"name"`
	ParentEnvironmentID              types.String `tfsdk:"parent_environment_id"`
	ProjectId                        types.String `tfsdk:"project_id"`
	Rule                             types.List   `tfsdk:"rule"`
	SpaceId                          types.String `tfsdk:"space_id"`
	TenantTags                       types.Set    `tfsdk:"tenant_tags"`
	Type                             types.String `tfsdk:"type"`

	ResourceModel
}

func (c ChannelSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(ChannelResourceDescription),
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"description": GetDescriptionResourceSchema(ChannelResourceDescription),
			"ephemeral_environment_name_template": resourceSchema.StringAttribute{
				Description: "The name template for ephemeral environments created from this channel.",
				Optional:    true,
			},
			"is_default": resourceSchema.BoolAttribute{
				Description: "Indicates whether this is the default channel for the associated project.",
				Optional:    true,
			},
			"lifecycle_id": resourceSchema.StringAttribute{
				Description: "The lifecycle ID associated with this channel.",
				Optional:    true,
			},
			"name": GetNameResourceSchema(true),
			"parent_environment_id": resourceSchema.StringAttribute{
				Description: "The parent environment ID for ephemeral environments.",
				Optional:    true,
			},
			"project_id": resourceSchema.StringAttribute{
				Description: "The project ID associated with this channel.",
				Required:    true,
			},
			"space_id": GetSpaceIdResourceSchema(ChannelResourceDescription),
			"tenant_tags": resourceSchema.SetAttribute{
				Description: "A set of tenant tags associated with this channel.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"type": resourceSchema.StringAttribute{
				Description: "The type of channel. Valid values are `\"Lifecycle\"` or `\"EphemeralEnvironment\"`. Defaults to `\"Lifecycle\"`.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]resourceSchema.Block{
			"rule": resourceSchema.ListNestedBlock{
				Description: "A list of rules associated with this channel.",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"id": resourceSchema.StringAttribute{
							Description: "The ID associated with this channel rule.",
							Computed:    true,
							Optional:    true,
						},
						"tag": resourceSchema.StringAttribute{
							Optional: true,
						},
						"version_range": resourceSchema.StringAttribute{
							Optional: true,
						},
					},
					Blocks: map[string]resourceSchema.Block{
						"action_package": resourceSchema.ListNestedBlock{
							NestedObject: resourceSchema.NestedBlockObject{
								Attributes: map[string]resourceSchema.Attribute{
									"deployment_action": resourceSchema.StringAttribute{
										Optional: true,
									},
									"package_reference": resourceSchema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
