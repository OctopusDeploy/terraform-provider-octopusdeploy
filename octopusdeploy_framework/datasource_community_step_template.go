package octopusdeploy_framework

import (
	"context"
	"strings"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func (d *stepTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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

	communityStepTemplates, err := d.Config.Client.CommunityActionTemplates.GetAll()

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
		if strings.TrimSpace(query.Name) != "" && communityStepTemplate.Name != query.Name {
			continue
		}
		matchingCommunityStepTemplates = append(matchingCommunityStepTemplates, communityStepTemplate)
	}

	flattenedSteps := []interface{}{}
	for _, step := range matchingCommunityStepTemplates {
		flattenedSteps = append(flattenedSteps, schemas.FlattenTenant(step))
	}

	util.DatasourceResultCount(ctx, "community_step_templates", len(flattenedSteps))

	data.ID = types.StringValue("Tenants " + time.Now().UTC().String())
	data.Steps, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.TenantObjectType()}, flattenedSteps)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
