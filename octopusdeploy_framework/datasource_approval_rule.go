package octopusdeploy_framework

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/approvalrules"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const approvalRuleDatasourceName = "approval_rules"

type approvalRuleDataSource struct {
	*Config
}

func NewApprovalRuleDataSource() datasource.DataSource {
	return &approvalRuleDataSource{}
}

var _ datasource.DataSource = &approvalRuleDataSource{}
var _ datasource.DataSourceWithConfigure = &approvalRuleDataSource{}

func (d *approvalRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = DataSourceConfiguration(req, resp)
}

func (d *approvalRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(approvalRuleDatasourceName)
}

func (d *approvalRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.ApprovalRuleSchema{}.GetDatasourceSchema()
}

func (d *approvalRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data schemas.ApprovalRulesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := approvalrules.ApprovalRulesQuery{
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}

	util.DatasourceReading(ctx, "approval rules", query)

	existingRules, err := approvalrules.Get(d.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load approval rules", err.Error())
		return
	}

	flattenedRules := []interface{}{}
	for _, rule := range existingRules.Items {
		flattenedRule, diags := mapApprovalRuleToAttribute(ctx, rule)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		flattenedRules = append(flattenedRules, flattenedRule)
	}

	data.ID = types.StringValue("Approval Rules " + time.Now().UTC().String())
	data.ApprovalRules, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: approvalRuleObjectType()}, flattenedRules)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapApprovalRuleToAttribute(ctx context.Context, rule *approvalrules.ApprovalRule) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	approvingUserIDs, listDiags := stringListOrNull(ctx, rule.ApprovingUserIds)
	diags.Append(listDiags...)

	approvingTeamIDs, listDiags := stringListOrNull(ctx, rule.ApprovingTeamIds)
	diags.Append(listDiags...)

	tagScopes := make([]attr.Value, 0, len(rule.TagScopes))
	for _, scope := range rule.TagScopes {
		projectTags, projectTagsDiags := stringListOrNull(ctx, scope.ProjectTags)
		diags.Append(projectTagsDiags...)

		environmentTags, environmentTagsDiags := stringListOrNull(ctx, scope.EnvironmentTags)
		diags.Append(environmentTagsDiags...)

		tagScope, tagScopeDiags := types.ObjectValue(approvalRuleTagScopeObjectType(), map[string]attr.Value{
			"project_tags":     projectTags,
			"environment_tags": environmentTags,
		})
		diags.Append(tagScopeDiags...)
		tagScopes = append(tagScopes, tagScope)
	}

	tagScopesList, tagScopesDiags := types.ListValue(types.ObjectType{AttrTypes: approvalRuleTagScopeObjectType()}, tagScopes)
	diags.Append(tagScopesDiags...)

	idScopes := make([]attr.Value, 0, len(rule.IdScopes))
	for _, scope := range rule.IdScopes {
		environmentIDs, environmentIDsDiags := stringListOrNull(ctx, scope.EnvironmentIds)
		diags.Append(environmentIDsDiags...)

		idScope, idScopeDiags := types.ObjectValue(approvalRuleIdScopeObjectType(), map[string]attr.Value{
			"project_id":      types.StringValue(scope.ProjectId),
			"environment_ids": environmentIDs,
		})
		diags.Append(idScopeDiags...)
		idScopes = append(idScopes, idScope)
	}

	idScopesList, idScopesDiags := types.ListValue(types.ObjectType{AttrTypes: approvalRuleIdScopeObjectType()}, idScopes)
	diags.Append(idScopesDiags...)

	if diags.HasError() {
		return nil, diags
	}

	attrs := map[string]attr.Value{
		"id":                         types.StringValue(rule.GetID()),
		"space_id":                   types.StringValue(rule.SpaceID),
		"name":                       types.StringValue(rule.Name),
		"description":                types.StringValue(rule.Description),
		"scoping_strategy":           types.StringValue(string(rule.ScopingStrategy)),
		"tenant_approval_strategy":   types.StringValue(string(rule.TenantApprovalStrategy)),
		"minimum_approvers_required": types.Int64Value(int64(rule.MinimumApproversRequired)),
		"allow_self_approval":        types.BoolValue(rule.AllowSelfApproval),
		"is_disabled":                types.BoolValue(rule.IsDisabled),
		"approving_user_ids":         approvingUserIDs,
		"approving_team_ids":         approvingTeamIDs,
		"tag_scopes":                 tagScopesList,
		"id_scopes":                  idScopesList,
	}

	return types.ObjectValueMust(approvalRuleObjectType(), attrs), diags
}

func approvalRuleTagScopeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"project_tags":     types.ListType{ElemType: types.StringType},
		"environment_tags": types.ListType{ElemType: types.StringType},
	}
}

func approvalRuleIdScopeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"project_id":      types.StringType,
		"environment_ids": types.ListType{ElemType: types.StringType},
	}
}

func approvalRuleObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                         types.StringType,
		"space_id":                   types.StringType,
		"name":                       types.StringType,
		"description":                types.StringType,
		"scoping_strategy":           types.StringType,
		"tenant_approval_strategy":   types.StringType,
		"minimum_approvers_required": types.Int64Type,
		"allow_self_approval":        types.BoolType,
		"is_disabled":                types.BoolType,
		"approving_user_ids":         types.ListType{ElemType: types.StringType},
		"approving_team_ids":         types.ListType{ElemType: types.StringType},
		"tag_scopes":                 types.ListType{ElemType: types.ObjectType{AttrTypes: approvalRuleTagScopeObjectType()}},
		"id_scopes":                  types.ListType{ElemType: types.ObjectType{AttrTypes: approvalRuleIdScopeObjectType()}},
	}
}
