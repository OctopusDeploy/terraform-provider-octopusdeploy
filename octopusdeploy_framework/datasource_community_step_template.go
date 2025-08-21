package octopusdeploy_framework

import (
	"context"
	"strings"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type communityStepTemplateDataSource struct {
	*Config
}

func NewCommunityStepTemplateDataSource() datasource.DataSource {
	return &communityStepTemplateDataSource{}
}

// Metadata defines the name of the data source
func (*communityStepTemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("community_step_template")
}

// Schema defines the schema of the data source
func (*communityStepTemplateDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.CommunityStepTemplateSchema{}.GetDatasourceSchema()
}

// The Configure function gives you access to a client used to interact with the Octopus Deploy API.
func (d *communityStepTemplateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = DataSourceConfiguration(req, resp)
}

// Read access the Octopus Deploy API to retrieve community step templates based on the provided configuration.
func (d *communityStepTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data schemas.CommunityStepTemplateTypeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := struct {
		ID      string
		Website string
		Name    string
	}{data.ID.ValueString(), data.Website.ValueString(), data.Name.ValueString()}

	util.DatasourceReading(ctx, "community_step_templates", query)

	communityStepTemplates, err := d.getCommunityStepTemplate(query.ID)

	if err != nil {
		resp.Diagnostics.AddError("Unable to query community step templates", err.Error())
		return
	}

	matchingCommunityStepTemplates := []*actions.CommunityActionTemplate{}

	for _, communityStepTemplate := range communityStepTemplates {
		if strings.TrimSpace(query.ID) != "" && communityStepTemplate.ID != query.ID {
			continue
		}
		if strings.TrimSpace(query.Website) != "" && communityStepTemplate.Website != query.Website {
			continue
		}
		if strings.TrimSpace(query.Name) != "" && communityStepTemplate.Name != data.Name.ValueString() {
			continue
		}
		matchingCommunityStepTemplates = append(matchingCommunityStepTemplates, communityStepTemplate)
	}

	util.DatasourceResultCount(ctx, "community_step_templates", len(matchingCommunityStepTemplates))

	data.ID = types.StringValue("CommunityActionTemplates " + time.Now().UTC().String())
	steps, stepsDiag := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.CommunityStepTemplateTypeObjectType()}, matchingCommunityStepTemplates)
	resp.Diagnostics.Append(stepsDiag...)
	data.Steps = steps

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *communityStepTemplateDataSource) getCommunityStepTemplate(id string) ([]*actions.CommunityActionTemplate, error) {
	queryIds := []string{}

	if strings.TrimSpace(id) != "" {
		queryIds = append(queryIds, id)
	}

	// Optimise the query by only requesting the IDs that are needed.
	// This does filtering server side.
	if len(queryIds) != 0 {
		return d.Config.Client.CommunityActionTemplates.GetByIDs(queryIds)
	}

	// Otherwise we have to filter client side.
	return d.Config.Client.CommunityActionTemplates.GetAll()
}

func mapCommunityStepTemplateToDatasourceModel(data *schemas.CommunityStepTemplateTypeDataSourceModel, at *actions.CommunityActionTemplate) diag.Diagnostics {
	resp := diag.Diagnostics{}

	data.ID = types.StringValue(at.ID)
	data.SpaceID = types.StringValue(at.SpaceID)
	stepTemplate, dg := convertStepTemplateAttributes(at)
	resp.Append(dg...)
	data.StepTemplate = stepTemplate
	return resp
}
