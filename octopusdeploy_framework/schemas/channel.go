package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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

func (c ChannelSchema) GetResourceSchema() schema.Schema {
	return schema.Schema{
		Description: util.GetResourceSchemaDescription(ChannelResourceDescription),
		Attributes: map[string]schema.Attribute{
			"id":          GetIdResourceSchema(),
			"description": GetDescriptionResourceSchema(ChannelResourceDescription),
			"ephemeral_environment_name_template": schema.StringAttribute{
				Description: "The name template for ephemeral environments created from this channel.",
				Optional:    true,
			},
			"is_default": schema.BoolAttribute{
				Description: "Indicates whether this is the default channel for the associated project.",
				Optional:    true,
			},
			"lifecycle_id": schema.StringAttribute{
				Description: "The lifecycle ID associated with this channel.",
				Optional:    true,
			},
			"name": GetNameResourceSchema(true),
			"parent_environment_id": schema.StringAttribute{
				Description: "The parent environment ID for ephemeral environments.",
				Optional:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID associated with this channel.",
				Required:    true,
			},
			"space_id": GetSpaceIdResourceSchema(ChannelResourceDescription),
			"tenant_tags": schema.SetAttribute{
				Description: "A set of tenant tags associated with this channel.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"type": schema.StringAttribute{
				Description: "The type of channel. Valid values are `\"Branch\"` or `\"Tag\"`. Defaults to `\"Lifecycle\"`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Lifecycle"),
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				Description: "A list of rules associated with this channel.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID associated with this channel rule.",
							Computed:    true,
							Optional:    true,
						},
						"tag": schema.StringAttribute{
							Optional: true,
						},
						"version_range": schema.StringAttribute{
							Optional: true,
						},
					},
					Blocks: map[string]schema.Block{
						"action_package": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"deployment_action": schema.StringAttribute{
										Optional: true,
									},
									"package_reference": schema.StringAttribute{
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
