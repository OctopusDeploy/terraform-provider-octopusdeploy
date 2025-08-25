package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	Website types.String `tfsdk:"website"`
	Name    types.String `tfsdk:"name"`
	Steps   types.List   `tfsdk:"steps"` // Steps used the type CommunityStepTemplateTypeObjectType()
}

// CommunityStepTemplateTypeObjectType returns the type mapping used to define the Steps attribute in the CommunityStepTemplateTypeDataSourceModel.
func CommunityStepTemplateTypeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"author":      types.StringType,
		"type":        types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"website":     types.StringType,
		"history_url": types.StringType,
		"version":     types.Int32Type,
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

// CommunityStepTemplateTypeResourceModel is the resource model for the community step template type. This only returned
// by the data source.
type CommunityStepTemplateTypeResourceModel struct {
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Author      types.String `tfsdk:"author"`
	Description types.String `tfsdk:"description"`
	Website     types.String `tfsdk:"website"`
	HistoryUrl  types.String `tfsdk:"history_url"`
	Version     types.Int32  `tfsdk:"version"`
	Packages    types.List   `tfsdk:"packages"`
	Parameters  types.List   `tfsdk:"parameters"`
	Properties  types.Map    `tfsdk:"properties"`

	ResourceModel
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
			"id": util.ResourceString().
				Description("Unique identifier of the community step template").
				Optional().
				Build(),
			"name": util.ResourceString().
				Description("Name of the Community Step Template").
				Optional().
				Build(),
			"website": util.ResourceString().
				Description("Website of the Community Step Template").
				Optional().
				Build(),
			"steps": ds.ListNestedAttribute{
				Computed: true,
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
		"author": util.ResourceString().
			Description("The author of this " + CommunityStepTemplateResourceDescription + ".").
			Computed().
			Build(),
		"type": util.ResourceString().
			Description("The action type of this " + CommunityStepTemplateResourceDescription + ".").
			Computed().
			Build(),
		"website": util.ResourceString().
			Description("The website of this " + CommunityStepTemplateResourceDescription + ".").
			Computed().
			Build(),
		"history_url": util.ResourceString().
			Description("The history url of this " + CommunityStepTemplateResourceDescription + ".").
			Computed().
			Build(),
		"version": util.ResourceInt32().
			Description("The version ID url of this " + CommunityStepTemplateResourceDescription + ".").
			Computed().
			Build(),
		"packages":   GetStepTemplatePackageResourceSchema(CommunityStepTemplateResourceDescription),
		"parameters": GetStepTemplateParameterResourceSchema(CommunityStepTemplateResourceDescription),
		"properties": util.ResourceMap(types.StringType).
			Description("Properties for the community step template").
			Computed().
			Build(),
	}
}

func (s CommunityStepTemplateSchema) GetResourceSchema() rs.Schema {
	return rs.Schema{
		Description: util.GetResourceSchemaDescription(CommunityStepTemplateResourceDescription),
		Attributes: map[string]rs.Attribute{
			"id": GetIdResourceSchema(),
			"name": util.ResourceString().
				Description("The name of the community step template.").
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"description": util.ResourceString().
				Description("The description of this " + CommunityStepTemplateResourceDescription + ".").
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"space_id": GetSpaceIdResourceSchema(CommunityStepTemplateResourceDescription),
			"version": util.ResourceInt32().
				Description("The version of the step template").
				Computed().
				PlanModifiers(int32planmodifier.UseStateForUnknown()).
				Build(),
			"action_type": util.ResourceString().
				Description("The action type of the step template").
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"community_action_template_id": util.ResourceString().
				Description("The ID of the community action template").
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()).
				Build(),
			"packages":   GetReadOnlyStepTemplatePackageResourceSchema(),
			"parameters": GetReadOnlyStepTemplateParameters(),
			"properties": util.ResourceMap(types.StringType).
				Description("Properties for the step template").
				Computed().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
		},
	}
}

func GetReadOnlyStepTemplateParameters() rs.ListNestedAttribute {
	return rs.ListNestedAttribute{
		Description: "List of parameters that can be used in the community step template.",
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		NestedObject: rs.NestedAttributeObject{
			Attributes: map[string]rs.Attribute{
				"default_value": util.ResourceString().
					Description("A default value for the parameter, if applicable. This can be a hard-coded value or a variable reference.").
					Computed().
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
				"default_sensitive_value": util.ResourceString().
					Description("Use this attribute to set a sensitive default value for the parameter when display settings are set to 'Sensitive'").
					Computed().
					Sensitive().
					Build(),
				"display_settings": util.ResourceMap(types.StringType).
					Description("The display settings for the parameter.").
					Computed().
					Build(),
				"help_text": util.ResourceString().
					Description("The help presented alongside the parameter input.").
					Computed().
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
				"id": util.ResourceString().
					Description("The id for the property.").
					Computed().
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
				"label": util.ResourceString().
					Description("The label shown beside the parameter when presented in the deployment process. Example: `Server name`.").
					Computed().
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
				"name": util.ResourceString().
					Description("The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods. Example: `ServerName`").
					Computed().
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
			},
		},
	}
}

func GetReadOnlyStepTemplatePackageResourceSchema() rs.ListNestedAttribute {
	return rs.ListNestedAttribute{
		Description: "Package information for the community step template",
		Computed:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		NestedObject: rs.NestedAttributeObject{
			Attributes: map[string]rs.Attribute{
				"acquisition_location": util.ResourceString().
					Description("Acquisition location for the package.").
					Computed().
					PlanModifiers(stringplanmodifier.UseStateForUnknown()).
					Build(),
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
				"properties": rs.SingleNestedAttribute{
					Description: "Properties for the package.",
					Computed:    true,
					Attributes: map[string]rs.Attribute{
						"extract": util.ResourceString().
							Description("If the package should extract.").
							Computed().
							PlanModifiers(stringplanmodifier.UseStateForUnknown()).
							Build(),
						"package_parameter_name": util.ResourceString().
							Description("The name of the package parameter").
							Computed().
							PlanModifiers(stringplanmodifier.UseStateForUnknown()).
							Build(),
						"purpose": util.ResourceString().
							Description("The purpose of this property.").
							Computed().
							PlanModifiers(stringplanmodifier.UseStateForUnknown()).
							Build(),
						"selection_mode": util.ResourceString().
							Description("The selection mode.").
							Computed().
							PlanModifiers(stringplanmodifier.UseStateForUnknown()).
							Build(),
					},
				},
			},
		},
	}
}
