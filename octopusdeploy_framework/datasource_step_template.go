package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stepTemplateDataSource struct {
	*Config
}

func NewStepTemplateDataSource() datasource.DataSource {
	return &stepTemplateDataSource{}
}
func (*stepTemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("step_template")
}

func (*stepTemplateDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.StepTemplateSchema{}.GetDatasourceSchema()
}

func (d *stepTemplateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = DataSourceConfiguration(req, resp)
}

func (d *stepTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data schemas.StepTemplateTypeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	search := actiontemplates.ActionTemplateSearch{
		ID:   data.ID.ValueString(),
		Name: data.Name.ValueString(),
	}

	if (data.ID.IsNull() || data.ID.IsUnknown() || data.ID.ValueString() == "") && (data.Name.IsNull() || data.Name.IsUnknown() || data.Name.ValueString() == "") {
		resp.Diagnostics.AddError("Invalid Step Template", "Either the 'id' or 'name' attribute must be specified.")
		return
	}

	util.DatasourceReading(ctx, "step_template", search)

	actionTemplates, err := actiontemplates.Get(d.Config.Client, data.SpaceID.ValueString(), search)
	if err != nil {
		resp.Diagnostics.AddError("Unable to load step template", err.Error())
		return
	}

	actionTemplateMatch := findActionTemplateByName(actionTemplates.Items, data.Name.ValueString())

	resp.Diagnostics.Append(mapStepTemplateToDatasourceModel(&data, actionTemplateMatch)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findActionTemplateByName(actionTemplates []*actiontemplates.ActionTemplate, name string) *actiontemplates.ActionTemplate {
	// If there was no name specified, the ID must have been specified, so return the first match
	if name == "" {
		return actionTemplates[0]
	}

	// Otherwise, find the first match by name
	for _, at := range actionTemplates {
		if at.Name == name {
			return at
		}
	}

	return nil
}

func mapStepTemplateToDatasourceModel(data *schemas.StepTemplateTypeDataSourceModel, at *actiontemplates.ActionTemplate) diag.Diagnostics {
	resp := diag.Diagnostics{}

	if at != nil {
		stepTemplate, dg := convertStepTemplateAttributes(at)
		resp.Append(dg...)
		data.StepTemplate = stepTemplate
	}

	return resp
}

func convertStepTemplateAttributes(at *actiontemplates.ActionTemplate) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	params := make([]attr.Value, len(at.Parameters))
	for i, param := range at.Parameters {
		p, dg := convertStepTemplateParameterAttribute(param, nil)
		diags.Append(dg...)
		params[i] = p
	}
	paramsListValue, dg := types.ListValue(types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()}, params)
	diags.Append(dg...)

	pkgs := make([]attr.Value, len(at.Packages))
	for i, pkg := range at.Packages {
		p, dg := convertStepTemplatePackageAttribute(pkg)
		diags.Append(dg...)
		pkgs[i] = p
	}
	packageListValue, dg := types.ListValue(types.ObjectType{AttrTypes: schemas.GetStepTemplatePackageTypeAttributes()}, pkgs)
	diags.Append(dg...)

	props := make(map[string]attr.Value, len(at.Properties))
	for key, val := range at.Properties {
		props[key] = types.StringValue(val.Value)
	}
	propertiesMap, dg := types.MapValue(types.StringType, props)
	diags.Append(dg...)

	gitDepends := make([]attr.Value, len(at.GitDependencies))
	for i, gitDep := range at.GitDependencies {
		gd, dg := convertStepTemplateGitDependencyAttribute(gitDep)
		diags.Append(dg...)
		gitDepends[i] = gd
	}
	gitDependsListValue, dg := types.ListValue(types.ObjectType{AttrTypes: schemas.GetStepTemplateGitDependencyTypeAttributes()}, gitDepends)
	diags.Append(dg...)

	if diags.HasError() {
		return types.ObjectNull(schemas.GetStepTemplateParameterTypeAttributes()), diags
	}

	stepTemplate, dg := types.ObjectValue(schemas.GetStepTemplateAttributes(), map[string]attr.Value{
		"id":                           types.StringValue(at.ID),
		"name":                         types.StringValue(at.Name),
		"description":                  types.StringValue(at.Description),
		"space_id":                     types.StringValue(at.SpaceID),
		"version":                      types.Int32Value(at.Version),
		"step_package_id":              types.StringValue(at.ActionType),
		"action_type":                  types.StringValue(at.ActionType),
		"community_action_template_id": types.StringValue(at.CommunityActionTemplateID),
		"packages":                     packageListValue,
		"git_dependencies":             gitDependsListValue,
		"parameters":                   paramsListValue,
		"properties":                   propertiesMap,
	})
	diags.Append(dg...)
	return stepTemplate, diags
}

func mapCommunityStepTemplateToCommunityResourceModel(ctx context.Context, data *schemas.CommunityStepTemplateTypeResourceModel, at *actions.CommunityActionTemplate) diag.Diagnostics {
	resp := diag.Diagnostics{}

	data.ID = types.StringValue(at.ID)
	data.Name = types.StringValue(at.Name)
	data.Version = types.Int32Value(at.Version)
	data.Description = types.StringValue(at.Description)
	data.Website = types.StringValue(at.Website)
	data.HistoryUrl = types.StringValue(at.HistoryURL)
	data.Type = types.StringValue(at.ActionType)
	data.Author = types.StringValue(at.Author)

	// Parameters
	sParams, dg := convertStepTemplateToParameterAttributes(ctx, at.Parameters, data.Parameters)
	resp.Append(dg...)
	data.Parameters = sParams

	// Properties
	stringProps := make(map[string]attr.Value, len(at.Properties))
	for keys, value := range at.Properties {
		stringProps[keys] = types.StringValue(value.Value)
	}
	props, dg := types.MapValue(types.StringType, stringProps)
	resp.Append(dg...)
	data.Properties = props

	// Packages
	pkgs, dg := convertStepTemplateToPackageAttributes(at.Packages)
	resp.Append(dg...)
	data.Packages = pkgs

	return resp
}
