package octopusdeploy_framework

import (
	"context"
	"strings"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type communityStepTemplateDataSource struct {
	*Config
}

func NewCommunityStepTemplateDataSource() datasource.DataSource {
	return &communityStepTemplateDataSource{}
}

func (*communityStepTemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("community_step_template")
}

func (*communityStepTemplateDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.CommunityStepTemplateSchema{}.GetDatasourceSchema()
}

func (d *communityStepTemplateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = DataSourceConfiguration(req, resp)
}

func (d *communityStepTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data schemas.CommunityStepTemplateTypeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ids []string
	resp.Diagnostics.Append(data.IDs.ElementsAs(ctx, &ids, false)...)

	query := struct {
		ID      string
		IDs     []string
		Website string
		Name    string
	}{data.ID.ValueString(), ids, data.Website.ValueString(), data.Name.ValueString()}

	util.DatasourceReading(ctx, "community_step_templates", query)

	communityStepTemplates, err := d.getCommunityStepTemplate(query.ID, query.IDs)

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

	data.ID = types.StringValue("Tenants " + time.Now().UTC().String())
	data.Steps = make([]schemas.CommunityStepTemplateTypeResourceModel, 0, len(matchingCommunityStepTemplates))
	for _, project := range matchingCommunityStepTemplates {
		data.Steps = append(data.Steps, flattenStep(project))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *communityStepTemplateDataSource) getCommunityStepTemplate(id string, ids []string) ([]*actions.CommunityActionTemplate, error) {
	queryIds := []string{}

	if len(ids) > 0 {
		queryIds = ids
	}

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

func flattenStep(step *actions.CommunityActionTemplate) schemas.CommunityStepTemplateTypeResourceModel {
	resourceModel := schemas.CommunityStepTemplateTypeResourceModel{
		Type:          types.StringValue(step.Type),
		Author:        types.StringValue(step.Author),
		Name:          types.StringValue(step.Name),
		Description:   types.StringValue(step.Description),
		Packages:      types.List{},
		Website:       types.StringValue(step.Website),
		HistoryUrl:    types.StringValue(step.HistoryURL),
		Parameters:    types.List{},
		Properties:    types.Map{},
		StepPackageId: types.StringValue(step.ActionType),
		Version:       types.Int32Value(step.Version),
		ResourceModel: schemas.ResourceModel{
			ID: types.StringValue(step.ID),
		},
	}

	return resourceModel
}
