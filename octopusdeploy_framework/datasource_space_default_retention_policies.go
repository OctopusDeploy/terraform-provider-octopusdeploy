package octopusdeploy_framework

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/retention"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type spaceDefaultRetentionPoliciesDataSource struct {
	*Config
}

func NewSpaceDefaultRetentionPoliciesDataSource() datasource.DataSource {
	return &spaceDefaultRetentionPoliciesDataSource{}
}

// Metadata implements datasource.DataSource.
func (s *spaceDefaultRetentionPoliciesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("space_default_retention_policy")
}

func (s *spaceDefaultRetentionPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	s.Config = DataSourceConfiguration(req, resp)
}

// Read implements datasource.DataSource.
func (s *spaceDefaultRetentionPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data schemas.SpaceDefaultRetentionPoliciesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := retention.SpaceDefaultRetentionPolicyQuery{
		SpaceID:       data.SpaceID.ValueString(),
		RetentionType: retention.RetentionType(data.RetentionType.ValueString()),
	}

	util.DatasourceReading(ctx, "space_default_retention_policy", query)
	existingPolicy, err := retention.Get(s.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load space default retention policies", err.Error())
		return
	}

	util.DatasourceResultCount(ctx, "space_default_retention_policy", 1)
	retentionPolicy, diag := types.ObjectValue(schemas.GetSpaceDefaultRetentionPolicyAttributes(), map[string]attr.Value{
		"id":               types.StringValue(existingPolicy.GetID()),
		"space_id":         types.StringValue(existingPolicy.GetSpaceID()),
		"retention_type":   types.StringValue(string(existingPolicy.RetentionType)),
		"strategy":         types.StringValue(string(existingPolicy.Strategy)),
		"quantity_to_keep": types.Int64Value(int64(existingPolicy.QuantityToKeep)),
		"unit":             types.StringValue(string(existingPolicy.Unit)),
	})
	data.RetentionPolicy = retentionPolicy
	data.ID = types.StringValue("Space Default Retention Policy " + time.Now().UTC().String())
	resp.Diagnostics.Append(diag...)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// Schema implements datasource.DataSource.
func (s *spaceDefaultRetentionPoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.SpaceDefaultRetentionPolicySchema{}.GetDatasourceSchema()
}
