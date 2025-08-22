package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
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
	updatedChannel, err := channels.Update(r.Client, channel)
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
	channel.Rules = expandChannelRules(model.Rule)
	channel.SpaceID = model.SpaceId.ValueString()
	channel.TenantTags = util.ExpandStringSet(model.TenantTags)

	return channel
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

	model.Rule = flattenChannelRules(channel.Rules, model.Rule)

	if channel.SpaceID == "" && model.SpaceId.IsNull() {
		model.SpaceId = types.StringNull()
	} else {
		model.SpaceId = types.StringValue(channel.SpaceID)
	}

	model.TenantTags = util.FlattenStringSet(channel.TenantTags, model.TenantTags)

	return model
}

func flattenChannelRules(rules []channels.ChannelRule, currentRules types.List) types.List {
	if rules == nil || len(rules) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()})
	}

	currentRulesMap := make(map[string]types.Object)
	if !currentRules.IsNull() && !currentRules.IsUnknown() {
		for _, elem := range currentRules.Elements() {
			ruleObj := elem.(types.Object)
			attrs := ruleObj.Attributes()
			if idAttr, ok := attrs["id"].(types.String); ok && !idAttr.IsNull() && !idAttr.IsUnknown() {
				currentRulesMap[idAttr.ValueString()] = ruleObj
			}
		}
	}

	flattenedRules := make([]attr.Value, 0, len(rules))
	for _, rule := range rules {
		currentRule := currentRulesMap[rule.ID]
		obj := flattenChannelRule(&rule, currentRule)
		flattenedRules = append(flattenedRules, obj)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getChannelRuleAttrTypes()}, flattenedRules)
}

func flattenChannelRule(rule *channels.ChannelRule, currentRule types.Object) types.Object {
	var currentTag, currentVersionRange types.String
	var currentActionPackages types.List
	if !currentRule.IsNull() && !currentRule.IsUnknown() {
		attrs := currentRule.Attributes()
		if v, ok := attrs["tag"]; ok {
			currentTag = v.(types.String)
		}
		if v, ok := attrs["version_range"]; ok {
			currentVersionRange = v.(types.String)
		}
		if v, ok := attrs["action_package"]; ok {
			currentActionPackages = v.(types.List)
		}
	}

	var tagValue, versionRangeValue types.String
	if rule.Tag == "" && currentTag.IsNull() {
		tagValue = types.StringNull()
	} else {
		tagValue = types.StringValue(rule.Tag)
	}

	if rule.VersionRange == "" && currentVersionRange.IsNull() {
		versionRangeValue = types.StringNull()
	} else {
		versionRangeValue = types.StringValue(rule.VersionRange)
	}

	return types.ObjectValueMust(getChannelRuleAttrTypes(), map[string]attr.Value{
		"action_package": flattenChannelRuleDeploymentActionPackages(rule.ActionPackages, currentActionPackages),
		"id":             types.StringValue(rule.ID),
		"tag":            tagValue,
		"version_range":  versionRangeValue,
	})

}

func flattenChannelRuleDeploymentActionPackages(actionPackages []packages.DeploymentActionPackage, currentActionPackages types.List) types.List {
	if actionPackages == nil || len(actionPackages) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: getChannelRuleDeploymentActionPackageAttrTypes()})
	}

	currentActionPackagesMap := make(map[string]types.Object)
	if !currentActionPackages.IsNull() && !currentActionPackages.IsUnknown() {
		for _, elem := range currentActionPackages.Elements() {
			packageObj := elem.(types.Object)
			attrs := packageObj.Attributes()

			var deploymentAction, packageReference string
			if deploymentActionAttr, ok := attrs["deployment_action"].(types.String); ok && !deploymentActionAttr.IsNull() {
				deploymentAction = deploymentActionAttr.ValueString()
			}
			if packageReferenceAttr, ok := attrs["package_reference"].(types.String); ok && !packageReferenceAttr.IsNull() {
				packageReference = packageReferenceAttr.ValueString()
			}

			key := deploymentAction + "|" + packageReference
			currentActionPackagesMap[key] = packageObj
		}
	}

	flattenedActionPackages := make([]attr.Value, 0, len(actionPackages))
	for _, actionPackage := range actionPackages {
		key := actionPackage.DeploymentAction + "|" + actionPackage.PackageReference
		currentActionPackage := currentActionPackagesMap[key]
		obj := flattenChannelRuleDeploymentActionPackage(&actionPackage, currentActionPackage)
		flattenedActionPackages = append(flattenedActionPackages, obj)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getChannelRuleDeploymentActionPackageAttrTypes()}, flattenedActionPackages)
}

func flattenChannelRuleDeploymentActionPackage(actionPackage *packages.DeploymentActionPackage, currentActionPackage types.Object) types.Object {
	var currentDeploymentAction, currentPackageReference types.String
	if !currentActionPackage.IsNull() && !currentActionPackage.IsUnknown() {
		attrs := currentActionPackage.Attributes()
		if v, ok := attrs["deployment_action"]; ok {
			currentDeploymentAction = v.(types.String)
		}
		if v, ok := attrs["package_reference"]; ok {
			currentPackageReference = v.(types.String)
		}
	}

	var deploymentActionValue, packageReferenceValue types.String
	if actionPackage.DeploymentAction == "" && currentDeploymentAction.IsNull() {
		deploymentActionValue = types.StringNull()
	} else {
		deploymentActionValue = types.StringValue(actionPackage.DeploymentAction)
	}

	if actionPackage.PackageReference == "" && currentPackageReference.IsNull() {
		packageReferenceValue = types.StringNull()
	} else {
		packageReferenceValue = types.StringValue(actionPackage.PackageReference)
	}

	return types.ObjectValueMust(getChannelRuleDeploymentActionPackageAttrTypes(), map[string]attr.Value{
		"deployment_action": deploymentActionValue,
		"package_reference": packageReferenceValue,
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

func getChannelRuleDeploymentActionPackageAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"deployment_action": types.StringType,
		"package_reference": types.StringType,
	}
}
