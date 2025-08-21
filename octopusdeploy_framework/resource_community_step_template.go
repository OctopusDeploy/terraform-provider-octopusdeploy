package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
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

type communityStepTemplateTypeResource struct {
	*Config
}

var (
	_ resource.ResourceWithImportState    = &communityStepTemplateTypeResource{}
	_ resource.ResourceWithValidateConfig = &communityStepTemplateTypeResource{}
)

func NewCommunityStepTemplateResource() resource.Resource {
	return &communityStepTemplateTypeResource{}
}

func (r *communityStepTemplateTypeResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("community_step_template")
}

func (r *communityStepTemplateTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.CommunityStepTemplateSchema{}.GetResourceSchema()
}

func (r *communityStepTemplateTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (*communityStepTemplateTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *communityStepTemplateTypeResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	var data schemas.StepTemplateFromCommunityStepTemplateTypeResourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *communityStepTemplateTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.StepTemplateFromCommunityStepTemplateTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newActionTemplate, dg := mapCommunityStepTemplateResourceModelToActionTemplate(ctx, data)
	resp.Diagnostics.Append(dg...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The end users work with action templates, not community step templates.
	// The difference with a community step template is that it is installed rather than created.
	// But the end result of an installed community step template is a regular (if read only) action template in the current space.
	communityStepTemplate := actions.NewCommunityActionTemplate(newActionTemplate.Name, newActionTemplate.ActionType)
	communityStepTemplate.ID = newActionTemplate.CommunityActionTemplateID

	// Installing a community step template essentially creates a read only step template in the current space.
	communityStepTemplate, err := r.Config.Client.CommunityActionTemplates.Install(*communityStepTemplate)

	if err != nil {
		resp.Diagnostics.AddError("unable to install community step template", err.Error())
		return
	}

	// read the details of the newly installed step template
	actionTemplate, err := actiontemplates.GetByID(r.Config.Client, data.SpaceID.ValueString(), communityStepTemplate.ID)
	if err != nil {
		resp.Diagnostics.AddError("unable to read the installed community step template", err.Error())
		return
	}

	resp.Diagnostics.Append(mapCommunityStepTemplateToResourceModel(ctx, &data, actionTemplate)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *communityStepTemplateTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.StepTemplateFromCommunityStepTemplateTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Unlike the Create function, we read a community step template just like a regular action template.
	// This is because the community step template is installed in the current space and is accessed via the action templates API.
	actionTemplate, err := actiontemplates.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "action template"); err != nil {
			resp.Diagnostics.AddError("unable to load environment", err.Error())
		}
		return
	}

	resp.Diagnostics.Append(mapCommunityStepTemplateToResourceModel(ctx, &data, actionTemplate)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *communityStepTemplateTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Step templates based on community step templates are read only.
}

func (r *communityStepTemplateTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.StepTemplateFromCommunityStepTemplateTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := actiontemplates.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete step template", err.Error())
		return
	}
}

func mapCommunityStepTemplateToResourceModel(ctx context.Context, data *schemas.StepTemplateFromCommunityStepTemplateTypeResourceModel, at *actiontemplates.ActionTemplate) diag.Diagnostics {
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

	return resp
}

func mapCommunityStepTemplateResourceModelToActionTemplate(ctx context.Context, data schemas.StepTemplateFromCommunityStepTemplateTypeResourceModel) (*actiontemplates.ActionTemplate, diag.Diagnostics) {
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

	if resp.HasError() {
		return at, resp
	}
	return at, resp
}
