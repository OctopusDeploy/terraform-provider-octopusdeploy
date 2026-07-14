package octopusdeploy_framework

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/approvalpolicies"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const approvalPolicyDatasourceName = "approval_policies"

type approvalPolicyDataSource struct {
	*Config
}

func NewApprovalPolicyDataSource() datasource.DataSource {
	return &approvalPolicyDataSource{}
}

var _ datasource.DataSource = &approvalPolicyDataSource{}
var _ datasource.DataSourceWithConfigure = &approvalPolicyDataSource{}

func (d *approvalPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = DataSourceConfiguration(req, resp)
}

func (d *approvalPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(approvalPolicyDatasourceName)
}

func (d *approvalPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.ApprovalPolicySchema{}.GetDatasourceSchema()
}

func (d *approvalPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data schemas.ApprovalPoliciesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := approvalpolicies.ApprovalPoliciesQuery{
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}

	util.DatasourceReading(ctx, "approval policies", query)

	existingPolicies, err := approvalpolicies.Get(d.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load approval policies", err.Error())
		return
	}

	flattenedPolicies := []interface{}{}
	for _, policy := range existingPolicies.Items {
		flattenedPolicy, diags := mapApprovalPolicyToAttribute(ctx, policy)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		flattenedPolicies = append(flattenedPolicies, flattenedPolicy)
	}

	data.ID = types.StringValue("Approval Policies " + time.Now().UTC().String())
	data.ApprovalPolicies, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: approvalPolicyObjectType()}, flattenedPolicies)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapApprovalPolicyToAttribute(ctx context.Context, policy *approvalpolicies.ApprovalPolicy) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	approvingUserIDs, listDiags := approvalPolicyStringListOrNull(ctx, policy.ApprovingUserIds)
	diags.Append(listDiags...)

	approvingTeamIDs, listDiags := approvalPolicyStringListOrNull(ctx, policy.ApprovingTeamIds)
	diags.Append(listDiags...)

	tagScopes := make([]attr.Value, 0, len(policy.TagScopes))
	for _, scope := range policy.TagScopes {
		projectTags, projectTagsDiags := approvalPolicyStringListOrNull(ctx, scope.ProjectTags)
		diags.Append(projectTagsDiags...)

		environmentTags, environmentTagsDiags := approvalPolicyStringListOrNull(ctx, scope.EnvironmentTags)
		diags.Append(environmentTagsDiags...)

		tagScope, tagScopeDiags := types.ObjectValue(approvalPolicyTagScopeObjectType(), map[string]attr.Value{
			"project_tags":     projectTags,
			"environment_tags": environmentTags,
		})
		diags.Append(tagScopeDiags...)
		tagScopes = append(tagScopes, tagScope)
	}

	tagScopesList, tagScopesDiags := types.ListValue(types.ObjectType{AttrTypes: approvalPolicyTagScopeObjectType()}, tagScopes)
	diags.Append(tagScopesDiags...)

	idScopes := make([]attr.Value, 0, len(policy.IdScopes))
	for _, scope := range policy.IdScopes {
		environmentIDs, environmentIDsDiags := approvalPolicyStringListOrNull(ctx, scope.EnvironmentIds)
		diags.Append(environmentIDsDiags...)

		idScope, idScopeDiags := types.ObjectValue(approvalPolicyIdScopeObjectType(), map[string]attr.Value{
			"project_id":      types.StringValue(scope.ProjectId),
			"environment_ids": environmentIDs,
		})
		diags.Append(idScopeDiags...)
		idScopes = append(idScopes, idScope)
	}

	idScopesList, idScopesDiags := types.ListValue(types.ObjectType{AttrTypes: approvalPolicyIdScopeObjectType()}, idScopes)
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

	return types.ObjectValueMust(approvalPolicyObjectType(), attrs), diags
}

// approvalPolicyStringListOrNull converts a Go string slice into a types.List,
// returning a null list when the slice is empty, mirroring the resource's
// mapApprovalPolicyToState behaviour.
func approvalPolicyStringListOrNull(ctx context.Context, values []string) (types.List, diag.Diagnostics) {
	if len(values) == 0 {
		return types.ListNull(types.StringType), nil
	}
	return types.ListValueFrom(ctx, types.StringType, values)
}

func approvalPolicyTagScopeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"project_tags":     types.ListType{ElemType: types.StringType},
		"environment_tags": types.ListType{ElemType: types.StringType},
	}
}

func approvalPolicyIdScopeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"project_id":      types.StringType,
		"environment_ids": types.ListType{ElemType: types.StringType},
	}
}

func approvalPolicyObjectType() map[string]attr.Type {
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
		"tag_scopes":                 types.ListType{ElemType: types.ObjectType{AttrTypes: approvalPolicyTagScopeObjectType()}},
		"id_scopes":                  types.ListType{ElemType: types.ObjectType{AttrTypes: approvalPolicyIdScopeObjectType()}},
	}
}
