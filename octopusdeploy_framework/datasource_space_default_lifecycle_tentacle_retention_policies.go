package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/retention"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type spaceDefaultLifecycleTentacleRetentionPoliciesDataSource struct {
	*Config
}

func NewSpaceDefaultLifecycleTentacleRetentionPoliciesDataSource() datasource.DataSource {
	return &spaceDefaultLifecycleTentacleRetentionPoliciesDataSource{}
}

// Metadata implements datasource.DataSource.
func (s *spaceDefaultLifecycleTentacleRetentionPoliciesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("space_default_lifecycle_tentacle_retention_policy")
}

func (s *spaceDefaultLifecycleTentacleRetentionPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	s.Config = DataSourceConfiguration(req, resp)
}

// Read implements datasource.DataSource.
func (s *spaceDefaultLifecycleTentacleRetentionPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data schemas.SpaceDefaultLifecycleTentacleRetentionPoliciesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := retention.SpaceDefaultRetentionPolicyQuery{
		SpaceID:       data.SpaceID.ValueString(),
		RetentionType: retention.LifecycleTentacleRetentionType,
	}

	util.DatasourceReading(ctx, "space_default_lifecycle_tentacle_retention_policy", query)
	existingPolicy, err := retention.Get(s.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load space default lifecycle tentacle retention policies", err.Error())
		return
	}

	updateLifecycleTentaclePolicyDatasourceModelFromResource(&data, existingPolicy)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// Schema implements datasource.DataSource.
func (s *spaceDefaultLifecycleTentacleRetentionPoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.SpaceDefaultLifecycleTentacleRetentionPolicySchema{}.GetDatasourceSchema()
}

func updateLifecycleTentaclePolicyDatasourceModelFromResource(data *schemas.SpaceDefaultLifecycleTentacleRetentionPoliciesDataSourceModel, resource *retention.SpaceDefaultRetentionPolicyResource) {
	data.ID = types.StringValue(resource.GetID())
	data.Strategy = types.StringValue(resource.Strategy)
	data.QuantityToKeep = util.Ternary(resource.QuantityToKeep == 0, types.Int64Null(), types.Int64Value(int64(resource.QuantityToKeep)))
	data.Unit = util.Ternary(resource.Unit == "", types.StringNull(), types.StringValue(resource.Unit))
}
