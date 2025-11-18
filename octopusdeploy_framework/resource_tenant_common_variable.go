package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &tenantCommonVariableResource{}
var _ resource.ResourceWithImportState = &tenantCommonVariableResource{}

type tenantCommonVariableResource struct {
	*Config
}

type tenantCommonVariableScopeModel struct {
	EnvironmentIDs types.Set `tfsdk:"environment_ids"`
}

type tenantCommonVariableResourceModel struct {
	SpaceID              types.String                     `tfsdk:"space_id"`
	TenantID             types.String                     `tfsdk:"tenant_id"`
	LibraryVariableSetID types.String                     `tfsdk:"library_variable_set_id"`
	TemplateID           types.String                     `tfsdk:"template_id"`
	Value                types.String                     `tfsdk:"value"`
	Scope                []tenantCommonVariableScopeModel `tfsdk:"scope"`

	schemas.ResourceModel
}

func NewTenantCommonVariableResource() resource.Resource {
	return &tenantCommonVariableResource{}
}

func (t *tenantCommonVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.TenantCommonVariableResourceName)
}

func (t *tenantCommonVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetTenantCommonVariableResourceSchema()
}

func (t *tenantCommonVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	t.Config = ResourceConfiguration(req, resp)
}

func (t *tenantCommonVariableResource) supportsV2() bool {
	if t.Config == nil || t.Config.FeatureToggles == nil {
		// If we can't check feature toggles, the server is too old for V2
		return false
	}
	return t.Config.FeatureToggleEnabled("CommonVariableScopingFeatureToggle")
}

func isCompositeID(id string) bool {
	return strings.Contains(id, ":")
}

func findCommonVariableTemplateByLibrarySetAndTemplate(variables []variables.TenantCommonVariable, missingVariables []variables.TenantCommonVariable, libraryVariableSetID, templateID string) (isSensitive bool, found bool) {
	for _, v := range append(variables, missingVariables...) {
		if v.LibraryVariableSetId == libraryVariableSetID && v.TemplateID == templateID {
			return isTemplateControlTypeSensitive(v.Template.DisplaySettings), true
		}
	}
	return false, false
}

func findCommonVariableByID(variables []variables.TenantCommonVariable, id string) (isSensitive bool, found bool) {
	for _, v := range variables {
		if v.GetID() == id {
			return isTemplateControlTypeSensitive(v.Template.DisplaySettings), true
		}
	}
	return false, false
}

func isTemplateControlTypeSensitive(displaySettings map[string]string) bool {
	return displaySettings["Octopus.ControlType"] == "Sensitive"
}

func (t *tenantCommonVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tenantCommonVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	internal.KeyedMutex.Lock(plan.TenantID.ValueString())
	defer internal.KeyedMutex.Unlock(plan.TenantID.ValueString())

	if len(plan.Scope) > 0 && !t.supportsV2() {
		resp.Diagnostics.AddError(
			"Scoped tenant variables are not supported",
			"The scope block requires V2 API support (CommonVariableScopingFeatureToggle). Your Octopus Server does not support this feature.",
		)
		return
	}

	tenant, err := tenants.GetByID(t.Client, plan.SpaceID.ValueString(), plan.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant", err.Error())
		return
	}

	spaceID := plan.SpaceID.ValueString()
	if spaceID == "" {
		spaceID = tenant.SpaceID
	}

	if t.supportsV2() {
		t.createV2(ctx, &plan, tenant, spaceID, resp)
	} else {
		t.createV1(ctx, &plan, tenant, resp)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (t *tenantCommonVariableResource) createV1(ctx context.Context, plan *tenantCommonVariableResourceModel, tenant *tenants.Tenant, resp *resource.CreateResponse) {
	id := fmt.Sprintf("%s:%s:%s", plan.TenantID.ValueString(), plan.LibraryVariableSetID.ValueString(), plan.TemplateID.ValueString())

	tenant, err := tenants.GetByID(t.Client, plan.SpaceID.ValueString(), plan.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant", err.Error())
		return
	}

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	err = checkIfCandidateVariableRequiredForTenant(tenant, tenantVariables, *plan)
	if err != nil {
		resp.Diagnostics.AddError("Tenant doesn't need a value for this Common Variable", "Tenants must be connected to a Project with an included Library Variable Set that defines Common Variable templates, before common variable values can be provided ("+err.Error()+")")
		return
	}

	isSensitive, err := checkIfCommonVariableIsSensitive(tenantVariables, *plan)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if err := updateTenantCommonVariable(tenantVariables, *plan, isSensitive); err != nil {
		resp.Diagnostics.AddError("Error updating tenant common variable", err.Error())
		return
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}

	plan.ID = types.StringValue(id)
	plan.SpaceID = types.StringValue(tenant.SpaceID)

	tflog.Debug(ctx, "Tenant common variable created with V1 API", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})
}

func (t *tenantCommonVariableResource) createV2(ctx context.Context, plan *tenantCommonVariableResourceModel, tenant *tenants.Tenant, spaceID string, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Using V2 API for tenant common variable")

	query := variables.GetTenantCommonVariablesQuery{
		TenantID:                plan.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: true,
	}

	getResp, err := tenants.GetCommonVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant common variables", err.Error())
		return
	}

	isSensitive, found := findCommonVariableTemplateByLibrarySetAndTemplate(
		getResp.Variables,
		getResp.MissingVariables,
		plan.LibraryVariableSetID.ValueString(),
		plan.TemplateID.ValueString(),
	)

	if !found {
		resp.Diagnostics.AddError("Template not found", fmt.Sprintf("Template %s not found in library variable set %s", plan.TemplateID.ValueString(), plan.LibraryVariableSetID.ValueString()))
		return
	}

	scope := variables.TenantVariableScope{}
	if len(plan.Scope) > 0 {
		envIDs, diags := util.SetToStringArray(ctx, plan.Scope[0].EnvironmentIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		scope.EnvironmentIds = envIDs
	}

	payloads := []variables.TenantCommonVariablePayload{}

	for _, v := range getResp.Variables {
		payloads = append(payloads, variables.TenantCommonVariablePayload{
			ID:                   v.GetID(),
			LibraryVariableSetId: v.LibraryVariableSetId,
			TemplateID:           v.TemplateID,
			Value:                v.Value,
			Scope:                v.Scope,
		})
	}

	newPayload := variables.TenantCommonVariablePayload{
		LibraryVariableSetId: plan.LibraryVariableSetID.ValueString(),
		TemplateID:           plan.TemplateID.ValueString(),
		Value:                core.NewPropertyValue(plan.Value.ValueString(), isSensitive),
		Scope:                scope,
	}
	payloads = append(payloads, newPayload)

	cmd := &variables.ModifyTenantCommonVariablesCommand{
		Variables: payloads,
	}

	updateResp, err := tenants.UpdateCommonVariables(t.Client, spaceID, plan.TenantID.ValueString(), cmd)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant common variables", err.Error())
		return
	}

	// Find the created variable and get its ID
	var createdID string
	for _, v := range updateResp.Variables {
		if v.LibraryVariableSetId == plan.LibraryVariableSetID.ValueString() && v.TemplateID == plan.TemplateID.ValueString() {
			// Also match on scope to handle multiple variables with same template but different scopes
			if len(plan.Scope) > 0 {
				if len(v.Scope.EnvironmentIds) > 0 {
					planEnvs := plan.Scope[0].EnvironmentIDs.Elements()
					if len(planEnvs) == len(v.Scope.EnvironmentIds) {
						match := true
						planEnvSet := make(map[string]bool)
						for _, e := range planEnvs {
							planEnvSet[e.String()] = true
						}
						for _, serverEnv := range v.Scope.EnvironmentIds {
							if !planEnvSet["\""+serverEnv+"\""] {
								match = false
								break
							}
						}
						if match {
							createdID = v.GetID()
							break
						}
					}
				}
			} else {
				// No scope in plan, match unscoped variable
				if len(v.Scope.EnvironmentIds) == 0 {
					createdID = v.GetID()
					break
				}
			}
		}
	}

	if createdID == "" {
		resp.Diagnostics.AddError("Failed to get variable ID", "Variable was created but ID not returned in response")
		return
	}

	plan.ID = types.StringValue(createdID)
	plan.SpaceID = types.StringValue(tenant.SpaceID)

	tflog.Debug(ctx, "Tenant common variable created with V2 API", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})
}

func (t *tenantCommonVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tenantCommonVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	internal.KeyedMutex.Lock(state.TenantID.ValueString())
	defer internal.KeyedMutex.Unlock(state.TenantID.ValueString())

	tenant, err := tenants.GetByID(t.Client, state.SpaceID.ValueString(), state.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant", err.Error())
		return
	}

	spaceID := state.SpaceID.ValueString()
	if spaceID == "" {
		spaceID = tenant.SpaceID
	}

	// Determine which API version to use based on ID format
	isV1ID := isCompositeID(state.ID.ValueString())

	if isV1ID && t.supportsV2() {
		t.migrateV1ToV2OnRead(ctx, &state, spaceID, resp)
	} else if !isV1ID {
		t.readV2(ctx, &state, spaceID, resp)
	} else {
		t.readV1(ctx, &state, tenant, resp)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (t *tenantCommonVariableResource) readV1(ctx context.Context, state *tenantCommonVariableResourceModel, tenant *tenants.Tenant, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading tenant common variable with V1 API")

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	isSensitive, err := checkIfCommonVariableIsSensitive(tenantVariables, *state)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if libraryVariable, ok := tenantVariables.LibraryVariables[state.LibraryVariableSetID.ValueString()]; ok {
		if value, ok := libraryVariable.Variables[state.TemplateID.ValueString()]; ok {
			if !isSensitive {
				state.Value = types.StringValue(value.Value)
			}
		} else {
			resp.State.RemoveResource(ctx)
			return
		}
	}
}

func (t *tenantCommonVariableResource) readV2(ctx context.Context, state *tenantCommonVariableResourceModel, spaceID string, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading tenant common variable with V2 API")

	query := variables.GetTenantCommonVariablesQuery{
		TenantID:                state.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: false,
	}

	getResp, err := tenants.GetCommonVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant common variables", err.Error())
		return
	}

	var found bool
	for _, v := range getResp.Variables {
		if v.GetID() == state.ID.ValueString() {
			if len(v.Scope.EnvironmentIds) > 0 {
				envSet := util.BuildStringSetOrEmpty(v.Scope.EnvironmentIds)
				state.Scope = []tenantCommonVariableScopeModel{{EnvironmentIDs: envSet}}
			} else {
				state.Scope = nil
			}

			isSensitive := isTemplateControlTypeSensitive(v.Template.DisplaySettings)
			if !isSensitive {
				state.Value = types.StringValue(v.Value.Value)
			}

			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
}

func (t *tenantCommonVariableResource) migrateV1ToV2OnRead(ctx context.Context, state *tenantCommonVariableResourceModel, spaceID string, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Migrating tenant common variable from V1 to V2")

	query := variables.GetTenantCommonVariablesQuery{
		TenantID:                state.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: false,
	}

	getResp, err := tenants.GetCommonVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant common variables", err.Error())
		return
	}

	// Find the variable
	var found bool
	for _, v := range getResp.Variables {
		if v.LibraryVariableSetId == state.LibraryVariableSetID.ValueString() && v.TemplateID == state.TemplateID.ValueString() {
			// Migrate to V2 ID
			state.ID = types.StringValue(v.GetID())

			// Update scope
			if len(v.Scope.EnvironmentIds) > 0 {
				envSet := util.BuildStringSetOrEmpty(v.Scope.EnvironmentIds)
				state.Scope = []tenantCommonVariableScopeModel{{EnvironmentIDs: envSet}}
			} else {
				state.Scope = nil
			}

			// Update value if not sensitive
			isSensitive := isTemplateControlTypeSensitive(v.Template.DisplaySettings)
			if !isSensitive {
				state.Value = types.StringValue(v.Value.Value)
			}

			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
}

func (t *tenantCommonVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tenantCommonVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	internal.KeyedMutex.Lock(plan.TenantID.ValueString())
	defer internal.KeyedMutex.Unlock(plan.TenantID.ValueString())

	// Validate scope block usage on unsupported servers
	if len(plan.Scope) > 0 && !t.supportsV2() {
		resp.Diagnostics.AddError(
			"V2 API not supported",
			"The 'scope' block requires V2 API support (CommonVariableScopingFeatureToggle). Your Octopus Server version does not support this feature. Please upgrade your server or remove the scope block.",
		)
		return
	}

	tenant, err := tenants.GetByID(t.Client, plan.SpaceID.ValueString(), plan.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant", err.Error())
		return
	}

	spaceID := plan.SpaceID.ValueString()
	if spaceID == "" {
		spaceID = tenant.SpaceID
	}

	isV1ID := isCompositeID(plan.ID.ValueString())
	if isV1ID && t.supportsV2() {
		t.migrateV1ToV2OnUpdate(ctx, &plan, spaceID, resp)
	} else if !isV1ID {
		t.updateV2(ctx, &plan, spaceID, resp)
	} else {
		t.updateV1(ctx, &plan, tenant, resp)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (t *tenantCommonVariableResource) updateV2(ctx context.Context, plan *tenantCommonVariableResourceModel, spaceID string, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating tenant common variable with V2 API")

	query := variables.GetTenantCommonVariablesQuery{
		TenantID:                plan.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: false,
	}

	getResp, err := tenants.GetCommonVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant common variables", err.Error())
		return
	}

	isSensitive, found := findCommonVariableByID(getResp.Variables, plan.ID.ValueString())

	if !found {
		resp.Diagnostics.AddError("Variable not found", fmt.Sprintf("Variable with ID %s not found", plan.ID.ValueString()))
		return
	}

	scope := variables.TenantVariableScope{}
	if len(plan.Scope) > 0 {
		envIDs, diags := util.SetToStringArray(ctx, plan.Scope[0].EnvironmentIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		scope.EnvironmentIds = envIDs
	}

	payloads := []variables.TenantCommonVariablePayload{}

	for _, v := range getResp.Variables {
		if v.GetID() == plan.ID.ValueString() {
			payloads = append(payloads, variables.TenantCommonVariablePayload{
				ID:                   plan.ID.ValueString(),
				LibraryVariableSetId: plan.LibraryVariableSetID.ValueString(),
				TemplateID:           plan.TemplateID.ValueString(),
				Value:                core.NewPropertyValue(plan.Value.ValueString(), isSensitive),
				Scope:                scope,
			})
		} else {
			payloads = append(payloads, variables.TenantCommonVariablePayload{
				ID:                   v.GetID(),
				LibraryVariableSetId: v.LibraryVariableSetId,
				TemplateID:           v.TemplateID,
				Value:                v.Value,
				Scope:                v.Scope,
			})
		}
	}

	cmd := &variables.ModifyTenantCommonVariablesCommand{
		Variables: payloads,
	}

	_, err = tenants.UpdateCommonVariables(t.Client, spaceID, plan.TenantID.ValueString(), cmd)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant common variables", err.Error())
		return
	}
}

func (t *tenantCommonVariableResource) updateV1(ctx context.Context, plan *tenantCommonVariableResourceModel, tenant *tenants.Tenant, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating tenant common variable with V1 API")

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	isSensitive, err := checkIfCommonVariableIsSensitive(tenantVariables, *plan)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if err := updateTenantCommonVariable(tenantVariables, *plan, isSensitive); err != nil {
		resp.Diagnostics.AddError("Error updating tenant common variable", err.Error())
		return
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}
}

func (t *tenantCommonVariableResource) migrateV1ToV2OnUpdate(ctx context.Context, plan *tenantCommonVariableResourceModel, spaceID string, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Migrating tenant common variable from V1 to V2 during update")

	query := variables.GetTenantCommonVariablesQuery{
		TenantID:                plan.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: true,
	}

	getResp, err := tenants.GetCommonVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant common variables", err.Error())
		return
	}

	var isSensitive bool
	var existingID string
	var foundExisting bool
	for _, v := range append(getResp.Variables, getResp.MissingVariables...) {
		if v.LibraryVariableSetId == plan.LibraryVariableSetID.ValueString() && v.TemplateID == plan.TemplateID.ValueString() {
			isSensitive = isTemplateControlTypeSensitive(v.Template.DisplaySettings)
			existingID = v.GetID()
			foundExisting = true
			break
		}
	}

	scope := variables.TenantVariableScope{}
	if len(plan.Scope) > 0 {
		envIDs, diags := util.SetToStringArray(ctx, plan.Scope[0].EnvironmentIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		scope.EnvironmentIds = envIDs
	}

	payloads := []variables.TenantCommonVariablePayload{}

	for _, v := range getResp.Variables {
		if v.LibraryVariableSetId == plan.LibraryVariableSetID.ValueString() && v.TemplateID == plan.TemplateID.ValueString() {
			payloads = append(payloads, variables.TenantCommonVariablePayload{
				ID:                   existingID,
				LibraryVariableSetId: plan.LibraryVariableSetID.ValueString(),
				TemplateID:           plan.TemplateID.ValueString(),
				Value:                core.NewPropertyValue(plan.Value.ValueString(), isSensitive),
				Scope:                scope,
			})
		} else {
			payloads = append(payloads, variables.TenantCommonVariablePayload{
				ID:                   v.GetID(),
				LibraryVariableSetId: v.LibraryVariableSetId,
				TemplateID:           v.TemplateID,
				Value:                v.Value,
				Scope:                v.Scope,
			})
		}
	}

	if !foundExisting {
		payloads = append(payloads, variables.TenantCommonVariablePayload{
			LibraryVariableSetId: plan.LibraryVariableSetID.ValueString(),
			TemplateID:           plan.TemplateID.ValueString(),
			Value:                core.NewPropertyValue(plan.Value.ValueString(), isSensitive),
			Scope:                scope,
		})
	}

	cmd := &variables.ModifyTenantCommonVariablesCommand{
		Variables: payloads,
	}

	updateResp, err := tenants.UpdateCommonVariables(t.Client, spaceID, plan.TenantID.ValueString(), cmd)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant common variables", err.Error())
		return
	}

	for _, v := range updateResp.Variables {
		if v.LibraryVariableSetId == plan.LibraryVariableSetID.ValueString() && v.TemplateID == plan.TemplateID.ValueString() {
			plan.ID = types.StringValue(v.GetID())
			break
		}
	}

	tflog.Debug(ctx, "Tenant common variable migrated to V2", map[string]interface{}{
		"new_id": plan.ID.ValueString(),
	})
}

func (t *tenantCommonVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tenantCommonVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	internal.KeyedMutex.Lock(state.TenantID.ValueString())
	defer internal.KeyedMutex.Unlock(state.TenantID.ValueString())

	tenant, err := tenants.GetByID(t.Client, state.SpaceID.ValueString(), state.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant", err.Error())
		return
	}

	spaceID := state.SpaceID.ValueString()
	if spaceID == "" {
		spaceID = tenant.SpaceID
	}

	isV1ID := isCompositeID(state.ID.ValueString())
	if !isV1ID {
		t.deleteV2(ctx, &state, spaceID, resp)
	} else {
		t.deleteV1(ctx, &state, tenant, resp)
	}
}

func (t *tenantCommonVariableResource) deleteV1(ctx context.Context, state *tenantCommonVariableResourceModel, tenant *tenants.Tenant, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting tenant common variable with V1 API")

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	isSensitive, err := checkIfCommonVariableIsSensitive(tenantVariables, *state)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if libraryVariable, ok := tenantVariables.LibraryVariables[state.LibraryVariableSetID.ValueString()]; ok {
		if isSensitive {
			libraryVariable.Variables[state.TemplateID.ValueString()] = core.PropertyValue{IsSensitive: true, SensitiveValue: &core.SensitiveValue{HasValue: false}}
		} else {
			delete(libraryVariable.Variables, state.TemplateID.ValueString())
		}
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}
}

func (t *tenantCommonVariableResource) deleteV2(ctx context.Context, state *tenantCommonVariableResourceModel, spaceID string, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting tenant common variable with V2 API")

	query := variables.GetTenantCommonVariablesQuery{
		TenantID:                state.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: false,
	}

	getResp, err := tenants.GetCommonVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant common variables", err.Error())
		return
	}

	isSensitive, found := findCommonVariableByID(getResp.Variables, state.ID.ValueString())

	if !found {
		return
	}

	payloads := []variables.TenantCommonVariablePayload{}

	for _, v := range getResp.Variables {
		if v.GetID() == state.ID.ValueString() {
			if isSensitive {
				payloads = append(payloads, variables.TenantCommonVariablePayload{
					ID:                   state.ID.ValueString(),
					LibraryVariableSetId: state.LibraryVariableSetID.ValueString(),
					TemplateID:           state.TemplateID.ValueString(),
					Value:                core.PropertyValue{IsSensitive: true, SensitiveValue: &core.SensitiveValue{HasValue: false}},
					Scope:                variables.TenantVariableScope{},
				})
			}
		} else {
			payloads = append(payloads, variables.TenantCommonVariablePayload{
				ID:                   v.GetID(),
				LibraryVariableSetId: v.LibraryVariableSetId,
				TemplateID:           v.TemplateID,
				Value:                v.Value,
				Scope:                v.Scope,
			})
		}
	}

	cmd := &variables.ModifyTenantCommonVariablesCommand{
		Variables: payloads,
	}

	_, err = tenants.UpdateCommonVariables(t.Client, spaceID, state.TenantID.ValueString(), cmd)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting tenant common variable", err.Error())
		return
	}
}

func (t *tenantCommonVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	// V1 format: TenantID:LibraryVariableSetID:TemplateID
	if len(idParts) == 3 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("library_variable_set_id"), idParts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("template_id"), idParts[2])...)
		return
	}

	// V2 format: TenantID:VariableID
	if len(idParts) == 2 {
		tenantID := idParts[0]
		variableID := idParts[1]

		tenant, err := tenants.GetByID(t.Client, "", tenantID)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving tenant", err.Error())
			return
		}

		query := variables.GetTenantCommonVariablesQuery{
			TenantID:                tenantID,
			SpaceID:                 tenant.SpaceID,
			IncludeMissingVariables: false,
		}

		getResp, err := tenants.GetCommonVariables(t.Client, query)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving tenant common variables", err.Error())
			return
		}

		var found bool
		for _, v := range getResp.Variables {
			if v.GetID() == variableID {
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), variableID)...)
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), tenantID)...)
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("library_variable_set_id"), v.LibraryVariableSetId)...)
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("template_id"), v.TemplateID)...)
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), tenant.SpaceID)...)

				if len(v.Scope.EnvironmentIds) > 0 {
					envSet := util.BuildStringSetOrEmpty(v.Scope.EnvironmentIds)
					scopeModel := []tenantCommonVariableScopeModel{{EnvironmentIDs: envSet}}
					resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope"), scopeModel)...)
				}

				isSensitive := isTemplateControlTypeSensitive(v.Template.DisplaySettings)
				if !isSensitive {
					resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("value"), v.Value.Value)...)
				}

				found = true
				break
			}
		}

		if !found {
			resp.Diagnostics.AddError(
				"Variable not found",
				fmt.Sprintf("Variable with ID %s not found for tenant %s", variableID, tenantID),
			)
		}
		return
	}

	resp.Diagnostics.AddError(
		"Incorrect Import Format",
		"ID must be in one of these formats:\n"+
			"  V1: TenantID:LibraryVariableSetID:TemplateID (e.g. Tenants-123:LibraryVariableSets-456:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c)\n"+
			"  V2: TenantID:VariableID (e.g. Tenants-123:TenantVariables-456)",
	)
}

func checkIfCommonVariableIsSensitive(tenantVariables *variables.TenantVariables, plan tenantCommonVariableResourceModel) (bool, error) {
	if libraryVariable, ok := tenantVariables.LibraryVariables[plan.LibraryVariableSetID.ValueString()]; ok {
		for _, template := range libraryVariable.Templates {
			if template.GetID() == plan.TemplateID.ValueString() {
				return isTemplateControlTypeSensitive(template.DisplaySettings), nil
			}
		}
	}
	return false, fmt.Errorf("unable to find template for tenant variable")
}

func checkIfCandidateVariableRequiredForTenant(tenant *tenants.Tenant, tenantVariables *variables.TenantVariables, plan tenantCommonVariableResourceModel) error {
	if len(tenant.ProjectEnvironments) == 0 {
		return fmt.Errorf("tenant not connected to any projects")
	}

	if libraryVariable, ok := tenantVariables.LibraryVariables[plan.LibraryVariableSetID.ValueString()]; ok {
		for _, template := range libraryVariable.Templates {
			if template.GetID() == plan.TemplateID.ValueString() {
				return nil
			}
		}
	} else {
		return fmt.Errorf("tenant not connected to a project that includes variable set %s", plan.LibraryVariableSetID)
	}

	return fmt.Errorf("common template %s not found in variable set %s", plan.TemplateID, plan.LibraryVariableSetID)
}

func updateTenantCommonVariable(tenantVariables *variables.TenantVariables, plan tenantCommonVariableResourceModel, isSensitive bool) error {
	if libraryVariable, ok := tenantVariables.LibraryVariables[plan.LibraryVariableSetID.ValueString()]; ok {
		libraryVariable.Variables[plan.TemplateID.ValueString()] = core.NewPropertyValue(plan.Value.ValueString(), isSensitive)
		return nil
	}
	return fmt.Errorf("unable to locate tenant variable for tenant ID %s", plan.TenantID.ValueString())
}
