package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type channelResource struct {
	*Config
}

func NewChannelResource() resource.Resource {
	return &channelResource{}
}

var _ resource.ResourceWithImportState = &channelResource{}

func (r *channelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("channel")
}

func (r *channelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ChannelSchema{}.GetResourceSchema()
}

func (r *channelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *channelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan schemas.ChannelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating channel", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	channel := expandChannel(ctx, plan)
	createdChannel, err := channels.Add(r.Config.Client, channel)
	if err != nil {
		resp.Diagnostics.AddError("Error creating channel", err.Error())
		return
	}

	state := flattenChannel(ctx, createdChannel, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *channelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state schemas.ChannelModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channel, err := channels.GetByID(r.Client, state.SpaceId.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "channelResource"); err != nil {
			resp.Diagnostics.AddError("unable to load channel", err.Error())
		}
		return
	}

	newState := flattenChannel(ctx, channel, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
	return
}

func (r *channelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan schemas.ChannelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channel := expandChannel(ctx, plan)
	updateReq := buildChannelUpdateRequest(channel, plan)
	updatedChannel, err := channels.UpdateChannel(r.Client, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating channel", err.Error())
		return
	}

	state := flattenChannel(ctx, updatedChannel, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	return
}

func (r *channelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state schemas.ChannelModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := channels.DeleteByID(r.Client, state.SpaceId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting channel", err.Error())
		return
	}
}

func (*channelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandChannel(ctx context.Context, model schemas.ChannelModel) *channels.Channel {
	var channelName = model.Name.ValueString()
	var projectId = model.ProjectId.ValueString()

	channel := channels.NewChannel(channelName, projectId)

	channel.ID = model.ID.ValueString()
	channel.Description = model.Description.ValueString()
	channel.IsDefault = model.IsDefault.ValueBool()
	channel.LifecycleID = model.LifecycleId.ValueString()
	channel.CustomFieldDefinitions = expandChannelCustomFieldDefinitions(model.CustomFieldDefinitions)
	channel.GitReferenceRules = expandChannelStringList(model.GitReferenceRules)
	channel.GitResourceRules = expandChannelGitResourceRules(model.GitResourceRule)
	channel.Rules = expandChannelRules(model.Rule)
	channel.SpaceID = model.SpaceId.ValueString()
	channel.TenantTags = util.ExpandStringSet(model.TenantTags)
	channel.Type = channels.ChannelType(model.Type.ValueString())
	channel.ParentEnvironmentID = model.ParentEnvironmentID.ValueString()
	channel.EphemeralEnvironmentNameTemplate = model.EphemeralEnvironmentNameTemplate.ValueString()

	return channel
}

func buildChannelUpdateRequest(channel *channels.Channel, plan schemas.ChannelModel) *newclient.UpdateRequest[channels.Channel] {
	updateReq := newclient.NewUpdateRequest(channel)
	if !plan.CustomFieldDefinitions.IsUnknown() &&
		(plan.CustomFieldDefinitions.IsNull() || len(plan.CustomFieldDefinitions.Elements()) == 0) {
		updateReq.Clear("CustomFieldDefinitions")
	}
	if !plan.Rule.IsUnknown() &&
		(plan.Rule.IsNull() || len(plan.Rule.Elements()) == 0) {
		updateReq.Clear("Rules")
	}
	if !plan.TenantTags.IsUnknown() &&
		(plan.TenantTags.IsNull() || len(plan.TenantTags.Elements()) == 0) {
		updateReq.Clear("TenantTags")
	}
	if !plan.GitReferenceRules.IsUnknown() &&
		(plan.GitReferenceRules.IsNull() || len(plan.GitReferenceRules.Elements()) == 0) {
		updateReq.Clear("GitReferenceRules")
	}
	if !plan.GitResourceRule.IsUnknown() &&
		(plan.GitResourceRule.IsNull() || len(plan.GitResourceRule.Elements()) == 0) {
		updateReq.Clear("GitResourceRules")
	}

	return updateReq
}

func expandChannelStringList(values types.List) []string {
	if values.IsNull() || values.IsUnknown() || len(values.Elements()) == 0 {
		return nil
	}

	var result []string
	values.ElementsAs(context.Background(), &result, false)
	if len(result) == 0 {
		return nil
	}

	return result
}

func expandChannelGitResourceRules(rules types.List) []channels.ChannelGitResourceRule {
	if rules.IsNull() || rules.IsUnknown() || len(rules.Elements()) == 0 {
		return nil
	}

	result := make([]channels.ChannelGitResourceRule, 0, len(rules.Elements()))
	for _, ruleElem := range rules.Elements() {
		ruleObj := ruleElem.(types.Object)
		ruleAttrs := ruleObj.Attributes()

		var gitResourceRule channels.ChannelGitResourceRule
		if v, ok := ruleAttrs["id"].(types.String); ok && !v.IsNull() && !v.IsUnknown() {
			gitResourceRule.Id = v.ValueString()
		}
		if v, ok := ruleAttrs["rules"].(types.List); ok && !v.IsNull() && !v.IsUnknown() {
			gitResourceRule.Rules = expandChannelStringList(v)
		}
		if v, ok := ruleAttrs["git_dependency_action"].(types.List); ok && !v.IsNull() && !v.IsUnknown() {
			gitResourceRule.GitDependencyActions = expandDeploymentActionGitDependencies(v)
		}

		result = append(result, gitResourceRule)
	}

	return result
}

func expandDeploymentActionGitDependencies(actions types.List) []gitdependencies.DeploymentActionGitDependency {
	if actions.IsNull() || actions.IsUnknown() || len(actions.Elements()) == 0 {
		return nil
	}

	result := make([]gitdependencies.DeploymentActionGitDependency, 0, len(actions.Elements()))
	for _, actionElem := range actions.Elements() {
		actionObj := actionElem.(types.Object)
		actionAttrs := actionObj.Attributes()

		var action gitdependencies.DeploymentActionGitDependency
		if v, ok := actionAttrs["deployment_action_slug"].(types.String); ok && !v.IsNull() && !v.IsUnknown() {
			action.DeploymentActionSlug = v.ValueString()
		}
		if v, ok := actionAttrs["git_dependency_name"].(types.String); ok && !v.IsNull() && !v.IsUnknown() {
			action.GitDependencyName = v.ValueString()
		}

		result = append(result, action)
	}

	return result
}

func expandChannelRules(rules types.List) []channels.ChannelRule {
	if rules.IsNull() || rules.IsUnknown() || len(rules.Elements()) == 0 {
		return nil
	}

	result := make([]channels.ChannelRule, 0, len(rules.Elements()))

	for _, ruleElem := range rules.Elements() {
		ruleObj := ruleElem.(types.Object)
		ruleAttrs := ruleObj.Attributes()

		channelRule := expandChannelRuleFromAttrs(ruleAttrs)
		result = append(result, channelRule)
	}

	return result
}

func expandChannelRuleFromAttrs(attrs map[string]attr.Value) channels.ChannelRule {
	var channelRule channels.ChannelRule

	if v, ok := attrs["id"].(types.String); ok && !v.IsNull() {
		channelRule.ID = v.ValueString()
	}

	if v, ok := attrs["tag"].(types.String); ok && !v.IsNull() {
		channelRule.Tag = v.ValueString()
	}

	if v, ok := attrs["version_range"].(types.String); ok && !v.IsNull() {
		channelRule.VersionRange = v.ValueString()
	}

	if v, ok := attrs["action_package"].(types.List); ok && !v.IsNull() {
		channelRule.ActionPackages = expandChannelRuleDeploymentActionPackagesFromList(v)
	}

	return channelRule
}

func expandChannelRuleDeploymentActionPackagesFromList(actionPackages types.List) []packages.DeploymentActionPackage {
	if actionPackages.IsNull() || actionPackages.IsUnknown() || len(actionPackages.Elements()) == 0 {
		return nil
	}

	result := make([]packages.DeploymentActionPackage, 0, len(actionPackages.Elements()))

	for _, packageElem := range actionPackages.Elements() {
		packageObj := packageElem.(types.Object)
		packageAttrs := packageObj.Attributes()

		var actionPackage packages.DeploymentActionPackage

		if v, ok := packageAttrs["deployment_action"].(types.String); ok && !v.IsNull() {
			actionPackage.DeploymentAction = v.ValueString()
		}

		if v, ok := packageAttrs["package_reference"].(types.String); ok && !v.IsNull() {
			actionPackage.PackageReference = v.ValueString()
		}

		result = append(result, actionPackage)
	}

	return result
}

func expandChannelRuleDeploymentActionPackages(actionPackages []map[string]interface{}) []packages.DeploymentActionPackage {
	var actionPackagesExpanded []packages.DeploymentActionPackage
	for _, actionPackage := range actionPackages {
		actionPackageExpanded := expandChannelRuleDeploymentActionPackage(actionPackage)
		actionPackagesExpanded = append(actionPackagesExpanded, actionPackageExpanded)
	}
	return actionPackagesExpanded
}

func expandChannelRuleDeploymentActionPackage(actionPackageMap map[string]interface{}) packages.DeploymentActionPackage {
	return packages.DeploymentActionPackage{
		DeploymentAction: actionPackageMap["deployment_action"].(string),
		PackageReference: actionPackageMap["package_reference"].(string),
	}

}

func flattenChannel(ctx context.Context, channel *channels.Channel, model schemas.ChannelModel) schemas.ChannelModel {
	model.ID = types.StringValue(channel.GetID())
	model.Description = types.StringValue(channel.Description)

	if !channel.IsDefault && model.IsDefault.IsNull() {
		model.IsDefault = types.BoolNull()
	} else {
		model.IsDefault = types.BoolValue(channel.IsDefault)
	}

	if channel.LifecycleID == "" && model.LifecycleId.IsNull() {
		model.LifecycleId = types.StringNull()
	} else {
		model.LifecycleId = types.StringValue(channel.LifecycleID)
	}

	model.Name = types.StringValue(channel.Name)
	model.ProjectId = types.StringValue(channel.ProjectID)

	model.CustomFieldDefinitions = flattenChannelCustomFieldDefinitions(channel.CustomFieldDefinitions)
	model.GitReferenceRules = flattenChannelStringList(channel.GitReferenceRules, model.GitReferenceRules)
	model.GitResourceRule = flattenChannelGitResourceRules(channel.GitResourceRules, model.GitResourceRule)
	model.Rule = flattenChannelRules(channel.Rules, model.Rule)

	if channel.SpaceID == "" && model.SpaceId.IsNull() {
		model.SpaceId = types.StringNull()
	} else {
		model.SpaceId = types.StringValue(channel.SpaceID)
	}

	model.TenantTags = util.FlattenStringSet(channel.TenantTags, model.TenantTags)

	model.Type = types.StringValue(string(channel.Type))
	model.ParentEnvironmentID = util.StringOrNull(channel.ParentEnvironmentID)
	model.EphemeralEnvironmentNameTemplate = util.StringOrNull(channel.EphemeralEnvironmentNameTemplate)

	return model
}

func flattenChannelStringList(values []string, current types.List) types.List {
	if len(values) == 0 {
		if current.IsNull() {
			return types.ListNull(types.StringType)
		}
		return types.ListValueMust(types.StringType, []attr.Value{})
	}

	flattened := make([]attr.Value, 0, len(values))
	for _, value := range values {
		flattened = append(flattened, types.StringValue(value))
	}

	return types.ListValueMust(types.StringType, flattened)
}

func flattenChannelRules(rules []channels.ChannelRule, currentRules types.List) types.List {
	if rules == nil || len(rules) == 0 {
		if currentRules.IsNull() {
			return types.ListNull(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()})
		}
		return types.ListValueMust(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()}, []attr.Value{})
	}

	flattenedRules := make([]attr.Value, 0, len(rules))
	for _, rule := range rules {
		obj := flattenChannelRule(&rule)
		flattenedRules = append(flattenedRules, obj)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()}, flattenedRules)
}

func flattenChannelGitResourceRules(rules []channels.ChannelGitResourceRule, currentRules types.List) types.List {
	if len(rules) == 0 {
		if currentRules.IsNull() {
			return types.ListNull(types.ObjectType{AttrTypes: getChannelGitResourceRuleAttrTypes()})
		}
		return types.ListValueMust(types.ObjectType{AttrTypes: getChannelGitResourceRuleAttrTypes()}, []attr.Value{})
	}

	flattenedRules := make([]attr.Value, 0, len(rules))
	for _, rule := range rules {
		flattenedRules = append(flattenedRules, types.ObjectValueMust(getChannelGitResourceRuleAttrTypes(), map[string]attr.Value{
			"id":                    util.StringOrNull(rule.Id),
			"rules":                 flattenChannelStringList(rule.Rules, types.ListNull(types.StringType)),
			"git_dependency_action": flattenDeploymentActionGitDependencies(rule.GitDependencyActions),
		}))
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getChannelGitResourceRuleAttrTypes()}, flattenedRules)
}

func flattenDeploymentActionGitDependencies(actions []gitdependencies.DeploymentActionGitDependency) types.List {
	if len(actions) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: getDeploymentActionGitDependencyAttrTypes()})
	}

	flattenedActions := make([]attr.Value, 0, len(actions))
	for _, action := range actions {
		flattenedActions = append(flattenedActions, types.ObjectValueMust(getDeploymentActionGitDependencyAttrTypes(), map[string]attr.Value{
			"deployment_action_slug": util.StringOrNull(action.DeploymentActionSlug),
			"git_dependency_name":    util.StringOrNull(action.GitDependencyName),
		}))
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getDeploymentActionGitDependencyAttrTypes()}, flattenedActions)
}

func flattenChannelRule(rule *channels.ChannelRule) types.Object {
	return types.ObjectValueMust(getChannelRuleAttrTypes(), map[string]attr.Value{
		"action_package": flattenChannelRuleDeploymentActionPackages(rule.ActionPackages),
		"id":             types.StringValue(rule.ID),
		"tag":            util.StringOrNull(rule.Tag),
		"version_range":  util.StringOrNull(rule.VersionRange),
	})

}

func flattenChannelRuleDeploymentActionPackages(actionPackages []packages.DeploymentActionPackage) types.List {
	if actionPackages == nil || len(actionPackages) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: getChannelRuleDeploymentActionPackageAttrTypes()})
	}

	flattenedActionPackages := make([]attr.Value, 0, len(actionPackages))
	for _, actionPackage := range actionPackages {
		obj := flattenChannelRuleDeploymentActionPackage(&actionPackage)
		flattenedActionPackages = append(flattenedActionPackages, obj)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getChannelRuleDeploymentActionPackageAttrTypes()}, flattenedActionPackages)
}

func flattenChannelRuleDeploymentActionPackage(actionPackage *packages.DeploymentActionPackage) types.Object {
	return types.ObjectValueMust(getChannelRuleDeploymentActionPackageAttrTypes(), map[string]attr.Value{
		"deployment_action": types.StringValue(actionPackage.DeploymentAction),
		"package_reference": util.StringOrNull(actionPackage.PackageReference),
	})
}

func getChannelRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"action_package": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: getChannelRuleDeploymentActionPackageAttrTypes(),
			},
		},
		"id":            types.StringType,
		"tag":           types.StringType,
		"version_range": types.StringType,
	}
}

func getChannelGitResourceRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":    types.StringType,
		"rules": types.ListType{ElemType: types.StringType},
		"git_dependency_action": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: getDeploymentActionGitDependencyAttrTypes(),
			},
		},
	}
}

func getDeploymentActionGitDependencyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"deployment_action_slug": types.StringType,
		"git_dependency_name":    types.StringType,
	}
}

func getChannelRuleDeploymentActionPackageAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"deployment_action": types.StringType,
		"package_reference": types.StringType,
	}
}

func expandChannelCustomFieldDefinitions(defs types.List) []channels.ChannelCustomFieldDefinition {
	if defs.IsNull() || defs.IsUnknown() || len(defs.Elements()) == 0 {
		return []channels.ChannelCustomFieldDefinition{}
	}

	result := make([]channels.ChannelCustomFieldDefinition, 0, len(defs.Elements()))
	for _, elem := range defs.Elements() {
		obj := elem.(types.Object)
		attrs := obj.Attributes()

		var def channels.ChannelCustomFieldDefinition
		if v, ok := attrs["field_name"].(types.String); ok && !v.IsNull() {
			def.FieldName = v.ValueString()
		}
		if v, ok := attrs["description"].(types.String); ok && !v.IsNull() {
			def.Description = v.ValueString()
		}
		result = append(result, def)
	}
	return result
}

func flattenChannelCustomFieldDefinitions(defs []channels.ChannelCustomFieldDefinition) types.List {
	if len(defs) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: getChannelCustomFieldDefinitionAttrTypes()})
	}

	elems := make([]attr.Value, 0, len(defs))
	for _, def := range defs {
		elems = append(elems, types.ObjectValueMust(getChannelCustomFieldDefinitionAttrTypes(), map[string]attr.Value{
			"field_name":  types.StringValue(def.FieldName),
			"description": types.StringValue(def.Description),
		}))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: getChannelCustomFieldDefinitionAttrTypes()}, elems)
}

func getChannelCustomFieldDefinitionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"field_name":  types.StringType,
		"description": types.StringType,
	}
}
