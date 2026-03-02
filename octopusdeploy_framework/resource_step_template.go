package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stepTemplateTypeResource struct {
	*Config
}

var (
	_ resource.ResourceWithImportState    = &stepTemplateTypeResource{}
	_ resource.ResourceWithValidateConfig = &stepTemplateTypeResource{}
)

func NewStepTemplateResource() resource.Resource {
	return &stepTemplateTypeResource{}
}

func (r *stepTemplateTypeResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("step_template")
}

func (r *stepTemplateTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.StepTemplateSchema{}.GetResourceSchema()
}

func (r *stepTemplateTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (*stepTemplateTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *stepTemplateTypeResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	var data schemas.StepTemplateTypeResourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(validateStepTemplateParameters(ctx, &data)...)
}

func (r *stepTemplateTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.StepTemplateTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newActionTemplate, dg := mapStepTemplateResourceModelToActionTemplate(ctx, data)
	resp.Diagnostics.Append(dg...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionTemplate, err := actiontemplates.Add(r.Config.Client, newActionTemplate)
	if err != nil {
		resp.Diagnostics.AddError("unable to create step template", err.Error())
		return
	}

	resp.Diagnostics.Append(mapStepTemplateToResourceModel(ctx, &data, actionTemplate)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stepTemplateTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.StepTemplateTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionTemplate, err := actiontemplates.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "action template"); err != nil {
			resp.Diagnostics.AddError("unable to load environment", err.Error())
		}
		return
	}

	resp.Diagnostics.Append(mapStepTemplateToResourceModel(ctx, &data, actionTemplate)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stepTemplateTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state schemas.StepTemplateTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	at, err := actiontemplates.GetByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load step template", err.Error())
		return
	}

	actionTemplateUpdate, dg := mapStepTemplateResourceModelToActionTemplate(ctx, data)
	resp.Diagnostics.Append(dg...)
	if resp.Diagnostics.HasError() {
		return
	}
	actionTemplateUpdate.ID = at.ID
	actionTemplateUpdate.SpaceID = at.SpaceID
	actionTemplateUpdate.Version = at.Version

	updatedActionTemplate, err := actiontemplates.Update(r.Config.Client, actionTemplateUpdate)
	if err != nil {
		resp.Diagnostics.AddError("unable to update step template", err.Error())
		return
	}

	resp.Diagnostics.Append(mapStepTemplateToResourceModel(ctx, &data, updatedActionTemplate)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stepTemplateTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.StepTemplateTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := actiontemplates.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete step template", err.Error())
		return
	}
}

func mapStepTemplateToResourceModel(ctx context.Context, data *schemas.StepTemplateTypeResourceModel, at *actiontemplates.ActionTemplate) diag.Diagnostics {
	resp := diag.Diagnostics{}

	data.ID = types.StringValue(at.ID)
	data.SpaceID = types.StringValue(at.SpaceID)
	data.Name = types.StringValue(at.Name)
	data.Version = types.Int32Value(at.Version)
	data.Description = types.StringValue(at.Description)
	data.CommunityActionTemplateId = types.StringValue(at.CommunityActionTemplateID)
	data.ActionType = types.StringValue(at.ActionType)

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

	//GitDependencies
	gitDepends, dg := convertStepTemplateToGitDependencyAttributes(at.GitDependencies)
	resp.Append(dg...)
	data.GitDependencies = gitDepends

	return resp
}

func mapStepTemplateResourceModelToActionTemplate(ctx context.Context, data schemas.StepTemplateTypeResourceModel) (*actiontemplates.ActionTemplate, diag.Diagnostics) {
	resp := diag.Diagnostics{}
	at := actiontemplates.NewActionTemplate(data.Name.ValueString(), data.ActionType.ValueString())

	at.SpaceID = data.SpaceID.ValueString()
	at.Description = data.Description.ValueString()
	if !data.CommunityActionTemplateId.IsNull() {
		at.CommunityActionTemplateID = data.CommunityActionTemplateId.ValueString()
	}

	pkgs := make([]schemas.StepTemplatePackageType, 0, len(data.Packages.Elements()))
	resp.Append(data.Packages.ElementsAs(ctx, &pkgs, false)...)
	if resp.HasError() {
		return at, resp
	}

	props := make(map[string]types.String, len(data.Properties.Elements()))
	resp.Append(data.Properties.ElementsAs(ctx, &props, false)...)
	if resp.HasError() {
		return at, resp
	}

	params := make([]schemas.StepTemplateParameterType, 0, len(data.Parameters.Elements()))
	resp.Append(data.Parameters.ElementsAs(ctx, &params, false)...)
	if resp.HasError() {
		return at, resp
	}

	gitDepends := make([]schemas.StepTemplateGitDependencyType, 0, len(data.GitDependencies.Elements()))
	resp.Append(data.GitDependencies.ElementsAs(ctx, &gitDepends, false)...)
	if resp.HasError() {
		return at, resp
	}

	if len(props) > 0 {
		templateProps := make(map[string]core.PropertyValue, len(props))
		for key, val := range props {
			templateProps[key] = core.NewPropertyValue(val.ValueString(), false)
		}
		at.Properties = templateProps
	} else {
		at.Properties = make(map[string]core.PropertyValue)
	}

	at.Packages = make([]packages.PackageReference, len(pkgs))
	if len(pkgs) > 0 {
		for i, val := range pkgs {
			pkgProps := convertAttributeStepTemplatePackageProperty(val.Properties.Attributes())
			pkgRef := packages.PackageReference{
				AcquisitionLocation: val.AcquisitionLocation.ValueString(),
				FeedID:              val.FeedID.ValueString(),
				Properties:          pkgProps,
				Name:                val.Name.ValueString(),
				PackageID:           val.PackageID.ValueString(),
			}
			pkgRef.ID = val.ID.ValueString()
			at.Packages[i] = pkgRef
		}
	}

	parameters, parameterDiags := mapStepTemplateParametersFromState(params)
	resp.Append(parameterDiags...)
	at.Parameters = parameters

	at.GitDependencies = make([]gitdependencies.GitDependency, len(gitDepends))
	if len(gitDepends) > 0 {
		for i, val := range gitDepends {
			var filePathFilters []string
			if err := val.FilePathFilters.ElementsAs(context.Background(), &filePathFilters, false); err != nil {
				return at, resp
			}
			at.GitDependencies[i] = gitdependencies.GitDependency{
				Name:              val.Name.ValueString(),
				RepositoryUri:     val.RepositoryUri.ValueString(),
				DefaultBranch:     val.DefaultBranch.ValueString(),
				GitCredentialType: val.GitCredentialType.ValueString(),
				FilePathFilters:   filePathFilters,
				GitCredentialId:   val.GitCredentialId.ValueString(),
			}
		}
	}

	if resp.HasError() {
		return at, resp
	}
	return at, resp
}

func mapStepTemplateParametersFromState(stateParameters []schemas.StepTemplateParameterType) ([]actiontemplates.ActionTemplateParameter, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	templateParameters := make([]actiontemplates.ActionTemplateParameter, len(stateParameters))
	if len(stateParameters) == 0 {
		return templateParameters, diags
	}

	paramIDMap := make(map[string]bool, len(stateParameters))
	for i, val := range stateParameters {
		templateParameter := actiontemplates.ActionTemplateParameter{
			Name:            val.Name.ValueString(),
			Label:           val.Label.ValueString(),
			HelpText:        val.HelpText.ValueString(),
			DisplaySettings: util.ConvertAttrStringMapToStringMap(val.DisplaySettings.Elements()),
		}

		// Determine default value
		defaultValue := val.DefaultValue.ValueString()
		isSensitive := false
		if !val.DefaultSensitiveValue.IsUnknown() && !val.DefaultSensitiveValue.IsNull() {
			// Is sensitive
			defaultValue = val.DefaultSensitiveValue.ValueString()
			isSensitive = true
		}
		value := core.NewPropertyValue(defaultValue, isSensitive)
		templateParameter.DefaultValue = &value

		// Confirm unique Id
		id := val.ID.ValueString()
		if _, ok := paramIDMap[id]; ok {
			diags.AddError("ID conflict", fmt.Sprintf("conflicting UUID's within parameters list: %s", id))
		}
		paramIDMap[id] = true
		templateParameter.ID = id

		templateParameters[i] = templateParameter
	}

	return templateParameters, diags
}

func convertStepTemplateToPackageAttributes(atPackage []packages.PackageReference) (types.List, diag.Diagnostics) {
	resp := diag.Diagnostics{}
	pkgs := make([]attr.Value, len(atPackage))
	for key, val := range atPackage {
		mapVal, dg := convertStepTemplatePackageAttribute(val)
		resp.Append(dg...)
		if resp.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplatePackageTypeAttributes()}), resp
		}
		pkgs[key] = mapVal
	}
	pkgSet, dg := types.ListValue(types.ObjectType{AttrTypes: schemas.GetStepTemplatePackageTypeAttributes()}, pkgs)
	resp.Append(dg...)
	if resp.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplatePackageTypeAttributes()}), resp
	}
	return pkgSet, dg
}

func convertStepTemplateToParameterAttributes(ctx context.Context, atParams []actiontemplates.ActionTemplateParameter, stateParameters types.List) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	lookup := make(map[string]*schemas.StepTemplateParameterType)
	if !(stateParameters.IsNull() || stateParameters.IsUnknown()) {
		stateElements := make([]schemas.StepTemplateParameterType, 0, len(stateParameters.Elements()))
		diags.Append(stateParameters.ElementsAs(ctx, &stateElements, false)...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()}), diags
		}

		for _, stateParameter := range stateElements {
			lookup[stateParameter.ID.ValueString()] = &stateParameter
		}
	}

	parameters := make([]attr.Value, len(atParams))
	for i, val := range atParams {
		objVal, dg := convertStepTemplateParameterAttribute(val, lookup[val.ID])
		diags.Append(dg...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()}), diags
		}
		parameters[i] = objVal
	}

	parametersList, dg := types.ListValue(types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()}, parameters)
	diags.Append(dg...)

	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()}), diags
	}

	return parametersList, diags
}

func convertStepTemplateToGitDependencyAttributes(atParams []gitdependencies.GitDependency) (types.List, diag.Diagnostics) {
	resp := diag.Diagnostics{}
	params := make([]attr.Value, len(atParams))
	for i, val := range atParams {
		objVal, dg := convertStepTemplateGitDependencyAttribute(val)
		resp.Append(dg...)
		if resp.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplateGitDependencyTypeAttributes()}), resp
		}
		params[i] = objVal
	}
	sParams, dg := types.ListValue(types.ObjectType{AttrTypes: schemas.GetStepTemplateGitDependencyTypeAttributes()}, params)
	resp.Append(dg...)
	if resp.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplateGitDependencyTypeAttributes()}), resp
	}
	return sParams, resp
}

func convertStepTemplateGitDependencyAttribute(atp gitdependencies.GitDependency) (types.Object, diag.Diagnostics) {
	filePathFilters, dg := types.ListValue(types.StringType, util.ToValueSlice(atp.FilePathFilters))
	if dg.HasError() {
		return types.ObjectNull(schemas.GetStepTemplateGitDependencyTypeAttributes()), dg
	}
	return types.ObjectValue(schemas.GetStepTemplateGitDependencyTypeAttributes(), map[string]attr.Value{
		"name":                types.StringValue(atp.Name),
		"repository_uri":      types.StringValue(atp.RepositoryUri),
		"default_branch":      types.StringValue(atp.DefaultBranch),
		"git_credential_type": types.StringValue(atp.GitCredentialType),
		"file_path_filters":   filePathFilters,
		"git_credential_id":   types.StringValue(atp.GitCredentialId),
	})
}

func convertStepTemplateParameterAttribute(atp actiontemplates.ActionTemplateParameter, stateParameter *schemas.StepTemplateParameterType) (types.Object, diag.Diagnostics) {
	displaySettings, dg := types.MapValue(types.StringType, util.ConvertStringMapToAttrStringMap(atp.DisplaySettings))

	if dg.HasError() {
		return types.ObjectNull(schemas.GetStepTemplateParameterTypeAttributes()), dg
	}

	// Set sensitive value to the value from current state.
	// This to avoid difference between planned and applied values ("unexpected new value"),
	// because Server never returns sensitive values from the API
	defaultValue := types.StringValue("")
	defaultSensitiveValue := types.StringNull()
	if atp.DefaultValue.IsSensitive && stateParameter != nil {
		defaultSensitiveValue = stateParameter.DefaultSensitiveValue
	} else {
		defaultValue = types.StringValue(atp.DefaultValue.Value)
	}

	return types.ObjectValue(schemas.GetStepTemplateParameterTypeAttributes(), map[string]attr.Value{
		"id":                      types.StringValue(atp.ID),
		"name":                    types.StringValue(atp.Name),
		"label":                   types.StringValue(atp.Label),
		"help_text":               types.StringValue(atp.HelpText),
		"default_value":           defaultValue,
		"default_sensitive_value": defaultSensitiveValue,
		"display_settings":        displaySettings,
	})
}

func convertStepTemplatePackageAttribute(atp packages.PackageReference) (types.Object, diag.Diagnostics) {
	props, dg := convertStepTemplatePackagePropertyAttribute(atp.Properties)
	if dg.HasError() {
		return types.ObjectNull(schemas.GetStepTemplatePackageTypeAttributes()), dg
	}
	return types.ObjectValue(schemas.GetStepTemplatePackageTypeAttributes(), map[string]attr.Value{
		"id":                   types.StringValue(atp.ID),
		"acquisition_location": types.StringValue(atp.AcquisitionLocation),
		"name":                 types.StringValue(atp.Name),
		"feed_id":              types.StringValue(atp.FeedID),
		"package_id":           types.StringValue(atp.PackageID),
		"properties":           props,
	})
}

func convertStepTemplatePackagePropertyAttribute(atpp map[string]string) (types.Object, diag.Diagnostics) {
	prop := make(map[string]attr.Value)
	diags := diag.Diagnostics{}

	// We need to manually convert the string map to ensure all fields are set.
	if extract, ok := atpp["Extract"]; ok {
		prop["extract"] = types.StringValue(extract)
	} else {
		diags.AddWarning("Package property missing value.", "extract value missing from package property")
		prop["extract"] = types.StringNull()
	}

	if purpose, ok := atpp["Purpose"]; ok {
		prop["purpose"] = types.StringValue(purpose)
	} else {
		diags.AddWarning("Package property missing value.", "purpose value missing from package property")
		prop["purpose"] = types.StringNull()
	}

	if purpose, ok := atpp["PackageParameterName"]; ok {
		prop["package_parameter_name"] = types.StringValue(purpose)
	} else {
		diags.AddWarning("Package property missing value.", "package_parameter_name value missing from package property")
		prop["package_parameter_name"] = types.StringNull()
	}

	if selectionMode, ok := atpp["SelectionMode"]; ok {
		prop["selection_mode"] = types.StringValue(selectionMode)
	} else {
		diags.AddWarning("Package property missing value.", "selection_mode value missing from package property")
		prop["selection_mode"] = types.StringNull()
	}

	propMap, dg := types.ObjectValue(schemas.GetStepTemplatePackagePropertiesTypeAttributes(), prop)
	if dg.HasError() {
		diags.Append(dg...)
		return types.ObjectNull(schemas.GetStepTemplatePackagePropertiesTypeAttributes()), diags
	}
	return propMap, diags
}

func convertAttributeStepTemplatePackageProperty(prop map[string]attr.Value) map[string]string {
	atpp := make(map[string]string)

	if extract, ok := prop["extract"]; ok {
		atpp["Extract"] = extract.(types.String).ValueString()
	}

	if purpose, ok := prop["purpose"]; ok {
		atpp["Purpose"] = purpose.(types.String).ValueString()
	}

	if purpose, ok := prop["package_parameter_name"]; ok {
		atpp["PackageParameterName"] = purpose.(types.String).ValueString()
	}

	if selectionMode, ok := prop["selection_mode"]; ok {
		atpp["SelectionMode"] = selectionMode.(types.String).ValueString()
	}
	return atpp
}

func validateStepTemplateParameters(ctx context.Context, data *schemas.StepTemplateTypeResourceModel) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if data.Parameters.IsUnknown() {
        return diags
    }

	if data.Parameters.IsNull() {
        return diags
    }

	parameters := make([]schemas.StepTemplateParameterType, 0, len(data.Parameters.Elements()))
	appendDiags := data.Parameters.ElementsAs(ctx, &parameters, false)
	diags.Append(appendDiags...)
	if diags.HasError() {
		return diags
	}

	for _, parameter := range parameters {
		diags.Append(validateStepTemplateParameterDefaultValue(parameter)...)
	}

	return diags
}

func validateStepTemplateParameterDefaultValue(param schemas.StepTemplateParameterType) diag.Diagnostics {
	diags := diag.Diagnostics{}

	isSensitive := false
	if !param.DisplaySettings.IsNull() && !param.DisplaySettings.IsUnknown() {
		displaySettings := param.DisplaySettings.Elements()
		if controlTypeValue, exists := displaySettings["Octopus.ControlType"]; exists {
			if ctrlTypeStr, ok := controlTypeValue.(types.String); ok {
				isSensitive = ctrlTypeStr.ValueString() == "Sensitive"
			}
		}
	}

	hasPlainValue := !param.DefaultValue.IsNull() && !param.DefaultValue.IsUnknown() && param.DefaultValue.ValueString() != ""
	hasSensitiveValue := !param.DefaultSensitiveValue.IsNull() && !param.DefaultSensitiveValue.IsUnknown()

	if isSensitive && hasPlainValue {
		diags.AddError(
			"Invalid step template parameter configuration",
			fmt.Sprintf("Parameter '%s' has display setting 'Octopus.ControlType=Sensitive' but uses 'default_value' instead of 'default_sensitive_value'. Sensitive parameters should use the 'default_sensitive_value' attribute.", param.Name.ValueString()),
		)

		return diags
	}

	if !isSensitive && hasSensitiveValue {
		diags.AddError(
			"Invalid step template parameter configuration",
			fmt.Sprintf("Parameter '%s' has non-sensitive display setting 'Octopus.ControlType', but uses 'default_sensitive_value'. Non-sensitive parameters should use the 'default_value' attribute.", param.Name.ValueString()),
		)
	}

	return diags
}
