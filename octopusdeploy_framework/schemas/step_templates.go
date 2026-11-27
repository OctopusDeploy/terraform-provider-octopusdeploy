package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const StepTemplatesDatasourceDescription = "step_templates"

type StepTemplatesDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	SpaceID       types.String `tfsdk:"space_id"`
	IDs           types.List   `tfsdk:"ids"`
	PartialName   types.String `tfsdk:"partial_name"`
	Skip          types.Int64  `tfsdk:"skip"`
	Take          types.Int64  `tfsdk:"take"`
	StepTemplates types.List   `tfsdk:"step_templates"`
}

type StepTemplatesSchema struct{}

var _ EntitySchema = StepTemplatesSchema{}

func (s StepTemplatesSchema) GetResourceSchema() rs.Schema {
	return rs.Schema{}
}

func (s StepTemplatesSchema) GetDatasourceSchema() ds.Schema {
	return ds.Schema{
		Description: util.GetDataSourceDescription(StepTemplatesDatasourceDescription),
		Attributes: map[string]ds.Attribute{
			"ids":          GetQueryIDsDatasourceSchema(),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take": ds.Int64Attribute{
				Description: "A filter to limit how many step templates are returned. Every match is returned when this is omitted. This bounds the number of results, not the number of requests.",
				Optional:    true,
			},
			"space_id": GetSpaceIdDatasourceSchema(StepTemplatesDatasourceDescription, false),

			// response
			"id": GetIdDatasourceSchema(true),
			"step_templates": util.DataSourceList(StepTemplateObjectType()).
				Description("The step templates matching the query. Use the singular octopusdeploy_step_template data source to look up one template by name or ID.").
				Computed().
				Build(),
		},
	}
}

func StepTemplateObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: GetStepTemplateAttributes(),
	}
}
