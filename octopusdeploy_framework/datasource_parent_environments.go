package octopusdeploy_framework

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments/v2/environments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type parentEnvironmentDataSource struct {
	*Config
}

type parentEnvironmentsDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	SpaceID      types.String `tfsdk:"space_id"`
	IDs          types.List   `tfsdk:"ids"`
	PartialName  types.String `tfsdk:"partial_name"`
	Name         types.String `tfsdk:"name"`
	Skip         types.Int64  `tfsdk:"skip"`
	Take         types.Int64  `tfsdk:"take"`
	Environments types.List   `tfsdk:"environments"`
}

func NewParentEnvironmentsDataSource() datasource.DataSource {
	return &parentEnvironmentDataSource{}
}

func (*parentEnvironmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("parent_environments")
}

func (*parentEnvironmentDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.EnvironmentSchema{}.GetDatasourceSchema()
}

func (e *parentEnvironmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	e.Config = DataSourceConfiguration(req, resp)
}

func (e *parentEnvironmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data parentEnvironmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := environments.EnvironmentQuery{
		IDs:         util.GetIds(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        util.GetNumber(data.Skip),
		Take:        util.GetNumber(data.Take),
		Type:        []string{"Parent"},
	}

	util.DatasourceReading(ctx, "parent_environments", query)

	existingEnvironments, err := environments.Get(e.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load parent environments", err.Error())
		return
	}

	mappedEnvironments := []schemas.ParentEnvironmentTypeResourceModel{}
	if data.Name.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf("parent environments returned from API: %#v", existingEnvironments))
		for _, environment := range existingEnvironments.Items {
			mappedEnvironments = append(mappedEnvironments, schemas.MapFromParentEnvironment(ctx, environment))
		}
	} else { // if name has been specified, match by exact name rather than partial name as the API does
		var matchedEnvironment *environments.Environment
		tflog.Debug(ctx, fmt.Sprintf("matching parent environment by name: %s", data.Name))
		for _, env := range existingEnvironments.Items {
			if strings.EqualFold(env.Name, data.Name.ValueString()) {
				matchedEnvironment = env
			}
		}
		if matchedEnvironment != nil {
			mappedEnvironments = append(mappedEnvironments, schemas.MapFromParentEnvironment(ctx, matchedEnvironment))
		}
	}

	util.DatasourceResultCount(ctx, "parent_environments", len(mappedEnvironments))

	data.Environments, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.EnvironmentObjectType()}, mappedEnvironments)
	data.ID = types.StringValue("Parent Environments " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
