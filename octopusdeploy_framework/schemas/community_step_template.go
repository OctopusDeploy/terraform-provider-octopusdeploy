package schemas

import (
	"fmt"
	"regexp"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
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

// CommunityStepTemplateTypeDataSourceModel represents the data source defined in the Terraform configuration.
type CommunityStepTemplateTypeDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	SpaceID types.String `tfsdk:"space_id"`
	Website types.String `tfsdk:"website"`
	Name    types.String `tfsdk:"name"`
	Steps   types.List   `tfsdk:"steps"`
}

// CommunityStepTemplateTypeObjectType returns the type mapping used to define the Steps attribute in the CommunityStepTemplateTypeDataSourceModel.
func CommunityStepTemplateTypeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"author":          types.StringType,
		"name":            types.StringType,
		"description":     types.StringType,
		"website":         types.StringType,
		"history_url":     types.StringType,
		"version":         types.Int32Type,
		"step_package_id": types.StringType,
		"parameters": types.ListType{
			ElemType: types.ObjectType{AttrTypes: ParametersObjectType()},
		},
		"properties": types.MapType{ElemType: types.StringType},
		"packages": types.ListType{
			ElemType: types.ObjectType{AttrTypes: PackagesObjectType()},
		},
	}
}

// ParametersObjectType returns the type mapping used to define the parameters attribute in the CommunityStepTemplateTypeObjectType function
func ParametersObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                      types.StringType,
		"default_value":           types.StringType,
		"display_settings":        types.MapType{ElemType: types.StringType},
		"default_sensitive_value": types.StringType,
		"help_text":               types.StringType,
		"label":                   types.StringType,
		"name":                    types.StringType,
	}
}

// PackagesObjectType returns the type mapping used to define the packages attribute in the CommunityStepTemplateTypeObjectType function
func PackagesObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"acquisition_location": types.StringType,
		"feed_id":              types.StringType,
		"id":                   types.StringType,
		"name":                 types.StringType,
		"package_id":           types.StringType,
		"properties": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"extract":                types.StringType,
				"package_parameter_name": types.StringType,
				"purpose":                types.StringType,
				"selection_mode":         types.StringType,
			},
		},
	}
}

// StepTemplateFromCommunityStepTemplateTypeResourceModel represents a step template generated from a community step template.
// Notably, it does not include the git dependencies, as these are not exposed on community step templates.
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
				Optional:    true,
			},
			"space_id": ds.StringAttribute{
				Description: "SpaceID of the Community Step Template",
				Optional:    true,
				Computed:    true,
			},
			"name": ds.StringAttribute{
				Description: "Name of the Community Step Template",
				Optional:    true,
			},
			"website": ds.StringAttribute{
				Description: "Website of the Community Step Template",
				Optional:    true,
			},
			"steps": ds.ListNestedAttribute{
				Computed: true,
				Optional: false,
				NestedObject: ds.NestedAttributeObject{
					Attributes: s.GetDataSourceStepsAttributes(),
				},
			},
		},
	}
}

func (s CommunityStepTemplateSchema) GetDataSourceStepsAttributes() map[string]ds.Attribute {
	return map[string]ds.Attribute{
		"id":          GetIdDatasourceSchema(true),
		"name":        GetReadonlyNameDatasourceSchema(),
		"description": GetReadonlyDescriptionDatasourceSchema(CommunityStepTemplateResourceDescription),
		"author": ds.StringAttribute{
			Description: "The author of this " + CommunityStepTemplateResourceDescription + ".",
			Computed:    true,
		},
		"website": ds.StringAttribute{
			Description: "The website of this " + CommunityStepTemplateResourceDescription + ".",
			Computed:    true,
		},
		"history_url": ds.StringAttribute{
			Description: "The history url of this " + CommunityStepTemplateResourceDescription + ".",
			Computed:    true,
		},
		"step_package_id": ds.StringAttribute{
			Description: "The step package ID url of this " + CommunityStepTemplateResourceDescription + ".",
			Computed:    true,
		},
		"version": ds.Int32Attribute{
			Description: "The version ID url of this " + CommunityStepTemplateResourceDescription + ".",
			Computed:    true,
		},
		"packages":   GetStepTemplatePackageResourceSchema(CommunityStepTemplateResourceDescription),
		"parameters": GetStepTemplateParameterResourceSchema(CommunityStepTemplateResourceDescription),
		"properties": rs.MapAttribute{
			Description: "Properties for the community step template",
			Computed:    true,
			ElementType: types.StringType,
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": rs.StringAttribute{
				Description: "The description of this " + CommunityStepTemplateResourceDescription + ".",
				Optional:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"space_id": GetSpaceIdResourceSchema(CommunityStepTemplateResourceDescription),
			"version": rs.Int32Attribute{
				Description: "The version of the step template",
				Optional:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"action_type": rs.StringAttribute{
				Description: "The action type of the step template",
				Optional:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"community_action_template_id": rs.StringAttribute{
				Description: "The ID of the community action template",
				Optional:    false,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"packages":   GetReadOnlyStepTemplatePackageResourceSchema(),
			"parameters": GetReadOnlyStepTemplateParameters(),
			"properties": rs.MapAttribute{
				Description: "Properties for the step template",
				Required:    false,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
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
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
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
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		NestedObject: rs.NestedAttributeObject{
			Attributes: map[string]rs.Attribute{
				"acquisition_location": rs.StringAttribute{
					Description: "Acquisition location for the package.",
					Default:     stringdefault.StaticString("Server"),
					Optional:    false,
					Computed:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"feed_id": util.ResourceString().
					Description("ID of the feed.").
					Computed().
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
				"id": GetIdResourceSchema(),
				"name": util.ResourceString().
					Description("Package name.").
					Computed().
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
				"package_id": util.ResourceString().
					Description("The ID of the package to use.").
					Computed().
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
				"properties": rs.MapAttribute{
					Description: "The display settings for the parameter.",
					Optional:    true,
					ElementType: types.StringType,
				},
			},
		},
	}
}
