// Package planmodifiers contains custom plan modifiers for the Terraform provider.
package planmodifiers

import (
	"context"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// OrderInsensitiveList returns a plan modifier that ignores the order of elements in a list.
// This fixes issues where APIs return scope values in a different order than submitted,
// causing unwanted Terraform plan differences.
func OrderInsensitiveList() planmodifier.List {
	return orderInsensitiveListModifier{}
}

type orderInsensitiveListModifier struct{}

func (m orderInsensitiveListModifier) Description(_ context.Context) string {
	return "Ignores order differences between planned and actual list values"
}

func (m orderInsensitiveListModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m orderInsensitiveListModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// If the plan value is null or unknown, no changes needed
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	// If the state value is null, this is a create operation, no changes needed
	if req.StateValue.IsNull() {
		return
	}

	// If the state value is unknown, no changes needed
	if req.StateValue.IsUnknown() {
		return
	}

	// Convert both lists to string slices for comparison
	planStrings := util.ExpandStringList(req.PlanValue)
	stateStrings := util.ExpandStringList(req.StateValue)

	// If the sorted content matches, keep the current state value to prevent drift
	if util.StringSlicesEqual(planStrings, stateStrings) {
		resp.PlanValue = req.StateValue
	}
}
