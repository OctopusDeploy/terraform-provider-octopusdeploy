package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	CommunityStepTemplateResourceType          = "community step template"
	CommunityStepTemplateResourceDescription   = "community_step_template"
	CommunityStepTemplateDatasourceDescription = "community_step_template"
)

type CommunityStepTemplateTypeDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Website types.String `tfsdk:"website"`
	Name    types.String `tfsdk:"name"`
	Steps   types.List   `tfsdk:"steps"`
}

// CommunityStepTemplateTypeResourceModel represents the resource model for a community step template.
// It is a little different to most other resources because a community step template is read only and
// installed rather than created.
type CommunityStepTemplateTypeResourceModel struct {
	SpaceID       types.String `tfsdk:"space_id"`
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
			"space_id":    GetSpaceIdResourceSchema(CommunityStepTemplateResourceDescription),
			"id":          GetIdResourceSchema(),
			"name":        GetNameResourceSchema(true),
			"description": GetDescriptionResourceSchema(CommunityStepTemplateResourceDescription),
			"type": rs.StringAttribute{
				Description: "The type of the community step template",
				Computed:    true,
				Optional:    false,
			},
			"author": rs.StringAttribute{
				Description: "The author of the community step template",
				Computed:    true,
				Optional:    false,
			},
			"website": rs.StringAttribute{
				Description: "The website link to the community step template",
				Computed:    true,
				Optional:    false,
			},
			"history_url": rs.StringAttribute{
				Description: "The website link to the history community step template",
				Computed:    true,
				Optional:    false,
			},
			"version": rs.Int32Attribute{
				Description: "The version of the step template",
				Optional:    false,
				Computed:    true,
			},
			"step_package_id": rs.StringAttribute{
				Description: "The ID of the step package",
				Required:    true,
			},
			"packages":   GetStepTemplatePackageResourceSchema(CommunityStepTemplateResourceType),
			"parameters": GetStepTemplateParameterResourceSchema(CommunityStepTemplateResourceType),
			"properties": rs.MapAttribute{
				Description: "Properties for the community step template",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}
