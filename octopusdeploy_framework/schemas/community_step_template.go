package schemas

import (
	"fmt"
	"regexp"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	CommunityStepTemplateResourceType          = "community step template"
	CommunityStepTemplateResourceDescription   = "community_step_template"
	CommunityStepTemplateDatasourceDescription = "community_step_template"
)

type CommunityStepTemplateTypeDataSourceModel struct {
	ID      types.String                             `tfsdk:"id"`
	IDs     types.List                               `tfsdk:"ids"`
	Website types.String                             `tfsdk:"website"`
	Name    types.String                             `tfsdk:"name"`
	Steps   []CommunityStepTemplateTypeResourceModel `tfsdk:"steps"`
}

// CommunityStepTemplateTypeResourceModel represents the resource model for a community step template.
// It is a little different to most other resources because a community step template is read only and
// installed rather than created.
type CommunityStepTemplateTypeResourceModel struct {
	Type          types.String `tfsdk:"type"`
	Author        types.String `tfsdk:"author"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Packages      types.List   `tfsdk:"packages"`
	Website       types.String `tfsdk:"website"`
	HistoryUrl    types.String `tfsdk:"history_url"`
	Parameters    types.List   `tfsdk:"parameters"`
	Properties    types.Map    `tfsdk:"properties"`
	StepPackageId types.String `tfsdk:"step_package_id"`
	Version       types.Int32  `tfsdk:"version"`

	ResourceModel
}

type StepTemplateFromCommunityStepTemplateTypeResourceModel struct {
	ActionType                types.String `tfsdk:"action_type"`
	SpaceID                   types.String `tfsdk:"space_id"`
	CommunityActionTemplateId types.String `tfsdk:"community_action_template_id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	Packages                  types.List   `tfsdk:"packages"`
	Parameters                types.List   `tfsdk:"parameters"`
	Properties                types.Map    `tfsdk:"properties"`
	Version                   types.Int32  `tfsdk:"version"`

	ResourceModel
}

type CommunityStepTemplateSchema struct{}

var _ EntitySchema = CommunityStepTemplateSchema{}

func (s CommunityStepTemplateSchema) GetDatasourceSchema() ds.Schema {
	return ds.Schema{
		Description: util.GetDataSourceDescription(CommunityStepTemplateDatasourceDescription),
		Attributes: map[string]ds.Attribute{
			"id": ds.StringAttribute{
				Description: "Unique identifier of the community step template",
				Required:    true,
			},
			"space_id": ds.StringAttribute{
				Description: "SpaceID of the Community Step Template",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (s CommunityStepTemplateSchema) GetResourceSchema() rs.Schema {
	return rs.Schema{
		Description: util.GetResourceSchemaDescription(CommunityStepTemplateResourceDescription),
		Attributes: map[string]rs.Attribute{
			"id": GetIdResourceSchema(),
			"name": rs.StringAttribute{
				Description: "The name of the community step template.",
				Optional:    false,
				Computed:    true,
			},
			"description": rs.StringAttribute{
				Description: "The description of this " + CommunityStepTemplateResourceDescription + ".",
				Optional:    false,
				Computed:    true,
			},
			"space_id": GetSpaceIdResourceSchema(CommunityStepTemplateResourceDescription),
			"version": rs.Int32Attribute{
				Description: "The version of the step template",
				Optional:    false,
				Computed:    true,
			},
			"action_type": rs.StringAttribute{
				Description: "The action type of the step template",
				Optional:    false,
				Computed:    true,
			},
			"community_action_template_id": rs.StringAttribute{
				Description: "The ID of the community action template",
				Optional:    false,
				Required:    true,
			},
			"packages":   GetReadOnlyStepTemplatePackageResourceSchema(),
			"parameters": GetReadOnlyStepTemplateParameters(),
			"properties": rs.MapAttribute{
				Description: "Properties for the step template",
				Required:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func GetReadOnlyStepTemplateParameters() rs.ListNestedAttribute {
	return rs.ListNestedAttribute{
		Description: "List of parameters that can be used in the community step template.",
		Required:    false,
		Optional:    false,
		Computed:    true,
		NestedObject: rs.NestedAttributeObject{
			Attributes: map[string]rs.Attribute{
				"default_value": util.ResourceString().
					Description("A default value for the parameter, if applicable. This can be a hard-coded value or a variable reference.").
					Computed().
					Default(stringdefault.StaticString("")).
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
				"default_sensitive_value": util.ResourceString().
					Description("Use this attribute to set a sensitive default value for the parameter when display settings are set to 'Sensitive'").
					Optional().
					Sensitive().
					Build(),
				"display_settings": rs.MapAttribute{
					Description: "The display settings for the parameter.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"help_text": rs.StringAttribute{
					Description: "The help presented alongside the parameter input.",
					Optional:    false,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"id": rs.StringAttribute{
					Description: "The id for the property.",
					Computed:    true,
					Optional:    false,
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"), fmt.Sprintf("must be a valid UUID, unique within this list. Here is one you could use: %s.\nExpect uuid", uuid.New())),
					},
				},
				"label": rs.StringAttribute{
					Description: "The label shown beside the parameter when presented in the deployment process. Example: `Server name`.",
					Computed:    true,
					Optional:    false,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"name": rs.StringAttribute{
					Description: "The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods. Example: `ServerName`",
					Computed:    true,
					Optional:    false,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
			},
		},
	}
}

func GetReadOnlyStepTemplatePackageResourceSchema() rs.ListNestedAttribute {
	return rs.ListNestedAttribute{
		Description: "Package information for the community step template",
		Optional:    false,
		Computed:    true,
		NestedObject: rs.NestedAttributeObject{
			Attributes: map[string]rs.Attribute{
				"acquisition_location": rs.StringAttribute{
					Description: "Acquisition location for the package.",
					Default:     stringdefault.StaticString("Server"),
					Optional:    false,
					Computed:    true,
				},
				"feed_id": util.ResourceString().
					Description("ID of the feed.").
					Computed().
					Build(),
				"id": GetIdResourceSchema(),
				"name": util.ResourceString().
					Description("Package name.").
					Computed().
					Build(),
				"package_id": util.ResourceString().
					Description("The ID of the package to use.").
					Computed().
					Build(),
				"properties": rs.SingleNestedAttribute{
					Description: "Properties for the package.",
					Optional:    false,
					Computed:    true,
					Attributes: map[string]rs.Attribute{
						"extract": rs.StringAttribute{
							Description: "If the package should extract.",
							Default:     stringdefault.StaticString("True"),
							Optional:    false,
							Computed:    true,
						},
						"package_parameter_name": rs.StringAttribute{
							Description: "The name of the package parameter",
							Default:     stringdefault.StaticString(""),
							Optional:    false,
							Computed:    true,
						},
						"purpose": rs.StringAttribute{
							Description: "The purpose of this property.",
							Default:     stringdefault.StaticString(""),
							Optional:    false,
							Computed:    true,
						},
						"selection_mode": rs.StringAttribute{
							Description: "The selection mode.",
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
