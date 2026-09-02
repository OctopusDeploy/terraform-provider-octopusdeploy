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

	existingPolicies, err := approvalrules.Get(d.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load approval rules", err.Error())
		return
	}

	flattenedPolicies := []interface{}{}
	for _, policy := range existingPolicies.Items {
		flattenedPolicy, diags := mapApprovalRuleToAttribute(ctx, policy)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		flattenedPolicies = append(flattenedPolicies, flattenedPolicy)
	}

	data.ID = types.StringValue("Approval Policies " + time.Now().UTC().String())
	data.ApprovalRules, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: approvalRuleObjectType()}, flattenedPolicies)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapApprovalRuleToAttribute(ctx context.Context, policy *approvalrules.ApprovalRule) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	approvingUserIDs, listDiags := stringListOrNull(ctx, policy.ApprovingUserIds)
	diags.Append(listDiags...)

	approvingTeamIDs, listDiags := stringListOrNull(ctx, policy.ApprovingTeamIds)
	diags.Append(listDiags...)

	tagScopes := make([]attr.Value, 0, len(policy.TagScopes))
	for _, scope := range policy.TagScopes {
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

	idScopes := make([]attr.Value, 0, len(policy.IdScopes))
	for _, scope := range policy.IdScopes {
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
		"id":                         types.StringValue(policy.GetID()),
		"space_id":                   types.StringValue(policy.SpaceID),
		"name":                       types.StringValue(policy.Name),
		"description":                types.StringValue(policy.Description),
		"scoping_strategy":           types.StringValue(string(policy.ScopingStrategy)),
		"minimum_approvers_required": types.Int64Value(int64(policy.MinimumApproversRequired)),
		"allow_self_approval":        types.BoolValue(policy.AllowSelfApproval),
		"is_disabled":                types.BoolValue(policy.IsDisabled),
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
		"minimum_approvers_required": types.Int64Type,
		"allow_self_approval":        types.BoolType,
		"is_disabled":                types.BoolType,
		"approving_user_ids":         types.ListType{ElemType: types.StringType},
		"approving_team_ids":         types.ListType{ElemType: types.StringType},
		"tag_scopes":                 types.ListType{ElemType: types.ObjectType{AttrTypes: approvalRuleTagScopeObjectType()}},
		"id_scopes":                  types.ListType{ElemType: types.ObjectType{AttrTypes: approvalRuleIdScopeObjectType()}},
	}
}
