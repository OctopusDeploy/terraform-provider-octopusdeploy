package octopusdeploy_framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tenantProjectVariableValidator struct{}

func (v tenantProjectVariableValidator) Description(ctx context.Context) string {
	return "Ensures either environment_id or scope block is provided, but not both"
}

func (v tenantProjectVariableValidator) MarkdownDescription(ctx context.Context) string {
	return "Ensures either `environment_id` or `scope` block is provided, but not both"
}

func (v tenantProjectVariableValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var environmentID types.String
	var scope types.List

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("environment_id"), &environmentID)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("scope"), &scope)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if environmentID.IsUnknown() || scope.IsUnknown() {
		return
	}

	hasEnvironmentID := !environmentID.IsNull() && environmentID.ValueString() != ""
	hasScope := !scope.IsNull() && len(scope.Elements()) > 0

	if hasEnvironmentID && hasScope {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment_id"),
			"Invalid Configuration",
			"Cannot specify both 'environment_id' and 'scope' block. Use 'environment_id' for V1 API or 'scope' block for V2 API with multiple environments.",
		)
		return
	}

	if !hasEnvironmentID && !hasScope {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Must specify either 'environment_id' or 'scope' block.",
		)
		return
	}
}

// TenantProjectVariableValidator returns a validator that ensures either environment_id or scope is provided, but not both.
func TenantProjectVariableValidator() resource.ConfigValidator {
	return tenantProjectVariableValidator{}
}
