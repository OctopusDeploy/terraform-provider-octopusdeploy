package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/approvalpolicies"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const approvalPolicyResourceName = "approval_policy"

type approvalPolicyResource struct {
	*Config
}

var _ resource.Resource = &approvalPolicyResource{}
var _ resource.ResourceWithImportState = &approvalPolicyResource{}
var _ resource.ResourceWithValidateConfig = &approvalPolicyResource{}
var _ resource.ResourceWithModifyPlan = &approvalPolicyResource{}

func NewApprovalPolicyResource() resource.Resource {
	return &approvalPolicyResource{}
}

func (r *approvalPolicyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(approvalPolicyResourceName)
}

func (r *approvalPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ApprovalPolicySchema{}.GetResourceSchema()
}

func (r *approvalPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

// ValidateConfig enforces that the configured scope matches the scoping strategy:
// tag_scopes with the "Tag" strategy and id_scopes with the "Id" strategy, and
// that the two scope types are never set at the same time.
func (r *approvalPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config schemas.ApprovalPolicyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasTagScopes := len(config.TagScopes) > 0
	hasIdScopes := len(config.IdScopes) > 0

	if hasTagScopes && hasIdScopes {
		resp.Diagnostics.AddAttributeError(
			path.Root("id_scopes"),
			"Conflicting approval policy scopes",
			`Only one of "tag_scopes" or "id_scopes" may be set on an approval policy. Use "tag_scopes" with scoping_strategy "Tag", or "id_scopes" with scoping_strategy "Id".`,
		)
		return
	}

	// When the scoping strategy is known, the provided scope must match it.
	if config.ScopingStrategy.IsNull() || config.ScopingStrategy.IsUnknown() {
		return
	}

	switch config.ScopingStrategy.ValueString() {
	case string(approvalpolicies.ApprovalPolicyScopingStrategyTag):
		if hasIdScopes {
			resp.Diagnostics.AddAttributeError(
				path.Root("id_scopes"),
				"Scope does not match scoping strategy",
				`"id_scopes" cannot be set when scoping_strategy is "Tag". Use "tag_scopes", or set scoping_strategy to "Id".`,
			)
		}
	case string(approvalpolicies.ApprovalPolicyScopingStrategyId):
		if hasTagScopes {
			resp.Diagnostics.AddAttributeError(
				path.Root("tag_scopes"),
				"Scope does not match scoping strategy",
				`"tag_scopes" cannot be set when scoping_strategy is "Id". Use "id_scopes", or set scoping_strategy to "Tag".`,
			)
		}
	}
}

// ModifyPlan resolves tag_scopes references that are given as canonical tag
// names (e.g. "TagSet/Tag") into their stable tag IDs (e.g. "TagSets-1/Tags-1"),
// which is the form the API returns. Canonical names are not immutable (they
// change when a tag or tag set is renamed), so storing IDs keeps state stable
// across renames and avoids perpetual diffs. Values that are already IDs, or are
// unknown at plan time, are left untouched.
func (r *approvalPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return // the resource is being destroyed; nothing to resolve
	}
	if r.Config == nil || r.Config.Client == nil {
		return
	}

	// Read tag_scopes as a framework type first: it may contain unknown values
	// (e.g. a tag created in the same apply and referenced by its id attribute),
	// which cannot be decoded into the Go model. When any value is unknown, skip
	// plan-time resolution and let it resolve naturally at apply.
	var tagScopesList types.List
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("tag_scopes"), &tagScopesList)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if tagScopesList.IsNull() || tagScopesList.IsUnknown() || len(tagScopesList.Elements()) == 0 {
		return
	}
	rawTagScopes, err := tagScopesList.ToTerraformValue(ctx)
	if err != nil {
		resp.Diagnostics.AddError("unable to inspect approval policy tag scopes", err.Error())
		return
	}
	if !rawTagScopes.IsFullyKnown() {
		return
	}

	var plan schemas.ApprovalPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || len(plan.TagScopes) == 0 {
		return
	}

	spaceID := plan.SpaceID.ValueString()
	if spaceID == "" {
		spaceID = r.Config.SpaceID
	}

	tagSets, err := tagsets.GetAll(r.Config.Client, spaceID)
	if err != nil {
		resp.Diagnostics.AddError("unable to resolve approval policy tags", err.Error())
		return
	}

	nameToID := make(map[string]string)
	validID := make(map[string]bool)
	for _, ts := range tagSets {
		for _, tag := range ts.Tags {
			nameToID[tag.CanonicalTagName] = tag.ID
			validID[tag.ID] = true
		}
	}

	for i := range plan.TagScopes {
		plan.TagScopes[i].ProjectTags = resolveApprovalPolicyTags(ctx, plan.TagScopes[i].ProjectTags, nameToID, validID, path.Root("tag_scopes").AtListIndex(i).AtName("project_tags"), &resp.Diagnostics)
		plan.TagScopes[i].EnvironmentTags = resolveApprovalPolicyTags(ctx, plan.TagScopes[i].EnvironmentTags, nameToID, validID, path.Root("tag_scopes").AtListIndex(i).AtName("environment_tags"), &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Only update tag_scopes; setting the whole plan would disturb other
	// Optional+Computed attributes and produce spurious "known after apply" diffs.
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("tag_scopes"), plan.TagScopes)...)
}

// resolveApprovalPolicyTags maps each element of a tag set to a stable tag ID,
// translating canonical tag names to IDs and leaving IDs unchanged. Null or
// unknown sets are returned as-is: an unknown set (for example a tag created in
// the same apply and referenced by its id attribute) is resolved naturally at
// apply. A known value that matches neither an existing tag ID nor canonical
// name is an error — Terraform does not allow a configured value to be planned
// as unknown, so a tag created in the same apply must be referenced by ID.
func resolveApprovalPolicyTags(ctx context.Context, tags types.Set, nameToID map[string]string, validID map[string]bool, attrPath path.Path, diags *diag.Diagnostics) types.Set {
	if tags.IsNull() || tags.IsUnknown() {
		return tags
	}

	var values []string
	diags.Append(tags.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return tags
	}

	resolved := make([]string, 0, len(values))
	changed := false
	for _, v := range values {
		switch {
		case validID[v]:
			resolved = append(resolved, v)
		default:
			if id, ok := nameToID[v]; ok {
				resolved = append(resolved, id)
				changed = true
			} else {
				diags.AddAttributeError(
					attrPath,
					"Unknown tag reference",
					fmt.Sprintf("%q is not a known tag ID or canonical tag name in this space. If the tag is created in the same apply, reference it by ID instead (for example octopusdeploy_tag.example.id).", v),
				)
				resolved = append(resolved, v)
			}
		}
	}

	if diags.HasError() || !changed {
		return tags
	}

	newSet, d := types.SetValueFrom(ctx, types.StringType, resolved)
	diags.Append(d...)
	return newSet
}

func (r *approvalPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan *schemas.ApprovalPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.SpaceID.IsNull() || plan.SpaceID.IsUnknown() {
		plan.SpaceID = types.StringValue(r.Config.SpaceID)
	}

	policy, diags := mapApprovalPolicyFromState(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdPolicy, err := approvalpolicies.Add(r.Config.Client, plan.SpaceID.ValueString(), policy)
	if err != nil {
		resp.Diagnostics.AddError("error while creating approval policy", err.Error())
		return
	}

	resp.Diagnostics.Append(mapApprovalPolicyToState(ctx, plan, createdPolicy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *approvalPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state *schemas.ApprovalPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.SpaceID.IsNull() || state.SpaceID.IsUnknown() {
		state.SpaceID = types.StringValue(r.Config.SpaceID)
	}

	policy, err := approvalpolicies.GetByID(r.Config.Client, state.SpaceID.ValueString(), state.GetID())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "approval policy"); err != nil {
			resp.Diagnostics.AddError("unable to load approval policy", err.Error())
		}
		return
	}

	resp.Diagnostics.Append(mapApprovalPolicyToState(ctx, state, policy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *approvalPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan *schemas.ApprovalPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.SpaceID.IsNull() || plan.SpaceID.IsUnknown() {
		plan.SpaceID = types.StringValue(r.Config.SpaceID)
	}

	policy, diags := mapApprovalPolicyFromState(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	policy.SetID(plan.ID.ValueString())

	updatedPolicy, err := approvalpolicies.Update(r.Config.Client, plan.SpaceID.ValueString(), policy)
	if err != nil {
		resp.Diagnostics.AddError("error while updating approval policy", err.Error())
		return
	}

	resp.Diagnostics.Append(mapApprovalPolicyToState(ctx, plan, updatedPolicy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *approvalPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state *schemas.ApprovalPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.SpaceID.IsNull() || state.SpaceID.IsUnknown() {
		state.SpaceID = types.StringValue(r.Config.SpaceID)
	}

	err := approvalpolicies.DeleteByID(r.Config.Client, state.SpaceID.ValueString(), state.GetID())
	if err != nil {
		resp.Diagnostics.AddError("unable to delete approval policy", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *approvalPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapApprovalPolicyFromState(state *schemas.ApprovalPolicyResourceModel) (*approvalpolicies.ApprovalPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	var approvingUserIDs []string
	if !state.ApprovingUserIDs.IsNull() {
		diags.Append(state.ApprovingUserIDs.ElementsAs(context.Background(), &approvingUserIDs, false)...)
	}

	var approvingTeamIDs []string
	if !state.ApprovingTeamIDs.IsNull() {
		diags.Append(state.ApprovingTeamIDs.ElementsAs(context.Background(), &approvingTeamIDs, false)...)
	}

	tagScopes := make([]approvalpolicies.ApprovalPolicyTagScope, 0, len(state.TagScopes))
	for _, s := range state.TagScopes {
		var projectTags []string
		if !s.ProjectTags.IsNull() {
			diags.Append(s.ProjectTags.ElementsAs(context.Background(), &projectTags, false)...)
		}

		var environmentTags []string
		if !s.EnvironmentTags.IsNull() {
			diags.Append(s.EnvironmentTags.ElementsAs(context.Background(), &environmentTags, false)...)
		}

		tagScopes = append(tagScopes, approvalpolicies.ApprovalPolicyTagScope{
			ProjectTags:     projectTags,
			EnvironmentTags: environmentTags,
		})
	}

	idScopes := make([]approvalpolicies.ApprovalPolicyIdScope, 0, len(state.IdScopes))
	for _, s := range state.IdScopes {
		var environmentIDs []string
		if !s.EnvironmentIDs.IsNull() {
			diags.Append(s.EnvironmentIDs.ElementsAs(context.Background(), &environmentIDs, false)...)
		}

		idScopes = append(idScopes, approvalpolicies.ApprovalPolicyIdScope{
			ProjectId:      s.ProjectID.ValueString(),
			EnvironmentIds: environmentIDs,
		})
	}

	if diags.HasError() {
		return nil, diags
	}

	policy := &approvalpolicies.ApprovalPolicy{
		SpaceID:                  state.SpaceID.ValueString(),
		Name:                     state.Name.ValueString(),
		Description:              state.Description.ValueString(),
		ScopingStrategy:          approvalpolicies.ApprovalPolicyScopingStrategy(state.ScopingStrategy.ValueString()),
		TagScopes:                tagScopes,
		IdScopes:                 idScopes,
		MinimumApproversRequired: int(state.MinimumApproversRequired.ValueInt64()),
		AllowSelfApproval:        state.AllowSelfApproval.ValueBool(),
		IsDisabled:               state.IsDisabled.ValueBool(),
		ApprovingUserIds:         approvingUserIDs,
		ApprovingTeamIds:         approvingTeamIDs,
	}
	policy.ID = state.ID.ValueString()

	return policy, diags
}

func mapApprovalPolicyToState(ctx context.Context, state *schemas.ApprovalPolicyResourceModel, policy *approvalpolicies.ApprovalPolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	state.ID = types.StringValue(policy.GetID())
	state.SpaceID = types.StringValue(policy.SpaceID)
	state.Name = types.StringValue(policy.Name)
	state.Description = types.StringValue(policy.Description)
	state.ScopingStrategy = types.StringValue(string(policy.ScopingStrategy))
	state.MinimumApproversRequired = types.Int64Value(int64(policy.MinimumApproversRequired))
	state.AllowSelfApproval = types.BoolValue(policy.AllowSelfApproval)
	state.IsDisabled = types.BoolValue(policy.IsDisabled)

	approvingUserIDs, listDiags := stringListOrNull(ctx, policy.ApprovingUserIds)
	diags.Append(listDiags...)
	state.ApprovingUserIDs = approvingUserIDs

	approvingTeamIDs, listDiags := stringListOrNull(ctx, policy.ApprovingTeamIds)
	diags.Append(listDiags...)
	state.ApprovingTeamIDs = approvingTeamIDs

	var tagScopes []schemas.ApprovalPolicyTagScopeModel
	for _, s := range policy.TagScopes {
		projectTags, projectTagsDiags := stringSetOrNull(ctx, s.ProjectTags)
		diags.Append(projectTagsDiags...)

		environmentTags, environmentTagsDiags := stringSetOrNull(ctx, s.EnvironmentTags)
		diags.Append(environmentTagsDiags...)

		tagScopes = append(tagScopes, schemas.ApprovalPolicyTagScopeModel{
			ProjectTags:     projectTags,
			EnvironmentTags: environmentTags,
		})
	}
	state.TagScopes = tagScopes

	var idScopes []schemas.ApprovalPolicyIdScopeModel
	for _, s := range policy.IdScopes {
		environmentIDs, environmentIDsDiags := stringListOrNull(ctx, s.EnvironmentIds)
		diags.Append(environmentIDsDiags...)

		idScopes = append(idScopes, schemas.ApprovalPolicyIdScopeModel{
			ProjectID:      types.StringValue(s.ProjectId),
			EnvironmentIDs: environmentIDs,
		})
	}
	state.IdScopes = idScopes

	return diags
}

// stringListOrNull converts a Go string slice into a types.List, returning a
// null list when the slice is empty so that Optional-only attributes don't
// flip from a null plan value to a known empty list in state (which the
// framework reports as an inconsistent-result error).
func stringListOrNull(ctx context.Context, values []string) (types.List, diag.Diagnostics) {
	if len(values) == 0 {
		return types.ListNull(types.StringType), nil
	}
	return types.ListValueFrom(ctx, types.StringType, values)
}

// stringSetOrNull is the set-typed counterpart of stringListOrNull, used for the
// order-insensitive tag scope collections.
func stringSetOrNull(ctx context.Context, values []string) (types.Set, diag.Diagnostics) {
	if len(values) == 0 {
		return types.SetNull(types.StringType), nil
	}
	return types.SetValueFrom(ctx, types.StringType, values)
}
