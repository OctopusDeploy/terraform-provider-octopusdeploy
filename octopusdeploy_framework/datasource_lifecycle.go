package octopusdeploy_framework

import (
	"context"
	"fmt"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	data.Lifecycles = flattenLifecyclesForDatasourceDEPRECATED(lifecyclesResult.Items)

	data.ID = types.StringValue("Lifecycles " + time.Now().UTC().String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenLifecyclesForDatasourceDEPRECATED(requestedLifecycles []*lifecycles.Lifecycle) types.List {
	var lifecycleAttrTypes = getDatasourceLifecycleAttrTypesIncludingPhaseAndRetentionDEPRECATED()
	lifecyclesList := make([]attr.Value, 0, len(requestedLifecycles))
	for _, lifecycle := range requestedLifecycles {
		lifecycleMap := map[string]attr.Value{
			"id":                        types.StringValue(lifecycle.ID),
			"space_id":                  types.StringValue(lifecycle.SpaceID),
			"name":                      types.StringValue(lifecycle.Name),
			"description":               types.StringValue(lifecycle.Description),
			"phase":                     flattenPhasesForDataSourceDEPRECATED(lifecycle.Phases),
			"release_retention_policy":  flattenRetentionWithoutStrategyForDataSourceDEPRECATED(lifecycle.ReleaseRetentionPolicy),
			"tentacle_retention_policy": flattenRetentionWithoutStrategyForDataSourceDEPRECATED(lifecycle.TentacleRetentionPolicy),
		}
		lifecyclesList = append(lifecyclesList, types.ObjectValueMust(lifecycleAttrTypes, lifecycleMap))

	}
	return types.ListValueMust(types.ObjectType{AttrTypes: lifecycleAttrTypes}, lifecyclesList)
}

func flattenPhasesForDataSourceDEPRECATED(requestedPhases []*lifecycles.Phase) types.List {
	var phaseAttrTypes = getDatasourcePhaseAttrTypesIncludingRetentionDEPRECATED()
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
			"release_retention_policy":              util.Ternary(goPhase.ReleaseRetentionPolicy != nil, flattenRetentionWithoutStrategyForDataSourceDEPRECATED(goPhase.ReleaseRetentionPolicy), ListNullRetentionWithoutStrategyDEPRECATED),
			"tentacle_retention_policy":             util.Ternary(goPhase.TentacleRetentionPolicy != nil, flattenRetentionWithoutStrategyForDataSourceDEPRECATED(goPhase.TentacleRetentionPolicy), ListNullRetentionWithoutStrategyDEPRECATED),
		}
		phasesList = append(phasesList, types.ObjectValueMust(phaseAttrTypes, attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: phaseAttrTypes}, phasesList)
}

func flattenRetentionWithoutStrategyForDataSourceDEPRECATED(requestedRetention *core.RetentionPeriod) types.List {
	var retentionAttrTypes = getDataSourceRetentionWithoutStrategyAttrTypesDEPRECATED()
	if requestedRetention == nil {
		return ListNullRetentionWithoutStrategyDEPRECATED
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: retentionAttrTypes},
		[]attr.Value{
			types.ObjectValueMust(
				retentionAttrTypes,
				map[string]attr.Value{
					"quantity_to_keep":    types.Int64Value(int64(requestedRetention.QuantityToKeep)),
					"should_keep_forever": types.BoolValue(requestedRetention.ShouldKeepForever),
					"unit":                types.StringValue(requestedRetention.Unit),
				},
			),
		},
	)
}

func getDatasourceLifecycleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"space_id":    types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
	}
}
func getDatasourceLifecycleAttrTypesIncludingPhaseAndRetentionDEPRECATED() map[string]attr.Type {

	var datasourceLifecycleAttrTypesDeprecated = getDatasourceLifecycleAttrTypes()
	datasourceLifecycleAttrTypesDeprecated["phase"] = types.ListType{ElemType: types.ObjectType{AttrTypes: getDatasourcePhaseAttrTypesIncludingRetentionDEPRECATED()}}
	datasourceLifecycleAttrTypesDeprecated["release_retention_policy"] = types.ListType{ElemType: types.ObjectType{AttrTypes: getDataSourceRetentionWithoutStrategyAttrTypesDEPRECATED()}}
	datasourceLifecycleAttrTypesDeprecated["tentacle_retention_policy"] = types.ListType{ElemType: types.ObjectType{AttrTypes: getDataSourceRetentionWithoutStrategyAttrTypesDEPRECATED()}}

	return datasourceLifecycleAttrTypesDeprecated
}

func getDatasourcePhaseAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"is_priority_phase":                     types.BoolType,
	}
}
func getDatasourcePhaseAttrTypesIncludingRetentionDEPRECATED() map[string]attr.Type {
	var datasourcePhaseAttrTypesDeprecated = getDatasourcePhaseAttrTypes()
	datasourcePhaseAttrTypesDeprecated["release_retention_policy"] = types.ListType{ElemType: types.ObjectType{AttrTypes: getDataSourceRetentionWithoutStrategyAttrTypesDEPRECATED()}}
	datasourcePhaseAttrTypesDeprecated["tentacle_retention_policy"] = types.ListType{ElemType: types.ObjectType{AttrTypes: getDataSourceRetentionWithoutStrategyAttrTypesDEPRECATED()}}
	return datasourcePhaseAttrTypesDeprecated
}

func getDataSourceRetentionWithoutStrategyAttrTypesDEPRECATED() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}
