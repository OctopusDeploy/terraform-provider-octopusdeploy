package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

type lifecyclesDataSource struct {
	*Config
}

type lifecyclesDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	SpaceID     types.String `tfsdk:"space_id"`
	IDs         types.List   `tfsdk:"ids"`
	PartialName types.String `tfsdk:"partial_name"`
	Skip        types.Int64  `tfsdk:"skip"`
	Take        types.Int64  `tfsdk:"take"`
	Lifecycles  types.List   `tfsdk:"lifecycles"`
}

var _ datasource.DataSource = &lifecyclesDataSource{}

func NewLifecyclesDataSource() datasource.DataSource {
	return &lifecyclesDataSource{}
}

func (l *lifecyclesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "lifecycles datasource Metadata")
	resp.TypeName = util.GetTypeName("lifecycles")
}

func (l *lifecyclesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "lifecycles datasource Schema")
	resp.Schema = schemas.LifecycleSchema{}.GetDatasourceSchema()
}

func (l *lifecyclesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "lifecycles datasource Configure")
	l.Config = DataSourceConfiguration(req, resp)
}

func (l *lifecyclesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "lifecycles datasource Read")
	var data lifecyclesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := lifecycles.Query{
		IDs:         util.ExpandStringList(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}
	util.DatasourceReading(ctx, "lifecycles", query)

	lifecyclesResult, err := lifecycles.Get(l.Config.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read lifecycles, got error: %s", err))
		return
	}

	util.DatasourceResultCount(ctx, "lifecycles", len(lifecyclesResult.Items))

	data.Lifecycles = flattenLifecyclesForDatasource(lifecyclesResult.Items)

	data.ID = types.StringValue("Lifecycles " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenLifecyclesForDatasource(requestedLifecycles []*lifecycles.Lifecycle) types.List {
	var lifecycleAttrTypes = getDatasourceLifecycleAttrTypes()
	lifecyclesList := make([]attr.Value, 0, len(requestedLifecycles))
	for _, lifecycle := range requestedLifecycles {
		lifecycleMap := map[string]attr.Value{
			"id":                               types.StringValue(lifecycle.ID),
			"space_id":                         types.StringValue(lifecycle.SpaceID),
			"name":                             types.StringValue(lifecycle.Name),
			"description":                      types.StringValue(lifecycle.Description),
			"phase":                            flattenPhasesForDataSource(lifecycle.Phases),
			"release_retention_with_strategy":  flattenRetentionForDataSource(lifecycle.ReleaseRetentionPolicy),
			"tentacle_retention_with_strategy": flattenRetentionForDataSource(lifecycle.TentacleRetentionPolicy),
		}
		lifecyclesList = append(lifecyclesList, types.ObjectValueMust(lifecycleAttrTypes, lifecycleMap))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: lifecycleAttrTypes}, lifecyclesList)
}

func flattenPhasesForDataSource(requestedPhases []*lifecycles.Phase) types.List {
	var phaseAttrTypes = getDatasourcePhaseAttrTypes()
	if requestedPhases == nil {
		return types.ListNull(types.ObjectType{AttrTypes: phaseAttrTypes})
	}
	phasesList := make([]attr.Value, 0, len(requestedPhases))

	for _, goPhase := range requestedPhases {
		attrs := map[string]attr.Value{
			"id":                                    types.StringValue(goPhase.ID),
			"name":                                  types.StringValue(goPhase.Name),
			"automatic_deployment_targets":          util.FlattenStringList(goPhase.AutomaticDeploymentTargets),
			"optional_deployment_targets":           util.FlattenStringList(goPhase.OptionalDeploymentTargets),
			"minimum_environments_before_promotion": types.Int64Value(int64(goPhase.MinimumEnvironmentsBeforePromotion)),
			"is_optional_phase":                     types.BoolValue(goPhase.IsOptionalPhase),
			"is_priority_phase":                     types.BoolValue(goPhase.IsPriorityPhase),
			"release_retention_with_strategy":       util.Ternary(goPhase.ReleaseRetentionPolicy != nil, flattenRetentionForDataSource(goPhase.ReleaseRetentionPolicy), ListNullRetention),
			"tentacle_retention_with_strategy":      util.Ternary(goPhase.TentacleRetentionPolicy != nil, flattenRetentionForDataSource(goPhase.TentacleRetentionPolicy), ListNullRetention),
		}
		phasesList = append(phasesList, types.ObjectValueMust(phaseAttrTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: phaseAttrTypes}, phasesList)
}

func flattenRetentionForDataSource(requestedRetention *core.RetentionPeriod) types.List {
	var retentionAttrTypes = getDatasourceRetentionAttrTypes()
	if requestedRetention == nil {
		return ListNullRetention
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: retentionAttrTypes},
		[]attr.Value{
			types.ObjectValueMust(
				retentionAttrTypes,
				map[string]attr.Value{
					"strategy":         types.StringValue(requestedRetention.Strategy),
					"unit":             types.StringValue(requestedRetention.Unit),
					"quantity_to_keep": types.Int64Value(int64(requestedRetention.QuantityToKeep)),
				},
			),
		},
	)
}

func getDatasourceLifecycleAttrTypes() map[string]attr.Type {
	var retentionAttrTypes = getDatasourceRetentionAttrTypes()
	return map[string]attr.Type{
		"id":                               types.StringType,
		"space_id":                         types.StringType,
		"name":                             types.StringType,
		"description":                      types.StringType,
		"phase":                            types.ListType{ElemType: types.ObjectType{AttrTypes: getDatasourcePhaseAttrTypes()}},
		"release_retention_with_strategy":  types.ListType{ElemType: types.ObjectType{AttrTypes: retentionAttrTypes}},
		"tentacle_retention_with_strategy": types.ListType{ElemType: types.ObjectType{AttrTypes: retentionAttrTypes}},
	}
}

func getDatasourcePhaseAttrTypes() map[string]attr.Type {
	var retentionAttrTypes = getDatasourceRetentionAttrTypes()
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"is_priority_phase":                     types.BoolType,
		"release_retention_with_strategy":       types.ListType{ElemType: types.ObjectType{AttrTypes: retentionAttrTypes}},
		"tentacle_retention_with_strategy":      types.ListType{ElemType: types.ObjectType{AttrTypes: retentionAttrTypes}},
	}
}

func getDatasourceRetentionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"strategy":         types.StringType,
		"quantity_to_keep": types.Int64Type,
		"unit":             types.StringType,
	}
}
