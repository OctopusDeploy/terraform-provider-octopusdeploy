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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &tenantProjectVariableResource{}
var _ resource.ResourceWithImportState = &tenantProjectVariableResource{}
var _ resource.ResourceWithConfigValidators = &tenantProjectVariableResource{}

type tenantProjectVariableResource struct {
	*Config
}

type tenantProjectVariableScopeModel struct {
	EnvironmentIDs types.Set `tfsdk:"environment_ids"`
}

type tenantProjectVariableResourceModel struct {
	SpaceID       types.String                      `tfsdk:"space_id"`
	TenantID      types.String                      `tfsdk:"tenant_id"`
	ProjectID     types.String                      `tfsdk:"project_id"`
	EnvironmentID types.String                      `tfsdk:"environment_id"`
	TemplateID    types.String                      `tfsdk:"template_id"`
	Value         types.String                      `tfsdk:"value"`
	Scope         []tenantProjectVariableScopeModel `tfsdk:"scope"`

	schemas.ResourceModel
}

func NewTenantProjectVariableResource() resource.Resource {
	return &tenantProjectVariableResource{}
}

func (t *tenantProjectVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.TenantProjectVariableResourceName)
}

func (t *tenantProjectVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.TenantProjectVariableSchema{}.GetResourceSchema()
}

func (t *tenantProjectVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	t.Config = ResourceConfiguration(req, resp)
}

func (t *tenantProjectVariableResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		TenantProjectVariableValidator(),
	}
}

func (t *tenantProjectVariableResource) supportsV2() bool {
	if t.Config == nil || t.Config.FeatureToggles == nil {
		// If we can't check feature toggles, the server is too old for V2
		return false
	}
	return t.Config.FeatureToggleEnabled("CommonVariableScopingFeatureToggle")
}

func (t *tenantProjectVariableResource) validateScopeSupport(planScope []tenantProjectVariableScopeModel, diags *diag.Diagnostics) bool {
	if len(planScope) > 0 && !t.supportsV2() {
		diags.AddError(
			"Scope block is not supported",
			"The 'scope' block requires V2 API support. Your Octopus Server does not support this feature.",
		)
		return false
	}
	return true
}

// findProjectVariableTemplateByProjectAndTemplate finds a project variable template and returns whether it's sensitive and its ID
func findProjectVariableTemplateByProjectAndTemplate(variables []variables.TenantProjectVariable, missingVariables []variables.TenantProjectVariable, projectID, templateID string) (isSensitive bool, variableID string, found bool) {
	for _, v := range append(variables, missingVariables...) {
		if v.ProjectID == projectID && v.TemplateID == templateID {
			return isTemplateControlTypeSensitive(v.Template.DisplaySettings), v.GetID(), true
		}
	}
	return false, "", false
}

// findProjectVariableByID finds a project variable by ID and returns whether it's sensitive
func findProjectVariableByID(variables []variables.TenantProjectVariable, id string) (isSensitive bool, found bool) {
	for _, v := range variables {
		if v.GetID() == id {
			return isTemplateControlTypeSensitive(v.Template.DisplaySettings), true
		}
	}
	return false, false
}

func projectVariableMatchesPlan(variable variables.TenantProjectVariable, planProjectID, planTemplateID string, planScope []tenantProjectVariableScopeModel, planEnvironmentID types.String) bool {
	if variable.ProjectID != planProjectID || variable.TemplateID != planTemplateID {
		return false
	}

	// Check if scope block matches
	if len(planScope) > 0 {
		return projectScopesMatch(planScope[0].EnvironmentIDs, variable.Scope.EnvironmentIds)
	}

	// Check if legacy environment_id matches
	if !planEnvironmentID.IsNull() && planEnvironmentID.ValueString() != "" {
		if len(variable.Scope.EnvironmentIds) == 1 && variable.Scope.EnvironmentIds[0] == planEnvironmentID.ValueString() {
			return true
		}
	}

	return false
}

func projectScopesMatch(planEnvIDs types.Set, serverEnvIDs []string) bool {
	if planEnvIDs.IsNull() || planEnvIDs.IsUnknown() {
		return len(serverEnvIDs) == 0
	}

	planEnvironments := make([]types.String, 0, len(planEnvIDs.Elements()))
	planEnvIDs.ElementsAs(context.Background(), &planEnvironments, false)

	if len(planEnvironments) != len(serverEnvIDs) {
		return false
	}

	planEnvSet := make(map[string]bool)
	for _, e := range planEnvironments {
		planEnvSet[e.ValueString()] = true
	}

	for _, serverEnv := range serverEnvIDs {
		if !planEnvSet[serverEnv] {
			return false
		}
	}

	return true
}

func (t *tenantProjectVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tenantProjectVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	internal.KeyedMutex.Lock(plan.TenantID.ValueString())
	defer internal.KeyedMutex.Unlock(plan.TenantID.ValueString())

	tflog.Debug(ctx, "Creating tenant project variable")

	hasEnvironmentID := !plan.EnvironmentID.IsNull() && plan.EnvironmentID.ValueString() != ""

	if !t.validateScopeSupport(plan.Scope, &resp.Diagnostics) {
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
		t.createV2(ctx, &plan, tenant, spaceID, hasEnvironmentID, resp)
	} else {
		t.createV1(ctx, &plan, tenant, hasEnvironmentID, resp)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (t *tenantProjectVariableResource) createV1(ctx context.Context, plan *tenantProjectVariableResourceModel, tenant *tenants.Tenant, hasEnvironmentID bool, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Using V1 API for tenant project variable")

	if !hasEnvironmentID {
		resp.Diagnostics.AddError("Invalid configuration", "environment_id is required for V1 API")
		return
	}

	id := fmt.Sprintf("%s:%s:%s:%s", plan.TenantID.ValueString(), plan.ProjectID.ValueString(), plan.EnvironmentID.ValueString(), plan.TemplateID.ValueString())

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	isSensitive, err := checkIfVariableIsSensitive(tenantVariables, *plan)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if err := updateTenantProjectVariable(tenantVariables, *plan, isSensitive); err != nil {
		resp.Diagnostics.AddError("Error updating tenant project variable", err.Error())
		return
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}

	plan.ID = types.StringValue(id)
	plan.SpaceID = types.StringValue(tenant.SpaceID)

	tflog.Debug(ctx, "Tenant project variable created with V1 API", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})
}

func (t *tenantProjectVariableResource) createV2(ctx context.Context, plan *tenantProjectVariableResourceModel, tenant *tenants.Tenant, spaceID string, hasEnvironmentID bool, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Using V2 API for tenant project variable")

	query := variables.GetTenantProjectVariablesQuery{
		TenantID:                plan.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: true,
	}

	getResp, err := tenants.GetProjectVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant project variables", err.Error())
		return
	}

	isSensitive, _, found := findProjectVariableTemplateByProjectAndTemplate(
		getResp.Variables,
		getResp.MissingVariables,
		plan.ProjectID.ValueString(),
		plan.TemplateID.ValueString(),
	)

	if !found {
		resp.Diagnostics.AddError("Template not found", fmt.Sprintf("Template %s not found in project %s", plan.TemplateID.ValueString(), plan.ProjectID.ValueString()))
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
	} else if hasEnvironmentID {
		scope.EnvironmentIds = []string{plan.EnvironmentID.ValueString()}
	}

	payloads := []variables.TenantProjectVariablePayload{}

	for _, v := range getResp.Variables {
		payloads = append(payloads, variables.TenantProjectVariablePayload{
			ID:         v.GetID(),
			ProjectID:  v.ProjectID,
			TemplateID: v.TemplateID,
			Value:      v.Value,
			Scope:      v.Scope,
		})
	}

	newPayload := variables.TenantProjectVariablePayload{
		ProjectID:  plan.ProjectID.ValueString(),
		TemplateID: plan.TemplateID.ValueString(),
		Value:      core.NewPropertyValue(plan.Value.ValueString(), isSensitive),
		Scope:      scope,
	}
	payloads = append(payloads, newPayload)

	cmd := &variables.ModifyTenantProjectVariablesCommand{
		Variables: payloads,
	}

	updateResp, err := tenants.UpdateProjectVariables(t.Client, spaceID, plan.TenantID.ValueString(), cmd)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant project variables", err.Error())
		return
	}

	var createdID string
	for _, v := range updateResp.Variables {
		if projectVariableMatchesPlan(v, plan.ProjectID.ValueString(), plan.TemplateID.ValueString(), plan.Scope, plan.EnvironmentID) {
			createdID = v.GetID()
			break
		}
	}

	if createdID == "" {
		resp.Diagnostics.AddError("Failed to get variable ID", "Variable was created but ID not returned in response")
		return
	}

	plan.ID = types.StringValue(createdID)
	plan.SpaceID = types.StringValue(tenant.SpaceID)

	tflog.Debug(ctx, "Tenant project variable created with V2 API", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})
}

func (t *tenantProjectVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tenantProjectVariableResourceModel
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

func (t *tenantProjectVariableResource) readV1(ctx context.Context, state *tenantProjectVariableResourceModel, tenant *tenants.Tenant, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading tenant project variable with V1 API")

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	if !checkIfTemplateExists(tenantVariables, *state) {
		// The template no longer exists, so the variable can no longer exist either
		return
	}

	isSensitive, err := checkIfVariableIsSensitive(tenantVariables, *state)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if projectVariable, ok := tenantVariables.ProjectVariables[state.ProjectID.ValueString()]; ok {
		if environment, ok := projectVariable.Variables[state.EnvironmentID.ValueString()]; ok {
			if value, ok := environment[state.TemplateID.ValueString()]; ok {
				if !isSensitive {
					state.Value = types.StringValue(value.Value)
				}
			} else {
				resp.State.RemoveResource(ctx)
				return
			}
		}
	}
}

func (t *tenantProjectVariableResource) readV2(ctx context.Context, state *tenantProjectVariableResourceModel, spaceID string, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading tenant project variable with V2 API")

	query := variables.GetTenantProjectVariablesQuery{
		TenantID:                state.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: false,
	}

	getResp, err := tenants.GetProjectVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant project variables", err.Error())
		return
	}

	var found bool
	for _, v := range getResp.Variables {
		if v.GetID() == state.ID.ValueString() {
			if len(v.Scope.EnvironmentIds) > 0 {
				envSet := util.BuildStringSetOrEmpty(v.Scope.EnvironmentIds)
				state.Scope = []tenantProjectVariableScopeModel{{EnvironmentIDs: envSet}}
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

func (t *tenantProjectVariableResource) migrateV1ToV2OnRead(ctx context.Context, state *tenantProjectVariableResourceModel, spaceID string, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Migrating tenant project variable from V1 to V2")

	query := variables.GetTenantProjectVariablesQuery{
		TenantID:                state.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: false,
	}

	getResp, err := tenants.GetProjectVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant project variables", err.Error())
		return
	}

	var found bool
	for _, v := range getResp.Variables {
		if v.ProjectID == state.ProjectID.ValueString() && v.TemplateID == state.TemplateID.ValueString() {
			if len(v.Scope.EnvironmentIds) > 0 {
				for _, envID := range v.Scope.EnvironmentIds {
					if envID == state.EnvironmentID.ValueString() {
						// Migrate to V2 ID
						state.ID = types.StringValue(v.GetID())

						envSet := util.BuildStringSetOrEmpty(v.Scope.EnvironmentIds)
						state.Scope = []tenantProjectVariableScopeModel{{EnvironmentIDs: envSet}}
						state.EnvironmentID = types.StringNull()

						isSensitive := isTemplateControlTypeSensitive(v.Template.DisplaySettings)
						if !isSensitive {
							state.Value = types.StringValue(v.Value.Value)
						}

						found = true
						break
					}
				}
			} else {
				// Variable exists but has no scope (applies to all environments)
				state.ID = types.StringValue(v.GetID())
				state.Scope = nil
				state.EnvironmentID = types.StringNull()

				isSensitive := isTemplateControlTypeSensitive(v.Template.DisplaySettings)
				if !isSensitive {
					state.Value = types.StringValue(v.Value.Value)
				}

				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
}

func (t *tenantProjectVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tenantProjectVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	internal.KeyedMutex.Lock(plan.TenantID.ValueString())
	defer internal.KeyedMutex.Unlock(plan.TenantID.ValueString())

	hasEnvironmentID := !plan.EnvironmentID.IsNull() && plan.EnvironmentID.ValueString() != ""

	if !t.validateScopeSupport(plan.Scope, &resp.Diagnostics) {
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
		t.migrateV1ToV2OnUpdate(ctx, &plan, spaceID, hasEnvironmentID, resp)
	} else if !isV1ID {
		t.updateV2(ctx, &plan, spaceID, hasEnvironmentID, resp)
	} else {
		t.updateV1(ctx, &plan, tenant, resp)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (t *tenantProjectVariableResource) updateV1(ctx context.Context, plan *tenantProjectVariableResourceModel, tenant *tenants.Tenant, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating tenant project variable with V1 API")

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	isSensitive, err := checkIfVariableIsSensitive(tenantVariables, *plan)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if err := updateTenantProjectVariable(tenantVariables, *plan, isSensitive); err != nil {
		resp.Diagnostics.AddError("Error updating tenant project variable", err.Error())
		return
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}
}

func (t *tenantProjectVariableResource) updateV2(ctx context.Context, plan *tenantProjectVariableResourceModel, spaceID string, hasEnvironmentID bool, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating tenant project variable with V2 API")

	query := variables.GetTenantProjectVariablesQuery{
		TenantID:                plan.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: false,
	}

	getResp, err := tenants.GetProjectVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant project variables", err.Error())
		return
	}

	isSensitive, found := findProjectVariableByID(getResp.Variables, plan.ID.ValueString())

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
	} else if hasEnvironmentID {
		scope.EnvironmentIds = []string{plan.EnvironmentID.ValueString()}
	}

	payloads := []variables.TenantProjectVariablePayload{}

	for _, v := range getResp.Variables {
		if v.GetID() == plan.ID.ValueString() {
			payloads = append(payloads, variables.TenantProjectVariablePayload{
				ID:         plan.ID.ValueString(),
				ProjectID:  plan.ProjectID.ValueString(),
				TemplateID: plan.TemplateID.ValueString(),
				Value:      core.NewPropertyValue(plan.Value.ValueString(), isSensitive),
				Scope:      scope,
			})
		} else {
			payloads = append(payloads, variables.TenantProjectVariablePayload{
				ID:         v.GetID(),
				ProjectID:  v.ProjectID,
				TemplateID: v.TemplateID,
				Value:      v.Value,
				Scope:      v.Scope,
			})
		}
	}

	cmd := &variables.ModifyTenantProjectVariablesCommand{
		Variables: payloads,
	}

	_, err = tenants.UpdateProjectVariables(t.Client, spaceID, plan.TenantID.ValueString(), cmd)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant project variables", err.Error())
		return
	}
}

func (t *tenantProjectVariableResource) migrateV1ToV2OnUpdate(ctx context.Context, plan *tenantProjectVariableResourceModel, spaceID string, hasEnvironmentID bool, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Migrating tenant project variable from V1 to V2 during update")

	query := variables.GetTenantProjectVariablesQuery{
		TenantID:                plan.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: true,
	}

	getResp, err := tenants.GetProjectVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant project variables", err.Error())
		return
	}

	isSensitive, existingID, foundExisting := findProjectVariableTemplateByProjectAndTemplate(
		getResp.Variables,
		getResp.MissingVariables,
		plan.ProjectID.ValueString(),
		plan.TemplateID.ValueString(),
	)

	scope := variables.TenantVariableScope{}
	if len(plan.Scope) > 0 {
		envIDs, diags := util.SetToStringArray(ctx, plan.Scope[0].EnvironmentIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		scope.EnvironmentIds = envIDs
	} else if hasEnvironmentID {
		scope.EnvironmentIds = []string{plan.EnvironmentID.ValueString()}
	}

	payloads := []variables.TenantProjectVariablePayload{}

	for _, v := range getResp.Variables {
		if v.ProjectID == plan.ProjectID.ValueString() && v.TemplateID == plan.TemplateID.ValueString() {
			payloads = append(payloads, variables.TenantProjectVariablePayload{
				ID:         existingID,
				ProjectID:  plan.ProjectID.ValueString(),
				TemplateID: plan.TemplateID.ValueString(),
				Value:      core.NewPropertyValue(plan.Value.ValueString(), isSensitive),
				Scope:      scope,
			})
		} else {
			payloads = append(payloads, variables.TenantProjectVariablePayload{
				ID:         v.GetID(),
				ProjectID:  v.ProjectID,
				TemplateID: v.TemplateID,
				Value:      v.Value,
				Scope:      v.Scope,
			})
		}
	}

	if !foundExisting {
		payloads = append(payloads, variables.TenantProjectVariablePayload{
			ProjectID:  plan.ProjectID.ValueString(),
			TemplateID: plan.TemplateID.ValueString(),
			Value:      core.NewPropertyValue(plan.Value.ValueString(), isSensitive),
			Scope:      scope,
		})
	}

	cmd := &variables.ModifyTenantProjectVariablesCommand{
		Variables: payloads,
	}

	updateResp, err := tenants.UpdateProjectVariables(t.Client, spaceID, plan.TenantID.ValueString(), cmd)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant project variables", err.Error())
		return
	}

	for _, v := range updateResp.Variables {
		if v.ProjectID == plan.ProjectID.ValueString() && v.TemplateID == plan.TemplateID.ValueString() {
			plan.ID = types.StringValue(v.GetID())
			break
		}
	}

	tflog.Debug(ctx, "Tenant project variable migrated to V2", map[string]interface{}{
		"new_id": plan.ID.ValueString(),
	})
}

func (t *tenantProjectVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tenantProjectVariableResourceModel
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

func (t *tenantProjectVariableResource) deleteV1(ctx context.Context, state *tenantProjectVariableResourceModel, tenant *tenants.Tenant, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting tenant project variable with V1 API")

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	isSensitive, err := checkIfVariableIsSensitive(tenantVariables, *state)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if projectVariable, ok := tenantVariables.ProjectVariables[state.ProjectID.ValueString()]; ok {
		if environment, ok := projectVariable.Variables[state.EnvironmentID.ValueString()]; ok {
			if isSensitive {
				environment[state.TemplateID.ValueString()] = core.PropertyValue{IsSensitive: true, SensitiveValue: &core.SensitiveValue{HasValue: false}}
			} else {
				delete(environment, state.TemplateID.ValueString())
			}
		}
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}
}

func (t *tenantProjectVariableResource) deleteV2(ctx context.Context, state *tenantProjectVariableResourceModel, spaceID string, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting tenant project variable with V2 API")

	query := variables.GetTenantProjectVariablesQuery{
		TenantID:                state.TenantID.ValueString(),
		SpaceID:                 spaceID,
		IncludeMissingVariables: false,
	}

	getResp, err := tenants.GetProjectVariables(t.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant project variables", err.Error())
		return
	}

	isSensitive, found := findProjectVariableByID(getResp.Variables, state.ID.ValueString())

	if !found {
		return
	}

	payloads := []variables.TenantProjectVariablePayload{}

	for _, v := range getResp.Variables {
		if v.GetID() == state.ID.ValueString() {
			if isSensitive {
				payloads = append(payloads, variables.TenantProjectVariablePayload{
					ID:         state.ID.ValueString(),
					ProjectID:  state.ProjectID.ValueString(),
					TemplateID: state.TemplateID.ValueString(),
					Value:      core.PropertyValue{IsSensitive: true, SensitiveValue: &core.SensitiveValue{HasValue: false}},
					Scope:      variables.TenantVariableScope{},
				})
			}
		} else {
			payloads = append(payloads, variables.TenantProjectVariablePayload{
				ID:         v.GetID(),
				ProjectID:  v.ProjectID,
				TemplateID: v.TemplateID,
				Value:      v.Value,
				Scope:      v.Scope,
			})
		}
	}

	cmd := &variables.ModifyTenantProjectVariablesCommand{
		Variables: payloads,
	}

	_, err = tenants.UpdateProjectVariables(t.Client, spaceID, state.TenantID.ValueString(), cmd)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting tenant project variable", err.Error())
		return
	}
}

func (t *tenantProjectVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	// V1 format: TenantID:ProjectID:EnvironmentID:TemplateID
	if len(idParts) == 4 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), idParts[2])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("template_id"), idParts[3])...)
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

		query := variables.GetTenantProjectVariablesQuery{
			TenantID:                tenantID,
			SpaceID:                 tenant.SpaceID,
			IncludeMissingVariables: false,
		}

		getResp, err := tenants.GetProjectVariables(t.Client, query)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving tenant project variables", err.Error())
			return
		}

		var found bool
		for _, v := range getResp.Variables {
			if v.GetID() == variableID {
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), variableID)...)
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), tenantID)...)
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), v.ProjectID)...)
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("template_id"), v.TemplateID)...)
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space_id"), tenant.SpaceID)...)

				if len(v.Scope.EnvironmentIds) > 0 {
					envSet := util.BuildStringSetOrEmpty(v.Scope.EnvironmentIds)
					scopeModel := []tenantProjectVariableScopeModel{{EnvironmentIDs: envSet}}
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
			"  V1: TenantID:ProjectID:EnvironmentID:TemplateID (e.g. Tenants-123:Projects-456:Environments-789:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c)\n"+
			"  V2: TenantID:VariableID (e.g. Tenants-123:TenantVariables-456)",
	)
}

func checkIfTemplateExists(tenantVariables *variables.TenantVariables, plan tenantProjectVariableResourceModel) bool {
	if projectVariable, ok := tenantVariables.ProjectVariables[plan.ProjectID.ValueString()]; ok {
		for _, template := range projectVariable.Templates {
			if template.GetID() == plan.TemplateID.ValueString() {
				return true
			}
		}
	}
	return false
}

func checkIfVariableIsSensitive(tenantVariables *variables.TenantVariables, plan tenantProjectVariableResourceModel) (bool, error) {
	if projectVariable, ok := tenantVariables.ProjectVariables[plan.ProjectID.ValueString()]; ok {
		for _, template := range projectVariable.Templates {
			if template.GetID() == plan.TemplateID.ValueString() {
				return isTemplateControlTypeSensitive(template.DisplaySettings), nil
			}
		}
	}
	return false, fmt.Errorf("unable to find template for tenant variable")
}

func updateTenantProjectVariable(tenantVariables *variables.TenantVariables, plan tenantProjectVariableResourceModel, isSensitive bool) error {
	if projectVariable, ok := tenantVariables.ProjectVariables[plan.ProjectID.ValueString()]; ok {
		if environment, ok := projectVariable.Variables[plan.EnvironmentID.ValueString()]; ok {
			environment[plan.TemplateID.ValueString()] = core.NewPropertyValue(plan.Value.ValueString(), isSensitive)
			return nil
		}
	}
	return fmt.Errorf("unable to locate tenant variable for tenant ID %s", plan.TenantID.ValueString())
}
