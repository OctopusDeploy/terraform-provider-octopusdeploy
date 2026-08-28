package octopusdeploy_framework

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stepTemplatesDataSource struct {
	*Config
}

func NewStepTemplatesDataSource() datasource.DataSource {
	return &stepTemplatesDataSource{}
}

func (*stepTemplatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("step_templates")
}

func (*stepTemplatesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.StepTemplatesSchema{}.GetDatasourceSchema()
}

func (d *stepTemplatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = DataSourceConfiguration(req, resp)
}

func (d *stepTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data schemas.StepTemplatesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := actiontemplates.Query{
		IDs:         util.GetIds(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        util.GetNumber(data.Skip),
		Take:        util.GetNumber(data.Take),
	}

	// The server pages at 30 by default, which would silently truncate the result set. Ask
	// for everything unless the configuration says otherwise; the server honours this in a
	// single request.
	if data.Take.IsNull() {
		query.Take = math.MaxInt32
	}

	util.DatasourceReading(ctx, "step_templates", query)

	actionTemplates, err := queryStepTemplates(d.Config.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("Unable to load step templates", err.Error())
		return
	}

	util.DatasourceResultCount(ctx, "step_templates", len(actionTemplates))

	stepTemplates := make([]attr.Value, len(actionTemplates))
	for i, at := range actionTemplates {
		stepTemplate, dg := convertStepTemplateAttributes(at)
		resp.Diagnostics.Append(dg...)
		stepTemplates[i] = stepTemplate
	}
	if resp.Diagnostics.HasError() {
		return
	}

	stepTemplateList, dg := types.ListValue(schemas.StepTemplateObjectType(), stepTemplates)
	resp.Diagnostics.Append(dg...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.StepTemplates = stepTemplateList
	data.ID = types.StringValue(stepTemplatesDataSourceID(data.SpaceID.ValueString(), query))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// stepTemplatesDataSourceID derives the data source ID from the query so that repeated reads
// of the same query keep the same ID, rather than producing a new one on every plan.
func stepTemplatesDataSourceID(spaceID string, query actiontemplates.Query) string {
	key := fmt.Sprintf("%s|%s|%s|%d|%d", spaceID, strings.Join(query.IDs, ","), query.PartialName, query.Skip, query.Take)
	return fmt.Sprintf("StepTemplates-%x", sha256.Sum256([]byte(key)))
}
