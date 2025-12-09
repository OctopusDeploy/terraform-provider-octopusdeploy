package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/retention"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type spaceDefaultRunbookRetentionPoliciesDataSource struct {
	*Config
}

func NewSpaceDefaultRunbookRetentionPoliciesDataSource() datasource.DataSource {
	return &spaceDefaultRunbookRetentionPoliciesDataSource{}
}

// Metadata implements datasource.DataSource.
func (s *spaceDefaultRunbookRetentionPoliciesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("space_default_runbook_retention_policy")
}

func (s *spaceDefaultRunbookRetentionPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	s.Config = DataSourceConfiguration(req, resp)
}

// Read implements datasource.DataSource.
func (s *spaceDefaultRunbookRetentionPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data schemas.SpaceDefaultRunbookRetentionPoliciesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := retention.SpaceDefaultRetentionPolicyQuery{
		SpaceID:       data.SpaceID.ValueString(),
		RetentionType: retention.RunbookRetentionType,
	}

	util.DatasourceReading(ctx, "space_default_runbook_retention_policy", query)
	existingPolicy, err := retention.Get(s.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load space default runbook retention policies", err.Error())
		return
	}

	updateRunbookPolicyDatasourceModelFromResource(&data, existingPolicy)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// Schema implements datasource.DataSource.
func (s *spaceDefaultRunbookRetentionPoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.SpaceDefaultRunbookRetentionPolicySchema{}.GetDatasourceSchema()
}

func updateRunbookPolicyDatasourceModelFromResource(data *schemas.SpaceDefaultRunbookRetentionPoliciesDataSourceModel, resource *retention.SpaceDefaultRetentionPolicyResource) {
	data.ID = types.StringValue(resource.GetID())
	data.Strategy = types.StringValue(resource.Strategy)
	data.QuantityToKeep = util.Ternary(resource.QuantityToKeep == 0, types.Int64Null(), types.Int64Value(int64(resource.QuantityToKeep)))
	data.Unit = util.Ternary(resource.Unit == "", types.StringNull(), types.StringValue(resource.Unit))
}
